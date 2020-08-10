package typechecker

type primitives struct {
	typeVOID         *typeObj
	typeBOOL         *typeObj
	typeINT          *typeObj
	typeINTLITERAL   *typeObj
	typeFLOATLITERAL *typeObj
	typeI8           *typeObj
	typeI16          *typeObj
	typeI32          *typeObj
	typeI64          *typeObj
	typeU8           *typeObj
	typeU16          *typeObj
	typeU32          *typeObj
	typeU64          *typeObj
	typeFLOAT        *typeObj
	typeF32          *typeObj
	typeF64          *typeObj
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
	p.typeI8 = createTypeObj("i8", false, true, false, true, false, false, false, false, false, false, false)
	p.typeI16 = createTypeObj("i16", false, true, false, true, false, false, false, false, false, false, false)
	p.typeI32 = createTypeObj("i32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeI64 = createTypeObj("i64", false, true, false, true, false, false, false, false, false, false, false)
	p.typeU8 = createTypeObj("u8", false, true, false, true, false, false, false, false, false, false, false)
	p.typeU16 = createTypeObj("u16", false, true, false, true, false, false, false, false, false, false, false)
	p.typeU32 = createTypeObj("u32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeU64 = createTypeObj("u64", false, true, false, true, false, false, false, false, false, false, false)
	p.typeFLOAT = createTypeObj("float", false, true, false, true, false, false, false, false, false, false, false)
	p.typeF32 = createTypeObj("f32", false, true, false, true, false, false, false, false, false, false, false)
	p.typeF64 = createTypeObj("f64", false, true, false, true, false, false, false, false, false, false, false)
	p.typeSTRING = createTypeObj("string", false, false, false, true, false, false, false, false, false, false, false)
	p.typeNIL = createTypeObj("nil", true, false, false, true, false, false, false, false, false, false, false)
}

func (p *primitives) IsPrimitiveType(text string) bool {
	return text == "bool" || text == "string" || text == "int" || text == "float" || text == "i8" || text == "i16" || text == "i32" || text == "i64" || text == "u8" || text == "u16" || text == "u32" || text == "u64" || text == "f32" || text == "f64"
}

func (p *primitives) IsNumberType(text string) bool {
	return text == "int" || text == "float" || text == "i8" || text == "i16" || text == "i32" || text == "i64" || text == "u8" || text == "u16" || text == "u32" || text == "u64" || text == "f32" || text == "f64"
}
