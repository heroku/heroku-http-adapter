package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
)

func main() {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal(err)
	}
	redirectURL := "http://127.0.0.1:8080"

	proxy := New(redirectURL)

	go func() {
		fmt.Println("==============================")
		fmt.Printf("Starting HTTP Proxy server on port %d\n", port)
		fmt.Printf("Sending Requests to %s\n", redirectURL)
		fmt.Println("==============================")

		http.HandleFunc("/", proxy.handle)
		log.Fatal(http.ListenAndServe(":"+strconv.Itoa(port), nil))
	}()

	command := exec.Command(os.Args[1], os.Args[2:]...)
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	command.Stderr = os.Stderr
	// set PORT=8080 always so the http server binds on the port we are sending traffic to
	os.Setenv("PORT", "8080")
	command.Env = os.Environ()

	if err := command.Start(); err != nil {
		panic(err)
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	command.Wait()

	// todo: proxy shutdown?
}

type HerokuProxy struct {
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func New(target string) *HerokuProxy {
	url, _ := url.Parse(target)
	return &HerokuProxy{target: url, proxy: httputil.NewSingleHostReverseProxy(url)}
}

func (p *HerokuProxy) handle(w http.ResponseWriter, r *http.Request) {
	p.proxy.Transport = &proxyTransport{}
	p.proxy.ServeHTTP(w, r)
}

type proxyTransport struct {
}

func (t *proxyTransport) RoundTrip(request *http.Request) (*http.Response, error) {
	if request.Header.Get("Content-Type") == "application/cloudevents" {
		dataContentType := request.Header.Get("ce-datacontenttype")
		if dataContentType == "" {
			dataContentType = "application/json"
		}
		request.Header.Set("Content-Type", dataContentType)
	}

	response, err := http.DefaultTransport.RoundTrip(request)
	return response, err
}
