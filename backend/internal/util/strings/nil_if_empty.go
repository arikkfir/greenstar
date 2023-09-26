package strings

func NilIfEmpty(s string) *string {
	if s == "" {
		return nil
	} else {
		return &s
	}
}

func EmptyIfNil(s any) string {
	if s == nil {
		return ""
	} else {
		return s.(string)
	}
}
