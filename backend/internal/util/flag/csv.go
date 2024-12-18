package flag

import "strings"

type CommaSeparatedValue []string

func (v *CommaSeparatedValue) String() string {
	if v == nil {
		return ""
	}
	return strings.Join(*v, ",")
}

func (v *CommaSeparatedValue) Set(s string) error {
	values := strings.Split(s, ",")
	*v = values
	return nil
}
