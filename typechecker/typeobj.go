package typechecker

// Internal representation of a type.
type typeObj struct {
	fullName    string    // Name of this type (and all inner types)
	name        string    // Name of this outer type
	isMutable   bool      // Is this mutable or immutable TODO: Unused right now.
	isReference bool      // Is this a reference type (should be opposite of isPrimitive?)
	isLiteral   bool      // Is is this a literal value
	isNumber    bool      // Is this a number type (i32, f64, etc.)
	isFunction  bool      // Is this a function
	isPrimitive bool      // Is this type a primitive (int, float, string, rune, byte, bool)
	isContainer bool      // Is this type a container (list, map, address, etc.)
	isList      bool      // Is this a builtin list structure
	isMap       bool      // Is this a builtin map structure
	isMulti     bool      // Is this a set of types (used for multiple return)
	isClass     bool      // Is this a user-defined class
	isEnum      bool      // Is this an enum
	isTypedef   bool      // Is this a typedef
	inner       []typeObj // Inner types. TODO: Make this a slice of pointers of typeObj.
}

func (t *typeObj) Init(name string, literal bool, number bool, function bool, primitive bool, container bool, listt bool, mapt bool, multi bool, classt bool, enumt bool, typedeft bool) {
	t = &typeObj{}
	t.name = name
	t.fullName = name
	t.isLiteral = literal
	t.isNumber = number
	t.isFunction = function
	t.isPrimitive = primitive
	t.isContainer = container
	t.isList = listt
	t.isMap = mapt
	t.isMulti = multi
	t.isClass = classt
	t.isEnum = enumt
	t.isTypedef = typedeft
}
