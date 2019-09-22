package array

type ArrString string

func (s ArrString) InArray(val string, array []string) (exists bool, index int) {
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

func (s ArrString) Remove(array []string, value string) []string {
	isExist, index := s.InArray(value, array)
	if isExist {
		array = append(array[:index], array[(index+1):]...)
	}

	return array
}
