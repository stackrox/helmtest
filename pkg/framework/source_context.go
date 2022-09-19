package framework

type sourceLocation struct {
	Line, Column int
}

type sourceContext struct {
	Key, Value sourceLocation
}

type sourceContextInfo struct {
	Elem   sourceContext
	Fields map[string]sourceContext
}
