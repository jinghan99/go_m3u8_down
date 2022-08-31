package service

import (
	"errors"
	"go_m3u8_down/models"
	"log"
	"os"
)

func AllM3u8Down() []*models.M3u8DownModel {
	return models.CacheGetAllHash()
}

// Runs 开始 m3u8下载
func Runs(dirPath string, m3u8Url string, name string) (*models.M3u8DownModel, error) {
	m := &models.M3u8DownModel{
		Name:            name,
		Url:             m3u8Url,
		UserAgent:       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.81 Safari/537.36 Edg/104.0.1293.54",
		DownloadTsCount: 0,
		TSRetry:         5,
		DownloadChan:    8,
		DownloadDir:     dirPath,
	}
	defer func() {
		if r := recover(); r != nil {
			log.Println("Runs 捕获异常 下载状态为 异常停止:", r)
		} else {
			log.Printf("Runs 下载 ：%s  启动协程下载 \n", name)
		}
	}()
	m.InitHash()
	//// 1、检测是否需要二次跳转
	err := m.CheckJump()
	if err != nil {
		return nil, err
	}
	//2、获取m3u8地址的host
	err = m.InitHost()
	if err != nil {
		return m, err
	}
	err = m.InitM3u8Body()
	if err != nil {
		return m, err
	}
	//3、初始化 ts 列表 及 加密 key
	err = m.InitKeyAndTs()
	//4、 初始化 下载文件夹
	m.InitDir()
	log.Printf("初始化完成  %v", m.Name)
	models.CacheSetNewM3u8(m)
	go goRuns(m)
	return m, nil
}

func goRuns(m *models.M3u8DownModel) (*models.M3u8DownModel, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("goRuns 捕获异常 下载状态为 异常停止:", r)
		}
	}()
	//5、开启下载
	m.StartDownloader()
	log.Printf("下载完成：%s \n", m.Name)
	//判断下载数量 是否正确
	if m.DownloadTsCount == m.TsLen {
		// 合并 多个 ts 文件
		err := m.MergeTs()
		log.Printf("合并完成：%s   \n", m.Name)
		if err != nil {
			log.Printf("合并失败")
			return m, err
		}
		//删除ts 文件
		m.DelTs()
	}
	return nil, nil
}

func DelM3u8ByHash(hashValue string) error {
	//先删除文件夹
	m, found := models.CacheGetM3u8ByHash(hashValue)
	if found {
		//删除文件夹ts
		if exists, _ := models.PathExists(m.DownloadDir); exists {
			os.RemoveAll(m.DownloadDir)
		}
		//删除mp4
		if exists, _ := models.PathExists(m.GetMp4File()); exists {
			os.RemoveAll(m.GetMp4File())
		}
		models.CacheRemoveByHash(hashValue)
	}
	return errors.New("文件已删除")
}

func ReDownM3u8(hashValue string) error {
	//先删除文件夹
	m, found := models.CacheGetM3u8ByHash(hashValue)
	if found {
		go goRuns(m)
	}
	return errors.New("文件已删除")
}
