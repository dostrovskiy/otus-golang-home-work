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
	min int64
	max int64
	in  []int64
}

func (rule intRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	switch rule.name {
	case "min":
		if field.value.Int() < rule.min {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   MinIntValidationError(rule.min),
			})
			return nil
		}
	case "max":
		if field.value.Int() > rule.max {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   MaxIntValidationError(rule.max),
			})
			return nil
		}
	case "in":
		if !slices.Contains(rule.in, field.value.Int()) {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   InIntValidationError(rule.in),
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
	len    int
	in     []string
}

func (rule strRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	switch rule.name {
	case "len":
		if field.value.Len() != rule.len {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   LenStrValidationError(rule.len),
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
		if !slices.Contains(rule.in, field.value.String()) {
			*vErrs = append(*vErrs, ValidationError{
				Field: field.name,
				Err:   InStrValidationError(rule.in),
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
	inner ruleValidator
}

func (rule sliceRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	if field.value.Len() > 0 && rule.inner != nil {
		for i := 0; i < field.value.Len(); i++ {
			err := rule.inner.Validate(fieldValue{field.name, field.value.Index(i)}, vErrs)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

type structRuleValidator struct {
	rule
	nested map[string][]ruleValidator
}

func (rule structRuleValidator) Validate(field fieldValue, vErrs *ValidationErrors) error {
	for i := 0; i < field.value.NumField(); i++ {
		fType := field.value.Type().Field(i)
		rules, ok := rule.nested[fType.Name]
		if !ok {
			continue
		}
		for _, rule := range rules {
			err := rule.Validate(fieldValue{fType.Name, field.value.Field(i)}, vErrs)
			if err != nil {
				return err
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
		return getIntRuleValidator(ruleName, ruleVal)
	case reflect.String:
		return getStrRuleValidator(ruleName, ruleVal)
	case reflect.Slice:
		return getSliceRuleValidator(fieldVal, ruleName, ruleVal)
	case reflect.Struct:
		return getStructRuleValidator(fieldVal, ruleName, ruleVal)
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

func getIntRuleValidator(ruleName, ruleVal string) (ruleValidator, error) {
	switch ruleName {
	case "min":
		vMin, err := strconv.ParseInt(ruleVal, 10, 64)
		if err != nil {
			return nil, err
		}
		return intRuleValidator{rule{ruleName, ruleVal}, vMin, 0, nil}, nil
	case "max":
		vMax, err := strconv.ParseInt(ruleVal, 10, 64)
		if err != nil {
			return nil, err
		}
		return intRuleValidator{rule{ruleName, ruleVal}, 0, vMax, nil}, nil
	case "in":
		nums := strings.Split(ruleVal, ",")
		ints := []int64{}
		for _, num := range nums {
			i, err := strconv.ParseInt(num, 10, 64)
			if err != nil {
				return nil, err
			}
			ints = append(ints, i)
		}
		return intRuleValidator{rule{ruleName, ruleVal}, 0, 0, ints}, nil
	default:
		return nil, ErrUnsupportedRule("int", ruleName)
	}
}

func getStrRuleValidator(ruleName, ruleVal string) (ruleValidator, error) {
	switch ruleName {
	case "len":
		sLen, err := strconv.Atoi(ruleVal)
		if err != nil {
			return nil, err
		}
		return strRuleValidator{rule{ruleName, ruleVal}, nil, sLen, nil}, nil
	case "regexp":
		re, err := regexp.Compile(ruleVal)
		if err != nil {
			return nil, err
		}
		return strRuleValidator{rule{ruleName, ruleVal}, re, 0, nil}, nil
	case "in":
		return strRuleValidator{rule{ruleName, ruleVal}, nil, 0, strings.Split(ruleVal, ",")}, nil
	default:
		return nil, ErrUnsupportedRule("str", ruleName)
	}
}

func getSliceRuleValidator(fieldVal reflect.Value, ruleName, ruleVal string) (ruleValidator, error) {
	if fieldVal.Len() > 0 {
		inner, err := getRuleValidator(fieldVal.Index(0), ruleName, ruleVal)
		if err != nil {
			return nil, err
		}
		return sliceRuleValidator{rule{ruleName, ruleVal}, inner}, nil
	}
	return nil, nil
}

func getStructRuleValidator(fieldVal reflect.Value, ruleName, ruleVal string) (ruleValidator, error) {
	if ruleName == "nested" {
		nested := make(map[string][]ruleValidator)
		for i := 0; i < fieldVal.NumField(); i++ {
			fVal := fieldVal.Field(i)
			fType := fieldVal.Type().Field(i)
			fTag := fType.Tag.Get("validate")
			if fTag == "" {
				continue
			}
			nested[fType.Name] = make([]ruleValidator, 0)
			strs := strings.Split(fTag, "|")
			for _, str := range strs {
				rName, rVal := str, "" //nolint
				idx := strings.Index(rName, ":")
				if idx > -1 {
					rName, rVal = rName[:idx], rName[idx+1:]
				}
				rule, err := getRuleValidator(fVal, rName, rVal)
				if err != nil {
					return nil, err
				}
				nested[fType.Name] = append(nested[fType.Name], rule)
			}
		}
		return structRuleValidator{rule{ruleName, ruleVal}, nested}, nil
	}
	return nil, nil
}
