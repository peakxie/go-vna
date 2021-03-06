package vna

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

//
// Author: 陈永佳 chenyongjia@parkingwang.com, yoojiachen@gmail.com
// 初始化检测环境
//

const DataDirName = "data"

var (
	gProvinceNames = make(map[string]string)
	gCitiesNames   = make(map[string]string)
)

var (
	logger = log.New(os.Stderr, "GoVNA", log.Lshortfile)
)

// 初始化检测环境
func InitDetectorEnv(base string) {
	// 检查data目录是否存在
	if !fileExists(base) {
		logger.Println("Data directory is not exists, create now.")
		os.MkdirAll(base, os.ModePerm)
	}

	// 加载或者获取省份/字头数据
	initProvinces(base,
		"prov-army_v1.csv",
		"prov-civil_v1.csv",
		"prov-spec_v1.csv")

	// 加载
	initCities(base,
		"city-civil_v1.csv",
		"city-army_v2.csv",
		"city-contries_v1.csv",
		"city-spec_v1.csv",
		"city-wj_v1.csv")
}

func initProvinces(base string, names ...string) {
	for _, name := range names {
		path := filepath.Join(base, name)
		logger.Println("Loading provinces data file: ", path)
		downloadIfNotExists(path, name)
		loadFileToMemory(path, gProvinceNames)
	}
}

func initCities(base string, names ...string) {
	for _, name := range names {
		path := filepath.Join(base, name)
		logger.Println("Loading cities data file: ", path)
		downloadIfNotExists(path, name)
		loadFileToMemory(path, gCitiesNames)
	}
}

func loadFileToMemory(path string, targetMap map[string]string) {
	pairs, err := ReadRecords(path)
	if nil != err {
		panic(err)
	}

	for _, kv := range pairs {
		targetMap[kv.Key] = kv.Value
	}
}

func downloadIfNotExists(path string, name string) {
	// 检查文件是否存在
	if !fileExists(path) {
		logger.Println("Data file is not exist, download from github server. file:", name)
		// 如果不存在，从GigHub中下载
		resp, he := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/parkingwang/go-vna/master/data/%s", name))
		if nil != he {
			logger.Println("Cannot download data file from github server:", name)
			panic(he)
		}

		f, fe := os.Create(path)
		if nil != fe {
			logger.Println("Cannot create data file:", path)
			panic(fe)
		}

		io.Copy(f, resp.Body)
	}
}

func fileExists(path string) bool {
	_, e := os.Stat(path)
	if nil != e && os.IsNotExist(e) {
		return false
	} else {
		return true
	}
}
