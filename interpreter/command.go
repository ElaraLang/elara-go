package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/ast"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/util"
	_ "github.com/ElaraLang/elara/util"
	"reflect"
	"strings"
)

type Command interface {
	Exec(ctx *Context) *ReturnedValue
}

type DefineVarCommand struct {
	Name        string
	Mutable     bool
	Type        ast.Type
	value       Command
	runtimeType Type

	hashedName uint64
}

func (c *DefineVarCommand) getType(ctx *Context) Type {
	if c.runtimeType == nil {
		if c.Type == nil {
			return nil
		}
		c.runtimeType = FromASTType(c.Type, ctx)
	}
	return c.runtimeType
}

func (c *DefineVarCommand) Exec(ctx *Context) *ReturnedValue {
	if c.hashedName == 0 {
		c.hashedName = util.Hash(c.Name)
	}
	var value *Value
	foundVar, _ := ctx.FindVariableMaxDepth(c.hashedName, 1)
	if foundVar != nil {
		asFunction, isFunction := foundVar.Value.Value.(*Function)
		if isFunction {
			value = c.value.Exec(ctx).Unwrap()
			valueAsFunction, valueIsFunction := value.Value.(*Function)
			if valueIsFunction && !asFunction.Signature.Accepts(&valueAsFunction.Signature, ctx, false) {
				//We'll allow it because the functions have different arity
			} else {
				//panic("Variable named " + c.Name + " already exists with the current signature") //TODO this might need to come back, maybe.
			}
		} else {
			panic("Variable named " + c.Name + " already exists")
		}
	}
	if value == nil {
		value = c.value.Exec(ctx).Unwrap()
	}

	if value == nil {
		panic("Command " + reflect.TypeOf(c.value).String() + " returned nil")
	}

	variableType := c.getType(ctx)
	if variableType != nil {
		if !variableType.Accepts(value.Type, ctx) {
			panic("Cannot use value of type " + value.Type.Name() + " in place of " + variableType.Name() + " for variable " + c.Name)
		}
	} else {
		variableType = value.Type
	}
	variable := &Variable{
		Name:    c.Name,
		Mutable: c.Mutable,
		Type:    variableType,
		Value:   value,
	}

	ctx.DefineVariable(variable)
	return NilValue()
}

type AssignmentCommand struct {
	Name  string
	value Command

	hashedName uint64
}

func (c *AssignmentCommand) Exec(ctx *Context) *ReturnedValue {
	if c.hashedName == 0 {
		c.hashedName = util.Hash(c.Name)
	}
	variable := ctx.FindVariable(c.hashedName)
	if variable == nil {
		panic("No such variable " + c.Name)
	}

	if !variable.Mutable {
		panic("Cannot reassign immutable variable " + c.Name)
	}

	value := c.value.Exec(ctx).Unwrap()

	if !variable.Type.Accepts(value.Type, ctx) {
		panic("Cannot reassign variable " + c.Name + " of type " + variable.Type.Name() + " to value " + value.String() + " of type " + value.Type.Name())
	}

	variable.Value = value
	return NonReturningValue(value)
}

type VariableCommand struct {
	Variable string

	hash      uint64
	cachedVar *Value
}

func (c *VariableCommand) findVariable(ctx *Context) *Variable {
	if c.hash == 0 {
		c.hash = util.Hash(c.Variable)
	}

	variable := ctx.FindVariable(c.hash)
	return variable
}

func (c *VariableCommand) Exec(ctx *Context) *ReturnedValue {
	if c.cachedVar != nil {
		return NonReturningValue(c.cachedVar)
	}
	paramIndex := -1
	fun := ctx.function
	if fun != nil {
		for i, parameter := range fun.Signature.Parameters {
			if parameter.Name == c.Variable {
				paramIndex = i
				break
			}
		}
	}
	if paramIndex != -1 {
		param := ctx.FindParameter(uint(paramIndex))
		if param != nil {
			return NonReturningValue(param)
		}
	}
	variable := c.findVariable(ctx)
	if variable != nil {
		return NonReturningValue(variable.Value)
	}

	constructor := ctx.FindConstructor(c.Variable)
	if constructor == nil {
		if ctx.function != nil && ctx.function.context != nil {
			return c.Exec(ctx.function.context)
		}
		panic("No such variable or parameter or constructor " + c.Variable)
	}
	c.cachedVar = constructor
	return NonReturningValue(constructor)
}

