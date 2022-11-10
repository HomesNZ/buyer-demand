package auth

type contextKey int

const (
	contextKeyRules contextKey = iota
	contextKeyToken
)
