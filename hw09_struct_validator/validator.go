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

type ValidationError struct {
	Field string
	Err   error
}

func (v ValidationError) Error() string {
	return fmt.Sprintf("field %s %s", v.Field, v.Err.Error())
}

// ValidationErrors is a slice of validation errors.
type ValidationErrors []ValidationError

func (v ValidationErrors) Error() string {
	errs := make([]string, 0, len(v))
	for _, err := range v {
		errs = append(errs, err.Error())
	}
	return fmt.Sprintf("Errors: %s", strings.Join(errs, ", "))
}

type validRule struct {
	name  string
	value string
}

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

// Validate validates the structure fields according to validation rules.
func Validate(v interface{}) error {
	val := reflect.ValueOf(v)
	if val.Kind() != reflect.Struct {
		return ErrInputNotStructure
	}
	validErrors := ValidationErrors{}[:]
	err := validateAll(val, &validErrors)
	if err != nil {
		return err
	}
	if len(validErrors) > 0 {
		return validErrors
	}
	return nil
}

func validateAll(val reflect.Value, validErrors *ValidationErrors) error {
	for i := 0; i < val.NumField(); i++ {
		fVal := val.Field(i)
		fType := val.Type().Field(i)
		fTag := fType.Tag.Get("validate")
		if fTag == "" {
			continue
		}
		ruleStrs := strings.Split(fTag, "|")
		for _, ruleStr := range ruleStrs {
			rule, err := getValidRule(ruleStr)
			if err != nil {
				return err
			}
			err = validateOne(fVal, fType.Name, *rule, validErrors)
			if err != nil && !addValidErr(err, validErrors) {
				return err
			}
		}
	}
	return nil
}

func validateOne(val reflect.Value, name string, rule validRule, validErrors *ValidationErrors) error {
	var err error
	switch val.Kind() {
	case reflect.Int:
		err = validateInt(val, name, rule)
	case reflect.String:
		err = validateString(val, name, rule)
	case reflect.Slice:
		err = validateSlice(val, name, rule, validErrors)
	case reflect.Struct:
		if rule.name == "nested" {
			err = validateAll(val, validErrors)
		} else {
			err = nil
		}
	case reflect.Invalid, reflect.Bool, reflect.Int8, reflect.Int16, reflect.Int32,
		reflect.Int64, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32,
		reflect.Uint64, reflect.Uintptr, reflect.Float32, reflect.Float64, reflect.Complex64,
		reflect.Complex128, reflect.Array, reflect.Chan, reflect.Func, reflect.Interface,
		reflect.Map, reflect.Pointer, reflect.UnsafePointer:
		err = ErrUnsupportedType(val.Kind().String())
	default:
		err = ErrUnsupportedType(val.Kind().String())
	}
	return err
}

func validateSlice(val reflect.Value, name string, rule validRule, validErrors *ValidationErrors) error {
	for i := 0; i < val.Len(); i++ {
		err := validateOne(val.Index(i), name, rule, validErrors)
		if err != nil && !addValidErr(err, validErrors) {
			return err
		}
	}
	return nil
}

func validateInt(val reflect.Value, name string, rule validRule) error {
	switch rule.name {
	case "min":
		vmin, err := strconv.ParseInt(rule.value, 10, 64)
		if err != nil {
			return err
		}
		if val.Int() < vmin {
			return ValidationError{
				Field: name,
				Err:   MinIntValidationError(vmin),
			}
		}
	case "max":
		vmax, err := strconv.ParseInt(rule.value, 10, 64)
		if err != nil {
			return err
		}
		if val.Int() > vmax {
			return ValidationError{
				Field: name,
				Err:   MaxIntValidationError(vmax),
			}
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
		if !slices.Contains(ints, val.Int()) {
			return ValidationError{
				Field: name,
				Err:   InIntValidationError(ints),
			}
		}
	default:
		return ErrUnsupportedRule("int", rule.name)
	}
	return nil
}

func validateString(val reflect.Value, name string, rule validRule) error {
	switch rule.name {
	case "len":
		slen, err := strconv.Atoi(rule.value)
		if err != nil {
			return err
		}
		if val.Len() != slen {
			return ValidationError{
				Field: name,
				Err:   LenStrValidationError(slen),
			}
		}
	case "regexp":
		re, err := regexp.Compile(rule.value)
		if err != nil {
			return err
		}
		if !re.MatchString(val.String()) {
			return ValidationError{
				Field: name,
				Err:   RegExpStrValidationError(rule.value),
			}
		}
	case "in":
		strs := strings.Split(rule.value, ",")
		if !slices.Contains(strs, val.String()) {
			return ValidationError{
				Field: name,
				Err:   InStrValidationError(strs),
			}
		}
	default:
		return ErrUnsupportedRule("string", rule.name)
	}
	return nil
}

func getValidRule(rule string) (*validRule, error) {
	if rule == "nested" {
		return &validRule{
			name:  rule,
			value: "",
		}, nil
	}
	i := strings.Index(rule, ":")
	if i > -1 {
		return &validRule{
			name:  rule[:i],
			value: rule[i+1:],
		}, nil
	}
	return nil, ErrRuleSyntax(rule)
}

func addValidErr(err error, validErrors *ValidationErrors) bool {
	var validErr ValidationError
	if errors.As(err, &validErr) {
		*validErrors = append(*validErrors, validErr)
		return true
	}
	return false
}
