package cursor

func MoveDown(index int, length int) int {
	if index+1 >= length {
		return 0
	}
	return index + 1
}

func MoveUp(index int, length int) int {
	if index == 0 {
		return length - 1
	}
	return index - 1
}
