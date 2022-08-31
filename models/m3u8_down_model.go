package models

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"errors"
	"fmt"
	"go_m3u8_down/conf"
	"go_m3u8_down/utils"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type M3u8DownModel struct {
	Name            string   `json:"name" form:"name"`                           // idnex 名称
	FromWebUrl      string   `json:"from_web_url" form:"from_web_url"`           // 来源网站
	Host            string   `json:"host" form:"host"`                           // 域名
	Url             string   `json:"url" form:"url"`                             // idnex url
	Body            string   `json:"body" form:"body"`                           // idnex 获取的内容
	Key             string   `json:"key" form:"key"`                             // key 时 加密key
	TsList          []string `json:"ts_list" form:"ts_list"`                     // ts 列表
	TsFileList      []string `json:"ts_file_list" form:"ts_file_list"`           //	ts 文件列表
	TsLen           int      `json:"ts_len" form:"ts_len"`                       //ts 个数
	Referer         string   `json:"referer" form:"referer"`                     //访问 Referer
	UserAgent       string   `json:"user_agent" form:"user_agent"`               //访问 UserAgent
	DownloadChan    int      `json:"download_chan" form:"download_chan"`         //下载 指定协程数量
	DownloadDir     string   `json:"download_dir" form:"download_dir"`           //下载文件夹
	DownloadTsCount int      `json:"download_ts_count" form:"download_ts_count"` //已经下载次数
	TSRetry         int      `json:"ts_retry" form:"ts_retry"`                   //下载重试次数ts
	Progress        string   `json:"progress" form:"progress"`                   //下载进度
	HashValue       uint32   `json:"hash_value" form:"hash_value"`               //idnex hash
	VideoFile       string   `json:"video_file" form:"video_file"`
	Status          string   `json:"status" form:"status"` //下载状态 0 下载中 1已完成
}

var mu sync.RWMutex

func CacheGetAllHash() []*M3u8DownModel {
	//所有 idnex 集合
	var arrayM3u8 []*M3u8DownModel
	//缓存 key 值
	var allHashKey []string

	if err := conf.NdbGetAny(conf.AllHash, &allHashKey); err == nil {
		if len(allHashKey) > 0 {
			for _, key := range allHashKey {
				if m3u8, found := CacheGetM3u8ByHash(key); found {
					arrayM3u8 = append(arrayM3u8, m3u8)
				} else {
					go CacheRemoveByHash(key)
				}
			}
		}
		return arrayM3u8
	}
	return nil
}

// CacheSetNewM3u8 设置缓存
func CacheSetNewM3u8(m *M3u8DownModel) {
	m.Status = "0"
	hashStr := strconv.FormatUint(uint64(m.HashValue), 10)
	conf.NdbPutAny(hashStr, m)
	CacheSetAllHash(hashStr)
}
func CacheSetM3u8(m *M3u8DownModel) {
	conf.NdbPutAny(strconv.Itoa(int(m.HashValue)), m)
}
func CacheSetAllHash(newHashStr string) {
	mu.Lock()
	defer mu.Unlock()
	// 获取一个 存放所有hash 的数组
	var arrayAllHash []string
	if err := conf.NdbGetAny(conf.AllHash, &arrayAllHash); err == nil {
		//判断是否已存在
		for _, hash := range arrayAllHash {
			if hash == newHashStr {
				// 存在 直接跳出
				return
			}
		}
		newArray := append(arrayAllHash, newHashStr)
		//新的 数组 存入
		conf.NdbPutAny(conf.AllHash, &newArray)
	}
}

func CacheRemoveByHash(hash string) {
	mu.Lock()
	defer mu.Unlock()
	conf.NdbDel(hash)
	//删除 缓存中
	//缓存 key 值
	var allHashKey []string
	if err := conf.NdbGetAny(conf.AllHash, &allHashKey); err == nil {
		//删除切片中 hash 值
		newArray := utils.DelSliceByStr(allHashKey, hash)
		if newArray != nil {
			//新的 数组 存入
			conf.NdbPutAny(conf.AllHash, &newArray)
		}
	}

}

func CacheGetM3u8ByHash(HashValue string) (*M3u8DownModel, bool) {
	var m *M3u8DownModel
	if err := conf.NdbGetAny(HashValue, &m); err == nil && m != nil {
		return m, true
	}
	return m, false
}

// InitHash 初始化 hash 存文件夹
func (m *M3u8DownModel) InitHash() {
	m.HashValue = utils.Hash(m.Url)
}

func (m *M3u8DownModel) JsonStr() string {
	jsonP, _ := json.Marshal(&m)
	return string(jsonP)
}

func (m *M3u8DownModel) DelTs() {
	if m.TsFileList != nil {
		os.RemoveAll(m.DownloadDir)
	}
}

