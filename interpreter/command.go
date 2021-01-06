package interpreter

import (
	"fmt"
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
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
	Type        parser.Type
	value       Command
	runtimeType Type
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
	var value *Value
	foundVar, _ := ctx.FindVariableMaxDepth(c.Name, 1)
	if foundVar != nil {
		asFunction, isFunction := foundVar.Value.Value.(*Function)
		if isFunction {
			value = c.value.Exec(ctx).Unwrap()
			valueAsFunction, valueIsFunction := value.Value.(*Function)
			if valueIsFunction && !asFunction.Signature.Accepts(&valueAsFunction.Signature, false) {
				//We'll allow it because the functions have different arity
			} else {
				panic("Variable named " + c.Name + " already exists with the current signature")
			}
		} else {
			panic("Variable named " + c.Name + " already exists")
		}
	}
	if value == nil {
		value = c.value.Exec(ctx).Unwrap()
	}

	if value == nil {
		panic("Command " + reflect.TypeOf(c.value).Name() + " returned nil")
	}

	variableType := c.getType(ctx)
	if variableType != nil {
		if !variableType.Accepts(value.Type) {
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

	ctx.DefineVariable(c.Name, variable)
	return NilValue()
}

type AssignmentCommand struct {
	Name  string
	value Command
}

func (c *AssignmentCommand) Exec(ctx *Context) *ReturnedValue {
	variable := ctx.FindVariable(c.Name)
	if variable == nil {
		panic("No such variable " + c.Name)
	}

	if !variable.Mutable {
		panic("Cannot reassign immutable variable " + c.Name)
	}

	value := c.value.Exec(ctx).Unwrap()

	if !variable.Type.Accepts(value.Type) {
		panic("Cannot reassign variable " + c.Name + " of type " + variable.Type.Name() + " to value " + value.String() + " of type " + value.Type.Name())
	}

	variable.Value = value
	return NilValue()
}

type VariableCommand struct {
	Variable string
}

func (c *VariableCommand) findVariable(ctx *Context) *Variable {
	variable := ctx.FindVariable(c.Variable)
	return variable
}
func (c *VariableCommand) Exec(ctx *Context) *ReturnedValue {
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
	return NonReturningValue(constructor)
}

type InvocationCommand struct {
	Invoking Command
	args     []Command

	cachedFun *Function
}

func (c *InvocationCommand) findReceiverFunction(ctx *Context, receiver *Value, argValues []*Value, functionName string) *Function {
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
	receiverFunction := ctx.FindFunction(functionName, receiverSignature)
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
			panic("Cannot invoke non-value")
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

	if c.cachedFun == nil {
		switch t := context.receiver.(type) {
		case *VariableCommand:
			variable := t.findVariable(ctx)
			if variable != nil && !variable.Mutable {
				receiver = variable.Value
				c.cachedFun = c.findReceiverFunction(ctx, receiver, argValues, context.variable)
			}
		}
	}
	receiver = context.receiver.Exec(ctx).Unwrap()

	if c.cachedFun != nil {
		argValuesAndSelf := []*Value{receiver}
		argValuesAndSelf = append(argValuesAndSelf, argValues...)
		return NonReturningValue(c.cachedFun.Exec(ctx, argValuesAndSelf))
	}
	structType, isStruct := receiver.Type.(*StructType)

	if isStruct {
		value, ok := structType.GetProperty(context.variable)
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

	//Look for a receiver

	receiverFunction := c.findReceiverFunction(ctx, receiver, argValues, context.variable)
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
	value Value
}

func (c *LiteralCommand) Exec(_ *Context) *ReturnedValue {
	return NonReturningValue(&c.value)
}

type FunctionLiteralCommand struct {
	name       *string
	parameters []parser.FunctionArgument
	returnType parser.Type //Can be nil - infer return type
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
			Name:     parameter.Name,
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
	lines []Command
}

func (c *BlockCommand) Exec(ctx *Context) *ReturnedValue {
	var last *ReturnedValue
	for _, line := range c.lines {
		val := line.Exec(ctx)
		if val.IsReturning {
			return val
		}
		last = val
	}
	return last
}

type ContextCommand struct {
	receiver Command
	variable string
}

func (c *ContextCommand) Exec(ctx *Context) *ReturnedValue {
	receiver := c.receiver.Exec(ctx).Unwrap()
	switch val := receiver.Value.(type) {
	case *Collection:
		{
			switch c.variable {
			case "size":
				return NonReturningValue(IntValue(int64(len(val.Elements))))
			default:
				panic("Unknown property" + c.variable)
			}
		}
	case *Instance:
		{
			return NonReturningValue(val.Values[c.variable])
		}
	default:
		panic("Unsupported receiver " + util.Stringify(receiver))
	}

	//instance, isInstance := receiver.Value.(Instance)
	//if isInstance {
	//	variable, ok := instance.Values[c.variable]
	//	if !ok {
	//		value, ok := receiver.Type.variables.m[c.variable]
	//		if !ok {
	//			panic("No such variable " + c.variable + " on type " + receiver.Type.Name)
	//		}
	//		return value.Value
	//	}
	//	return variable
	//}
	//variable, ok := receiver.Type.variables.m[c.variable]
	//if !ok {
	//	panic("No such variable " + c.variable + " on type " + receiver.Type.Name)
	//}
	//
	//return variable.Value
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
	ifBranch   []Command
	ifResult   Command
	elseBranch []Command
	elseResult Command
}

func (c *IfElseExpressionCommand) Exec(ctx *Context) *ReturnedValue {
	condition := c.condition.Exec(ctx)
	value, ok := condition.Unwrap().Value.(bool)
	if !ok {
		panic("If statements requires boolean value")
	}

	if value {
		if c.ifBranch != nil {
			for _, cmd := range c.ifBranch {
				cmd.Exec(ctx)
			}
		}
		return c.ifResult.Exec(ctx)
	} else {
		if c.elseBranch != nil {
			for _, cmd := range c.elseBranch {
				cmd.Exec(ctx)
			}
		}
		return c.elseResult.Exec(ctx)
	}
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
	imports []string
}

func (c *ImportCommand) Exec(ctx *Context) *ReturnedValue {
	for _, s := range c.imports {
		ctx.Import(s)
	}
	return NilValue()
}

type StructDefCommand struct {
	name   string
	fields []parser.StructField
}

func (c *StructDefCommand) Exec(ctx *Context) *ReturnedValue {

	properties := make([]Property, len(c.fields))
	propertyPositions := map[string]int{}

	for i, field := range c.fields {
		var Type Type
		if field.FieldType == nil {
			Type = AnyType
		} else {
			Type = FromASTType(*field.FieldType, ctx)
		}

		var defaultValue *Value
		if field.Default != nil {
			defaultValue = ExpressionToCommand(field.Default).Exec(ctx).Unwrap()
		}

		modifiers := uint(0)
		if field.Mutable {
			modifiers |= Mut
		}
		properties[i] = Property{
			Name:         field.Identifier,
			Modifiers:    modifiers,
			Type:         Type,
			DefaultValue: defaultValue,
		}
		propertyPositions[field.Identifier] = i
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
}

func (c *ExtendCommand) Exec(ctx *Context) *ReturnedValue {
	panic("TODO ExtendCommand")
	//extending := ctx.FindType(c.Type)
	//if extending == nil {
	//	panic("No such type " + c.Type)
	//}
	//for _, statement := range c.statements {
	//	switch statement := statement.(type) {
	//	case *DefineVarCommand:
	//		//defVar := statement
	//		//value := defVar.value.Exec(ctx)
	//		//variable := &Variable{
	//		//	Name:    defVar.Name,
	//		//	Mutable: defVar.Mutable,
	//		//	Type:    defVar.getType(ctx),
	//		//	Value:   value,
	//		//}
	//
	//		//extending.variables.Set(defVar.Name, variable)
	//	}
	//}
	//return nil
}

type TypeCheckCommand struct {
	expression Command
	checkType  parser.Type
}

func (c *TypeCheckCommand) Exec(ctx *Context) *ReturnedValue {
	checkAgainst := FromASTType(c.checkType, ctx)
	if checkAgainst == nil {
		panic("No such type " + util.Stringify(c.checkType))
	}
	res := c.expression.Exec(ctx).Unwrap()
	is := res.Type.Accepts(checkAgainst)
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
	coll, isColl := c.checking.Exec(ctx).Unwrap().Value.(*Collection)
	if !isColl {
		panic("Indexed access not supported for non-collection type")
	}
	index, isInt := c.index.Exec(ctx).Unwrap().Value.(int64)
	if !isInt {
		panic("Index was not an integer")
	}
	return NonReturningValue(coll.Elements[index])
}

func ToCommand(statement parser.Stmt) Command {
	switch t := statement.(type) {
	case parser.VarDefStmt:
		valueExpr := NamedExpressionToCommand(t.Value, &t.Identifier)
		return &DefineVarCommand{
			Name:    t.Identifier,
			Mutable: t.Mutable,
			Type:    t.Type,
			value:   valueExpr,
		}
	case parser.ExpressionStmt:
		return ExpressionToCommand(t.Expr)

	case parser.BlockStmt:
		commands := make([]Command, len(t.Stmts))
		for i, stmt := range t.Stmts {
			commands[i] = ToCommand(stmt)
			_, isReturn := commands[i].(*ReturnCommand)
			//Small optimisation, it's not worth transforming anything that won't ever be reached
			if isReturn {
				return &BlockCommand{lines: commands[:i+1]}
			}
		}

		return &BlockCommand{lines: commands}

	case parser.IfElseStmt:
		condition := ExpressionToCommand(t.Condition)
		ifBranch := ToCommand(t.MainBranch)
		var elseBranch Command
		if t.ElseBranch != nil {
			elseBranch = ToCommand(t.ElseBranch)
		}

		return &IfElseCommand{
			condition:  condition,
			ifBranch:   ifBranch,
			elseBranch: elseBranch,
		}

	case parser.ReturnStmt:
		if t.Returning != nil {
			return &ReturnCommand{
				ExpressionToCommand(t.Returning),
			}
		}
		return &ReturnCommand{
			nil,
		}
	case parser.NamespaceStmt:
		return &NamespaceCommand{
			namespace: t.Namespace,
		}
	case parser.ImportStmt:
		return &ImportCommand{
			imports: t.Imports,
		}
	case parser.StructDefStmt:
		name := t.Identifier
		return &StructDefCommand{
			name:   name,
			fields: t.StructFields,
		}

	case parser.ExtendStmt:
		commands := make([]Command, len(t.Body.Stmts))
		for i, stmt := range t.Body.Stmts {
			commands[i] = ToCommand(stmt)
		}
		return &ExtendCommand{
			Type:       t.Identifier,
			statements: commands,
		}

	case parser.WhileStmt:
		return &WhileCommand{
			condition: ExpressionToCommand(t.Condition),
			body:      ToCommand(t.Body),
		}
	}

	panic("Could not handle " + reflect.TypeOf(statement).Name())
}

func ExpressionToCommand(expr parser.Expr) Command {
	return NamedExpressionToCommand(expr, nil)
}

func NamedExpressionToCommand(expr parser.Expr, name *string) Command {

	switch t := expr.(type) {
	case parser.VariableExpr:
		return &VariableCommand{Variable: t.Identifier}

	case parser.InvocationExpr:
		fun := ExpressionToCommand(t.Invoker)
		args := make([]Command, 0)
		for _, arg := range t.Args {
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

	case parser.StringLiteralExpr:
		str := t.Value
		value := Value{
			Type:  StringType,
			Value: str,
		}
		return &LiteralCommand{value: value}

	case parser.IntegerLiteralExpr:
		integer := t.Value
		value := Value{
			Type:  IntType,
			Value: integer,
		}
		return &LiteralCommand{value: value}
	case parser.FloatLiteralExpr:
		float := t.Value
		value := Value{
			Type:  FloatType,
			Value: float,
		}
		return &LiteralCommand{value: value}
	case parser.BooleanLiteralExpr:
		boolean := t.Value
		value := Value{
			Type:  BooleanType,
			Value: boolean,
		}
		return &LiteralCommand{value: value}

	case parser.BinaryExpr:
		lhs := t.Lhs
		lhsCmd := ExpressionToCommand(lhs)
		op := t.Op
		rhs := t.Rhs
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
	case parser.FuncDefExpr:
		return &FunctionLiteralCommand{
			name:       name,
			parameters: t.Arguments,
			returnType: t.ReturnType,
			body:       ToCommand(t.Statement),
		}

	case parser.ContextExpr:
		contextCmd := ExpressionToCommand(t.Context)
		varName := t.Variable.Identifier
		return &ContextCommand{
			contextCmd,
			varName,
		}

	case parser.AssignmentExpr:
		//TODO contexts
		name := t.Identifier
		valueCmd := NamedExpressionToCommand(t.Value, &name)
		return &AssignmentCommand{
			Name:  name,
			value: valueCmd,
		}

	case parser.IfElseExpr:
		condition := ExpressionToCommand(t.Condition)
		var ifBranch []Command
		if t.IfBranch != nil {
			ifBranch = make([]Command, len(t.IfBranch))
			for i, stmt := range t.IfBranch {
				ifBranch[i] = ToCommand(stmt)
			}
		}
		ifResult := ExpressionToCommand(t.IfResult)

		var elseBranch []Command
		if t.ElseBranch != nil {
			elseBranch = make([]Command, len(t.ElseBranch))
			for i, stmt := range t.ElseBranch {
				elseBranch[i] = ToCommand(stmt)
			}
		}
		elseResult := ExpressionToCommand(t.ElseResult)

		return &IfElseExpressionCommand{
			condition:  condition,
			ifBranch:   ifBranch,
			ifResult:   ifResult,
			elseBranch: elseBranch,
			elseResult: elseResult,
		}

	case parser.TypeCheckExpr:
		return &TypeCheckCommand{
			expression: ExpressionToCommand(t.Expr),
			checkType:  t.Type,
		}
	case parser.GroupExpr:
		return ExpressionToCommand(t.Group)

	case parser.CollectionExpr:
		elements := make([]Command, len(t.Elements))
		for i, element := range t.Elements {
			elements[i] = ExpressionToCommand(element)
		}
		return &CollectionCommand{Elements: elements}

	case parser.AccessExpr:
		return &AccessCommand{
			checking: ExpressionToCommand(t.Expr),
			index:    ExpressionToCommand(t.Index),
		}
	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
