package typechecker

type primitives struct {
	typeVOID    *typeObj
	typeBOOL    *typeObj
	typeINT     *typeObj
	typeFLOAT   *typeObj
	typeSTRING  *typeObj
	typeNIL     *typeObj
	typeLIST    *typeObj
	typeMAP     *typeObj
	typeADDRESS *typeObj
}

func (p *primitives) Init() {
	p.typeVOID = &typeObj{}
	p.typeBOOL = &typeObj{}
	p.typeINT = &typeObj{}
	p.typeFLOAT = &typeObj{}
	p.typeSTRING = &typeObj{}
	p.typeNIL = &typeObj{}
	p.typeBOOL.isPrimitive = true
	p.typeINT.isPrimitive = true
	p.typeFLOAT.isPrimitive = true
	p.typeSTRING.isPrimitive = true
	p.typeVOID.fullName = "void"
	p.typeBOOL.fullName = "bool"
	p.typeINT.fullName = "int"
	p.typeFLOAT.fullName = "float"
	p.typeSTRING.fullName = "string"
	p.typeNIL.fullName = "nil"
}