// MergeTs 合并 多个 ts 文件
func (m *M3u8DownModel) MergeTs() error {
	//新建一个合并文件
	mp4FilePath := m.GetMp4File()
	//删除旧文件
	if exists, _ := PathExists(mp4FilePath); exists {
		os.Remove(mp4FilePath)
	}
	m.VideoFile = m.GetVideoFile()
	m.Status = "1"
	CacheSetM3u8(m)
	mp4File, err := os.Create(mp4FilePath)
	defer func() {
		mp4File.Close()
		if r := recover(); r != nil {
			log.Println("合并 多个 ts 文件 捕获异常:", r)
			os.Remove(mp4FilePath)
		}
	}()
	defer mp4File.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	//将ts 文件以名称排序
	sort.Strings(m.TsFileList)
	//遍历所有ts 文件
	for i, filePath := range m.TsFileList {
		log.Printf("mergeTs file  %s :  %5.2f %%    \n", filePath, (float32(i+1)/float32(m.TsLen))*100)
		tsFile, err := os.ReadFile(filePath)
		if err != nil {
			log.Println(err)
			return err
		}
		_, err = mp4File.Write(tsFile)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	return nil
}

// StartDownloader 下载ts切割文件到download_dir
func (m *M3u8DownModel) StartDownloader() {
	var wg sync.WaitGroup
	//下载线程数 m
	limiter := make(chan struct{}, m.DownloadChan)
	wg.Add(m.TsLen)
	for i, ts := range m.TsList {
		//协程数量保持在 m.DownloadChan
		limiter <- struct{}{}
		go m.goDownloadFile(&wg, limiter, ts, i)
	}
	wg.Wait()
}
func (m *M3u8DownModel) goDownloadFile(wg *sync.WaitGroup, limiter chan struct{}, ts string, indexTs int) {
	defer func() {
		wg.Done()
		<-limiter
	}()
	//下载ts切割文件到download_dir
	err := m.DownloadTsFile(ts, indexTs, m.TSRetry)
	if err != nil {
		log.Printf("DownloadTsFile 下载异常:  %v", err)
		return
	}
	m.DownloadTsCount++
	m.Progress = fmt.Sprintf("%5.2f %%", (float32(m.DownloadTsCount)/float32(m.TsLen))*100)
	//更新 下载进度
	CacheSetM3u8(m)
	log.Printf("Downloading  %s  %s  \n", fmt.Sprint(indexTs)+".ts", m.Progress)
}

// PathExists 判断文件是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// DownloadTsFile 下载ts文件
// @modify: 2020-08-13 修复ts格式SyncByte合并不能播放问题
func (m *M3u8DownModel) DownloadTsFile(ts string, indexTs int, retries int) error {
	//判断当前文件是否存在  格式化 int，位数不够0补齐
	tsFilePath := filepath.Join(m.DownloadDir, fmt.Sprintf("%07d", indexTs)+".ts")
	defer func() {
		if r := recover(); r != nil {
			log.Printf("下载ts文件 %s ; 循环次数：%d 捕获异常: %v  \n", ts, retries, r)
			os.Remove(tsFilePath)
			if retries > 0 {
				m.DownloadTsFile(ts, indexTs, retries-1)
			}
		}
	}()
	//判断是否有删除指令
	if _, found := CacheGetM3u8ByHash(strconv.Itoa(int(m.HashValue))); !found {
		return nil
	}
	if isExist, _ := PathExists(tsFilePath); isExist {
		m.TsFileList = append(m.TsFileList, tsFilePath)
		return nil
	}
	req, _ := http.NewRequest(http.MethodGet, ts, nil)
	req.Header.Add("user-agent", m.UserAgent)
	//发起 请求 响应
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		if retries > 0 {
			return m.DownloadTsFile(ts, indexTs, retries-1)
		} else {
			//logger.Printf("[warn] File :%s", ts.Url)
			return err
		}
	}
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()
	// 校验长度是否合法
	var origData []byte

	// body 正确响应
	origData, err = ioutil.ReadAll(resp.Body)
	if len(origData) == 0 || resp.ContentLength == 0 {
		//logger.Println("[warn] File: " + ts.Name + "res origData invalid or err：", res.Error)
		return m.DownloadTsFile(ts, indexTs, retries-1)

	}
	// 解密出视频 ts 源文件
	if m.Key != "" {
		//解密 ts 文件，算法：aes 128 cbc pack5
		origData, err = AesDecrypt(origData, []byte(m.Key))
		if err != nil {
			return m.DownloadTsFile(ts, indexTs, retries-1)

		}
	}
	syncByte := uint8(71) //0x47
	bLen := len(origData)
	for j := 0; j < bLen; j++ {
		if origData[j] == syncByte {
			origData = origData[j:]
			break
		}
	}
	err = ioutil.WriteFile(tsFilePath, origData, 0666)
	if err != nil {
		//删除 等待重试
		os.Remove(tsFilePath)
		return m.DownloadTsFile(ts, indexTs, retries-1)
	}
	//存储 ts 文件目录
	m.TsFileList = append(m.TsFileList, tsFilePath)
	return nil
}

