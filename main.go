package main

import (
	"bili-down-fav/src/bili/fav"
	"bili-down-fav/src/bili/outdir"
	"bili-down-fav/src/bili/user"
	"bili-down-fav/src/bili/video"
	"bili-down-fav/src/conf"
	"bili-down-fav/src/util"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

func main() {
	curUser := user.CurrentUser().Data
	if curUser.IsLogin {
		fmt.Println("当前用户：" + curUser.Uname)
		fmt.Println("当前收藏夹：" + conf.Get("fav", "d_folder"))

		favs := fav.ListForDownloads(fav.ListFavFolder(curUser.Mid))
		if len(favs) == 0 {
			fmt.Println("当前收藏夹无视频可下载")
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		fmt.Println("视频数量：", len(favs))
		go util.Log(ctx)

		mid2Name := outdir.MakeJson()
		threads, _ := strconv.Atoi(conf.Get("file", "multi_download"))
		if threads > 0 {
			// 限制并发数量
			ch := make(chan any, threads)
			// 多线程下载
			wg := new(sync.WaitGroup)
			for _, v := range favs {
				wg.Add(1)
				ch <- v
				bvid := v.Bvid
				go func() {
					video.Download(bvid, mid2Name)
					wg.Done()
					<-ch
				}()
			}
			wg.Wait()
		} else {
			fmt.Println("线程数必须大于0")
		}
		cancel()
	} else {
		fmt.Println("当前未登录，请登录")
		user.Login(func() {
			os.Remove(filepath.Join(conf.ExecDir, conf.QrPath))
			fmt.Println("登录")
			main()
		})
	}
}
