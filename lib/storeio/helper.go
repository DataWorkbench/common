package storeio

const (
	filePathPrefix = "/dataomnis/resource"
)

// GenerateWorkspaceDir generate the directory path of specified workspace.
func GenerateWorkspaceDir(spaceId string) string {
	if spaceId == "" {
		panic("GenerateWorkspaceDir: spaceId cannot be empty")
	}
	return filePathPrefix + "/" + spaceId
}

// GenerateResourceFileDir generate the root directory of resource file.
func GenerateResourceFileDir(spaceId, fileId string) string {
	if spaceId == "" {
		panic("GenerateResourceFileDir: spaceId cannot be empty")
	}
	if fileId == "" {
		panic("GenerateResourceFileDir: fileId cannot be empty")
	}
	return filePathPrefix + "/" + spaceId + "/" + fileId
}

// GenerateResourceFilePath generate the file path of resource file.
func GenerateResourceFilePath(spaceId, fileId, version string) string {
	if spaceId == "" {
		panic("GenerateResourceFilePath: spaceId cannot be empty")
	}
	if fileId == "" {
		panic("GenerateResourceFilePath: fileId cannot be empty")
	}
	if version == "" {
		panic("GenerateResourceFilePath: version cannot be empty")
	}
	return filePathPrefix + "/" + spaceId + "/" + fileId + "/" + version
}
