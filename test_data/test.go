package test_data

import (
	"go/build"
	"io/ioutil"
	"os"
	"strings"

	"github.com/sipt/gopc"
)

//LineInfo : 代码行信息
type LineInfo struct {
	ImportPath string
	FileName   string
	FilePath   string
	FilterMap  map[string]interface{}
}

func (l *LineInfo) Clear() {
	l.ImportPath = ""
	l.FileName = ""
	l.FilePath = ""
	l.ClearFilterMap()
}

func (l *LineInfo) ClearFilterMap() {
	for k := range l.FilterMap {
		delete(l.FilterMap, k)
	}
}

//Compiler : go简单语法编译器
type Compiler struct {
	filterManager *gopc.FilterManager
}

//NewCompiler : 初始化Compilter
func NewCompiler() *Compiler {
	return &Compiler{
		filterManager: gopc.NewFilterManager(),
	}
}

//AddFilter : 添加filter
func (c *Compiler) AddFilter(filter gopc.Filter) {
	c.filterManager.AddFilter(filter)
}

//Entrance 入口函数
func (c *Compiler) Entrance(path string, suffixes []string) error {
	var err error
	if isDir(path) {
		//遍历出目录下所有文件
		fs, err := walkDir(path, suffixes)
		if err != nil {
			return err
		}
		err = c.walkFile(fs...)
	} else {
		//处理文件
		err = c.walkFile(path)
	}
	return err
}

func (c *Compiler) walkFile(pathes ...string) error {
	PthSep := string(os.PathSeparator)
	for _, path := range pathes {
		pkg, err := build.ImportDir(path, 0)
		if err != nil {
			return err
		}
		lineInfo := new(gopc.LineInfo)
		for _, file := range pkg.GoFiles {
			lineInfo.Clear()
			lineInfo.FilePath = pkg.Dir + PthSep + file
			lineInfo.FileName = file
			lineInfo.ImportPath = pkg.ImportPath
			c.handleFile(lineInfo)
		}
	}
	return nil
}

func (c *Compiler) handleFile(l *gopc.LineInfo) error {
	bytes, err := ioutil.ReadFile(l.FileName)
	if err != nil {
		return err
	}
	lineDatas := strings.Split(string(bytes), "\n")
	for i, line := range lineDatas {
		line, err = c.filterManager.Start(l, line, c.filterManager)
		if err != nil {
			return err
		}
		lineDatas[i] = line
	}
	data := strings.Join(lineDatas, "\n")
	err = ioutil.WriteFile(l.FileName, []byte(data), 0666)
	return err
}

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
