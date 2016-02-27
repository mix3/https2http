package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Conf struct {
	Addr string `envconfig:"ADDR" default:":19300"`
}

var conf Conf

func init() {
	err := envconfig.Process("", &conf)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url := "https://" + strings.TrimPrefix(r.URL.Path, "/")
		resp, err := http.Get(url)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Fprintln(w, err)
			return
		}
		for k, vs := range resp.Header {
			for _, v := range vs {
				w.Header().Add(k, v)
			}
		}
		fmt.Fprint(w, string(body))
	})
	log.Fatal(http.ListenAndServe(conf.Addr, nil))
}
