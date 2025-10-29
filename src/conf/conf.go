package conf

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/ini.v1"
)

const (
	IniPath = "assets/config.ini"
	QrPath  = "assets/qr.png"
)

var ExecDir, _ = getExecDir()
var file, _ = ini.Load(filepath.Join(ExecDir, IniPath))

func List(section string) map[string]string {
	return file.Section(section).KeysHash()
}

func Get(section string, key string) string {
	return file.Section(section).Key(key).String()
}

func Save(section string, key string, value string) error {
	file.Section(section).Key(key).SetValue(value)
	return file.SaveTo(IniPath)
}

func getExecDir() (string, error) {
	// 1. 获取当前可执行文件的完整路径
	// 例如: /path/to/your/project/myapp 或 C:\path\to\your\project\myapp.exe
	execPath, err := os.Executable()
	if err != nil {
		fmt.Println("Error getting executable path:", err)
		return "", err
	}

	// 2. 使用 filepath.EvalSymlinks 来解析符号链接，获取真实路径
	//    这对于开发环境（如使用 `go run`）尤其重要，因为 `go run` 会在临时目录创建一个可执行文件
	// 如果解析失败，就使用原始路径
	// realPath, err := filepath.EvalSymlinks(execPath)
	// if err != nil {
	// 	realPath = execPath
	// }

	// 3. 使用 filepath.Dir() 获取可执行文件所在的目录
	// 例如: /path/to/your/project 或 C:\path\to\your\project
	execDir := filepath.Dir(execPath)

	return execDir, nil
}