// InitDir 初始化 下载文件夹
func (m *M3u8DownModel) InitDir() {
	m.DownloadDir = filepath.Join(m.DownloadDir, strconv.Itoa(int(m.HashValue)))
	if isExist, _ := PathExists(m.DownloadDir); !isExist {
		os.MkdirAll(m.DownloadDir, os.ModePerm)
	}
}

// InitKeyAndTs 初始化 ts 列表 及 加密 key
func (m *M3u8DownModel) InitKeyAndTs() error {
	lines := strings.Split(m.Body, "\n")
	var m3u8TsSlice []string
	for _, line := range lines {
		if strings.Contains(line, "#EXT-X-KEY") {
			uriPos := strings.Index(line, "URI")
			quotationMarkPos := strings.LastIndex(line, "\"")
			keyUrl := strings.Split(line[uriPos:quotationMarkPos], "\"")[1]
			if !strings.Contains(line, "http") {
				keyUrl = m.Host + keyUrl
			}
			req, _ := http.NewRequest(http.MethodGet, keyUrl, nil)
			req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.81 Safari/537.36 Edg/104.0.1293.54")
			//发起 请求 响应
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				return err
			}
			if res.StatusCode == http.StatusOK {
				// body正确响应
				body, _ := ioutil.ReadAll(res.Body)
				m.Key = string(body)
			}
		} else if strings.HasPrefix(line, "http") {
			m3u8TsSlice = append(m3u8TsSlice, line)
		} else if strings.HasPrefix(line, "/") {
			newUrl := m.Host + line
			m3u8TsSlice = append(m3u8TsSlice, newUrl)
		}
	}
	m.TsLen = len(m3u8TsSlice)
	m.TsList = m3u8TsSlice
	return nil
}

// InitHost 获取m3u8地址的host
func (m *M3u8DownModel) InitHost() error {
	uri, err := url.Parse(m.Url)
	if err != nil {
		return err
	}
	m.Referer = uri.Scheme + "://" + uri.Host + filepath.Dir(uri.EscapedPath())
	m.Host = uri.Scheme + "://" + uri.Host
	return nil
}

// CheckJump 检测是否需要二次跳转
func (m *M3u8DownModel) CheckJump() error {
	uri, err := url.Parse(m.Url)
	if err != nil {
		return err
	}
	req, _ := http.NewRequest(http.MethodGet, m.Url, nil)
	req.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.81 Safari/537.36 Edg/104.0.1293.54")
	//发起 请求 响应
	resp, err := http.DefaultClient.Do(req)
	if http.StatusOK == resp.StatusCode {
		//逐行读取 文本
		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			//判断 是否为 idnex
			if strings.HasPrefix(line, "/") || (strings.HasSuffix(line, ".m3u8") || strings.HasSuffix(line, ".m3u")) {
				m.Url = uri.Scheme + "://" + uri.Host + line
				return nil
			}
		}
	} else if resp != nil {
		return errors.New(resp.Status)
	}
	return errors.New("网络访问错误")
}

func (m *M3u8DownModel) InitM3u8Body() error {
	req, _ := http.NewRequest(http.MethodGet, m.Url, nil)
	req.Header.Add("user-agent", m.UserAgent)
	req.Header.Add("referer", m.Referer)
	//发起 请求 响应
	resp, _ := http.DefaultClient.Do(req)
	if http.StatusOK == resp.StatusCode {
		// body 正确响应
		body, err := ioutil.ReadAll(resp.Body)
		m.Body = string(body)
		return err
	} else if resp != nil {
		return errors.New(resp.Status)
	}
	return errors.New("网络访问错误")
}

func (m *M3u8DownModel) GetMp4File() string {
	return filepath.Join(conf.DownDir, m.Name+".mp4")
}
func (m *M3u8DownModel) GetVideoFile() string {
	return filepath.Join(conf.VideoDir, m.Name+".mp4")
}

func AesDecrypt(crypted, key []byte, ivs ...[]byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	var iv []byte
	if len(ivs) == 0 {
		iv = key
	} else {
		iv = ivs[0]
	}
	blockMode := cipher.NewCBCDecrypter(block, iv[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS7UnPadding(origData)
	return origData, nil
}

func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
