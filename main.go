package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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

func main() {
	// 指定目标目录
	dir := "./fuzz"
	// 基URL
	baseURL := "http://example.com"

	// 读取YAML文件
	configs, err := readYAMLFiles(dir)
	if err != nil {
		log.Fatalf("读取YAML文件失败: %v\n", err)
	}

	// 对每个URL发送HTTP请求
	for _, config := range configs {
		sendHTTPRequest(config.Name, baseURL, config.Poc, config.StatusCode)
	}
}
