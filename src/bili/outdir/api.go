package outdir

import (
	"bili-down-fav/src/conf"
	"os"
	"path/filepath"
	"regexp"
)

func MakeJson() map[string]string {
	context := filepath.Join(conf.ExecDir, conf.Get("file", "out_dir"))
	files, err := os.ReadDir(context)
	if err != nil {
		return nil
	}
	pattern := `\((\d+)\)`

	// 2. 编译正则表达式
	// regexp.Compile 在编译失败时会返回一个 error
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil
	}
	mid2Name := make(map[string]string)
	for _, file := range files {
		fileName := file.Name()
		matches := re.FindStringSubmatch(fileName)
		if len(matches) > 1 {
			mid := matches[1]
			mid2Name[mid] = fileName
		}
	}
	return mid2Name
}
