package language

import (
	"callmemaybe/language/common"
	"fmt"
	"io"
	"strconv"
)

type Parser struct {
	tokenizer *Tokenizer
	buffer    struct {
		kind  Token
		token string
		full  bool
	}
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		tokenizer: NewTokenizer(reader),
	}
}

func (parser *Parser) Parse() (Stmt, error) {
	stmt, err := parser.parseSeq()
	nextKind, nextStr := parser.readIgnoreWhiteSpace()
	if err == nil && nextKind != EOF {
		return nil, fmt.Errorf("failed to parse the entire program: %s", nextStr)
	}
	return stmt, err
}

func (parser *Parser) read() (Token, string) {
	if parser.buffer.full {
		parser.buffer.full = false
		return parser.buffer.kind, parser.buffer.token
	}
	kind, token := parser.tokenizer.NextToken()
	parser.buffer.kind = kind
	parser.buffer.token = token
	return kind, token
}

func (parser *Parser) unread() {
	parser.buffer.full = true
}

func (parser *Parser) readIgnoreWhiteSpace() (Token, string) {
	kind, token := parser.read()
	if kind == Whitespace {
		kind, token = parser.read()
	}
	return kind, token
}

func (parser *Parser) ParseExp() (Exp, error) {
	nextKind, _ := parser.readIgnoreWhiteSpace()
	parser.unread()

	if nextKind == Call {
		return parser.parseCall()
	}

	if nextKind == AngleBracketStart {
		return parser.parseFunction()
	}

	if nextKind == BoxBracketStart {
		return parser.parseList()
	}

	return parser.ParseCalculation()
}

func (parser *Parser) ParseCalculation() (Exp, error) {
	left, err := parser.parseVal()
	if err != nil {
		return nil, fmt.Errorf("failed to parse first val in exp: %w", err)
	}
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == Plus {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of plus exp: %w", err)
			}
			left = ExpPlus{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Multiply {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of multiply exp: %w", err)
			}
			left = ExpMultiply{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Divide {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of divide exp: %w", err)
			}
			left = ExpDivide{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Modulo {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of modulo exp: %w", err)
			}
			left = ExpModulo{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == Minus {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of minus exp: %w", err)
			}
			left = ExpMinus{
				Left:  left,
				Right: right,
			}
			continue
		}
		if nextKind == AngleBracketStart {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of less expression")
			}
			left = ExpLess{
				Left: left,
				Right: right,
			}
			continue
		}
		if nextKind == AngleBracketEnd {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of greater expression")
			}
			left = ExpGreater{
				Left: left,
				Right: right,
			}
			continue
		}
		if nextKind == Equals {
			right, err := parser.parseVal()
			if err != nil {
				return nil, fmt.Errorf("failed to parse right side of equals expression")
			}
			left = ExpEquals{
				Left: left,
				Right: right,
			}
			continue
	}
		parser.unread()
		break
	}
	return left, nil
}

func (parser *Parser) parseVal() (Exp, error) {
	nextKind, nextToken := parser.readIgnoreWhiteSpace()
	if nextKind == Number {
		value, _ := strconv.Atoi(nextToken)
		return ExpNum{
			Value: value,
		}, nil
	}
	if nextKind == True {
		return ExpBool{
			Value: true,
		}, nil
	}
	if nextKind == False {
		return ExpBool{
			Value: false,
		}, nil
	}
	if nextKind == RoundBracketStart {
		inside, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse exp in parentheses: %w", err)
		}
		nextKind, _ = parser.readIgnoreWhiteSpace()
		if nextKind != RoundBracketEnd {
			return nil, fmt.Errorf("missing closing parentheses")
		}
		return ExpParentheses{
			Inside: inside,
		}, nil
	}
	if nextKind == Identifier {
		return ExpIdentifier{
			Name: nextToken,
		}, nil
	}
	if nextKind == Minus {
		inside, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse exp in negative expression: %w", err)
		}
		return ExpNegative{
			Inside: inside,
		}, nil
	}
	if nextKind == Character {
		return ExpChar{
			Value: nextToken,
		}, nil
	}
	if nextKind == Get {
		parser.unread()
		return parser.parseGetFromList()
	}
	return nil, fmt.Errorf("unexpected token while parsing val")
}

