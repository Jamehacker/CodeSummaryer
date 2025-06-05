package CodeSummary

import (
	"regexp"
	"strconv"
)

type Stack struct {
	data   []string
	number []int
	syms   []Symbol
}

func (s *Stack) Push(str string, num int, sym Symbol) {
	s.data = append(s.data, str)
	s.number = append(s.number, num)
	s.syms = append(s.syms, sym)
}
func (s *Stack) Pop() (string, int, Symbol) {
	length := len(s.data)
	if length == 0 {
		return "", -1, Symbol{}
	}
	str := s.data[length-1]
	s.data = s.data[:length-1]
	num := s.number[length-1]
	s.number = s.number[:length-1]

	sym := s.syms[length-1]
	s.syms = s.syms[:length-1]

	return str, num, sym
}

type MatchResultText struct {
	text []string
	sym  Symbol
}

func (m *MatchResultText) Append(str string) {
	m.text = append(m.text, str)
}
func (m *MatchResultText) AppendSlice(str []string) {
	m.text = append(m.text, str...)
}

type Symbol struct {
	symbolText     string
	commentFlag    string
	commentContent string
}

func (s *Symbol) IsValid() bool {
	return s.symbolText != ""
}

// 匹配嵌套的标注
type Matcher struct {
	matchSymbols []Symbol
	matchText    []string
	stack        Stack
	startMatch   *regexp.Regexp
	endMatch     *regexp.Regexp
	resultText   []MatchResultText
}

// ReadIn
//
//	@Description: 读入一条数据
//	@receiver m
//	@param symbol  含有标记和数字的字符串
//	@param lineText  文件中的一行文本
//	@param commentFlag
//	@param commentContent
func (m *Matcher) ReadIn(symbol string, lineText string, commentFlag string, commentContent string) {
	sym := Symbol{
		symbolText:     symbol,
		commentFlag:    commentFlag,
		commentContent: commentContent,
	}
	m.matchSymbols = append(m.matchSymbols, sym)
	m.matchText = append(m.matchText, lineText)
}

func (m *Matcher) SetRegMatch(startFlag string, endFlag string) {
	m.startMatch = regexp.MustCompile(startFlag + "\\d*")
	m.endMatch = regexp.MustCompile(endFlag + "\\d*")
}

var regIntMatch = regexp.MustCompile("\\d+")

// Reformat 把嵌套语句提取出来
func (m *Matcher) Reformat() {
	for i, val := range m.matchSymbols {
		if !val.IsValid() {
			continue
		}
		resultStart := m.startMatch.FindString(val.symbolText)
		resultEnd := m.endMatch.FindString(val.symbolText)
		if resultStart != "" && resultEnd == "" {
			//开始符号
			m.stack.Push(resultStart, i, val)
		} else if resultStart == "" && resultEnd != "" {
			m.Store(resultEnd, i)
		}
	}
	for {
		_, lineNum, sym := m.stack.Pop()
		if lineNum == -1 {
			return
		}
		m.StoreLine(lineNum, lineNum, sym)
	}
}

// StoreLine
//
//	@Description: 遍历栈，找到结束符号对应的起始符号，并保存遍历过的内容
//	@receiver m
//	@param endSymbol
func (m *Matcher) Store(endSymbol string, endLine int) {
	//遇到结束符号，查找对应的起始符号
	symbolEndNo := GetNumber(endSymbol)
	symbol, lineNum, sym := m.stack.Pop()
	if lineNum == -1 {
		return
	}

	symbolCurNo := GetNumber(symbol)
	for {
		if symbolCurNo == symbolEndNo {
			//找到开始符号
			m.StoreLine(lineNum, endLine, sym)
			break
		} else {
			m.StoreLine(lineNum, lineNum, sym)
			symbol, lineNum, sym = m.stack.Pop()
			if lineNum == -1 {
				return
			}
		}
	}

}

func GetNumber(str string) int {
	res := regIntMatch.FindString(str)
	num, err := strconv.Atoi(res)
	if err != nil {
		return -1
	}
	return num
}
func (m *Matcher) StoreLine(from int, end int, symbol Symbol) {
	if from == end {
		rt := MatchResultText{sym: symbol}
		rt.Append(m.matchText[from])
		m.resultText = append(m.resultText, rt)
	} else {
		rt := MatchResultText{sym: symbol}
		rt.AppendSlice(m.matchText[from : end+1])
		m.resultText = append(m.resultText, rt)
	}
}
func (m *Matcher) WriteToParser(parser *Parser) {
	for i := 0; i < len(m.resultText); i++ {
		rt := m.resultText[i]
		parser.AppendText(rt.sym.symbolText, true, rt.sym.commentFlag, rt.sym.commentContent)
		for j := 1; j < len(rt.text); j++ {
			parser.AppendText(rt.text[j], false, "", "")
		}
	}
}
