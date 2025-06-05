package CodeSummary

import (
	"strconv"
	"strings"
)

// 用于解析文本
type text struct {
	multiLinetext  []string //含有注释的多行文本
	temptext       []string //缓存未确认的多行文本（如果匹配到停止符号就存入multiline中，否则后续转换时候只取第一条记录）
	commentFlag    string   //注释标记符号
	commentContent string   //注释内容
	result         string   //处理结果的文本
	id             int      //该文本处于第几条
}

type Parser struct {
	AllText    []text
	IsMarkdown bool
	Count      int //统计有多少条文本内容，便于编号
}

func (p *Parser) EnableMarkdown(isenable bool) {
	p.IsMarkdown = isenable
}

// AppendText
//
//		@Description: 添加一条文本
//		@receiver p
//		@param oneline  包含注释的全部文本
//		@param isFirstLine  是否是第一行
//		@param commentFlag  注释标记符号
//	 	@param commentContent 具体的注释内容
func (p *Parser) AppendText(oneline string, isFirstLine bool, commentFlag string, commentContent string) {
	//判断是新建还是在已有的多行文本上增加
	if isFirstLine {
		p.Count++
		t := text{commentFlag: commentFlag, commentContent: commentContent, id: p.Count}
		t.Store(oneline)
		t.PushTotemp(oneline)
		p.AllText = append(p.AllText, t)

	} else {
		length := len(p.AllText)
		if length == 0 {
			return
		}
		p.AllText[length-1].PushTotemp(oneline)
	}
}

// Store 如果匹配到停止符号，就调用该函数，把缓存区的内容都保存
func (p *Parser) Store() {
	length := len(p.AllText)
	if length == 0 {
		return
	}
	p.AllText[length-1].StoreAll()
}

func (p *Parser) ConvertToMarkdown() {
	if p.IsMarkdown {
		for i := 0; i < len(p.AllText); i++ {
			p.AllText[i].ConvertToMarkdown()
		}
	} else {
		for i := 0; i < len(p.AllText); i++ {
			p.AllText[i].ConvertToNormalText()
		}
	}
}

func (t *text) ConvertToMarkdown() {
	rawText := strings.Join(t.multiLinetext, "\n")
	rawText = strings.Replace(rawText, t.commentFlag+t.commentContent, "", -1)

	// 根据注释内容创建加粗标题
	styleFlag := "" //文本样式标记符号 ** == * %%
	comment := strconv.Itoa(t.id) + ". " + styleFlag + t.commentContent + styleFlag + ": "

	var newString string
	//判断是否是多行文本
	if len(t.multiLinetext) > 1 {
		newString += comment + "\n"
		newString += "```js" + "\n"
		newString += rawText + "\n"
		newString += "```"
	} else {
		newString += comment + " " + "`" + rawText + "`"
	}
	t.result = newString
}
func (t *text) ConvertToNormalText() {

}

// GetText
//
//	@Description:
//	@receiver p
//	@param IsRemoveEnterKeys 是否移除连续的多个换行符
//	@return string
func (p *Parser) GetText(IsRemoveEnterKeys bool) string {
	var resultList []string
	for _, val := range p.AllText {
		resultList = append(resultList, val.result)
	}
	res := strings.Join(resultList, "\n")
	if IsRemoveEnterKeys {
		res = strings.Replace(res, "\n\n", "", -1)
	}
	return res
}
func (p *Parser) Clear() {
	p.AllText = []text{} //SUM 清空切片 重新分配一个切片
}
func (p *Parser) IsEmpty() bool {
	return len(p.AllText) == 0
}
func (t *text) Store(str string) {
	t.multiLinetext = append(t.multiLinetext, str)
}
func (t *text) Flush() {
	t.temptext = []string{}
}

// 把临时数组中的内容存到multilinetext中，同时清空temp中的内容
func (t *text) StoreAll() {
	t.multiLinetext = make([]string, len(t.temptext))
	for i, _ := range t.temptext {
		t.multiLinetext[i] = t.temptext[i]
	}
	t.temptext = []string{}
}
func (t *text) PushTotemp(str string) {
	t.temptext = append(t.temptext, str)
}
