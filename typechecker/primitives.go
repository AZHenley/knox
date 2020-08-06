package typechecker

type primitives struct {
	typeVOID       *typeObj
	typeBOOL       *typeObj
	typeINT        *typeObj
	typeINTLITERAL *typeObj
	typeI32        *typeObj
	typeFLOAT      *typeObj
	typeSTRING     *typeObj
	typeNIL        *typeObj
	typeLIST       *typeObj
	typeMAP        *typeObj
	typeADDRESS    *typeObj
}

func (p *primitives) Init() {
	p.typeVOID.Init("void", false, false, false, true, false, false, false, false, false, false, false)
	p.typeBOOL.Init("bool", false, false, false, true, false, false, false, false, false, false, false)
	p.typeINT.Init("int", false, true, false, true, false, false, false, false, false, false, false)
	p.typeINTLITERAL.Init("INT_LITERAL", true, true, false, true, false, false, false, false, false, false, false)
	p.typeI32.Init("i32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeFLOAT.Init("float", false, true, false, true, false, false, false, false, false, false, false)
	p.typeSTRING.Init("string", false, false, false, true, false, false, false, false, false, false, false)
	p.typeNIL.Init("nil", true, false, false, true, false, false, false, false, false, false, false)
}

func (p *primitives) IsPrimitiveType(literal string) bool {
	return literal == "bool" || literal == "string" || literal == "int" || literal == "float" || literal == "i32"
}

func (p *primitives) IsNumberType(literal string) bool {
	return literal == "int" || literal == "float" || literal == "i32"
}
