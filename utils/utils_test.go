package utils

import "testing"

func TestRandomOneOf(t *testing.T) {
	v := RandomOneOf(1, 2, 3, 4, 5, 6)
	println(v)

	arr := []string{"Vijay", "Sabi", "Shru", "Anan", "Priya", "Rob", "Abc"}
	for i := 0; i < 10; i++ {
		s := RandomOneOf(arr...)
		println(s)
	}
}