type InvocationCommand struct {
	Invoking Command
	args     []Command

	cachedFun *Function
}

func (c *InvocationCommand) findReceiverFunction(ctx *Context, receiver *Value, argValues []*Value, functionName string, nameHash uint64) *Function {
	receiverType := receiver.Type
	parameters := []Parameter{{
		Name: "this",
		Type: receiverType,
	}}
	for i, value := range argValues {
		parameters = append(parameters, Parameter{
			Name: fmt.Sprintf("<param%d>", i),
			Type: value.Type,
		})
	}
	receiverSignature := &Signature{
		Parameters: parameters,
		ReturnType: AnyType, //can't infer this rn
	}

	receiverFunction := ctx.FindFunction(nameHash, receiverSignature)
	if receiverFunction == nil {
		paramTypes := make([]string, 0)
		for _, value := range argValues {
			paramTypes = append(paramTypes, value.Type.Name())
		}
		panic("Unknown function " + receiverType.Name() + "::" + functionName + "(" + strings.Join(paramTypes, ",") + ")")
	}

	return receiverFunction
}
func (c *InvocationCommand) Exec(ctx *Context) *ReturnedValue {
	context, usingReceiver := c.Invoking.(*ContextCommand)

	argValues := make([]*Value, len(c.args))
	for i, arg := range c.args {
		argValues[i] = arg.Exec(ctx).UnwrapNotNil()
	}

	if !usingReceiver {
		if c.cachedFun != nil {
			return NonReturningValue(c.cachedFun.Exec(ctx, argValues)) //Avoid unnecessary lookup
		}
		val := c.Invoking.Exec(ctx).Unwrap()
		fun, ok := val.Value.(*Function)
		if !ok {
			panic("Cannot invoke value that isn't a function ")
		}
		switch t := c.Invoking.(type) {
		case *VariableCommand:
			variable := t.findVariable(ctx)
			if variable != nil && !variable.Mutable {
				c.cachedFun = fun
			}
		}

		return NonReturningValue(fun.Exec(ctx, argValues))
	}

	//ContextCommand seems to think it's a special case... because it is.
	var receiver *Value

	functionName := context.variable

	receiver = context.receiver.Exec(ctx).Unwrap()

	if c.cachedFun != nil {
		argValuesAndSelf := []*Value{receiver}
		argValuesAndSelf = append(argValuesAndSelf, argValues...)
		return NonReturningValue(c.cachedFun.Exec(ctx, argValuesAndSelf))
	}

	structType, isStruct := receiver.Type.(*StructType)

	if isStruct {
		value, ok := structType.GetProperty(functionName)
		if ok {
			function, ok := value.DefaultValue.Value.(*Function)
			if !ok {
				panic("Cannot invoke non-function " + value.Name)
			}
			this := context.receiver.Exec(ctx).Unwrap()

			argValues := make([]*Value, len(c.args))
			argValues = append(argValues, this)
			for i, arg := range c.args {
				argValues[i] = arg.Exec(ctx).Unwrap()
			}
			return NonReturningValue(function.Exec(ctx, argValues))
		}
	}

	extension := ctx.FindExtension(receiver.Type, functionName)
	if extension != nil {
		fun := extension.Value.Value.Value.(*Function)
		c.cachedFun = fun
		argValuesAndSelf := []*Value{receiver}
		argValuesAndSelf = append(argValuesAndSelf, argValues...)
		return NonReturningValue(fun.Exec(ctx, argValuesAndSelf))
	}

	//Look for a receiver
	receiverFunction := c.findReceiverFunction(ctx, receiver, argValues, functionName, context.hash())
	argValuesAndSelf := []*Value{receiver}
	argValuesAndSelf = append(argValuesAndSelf, argValues...)
	return NonReturningValue(receiverFunction.Exec(ctx, argValuesAndSelf))
}

