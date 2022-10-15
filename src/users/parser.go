package users

import (
	"encoding/json"
	"math"
	"sort"
	"strconv"
	"strings"

	"github.com/abishekmuthian/bonehealthtracker/src/lib/server/log"
)

type Direction struct {
	text        string
	beginOffset float64
	endOffset   float64
}

type BMD struct {
	text        float64
	beginOffset float64
	endOffset   float64
}

type tScore struct {
	text        float64
	beginOffset float64
	endOffset   float64
}

type zScore struct {
	text        float64
	beginOffset float64
	endOffset   float64
}

// Parse parses the medical data from AWS Comprehend Medical and returns organs
func Parse(dexaData []byte) []Organ {
	var result interface{}
	err := json.Unmarshal(dexaData, &result)

	if err != nil {
		log.Error(log.V{"Parser, Error unmarshaling dexaData": err})
	}

	var organs []Organ
	var directions []Direction
	var bmds []BMD
	var tScores []tScore
	var zScores []zScore

	m := result.(map[string]interface{})

	for k, v := range m {
		switch vv := v.(type) {
		case string:
			// log.Info(log.V{"Parser":k, "is string": vv})
		case int:
			// log.Info(log.V{"Parser":k, "is int": vv})
		case []interface{}:
			// log.Info(log.V{"Parser":k, "is array": vv})
			for _, u := range vv {
				if u.(map[string]interface{})["Type"] == "ANATOMY" {

					if u.(map[string]interface{})["Attribute"] != nil {
						attribute := u.(map[string]interface{})["Attribute"]

						if attribute.(map[string]interface{})["Type"].(string) == "DIRECTION" {

							// log.Info(log.V{"Parser, Direction":attribute.(map[string]interface{})["Text"]})

							text := strings.ToLower(attribute.(map[string]interface{})["Text"].(string))

							if text == "ap" ||
								text == "left" ||
								text == "right" {
								d := Direction{
									text:        attribute.(map[string]interface{})["Text"].(string),
									beginOffset: attribute.(map[string]interface{})["BeginOffset"].(float64),
									endOffset:   attribute.(map[string]interface{})["EndOffset"].(float64),
								}

								directions = append(directions, d)
							}
						}

					}
				}

				if u.(map[string]interface{})["Type"] == "TEST_TREATMENT_PROCEDURE" {

					attribute := u.(map[string]interface{})["Attribute"]

					if attribute.(map[string]interface{})["Type"] == "TEST_VALUE" {

						// log.Info(log.V{"Paser, BMD":attribute.(map[string]interface{})["Text"]})

						text := strings.ToLower(attribute.(map[string]interface{})["Text"].(string))

						if !strings.Contains(text, "low") &&
							!strings.Contains(text, "high") &&
							!strings.Contains(text, "normal") &&
							!strings.Contains(text, "g/cm2") &&
							!strings.Contains(text, "gm/cm2") &&
							!strings.Contains(text, "gm per cm2") &&
							!strings.Contains(text, "g/cm²") {

							if f, err := strconv.ParseFloat(text, 64); err == nil && f < 1.5 {

								b := BMD{
									text:        f,
									beginOffset: attribute.(map[string]interface{})["BeginOffset"].(float64),
									endOffset:   attribute.(map[string]interface{})["EndOffset"].(float64),
								}

								bmds = append(bmds, b)
							} else {
								log.Error(log.V{"Parser, Error converting text to float": err})
							}
						}

					}

				}

				if u.(map[string]interface{})["Category"] == "ANATOMY" {
					if u.(map[string]interface{})["Type"] == "SYSTEM_ORGAN_SITE" {
						log.Info(log.V{"Paser, Organ": u.(map[string]interface{})["Text"]})

						text := strings.ToLower(u.(map[string]interface{})["Text"].(string))

						if strings.Contains(text, "spine") ||
							strings.Contains(text, "hip") ||
							strings.Contains(text, "l1-l4") ||
							strings.Contains(text, "l1 through l4") ||
							strings.Contains(text, "femur") ||
							strings.Contains(text, "neck") ||
							strings.Contains(text, "forearm") {

							o := Organ{
								Site:        u.(map[string]interface{})["Text"].(string),
								BeginOffset: u.(map[string]interface{})["BeginOffset"].(float64),
								EndOffset:   u.(map[string]interface{})["EndOffset"].(float64),
							}

							organs = append(organs, o)

						}
					}
				}

				if u.(map[string]interface{})["Category"] == "TEST_TREATMENT_PROCEDURE" {

					text := strings.ToLower(u.(map[string]interface{})["Text"].(string))

					if strings.Contains(text, "femur total") {
						o := Organ{
							Site:        u.(map[string]interface{})["Text"].(string),
							BeginOffset: u.(map[string]interface{})["BeginOffset"].(float64),
							EndOffset:   u.(map[string]interface{})["EndOffset"].(float64),
						}

						organs = append(organs, o)

						if u.(map[string]interface{})["Attributes"] != nil {
							attributes := u.(map[string]interface{})["Attributes"].([]interface{})

							for _, a := range attributes {
								// log.Info(log.V{"Parser, BMD": a.(map[string]interface{})["Text"]})

								text := strings.ToLower(a.(map[string]interface{})["Text"].(string))

								if !strings.Contains(text, "low") &&
									!strings.Contains(text, "high") &&
									!strings.Contains(text, "normal") &&
									!strings.Contains(text, "g/cm2") &&
									!strings.Contains(text, "gm/cm2") &&
									!strings.Contains(text, "gm per cm2") &&
									!strings.Contains(text, "g/cm²") {

									if f, err := strconv.ParseFloat(text, 64); err == nil && f < 1.5 {

										b := BMD{
											text:        f,
											beginOffset: a.(map[string]interface{})["BeginOffset"].(float64),
											endOffset:   a.(map[string]interface{})["EndOffset"].(float64),
										}

										bmds = append(bmds, b)
									} else {
										log.Error(log.V{"Parser, Error converting text to float": err})
									}
								}
							}
						}

					}

					if strings.Contains(text, "bmd") || strings.Contains(text, "bone mineral density") {

						if u.(map[string]interface{})["Attributes"] != nil {

							attributes := u.(map[string]interface{})["Attributes"].([]interface{})

							for _, a := range attributes {

								if a.(map[string]interface{})["Type"] == "TEST_VALUE" {
									// log.Info(log.V{"Parser, BMD":a.(map[string]interface{})["Text"]})

									text := a.(map[string]interface{})["Text"].(string)

									if !strings.Contains(text, "low") &&
										!strings.Contains(text, "high") &&
										!strings.Contains(text, "normal") &&
										!strings.Contains(text, "g/cm2") &&
										!strings.Contains(text, "gm/cm2") &&
										!strings.Contains(text, "gm per cm2") &&
										!strings.Contains(text, "g/cm²") {

										if f, err := strconv.ParseFloat(text, 64); err == nil && f < 1.5 {

											b := BMD{
												text:        f,
												beginOffset: a.(map[string]interface{})["BeginOffset"].(float64),
												endOffset:   a.(map[string]interface{})["EndOffset"].(float64),
											}

											bmds = append(bmds, b)
										} else {
											log.Error(log.V{"Parser, Error converting text to float": err})
										}
									}
								}

							}
						}
					}

					if u.(map[string]interface{})["Text"] == "T-score" {

						if u.(map[string]interface{})["Attributes"] != nil {
							attributes := u.(map[string]interface{})["Attributes"].([]interface{})

							for _, a := range attributes {
								// log.Info(log.V{"Parser, T-score":a.(map[string]interface{})["Text"]})

								text := strings.ToLower(a.(map[string]interface{})["Text"].(string))

								if !strings.Contains(text, "low") &&
									!strings.Contains(text, "high") &&
									!strings.Contains(text, "normal") &&
									!strings.Contains(text, "g/cm2") &&
									!strings.Contains(text, "gm/cm2") &&
									!strings.Contains(text, "gm per cm2") &&
									!strings.Contains(text, "g/cm²") {

									if f, err := strconv.ParseFloat(text, 64); err == nil && (f >= -5.0 && f <= 2.5) {

										t := tScore{
											text:        f,
											beginOffset: a.(map[string]interface{})["BeginOffset"].(float64),
											endOffset:   a.(map[string]interface{})["EndOffset"].(float64),
										}
										tScores = append(tScores, t)
									} else {
										log.Error(log.V{"Parser, Error converting text to float": err})
									}

								}

							}
						}
					}

					if u.(map[string]interface{})["Text"] == "Z-score" {

						if u.(map[string]interface{})["Attributes"] != nil {
							attributes := u.(map[string]interface{})["Attributes"].([]interface{})

							for _, a := range attributes {
								// log.Info(log.V{"Parser, Z-score":a.(map[string]interface{})["Text"]})

								text := strings.ToLower(a.(map[string]interface{})["Text"].(string))

								if !strings.Contains(text, "low") &&
									!strings.Contains(text, "high") &&
									!strings.Contains(text, "normal") &&
									!strings.Contains(text, "g/cm2") &&
									!strings.Contains(text, "gm/cm2") &&
									!strings.Contains(text, "gm per cm2") &&
									!strings.Contains(text, "g/cm²") {

									if f, err := strconv.ParseFloat(text, 64); err == nil && (f >= -2.5 && f <= 2.5) {

										z := zScore{
											text:        f,
											beginOffset: a.(map[string]interface{})["BeginOffset"].(float64),
											endOffset:   a.(map[string]interface{})["EndOffset"].(float64),
										}
										zScores = append(zScores, z)
									} else {
										log.Error(log.V{"Parser, Error converting text to float": err})
									}
								}
							}
						}
					}
				}

			}
		default:
			log.Info(log.V{"Parser, Don't know how to handle type": k})
		}

	}

	/* 	log.Info(log.V{"Parser, Organs": organs})
	   	log.Info(log.V{"Parser, Directions": directions})
	   	log.Info(log.V{"Parser, BMDs": bmds})
	   	log.Info(log.V{"Parser, T-Scores": tScores})
	   	log.Info(log.V{"Parser, Z-Scores": zScores}) */

	return setOrganValues(organs, directions, tScores, zScores, bmds)
}

