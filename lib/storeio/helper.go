package storeio

import (
	"fmt"
)

const (
	filePathPrefix = "/dataomnis/resource"
)

// GenerateFileRootDir generate the root dir of resource file.
func GenerateFileRootDir(spaceId string) string {
	return fmt.Sprintf("%s/%s", filePathPrefix, spaceId)
}

// GenerateFilePath to generate a path that store the file.
func GenerateFilePath(spaceId string, fileId string, version string) string {
	return fmt.Sprintf("%s/%s.%s", GenerateFileRootDir(spaceId), fileId, version)
}
