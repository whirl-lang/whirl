package parser

import (
	"testing"

	"github.com/whirl-lang/whirl/pkg/lexer"
)

func TestParserEmptyMain(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestParserVariables(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { let a: int = 5; escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestParserArray(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { let a: int[] = [1, 2, 3, 4, 5]; escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestParserIf(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { if 5 == 5 { escape 0; } escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func CheckForErrorsInIterator(input []byte) error {
	lexerIterator := lexer.Iterator([]byte(input))
	instructionIterator := Iterator(lexerIterator)

	for {
		token, err := instructionIterator.Next()

		if err != nil {
			return err
		}

		if token == nil {
			return nil
		}
	}

}
