package jsonparser

import (
	"errors"
	"strconv"
	"unicode"
)

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

/*
string = quotation-mark *char quotation-mark

char = unescaped /

	escape (
		%x22 /           ; " quotation mark U+0022
		%x5C /           ; \ reverse solidus U+005C
		%x2F /           ; / solidus U+002F
		%x62 /           ; b backspace U+0008
		%x66 /           ; f form feed U+000C
		%x6E /           ; n line feed U+000A
		%x72 /           ; r carriage return U+000D
		%x74 /           ; t tab U+0009
		%x75 /           ; uXXXX U+XXXX
	)

escape = %x5C            ; \

quotation-mark = %x22    ; "

unescaped = %x20-21 / %x23-5B / %x5D-10FFFF
*/

func (v *jsonValidator) parseString() (bool, error) {
	if v.input[v.cursor] == '"' {
		v.cursor++
		v.skipWhiteSpace()
		for v.input[v.cursor] != '"' {
			if v.input[v.cursor] == '\\' {
				char := v.input[v.cursor+1]
				if contains([]byte{'"', '\\', '/', 'b', 'f', 'n', 'r', 't'}, char) {
					v.cursor++
				} else if char == 'u' {
					if isHexadecimalDigit(rune(v.input[v.cursor+2])) && isHexadecimalDigit(rune(v.input[v.cursor+3])) && isHexadecimalDigit(rune(v.input[v.cursor+4])) && isHexadecimalDigit(rune(v.input[v.cursor+5])) {
						v.cursor += 5
					}
				} else {
					return false, errors.New("invalid json: Illegal backslash escape sequence")
				}
			} else {
				if v.input[v.cursor] == '\t' || v.input[v.cursor] == '\n' {
					return false, errors.New("invalid json: Illegal character tab character or new line character")
				}
			}
			v.cursor++
		}
		v.cursor++
		return true, nil
	}
	return false, nil
}

func isHexadecimalDigit(char rune) bool {
	return isDigit(byte(char)) || ('a' <= char && char <= 'f') || ('A' <= char && char <= 'F')
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

			result, err := v.parseString()
			if err != nil {
				return result, err
			}

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

	return false, nil
}

/*
	number = [ minus ] int [ frac ] [ exp ]
	decimal-point = %x2E ; .
	digit1-9 = %x31-39   ; 1-9
	e = %x65 / %x45      ; e E
	exp = e [ minus / plus ] 1*DIGIT
	frac = decimal-point 1*DIGIT
	int = zero / (digit1-9 *DIGIT)
	minus = %x2D         ; -
	plus = %x2B          ; +
	zero = %x30          ; 0
*/

func isDigit(char byte) bool {
	return unicode.IsDigit(rune(char))
}

func (v *jsonValidator) parseNumber() (bool, error) {
	start := v.cursor

	if v.input[v.cursor] == '-' {
		v.cursor++
	}

	if v.input[v.cursor] == '0' {
		v.cursor++
	}

	if isDigit(v.input[v.cursor]) {
		v.cursor++
		for isDigit(v.input[v.cursor]) {
			v.cursor++
		}
	}

	if v.input[v.cursor] == '.' {
		v.cursor++
		for isDigit(v.input[v.cursor]) {
			v.cursor++
		}
	}

	if v.input[v.cursor] == 'e' || v.input[v.cursor] == 'E' {
		v.cursor++
		if v.input[v.cursor] == '-' || v.input[v.cursor] == '+' {
			v.cursor++
			for isDigit(v.input[v.cursor]) {
				v.cursor++
			}
		}
	}

	if v.cursor > start {
		_, err := strconv.ParseFloat(v.input[start:v.cursor], 64)

		if err != nil {
			return false, err
		}

		return true, nil
	}

	return false, nil
}

//array = begin-array [ value *( value-seprator value ) ] end-array

func (v *jsonValidator) parseArray() (bool, error) {
	if v.input[v.cursor] == '[' {
		v.cursor++
		initial := true
		v.skipWhiteSpace()

		for v.cursor < len(v.input) && v.input[v.cursor] != ']' {
			if !initial {
				v.parseComma()
				v.skipWhiteSpace()
			}

			result, err := v.parseValue()
			if err != nil {
				return result, err
			}

			v.skipWhiteSpace()
			initial = false
		}

		if v.cursor == len(v.input) {
			return false, errors.New("closing bracket for the array is missing")
		}

		v.cursor++
		return true, nil
	}
	return false, nil
}

// - Following three literal names:
//   - false
//   - null
//   - true

func (v *jsonValidator) parseKeyword(name string) (bool, error) {
	if v.input[v.cursor:v.cursor+len(name)] == name {
		v.cursor += len(name)
		return true, nil
	}

	if name == "null" {
		return false, errors.New("invalid json: missing value")
	}
	return false, nil
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

func (v *jsonValidator) parseValue() (bool, error) {
	result, err := v.parseObject()

	if !result {
		result, err = v.parseArray()
	}

	if !result {
		result, err = v.parseNumber()
	}

	if !result {
		result, err = v.parseString()
	}

	if !result {
		result, err = v.parseKeyword("true")
	}

	if !result {
		result, err = v.parseKeyword("false")
	}

	if !result {
		result, err = v.parseKeyword("null")
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
