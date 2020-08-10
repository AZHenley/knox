package emitter

func initDataTypes() map[string]string {
	m := map[string]string{
		"bool":   "bool",
		"int":    "int",
		"float":  "float",
		"i8":     "int8_t",
		"i16":    "int16_t",
		"i32":    "int32_t",
		"i64":    "int64_t",
		"u8":     "uint8_t",
		"u16":    "uint16_t",
		"u32":    "uint32_t",
		"u64":    "uint64_t",
		"f32":    "float",
		"f64":    "double",
		"string": "const char *"}
	return m
}
