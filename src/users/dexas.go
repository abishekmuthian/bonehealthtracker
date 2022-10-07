package users

type Dexa struct {
	Id     int     `json: id`
	Year   string  `json: year`
	Organs []Organ `json: organs`
}
