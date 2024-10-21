package utils

import (
	"path"
	"strings"
)

func BuildObjectKey(dirs []string, objectName, extension string) string {

	fullPath := path.Join(append(dirs, objectName)...)

	if extension != "" {
		fullPath += "." + strings.TrimPrefix(extension, ".")
	}

	return fullPath
}