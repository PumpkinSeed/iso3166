package generator

import (
	"testing"
)

func TestGenerateCountries(t *testing.T) {
	var cw countryWrapper
	mustParseJSONFile("../../data/iso3166-1.json", &cw)

	csvLines := mustParseCsv("../../data/iso3166-2.csv")
	subDivisions := csvToSubdivision(csvLines)

	countries := addSubdivisions(cw.Countries,subDivisions)
	generateCountries(countries)
}

func TestGenerateCountryMapsFromJSON(t *testing.T) {
	var cw countryWrapper
	mustParseJSONFile("../../data/iso3166-1.json", &cw)

	data := getCountryToAlphaMap(cw)
	generateCountryMap(data)
}
