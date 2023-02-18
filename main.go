package main

import (
	"bufio"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func main() {
	// 设置正则表达式规则
	reg := regexp.MustCompile(`v.douyin.com/(.*?)/`)
	// 临时地址
	temp := ""
	var videoid string
	for {
		// 获取视频分享链接
		fmt.Println("视频分享url链接:")
		// 使用标准输出获取用户输入,否则如果有空格会多次输出下面那句话
		reader := bufio.NewReader(os.Stdin)
		// 以换行结尾
		temp, _ = reader.ReadString('\n')
		// 字符串中是否包含这个域名
		if !strings.Contains(temp, "v.douyin.com") {
			fmt.Println("请输入v.douyin.com的分享链接")
			continue
		}
		// 正则获取视频短链接id
		videoid = reg.FindAllStringSubmatch(temp, -1)[0][1]
		// 链接正确,下一步
		break
	}

	// 获取真实url
	realUrl := getRealUrl(videoid)

	// 获取视频地址
	videoUrl := getVideo(realUrl)

	// 下载视频
	downloadVideo(videoUrl)
}

// 通过分享的url链接获取视频真实链接
func getRealUrl(url string) string {
	get, _ := http.Get("https://v.douyin.com/" + url + "/")
	reg := regexp.MustCompile(`/video/(.*?)/`)
	// 获取抖音视频id
	videoId := reg.FindAllStringSubmatch(get.Request.URL.String(), 1)[0][1]

	return "https://www.douyin.com/video/" + videoId
}

// 获取视频下载地址
func getVideo(realUrl string) string {
	// 发送请求内容
	url1, _ := url.Parse(realUrl)

	req := &http.Request{
		Method: "GET",
		URL:    url1,
		Header: map[string][]string{},
	}

	// 伪装正常用户
	req.Header.Set("Host", "www.douyin.com")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:109.0) Gecko/20100101 Firefox/110.0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("TE", "trailers")
	req.Header.Set("Sec-Fetch-Site", "cross-site")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Cookie", "")

	// 发送请求
	client := http.Client{}
	get, _ := client.Do(req)
	defer get.Body.Close()

	// 用于保存页面内容
	str := ""
	buf := make([]byte, 1024)
	for {
		n, _ := get.Body.Read(buf)
		if n == 0 {
			break
		}
		// 整个页面是url编码后的数据,然后用js变成页面,获取到编码后将其解码存入字符串
		query, _ := url.QueryUnescape(string(buf[:n]))
		str += query
	}

	// 获取视频地址
	reg1 := regexp.MustCompile(`"playApi":"(.*?)"`)
	return reg1.FindAllStringSubmatch(str, 1)[0][1]
}

// 下载视频
func downloadVideo(url string) {
	// 创建文件
	file, _ := os.Create(strconv.FormatInt(time.Now().UnixMilli(), 10) + ".mp4")

	// 读取链接
	get, _ := http.Get("https:" + url)
	defer get.Body.Close()
	// 将读取的内容写入创建的文件
	buf := make([]byte, 1024)
	for {
		n, _ := get.Body.Read(buf)
		if n == 0 {
			break
		}
		file.Write(buf[:n])
	}
}
