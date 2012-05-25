package govalidations

import (
	"regexp"
	"strings"
)

type ValueGetter func(object interface{}) interface{}

type Validator func(object interface{}) []*Error

func CustomValidator(vg ValueGetter, vd func(value interface{}) bool, name string, message string) Validator {
	return func(object interface{}) (r []*Error) {
		val := vg(object)
		if vd(val) {
			return
		}

		r = append(r, &Error{
			Name:    name,
			Message: message,
		})
		return
	}
}

func FormatValidator(vg ValueGetter, matcher *regexp.Regexp, name string, message string) Validator {
	return CustomValidator(vg, func(value interface{}) bool {
		return matcher.MatchString(value.(string))
	}, name, message)
}

func PresenceValidator(vg ValueGetter, name string, message string) Validator {
	return CustomValidator(vg, func(value interface{}) bool {
		return strings.Trim(value.(string), " 　") != ""
	}, name, message)
}
