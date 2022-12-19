package utils

func IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func IsAlphaOrUnderscore(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '_'
}

func IsWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\n' || ch == '\t' || ch == '\r'
}
