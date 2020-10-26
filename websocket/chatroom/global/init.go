package global

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	Init()
}

var RootDir string

var once = new(sync.Once)

func Init() {
	once.Do(func() {
		inferRootDir()
		initConfig()
	})
}

// inferRootDir 推断出项目根目录
func inferRootDir() {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	fmt.Println("os.Getwd()", cwd)
	var infer func(d string) string
	infer = func(d string) string {
		// 这里要确保项目根目录下存在 template 目录
		if exists(d + "/template") {
			return d
		}

		return infer(filepath.Dir(d))
	}

	RootDir = infer(cwd)
	fmt.Println(RootDir)
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}
