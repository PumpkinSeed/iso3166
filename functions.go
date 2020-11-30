package iso3166

import (
	"errors"
	"strings"
)

var (
	ErrInvalidCountryName   = errors.New("invalid country name")
	ErrInvalidCountryAlpha2 = errors.New("invalid country alpha2")
	ErrInvalidSubDivName    = errors.New("invalid state name")
	ErrInvalidSubDivCode    = errors.New("invalid state code")
)

// CountryNameToAlpha2 returns the countries alpha2 representation
func CountryNameToAlpha2(name string) (string, error) {
	if alpha2, ok := CountryToAlpha2[name]; ok {
		return alpha2, nil
	} else {
		name = strings.ToUpper(name)
		for _, country := range CountryStates {
			if strings.ToUpper(country.Name) == name || strings.ToUpper(country.OfficialName) == name || strings.ToUpper(country.CommonName) == name {
				return country.Alpha2, nil
			}
		}
	}
	return "", ErrInvalidCountryName
}

// ValidCountryName is validate whether the country name is a valid name,
// looking in the alpha2 representation, country's name, country's official
// name and country's common name
func ValidCountryName(name string) bool {
	if _, ok := CountryToAlpha2[name]; ok {
		return true
	} else {
		name = strings.ToUpper(name)
		for _, country := range CountryStates {
			if strings.ToUpper(country.Name) == name || strings.ToUpper(country.OfficialName) == name || strings.ToUpper(country.CommonName) == name {
				return true
			}
		}
	}

	return false
}

// CountryAlpha2ToName returns the country's name from alpha2 representation
func CountryAlpha2ToName(alpha2 string) (string, error) {
	alpha2 = strings.ToUpper(alpha2)
	if country, ok := CountryStates[alpha2]; ok {
		return country.Name, nil
	}

	return "", ErrInvalidCountryAlpha2
}

// CountryAlpha2ToOfficalName returns the country's offical name from alpha2 representation
func CountryAlpha2ToOfficialName(alpha2 string) (string, error) {
	alpha2 = strings.ToUpper(alpha2)
	if country, ok := CountryStates[alpha2]; ok {
		return country.OfficialName, nil
	}

	return "", ErrInvalidCountryAlpha2
}

// CountryAlpha2ToCommonName returns the country's common name from alpha2 representation
func CountryAlpha2ToCommonName(alpha2 string) (string, error) {
	alpha2 = strings.ToUpper(alpha2)
	if country, ok := CountryStates[alpha2]; ok {
		return country.CommonName, nil
	}

	return "", ErrInvalidCountryAlpha2
}

// ValidCountryAlpha2 validates the alpha2 representation
func ValidCountryAlpha2(alpha2 string) bool {
	alpha2 = strings.ToUpper(alpha2)
	_, ok := CountryStates[alpha2]

	return ok
}

// SubdivisionNameToCode returns the subdivision's code from it's name
func SubdivisionNameToCode(countryAlpha2, subDivName string) (string, error) {
	countryAlpha2 = strings.ToUpper(countryAlpha2)
	subDivName = strings.ToUpper(subDivName)
	if !ValidCountryAlpha2(countryAlpha2) {
		var err error
		countryAlpha2, err = CountryNameToAlpha2(countryAlpha2)
		if err != nil {
			return "", err
		}
	}
	if c, ok := CountryStates[countryAlpha2].SubDivNameToCode[subDivName]; ok {
		return c.Code, nil
	}
	for _, subDiv := range CountryStates[countryAlpha2].SubDivNameToCode {
		if codeWrapper, ok := subDiv.SubDivNameToCode[subDivName]; ok {
			return codeWrapper.Code, nil
		}
	}
	for subDivCode, subDivWrapper := range CountryStates[countryAlpha2].SubDivCodeToName {
		if subDivWrapper.Name == subDivName || subDivWrapper.LocalName == subDivName {
			return subDivCode, nil
		}
	}

	return "", ErrInvalidSubDivName
}

// SubdivisionCodeToName returns the subdivison's name from it's code
func SubdivisionCodeToName(countryAlpha2, subDivCode string) (string, error) {
	countryAlpha2 = strings.ToUpper(countryAlpha2)
	if !ValidCountryAlpha2(countryAlpha2) {
		var err error
		countryAlpha2, err = CountryNameToAlpha2(countryAlpha2)
		if err != nil {
			return "", err
		}
	}
	if c, ok := CountryStates[countryAlpha2].SubDivCodeToName[subDivCode]; ok {
		return c.Name, nil
	}
	for _, subDiv := range CountryStates[countryAlpha2].SubDivCodeToName {
		if codeWrapper, ok := subDiv.SubDivCodeToName[subDivCode]; ok {
			return codeWrapper.Name, nil
		}
	}

	return "", ErrInvalidSubDivCode
}

// ValidSubdivisonCode validate the code of the subdivison
func ValidSubdivisionCode(countryAlpha2, subDivCode string) bool {
	countryAlpha2 = strings.ToUpper(countryAlpha2)
	if !ValidCountryAlpha2(countryAlpha2) {
		var err error
		countryAlpha2, err = CountryNameToAlpha2(countryAlpha2)
		if err != nil {
			return false
		}
	}
	if _, ok := CountryStates[countryAlpha2].SubDivCodeToName[subDivCode]; ok {
		return true
	}
	for _, subDiv := range CountryStates[countryAlpha2].SubDivCodeToName {
		if _, ok := subDiv.SubDivCodeToName[subDivCode]; ok {
			return true
		}
	}
	return false
}
