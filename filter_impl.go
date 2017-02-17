package gopc

const (
	//LineNum 存lineNum值的名
	LineNum string = "LineNum"

	//PackageName 存packageName值的名
	PackageName string = "PackageName"

	//TypeName 存TypeName值的名
	TypeName string = "TypeName"

	//TypeType 存TypeType值的名
	TypeType string = "TypeType"

	//FuncName 存FuncName值的名
	FuncName string = "FuncName"
)

//LineNumFilter : 行号记录
type LineNumFilter struct {
	lineNum int
}

//NewLineNumFilter 新建行号过滤器
func NewLineNumFilter() Filter {
	return &LineNumFilter{
		lineNum: 0,
	}
}

//Filter : 行号+1
func (l *LineNumFilter) Filter(data *LineInfo, line string, fm *FilterManager) (string, error) {
	l.lineNum = l.lineNum + 1
	data.FilterMap[LineNum] = l.lineNum
	return fm.Next(data, line, fm)
}

//FileChange : 文件改变时行号清0
func (l *LineNumFilter) FileChange(old, new string) error {
	l.lineNum = 0
	return nil
}

//LexicalFilter : 词法过滤器
type LexicalFilter struct {
	lexicalAnalyzer *LexicalAnalyzer
}

//NewLexicalFilter 新建词法过滤器
func NewLexicalFilter() Filter {
	return &LexicalFilter{
		lexicalAnalyzer: NewLexicalAnalyzer(),
	}
}

//Filter : 词法分析
func (p *LexicalFilter) Filter(data *LineInfo, line string, fm *FilterManager) (string, error) {
	p.lexicalAnalyzer.HandleLine(line)
	data.FilterMap[PackageName] = p.lexicalAnalyzer.PackageName
	data.FilterMap[TypeName] = p.lexicalAnalyzer.StructName
	data.FilterMap[TypeType] = p.lexicalAnalyzer.StructType
	data.FilterMap[FuncName] = p.lexicalAnalyzer.FuncName
	return fm.Next(data, line, fm)
}

//FileChange : 文件改变时数据清空
func (p *LexicalFilter) FileChange(old, new string) error {
	p.lexicalAnalyzer.Reset()
	return nil
}
