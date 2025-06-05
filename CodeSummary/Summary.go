package CodeSummary

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// 一些配置
// Flag 标记词,指示标记器的开始位置
type Config struct {
	StartFlag      string `json:"start_flag"`
	EndFlag        string `json:"end_flag"`
	RegFormat      string `json:"reg_format"`
	RegFormatEnd   string `json:"reg_format_end"`
	IsMarkdownType bool   `json:"markdown_type"`
}

var generalConfig Config
var exp *regexp.Regexp
var expEnd *regexp.Regexp
var parser Parser
var matcher Matcher

//	func (c *Config) GetFilledExpStr() string {
//		return fmt.Sprintf(c.RegFormat, c.StartFlag)
//	}
//
//	func (c *Config) GetFilledExpEndStr() string {
//		return fmt.Sprintf(c.RegFormatEnd, c.EndFlag)
//	}
func (c *Config) IsValid() bool {
	if c.RegFormat == "" || c.RegFormatEnd == "" {
		return false
	}
	return true
}
func GetDefaultConfigObj() *Config {
	cfg := Config{}
	cfg.StartFlag = "SUM"
	cfg.EndFlag = "SUMEND"
	cfg.RegFormat = "(//\\s*SUM\\d*\\s+)(.+)"
	cfg.RegFormatEnd = "(//\\s*SUMEND\\d*)"
	cfg.IsMarkdownType = true
	return &cfg
}

// 从json文件中设置工具的配置选项
func SetConfig(jsonFilePath string) {
	bytes, err := os.ReadFile(jsonFilePath) //SUM 读取文件全部内容
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(bytes, &generalConfig)
	if err != nil {
		fmt.Println(err)
		return
	}
	if !generalConfig.IsValid() {
		panic("json模板无效")
	}
	exp = regexp.MustCompile(generalConfig.RegFormat)
	expEnd = regexp.MustCompile(generalConfig.RegFormatEnd)
}

// 创建一个工具配置的json模板
func CreateJsonFormat(dstpath string) {
	exp_cfg := GetDefaultConfigObj()
	b, _ := json.Marshal(exp_cfg)
	file, err := os.OpenFile(dstpath, os.O_CREATE|os.O_WRONLY, 666)
	if err != nil {
		fmt.Println(err)
		return
	}
	file.Write(b)
	file.Close()
}
func SetDefaultConfig() {
	generalConfig = *GetDefaultConfigObj()
	exp = regexp.MustCompile(generalConfig.RegFormat)
	expEnd = regexp.MustCompile(generalConfig.RegFormatEnd)
}

// SUM
func ExtractFrommFile(filename string, dstFile string) {
	parser = Parser{
		AllText:    nil,
		IsMarkdown: false,
		Count:      parser.Count,
	}
	parser.EnableMarkdown(generalConfig.IsMarkdownType)

	matcher = Matcher{}
	matcher.SetRegMatch(generalConfig.StartFlag, generalConfig.EndFlag)

	//遍历文件，查找标记语法
	file, err := os.OpenFile(filename, os.O_RDONLY, 777)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file) //SUM 把文件放入扫描器，可以一行一行地读取
	//var matchedTexts []string
	var isMatchedResult bool

	for scanner.Scan() {
		txt := scanner.Text()
		//找到开始符号
		result := exp.FindStringSubmatch(txt) //SUM 正则表达式找到不同的组,result第0位是匹配的全部，后面几位是分组的结果
		if len(result) == 3 {
			//如果存在
			isMatchedResult = true
			parser.AppendText(txt, true, result[1], result[2])
			//matcher.ReadIn(result[1], txt, result[1], result[2])
			continue
		}

		//找到结束符号
		resultEnd := expEnd.FindStringSubmatch(txt)
		if resultEnd != nil {
			//替换结束符号为空白
			txt = strings.Replace(txt, resultEnd[0], "", -1)
			parser.AppendText(txt, false, "", "")
			parser.Store()
			isMatchedResult = false
			continue
		}

		//处于多行文本之间时候，添加文本
		if isMatchedResult {
			parser.AppendText(txt, false, "", "")
		}
		//if isMatchedResult {
		//	matcher.ReadIn("", txt, "", "")
		//}
		//resultEnd := expEnd.FindStringSubmatch(txt)
		//if resultEnd != nil {
		//	//替换结束符号为空白
		//	isMatchedResult = false
		//}
	}
	//matcher.Reformat()
	//matcher.WriteToParser(&parser)
	parser.ConvertToMarkdown()
	WriteInFile(parser.GetText(true), dstFile)
	file.Close()

}
func WriteInFile(data string, dstFile string) {
	//比较大小,先替换较长的字符串,如SUMEND,再替换SUM,防止出现替换了一半的情况
	//if len(generalConfig.EndFlag) > len(generalConfig.StartFlag) {
	//	data = strings.Replace(data, generalConfig.EndFlag, "", -1)
	//	data = strings.Replace(data, generalConfig.StartFlag, "", -1)
	//} else {
	//	data = strings.Replace(data, generalConfig.StartFlag, "", -1)
	//	data = strings.Replace(data, generalConfig.EndFlag, "", -1)
	//}
	if data == "" {
		return
	}
	data = strings.Replace(data, "\t", "", -1)

	file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_APPEND, 666) //SUM 创建文件并且可读可写
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewWriter(file)

	_, err = scanner.WriteString(data) //SUM 每次写入都会清除缓存区
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("WriteString: ", data)
		scanner.WriteString("\n")
		scanner.Flush() //SUM 使用Flush把缓存区内容写入文件
	}
	file.Close()
}
func ExtractFromFolder(folderPath string, dstFile string) {
	var files []string
	filepath.Walk(folderPath, func(path string, info fs.FileInfo, err error) error { //SUM 遍历文件夹
		//判断路径是否是隐藏路径
		if strings.HasPrefix(path, ".") { //SUM 判断字符串是否以...开头
			return nil
		}
		if !info.IsDir() && path != dstFile {
			files = append(files, path)
		}
		return nil

	}) //SUMEND
	for _, file := range files {
		ExtractFrommFile(file, dstFile)
	}

}
