package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Config struct {
	Version string `json:"version"`
}

func Version() {
	jsonData, err := os.ReadFile("config/config.json")
	if err != nil {
		log.Fatalf("无法读取配置文件: %v", err)
	}

	// 创建 Config 结构体实例
	var config Config

	// 反序列化 JSON 数据到 Config 结构体
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		log.Fatalf("解析 JSON 数据时出错: %v", err)
	}

	// 打印 version 信息
	log.Println(config.Version)

	// 提升 version 版本
	err = config.VersionIncrement()
	if err != nil {
		log.Printf("提升 version 时出错: %s", err)
	}

	// 将新的 JSON 写回文件
	updatedJsonData, _ := json.MarshalIndent(config, "", "  ")
	err = os.WriteFile("config/config.json", updatedJsonData, 0644)
	if err != nil {
		log.Fatalf("写入文件出错: %v", err)
	}
}

func (config *Config) VersionIncrement() error {
	// 使用正则表达式提取版本号中的数字部分
	re := regexp.MustCompile(`\d+(\.\d+)*$`)
	versionPart := re.FindString(config.Version)

	if versionPart == "" {
		return fmt.Errorf("版本号格式无效: %s", config.Version)
	}

	// 将版本号按 "." 分割
	parts := strings.Split(versionPart, ".")

	// 将最后一部分转为整数并自增
	lastIndex := len(parts) - 1
	lastNum, _ := strconv.Atoi(parts[lastIndex])
	parts[lastIndex] = strconv.Itoa(lastNum + 1)

	// 重新拼接版本号
	newVersion := strings.Join(parts, ".")

	// 将原始的非数字部分加上新的版本号
	config.Version = re.ReplaceAllString(config.Version, newVersion)

	return nil
}
