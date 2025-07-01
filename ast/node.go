package ast

import (
	"github.com/mredencom/expr/lexer"
	"github.com/mredencom/expr/types"
)

// Node represents any node in the AST
type Node interface {
	Type() types.TypeInfo
	Position() lexer.Position
	String() string
}

// Statement represents a statement node
type Statement interface {
	Node
	statementNode()
}

// Expression represents an expression node
type Expression interface {
	Node
	expressionNode()
}

// Program represents the root of every AST
type Program struct {
	Statements []Statement
	Pos        lexer.Position
}

func (p *Program) Type() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindNil, Name: "program"}
}

func (p *Program) Position() lexer.Position {
	return p.Pos
}

func (p *Program) String() string {
	return "program"
}

func (p *Program) statementNode() {}

// ExpressionStatement represents an expression used as a statement
type ExpressionStatement struct {
	Expression Expression
	Pos        lexer.Position
}

func (es *ExpressionStatement) Type() types.TypeInfo {
	if es.Expression != nil {
		return es.Expression.Type()
	}
	return types.TypeInfo{Kind: types.KindNil, Name: "void"}
}

func (es *ExpressionStatement) Position() lexer.Position {
	return es.Pos
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

func (es *ExpressionStatement) statementNode() {}

// Identifier represents an identifier expression
type Identifier struct {
	Value    string
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (i *Identifier) Type() types.TypeInfo {
	return i.TypeInfo
}

func (i *Identifier) Position() lexer.Position {
	return i.Pos
}

func (i *Identifier) String() string {
	return i.Value
}

func (i *Identifier) expressionNode() {}

// Literal represents a literal value
type Literal struct {
	Value types.Value
	Pos   lexer.Position
}

func (l *Literal) Type() types.TypeInfo {
	if l.Value != nil {
		return l.Value.Type()
	}
	return types.TypeInfo{Kind: types.KindNil, Name: "nil"}
}

func (l *Literal) Position() lexer.Position {
	return l.Pos
}

func (l *Literal) String() string {
	if l.Value != nil {
		return l.Value.String()
	}
	return "nil"
}

func (l *Literal) expressionNode() {}

// InfixExpression represents an infix expression (e.g., a + b)
type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (ie *InfixExpression) Type() types.TypeInfo {
	return ie.TypeInfo
}

func (ie *InfixExpression) Position() lexer.Position {
	return ie.Pos
}

func (ie *InfixExpression) String() string {
	return "(" + ie.Left.String() + " " + ie.Operator + " " + ie.Right.String() + ")"
}

func (ie *InfixExpression) expressionNode() {}

// PrefixExpression represents a prefix expression (e.g., !x, -x)
type PrefixExpression struct {
	Operator string
	Right    Expression
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (pe *PrefixExpression) Type() types.TypeInfo {
	return pe.TypeInfo
}

func (pe *PrefixExpression) Position() lexer.Position {
	return pe.Pos
}

func (pe *PrefixExpression) String() string {
	return "(" + pe.Operator + pe.Right.String() + ")"
}

func (pe *PrefixExpression) expressionNode() {}

// CallExpression represents a function call expression
type CallExpression struct {
	Function  Expression // identifier or function literal
	Arguments []Expression
	TypeInfo  types.TypeInfo
	Pos       lexer.Position
}

func (ce *CallExpression) Type() types.TypeInfo {
	return ce.TypeInfo
}

func (ce *CallExpression) Position() lexer.Position {
	return ce.Pos
}

func (ce *CallExpression) String() string {
	result := ce.Function.String() + "("
	for i, arg := range ce.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

func (ce *CallExpression) expressionNode() {}

// IndexExpression represents indexing expressions (e.g., arr[0], obj["key"])
type IndexExpression struct {
	Left     Expression
	Index    Expression
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (ie *IndexExpression) Type() types.TypeInfo {
	return ie.TypeInfo
}

func (ie *IndexExpression) Position() lexer.Position {
	return ie.Pos
}

func (ie *IndexExpression) String() string {
	return "(" + ie.Left.String() + "[" + ie.Index.String() + "])"
}

func (ie *IndexExpression) expressionNode() {}

// MemberExpression represents member access expressions (e.g., obj.field)
type MemberExpression struct {
	Object   Expression
	Property Expression // Can be Identifier or WildcardExpression
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (me *MemberExpression) Type() types.TypeInfo {
	return me.TypeInfo
}

func (me *MemberExpression) Position() lexer.Position {
	return me.Pos
}

func (me *MemberExpression) String() string {
	return me.Object.String() + "." + me.Property.String()
}

func (me *MemberExpression) expressionNode() {}

// ConditionalExpression represents ternary conditional expressions (e.g., a ? b : c)
type ConditionalExpression struct {
	Test        Expression
	Consequent  Expression
	Alternative Expression
	TypeInfo    types.TypeInfo
	Pos         lexer.Position
}

func (ce *ConditionalExpression) Type() types.TypeInfo {
	return ce.TypeInfo
}

func (ce *ConditionalExpression) Position() lexer.Position {
	return ce.Pos
}

func (ce *ConditionalExpression) String() string {
	return "(" + ce.Test.String() + " ? " + ce.Consequent.String() + " : " + ce.Alternative.String() + ")"
}

func (ce *ConditionalExpression) expressionNode() {}

// ArrayLiteral represents array literal expressions (e.g., [1, 2, 3])
type ArrayLiteral struct {
	Elements []Expression
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (al *ArrayLiteral) Type() types.TypeInfo {
	return al.TypeInfo
}

func (al *ArrayLiteral) Position() lexer.Position {
	return al.Pos
}

func (al *ArrayLiteral) String() string {
	result := "["
	for i, elem := range al.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

func (al *ArrayLiteral) expressionNode() {}

// MapLiteral represents map literal expressions (e.g., {"key": "value"})
type MapLiteral struct {
	Pairs    []MapPair
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

type MapPair struct {
	Key   Expression
	Value Expression
}

func (ml *MapLiteral) Type() types.TypeInfo {
	return ml.TypeInfo
}

func (ml *MapLiteral) Position() lexer.Position {
	return ml.Pos
}

func (ml *MapLiteral) String() string {
	result := "{"
	for i, pair := range ml.Pairs {
		if i > 0 {
			result += ", "
		}
		result += pair.Key.String() + ": " + pair.Value.String()
	}
	result += "}"
	return result
}

func (ml *MapLiteral) expressionNode() {}

// BuiltinExpression represents built-in function calls
type BuiltinExpression struct {
	Name      string
	Arguments []Expression
	TypeInfo  types.TypeInfo
	Pos       lexer.Position
}

func (be *BuiltinExpression) Type() types.TypeInfo {
	return be.TypeInfo
}

func (be *BuiltinExpression) Position() lexer.Position {
	return be.Pos
}

func (be *BuiltinExpression) String() string {
	result := be.Name + "("
	for i, arg := range be.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

func (be *BuiltinExpression) expressionNode() {}

// VariableExpression represents variable references with resolved type
type VariableExpression struct {
	Name     string
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (ve *VariableExpression) Type() types.TypeInfo {
	return ve.TypeInfo
}

func (ve *VariableExpression) Position() lexer.Position {
	return ve.Pos
}

func (ve *VariableExpression) String() string {
	return ve.Name
}

func (ve *VariableExpression) expressionNode() {}

// LambdaExpression represents lambda functions (e.g., x => x * 2)
type LambdaExpression struct {
	Parameters []string   // Parameter names
	Body       Expression // Lambda body expression
	TypeInfo   types.TypeInfo
	Pos        lexer.Position
}

func (le *LambdaExpression) Type() types.TypeInfo {
	return le.TypeInfo
}

func (le *LambdaExpression) Position() lexer.Position {
	return le.Pos
}

func (le *LambdaExpression) String() string {
	params := ""
	for i, param := range le.Parameters {
		if i > 0 {
			params += ", "
		}
		params += param
	}
	if len(le.Parameters) == 1 {
		return params + " => " + le.Body.String()
	}
	return "(" + params + ") => " + le.Body.String()
}

func (le *LambdaExpression) expressionNode() {}

// PipeExpression represents pipeline operations (e.g., data | filter(...) | map(...))
type PipeExpression struct {
	Left     Expression // Left side of pipe
	Right    Expression // Right side (usually function call)
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (pe *PipeExpression) Type() types.TypeInfo {
	return pe.TypeInfo
}

func (pe *PipeExpression) Position() lexer.Position {
	return pe.Pos
}

func (pe *PipeExpression) String() string {
	return pe.Left.String() + " | " + pe.Right.String()
}

func (pe *PipeExpression) expressionNode() {}

// WildcardExpression represents wildcard expressions (e.g., user.*, *.field)
type WildcardExpression struct {
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (we *WildcardExpression) Type() types.TypeInfo {
	return we.TypeInfo
}

func (we *WildcardExpression) Position() lexer.Position {
	return we.Pos
}

func (we *WildcardExpression) String() string {
	return "*"
}

func (we *WildcardExpression) expressionNode() {}

// PlaceholderExpression represents a pipeline placeholder expression (#)
type PlaceholderExpression struct {
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (pe *PlaceholderExpression) Type() types.TypeInfo {
	return pe.TypeInfo
}

func (pe *PlaceholderExpression) Position() lexer.Position {
	return pe.Pos
}

func (pe *PlaceholderExpression) String() string {
	return "#"
}

func (pe *PlaceholderExpression) expressionNode() {}

// OptionalChainingExpression represents an optional chaining expression (e.g., obj?.property)
type OptionalChainingExpression struct {
	Object   Expression
	Property Expression // Can be Identifier or computed property
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (oce *OptionalChainingExpression) Type() types.TypeInfo {
	return oce.TypeInfo
}

func (oce *OptionalChainingExpression) Position() lexer.Position {
	return oce.Pos
}

func (oce *OptionalChainingExpression) String() string {
	return oce.Object.String() + "?." + oce.Property.String()
}

func (oce *OptionalChainingExpression) expressionNode() {}

// NullCoalescingExpression represents a null coalescing expression (e.g., a ?? b)
type NullCoalescingExpression struct {
	Left     Expression // Left operand
	Right    Expression // Right operand (default value)
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (nce *NullCoalescingExpression) Type() types.TypeInfo {
	return nce.TypeInfo
}

func (nce *NullCoalescingExpression) Position() lexer.Position {
	return nce.Pos
}

func (nce *NullCoalescingExpression) String() string {
	return "(" + nce.Left.String() + " ?? " + nce.Right.String() + ")"
}

func (nce *NullCoalescingExpression) expressionNode() {}

// ImportStatement represents import statements (e.g., import "math" as m)
type ImportStatement struct {
	ModuleName string // The module name to import (e.g., "math")
	Alias      string // The alias for the module (e.g., "m")
	Pos        lexer.Position
}

func (is *ImportStatement) Type() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindNil, Name: "nil", Size: 0}
}

func (is *ImportStatement) Position() lexer.Position {
	return is.Pos
}

func (is *ImportStatement) String() string {
	if is.Alias != "" {
		return "import '" + is.ModuleName + "' as " + is.Alias
	}
	return "import '" + is.ModuleName + "'"
}

func (is *ImportStatement) statementNode() {}

// ModuleCallExpression represents module function calls (e.g., m.sqrt(16))
type ModuleCallExpression struct {
	Module    string       // Module alias (e.g., "m")
	Function  string       // Function name (e.g., "sqrt")
	Arguments []Expression // Function arguments
	TypeInfo  types.TypeInfo
	Pos       lexer.Position
}

func (mce *ModuleCallExpression) Type() types.TypeInfo {
	return mce.TypeInfo
}

func (mce *ModuleCallExpression) Position() lexer.Position {
	return mce.Pos
}

func (mce *ModuleCallExpression) String() string {
	result := mce.Module + "." + mce.Function + "("
	for i, arg := range mce.Arguments {
		if i > 0 {
			result += ", "
		}
		result += arg.String()
	}
	result += ")"
	return result
}

func (mce *ModuleCallExpression) expressionNode() {}

// DestructuringAssignment represents destructuring assignment (e.g., [a, b] = [1, 2])
type DestructuringAssignment struct {
	Left  DestructuringPattern // The destructuring pattern (left side)
	Right Expression           // The value expression (right side)
	Pos   lexer.Position
}

func (da *DestructuringAssignment) Type() types.TypeInfo {
	return types.TypeInfo{Kind: types.KindNil, Name: "nil", Size: 0}
}

func (da *DestructuringAssignment) Position() lexer.Position {
	return da.Pos
}

func (da *DestructuringAssignment) String() string {
	return da.Left.String() + " = " + da.Right.String()
}

func (da *DestructuringAssignment) statementNode() {}

// DestructuringPattern interface for all destructuring patterns
type DestructuringPattern interface {
	Node
	destructuringPatternNode()
}

// ArrayDestructuringPattern represents array destructuring pattern [a, b, c]
type ArrayDestructuringPattern struct {
	Elements []DestructuringElement // Elements in the array pattern
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (adp *ArrayDestructuringPattern) Type() types.TypeInfo {
	return adp.TypeInfo
}

func (adp *ArrayDestructuringPattern) Position() lexer.Position {
	return adp.Pos
}

func (adp *ArrayDestructuringPattern) String() string {
	result := "["
	for i, elem := range adp.Elements {
		if i > 0 {
			result += ", "
		}
		result += elem.String()
	}
	result += "]"
	return result
}

func (adp *ArrayDestructuringPattern) destructuringPatternNode() {}

// ObjectDestructuringPattern represents object destructuring pattern {name, age}
type ObjectDestructuringPattern struct {
	Properties []ObjectDestructuringProperty // Properties in the object pattern
	TypeInfo   types.TypeInfo
	Pos        lexer.Position
}

func (odp *ObjectDestructuringPattern) Type() types.TypeInfo {
	return odp.TypeInfo
}

func (odp *ObjectDestructuringPattern) Position() lexer.Position {
	return odp.Pos
}

func (odp *ObjectDestructuringPattern) String() string {
	result := "{"
	for i, prop := range odp.Properties {
		if i > 0 {
			result += ", "
		}
		result += prop.String()
	}
	result += "}"
	return result
}

func (odp *ObjectDestructuringPattern) destructuringPatternNode() {}

// DestructuringElement interface for elements in destructuring patterns
type DestructuringElement interface {
	Node
	destructuringElementNode()
}

// IdentifierElement represents a simple identifier in destructuring (e.g., 'a' in [a, b])
type IdentifierElement struct {
	Name     string     // Variable name
	Default  Expression // Default value (optional)
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (ie *IdentifierElement) Type() types.TypeInfo {
	return ie.TypeInfo
}

func (ie *IdentifierElement) Position() lexer.Position {
	return ie.Pos
}

func (ie *IdentifierElement) String() string {
	if ie.Default != nil {
		return ie.Name + " = " + ie.Default.String()
	}
	return ie.Name
}

func (ie *IdentifierElement) destructuringElementNode() {}

// RestElement represents rest element in destructuring (e.g., '...rest' in [a, ...rest])
type RestElement struct {
	Name     string // Variable name for rest elements
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (re *RestElement) Type() types.TypeInfo {
	return re.TypeInfo
}

func (re *RestElement) Position() lexer.Position {
	return re.Pos
}

func (re *RestElement) String() string {
	return "..." + re.Name
}

func (re *RestElement) destructuringElementNode() {}

// ObjectDestructuringProperty represents a property in object destructuring
type ObjectDestructuringProperty struct {
	Key      string     // Property key
	Value    string     // Variable name (can be different from key)
	Default  Expression // Default value (optional)
	TypeInfo types.TypeInfo
	Pos      lexer.Position
}

func (odp *ObjectDestructuringProperty) Type() types.TypeInfo {
	return odp.TypeInfo
}

func (odp *ObjectDestructuringProperty) Position() lexer.Position {
	return odp.Pos
}

func (odp *ObjectDestructuringProperty) String() string {
	result := odp.Key
	if odp.Value != odp.Key {
		result += ": " + odp.Value
	}
	if odp.Default != nil {
		result += " = " + odp.Default.String()
	}
	return result
}

func (odp *ObjectDestructuringProperty) destructuringElementNode() {}