func setOrganValues(organs []Organ, directions []Direction, tScores []tScore, zScores []zScore, bmds []BMD) []Organ {

	var organBeginOffsets, organEndOffsets []float64
	var tempOrgans []Organ

	organBeginOffsets = []float64{}
	organEndOffsets = []float64{}

	for _, organ := range organs {

		organBeginOffsets = append(organBeginOffsets, organ.BeginOffset)
		organEndOffsets = append(organEndOffsets, organ.EndOffset)

	}

	for _, direction := range directions {

		closestBeginOffsetIndex := findClosestElementIndex(organBeginOffsets, 1, direction.endOffset)
		closestEndOffsetIndex := findClosestElementIndex(organEndOffsets, 1, direction.beginOffset)

		if direction.endOffset-organBeginOffsets[closestBeginOffsetIndex] > organEndOffsets[closestEndOffsetIndex]-direction.beginOffset {
			organs[closestEndOffsetIndex].Direction = direction.text
		} else {
			organs[closestBeginOffsetIndex].Direction = direction.text
		}

	}

	for i, organ := range organs {

		if i != len(organs)-1 && (organ.EndOffset+1 == organs[i+1].BeginOffset) {
			if organ.Direction != "" {
				organs[i].Site = organ.Site + " " + organs[i+1].Site

				if !strings.Contains(strings.ToLower(organs[i+1].Site), "forearm") {
					organs[i].EndOffset = organs[i+1].EndOffset
					organs[i+1].Direction = "Remove"
				}
			} else {
				organs[i+1].Site = organ.Site + " " + organs[i+1].Site
				if !strings.Contains(strings.ToLower(organs[i].Site), "forearm") {
					organs[i+1].BeginOffset = organ.BeginOffset
					organs[i].Direction = "Remove"
				}
			}
		} else if organ.Direction == "" && !strings.Contains(strings.ToLower(organs[i].Site), "forearm") &&
			!strings.Contains(strings.ToLower(organs[i].Site), "l1 through l4") {
			organs[i].Direction = "Remove"
		}

	}

	// Remove duplicate organs

	tempOrgans = []Organ{}

	for _, organ := range organs {

		if organ.Direction != "Remove" {
			tempOrgans = append(tempOrgans, organ)
		}
	}

	organs = tempOrgans

	tempOrgans = []Organ{}

	inResult := make(map[string]Organ)

	for _, organ := range organs {
		if _, ok := inResult[organ.Direction+""+organ.Site]; !ok {
			inResult[organ.Direction+organ.Site] = organ
			tempOrgans = append(tempOrgans, organ)
		}
	}

	organs = tempOrgans

	// Find T-scores for the organs
	if len(tScores) != len(organs) {
		organBeginOffsets = []float64{}
		organEndOffsets = []float64{}

		for _, organ := range organs {

			organBeginOffsets = append(organBeginOffsets, organ.BeginOffset)
			organEndOffsets = append(organEndOffsets, organ.EndOffset)

		}

		for _, tScore := range tScores {

			closestBeginOffsetIndex := findClosestElementIndex(organBeginOffsets, 1, tScore.endOffset)
			closestEndOffsetIndex := findClosestElementIndex(organEndOffsets, 1, tScore.beginOffset)

			if math.Abs(organBeginOffsets[closestBeginOffsetIndex]-tScore.endOffset) > math.Abs(tScore.beginOffset-organEndOffsets[closestEndOffsetIndex]) {
				organs[closestEndOffsetIndex].TScore = tScore.text
			} else {
				organs[closestBeginOffsetIndex].TScore = tScore.text
			}

		}
	} else {
		for i, tScore := range tScores {
			organs[i].TScore = tScore.text
		}
	}

	// Find Z-scores for the organs
	if len(zScores) != len(organs) {
		organBeginOffsets = []float64{}
		organEndOffsets = []float64{}

		for _, organ := range organs {

			organBeginOffsets = append(organBeginOffsets, organ.BeginOffset)
			organEndOffsets = append(organEndOffsets, organ.EndOffset)

		}

		for _, zScore := range zScores {

			closestBeginOffsetIndex := findClosestElementIndex(organBeginOffsets, 1, zScore.endOffset)
			closestEndOffsetIndex := findClosestElementIndex(organEndOffsets, 1, zScore.beginOffset)

			if math.Abs(organBeginOffsets[closestBeginOffsetIndex]-zScore.endOffset) > math.Abs(zScore.beginOffset-organEndOffsets[closestEndOffsetIndex]) {
				organs[closestEndOffsetIndex].ZScore = zScore.text
			} else {
				organs[closestBeginOffsetIndex].ZScore = zScore.text
			}

		}
	} else {
		for i, zScore := range zScores {
			organs[i].ZScore = zScore.text
		}
	}

	if len(bmds) != len(organs) {
		organBeginOffsets = []float64{}
		organEndOffsets = []float64{}

		for _, organ := range organs {

			organBeginOffsets = append(organBeginOffsets, organ.BeginOffset)
			organEndOffsets = append(organEndOffsets, organ.EndOffset)

		}

		for _, bmd := range bmds {

			closestBeginOffsetIndex := findClosestElementIndex(organBeginOffsets, 1, bmd.endOffset)
			closestEndOffsetIndex := findClosestElementIndex(organEndOffsets, 1, bmd.beginOffset)

			if bmd.endOffset-organBeginOffsets[closestBeginOffsetIndex] > organEndOffsets[closestEndOffsetIndex]-bmd.beginOffset {
				organs[closestEndOffsetIndex].Bmd = bmd.text
			} else {
				organs[closestBeginOffsetIndex].Bmd = bmd.text
			}

		}
	} else {
		for i, bmd := range bmds {
			organs[i].Bmd = bmd.text
		}
	}

	// Add Id and remove offsets for privacy

	tempOrgans = []Organ{}

	for i, organ := range organs {
		organ.Id = i
		organ.BeginOffset = 0
		organ.EndOffset = 0

		tempOrgans = append(tempOrgans, organ)
	}

	organs = tempOrgans

	log.Info(log.V{"Parser, Organs after setting values": organs})

	return organs
}

func findClosestElementIndex(arr []float64, k int, x float64) int {
	return sort.Search(len(arr)-k, func(i int) bool { return x-arr[i] <= arr[i+k]-x })
}
