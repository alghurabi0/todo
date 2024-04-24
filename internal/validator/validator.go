package validator

import (
	"strconv"
	"strings"
)

type Validator struct {
	Errors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}
func (v *Validator) AddError(key, message string) {
	if v.Errors == nil {
		v.Errors = make(map[string]string)
	}
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}
func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}
func (v *Validator) NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}
func (v *Validator) VerifyColType(value interface{}, colType string) bool {
	switch colType {
	case "Text":
		_, ok := value.(string)
		return ok
	case "Number":
		_, err := strconv.Atoi(value.(string))
		if err != nil {
			return false
		}
		return true
	default:
		return false
	}
}
