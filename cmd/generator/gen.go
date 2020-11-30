package generator

import (
	"bytes"
	"go/format"
	"os"
	"strings"
	"text/template"
)

type CountryWrapper struct {
	Countries []Country `json:"3166-1"`
}

type Country struct {
	Alpha2           string                            `json:"alpha_2"`
	Alpha3           string                            `json:"alpha_3"`
	Name             string                            `json:"name"`
	OfficialName     string                            `json:"official_name"`
	CommonName       string                            `json:"common_name"`
	Numeric          string                            `json:"numeric"`
	SubDivCodeToName map[string]SubDivisionNameWrapper `json:"sub_div_code_to_name"`
	SubDivNameToCode map[string]SubDivisionCodeWrapper `json:"sub_div_name_to_code"`
}

type SubDivisionWrapper struct {
	SubDivisions []SubDivision `json:"3166-2"`
}

type SubDivision struct {
	Name         string `json:"name"`
	LocalName    string `json:"local_name"`
	LanguageCode string `json:"language_code"`
	Code         string `json:"code"`
	Parent       string `json:"parent"`
	Type         string `json:"type"`
}

type SubDivisionNameWrapper struct {
	Name             string                            `json:"name"`
	LocalName        string                            `json:"local_name"`
	LanguageCode     string                            `json:"language_code"`
	Type             string                            `json:"type"`
	SubDivCodeToName map[string]SubDivisionNameWrapper `json:"subdivision"`
}

type SubDivisionCodeWrapper struct {
	Code             string                            `json:"code"`
	SubDivNameToCode map[string]SubDivisionCodeWrapper `json:"subdivision"`
}

func SlicesToMap(cw CountryWrapper, sw SubDivisionWrapper) map[string]Country {
	countryCodes := getAlpha2CountryCodes(cw)
	countryMap := make(map[string]Country)

	for _, countryCode := range countryCodes {
		parentStructure := getParentStructure(countryCode, sw)
		country := getCountry(countryCode, cw)
		country.SubDivCodeToName, country.SubDivNameToCode = buildSubDiv(parentStructure, countryCode, sw)
		countryMap[countryCode] = country
	}

	return countryMap
}

func buildSubDiv(parent map[string][]string,
	countryCode string, wrapper SubDivisionWrapper) (map[string]SubDivisionNameWrapper, map[string]SubDivisionCodeWrapper) {
	parentSubDivCodeMap := make(map[string]SubDivisionNameWrapper)
	parentSubDivNameMap := make(map[string]SubDivisionCodeWrapper)
	for subCode, parent := range parent {
		subSubDivCode, subSubDivName := buildSubSubDiv(parent, countryCode, wrapper)

		subDiv := getSubDivision(countryCode, subCode, wrapper)
		parentSubDivCodeMap[subCode] = SubDivisionNameWrapper{
			Name:             subDiv.Name,
			LocalName:        subDiv.LocalName,
			LanguageCode:     subDiv.LanguageCode,
			Type:             subDiv.Type,
			SubDivCodeToName: subSubDivCode,
		}
		parentSubDivNameMap[strings.ToUpper(subDiv.Name)] = SubDivisionCodeWrapper{
			Code:             strings.Split(subDiv.Code, "-")[1],
			SubDivNameToCode: subSubDivName,
		}
	}

	return parentSubDivCodeMap, parentSubDivNameMap
}

func buildSubSubDiv(parent []string, countryCode string, wrapper SubDivisionWrapper) (map[string]SubDivisionNameWrapper, map[string]SubDivisionCodeWrapper) {
	tmpSubDivCode := make(map[string]SubDivisionNameWrapper)
	tmpSubDivName := make(map[string]SubDivisionCodeWrapper)
	for _, subCode := range parent {
		subDiv := getSubDivision(countryCode, subCode, wrapper)
		tmpSubDivCode[subCode] = SubDivisionNameWrapper{
			Name:         subDiv.Name,
			LocalName:    subDiv.LocalName,
			LanguageCode: subDiv.LanguageCode,
			Type:         subDiv.Type,
		}
		tmpSubDivName[strings.ToUpper(subDiv.Name)] = SubDivisionCodeWrapper{
			Code: strings.Split(subDiv.Code, "-")[1],
		}
	}

	return tmpSubDivCode, tmpSubDivName
}

func getCountryNameToAlpha2Map(cw CountryWrapper) map[string]string {
	cNameToAlpha2Map := make(map[string]string)
	for _, country := range cw.Countries {
		cNameToAlpha2Map[country.Name] = country.Alpha2
	}

	return cNameToAlpha2Map
}

func getCountryNameToAlpha3Map(cw CountryWrapper) map[string]string {
	cNameToAlpha3Map := make(map[string]string)
	for _, country := range cw.Countries {
		cNameToAlpha3Map[country.Name] = country.Alpha3
	}

	return cNameToAlpha3Map
}

func getCountry(countryCode string, cw CountryWrapper) Country {
	for _, country := range cw.Countries {
		if country.Alpha2 == countryCode {
			return country
		}
	}
	panic("country not found")
}

func getSubDivision(countryCode, subCode string, sw SubDivisionWrapper) SubDivision {
	for _, subDivision := range sw.SubDivisions {
		if subDivision.Code == countryCode+"-"+subCode {
			return subDivision
		}
	}
	panic("subdiv not found" + countryCode + "subcode:" + subCode)
}

