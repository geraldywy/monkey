package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/geraldywy/monkey/token"

	"github.com/geraldywy/monkey/lexer"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer, filename string) {
	scanner := bufio.NewScanner(in)
	for {
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line, filename)
		for tok, err := l.NextToken(); tok.Type != token.EOF; tok, err = l.NextToken() {
			if err != nil {
				fmt.Printf("error parsing token: %v\n", err)
				continue
			}
			fmt.Printf("%+v\n", tok)
		}
	}
}