type AbstractCommand struct {
	content func(ctx *Context) *ReturnedValue
}

func (c *AbstractCommand) Exec(ctx *Context) *ReturnedValue {
	return c.content(ctx)
}

func NewAbstractCommand(content func(ctx *Context) *ReturnedValue) *AbstractCommand {
	return &AbstractCommand{
		content: content,
	}
}

type LiteralCommand struct {
	value *Value
}

func (c *LiteralCommand) Exec(_ *Context) *ReturnedValue {
	return NonReturningValue(c.value)
}

type FunctionLiteralCommand struct {
	name       *string
	parameters []ast.Parameter
	returnType ast.Type //Can be nil - infer return type
	body       Command

	currentContext *Context
}

func (c *FunctionLiteralCommand) Exec(ctx *Context) *ReturnedValue {
	if c.currentContext == nil {
		c.currentContext = ctx.Clone()
		//Function literals take a snapshot of their current context to avoid scoping issues
		//This one will be cached forever, so we don't need to cleanup
	}
	params := make([]Parameter, len(c.parameters))

	for i, parameter := range c.parameters {
		paramType := FromASTType(parameter.Type, c.currentContext)
		params[i] = Parameter{
			Type:     paramType,
			Name:     parameter.Identifier.Name,
			Position: uint(i),
		}
	}

	astReturnType := c.returnType
	var returnType Type
	if astReturnType == nil {
		returnType = AnyType
	} else {
		returnType = FromASTType(c.returnType, c.currentContext)
	}

	fun := &Function{
		name: c.name,
		Signature: Signature{
			Parameters: params,
			ReturnType: returnType,
		},
		Body:    c.body,
		context: c.currentContext,
	}

	functionType := NewFunctionType(fun)

	return NonReturningValue(&Value{
		Type:  functionType,
		Value: fun,
	})
}

type BinaryOperatorCommand struct {
	lhs Command
	op  func(ctx *Context, lhs *Value, rhs *Value) *ReturnedValue
	rhs Command
}

func (c *BinaryOperatorCommand) Exec(ctx *Context) *ReturnedValue {
	lhs := c.lhs.Exec(ctx).Unwrap()
	rhs := c.rhs.Exec(ctx).Unwrap()

	return c.op(ctx, lhs, rhs)
}

type BlockCommand struct {
	lines []*Command
}

func (c *BlockCommand) Exec(ctx *Context) *ReturnedValue {
	var last = NonReturningValue(UnitValue())
	for _, lineRef := range c.lines {
		line := *lineRef
		val := line.Exec(ctx)
		if val.IsReturning {
			return val
		}
		last = val
	}
	return last
}

type ContextCommand struct {
	receiver       Command
	variable       string
	hashedVariable uint64
}

func (c *ContextCommand) hash() uint64 {
	if c.hashedVariable == 0 {
		c.hashedVariable = util.Hash(c.variable)
	}
	return c.hashedVariable
}
func (c *ContextCommand) Exec(ctx *Context) *ReturnedValue {
	if c.hashedVariable == 0 {
		c.hashedVariable = util.Hash(c.variable)
	}
	receiver := c.receiver.Exec(ctx).Unwrap()

	var value *ReturnedValue
	switch val := receiver.Value.(type) {
	case *Collection:
		switch c.variable {
		case "size":
			value = NonReturningValue(IntValue(int64(len(val.Elements))))
		}
	case *Map:
		switch c.variable {
		case "keys":
			keySet := make([]*Value, len(val.Elements))
			for i, element := range val.Elements {
				keySet[i] = element.Key
			}
			collection := &Collection{
				ElementType: val.MapType.KeyType,
				Elements:    keySet,
			}
			collectionType := NewCollectionType(collection)

			value = NonReturningValue(&Value{
				Type:  collectionType,
				Value: collection,
			})
		case "values":
			valueSet := make([]*Value, len(val.Elements))
			for i, element := range val.Elements {
				valueSet[i] = element.Value
			}
			collection := &Collection{
				ElementType: val.MapType.ValueType,
				Elements:    valueSet,
			}
			collectionType := NewCollectionType(collection)

			value = NonReturningValue(&Value{
				Type:  collectionType,
				Value: collection,
			})
		}
	case *Instance:
		{
			value = NonReturningValue(val.Values[c.variable])
		}
	default:
		panic("Unsupported receiver " + util.Stringify(receiver))
	}
	if value != nil && value.Value != nil {
		return value
	}

	//Search for an extension
	extension := ctx.FindExtension(receiver.Type, c.variable)
	if extension == nil {
		panic("Unknown property or extension for " + receiver.String() + " with name " + c.variable)
	}
	return NonReturningValue(extension.Value.Value)
}

