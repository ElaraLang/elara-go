package interpreter

import (
	"github.com/ElaraLang/elara/lexer"
	"github.com/ElaraLang/elara/parser"
	_ "github.com/ElaraLang/elara/util"
	"reflect"
	"strings"
)

type Command interface {
	Exec(ctx *Context) *Value
}

type DefineVarCommand struct {
	Name        string
	Mutable     bool
	Type        *parser.Type
	value       Command
	runtimeType *Type
}

func (c *DefineVarCommand) getType(ctx *Context) *Type {
	if c.runtimeType == nil {
		c.runtimeType = FromASTType(*c.Type, ctx)
	}
	return c.runtimeType
}
func (c *DefineVarCommand) Exec(ctx *Context) *Value {
	if ctx.FindVariableMaxDepth(c.Name, 1) != nil {
		panic("Variable named " + c.Name + " already exists")
	}
	value := c.value.Exec(ctx)

	if value == nil {
		panic("Command " + reflect.TypeOf(c.value).Name() + " returned nil")
	}

	variableType := c.getType(ctx)
	if !variableType.Accepts(*value.Type) {
		panic("Cannot use value of type " + value.Type.Name + " in place of " + variableType.Name + " for variable " + c.Name)
	}
	variable := Variable{
		Name:    c.Name,
		Mutable: c.Mutable,
		Type:    *variableType,
		Value:   value,
	}

	ctx.DefineVariable(c.Name, variable)
	return nil
}

type AssignmentCommand struct {
	Name  string
	value Command
}

func (c *AssignmentCommand) Exec(ctx *Context) *Value {
	variable := ctx.FindVariable(c.Name)
	if variable == nil {
		panic("No such variable " + c.Name)
	}

	if !variable.Mutable {
		panic("Cannot reassign immutable variable " + c.Name)
	}

	value := c.value.Exec(ctx)

	if !variable.Type.Accepts(*value.Type) {
		panic("Cannot reassign variable " + c.Name + " of type " + variable.Type.Name + " to value " + *value.String() + " of type " + value.Type.Name)
	}

	variable.Value = value
	return nil
}

type VariableCommand struct {
	Variable string
}

func (c *VariableCommand) Exec(ctx *Context) *Value {
	if ctx.receiver != nil {
		receiver := ctx.receiver
		asInstance, isInstance := ctx.receiver.Value.(Instance)
		if isInstance {
			instanceVariable, exists := asInstance.Values[c.Variable]
			if exists {
				return instanceVariable
			}
		}
		typeVariable, exists := receiver.Type.variables.m[c.Variable]
		if exists {
			return typeVariable.Value
		}
	}

	variable := ctx.FindVariable(c.Variable)
	if variable == nil {
		param := ctx.FindParameter(c.Variable)
		if param == nil {
			constructor := ctx.FindConstructor(c.Variable)
			if constructor == nil {
				panic("No such variable or parameter or constructor " + c.Variable)
			}
			return constructor
		}
		return param
	}
	return variable.Value
}

type InvocationCommand struct {
	Invoking Command
	args     []Command
}

func (c *InvocationCommand) Exec(ctx *Context) *Value {
	context, usingReceiver := c.Invoking.(*ContextCommand)

	if !usingReceiver {
		val := c.Invoking.Exec(ctx)
		fun, ok := val.Value.(Function)
		if !ok {
			panic("Cannot invoke non-value")
		}

		argValues := make([]*Value, len(c.args))
		for i, arg := range c.args {
			argValues[i] = arg.Exec(ctx)
		}
		return fun.Exec(ctx, nil, argValues)
	}

	//ContextCommand seems to think it's a special case... because it is.
	receiver := context.receiver.Exec(ctx)
	value, ok := receiver.Type.variables.m[context.variable]
	if !ok {

		argTypes := make([]string, len(c.args))
		for i, arg := range c.args {
			argTypes[i] = arg.Exec(ctx).Type.Name
		}
		panic("No such function " + receiver.Type.Name + "::" + context.variable + "(" + strings.Join(argTypes, ", ") + ")")
	}
	function, ok := value.Value.Value.(Function)
	if !ok {
		panic("Cannot invoke non-function " + value.string())
	}
	exec := context.receiver.Exec(ctx)

	argValues := make([]*Value, len(c.args))
	for i, arg := range c.args {
		argValues[i] = arg.Exec(ctx)
	}
	return function.Exec(ctx, exec, argValues)

}

type AbstractCommand struct {
	content func(ctx *Context) *Value
}

func (c *AbstractCommand) Exec(ctx *Context) *Value {
	return c.content(ctx)
}

func NewAbstractCommand(content func(ctx *Context) *Value) *AbstractCommand {
	return &AbstractCommand{
		content: content,
	}
}

type LiteralCommand struct {
	value Value
}

func (c *LiteralCommand) Exec(_ *Context) *Value {
	return &c.value
}

type FunctionLiteralCommand struct {
	name       *string
	parameters []parser.FunctionArgument
	returnType parser.Type
	body       Command
}

