package users

type Report struct {
	Age   int    `json:"age"`
	Dexas []Dexa `json:"dexas"`
}
