package users

type Organ struct {
	Id          int     `json: id`
	Direction   string  `json:direction`
	Site        string  `json: site`
	Bmd         float64 `json:bmd`
	TScore      float64 `json:tScore`
	ZScore      float64 `json:zScore`
	BeginOffset float64 `json:beginOffset`
	EndOffset   float64 `json:endOffset`
}