type IfElseCommand struct {
	condition  Command
	ifBranch   Command
	elseBranch Command
}

func (c *IfElseCommand) Exec(ctx *Context) *ReturnedValue {
	condition := c.condition.Exec(ctx)
	value, ok := condition.Unwrap().Value.(bool)
	if !ok {
		panic("If statements requires boolean value")
	}

	if value {
		return c.ifBranch.Exec(ctx)
	} else if c.elseBranch != nil {
		return c.elseBranch.Exec(ctx)
	} else {
		return NilValue()
	}
}

type IfElseExpressionCommand struct {
	condition  Command
	ifBranch   Command
	elseBranch Command
}

func (c *IfElseExpressionCommand) Exec(ctx *Context) *ReturnedValue {
	condition := c.condition.Exec(ctx)
	value, ok := condition.Unwrap().Value.(bool)
	if !ok {
		panic("If statements requires boolean value")
	}

	if value {
		if c.ifBranch != nil {
			return c.ifBranch.Exec(ctx)
		}
	} else {
		if c.elseBranch != nil {
			return c.elseBranch.Exec(ctx)
		}
	}
	return NilValue()
}

type ReturnCommand struct {
	returning Command
}

func (c *ReturnCommand) Exec(ctx *Context) *ReturnedValue {
	if c.returning == nil {
		return ReturningValue(UnitValue())
	}
	return ReturningValue(c.returning.Exec(ctx).Unwrap())
}

type NamespaceCommand struct {
	namespace string
}

func (c *NamespaceCommand) Exec(ctx *Context) *ReturnedValue {
	ctx.Init(c.namespace)
	return NilValue()
}

type ImportCommand struct {
	module string
}

func (c *ImportCommand) Exec(ctx *Context) *ReturnedValue {
	ctx.Import(c.module)
	return NilValue()
}

type StructDefCommand struct {
	name   string
	fields []ast.StructField
}

func (c *StructDefCommand) Exec(ctx *Context) *ReturnedValue {

	properties := make([]Property, len(c.fields))
	propertyPositions := map[string]int{}

	for i, field := range c.fields {
		var Type Type
		if field.Type == nil {
			Type = AnyType
		} else {
			Type = FromASTType(field.Type, ctx)
		}

		var defaultValue *Value
		//if field != nil {
		//	defaultValue = ExpressionToCommand(field.Default).Exec(ctx).Unwrap()
		//}

		modifiers := uint(0)
		//if field.Mutable {
		//	modifiers |= Mut
		//}
		properties[i] = Property{
			Name:         field.Identifier.Name,
			Modifiers:    modifiers,
			Type:         Type,
			DefaultValue: defaultValue,
		}
		propertyPositions[field.Identifier.Name] = i
	}

	ctx.types[c.name] = &StructType{
		TypeName:          c.name,
		Properties:        properties,
		propertyPositions: propertyPositions,
	}

	return NilValue()
}

type ExtendCommand struct {
	Type       string
	statements []Command
	alias      string
}

