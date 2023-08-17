package find

// Signed 有符合整数
type Signed interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

// Unsigned 无符号整数
type Unsigned interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr
}

// Integer 整数
type Integer interface {
	Signed | Unsigned
}

// Float 浮点数
type Float interface {
	~float32 | ~float64
}

// Number 数字
type Number interface {
	Integer | Float
}

// Ordered 数字或字符串
type Ordered interface {
	Number | ~string
}
