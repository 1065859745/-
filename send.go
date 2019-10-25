package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

var p = flag.String("p", "80", "Transfer server prot,Default port 80")

func main() {
	var wg sync.WaitGroup
	flag.Parse()
	u := flag.Args()
	if len(u) < 1 {
		fmt.Println("DingTalk webhook url has at least one, for example: /path/send.go http://dingtalk.com http://56a4dsasgr4qw78654 ...")
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.Method == "POST" {
				defer r.Body.Close()
				result, _ := ioutil.ReadAll(r.Body)

				fmt.Fprint(w, string(result))
				fmt.Println(string(result))
			} else {
				fmt.Fprint(w, "not post")
			}
		})
		http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			http.Post("http://localhost", "application/json;charset=utf-8", bytes.NewBufferString("{\"msgtype\":\"text\",\"text\":{\"content\":\"Test DingTalk webhook url\"}}"))
			fmt.Fprint(w, "ok")
		})
		wg.Add(1)
		go http.ListenAndServe(":"+*p, nil)
		fmt.Printf("Server starting at http://localhost:%s\n", *p)
		if res, err := http.Get("http://localhost/test"); err != nil {
			log.Fatal(err)
		} else {
			fmt.Println(res.StatusCode)
		}
	}
	wg.Wait()
}
