# validate-json Tool

The validate-json is a go-based tool to validate if the json provided in a file is a valid json.


## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Description](#description)
- [Contributing](#contributing)

## Installation

```bash
# Clone the repository
git clone https://github.com/maniSHarma7575/validate-json

# Change directory
cd validate-json/cmd/jsonparser

# Build
go build

# Run Tests

cd validate-json/test/jsonparser
go test

# Add as linux command

ln -s validate-json/cmd/jsonparser/main /usr/local/bin/validatejson
```

## Usage

You can use validatejson utility by running command in your terminal:

`validatejson -file [file]`

### Description

#### What is JSON?

Javascript Object Notation

- Json can support four primitive data types:
	- strings
	- numbers
	- booleans
	- null

- Two structured data types:
	- Objects
	- Array

- String: Sequence of zero or more unicode characters
- Object: unordered collection of zero or more name/value pairs
	- Name: string
	- Value: string, number, boolean, null, object, or array
- Array: unordered sequence or zero or more values.

#### JSON Grammar

**Lexer**

A JSON text is a sequence of tokens. The set to tokens includes:
- Structural characters
- Strings
- Numbers
- Three literal names

`JSON-text = ws value ws`

Six structural characters:

- begin-array = ws %x5B ws ; [ left square bracket
- begin-object = ws %x7B ws ; { left curly bracket
- end-array = ws %x5D ws ; ] right square bracket
- end-object = ws %x7D ws ; } right curly bracket
- name-separator = ws %x3A ws ; : colon
- value-seprator = ws %x2C ws ; , comma

```
ws = *(
        %x20 /            ; Space
        %x09 /            ; Horizontal tab
        %x0A /            ; Line Feed or New line
        %x0D )            ; Carriage return
      )
```

**Values**

A JSON value must be:

- Object
- Array
- Number
- String
- Following three literal names:

  - false
  - null
  - true

**Objects**

```
object = begin-object [ member *( value-seprator member ) ]
         end-object

member = string name-seprator value
```

**Array**

```
array = begin-array [ value *( value-seprator value ) ] end-array
```

**Numbers**

```
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
```

**Strings**

```
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
```

## Contributing

Thank you for your interest in contributing to our project! We welcome your suggestions, improvements, or contributions. To get started, follow these steps:

### 1. Fork the Project

Click the "Fork" button on the top-right corner of this repository to create your own copy of the project.

### 2. Create a New Branch

Once you've forked the project, it's a good practice to create a new branch for your changes. This keeps your changes isolated and makes it easier to manage multiple contributions.

```bash
git checkout -b your-new-branch
```
