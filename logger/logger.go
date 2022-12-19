package logger

import "fmt"

func PrettyPrintErr(fileName string, lineNum int, linePos int, err error) {
	fmt.Printf("file: %s line: %d column: %d, error: %s", fileName, lineNum, linePos, err.Error())
}
