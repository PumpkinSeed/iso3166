package generator

var countryToAlphaTmpl = `
package iso3166

// Code generated by country-states 0.0.1 (https://github.com/hajnalandor/country-states). DO NOT EDIT.
// This file is meant to be re-generated in place and/or deleted at any time.

var CountryToAlpha2 = map[string]string {
	{{ range $country, $alpha_codes := . }}
		"{{$country}}":"{{$alpha_codes.Alpha2}}",
	{{- end}}
}

var CountryToAlpha3 = map[string]string {
	{{ range $country, $alpha_codes := . }}
		"{{$country}}":"{{$alpha_codes.Alpha3}}",
	{{- end}}
}

var Alpha2ToCountry = map[string]string {
	{{ range $country, $alpha_codes := . }}
		"{{$alpha_codes.Alpha2}}":"{{$country}}",
	{{- end}}
}

var Alpha3ToCountry = map[string]string {
	{{ range $country, $alpha_codes := . }}
		"{{$alpha_codes.Alpha3}}":"{{$country}}",
	{{- end}}
}
`

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
		Subdivisions  []Subdivision
	}

	type Subdivision struct {
		Name             string
		Code 			 string
		Type             string
		LocalName        string
		ParentCode		 string
	}

	var Countries = [{{len .}}]Country {
		{{ range $index, $country := . }}
		{
			Alpha2:			"{{$country.Alpha2}}",
			Alpha3:			"{{$country.Alpha3}}",
			Name:			"{{$country.Name}}",
			OfficialName:	"{{$country.OfficialName}}",
			CommonName:		"{{$country.CommonName}}",
			Numeric:		"{{$country.Numeric}}",
			Subdivisions:   []Subdivision {
				{{ range $subdivision := $country.Subdivisions }}
				{
					Name:"{{$subdivision.Name}}",
					Code:"{{$subdivision.Code}}",
					Type:"{{$subdivision.Type}}",
					LocalName:"{{$subdivision.LocalName}}",
					ParentCode:"{{$subdivision.Parent}}",
				},
				{{- end}}
			},
		},
		{{- end}}
	}
`
