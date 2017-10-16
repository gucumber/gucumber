package gucumber

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackageImportPathIsExtractedFromFilePath(t *testing.T) {
	//arrange
	packageName := "github.com/username/reponame"
	gopath := os.Getenv("GOPATH")
	for _, p := range filepath.SplitList(gopath) {
		file := filepath.Join(p, "src/github.com/username/reponame/main.go")
		//act
		importPath := assembleImportPath(file)
		//assert
		assert.Equal(t, packageName, importPath)
	}
}
