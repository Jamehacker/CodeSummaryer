package CodeSummary

import "regexp"

type TextLayer struct {
	matchs     []string    //该层的信息
	layersInfo SummaryInfo //该层SUM-END结构下的SUMMARY的信息。和TextResult中的字段有关
}
type RegFormat struct {
	MatchBegin *regexp.Regexp
	MatchEnd   *regexp.Regexp
}
type SummaryInfo struct {
	startFlag      string //开始符号
	endFlag        string //结束符号
	comment        string //代码段解释
	flagAndComment string //开始符号加代码段解释
}
type TextResult struct {
	info    SummaryInfo
	content []string //代码内容
	id      int      //序号
}

func (t *TextLayer) AppendMatchs(string2 string) {
	t.matchs = append(t.matchs, string2)
}
func (t *TextLayer) GetMatchs() []string {
	return t.matchs
}
func (t *TextLayer) SetStartInfo(startFlag string, comment string) {
	t.layersInfo.startFlag = startFlag
	t.layersInfo.comment = comment
}

func (t *TextLayer) SetEndInfo(endFlag string) {
	t.layersInfo.endFlag = endFlag
}
func (t *TextLayer) GetFirstMatch() string {
	return t.matchs[0]
}
