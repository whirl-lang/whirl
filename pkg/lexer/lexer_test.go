package lexer

import (
	"testing"
)

func TestLexerEmptyMain(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLexerVariables(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { let a: int = 5; escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLexerArray(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { let a: int[] = [1, 2, 3, 4, 5]; escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestLexerIf(t *testing.T) {
	err := CheckForErrorsInIterator([]byte("proc main() :: int { if 5 == 5 { escape 0; } escape 0; }"))

	if err != nil {
		t.Fatalf(err.Error())
	}
}

func CheckForErrorsInIterator(input []byte) error {
	i := Iterator([]byte(input))

	for {
		token, err := i.Next()

		if err != nil {
			return err
		}

		if token.Kind == EOF {
			return nil
		}
	}

}
