package depthparser

import (
	"fmt"
	"strconv"
)

func Parse(input string) (string, error) {
	lex := newLexer(input)

	token := lex.lex()
	err := checkErr(token)
	if err != nil {
		return "", err
	}
	if token.typ != tokenNumber {
		return "", parserErr(incorrectType(tokenNumber, token))
	}
	numSpaces, err := strconv.Atoi(token.value)
	if err != nil {
		return "", err
	}

	token = lex.lex()
	err = checkErr(token)
	if err != nil {
		return "", err
	}
	if token.typ != tokenSpaceType {
		return "", parserErr(incorrectType(tokenSpaceType, token))
	}

	var char rune
	switch token.value {
	case "s":
		char = ' '
	case "t":
		char = '	'
	default:
		panic(
			fmt.Sprintf(
				"panic: somehow things lexed correctly but you have submitted a letter that is not covered by the switch case: %s",
				token.value,
			),
		)
	}

	spaceStop := make([]byte, numSpaces)
	for i := 0; i < numSpaces; i++ {
		spaceStop[i] = byte(char)
	}

	return string(spaceStop), nil
}

func checkErr(t token) error {
	if t.typ == tokenError {
		return fmt.Errorf("lexer error: %s", t.value)
	}
	return nil
}

func parserErr(msg string) error {
	return fmt.Errorf("parse error: %s", msg)
}

func incorrectType(tt tokenType, actualToken token) string {
	return fmt.Sprintf("expected a %s but got %s with value %q at position %d", tt, actualToken.typ, actualToken.value, actualToken.pos)
}
