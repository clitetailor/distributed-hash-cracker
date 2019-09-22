package lib

// DataTransfer contains transfer data.
type DataTransfer struct {
	Type   string `json:"type"`
	Start  []rune `json:"start"`
	End    []rune `json:"end"`
	Code   string `json:"code"`
	Result string `json:"result"`
}
