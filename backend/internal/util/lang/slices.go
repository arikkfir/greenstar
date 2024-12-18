package lang

func SliceFrom(startIndex int, args []string) []string {
	if len(args) <= startIndex {
		return []string{}
	} else {
		return args[startIndex:]
	}
}
