package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/jroimartin/monmq"
	"github.com/jroimartin/rpcmq"
	"github.com/zlowram/gsd"
	"golang.org/x/net/proxy"
)

var (
	a      *monmq.Agent
	config Config
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("usage: godan-worker <configuration file>\n")
		os.Exit(2)
	}
	var err error
	config, err = loadConfig(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	cmds := make(chan monmq.Command)
	monmq.Log = log.New(os.Stderr, "", log.LstdFlags)
	a = monmq.NewAgent("amqp://"+config.Monmq.Host+":"+config.Monmq.Port, config.Monmq.Exchange, config.Name)

	a.HardShutdownFunc = func(data []byte) ([]byte, error) {
		cmds <- monmq.HardShutdown
		return nil, nil
	}
	a.SoftShutdownFunc = func(data []byte) ([]byte, error) {
		cmds <- monmq.SoftShutdown
		return nil, nil
	}
	a.ResumeFunc = func(data []byte) ([]byte, error) {
		cmds <- monmq.Resume
		return nil, nil
	}
	a.PauseFunc = func(data []byte) ([]byte, error) {
		cmds <- monmq.Pause
		return nil, nil
	}

	if err := a.Init(); err != nil {
		log.Fatalf("Init: %v", err)
	}

	s := rpcmq.NewServer("amqp://"+config.Rpcmq.Host+":"+config.Rpcmq.Port, config.Rpcmq.Queue, config.Rpcmq.Exchange, config.Rpcmq.ExchangeType)
	if err := s.Register("getService", getService); err != nil {
		log.Fatalf("Register: %v", err)
	}

	s.DeliveryMode = rpcmq.Transient
	s.Parallel = config.Rpcmq.Parallel

	if err := s.Init(); err != nil {
		log.Fatalf("Init: %v", err)
	}

	paused := false
loop:
	for {
		switch <-cmds {
		case monmq.HardShutdown:
			log.Println("Hard shutdown...")
			break loop
		case monmq.SoftShutdown:
			log.Println("Soft shutdown...")
			s.Shutdown()
			a.Shutdown()
			break loop
		case monmq.Pause:
			if paused {
				continue
			}
			log.Println("Pause...")
			s.Shutdown()
			paused = true
		case monmq.Resume:
			if !paused {
				continue
			}
			log.Println("Resume...")
			if err := s.Init(); err != nil {
				log.Fatalln("Server init:", err)
			}
			paused = false
		}
	}
}

func getService(id string, data []byte) ([]byte, error) {
	a.RegisterTask(id)
	defer a.RemoveTask(id)

	ip := []string{strings.Split(string(data), ":")[0]}
	port := []string{strings.Split(string(data), ":")[1]}

	g := gsd.NewGsd(ip, port)

	userAgents := []string{
		"Mozilla/4.0 (compatible; MSIE 6.0; Windows NT 5.1; FSL 7.0.6.01001)",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; rv:12.0) Gecko/20100101 Firefox/12.0",
		"Mozilla/5.0 (X11; U; Linux x86_64; de; rv:1.9.2.8) Gecko/20100723 Ubuntu/10.04 (lucid) Firefox/3.6.8",
		"Opera/9.80 (Windows NT 5.1; U; en) Presto/2.10.289 Version/12.01",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_2) AppleWebKit/600.3.18 (KHTML, like Gecko) Version/8.0.3 Safari/600.3.18",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 8_1_3 like Mac OS X) AppleWebKit/600.1.4 (KHTML, like Gecko) Version/8.0 Mobile/12B466 Safari/600.1.4",
		"Mozilla/5.0 (Windows NT 6.1; WOW64; Trident/7.0; rv:11.0) like Gecko",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_5) AppleWebKit/600.3.18 (KHTML, like Gecko) Version/7.1.3 Safari/537.85.12",
		"Mozilla/5.0 (Linux; U; Android 4.0.3; ko-kr; LG-L160L Build/IML74K) AppleWebkit/534.30 (KHTML, like Gecko) Version/4.0 Mobile Safari/534.30",
		"Mozilla/6.0 (iPhone; U; CPU like Mac OS X; en) AppleWebKit/420+ (KHTML, like Gecko) Version/3.0 Mobile/1A543a Safari/419.3",
	}

	https := gsd.NewHttpsService()
	https.SetHeader("User-Agent", userAgents[rand.Intn(len(userAgents))])

	http := gsd.NewHttpService()
	http.SetHeader("User-Agent", userAgents[rand.Intn(len(userAgents))])

	services := []gsd.Service{
		https,
		http,
		gsd.NewTCPService(),
		gsd.NewTCPTLSService(),
	}
	g.AddServices(services)

	if config.Proxy.Host != "" && config.Proxy.Port != "" {
		var auth *proxy.Auth
		if config.Proxy.Username != "" && config.Proxy.Password != "" {
			auth = &proxy.Auth{
				User:     config.Proxy.Username,
				Password: config.Proxy.Password,
			}
		} else {
			auth = &proxy.Auth{}
		}
		g.SetProxy(config.Proxy.Host+":"+config.Proxy.Port, auth)
	}

	results := g.Run(50)

	banners := make([]gsd.Banner, 0)
	for r := range results {
		banners = append(banners, r)
	}

	return json.Marshal(banners)
}
