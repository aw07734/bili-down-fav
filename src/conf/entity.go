package conf

type Fav struct {
	Folder  string `json:"folder"`  // 收藏夹名称
	OutDir  string `json:"out_dir"` // 输出目录
	Flatmap bool   `json:"flatmap"` // 是否平铺(true: 不按投稿人分目录, false: 按投稿人分目录)
}

type Log struct {
	Level   string `json:"level"`   // 日志级别
	Success string `json:"success"` // 成功日志路径
	Fail    string `json:"fail"`    // 失败日志路径
}

type Config struct {
	MultiDownload int    `json:"multi_download"` // 多线程下载数量
	RemoveTemp    bool   `json:"remove_temp"`    // 下载完成后是否删除临时文件
	JsonDir       string `json:"json_dir"`       // JSON文件存放目录
	TempDir       string `json:"temp_dir"`       // 临时文件存放目录
	Favs          []Fav  `json:"favs"`           // 收藏夹列表
	Log           Log    `json:"log"`            // 日志配置
}
