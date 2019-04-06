package stredit

// Pluralize uses the check count to see if a string should be pluralized or not
func Pluralize(checkCount uint) string {
	if checkCount == 1 {
		return ""
	}
	return "s"
}
