package main

import (
	"MyTools/CodeSummary"
	"flag"
	"fmt"
	"os"
)

// SUM 判断文件是否存在
func IsExist(filePath string) bool {
	if filePath == "" {
		return false
	}
	_, err := os.Stat(filePath)
	if err == nil {
		return true
	}
	return false
}

//SUMEND

func main() {
	//SUM 处理命令行参数
	//createJsonPath := flag.String("create", "", "创建该工具的配置模板")
	configPath := flag.String("cfg", "", "使用指定json文件中的配置")
	rootDir := flag.String("dir", ".", "指定要递归查找的根目录")
	output := flag.String("o", "", "指定结果保存文件")
	//
	//flag.Parse()
	//
	////设置配置
	if IsExist(*configPath) {
		CodeSummary.SetConfig(*configPath)
	} else {
		CodeSummary.SetDefaultConfig()
	}
	////如果输出目录不为空，默认开始
	*output = "1.txt"
	*rootDir = "Tests"
	if *output != "" {
		fmt.Println("使用默认的配置")
		fmt.Println("提取开始")
		os.Truncate(*output, 0) //SUM 把文件的大小设置为0
		CodeSummary.ExtractFromFolder(*rootDir, *output)
		fmt.Println("提取结束")
		return
	}
}