func (c *FunctionLiteralCommand) Exec(ctx *Context) *Value {
	params := make([]Parameter, len(c.parameters))

	for i, parameter := range c.parameters {
		paramType := FromASTType(parameter.Type, ctx)
		params[i] = Parameter{
			Type: *paramType,
			Name: parameter.Name,
		}
	}

	returnType := FromASTType(c.returnType, ctx)

	fun := Function{
		name: c.name,
		Signature: Signature{
			Parameters: params,
			ReturnType: *returnType,
		},
		Body: c.body,
	}

	functionType := FunctionType(fun)

	return &Value{
		Type:  functionType,
		Value: fun,
	}
}

type BinaryOperatorCommand struct {
	lhs Command
	op  func(ctx *Context, lhs *Value, rhs *Value) *Value
	rhs Command
}

func (c *BinaryOperatorCommand) Exec(ctx *Context) *Value {
	lhs := c.lhs.Exec(ctx)
	rhs := c.rhs.Exec(ctx)

	return c.op(ctx, lhs, rhs)
}

type BlockCommand struct {
	lines []Command
}

func (c *BlockCommand) Exec(ctx *Context) *Value {
	var last *Value
	for _, line := range c.lines {
		last = line.Exec(ctx)
	}
	return last
}

type ContextCommand struct {
	receiver Command
	variable string
}

func (c *ContextCommand) Exec(ctx *Context) *Value {
	receiver := c.receiver.Exec(ctx)
	instance, isInstance := receiver.Value.(Instance)
	if isInstance {
		variable, ok := instance.Values[c.variable]
		if !ok {
			value, ok := receiver.Type.variables.m[c.variable]
			if !ok {
				panic("No such variable " + c.variable + " on type " + receiver.Type.Name)
			}
			return value.Value
		}
		return variable
	}
	variable, ok := receiver.Type.variables.m[c.variable]
	if !ok {
		panic("No such variable " + c.variable + " on type " + receiver.Type.Name)
	}

	return variable.Value
}

type IfElseCommand struct {
	condition  Command
	ifBranch   Command
	elseBranch Command
}

func (c *IfElseCommand) Exec(ctx *Context) *Value {
	condition := c.condition.Exec(ctx)
	value, ok := condition.Value.(bool)
	if !ok {
		panic("If statements requires boolean value")
	}

	if value {
		return c.ifBranch.Exec(ctx)
	} else if c.elseBranch != nil {
		return c.elseBranch.Exec(ctx)
	} else {
		return nil
	}
}

type IfElseExpressionCommand struct {
	condition  Command
	ifBranch   []Command
	ifResult   Command
	elseBranch []Command
	elseResult Command
}

func (c *IfElseExpressionCommand) Exec(ctx *Context) *Value {
	condition := c.condition.Exec(ctx)
	value, ok := condition.Value.(bool)
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

func (c *ReturnCommand) Exec(ctx *Context) *Value {
	panic(c.returning.Exec(ctx))
}

type NamespaceCommand struct {
	namespace string
}

func (c *NamespaceCommand) Exec(ctx *Context) *Value {
	ctx.Init(c.namespace)
	return nil
}

type ImportCommand struct {
	imports []string
}

func (c *ImportCommand) Exec(ctx *Context) *Value {
	for _, s := range c.imports {
		ctx.Import(s)
	}
	return nil
}

type StructDefCommand struct {
	name   string
	fields []parser.StructField
}

func (c *StructDefCommand) Exec(ctx *Context) *Value {
	variables := NewVariableMap()

	for _, field := range c.fields {
		var Type Type
		if field.FieldType == nil {
			Type = *AnyType
		} else {
			Type = *FromASTType(*field.FieldType, ctx)
		}
		variables.Set(field.Identifier, Variable{
			Name:    field.Identifier,
			Mutable: field.Mutable,
			Type:    Type,
			Value:   nil,
		})
	}
	ctx.types[c.name] = Type{
		Name:      c.name,
		variables: *variables,
	}
	return nil
}

type ExtendCommand struct {
	Type       string
	statements []Command
}

func (c *ExtendCommand) Exec(ctx *Context) *Value {
	extending := ctx.FindType(c.Type)
	if extending == nil {
		panic("No such type " + c.Type)
	}
	for _, statement := range c.statements {
		switch statement := statement.(type) {
		case *DefineVarCommand:
			defVar := statement
			value := defVar.value.Exec(ctx)
			variable := &Variable{
				Name:    defVar.Name,
				Mutable: defVar.Mutable,
				Type:    *defVar.getType(ctx),
				Value:   value,
			}
			extending.variables.Set(defVar.Name, *variable)
		}
	}
	return nil
}

func ToCommand(statement parser.Stmt) Command {
	switch t := statement.(type) {
	case parser.VarDefStmt:
		valueExpr := NamedExpressionToCommand(t.Value, &t.Identifier)
		return &DefineVarCommand{
			Name:    t.Identifier,
			Mutable: t.Mutable,
			Type:    &t.Type,
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
		return &ReturnCommand{
			ExpressionToCommand(t.Returning),
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
				panic("Could not convert expression " + reflect.TypeOf(arg).Name() + " to command")
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
	}

	panic("Could not handle " + reflect.TypeOf(expr).Name())
}