func (c *ExtendCommand) Exec(ctx *Context) *ReturnedValue {
	extending := ctx.FindType(c.Type)
	if extending == nil {
		panic("No such type " + c.Type)
	}

	for _, statement := range c.statements {
		switch statement := statement.(type) {
		case *DefineVarCommand:

			value := statement.value.Exec(ctx).Unwrap()

			asFunction, isFunction := value.Value.(*Function)
			if isFunction {
				//Need to append the receiver to the function's signature
				receiverParameter := Parameter{
					Name:     c.alias,
					Position: 0,
					Type:     extending,
				}
				signature := asFunction.Signature
				params := make([]Parameter, len(signature.Parameters)+1)
				params[0] = receiverParameter
				for i := range signature.Parameters {
					p := signature.Parameters[i]
					p.Position++
					params[i] = p
				}
				signature.Parameters = params
				asFunction.Signature = signature
			}

			variable := &Variable{
				Name:    statement.Name,
				Mutable: statement.Mutable,
				Type:    statement.getType(ctx),
				Value:   value,
			}

			extension := &Extension{
				ReceiverName: c.alias,
				Value:        variable,
			}

			ctx.DefineExtension(extending, statement.Name, extension)
		}
	}
	return NilValue()
}

type TypeCheckCommand struct {
	expression Command
	checkType  ast.Type
}

func (c *TypeCheckCommand) Exec(ctx *Context) *ReturnedValue {
	checkAgainst := FromASTType(c.checkType, ctx)
	if checkAgainst == nil {
		panic("No such type " + util.Stringify(c.checkType))
	}
	res := c.expression.Exec(ctx).Unwrap()
	is := checkAgainst.Accepts(res.Type, ctx)
	return NonReturningValue(BooleanValue(is))
}

type WhileCommand struct {
	condition Command
	body      Command
}

func (c *WhileCommand) Exec(ctx *Context) *ReturnedValue {
	for {
		val := c.condition.Exec(ctx)
		condition, ok := val.Unwrap().Value.(bool)
		if !ok {
			panic("If statements requires boolean condition")
		}
		if !condition {
			break
		}
		c.body.Exec(ctx)
	}
	return NilValue()
}

type CollectionCommand struct {
	Elements []Command
}

func (c *CollectionCommand) Exec(ctx *Context) *ReturnedValue {
	elements := make([]*Value, len(c.Elements))
	for i, element := range c.Elements {
		elements[i] = element.Exec(ctx).Unwrap()
	}
	//In a proper type system we might try and find a union of all elements, but this is dynamic, and I'm lazy
	var collType Type
	if len(elements) == 0 {
		collType = AnyType
	} else {
		collType = elements[0].Type
	}
	collection := &Collection{
		ElementType: collType,
		Elements:    elements,
	}
	collectionType := NewCollectionType(collection)

	return NonReturningValue(&Value{
		Type:  collectionType,
		Value: collection,
	})
}

type AccessCommand struct {
	checking Command
	index    Command
}

func (c *AccessCommand) Exec(ctx *Context) *ReturnedValue {
	checking := c.checking.Exec(ctx).Unwrap().Value
	switch accessingType := checking.(type) {
	case *Collection:
		index, isInt := c.index.Exec(ctx).Unwrap().Value.(int64)
		if !isInt {
			panic("Index was not an integer")
		}
		return NonReturningValue(accessingType.Elements[index])

	case *Map:
		index := c.index.Exec(ctx).Unwrap()
		return NonReturningValue(accessingType.Get(ctx, index))
	}
	panic("Indexed access not supported for non-collection type")
}

type TypeCommand struct {
	name  string
	value ast.Type
}

func (c *TypeCommand) Exec(ctx *Context) *ReturnedValue {
	runtimeType := FromASTType(c.value, ctx)
	existing := ctx.FindType(c.name)
	if existing != nil {
		panic("TokenType with name " + c.name + " already exists in current scope")
	}
	ctx.types[c.name] = runtimeType
	return NilValue()
}

type MapCommand struct {
	entries []MapEntry
}
type MapEntry struct {
	key   Command
	value Command
}

func (c *MapCommand) Exec(ctx *Context) *ReturnedValue {
	elements := make([]*Entry, 0)
	for _, entry := range c.entries {
		key := entry.key.Exec(ctx).Unwrap()
		value := entry.value.Exec(ctx).Unwrap()
		entry := &Entry{
			Key:   key,
			Value: value,
		}
		elements = append(elements, entry)
	}

	mapValue := MapOf(elements)
	mapType := mapValue.MapType
	value := NewValue(mapType, mapValue)
	return NonReturningValue(value)
}

