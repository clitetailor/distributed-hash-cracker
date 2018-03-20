package charset

// Charset defines a set of common characters to be used.
var Charset = []rune{}
// CharsetMap maps between character and its index position in Charset.
var CharsetMap map[rune]int
// CharsetSize is the size of Charset.
var CharsetSize int


func init() {
	Charset = append(Charset, rune(0))

	for char := 'a'; char < 'z'; char++ {
		Charset = append(Charset, char)
	}

	Charset = append(Charset, ' ')

	for char := '0'; char < '9'; char++ {
		Charset = append(Charset, char)
	}

	CharsetMap = make(map[rune]int)
	for i, char := range Charset {
		CharsetMap[char] = i
	}

	CharsetSize = len(Charset)
}

// IncRuneArr helps increment rune array and returns false if contains null character.
func IncRuneArr(arr []rune) ([]rune, bool) {
	newArr := make([]rune, len(arr))
	copy(newArr, arr)
	
	inc := 1
	for i, char := range newArr {
		newArr[i] = Charset[(CharsetMap[char] + inc) % CharsetSize]
		inc = (CharsetMap[char] + inc) / CharsetSize

		if inc == 0 {
			break
		}
	}

	for _, char := range newArr {
		if char == 0 {
			return newArr, false
		}
	}

	return newArr, true
}

// Sign helps compare two rune array.
func Sign(a []rune, b []rune) int {
	if (len(a) > len(b)) {
		return 1
	}

	if (len(a) < len(b)) {
		return -1
	}

	for i := len(a); i > -1; i-- {
		if (a[i] > b[i]) {
			return 1
		}
		if (a[i] < b[i]) {
			return -1
		}
	}

	return 0
}
