package CodeSummary //SUM testcode

import (
	"bufio"
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
	StartFlag    string
	EndFlag      string
	RegFormat    string
	RegFormatEnd string
}

var generalConfig Config
var exp *regexp.Regexp
var expEnd *regexp.Regexp

func (c *Config) GetFilledExpStr() string {
	return fmt.Sprintf(c.RegFormat, c.StartFlag)
}
func (c *Config) GetFilledExpEndStr() string {
	return fmt.Sprintf(c.RegFormatEnd, c.EndFlag)
}
func SetDefaultConfig() {
	generalConfig.StartFlag = "SUM"
	generalConfig.EndFlag = "SUMEND"
	generalConfig.RegFormat = "(//\\s*%s)(.+)"
	generalConfig.RegFormatEnd = "(//\\s*%s)"
	exp = regexp.MustCompile(generalConfig.GetFilledExpStr())
	expEnd = regexp.MustCompile(generalConfig.GetFilledExpEndStr())
}

// SUM
func ExtractFrommFile(filename string, dstFile string) {
	//遍历文件，查找标记语法
	file, err := os.OpenFile(filename, os.O_RDONLY, 777)
	if err != nil {
		fmt.Println(err)
	}
	scanner := bufio.NewScanner(file) //SUM 把文件放入扫描器，可以一行一行地读取
	var matchedTexts []string
	var isMatchedResult bool

	for scanner.Scan() {
		text := scanner.Text()
		//找到开始符号
		result := exp.FindStringSubmatch(text) //SUM 找到不同的组
		if len(result) != 0 {
			//如果存在
			isMatchedResult = true
			if len(matchedTexts) != 0 {
				WriteInFile(matchedTexts[0], dstFile)
			}
			matchedTexts = []string{} //SUM 清空切片 重新分配一个切片
		}
		if isMatchedResult {
			matchedTexts = append(matchedTexts, text)
		}
		//找到结束符号
		resultEnd := expEnd.FindStringSubmatch(text)
		if resultEnd != nil {
			//写入开始符号到结束符号之间的全部信息
			WriteInFile(strings.Join(matchedTexts, "\n"), dstFile)
		}

	}

	if len(matchedTexts) != 0 {
		WriteInFile(matchedTexts[0], dstFile)
	}
	file.Close()

}
func WriteInFile(data string, dstFile string) {
	//os.Mkdir("Tests", 666)
	data = strings.Replace(data, generalConfig.StartFlag, "", -1)
	data = strings.Replace(data, "\t", "", -1)

	file, err := os.OpenFile(dstFile, os.O_CREATE|os.O_APPEND, 666) //创建文件并且可读可写
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
		files = append(files, path) //SUM 列表添加元素
		return nil                  //SUMEND
	}) //SUMEND

}
