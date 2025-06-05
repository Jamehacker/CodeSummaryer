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
var textProcessor *TextProcessor

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
	textProcessor = &TextProcessor{}
	textProcessor.Init(generalConfig.RegFormat, generalConfig.RegFormatEnd)
	//遍历文件，查找标记语法
	file, err := os.OpenFile(filename, os.O_RDONLY, 777)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file) //SUM 把文件放入扫描器，可以一行一行地读取

	for scanner.Scan() {
		txt := scanner.Text()
		textProcessor.ReadIn(txt)
	}
	str := textProcessor.GetResult()
	fmt.Println(str)
	WriteInFile(str, dstFile)
	file.Close()

}
func WriteInFile(data string, dstFile string) {
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
		if strings.HasSuffix(path, ".exe") {
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
