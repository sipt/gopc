package gopc

//Filter 过滤器
type Filter interface {
	Filter(*LineInfo, string, *FilterManager) (string, error)
	FileChange(old string, new string) error
}

//FilterManager 过滤器注册以及顺序执行
type FilterManager struct {
	filters []Filter
	index   int
}

//NewFilterManager 新建一个FilterManager
func NewFilterManager() *FilterManager {
	return &FilterManager{
		filters: make([]Filter, 0),
		index:   0,
	}
}

//AddFilter 添回过滤器
func (f *FilterManager) AddFilter(filter Filter) {
	f.filters = append(f.filters, filter)
}

//Start 从头开始执行过滤器
func (f *FilterManager) Start(l *LineInfo, line string, fm *FilterManager) (string, error) {
	f.index = -1
	return f.Next(l, line, fm)
}

//Next 执行下一个过滤器，如果没有了就直接返回
func (f *FilterManager) Next(l *LineInfo, line string, fm *FilterManager) (string, error) {
	if f.index+1 >= len(f.filters) {
		return line, nil
	}
	f.index++
	return f.filters[f.index].Filter(l, line, fm)
}
