package fav

import (
	"bili-down-fav/src/conf"
	"bili-down-fav/src/util"
	"encoding/json"
	"strconv"
)

func ListFavFolder(mid int) *FolderResp {
	url := "https://api.bilibili.com/x/v3/fav/folder/created/list-all"
	folders := new(FolderResp)
	if resp, err := util.Http.R().
		SetQueryParam("up_mid", strconv.Itoa(mid)).
		Get(url); err == nil {
		json.Unmarshal(resp.Body(), folders)
	}
	return folders
}

func filterFolder(folders *FolderResp, dFolder string) int {
	for _, v := range folders.Data.List {
		if v.Title == dFolder {
			return v.Id
		}
	}
	return 0
}

func ListForDownloads(folders *FolderResp) []Media {
	mlid := filterFolder(folders, conf.Get("fav", "d_folder"))
	if mlid == 0 {
		return []Media{}
	}
	url := "https://api.bilibili.com/x/v3/fav/resource/list"
	content := new(ContentResp)
	var array []Media
	for i := 1; ; i++ {
		if resp, err := util.Http.R().
			SetQueryParam("media_id", strconv.Itoa(mlid)).
			SetQueryParam("ps", strconv.Itoa(20)).
			SetQueryParam("pn", strconv.Itoa(i)).
			Get(url); err == nil {
			json.Unmarshal(resp.Body(), content)
			for _, v := range content.Data.Medias {
				array = append(array, v)
			}
			if !content.Data.HasMore {
				break
			}
		} else {
			break
		}
	}
	return array
}