func getParentStructure(countryCode string, sw SubDivisionWrapper) map[string][]string {
	parents := make(map[string][]string)
	for _, sd := range sw.SubDivisions {
		if code := strings.Split(sd.Code, "-"); code[0] == countryCode {
			if sd.Parent != "" {
				parentCode := sd.Parent
				if splitted := strings.Split(sd.Parent, "-"); len(splitted) == 2 {
					parentCode = splitted[1]
				}
				parents[parentCode] = append(parents[parentCode], code[1])
			} else if _, ok := parents[code[1]]; !ok {
				parents[code[1]] = make([]string, 0)
			}
		}
	}

	return parents
}

func getAlpha2CountryCodes(w CountryWrapper) []string {
	countries := make([]string, 0)
	for _, c := range w.Countries {
		countries = append(countries, c.Alpha2)
	}

	return countries
}

var countrySubDivtmpl = `
package iso3166

// Code generated by country-states 0.0.1 (https://github.com/hajnalandor/country-states). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

	type Country struct {
		Alpha2       string 
		Alpha3       string 
		Name         string 
		OfficialName string 
		CommonName   string 
		Numeric      string 
		SubDivCodeToName map[string]SubDivisionNameWrapper
		SubDivNameToCode map[string]SubDivisionCodeWrapper
	}

	type SubDivisionNameWrapper struct {
		Name             string                            
		Type             string
		LocalName        string
		LanguageCode     string 
		SubDivCodeToName map[string]SubDivisionNameWrapper
	}

		type SubDivisionCodeWrapper struct {
		Code             string                            
		SubDivNameToCode map[string]SubDivisionCodeWrapper 
	}

	var CountryStates = map[string]Country {
	{{ range $key, $value := . }}
		"{{$key}}": {
			Alpha2:		"{{$value.Alpha2}}",
			Alpha3:		"{{$value.Alpha3}}",
			Name:			"{{$value.Name}}",
			OfficialName:	"{{$value.OfficialName}}",
			CommonName:	"{{$value.CommonName}}",
			Numeric:		"{{$value.Numeric}}",
			{{ if ne (len $value.SubDivCodeToName) 0}}
			SubDivCodeToName:       map[string]SubDivisionNameWrapper{
							{{ range $sk1, $sk2 := $value.SubDivCodeToName}}
								"{{$sk1}}": {
									Name:   "{{$sk2.Name}}",
									LocalName: "{{$sk2.LocalName}}",
									LanguageCode: "{{$sk2.LanguageCode}}",
									Type:   "{{$sk2.Type}}",
									{{ if ne (len $sk2.SubDivCodeToName) 0}}
									SubDivCodeToName:  map[string]SubDivisionNameWrapper{
										{{ range $childKey, $childVal := $sk2.SubDivCodeToName}}
											"{{$childKey}}": {
												Name:   "{{$childVal.Name}}",
												LocalName: "{{$childVal.LocalName}}",
												LanguageCode: "{{$childVal.LanguageCode}}",
												Type:   "{{$childVal.Type}}",
											},
										{{- end}}
										},
									{{- end}}
								},
							{{- end}}
							},
			{{- end}}
			{{ if ne (len $value.SubDivNameToCode) 0}}
			SubDivNameToCode:       map[string]SubDivisionCodeWrapper{
							{{ range $sk1, $sk2 := $value.SubDivNameToCode}}
								"{{$sk1}}": {
									Code:   "{{$sk2.Code}}",
									{{ if ne (len $sk2.SubDivNameToCode) 0}}
									SubDivNameToCode:  map[string]SubDivisionCodeWrapper{
										{{ range $childKey, $childVal := $sk2.SubDivNameToCode}}
											"{{$childKey}}": {
												Code:   "{{$childVal.Code}}",
											},
										{{- end}}
										},
									{{- end}}
								},
							{{- end}}
							},
			{{- end}}
		},
	{{ end }}

	}
`

var countryToAlpha2Tmpl = `
package iso3166

// Code generated by country-states 0.0.1 (https://github.com/hajnalandor/country-states). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

var CountryToAlpha2 = map[string]string {
	{{ range $key, $value := . }}
		"{{$key}}":"{{$value}}",
	{{- end}}
}
`

var countryToAlpha3Tmpl = `
package iso3166

// Code generated by country-states 0.0.1 (https://github.com/hajnalandor/country-states). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

var CountryToAlpha3 = map[string]string {
	{{ range $key, $value := . }}
		"{{$key}}":"{{$value}}",
	{{- end}}
}
`

func GenerateCountryStates(data map[string]Country) {
	tmpl, err := template.New("country-state-generator").Parse(countrySubDivtmpl)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../../country-states.go")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = f.Write(formatted)
	if err != nil {
		panic(err)
	}
}

func GenerateCountryToAlpha2(data map[string]string) {
	tmpl, err := template.New("country-to-alpha2-generator").Parse(countryToAlpha2Tmpl)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../../country-alpha2.go")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = f.Write(formatted)
	if err != nil {
		panic(err)
	}
}

func GenerateCountryToAlpha3(data map[string]string) {
	tmpl, err := template.New("country-to-alpha3-generator").Parse(countryToAlpha3Tmpl)
	if err != nil {
		panic(err)
	}

	f, err := os.Create("../../country-alpha3.go")
	if err != nil {
		panic(err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		panic(err)
	}
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		panic(err)
	}
	_, err = f.Write(formatted)
	if err != nil {
		panic(err)
	}
}
