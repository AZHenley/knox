package typechecker

type primitives struct {
	typeVOID         *typeObj
	typeBOOL         *typeObj
	typeINT          *typeObj
	typeINTLITERAL   *typeObj
	typeFLOATLITERAL *typeObj
	typeI32          *typeObj
	typeFLOAT        *typeObj
	typeSTRING       *typeObj
	typeNIL          *typeObj
	typeLIST         *typeObj
	typeMAP          *typeObj
	typeADDRESS      *typeObj
}

func (p *primitives) Init() {
	p.typeVOID = createTypeObj("void", false, false, false, true, false, false, false, false, false, false, false)
	p.typeBOOL = createTypeObj("bool", false, false, false, true, false, false, false, false, false, false, false)
	p.typeINT = createTypeObj("int", false, true, false, true, false, false, false, false, false, false, false)
	p.typeINTLITERAL = createTypeObj("INT_LITERAL", true, true, false, true, false, false, false, false, false, false, false)
	p.typeFLOATLITERAL = createTypeObj("FLOAT_LITERAL", true, true, false, true, false, false, false, false, false, false, false)
	p.typeI32 = createTypeObj("i32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeFLOAT = createTypeObj("float", false, true, false, true, false, false, false, false, false, false, false)
	p.typeSTRING = createTypeObj("string", false, false, false, true, false, false, false, false, false, false, false)
	p.typeNIL = createTypeObj("nil", true, false, false, true, false, false, false, false, false, false, false)
}

func (p *primitives) IsPrimitiveType(text string) bool {
	return text == "bool" || text == "string" || text == "int" || text == "float" || text == "i32"
}

func (p *primitives) IsNumberType(text string) bool {
	return text == "int" || text == "float" || text == "i32"
}
