package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ian-kent/goose"
)

var stream = goose.NewEventStream()
var proxying = ""
var url = ""

func main() {
	http.Handle("/!-stream", http.HandlerFunc(browserAgentStream))
	http.Handle("/!-switch", http.HandlerFunc(browserAgentSwitch))
	http.Handle("/", http.HandlerFunc(browserAgentProxy))

	err := http.ListenAndServe(":3123", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func browserAgentStream(w http.ResponseWriter, req *http.Request) {
	stream.AddReceiver(w)
}

func browserAgentSwitch(w http.ResponseWriter, req *http.Request) {
	proxying = req.URL.Query().Get("proxy")
	url = req.URL.Query().Get("url")
	eval := req.URL.Query().Get("eval")

	if len(proxying) == 0 && len(url) == 0 && len(eval) > 0 {
		r := map[string]string{
			"eval": eval,
		}
		b, _ := json.Marshal(&r)
		stream.Notify("data", b)
		return
	}

	if len(url) == 0 {
		url = "/"
	}
	r := map[string]string{
		"url": url,
	}
	b, _ := json.Marshal(&r)
	stream.Notify("data", b)
}

func browserAgentProxy(w http.ResponseWriter, req *http.Request) {
	if len(proxying) == 0 {
		w.Write([]byte(page(script, "")))
		return
	}

	p := req.URL.Path
	q := req.URL.Query().Encode()

	rb, err := ioutil.ReadAll(req.Body)
	req.Body.Close()
	if err != nil {
		message := fmt.Sprintf("error proxying %s: %s", proxying, err)
		log.Println(message)
		w.Write([]byte(page(script, message)))
		return
	}

	u := proxying
	if len(p) > 0 {
		if !strings.HasSuffix(u, "/") {
			u += "/"
		}
		p = strings.TrimPrefix(p, "/")
		u += p
	}
	if len(q) > 0 {
		u += "?" + q
	}

	r, err := http.NewRequest(req.Method, u, bytes.NewReader(rb))

	for h, v := range req.Header {
		for _, v2 := range v {
			r.Header.Add(h, v2)
		}
	}

	res, err := http.DefaultClient.Do(r)
	if err != nil {
		message := fmt.Sprintf("error proxying %s: %s", proxying, err)
		log.Println(message)
		w.Write([]byte(page(script, message)))
		return
	}

	b, err := ioutil.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		message := fmt.Sprintf("error proxying %s: %s", proxying, err)
		log.Println(message)
		w.Write([]byte(page(script, message)))
		return
	}

	sb := string(b)

	sb = strings.Replace(sb, "</head>", script+"</head>", -1)

	for h, v := range res.Header {
		for _, v2 := range v {
			w.Header().Add(h, v2)
		}
	}

	clh := w.Header().Get("Content-Length")
	if len(clh) > 0 {
		cl, _ := strconv.Atoi(clh)
		cl += len(script)
		w.Header().Set("Content-Length", fmt.Sprintf("%d", cl))
	}

	w.Header().Del("Content-Security-Policy")
	w.Header().Set("Cache-Control", "private, max-age=0, no-cache")

	w.WriteHeader(res.StatusCode)
	w.Write([]byte(sb))
}

func page(script, message string) string {
	return `<html>
<head>
	<title>Peroxy: Surreptitious MITM</title>
` + script + `
	<style>
		html, body, iframe {
			width: 100%;
			height: 100%;
		}
	</style>
</head>
<body>
	<h1>Peroxy: Surreptitious MITM</h1>
` + message + `
</body>
</html>
`
}

var script = `
<script>
  src = new EventSource('/!-stream');
  src.addEventListener('message', function(e) {
    console.log("Event source message:")
    console.log(e)
		ev = JSON.parse(e.data)
		if (ev.url) {
			window.location.href=ev.url;
		}
		if (ev.eval) {
			eval(ev.eval);
		}
  }, false)
  src.addEventListener('error', function(e) {
    //console.log("Event source error:")
    //console.log(e)
  }, false)
  src.addEventListener('open', function(e) {
    //console.log("Event source open:")
    //console.log(e)
  }, false)
</script>
`
