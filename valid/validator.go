// Package validator
package valid

// 校验器
type Validator interface {
	Generate(map[string]interface{}) string
	Verify(string) (map[string]interface{}, bool)
}
