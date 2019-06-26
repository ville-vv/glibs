// Package validator
package valid

// 校验器
type Validator interface {
	Generate() string
	Verify(string) bool
}
