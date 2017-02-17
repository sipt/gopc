package gopc

import (
	"regexp"
	"strings"
)

//LexicalAnalyzer 词法分析器
type LexicalAnalyzer struct {
	PackageName    string
	StructName     string
	StructType     string
	FuncName       string
	stack          *Stack
	KeyWorkVisible bool
	InSpace        bool
	cutRegs        []*regexp.Regexp
	Line           string
}

//NewLexicalAnalyzer 新建词法分析器
func NewLexicalAnalyzer() *LexicalAnalyzer {
	c := &LexicalAnalyzer{
		KeyWorkVisible: true,
		InSpace:        false,
		cutRegs:        make([]*regexp.Regexp, 0),
		stack:          NewStringStack(),
	}
	c.cutRegs = append(c.cutRegs, regexp.MustCompile("\".*\""))
	c.cutRegs = append(c.cutRegs, regexp.MustCompile("//.*"))
	return c
}

//Clear 清空数据
func (c *LexicalAnalyzer) Clear() {
	c.StructName = ""
	c.StructType = ""
	c.FuncName = ""
	c.KeyWorkVisible = true
	c.InSpace = false
}

//Reset 清空数据
func (c *LexicalAnalyzer) Reset() {
	c.PackageName = ""
	c.Clear()
	c.stack.Clear()
}

//HandleLine 处理单行信息
func (c *LexicalAnalyzer) HandleLine(line string) {
	c.Line = line
	if c.KeyWorkVisible {
		line = c.cutLine(line)
		if c.InSpace {
			c.lexicalAnalysis(line)
		} else {
			var reg *regexp.Regexp
			if strings.HasPrefix(line, "package") {
				reg = regexp.MustCompile("package (\\w+)")
				result := reg.FindAllStringSubmatch(line, -1)
				if len(result) >= 1 {
					if len(result[0]) == 2 {
						c.Reset()
						c.PackageName = result[0][1]
					}
				}
			} else if strings.HasPrefix(line, "type") {
				reg = regexp.MustCompile("type ([^ ]+) ([A-Za-z0-9_]\\w+)")
				result := reg.FindAllStringSubmatch(line, -1)
				if len(result) >= 1 {
					if len(result[0]) == 3 {
						c.Clear()
						c.StructName = result[0][1]
						c.StructType = result[0][2]
						c.InSpace = true
						c.stack.Push("{")
					}
				}
			} else if strings.HasPrefix(line, "func") {
				reg = regexp.MustCompile("func (?:[\\(](?:[\\w]+[ ])?[\\*]?([\\w]+)?[\\)][ ]?)?([A-Za-z0-9_]\\w*).*\\{")
				result := reg.FindAllStringSubmatch(line, -1)
				if len(result) >= 1 {
					if len(result[0]) == 3 {
						c.Clear()
						c.StructName = result[0][1]
						c.FuncName = result[0][2]
						c.InSpace = true
						c.stack.Push("{")
					}
				}
			}
		}
	} else {
		c.lexicalAnalysis(line)
	}
}

func (c *LexicalAnalyzer) lexicalAnalysis(line string) {
	lastChar, err := c.stack.Peep()
	if err != nil {
		lastChar = ""
	}
	charse := []rune(line)
	for i := 0; i < len(charse); i++ {
		switch charse[i] {
		case '{':
			if c.KeyWorkVisible {
				c.stack.Push("{")
			}
		case '}':
			if c.InSpace && lastChar == "{" {
				c.stack.Pop()
				lastChar, err = c.stack.Peep()
				if err != nil {
					c.Clear()
				}
			}
		case '`':
			if !c.KeyWorkVisible && lastChar == "`" {
				c.stack.Pop()
				c.KeyWorkVisible = true
			} else if c.KeyWorkVisible {
				c.stack.Push("`")
				c.KeyWorkVisible = false
			}
		case '/':
			if c.KeyWorkVisible && len(charse) > i+1 && charse[i+1] == '*' {
				i++
				c.KeyWorkVisible = false
				c.stack.Push("/*")
			}
		case '*':
			if !c.KeyWorkVisible && len(charse) > i+1 && charse[i+1] == '/' && lastChar == "*/" {
				i++
				c.KeyWorkVisible = true
				c.stack.Pop()
			}
		default:
		}
	}
}

func (c *LexicalAnalyzer) cutLine(line string) string {
	result := line
	for _, reg := range c.cutRegs {
		result = reg.ReplaceAllString(line, "")
	}
	return result
}
