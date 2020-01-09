package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

var port, checkIp, checkIpMsg = flag.String("p", "5001", "Transfer server prot"), flag.String(`i`, ``, "Check connect status by remote ip"), flag.String(`m`, `监控异常，请检查连接`, "Send mssages when connect faild")
var wg sync.WaitGroup

// 发送到钉钉的数据结构
type Content struct {
	Content string `json:"content"`
}
type At struct {
	AtMobiles []string `json:"atMobiles"`
}
type Ding struct {
	Msgtype string  `json:"msgtype"`
	Text    Content `json:"text"`
	At      At      `json:"at"`
	IsAtAll bool    `json:"isAtAll"`
}

// 定义接收 prometheus json数据
type Annotation struct {
}
type Label struct {
	Alertname, Class, Instance, Job, Project, Severity string
}
type Alert struct {
	Status                                     string
	Labels                                     Label
	Annotations                                Annotation
	StartsAt, EndsAt, GeneratorURL, Ingerprint string
}

func main() {
	flag.Parse()
	if len(flag.Args()) == 0 {
		log.Fatal(`No DingTalk url, At last one DingTalk url in the parameter`)
	}

	// 如果使用vpn连接，检查vpn状态
	if *checkIp != `` {
		checkConn(*checkIp)
	}
	// 测试钉钉url
	sendText(`我来啦，我来啦`, []string{}, false)

	// 处理请求
	http.HandleFunc("/", middleWare(receiveData))

	wg.Add(1)
	go listenAndServer()
	wg.Wait()
}

// 检测vpn连接
func checkConn(vpn string) {
	for {
		conn, e := net.Dial("ip4:icmp", vpn)
		checkErr(e)
		conn.SetReadDeadline((time.Now().Add(time.Second * 2)))
		var msg [32]byte
		msg[0] = 8
		check := checkSum(msg[0:9])
		msg[2] = byte(check >> 8)
		msg[3] = byte(check & 0xff)
		_, e = conn.Write(msg[0:9])
		checkErr(e)
		_, e = conn.Read(msg[0:])
		checkErr(e)
		time.Sleep(time.Second)
	}
}

//检查ip数据包的和
func checkSum(msg []byte) uint16 {
	sum := 0
	for i := 0; i < len(msg)-1; i += 2 {
		sum += int(msg[i]) * 256
	}
	// if len%2 == 1 {
	// 	sum += int(msg[len-1]) * 256 // notice here,why *256?
	// }
	sum = sum & 0xffff
	var answer uint16 = uint16(^sum)
	return answer
}

//检查连接错误
func checkErr(e error) {

	if e != nil {
		sendText(*checkIpMsg, regexp.MustCompile(`@\d{11}`).FindAllString(*checkIpMsg, -1), false)
		log.Fatal(e)
	}
}

// 定义中间件
func middleWare(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println(r.RemoteAddr, r.Method, r.URL)
		f(w, r)
	}
}

// 处理prometheus数据
func receiveData(w http.ResponseWriter, r *http.Request) {
	if alert_byte, e := ioutil.ReadAll(r.Body); e != nil {
		log.Fatal(e)
	} else {
		var (
			alerts []Alert
			text   string
		)
		serious := false
		content := regexp.MustCompile(`\[.*\]`).Find(alert_byte)
		if e = json.Unmarshal(content, &alerts); e != nil {
			log.Fatal(e)
		}
		for _, v := range alerts {
			text += fmt.Sprintf("%s:  %s  %s\n", v.Labels.Class, v.Labels.Instance, v.Labels.Alertname)
			if serious != true {
				if v.Labels.Severity == `serious` {
					serious = true
				}
			}
		}
		sendText(text, []string{}, serious)
	}

	// 函数结束时关闭流
	defer r.Body.Close()
}

func sendText(msg string, at []string, atAll bool) {
	// 将发送的信息转成 json 数据
	var dings Ding
	dings.Msgtype = `text`
	dings.Text.Content = msg
	dings.IsAtAll = atAll
	dings.At = At{AtMobiles: at}
	b, _ := json.Marshal(&dings)
	fmt.Println(string(b))
	client := &http.Client{Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}}
	for _, v := range flag.Args() {
		if r, e := client.Post(v, `application/json`, strings.NewReader(string(b))); e != nil {
			log.Printf("%s ERROR: %s", v, e)
		} else {
			content, e := ioutil.ReadAll(r.Body)
			if e != nil {
				log.Fatal(e)
				return
			}
			log.Printf("%s SUCCESS: %s", v, string(content))
			defer r.Body.Close()
		}
	}
}

// 启动 http 监听
func listenAndServer() {
	log.Println("Server will start at 127.0.0.1:" + *port)
	e := http.ListenAndServe(":"+*port, nil)
	if e != nil {
		log.Fatal(e)
		wg.Done()
	}
}
