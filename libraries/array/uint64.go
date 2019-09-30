package array

// ArrUint64 type of array uint 64
type ArrUint64 uint64

// InArray checking, return bool and index
func (s ArrUint64) InArray(val uint64, array []uint64) (exists bool, index int) {
	exists = false
	index = -1

	for i, s := range array {
		if s == val {
			exists = true
			index = i
			return
		}
	}

	return
}

// Remove array by value
func (s ArrUint64) Remove(array []uint64, value uint64) []uint64 {
	isExist, index := s.InArray(value, array)
	if isExist {
		array = append(array[:index], array[(index+1):]...)
	}

	return array
}

// RemoveByIndex is remove array by index
func (s ArrUint64) RemoveByIndex(array []uint64, index int) []uint64 {
	return append(array[:index], array[(index+1):]...)
}
