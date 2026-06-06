package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// CurrentVersion 当前版本号
// 在 init() 中从 wails.json 的 version 字段读取
var CurrentVersion = "0.0.0"

// GitHub 仓库地址
const RepoOwner = "smile-yan"
const RepoName = "easy-flip-clock"

// wailsConfig 用于解析 wails.json 中的版本字段
type wailsConfig struct {
	Version string `json:"version"`
}

// init 在程序启动时从 wails.json 读取 version 作为 CurrentVersion 的初始值
func init() {
	data, err := os.ReadFile("wails.json")
	if err != nil {
		log.Printf("读取 wails.json 失败，使用默认版本 %s: %v", CurrentVersion, err)
		return
	}

	var cfg wailsConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		log.Printf("解析 wails.json 失败，使用默认版本 %s: %v", CurrentVersion, err)
		return
	}

	if cfg.Version != "" {
		CurrentVersion = cfg.Version
		log.Printf("CurrentVersion 初始化为: %s", CurrentVersion)
	}
}

// ReleaseInfo 存储 GitHub release 信息
type ReleaseInfo struct {
	TagName     string `json:"tag_name"`
	Name        string `json:"name"`
	Body        string `json:"body"`
	HTMLURL     string `json:"html_url"`
	DownloadURL string `json:"download_url"`
	Prerelease  bool   `json:"prerelease"`
}

// UpdateResult 检查更新结果
type UpdateResult struct {
	HasUpdate   bool   `json:"hasUpdate"`
	CurrentVer  string `json:"currentVer"`
	LatestVer   string `json:"latestVer"`
	ReleaseNote string `json:"releaseNote"`
	DownloadURL string `json:"downloadUrl"`
	Message     string `json:"message"`
}

// parseVersion 将版本号字符串解析为可比较的整数数组
// 例如 "v1.2.3" -> [1, 2, 3]
func parseVersion(version string) ([]int, error) {
	// 移除前缀 "v"
	version = strings.TrimPrefix(version, "v")
	version = strings.TrimPrefix(version, "V")

	// 使用正则表达式提取数字部分
	re := regexp.MustCompile(`(\d+)`)
	matches := re.FindAllStringSubmatch(version, -1)

	var nums []int
	for _, match := range matches {
		num, err := strconv.Atoi(match[1])
		if err != nil {
			return nil, err
		}
		nums = append(nums, num)
	}

	return nums, nil
}

// compareVersions 比较两个版本号
// 返回值：0 表示相等，1 表示 v1 > v2，-1 表示 v1 < v2
func compareVersions(v1, v2 string) (int, error) {
	nums1, err := parseVersion(v1)
	if err != nil {
		return 0, err
	}
	nums2, err := parseVersion(v2)
	if err != nil {
		return 0, err
	}

	// 补齐长度
	maxLen := len(nums1)
	if len(nums2) > maxLen {
		maxLen = len(nums2)
	}
	for len(nums1) < maxLen {
		nums1 = append(nums1, 0)
	}
	for len(nums2) < maxLen {
		nums2 = append(nums2, 0)
	}

	// 逐个比较
	for i := 0; i < maxLen; i++ {
		if nums1[i] > nums2[i] {
			return 1, nil
		}
		if nums1[i] < nums2[i] {
			return -1, nil
		}
	}

	return 0, nil
}

// CheckForUpdate 检查应用更新
func CheckForUpdate() *UpdateResult {
	result := &UpdateResult{
		CurrentVer: CurrentVersion,
		Message:    "检查更新中...",
	}

	// 调用 GitHub API 获取最新 release
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", RepoOwner, RepoName)

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		result.Message = fmt.Sprintf("创建请求失败: %v", err)
		return result
	}

	// 设置 User-Agent（GitHub API 要求）
	req.Header.Set("User-Agent", "easy-flip-clock")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		result.Message = fmt.Sprintf("网络请求失败: %v", err)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.Message = fmt.Sprintf("API 请求失败，状态码: %d", resp.StatusCode)
		return result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.Message = fmt.Sprintf("读取响应失败: %v", err)
		return result
	}

	var release ReleaseInfo
	if err := json.Unmarshal(body, &release); err != nil {
		result.Message = fmt.Sprintf("解析响应失败: %v", err)
		return result
	}

	result.LatestVer = release.TagName
	result.ReleaseNote = release.Body
	result.DownloadURL = release.HTMLURL

	// 比较版本
	cmp, err := compareVersions(CurrentVersion, release.TagName)
	if err != nil {
		result.Message = fmt.Sprintf("版本比较失败: %v", err)
		return result
	}

	if cmp < 0 {
		result.HasUpdate = true
		result.Message = fmt.Sprintf("发现新版本 %s", release.TagName)
	} else if cmp == 0 {
		result.HasUpdate = false
		result.Message = "当前已是最新版本"
	} else {
		result.HasUpdate = false
		result.Message = fmt.Sprintf("当前版本比发布版本还新 (%s)", CurrentVersion)
	}

	return result
}
