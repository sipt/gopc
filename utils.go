package gopc

import (
	"io/ioutil"
	"os"
	"strings"
)

func isDir(filename string) bool {
	fi, err := os.Stat(filename)
	return err == nil && fi.IsDir()
}

func exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

//ListDir : 获取指定目录下的所有目录，可以排除指定目录
func walkDir(dirPth string, suffixes []string) (files []string, err error) {
	for i := 0; i < len(suffixes); i++ {
		suffixes[i] = strings.ToUpper(suffixes[i]) //忽略后缀匹配的大小写
	}
	return listDir(dirPth, suffixes)
}

//获取指定目录下的所有文件，不进入下一级目录搜索，可以匹配后缀过滤。
func listDir(dirPth string, suffixes []string) (files []string, err error) {
	files = make([]string, 0, 10)
	dir, err := ioutil.ReadDir(dirPth)
	if err != nil {
		return nil, err
	}
	PthSep := string(os.PathSeparator)
	allow := true
	for _, fi := range dir {
		if fi.IsDir() { // 忽略目录
			for _, suffix := range suffixes {
				if !strings.HasSuffix(strings.ToUpper(fi.Name()), suffix) { //匹配文件
					allow = false
				}
			}
			if allow {
				subDir := dirPth + PthSep + fi.Name()
				files = append(files, subDir)
				fs, err := listDir(subDir, suffixes)
				if err != nil {
					return nil, err
				}
				files = append(files, fs...)
			}
		}
	}
	return files, nil
}

//浅考备map
func copyMap(src map[string]interface{}) map[string]interface{} {
	dst := make(map[string]interface{}, len(src))
	for k, v := range src {
		dst[k] = v
	}
	return dst
}
