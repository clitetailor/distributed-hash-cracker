package charset

import "testing"

func TestInc(t *testing.T) {
	arr := []rune("gello world!")
	arr, _ = IncRuneArr(arr)

	if string(arr) != "hello world!" {
		t.Error("Expected 'hello world!' got ", string(arr))
	}
}
