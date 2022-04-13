package funcutil

func IntMin(x, y int) int {
	if x < y {
		return x
	}
	return y
}

func UintMin(x, y uint) uint {
	if x < y {
		return x
	}
	return y
}

func Int64Min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}

func Uint64Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}