func ToCommand(statement ast.Statement) Command {
	switch t := statement.(type) {
	case *ast.DeclarationStatement:
		valueExpr := NamedExpressionToCommand(t.Value, &t.Identifier.Name)
		return &DefineVarCommand{
			Name:    t.Identifier.Name,
			Mutable: t.Mutable,
			Type:    t.Type,
			value:   valueExpr,
		}

	case *ast.ExpressionStatement:
		return ExpressionToCommand(t.Expression)

	case *ast.BlockStatement:
		commands := make([]*Command, len(t.Block))
		for i, stmt := range t.Block {
			cmd := ToCommand(stmt)
			commands[i] = &cmd
			_, isReturn := cmd.(*ReturnCommand)
			//Small optimisation, it's not worth transforming anything that won't ever be reached
			if isReturn {
				return &BlockCommand{lines: commands[:i+1]}
			}
		}

		return &BlockCommand{lines: commands}

	//case *ast.Condi:
	//	condition := ExpressionToCommand(t.Condition)
	//	ifBranch := ToCommand(t.MainBranch)
	//	var elseBranch Command
	//	if t.ElseBranch != nil {
	//		elseBranch = ToCommand(t.ElseBranch)
	//	}
	//
	//	return &IfElseCommand{
	//		condition:  condition,
	//		ifBranch:   ifBranch,
	//		elseBranch: elseBranch,
	//	}

	case *ast.ReturnStatement:
		if t.Value != nil {
			return &ReturnCommand{
				ExpressionToCommand(t.Value),
			}
		}
		return &ReturnCommand{
			nil,
		}
	case *ast.NamespaceStatement:
		return &NamespaceCommand{
			namespace: t.Module.Pkg, //TODO full module support
		}
	case *ast.ImportStatement:
		return &ImportCommand{
			module: t.Module.Pkg,
		}
	case *ast.StructDefStatement:
		name := t.Id.Name
		return &StructDefCommand{
			name:   name,
			fields: t.Fields,
		}

	case *ast.ExtendStatement:
		commands := make([]Command, len(t.Body.Block))
		for i, stmt := range t.Body.Block {
			commands[i] = ToCommand(stmt)
		}
		return &ExtendCommand{
			Type:       t.Identifier.Name,
			statements: commands,
			alias:      t.Alias.Name,
		}

	case *ast.WhileStatement:
		return &WhileCommand{
			condition: ExpressionToCommand(t.Condition),
			body:      ToCommand(t.Body),
		}
	case *ast.TypeStatement:
		return &TypeCommand{
			name:  t.Identifier.Name,
			value: t.Contract,
		}
	}

	panic("Could not handle " + reflect.TypeOf(statement).Name())
}

func ExpressionToCommand(expr ast.Expression) Command {
	return NamedExpressionToCommand(expr, nil)
}

