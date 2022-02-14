package constants

import (
	"fmt"
)

const (
	resourceFilePathPrefix = "/dataomnis/resource"
)

// GenResourceFileRootDir generate the root dir of resource file.
func GenResourceFileRootDir(spaceId string) string {
	return fmt.Sprintf("%s/%s", resourceFilePathPrefix, spaceId)
}

// GenResourceFilePath to generate a path that store the file.
func GenResourceFilePath(spaceId string, resourceId string, version string) string {
	return fmt.Sprintf("%s/%s.%s", GenResourceFileRootDir(spaceId), resourceId, version)
}
