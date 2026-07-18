# bilibili收藏夹下载器
一个能快速下载 bilibili 原始视频的脚本。

go 语言编写，支持 linux，方便归档在服务器中。

## 用法
- 将 assets 目录放在可执行文件**同级目录**，然后执行即可
- assets目录中的 config.json.temp 请自行参考，使用时可将文件**重命名为 config.json**
- 首次执行手机 bilibili 扫一下登录就好，之后他自己会保存 cookie 信息，除非强制下线或者过期，通常下次就不需要再扫了
- config.json 中的字段含义参考 `src/conf/entity.go`

## ffmpeg
- 在linux中使用请**保证已安装ffmpeg**，我的版本是5.1.4，理论上大于这个即可
- 在windows中就去官网下载ffmpeg，然后将其添加到环境变量
