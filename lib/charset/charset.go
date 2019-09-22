package charset

import (
	"crypto/md5"
	"encoding/hex"
	"math/big"
)

// Charset defines a set of common characters to be used.
var Charset = []rune{}

// CharsetMap maps between character and its index position in Charset.
var CharsetMap map[rune]int

// CharsetSize is the size of Charset.
var CharsetSize int

func init() {
	Charset = append(Charset, rune(0))

	for char := 'a'; char <= 'z'; char++ {
		Charset = append(Charset, char)
	}

	Charset = append(Charset, ' ')

	for char := '0'; char <= '9'; char++ {
		Charset = append(Charset, char)
	}

	CharsetMap = make(map[rune]int)
	for i, char := range Charset {
		CharsetMap[char] = i
	}

	CharsetSize = len(Charset)
}

// IncRuneArr increments rune arr.
func IncRuneArr(arr []rune) []rune {
	newArr := make([]rune, len(arr))
	copy(newArr, arr)
	inc := 1

	for i, char := range arr {
		newArr[i] = Charset[(CharsetMap[char]+inc)%CharsetSize]
		inc = (CharsetMap[char] + inc) / CharsetSize

		if inc == 0 {
			break
		}
	}

	if inc != 0 {
		newArr = append(newArr, Charset[rune(inc)])
	}

	return newArr
}

// IsValid checks whether an array is valid.
func IsValid(arr []rune) bool {
	length := len(arr)

	for i := 0; i < length; i++ {
		if arr[i] == rune(0) {
			return false
		}
	}

	return true
}

// Sign helps compare two rune array.
func Sign(a []rune, b []rune) int {
	if len(a) > len(b) {
		return 1
	}

	if len(a) < len(b) {
		return -1
	}

	for i := len(a) - 1; i > -1; i-- {
		if a[i] > b[i] {
			return 1
		}
		if a[i] < b[i] {
			return -1
		}
	}

	return 0
}

// HashString hashes a string into a string and returns.
func HashString(str string) string {
	hash := md5.Sum([]byte(str))
	return hex.EncodeToString(hash[:])
}

// RuneArrToBigInt converts rune array to big.Int and returns.
func RuneArrToBigInt(runeArr []rune) *big.Int {

	accum := big.NewInt(0)
	charsetSize := big.NewInt(int64(CharsetSize))

	runeArr = Reverse(runeArr)

	for _, ch := range runeArr {
		accum.Mul(accum, charsetSize)
		accum.Add(accum, RuneToBigInt(ch))
	}

	return accum
}

// RuneToBigInt converts rune to big.Int.
func RuneToBigInt(r rune) *big.Int {
	n := CharsetMap[r]
	return big.NewInt(int64(n))
}

// BigIntToRuneArr converts big.Int to rune array.
func BigIntToRuneArr(number *big.Int) []rune {
	number = new(big.Int).Set(number)
	runeArr := []rune{}

	charsetSize := big.NewInt(int64(CharsetSize))
	remainder := big.NewInt(0)
	for {
		number.DivMod(number, charsetSize, remainder)

		ch := Charset[int(remainder.Int64())]
		runeArr = append(runeArr, ch)

		if number.Cmp(charsetSize) < 0 {
			number.DivMod(number, charsetSize, remainder)

			ch := Charset[int(remainder.Int64())]
			runeArr = append(runeArr, ch)
			break
		}
	}

	return runeArr
}

// Reverse reverses the rune array.
func Reverse(runeArr []rune) []rune {
	newRuneArr := make([]rune, len(runeArr))

	length := len(runeArr)
	for i := 0; i < length/2; i++ {
		j := len(runeArr) - i - 1
		newRuneArr[i], newRuneArr[j] = runeArr[j], runeArr[i]
	}

	if length%2 != 0 {
		newRuneArr[length/2] = runeArr[length/2]
	}

	return newRuneArr
}

// Range splits range into a number of ranges.
func Range(start []rune, end []rune, count int) [][2][]rune {
	ranges := [][2][]rune{}

	startInt := RuneArrToBigInt(start)
	endInt := RuneArrToBigInt(end)

	diff := new(big.Int).Sub(endInt, startInt)
	chunkSize := new(big.Int).Div(diff, big.NewInt(int64(count)))

	nLoop := count - 1
	for i := 0; i < nLoop; i++ {
		midInt := new(big.Int).Add(startInt, chunkSize)

		ranges = append(ranges, [2][]rune{
			BigIntToRuneArr(startInt),
			BigIntToRuneArr(midInt)})

		startInt = midInt
	}

	ranges = append(ranges, [2][]rune{
		BigIntToRuneArr(startInt),
		BigIntToRuneArr(endInt)})

	return ranges
}
