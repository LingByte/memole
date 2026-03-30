package lexer

import (
	"fmt"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `= == ! !=`

	tests := []struct {
		expectedType TokenType
		expectedLit  string
	}{
		{Assign, "="},
		{EQ, "=="},
		{Bang, "!"},
		{NotEQ, "!="},
		{EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLit {
			t.Fatalf("test[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLit, tok.Literal)
		}

		fmt.Printf("Token [%d]: %v '%s'\n", i, tok.Type, tok.Literal)
	}
}
