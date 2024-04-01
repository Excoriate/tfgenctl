package generate

import "strings"

// ConvertNameToCanonical converts a name to canonical format.
// It replaces all '-' with '_'.
func ConvertNameToCanonical(name string) string {
	if name == "" {
		return ""
	}

	return strings.ReplaceAll(name, "-", "_")
}

func CleanBeforeCopy(relPath string) string {
	if strings.HasSuffix(relPath, ".tmpl") {
		relPath = strings.TrimSuffix(relPath, ".tmpl")
	}

	return relPath
}
