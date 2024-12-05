package hw09structvalidator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

var (
	ErrInputNotStructure = errors.New("input param is not a structure")
	ErrUnsupportedRule   = func(dataType, rule string) error {
		return fmt.Errorf("unsupported %s validation rule: %s", dataType, rule)
	}
	ErrRuleSyntax      = func(rule string) error { return fmt.Errorf("unsupported validation rule syntax: %s", rule) }
	ErrUnsupportedType = func(dataType string) error { return fmt.Errorf("unsupported type: %s", dataType) }
)

var (
	// MinIntValidationError is a validation error for int64 min value rule.
	MinIntValidationError = func(vmin int64) error { return fmt.Errorf("must be greater than or equal to %d", vmin) }
	// MaxIntValidationError is a validation error for int64 max value rule.
	MaxIntValidationError = func(vmax int64) error { return fmt.Errorf("must be less than or equal to %d", vmax) }
	// InIntValidationError is a validation error for int64 array containing value rule.
	InIntValidationError = func(in []int64) error { return fmt.Errorf("must be in %v", in) }
	// InStrValidationError is a validation error for string array containing value rule.
	InStrValidationError = func(in []string) error { return fmt.Errorf("must be in %v", in) }
	// LenStrValidationError is a validation error for string len rule.
	LenStrValidationError = func(slen int) error { return fmt.Errorf("must be of len %d", slen) }
	// RegExpStrValidationError is a validation error for string regexp rule.
	RegExpStrValidationError = func(regExp string) error { return fmt.Errorf("must match regexp %s", regExp) }
)

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %s %s", v.Field, v.Err.Error())
}

type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("Errors: %s", strings.Join(errs, ", "))
}

type fieldValue struct {
	name  string
	value reflect.Value
}

type ruleValidator interface {
	Validate(field fieldValue, vErrs *ValidationErrors) error
}

type rule struct {
	name  string
	value string
}

type intRuleValidator struct {
	rule
}

func (rule intRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	switch rule.name {
	case "min":
		vMin, err := strconv.ParseInt(rule.value, 10, 64)
		if err != nil {
			return err
		}
		if field.value.Int() < vMin {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   MinIntValidationError(vMin),
			})
			return nil
		}
	case "max":
		vMax, err := strconv.ParseInt(rule.value, 10, 64)
		if err != nil {
			return err
		}
		if field.value.Int() > vMax {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   MaxIntValidationError(vMax),
			})
			return nil
		}
	case "in":
		nums := strings.Split(rule.value, ",")
		var ints []int64
		for _, num := range nums {
			i, err := strconv.ParseInt(num, 10, 64)
			if err != nil {
				return err
			}
			ints = append(ints, i)
		}
		if !slices.Contains(ints, field.value.Int()) {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   InIntValidationError(ints),
			})
			return nil
		}
	default:
		return ErrUnsupportedRule("int", rule.name)
	}
	return nil
}

type strRuleValidator struct {
	rule
	regExp *regexp.Regexp
}

func (rule strRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	switch rule.name {
	case "len":
		slen, err := strconv.Atoi(rule.value)
		if err != nil {
			return err
		}
		if field.value.Len() != slen {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   LenStrValidationError(slen),
			})
			return nil
		}
	case "regexp":
		if !rule.regExp.MatchString(field.value.String()) {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   RegExpStrValidationError(rule.value),
			})
			return nil
		}
	case "in":
		strs := strings.Split(rule.value, ",")
		if !slices.Contains(strs, field.value.String()) {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   InStrValidationError(strs),
			})
			return nil
		}
	default:
		return ErrUnsupportedRule("string", rule.name)
	}
	return nil
}

type sliceRuleValidator struct {
	rule
}

func (rule sliceRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	if field.value.Len() == 0 {
		return nil
	}
	vRule, err := getRuleValidator(field.value.Index(0), rule.name, rule.value)
	if err != nil {
		return err
	}
	if vRule != nil {
		for i := 0; i < field.value.Len(); i++ {
			err = vRule.Validate(fieldValue{field.name, field.value.Index(i)}, vErrs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type structRuleValidator struct {
	rule
}

func (rule structRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	for i := 0; i < field.value.NumField(); i++ {
		fVal := field.value.Field(i)
		fType := field.value.Type().Field(i)
		fTag := fType.Tag.Get("validate")
		if fTag == "" {
			continue
		}
		ruleStrs := strings.Split(fTag, "|")
		for _, ruleStr := range ruleStrs {
			rName := ruleStr
			rVal := ""
			i := strings.Index(ruleStr, ":")
			if i > -1 {
				rName = ruleStr[:i]
				rVal = ruleStr[i+1:]
			}
			vRule, err := getRuleValidator(fVal, rName, rVal)
			if err != nil {
				return err
			}
			if vRule != nil {
				err = vRule.Validate(fieldValue{fType.Name, fVal}, vErrs)
				if err != nil {
					return err
				}
			}
		}
	}
	return nil
}

// Validate validates the structure fields according to validation rules.
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrInputNotStructure
	}
	validErrors := ValidationErrors{}[:]
	rule, err := getRuleValidator(val, "nested", "")
	if err != nil {
		return err
	}
	err = rule.Validate(fieldValue{val.Type().Name(), val}, &validErrors)
	if err != nil {
		return err
	}
	if len(validErrors) > 0 {
		return validErrors
	}
	return nil
}

func getRuleValidator(fieldVal reflect.Value, ruleName string, ruleVal string) (ruleValidator, error) {
	switch fieldVal.Kind() {
	case reflect.Int:
		return intRuleValidator{rule{ruleName, ruleVal}}, nil
	case reflect.String:
		if ruleName == "regexp" {
			re, err := regexp.Compile(ruleVal)
			if err != nil {
				return nil, err
			}
			return strRuleValidator{rule{ruleName, ruleVal}, re}, nil
		}
		return strRuleValidator{rule{ruleName, ruleVal}, nil}, nil
	case reflect.Slice:
		if fieldVal.Len() > 0 {
			return sliceRuleValidator{rule{ruleName, ruleVal}}, nil
		} else {
			return nil, nil
		}
	case reflect.Struct:
		if ruleName == "nested" {
			return structRuleValidator{rule{ruleName, ruleVal}}, nil
		} else {
			return nil, nil
		}
	case reflect.Invalid, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.UnsafePointer:
		return nil, ErrUnsupportedType(fieldVal.Kind().String())
	default:
		return nil, ErrUnsupportedType(fieldVal.Kind().String())
	}
}