func NamedExpressionToCommand(expr ast.Expression, name *string) Command {

	switch t := expr.(type) {
	case *ast.IdentifierLiteral:
		return &VariableCommand{Variable: t.Name}

	case *ast.CallExpression:
		fun := ExpressionToCommand(t.Expression)
		args := make([]Command, 0)
		for _, arg := range t.Arguments {
			command := ExpressionToCommand(arg)
			if command == nil {
				panic("Could not convert expression " + reflect.TypeOf(arg).Name() + " to condition")
			}
			args = append(args, command)
		}

		return &InvocationCommand{
			Invoking: fun,
			args:     args,
		}

	case *ast.StringLiteral:
		str := t.Value
		value := StringValue(str)
		return &LiteralCommand{value: value}

	case *ast.IntegerLiteral:
		integer := t.Value
		value := IntValue(integer)
		return &LiteralCommand{value: value}
	case *ast.FloatLiteral:
		float := t.Value
		value := FloatValue(float)
		return &LiteralCommand{value: value}
	case *ast.BooleanLiteral:
		boolean := t.Value
		value := BooleanValue(boolean)
		return &LiteralCommand{value: value}
	case *ast.CharLiteral:
		char := t.Value
		value := CharValue(char)
		return &LiteralCommand{value: value}
	case *ast.BinaryExpression:
		lhs := t.Left
		lhsCmd := ExpressionToCommand(lhs)
		op := t.Operator.TokenType
		rhs := t.Right
		rhsCmd := ExpressionToCommand(rhs)

		switch op {
		case lexer.Add:
			return &InvocationCommand{
				Invoking: &ContextCommand{receiver: lhsCmd, variable: "plus"},
				args:     []Command{rhsCmd},
			}
		case lexer.Subtract:
			return &InvocationCommand{
				Invoking: &ContextCommand{receiver: lhsCmd, variable: "minus"},
				args:     []Command{rhsCmd},
			}
		case lexer.Multiply:
			return &InvocationCommand{
				Invoking: &ContextCommand{receiver: lhsCmd, variable: "times"},
				args:     []Command{rhsCmd},
			}
		case lexer.Slash:
			return &InvocationCommand{
				Invoking: &ContextCommand{receiver: lhsCmd, variable: "divide"},
				args:     []Command{rhsCmd},
			}
		case lexer.Equals:
			return &InvocationCommand{Invoking: &ContextCommand{receiver: lhsCmd, variable: "equals"},
				args: []Command{rhsCmd},
			}
		case lexer.NotEquals:
			return NewAbstractCommand(func(ctx *Context) *ReturnedValue {
				command := &InvocationCommand{Invoking: &ContextCommand{receiver: lhsCmd, variable: "equals"},
					args: []Command{rhsCmd},
				}

				val := command.Exec(ctx).Unwrap()

				asBool, ok := (*val).Value.(bool)
				if !ok {
					panic("equals function did not return bool")
				}
				return NonReturningValue(BooleanValue(!asBool))
			})

		case lexer.Mod:
			return &InvocationCommand{Invoking: &ContextCommand{receiver: lhsCmd, variable: "mod"},
				args: []Command{rhsCmd},
			}
		}
	case *ast.FunctionLiteral:
		return &FunctionLiteralCommand{
			name:       name,
			parameters: t.Parameters,
			returnType: t.ReturnType,
			body:       ToCommand(t.Body),
		}

	case *ast.PropertyExpression:
		contextCmd := ExpressionToCommand(t.Context)
		varName := t.Variable.Name
		return &ContextCommand{
			contextCmd,
			varName,
			util.Hash(varName),
		}

	case *ast.AssignmentExpression:
		//TODO contexts
		name := t.Variable.Name
		valueCmd := NamedExpressionToCommand(t.Value, &name)
		return &AssignmentCommand{
			Name:  name,
			value: valueCmd,
		}

	case *ast.IfExpression:
		condition := ExpressionToCommand(t.Condition)
		var ifBranch Command
		if t.MainBranch != nil {
			ifBranch = ToCommand(t.MainBranch)
		}

		var elseBranch Command
		if t.ElseBranch != nil {
			elseBranch = ToCommand(t.ElseBranch)
		}

		return &IfElseExpressionCommand{
			condition:  condition,
			ifBranch:   ifBranch,
			elseBranch: elseBranch,
		}

	case *ast.TypeOperationExpression:
		switch t.Operation.TokenType {
		case lexer.Is:
			return &TypeCheckCommand{
				expression: ExpressionToCommand(t.Expression),
				checkType:  t.Type,
			}
		case lexer.As:
			return &TypeCheckCommand{ // TODO:: Add the type cast command
				expression: ExpressionToCommand(t.Expression),
				checkType:  t.Type,
			}

		}

	case *ast.CollectionLiteral:
		elements := make([]Command, len(t.Elements))
		for i, element := range t.Elements {
			elements[i] = ExpressionToCommand(element)
		}
		return &CollectionCommand{Elements: elements}

	case *ast.AccessExpression:
		return &AccessCommand{
			checking: ExpressionToCommand(t.Expression),
			index:    ExpressionToCommand(t.Index),
		}
	case *ast.MapLiteral:
		entries := make([]MapEntry, len(t.Entries))
		for i, entry := range t.Entries {
			mapEntry := MapEntry{
				key:   ExpressionToCommand(entry.Key),
				value: ExpressionToCommand(entry.Value),
			}
			entries[i] = mapEntry
		}
		return &MapCommand{
			entries: entries,
		}
	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
