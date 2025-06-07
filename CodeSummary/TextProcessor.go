package CodeSummary

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type TextProcessor struct {
	reg          *RegFormat
	layers       map[int]*TextLayer //保存当前正在匹配的层，匹配好后对应键值移动到layers_Taken
	layers_Taken map[int]*TextLayer //匹配好SUM-END对儿的层
	layerCount   int                //记录层的个数
	//result       [][]string
}

func (r *RegFormat) init(begin string, end string) {
	r.MatchBegin = regexp.MustCompile(begin)
	r.MatchEnd = regexp.MustCompile(end)
}
func (t *TextProcessor) Init(begin string, end string) {
	t.reg = &RegFormat{MatchBegin: nil, MatchEnd: nil}
	t.reg.init(begin, end)
	t.layers = make(map[int]*TextLayer)
	t.layers_Taken = make(map[int]*TextLayer)
	t.layerCount = 0
}

// 读入一行数据，保存匹配内容
func (t *TextProcessor) ReadIn(lineText string) {
	resultBegin := t.reg.MatchBegin.FindStringSubmatch(lineText)
	resultEnd := t.reg.MatchEnd.FindStringSubmatch(lineText)
	if resultBegin == nil && resultEnd == nil {
		//不是开始和结束符号的行
		t.insertTextToLayers(lineText)
	} else if resultBegin != nil {
		t.insertNewLayerToMap(lineText, resultBegin)
	} else {
		t.takeCurrentLayerToResult(lineText, resultEnd)
	}

}
func (t *TextProcessor) insertNewLayerToMap(string2 string, matchString []string) {
	if len(matchString) != 3 {
		fmt.Println("标记开始（reg_format）的正则表达式格式错误，无法匹配到两个需要的内容！")
		return
	}
	t.insertTextToLayers(string2)
	layersInfo := SummaryInfo{
		startFlag:      matchString[1],
		endFlag:        "",
		comment:        matchString[2],
		flagAndComment: matchString[0],
	}
	t.layers[t.layerCount] = &TextLayer{
		matchs:     []string{string2},
		layersInfo: layersInfo,
	}
	t.layerCount++

}

// insertTextToLayers
//
//	@Description: 给当前层添加一个string，如果没有层就不添加
//	@receiver t
//	@param string2
func (t *TextProcessor) insertTextToLayers(string2 string) {
	for key, _ := range t.layers {
		t.layers[key].AppendMatchs(string2)
	}

}
func (t *TextProcessor) takeCurrentLayerToResult(string2 string, matchString []string) {
	if len(t.layers) == 0 {
		fmt.Println("====Error 无法匹配的终止符号：'" + matchString[1] + "'=========")
		return
	}

	t.insertTextToLayers(string2)
	//t.result = append(t.result, t.layers[len1-1].GetMatchs()) //把最后一层加到结构体处理结果中
	if len(matchString) != 2 {
		fmt.Println("标记开始（reg_format）的正则表达式格式错误，无法匹配到两个需要的内容！")
		return
	}
	var keys []int
	for key, _ := range t.layers {
		keys = append(keys, key)
	}
	sort.Ints(keys) //SUM 升序排序
	lastKey := keys[len(keys)-1]
	t.layers[lastKey].SetEndInfo(matchString[1])
	t.layers_Taken[lastKey] = t.layers[lastKey]
	delete(t.layers, lastKey)
}

// GetResult
//
//	@Description: 处理并返回结果
//	@receiver t
//	@return string
func (t *TextProcessor) GetResult() string {
	//没有匹配到end符号的Layers，对应层的第一个，其他舍去
	for key, val := range t.layers {
		//t.result = append(t.result, []string{val.GetFirstMatch()})
		t.layers_Taken[key] = &TextLayer{
			layersInfo: val.layersInfo,
			matchs:     []string{val.GetFirstMatch()},
		}
	}
	var returnString string
	// 遍历处理结果
	i := 0
	for _, val := range t.layers_Taken {
		te := TextResult{
			info:    val.layersInfo,
			content: val.GetMatchs(),
			id:      i + 1,
		}
		returnString += te.ConvertToMarkdown() + "\n"
		i++
	}
	return returnString
}

func (te *TextResult) ConvertToMarkdown() string {
	rawText := strings.Join(te.content, "\n")
	rawText = strings.Replace(rawText, te.info.flagAndComment, "", -1)
	rawText = strings.Replace(rawText, te.info.endFlag, "", -1)

	// 根据注释内容创建加粗标题
	styleFlag := "" //文本样式标记符号 ** == * %%
	var comment string
	if te.id > -1 {
		comment = strconv.Itoa(te.id) + ". " + styleFlag + te.info.comment + styleFlag + ": "
	} else {
		comment = styleFlag + te.info.comment + styleFlag + ": "
	}

	var newString string
	//判断是否是多行文本
	if len(te.content) > 1 {
		newString += comment + "\n"
		newString += "```js" + "\n"
		newString += rawText + "\n"
		newString += "```"
	} else {
		newString += comment + " " + "`" + rawText + "`"
	}
	return newString
}
