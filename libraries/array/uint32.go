package array

type ArrUint32 uint32

func (s ArrUint32) InArray(val uint32, array []uint32) (exists bool, index int) {
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

func (s ArrUint32) Remove(array []uint32, value uint32) []uint32 {
	isExist, index := s.InArray(value, array)
	if isExist {
		array = append(array[:index], array[(index+1):]...)
	}

	return array
}

func (s ArrUint32) RemoveByIndex(array []uint32, index int) []uint32 {
	return append(array[:index], array[(index+1):]...)
}
