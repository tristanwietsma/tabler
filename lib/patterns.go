package lib

import "regexp"

var (
	foreignKeyPattern = regexp.MustCompile("([A-Za-z][A-Za-z0-9]*)ID")
	tagPattern        = regexp.MustCompile(`tabler:"([0-9a-zA-Z=&\(\)]*)"`)
)
