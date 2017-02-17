package gopc

import (
	"fmt"
	"testing"
)

func TestPrecompile(t *testing.T) {
	path := "/Users/sipt/Documents/GOPATH/src/github.com/sipt/gopc/test_data/test.go"
	c := NewCompiler()

	//addfilters
	c.AddFilter(NewLineNumFilter())
	c.AddFilter(NewLexicalFilter())
	c.AddFilter(NewTestFilter())
	err := c.Entrance(path, []string{".git"})
	if err != nil {
		panic(err)
	}

}

//TestFilter : 过滤器
type TestFilter struct {
}

//NewTestFilter 新建过滤器
func NewTestFilter() Filter {
	return &TestFilter{}
}

//Filter : 词法分析
func (p *TestFilter) Filter(data *LineInfo, line string, fm *FilterManager) (string, error) {
	t := ToTestModel(data.FilterMap)
	switch t.LineNum {
	case 1:
		if t.FuncName != "" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 13:
		if t.FuncName != "" || t.PackageName != "test_data" || t.TypeName != "LineInfo" || t.TypeType != "struct" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 20:
		if t.FuncName != "Clear" || t.PackageName != "test_data" || t.TypeName != "LineInfo" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 27:
		if t.FuncName != "ClearFilterMap" || t.PackageName != "test_data" || t.TypeName != "LineInfo" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 34:
		if t.FuncName != "" || t.PackageName != "test_data" || t.TypeName != "Compiler" || t.TypeType != "struct" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 39:
		if t.FuncName != "NewCompiler" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 46:
		if t.FuncName != "AddFilter" || t.PackageName != "test_data" || t.TypeName != "Compiler" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 56:
		if t.FuncName != "Entrance" || t.PackageName != "test_data" || t.TypeName != "Compiler" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 73:
		if t.FuncName != "walkFile" || t.PackageName != "test_data" || t.TypeName != "Compiler" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 96:
		if t.FuncName != "handleFile" || t.PackageName != "test_data" || t.TypeName != "Compiler" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 104:
		if t.FuncName != "isDir" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 111:
		if t.FuncName != "exists" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 119:
		if t.FuncName != "walkDir" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 138:
		if t.FuncName != "listDir" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	case 158:
		if t.FuncName != "copyMap" || t.PackageName != "test_data" || t.TypeName != "" || t.TypeType != "" {
			fmt.Errorf("line: %d ! lexical analysis error!")
		}
	default:
	}

	return fm.Next(data, line, fm)
}

//FileChange : 文件改变时数据清空
func (p *TestFilter) FileChange(old, new string) error {
	return nil
}

type TestModel struct {
	FuncName    string
	LineNum     int
	PackageName string
	TypeName    string
	TypeType    string
}

func ToTestModel(data map[string]interface{}) *TestModel {
	t := new(TestModel)
	t.FuncName = data["FuncName"].(string)
	t.LineNum = data["LineNum"].(int)
	t.PackageName = data["PackageName"].(string)
	t.TypeName = data["TypeName"].(string)
	t.TypeType = data["TypeType"].(string)
	return t
}