func (parser *Parser) parseAssign() (Stmt, error) {
	kind, identifier := parser.readIgnoreWhiteSpace()
	if kind != Identifier && kind != Placeholder {
		return nil, fmt.Errorf("failed to parse identifier at start of assign statement")
	}
	kind, token := parser.readIgnoreWhiteSpace()
	if kind != Assign {
		return nil, fmt.Errorf("expected assign operator in assign stmt but got: %s", token)
	}
	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in assign stmt: %w", err)
	}
	return StmtAssign{Identifier: identifier, Expression: expr}, nil
}

func (parser *Parser) parsePrintln() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != PrintLn {
		return nil, fmt.Errorf("expected println keyword at start of println stmt")
	}
	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in print stmt: %w", err)
	}
	return StmtPrintln{Expression: expr}, nil
}

func (parser *Parser) parseSeq() (Stmt, error) {
	var statements []Stmt
	for {
		nextKind, _ := parser.readIgnoreWhiteSpace()
		if nextKind == Identifier || nextKind == Placeholder {
			parser.unread()
			statement, err := parser.parseAssign()
			if err != nil {
				return nil, fmt.Errorf("failed to parse assign expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == PrintLn {
			parser.unread()
			statement, err := parser.parsePrintln()
			if err != nil {
				return nil, fmt.Errorf("failed to parse println expression: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Return {
			expr, err := parser.ParseExp()
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression after return: %w", err)
			}
			statement := StmtReturn{Expression: expr}
			statements = append(statements, statement)
			continue
		}
		if nextKind == If {
			parser.unread()
			statement, err := parser.parseIf()
			if err != nil {
				return nil, fmt.Errorf("failed to parse if statement: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		if nextKind == Append {
			parser.unread()
			statement, err := parser.parseAppendToList()
			if err != nil {
				return nil, fmt.Errorf("failed to parse append statement: %w", err)
			}
			statements = append(statements, statement)
			continue
		}
		parser.unread()
		break
	}
	return StmtSeq{Statements: statements}, nil
}

func (parser *Parser) parseIf() (Stmt, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != If {
		return nil, fmt.Errorf("expected if keyword at start of if statement")
	}

	expr, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse condition expression in if statement: %w", err)
	}

	kind, text := parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected { when parsing if statement, but got: %s", text)
	}

	seq, err := parser.parseSeq()
	if err != nil {
		return nil, fmt.Errorf("failed to parse sequence in if statement: %w", err)
	}

	kind, text = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketEnd {
		return nil, fmt.Errorf("expected } when parsing if statement, but got: %s", text)
	}

	return StmtIf{
		Expression: expr,
		Body: seq,
	}, nil
}

func (parser *Parser) parseCall() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Call {
		return nil, fmt.Errorf("expected call keyword at start of function call")
	}

	kind, identifier := parser.readIgnoreWhiteSpace()
	if kind != Identifier {
		return nil, fmt.Errorf("expected identifier after keyword call in function call")
	}

	call := FunctionCall{}
	call.Name = identifier

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind == With {
		for {
			expr, err := parser.ParseExp()
			if err != nil {
				return nil, fmt.Errorf("failed to parse expression in function call")
			}
			call.Arguments = append(call.Arguments, expr)
			kind, _ = parser.readIgnoreWhiteSpace()
			if kind != Comma {
				parser.unread()
				break
			}
		}
	} else {
		parser.unread()
	}

	return call, nil
}

func (parser *Parser) parseFunction() (Exp, error) {
	function := ExpFunction{}
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != AngleBracketStart {
		return nil, fmt.Errorf("expected < at start of function expression")
	}

	first := true
	for {
		kind, identifier := parser.readIgnoreWhiteSpace()
		if kind == AngleBracketEnd {
			break
		}

		if kind != Identifier {
			return nil, fmt.Errorf("expected identifier when parsing argument list, but got %s", identifier)
		}

		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == TypeEmpty {
			return nil, fmt.Errorf("only non-empty types allowed in function arguments")
		}

		if kind == Comma && first {
			function.Recurse = identifier
			if identifier != "me" {
				return nil, fmt.Errorf("the first recurse argument in a function has to be named me")
			}
			continue
		}
		first = false


		contextKind := kindFromType(kind)
		if contextKind == common.ContextElementKindInvalid {
			return nil, fmt.Errorf("expected valid type when parsing argument in argument list")
		}

		kind, _ = parser.readIgnoreWhiteSpace()

		function.Args = append(function.Args, common.Arg{Identifier: identifier, Type: contextKind})

		if kind == AngleBracketEnd {
			break
		}

		if kind == Comma {
			continue
		}

		return nil, fmt.Errorf("expected comma or end of arguemnt list")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != Arrow {
		return nil, fmt.Errorf("expected arrow after argument list when parsing function")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind == CurlyBracketStart {
		function.ReturnType = common.ContextElementKindEmpty
	} else {
		function.ReturnType = kindFromType(kind)
		if function.ReturnType == common.ContextElementKindInvalid {
			return nil, fmt.Errorf("invalid return type while parsing function definition")
		}
		kind, _ = parser.readIgnoreWhiteSpace()
	}


	if kind != CurlyBracketStart {
		return nil, fmt.Errorf("expected opening curly bracket when parsing function")
	}

	seq, err := parser.parseSeq()
	if err != nil {
		return nil, fmt.Errorf("failed to parse statements in function: %w", err)
	}

	function.Body = seq

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != CurlyBracketEnd {
		return nil, fmt.Errorf("expected closing curly bracker when parsing function")
	}

	return function, nil
}

func (parser *Parser) parseList() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != BoxBracketStart {
		return nil, fmt.Errorf("expected box bracket at start of list declaration")
	}

	list := ExpList{}
	first := true

	for {
		kind, _ = parser.readIgnoreWhiteSpace()
		if first && kind == BoxBracketEnd {
			break
		}
		first = false
		parser.unread()
		expr, err := parser.ParseExp()
		if err != nil {
			return nil, fmt.Errorf("failed to parse expression in list declaration: %w", err)
		}
		list.Elements = append(list.Elements, expr)
		kind, _ = parser.readIgnoreWhiteSpace()
		if kind == Comma {
			continue
		}
		if kind == BoxBracketEnd {
			break
		}
		return nil, fmt.Errorf("unexpected token when parsing list declareation")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != Colon {
		return nil, fmt.Errorf("expected comma when parsing list declaration")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	contextKind := kindFromType(kind)
	if contextKind == common.ContextElementKindInvalid || contextKind == common.ContextElementKindEmpty {
		return nil, fmt.Errorf("invalid list type")
	}

	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != Colon {
		return nil, fmt.Errorf("expected comma when parsing list declaration")
	}

	kind, number := parser.readIgnoreWhiteSpace()
	if kind != Number {
		return nil, fmt.Errorf("expected number as list size")
	}

	parsed, _ := strconv.Atoi(number)

	list.Type = contextKind
	list.Size = parsed

	return list, nil
}

func (parser *Parser) parseGetFromList() (Exp, error) {
	kind, _ := parser.readIgnoreWhiteSpace()
	if kind != Get {
		return nil, fmt.Errorf("expected get when parsing get from list")
	}
	numExp, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse index expression in get from list")
	}
	kind, _ = parser.readIgnoreWhiteSpace()
	if kind != From {
		return nil, fmt.Errorf("expected from keyword when parsing get from list")
	}
	exp, err := parser.ParseExp()
	if err != nil {
		return nil, fmt.Errorf("failed to parse expression in get from list: %w", err)
	}
	return ExpGetFromList{
		List: exp,
		Index: numExp,
	}, nil
}

func (parser *Parser) parseAppendToList() (Stmt, error) {
	// TODO : implement
	return nil, fmt.Errorf("not implemented")
}

func kindFromType(token Token) common.ContextElementKind {
	switch token {
	case TypeInt:
		return common.ContextElementKindNumber
	case TypeChar:
		return common.ContextElementKindChar
	case TypeBool:
		return common.ContextElementKindBoolean
	case TypeEmpty:
		return common.ContextElementKindEmpty
	default:
		return common.ContextElementKindInvalid
	}
}
