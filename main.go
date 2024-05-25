package main

import (
	"flag"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

// URLConfig 结构体用于表示YAML文件中的配置
type URLConfig struct {
	Name       string `yaml:"name"`
	Poc        string `yaml:"poc"`
	StatusCode int    `yaml:"status_code"`
}

// readYAMLFiles 读取目录下的所有.yaml文件并返回URL配置
func readYAMLFiles(dir string) ([]URLConfig, error) {
	var configs []URLConfig

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".yaml" {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}

			var config URLConfig
			err = yaml.Unmarshal(data, &config)
			if err != nil {
				return err
			}

			configs = append(configs, config)
		}

		return nil
	})

	return configs, err
}

// sendHTTPRequest 发送HTTP GET请求并打印响应
func sendHTTPRequest(name, baseURL, poc string, expectedStatusCode int) {
	fullURL, err := url.Parse(baseURL)
	if err != nil {
		log.Printf("无法解析基URL: %v\n", err)
		return
	}

	pocURL, err := url.Parse(poc)
	if err != nil {
		log.Printf("无法解析POC URL: %v\n", err)
		return
	}

	fullURL = fullURL.ResolveReference(pocURL)

	resp, err := http.Get(fullURL.String())
	if err != nil {
		log.Printf("无法发送请求到 %s (%s): %v\n", name, fullURL.String(), err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == expectedStatusCode {
		fmt.Printf("响应来自 %s (%s): 状态码匹配 (%d)\n", name, fullURL.String(), resp.StatusCode)
	} else {
		fmt.Printf("响应来自 %s (%s): 状态码不匹配 (预期: %d, 实际: %d)\n", name, fullURL.String(), expectedStatusCode, resp.StatusCode)
	}
}

// writePocsToFile 将所有POC写入文本文件
func writePocsToFile(configs []URLConfig, filePath string) error {
	var pocs []string
	for _, config := range configs {
		pocs = append(pocs, config.Poc)
	}

	content := strings.Join(pocs, "\n")

	return ioutil.WriteFile(filePath, []byte(content), 0644)
}

func main() {
	// 从命令行解析参数
	baseURL := flag.String("u", "", "指定基URL")
	outputFile := flag.String("t", "", "将POC输出到指定的文本文件")
	flag.Parse()

	// 指定目标目录
	dir := "./fuzz"

	// 读取YAML文件
	configs, err := readYAMLFiles(dir)
	if err != nil {
		log.Fatalf("读取YAML文件失败: %v\n", err)
	}

	// 如果指定了输出文件，则写入POC到文件
	if *outputFile != "" {
		err := writePocsToFile(configs, *outputFile)
		if err != nil {
			log.Fatalf("写入POC到文件失败: %v\n", err)
		}
		fmt.Printf("所有POC已写入文件 %s\n", *outputFile)
		return
	}

	// 如果未指定基URL，则提示错误并退出
	if *baseURL == "" {
		log.Fatal("必须指定基URL (-u 参数)")
	}

	// 对每个URL发送HTTP请求
	for _, config := range configs {
		sendHTTPRequest(config.Name, *baseURL, config.Poc, config.StatusCode)
	}
}
