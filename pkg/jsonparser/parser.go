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

func (v *jsonValidator) parseColon() (bool, error) {
	if v.isArrayParsed() {
		return false, nil
	}

	if v.input[v.cursor] != ':' {
		return false, errors.New("invalid json: Expected ':'")
	}
	v.cursor++
	return true, nil
}

func (v *jsonValidator) parseComma() (bool, error) {
	if v.isArrayParsed() {
		return false, nil
	}

	if v.input[v.cursor] != ',' {
		return false, errors.New(("invalid json: Expected ','"))
	}
	v.cursor++
	return true, nil
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
	if v.isArrayParsed() {
		return false, nil
	}

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
	if v.isArrayParsed() {
		return false, nil
	}

	if v.input[v.cursor] == '{' {
		v.cursor++
		initial := true
		v.skipWhiteSpace()

		for v.cursor < len(v.input) && v.input[v.cursor] != '}' {
			if !initial {
				v.skipWhiteSpace()
				if result, err := v.parseComma(); err != nil {
					return result, err
				}

				v.skipWhiteSpace()
			}

			if result, err := v.parseString(); err != nil || !result {
				if err == nil {
					return result, errors.New("object key name is not a string")
				}
				return result, err
			}

			v.skipWhiteSpace()
			if result, err := v.parseColon(); err != nil {
				return result, err
			}
			v.skipWhiteSpace()
			if result, err := v.parseValue(); err != nil || !result {
				if err == nil {
					return result, errors.New("object value is not in correct format")
				}
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

func (v *jsonValidator) isArrayParsed() bool {
	return v.cursor >= len(v.input)
}

func isValidFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)

	return err == nil
}

func (v *jsonValidator) parseNumber() (bool, error) {
	if v.isArrayParsed() {
		return false, nil
	}

	start := v.cursor

	if v.input[v.cursor] == '-' {
		v.cursor++
	}

	if v.input[v.cursor] == '0' {
		if v.cursor < len(v.input) && isDigit(v.input[v.cursor+1]) {
			return false, errors.New("number cannot have leading zeros")
		}
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
		}
		for isDigit(v.input[v.cursor]) {
			v.cursor++
		}
	}

	if v.cursor > start {
		result := isValidFloat(v.input[start:v.cursor])
		if !result {
			return false, errors.New("not a valid number")
		}

		return true, nil
	}

	return false, nil
}

//array = begin-array [ value *( value-seprator value ) ] end-array

func (v *jsonValidator) parseArray() (bool, error) {
	if v.isArrayParsed() {
		return false, nil
	}

	if v.input[v.cursor] == '[' {
		v.cursor++
		initial := true
		v.skipWhiteSpace()

		for v.cursor < len(v.input) && v.input[v.cursor] != ']' {
			if !initial {
				if result, err := v.parseComma(); err != nil {
					return result, err
				}
				v.skipWhiteSpace()
			}

			if result, err := v.parseValue(); err != nil || !result {
				if err == nil {
					return result, errors.New("array value is not correct")
				}
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
	if v.isArrayParsed() {
		return false, nil
	}

	if v.cursor+len(name) < len(v.input) && v.input[v.cursor:v.cursor+len(name)] == name {
		v.cursor += len(name)
		return true, nil
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

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseArray()
	}

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseNumber()
	}

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseString()
	}

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseKeyword("true")
	}

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseKeyword("false")
	}

	if err != nil {
		return false, err
	}

	if !result {
		result, err = v.parseKeyword("null")
	}

	return result, err
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

func (v *jsonValidator) ParseJson() (bool, error) {
	v.skipWhiteSpace()

	if contains([]byte{'[', '{'}, v.input[v.cursor]) {
		result, err := v.parseValue()

		if err != nil {
			return false, err
		}

		v.skipWhiteSpace()
		if v.cursor < len(v.input) {
			return false, errors.New("invalid json: extra character after closing bracket")
		}
		return result, err
	} else {
		return false, errors.New("a json payload should be an object or array")
	}
}
