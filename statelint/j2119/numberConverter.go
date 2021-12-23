package j2119

func GetFloat(unk interface{}) float64 {
	if f, ok := unk.(float64); ok {
		return f
	}

	if i, ok := unk.(int); ok {
		return float64(i)
	}

	panic("type is not float or int")
}
