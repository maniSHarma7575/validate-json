package jsonparser

import "errors"

// Whitespace constant are declared here
const (
	space          = ' '
	horizontalTab  = '\t'
	newLine        = '\n'
	carriageReturn = '\r'
)

type jsonValidator struct {
	input  string
	cursor int
}

type nilableBool struct {
	value bool
}

var whitespaces = []byte{space, horizontalTab, newLine, carriageReturn}

func (v *jsonValidator) skipWhiteSpace() {
	for v.cursor < len(v.input) && contains(whitespaces, v.input[v.cursor]) {
		v.cursor++
	}
}

func (v *jsonValidator) parseColon() {
	if v.input[v.cursor] != ':' {
		errors.New("invalid json: Expected ':'")
	}
	v.cursor++
}

func (v *jsonValidator) parseComma() {
	if v.input[v.cursor] != ',' {
		errors.New(("invalid json: Expected ','"))
	}
	v.cursor++
}

func (v *jsonValidator) parseString() {

}

// Object Grammar
// object = begin-object [ member *( value-seprator member ) ]
//          end-object
// member = string name-seprator value

func (v *jsonValidator) parseObject() (bool, error) {
	if v.input[v.cursor] == '{' {
		v.cursor++
		initial := true
		v.skipWhiteSpace()

		for v.cursor < len(v.input) && v.input[v.cursor] != '}' {
			if !initial {
				v.skipWhiteSpace()
				v.parseComma()
				v.skipWhiteSpace()
			}

			v.parseString()
			v.skipWhiteSpace()
			v.parseColon()
			v.skipWhiteSpace()
			v.parseValue()
			v.skipWhiteSpace()
			initial = false
		}

		if v.cursor == len(v.input) {
			return false, errors.New("closing bracket for the array is missing")
		}

		v.cursor++

		return true, nil
	}

	return nil, nil
}

func (v *jsonValidator) parseNumber() {

}

//array = begin-array [ value *( value-seprator value ) ] end-array

func (v *jsonValidator) parseArray() ([]int, error) {
	if v.input[v.cursor] == '[' {
		v.cursor++
		initial := true
		v.skipWhiteSpace()
		result := []int{}

		for v.cursor < len(v.input) && v.input[v.cursor] != ']' {
			if !initial {
				v.parseComma()
				v.skipWhiteSpace()
			}

			value, err := v.parseValue()
			v.skipWhiteSpace()
			result = append(result, value)
			initial = false
		}

		if v.cursor == len(v.input) {
			return nil, errors.New("closing bracket for the array is missing")
		}

		v.cursor++
		return result, nil
	}
	return nil, nil
}

// - Following three literal names:
//   - false
//   - null
//   - true

func (v *jsonValidator) parseKeyword(name string, value nilableBool) (nilableBool, error) {
	if v.input[v.cursor:v.cursor+len(name)] == name {
		v.cursor += len(name)
		return value, nil
	}

	if name == "null" {
		return value, errors.New("invalid json: missing value")
	}
	return nilableBool{}, nil
}

// A JSON value must be:

// - Object
// - Array
// - Number
// - String
// - Following three literal names:

//   - false
//   - null
//   - true

func (v *jsonValidator) parseValue() {
	result, err := v.parseObject()

	if result == nil {
		result = v.parseArray()
	}

	if result == nil {
		result = v.parseNumber()
	}

	if result == nil {
		result = v.parseString()
	}

	if result == nil {
		result, err = v.parseKeyword("true", nilableBool{value: true})
	}

	if result == nil {
		result, err = v.parseKeyword("false", nilableBool{value: false})
	}

	if result == nil {
		result, err = v.parseKeyword("null", nilableBool{})
	}

	return result, err
}

func (v *jsonValidator) ParseJson() (bool, error) {
	v.skipWhiteSpace()

	if contains([]byte{'[', '{'}, v.input[v.cursor]) {
		result, err := v.parseValue()
		return result, err
	} else {
		return false, errors.New("a json payload should be an object or array")
	}
}

func contains(array []byte, value byte) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func NewJSONValidator(input string) *jsonValidator {
	return &jsonValidator{
		input:  input,
		cursor: 0,
	}
}
