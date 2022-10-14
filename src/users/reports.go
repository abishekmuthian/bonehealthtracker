package users

type Report struct {
	FirstReportAge  int    `json:"first-report-age"`
	LatestReportAge int    `json:"latest-report-age"`
	Treatment       string `json:"treatment"`
	RaceEthinicity  string `json:"race-ethnicity"`
	Dexas           []Dexa `json:"dexas"`
}
