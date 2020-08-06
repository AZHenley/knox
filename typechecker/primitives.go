package typechecker

import "fmt"

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
	p.typeVOID = createTypeObj("void", false, false, false, true, false, false, false, false, false, false, false)
	fmt.Println("%%%")
	fmt.Println(p.typeVOID)
	p.typeBOOL = createTypeObj("bool", false, false, false, true, false, false, false, false, false, false, false)
	p.typeINT = createTypeObj("int", false, true, false, true, false, false, false, false, false, false, false)
	p.typeINTLITERAL = createTypeObj("INT_LITERAL", true, true, false, true, false, false, false, false, false, false, false)
	p.typeI32 = createTypeObj("i32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeFLOAT = createTypeObj("float", false, true, false, true, false, false, false, false, false, false, false)
	p.typeSTRING = createTypeObj("string", false, false, false, true, false, false, false, false, false, false, false)
	p.typeNIL = createTypeObj("nil", true, false, false, true, false, false, false, false, false, false, false)
}

func (p *primitives) IsPrimitiveType(literal string) bool {
	return literal == "bool" || literal == "string" || literal == "int" || literal == "float" || literal == "i32"
}

func (p *primitives) IsNumberType(literal string) bool {
	return literal == "int" || literal == "float" || literal == "i32"
}
