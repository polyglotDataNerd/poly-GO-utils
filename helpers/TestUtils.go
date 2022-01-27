package helpers

import (
	log "github.com/polyglotDataNerd/poly-Go-utils/utils"
	"path"
	"runtime"
)

func GetTestDir() (pathDir string, err error) {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		log.Error.Fatal("no such file or directory")
		return "", err
	}
	pathDir = path.Join(path.Dir(fileName), "..", "testdata")
	return pathDir, nil
}
