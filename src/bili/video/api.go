package video

import (
	"bili-down-fav/src/conf"
	"bili-down-fav/src/util"
	"encoding/json"
	"errors"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func videoInfo(bvid string) *Info {
	url := "https://api.bilibili.com/x/web-interface/view"
	info := new(Info)
	if resp, err := util.Http.R().
		SetQueryParam("bvid", bvid).
		Get(url); err == nil {
		json.Unmarshal(resp.Body(), info)
	}
	return info
}

func Download(bvid string, mid2Name map[string]string, currentFav *conf.Fav) {
	vinfo := videoInfo(bvid)
	if vinfo.Code != 0 {
		util.LogFail(bvid, vinfo.Message, nil)
		return
	}
	if len(vinfo.Data.Pages) > 1 {
		for index := range vinfo.Data.Pages {
			page := vinfo.Data.Pages[index]
			oFile, title := makeFile(bvid, vinfo, mid2Name, currentFav)
			download_page(bvid, oFile, title, page.Cid)
		}
	} else {
		oFile, title := makeFile(bvid, vinfo, mid2Name, currentFav)
		download_page(bvid, oFile, title, vinfo.Data.Cid)
	}
}

func makeFile(bvid string, vInfo *Info, mid2Name map[string]string, currentFav *conf.Fav) (string, string) {
	title := makeTitleLegal(vInfo.Data.Title) + "(" + bvid + ")"
	var relativePath string
	if currentFav.Flatmap {
		relativePath = title + ".mkv"
	} else {
		owner := vInfo.Data.Owner
		if name, ok := mid2Name[strconv.Itoa(owner.Mid)]; ok {
			relativePath = name + "/" + title + ".mkv"
		} else {
			relativePath = owner.Name + "(" + strconv.Itoa(owner.Mid) + ")" + "/" + title + ".mkv"
		}
	}
	return filepath.Join(conf.ExecDir, currentFav.OutDir, relativePath), title
}

func download_page(bvid, oFile, title string, cid int) {
	if fileExist(oFile) {
		log.Println("跳过:", oFile)
		return
	}
	url := "https://api.bilibili.com/x/player/playurl"
	info := new(PlayInfo)
	var fnval uint = 16 | 64 | 256 | 512 | 1024 | 2048
	if resp, err := util.Http.R().
		SetQueryParam("bvid", bvid).
		SetQueryParam("cid", strconv.Itoa(cid)).
		SetQueryParam("fnval", strconv.Itoa(int(fnval))).
		SetQueryParam("fourk", "1").
		SetResult(info).
		Get(url); err == nil {

		// 保存json文件
		jsonDir := filepath.Join(conf.ExecDir, conf.ConfigInfo.JsonDir)
		if jsonDir != "" {
			os.MkdirAll(jsonDir, os.ModePerm)
			os.WriteFile(jsonDir+"/"+bvid+"_video.json", resp.Body(), fs.ModePerm)
		}
		if info.Code != 0 {
			util.LogFail(bvid, title, info.Message)
			return
		}
		// log.Println("视频清晰度：", info.Data.AcceptDescription[0])

		// 获取最佳video
		if len(info.Data.Dash.Video) > 0 {
			dash := info.Data.Dash
			bandwidth := 0
			bestVideo := dash.Video[0]
			for _, v := range dash.Video {
				if v.Bandwidth > bandwidth {
					bestVideo = v
					bandwidth = v.Bandwidth
				}
			}
			vUrl := bestVideo.BaseUrl
			// log.Println("最佳视频带宽", bestVideo.Bandwidth)

			// 获取最佳音频
			var aUrl string
			var aBackUrl []string
			hasAudio := true
			if dash.Flac.Audio.Id > 0 {
				// log.Println("音频id：", dash.Flac.Audio.Id)
				aUrl = dash.Flac.Audio.BaseUrl
				aBackUrl = dash.Flac.Audio.BackupUrl
			} else if dash.Dolby.Audio != nil && dash.Dolby.Audio[0].Id > 0 {
				// log.Println("音频id：", dash.Dolby.Audio[0].Id)
				aUrl = dash.Dolby.Audio[0].BaseUrl
				aBackUrl = dash.Dolby.Audio[0].BackupUrl
			} else if dash.Audio != nil {
				// log.Println("音频id：", dash.Audio[0].Id)
				aUrl = dash.Audio[0].BaseUrl
				aBackUrl = dash.Audio[0].BackupUrl
			} else {
				hasAudio = false
			}
			toFile(bvid, oFile, title, vUrl, bestVideo.BackupUrl, hasAudio, aUrl, aBackUrl)
		} else if len(info.Data.Durl) > 0 {
			dual := info.Data.Durl[0]
			toFile(bvid, oFile, title, dual.URL, dual.BackupURL, false, "", nil)
		} else {
			util.LogFail(bvid, oFile, "无可用视频")
			// return
		}
	}
}

func toFile(bvid, oFile, title, vUrl string, vBackupUrl []string, hasAudio bool, aUrl string, aBackupUrl []string) {
	prefix := filepath.Join(conf.ExecDir, conf.ConfigInfo.TempDir)
	vFile := prefix + "/video_" + title + ".m4s"
	aFile := prefix + "/audio_" + title + ".m4s"

	// 下载视频
	if err := download(vUrl, vFile); err != nil {
		log.Println("视频下载失败，使用备用地址：", title)

		success := false
		for _, u := range vBackupUrl {
			if errb := download(u, vFile); errb == nil {
				success = true
				log.Println("备用地址下载成功：", title)
				break
			}
		}

		if !success {
			util.LogFail(bvid, title, err)
			return
		}
	}

	// 下载音频
	if hasAudio {
		if err := download(aUrl, aFile); err != nil {
			log.Println("音频下载失败，使用备用地址：", title)

			success := false
			for _, u := range aBackupUrl {
				if errb := download(u, aFile); errb == nil {
					success = true
					log.Println("备用地址下载成功：", title)
					break
				}
			}

			if !success {
				util.LogFail(bvid, title, err)
				return
			}
		}
	}

	// log.Println("混流")
	if hasAudio {
		util.Combine(vFile, aFile, oFile)
	} else {
		util.Convert(vFile, oFile)
	}
	if conf.ConfigInfo.RemoveTemp {
		// log.Println("移除临时文件")
		os.Remove(vFile)
		os.Remove(aFile)
	}
	log.Println("文件已下载至", oFile)
	util.LogSuccess(bvid, title)
}

func download(url, name string) error {
	// log.Println(name, "下载开始")
	if resp, err := util.Http.R().
		SetOutput(name).
		SetHeader("Referer", "https://www.bilibili.com").
		Get(url); err == nil {
		if resp.StatusCode() == 200 {
			// log.Println(name, "下载完成")
			return nil
		}
		return errors.New(resp.Status())
	} else {
		log.Println(name, "下载失败", err)
		return err
	}
}

func makeTitleLegal(str string) string {
	result := ""
	for _, v := range str {
		letter := string(v)
		switch letter {
		case "/", "\\", ":", "*", "?", "\"", "|":
			result += "_"
		case "<":
			result += "["
		case ">":
			result += "]"
		default:
			result += letter
		}
	}
	return result
}

func fileExist(filePath string) bool {
	if _, err := os.Stat(filePath); err == nil {
		return true
	}
	return false
}
