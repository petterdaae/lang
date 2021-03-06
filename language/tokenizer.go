package language

import (
	"bufio"
	"bytes"
	"io"
	"unicode"
)

const (
	Number Token = iota
	Character
	Plus
	Minus
	Modulo
	Divide
	Multiply
	Assign
	Question
	Hash
	Dot
	Loop
	Colon
	Pipe
	NotEqual
	Length
	At
	Struct
	Not
	String
	RoundBracketStart
	RoundBracketEnd
	BoxBracketStart
	BoxBracketEnd
	CurlyBracketStart
	CurlyBracketEnd
	AngleBracketStart
	AngleBracketEnd
	Comma
	TypeInt
	TypeBool
	TypeChar
	TypeString
	TypeList
	TypeFunc
	Whitespace
	PrintLn
	Return
	Identifier
	Placeholder
	True
	False
	If
	Equals
	EOF
	Error
)

const eof = rune(0)

type Token int

type Tokenizer struct {
	reader *bufio.Reader
}

func NewTokenizer(reader io.Reader) *Tokenizer {
	return &Tokenizer{
		reader: bufio.NewReader(reader),
	}
}

func (tokenizer *Tokenizer) read() rune {
	character, _, err := tokenizer.reader.ReadRune()
	if err != nil {
		return eof
	}
	return character
}

func (tokenizer *Tokenizer) unread() {
	tokenizer.reader.UnreadRune()
}

func validIdentifierChar(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || r == '_'
}

func (tokenizer *Tokenizer) NextToken() (Token, string) {
	character := tokenizer.read()

	if unicode.IsSpace(character) {
		tokenizer.unread()
		return tokenizer.whitespace()
	}

	if unicode.IsDigit(character) {
		tokenizer.unread()
		return tokenizer.number()
	}

	if validIdentifierChar(character) {
		tokenizer.unread()
		return tokenizer.identifier()
	}

	if character == '\'' {
		tokenizer.unread()
		return tokenizer.character()
	}

	if character == '"' {
		tokenizer.unread()
		return tokenizer.string()
	}

	switch character {
	case eof:
		return EOF, ""
	case '*':
		return Multiply, string(character)
	case '+':
		return Plus, string(character)
	case '(':
		return RoundBracketStart, string(character)
	case ')':
		return RoundBracketEnd, string(character)
	case '=':
		next := tokenizer.read()
		if next == '=' {
			return Equals, "=="
		}
		tokenizer.unread()
		return Assign, string(character)
	case '{':
		return CurlyBracketStart, string(character)
	case '}':
		return CurlyBracketEnd, string(character)
	case '[':
		return BoxBracketStart, string(character)
	case ']':
		return BoxBracketEnd, string(character)
	case '<':
		return AngleBracketStart, string(character)
	case '>':
		return AngleBracketEnd, string(character)
	case ',':
		return Comma, string(character)
	case '_':
		return Placeholder, string(character)
	case '-':
		return Minus, string(character)
	case '/':
		return Divide, string(character)
	case '%':
		return Modulo, string(character)
	case '?':
		return Question, string(character)
	case ':':
		return Colon, string(character)
	case '|':
		return Pipe, string(character)
	case '!':
		next := tokenizer.read()
		if next == '=' {
			return NotEqual, "!="
		}
		return Not, "!"
	case '@':
		return At, string(character)
	case '.':
		return Dot, string(character)
	case '#':
		return Hash, string(character)
	}
	
	return Error, ""
}

func (tokenizer *Tokenizer) string() (Token, string) {
	result := ""

	character := tokenizer.read()
	if character != '"' {
		return Error, ""
	}

	for {
		character = tokenizer.read()
		if character == '"' {
			break
		}

		if character == '\\' {
			character = tokenizer.read()
			if character != '"' && character != '\\' {
				return Error, ""
			}
		}

		result = result + string(character)
	}

	return String, result
}

func (tokenizer *Tokenizer) character() (Token, string) {
	character := tokenizer.read()
	if character != '\'' {
		return Error, ""
	}

	character = tokenizer.read()

	if character == '\'' {
		return Error, ""
	}

	if character == '\\' {
		character = tokenizer.read()
		if character != '\'' && character != '\\' {
			return Error, ""
		}
	}

	end := tokenizer.read()
	if end != '\'' {
		return Error, ""
	}

	return Character, string(character)
}

func (tokenizer *Tokenizer) whitespace() (Token, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(tokenizer.read())

	for {
		character := tokenizer.read()
		if character == eof || !unicode.IsSpace(character) {
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}

	return Whitespace, buffer.String()
}

func (tokenizer *Tokenizer) number() (Token, string) {
	var buffer bytes.Buffer
	first := tokenizer.read()
	buffer.WriteRune(first)

	if first == '0' {
		return Number, "0"
	}

	for {
		character := tokenizer.read()
		if !unicode.IsDigit(character) {
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}
	return Number, buffer.String()
}

func (tokenizer *Tokenizer) identifier() (Token, string) {
	var buffer bytes.Buffer
	buffer.WriteRune(tokenizer.read())

	for {
		character := tokenizer.read()
		if !validIdentifierChar(character) && !unicode.IsDigit(character){
			tokenizer.unread()
			break
		}
		buffer.WriteRune(character)
	}

	word := buffer.String()

	switch word {
	case "println":
		return PrintLn, word
	case "int":
		return TypeInt, word
	case "char":
		return TypeChar, word
	case "list":
		return TypeList, word
	case "return":
		return Return, word
	case "if":
		return If, word
	case "true":
		return True, word
	case "false":
		return False, word
	case "bool":
		return TypeBool, word
	case "func":
		return TypeFunc, word
	case "loop":
		return Loop, word
	case "struct":
		return Struct, word
	case "string":
		return TypeString, word
	case "len":
		return Length, word
	}

	return Identifier, word
}
