package govalidations

import (
	"regexp"
	"strings"
)

type ValueGetter func(object interface{}) interface{}

type Validator func(object interface{}) []*Error

func Custom(vd func(object interface{}) bool, name string, message string) Validator {
	return func(object interface{}) (r []*Error) {
		if vd(object) {
			return
		}

		r = append(r, &Error{
			Name:    name,
			Message: message,
		})
		return
	}
}

func Regexp(vg ValueGetter, matcher *regexp.Regexp, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return matcher.MatchString(value)
	}, name, message)
}

func Presence(vg ValueGetter, name string, message string) Validator {
	return Custom(func(object interface{}) bool {
		value := vg(object).(string)
		return strings.Trim(value, " 　") != ""
	}, name, message)
}

func Limitation(vg ValueGetter, min int, max int, name string, message string) (Validator){
        return Custom(func(object interface{}) bool {
                value := vg(object).(string)
                return len(value) >= min && len(value) <= max
        }, name, message)
}

func Equation(vg ValueGetter, text string, name string, message string) (Validator){
        return Custom(func(object interface{}) bool {
                value := vg(object).(string)
                return value == text
        }, name, message)
}
