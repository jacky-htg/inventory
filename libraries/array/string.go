package array

// ArrString type of array string
type ArrString string

// InArray checking, return bool and index
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

// Remove array by value
func (s ArrString) Remove(array []string, value string) []string {
	isExist, index := s.InArray(value, array)
	if isExist {
		array = append(array[:index], array[(index+1):]...)
	}

	return array
}
