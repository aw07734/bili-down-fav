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
	"sync"
)

func main() {
	curUser := user.CurrentUser().Data
	if curUser.IsLogin {
		fmt.Println("当前用户：" + curUser.Uname)
		threads := conf.ConfigInfo.MultiDownload
		if threads <= 0 {
			fmt.Println("线程数必须大于0")
		}
		folders := fav.ListFavFolder(curUser.Mid)
		ctx, cancel := context.WithCancel(context.Background())
		for index := range conf.ConfigInfo.Favs {
			currentFav := conf.ConfigInfo.Favs[index]

			fmt.Println("当前收藏夹：" + currentFav.Folder)

			favs := fav.ListForDownloads(folders, currentFav.Folder)
			if len(favs) == 0 {
				fmt.Println("当前收藏夹无视频可下载")
				continue
			}
			fmt.Println("视频数量：", len(favs))
			go util.Log(ctx)

			mid2Name := outdir.MakeOutDirCache(currentFav.OutDir)
			if threads > 0 {
				// 限制并发数量
				ch := make(chan any, threads)
				// 多线程下载
				wg := new(sync.WaitGroup)
				for index := range favs {
					v := favs[index]
					wg.Add(1)
					ch <- v
					bvid := v.Bvid
					go func() {
						video.Download(bvid, mid2Name, &currentFav)
						wg.Done()
						<-ch
					}()
				}
				wg.Wait()
			} else {
				fmt.Println("线程数必须大于0")
			}
		}

		cancel()
	} else {
		fmt.Println("当前未登录，请登录")
		user.Login(func() {
			fmt.Println("登录")
			main()
		})
	}
}
