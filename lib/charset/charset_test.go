package charset

import "testing"

func TestInc(t *testing.T) {
	arr := []rune("hello world")
	arr = IncRuneArr(arr)

	if string(arr) != "iello world" {
		t.Error("Expected 'hello world' got", string(arr))
	}
}

func TestBoundary(t *testing.T) {
	arr := []rune("99")
	arr = IncRuneArr(arr)

	if string(arr) != string([]rune{ rune(0), rune(0), 'a' }) {
		t.Error("Expected '00a' got", string(arr))
	}
}

func TestBigInt(t *testing.T) {
	arr := []rune("hello world")

	bigInt := RuneArrToBigInt(arr)
	str := string(BigIntToRuneArr(bigInt))
	if str != "hello world" {
		t.Error("Expected 'hello world' got", str)
	}
}
