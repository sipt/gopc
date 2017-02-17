package gopc

import (
	"go/build"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
	filterManager *FilterManager
}

//NewCompiler : 初始化Compilter
func NewCompiler() *Compiler {
	return &Compiler{
		filterManager: NewFilterManager(),
	}
}

//AddFilter : 添加filter
func (c *Compiler) AddFilter(filter Filter) {
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

	lineInfo := new(LineInfo)
	lineInfo.FilterMap = make(map[string]interface{})
	for _, path := range pathes {
		if isDir(path) {
			pkg, err := build.ImportDir(path, 0)
			if err != nil {
				return err
			}
			for _, file := range pkg.GoFiles {
				lineInfo.Clear()
				lineInfo.FilePath = pkg.Dir + PthSep + file
				lineInfo.FileName = file
				lineInfo.ImportPath = pkg.ImportPath
				c.handleFile(lineInfo)
			}
		} else {
			dir := filepath.Dir(path)
			pkg, err := build.ImportDir(dir, 0)
			if err != nil {
				return err
			}
			lineInfo.FilePath = path
			lineInfo.FileName = path[len(dir)+1:]
			lineInfo.ImportPath = pkg.ImportPath
			c.handleFile(lineInfo)
		}
	}
	return nil
}

func (c *Compiler) handleFile(l *LineInfo) error {
	bytes, err := ioutil.ReadFile(l.FilePath)
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
	err = ioutil.WriteFile(l.FilePath, []byte(data), 0666)
	return err
}
