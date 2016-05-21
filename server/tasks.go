package main

import (
	"encoding/json"
	"errors"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/zlowram/godan/model"
	"github.com/zlowram/godan/persistence"

	"github.com/jroimartin/rpcmq"
)

const (
	taskLimit = 0
)

type taskManager struct {
	ctrl   *ctrlTable
	wg     sync.WaitGroup
	client *rpcmq.Client
	pm     persistence.PersistenceManager
}

type Task struct {
	IPs   []string `json:"ips"`
	Ports []string `json:"ports"`
}

func newTaskManager(c *rpcmq.Client, p persistence.PersistenceManager) *taskManager {
	tm := &taskManager{
		ctrl:   &ctrlTable{m: make(map[string]bool)},
		client: c,
		pm:     p,
	}
	return tm
}

func (tm *taskManager) runTasks(task Task) {
	ports := expandPorts(task.Ports)

	tm.wg.Add(1)
	go func() {
		var ipExpander *IPExpander
		for _, cidr := range task.IPs {
			var ip net.IP
			var err error
			var ok bool
			if ipExpander, err = NewIPExpander(cidr); err != nil {
				ip = net.ParseIP(cidr)
				for _, port := range ports {
					task := ip.String() + ":" + port
					tm.sendTask([]byte(task))
				}
				continue
			}
			for {
				ip, ok = ipExpander.Next()
				if !ok {
					break
				}
				for _, port := range ports {
					task := ip.String() + ":" + port
					tm.sendTask([]byte(task))
				}
			}
		}

		tm.ctrl.setFull(true)
		tm.wg.Done()
	}()

	tm.wg.Add(1)
	go func() {
		for !tm.ctrl.done() {
			r := <-tm.client.Results()
			banners := make([]model.Banner, 100)
			if err := json.Unmarshal(r.Data, &banners); err != nil {
				log.Println(err)
			}
			if r.Err != "" {
				log.Println(r.Err)
			}

			for _, b := range banners {
				if b.Error != "" {
					continue
				}
				tm.pm.SaveBanner(b)
			}

			tm.ctrl.insert(r.UUID, true)
		}
		tm.wg.Done()
	}()
	tm.wg.Wait()
}

func (tm *taskManager) sendTask(task []byte) {
	if taskLimit > 0 {
		for {
			sent, _ := tm.ctrl.status()
			if sent < taskLimit {
				break
			}
			<-time.After(2 * time.Second)
		}
	}
	uuid, err := tm.client.Call("getService", task, 0)
	if err != nil {
		log.Println("Call:", err)
	}

	tm.ctrl.insert(uuid, false)
}

type ctrlTable struct {
	sync.RWMutex
	isFull bool
	m      map[string]bool
}

func (ct *ctrlTable) insert(uuid string, v bool) {
	ct.Lock()
	defer ct.Unlock()

	ct.m[uuid] = v
}

func (ct *ctrlTable) done() bool {
	ct.RLock()
	defer ct.RUnlock()

	all := true
	for _, v := range ct.m {
		if !v {
			all = false
			break
		}
	}
	return all && ct.isFull
}

func (ct *ctrlTable) setFull(b bool) {
	ct.Lock()
	defer ct.Unlock()

	ct.isFull = b
}

func (ct *ctrlTable) status() (int, int) {
	ct.Lock()
	defer ct.Unlock()

	sent := 0
	completed := 0
	for _, v := range ct.m {
		if v {
			completed++
		} else {
			sent++
		}
	}
	return sent, completed
}

func expandPorts(ports []string) []string {
	var portList []string

	for _, i := range ports {
		if strings.Contains(i, "-") {
			sp := strings.Split(i, "-")
			prange, err := portRange(sp[0], sp[1])
			if err != nil {
				log.Fatal(err)
			}
			portList = append(portList, prange...)
		} else {
			portList = append(portList, i)
		}
	}
	return portList
}

func portRange(a, b string) ([]string, error) {
	var ports []string

	n, _ := strconv.Atoi(a)
	m, _ := strconv.Atoi(b)

	if n >= m {
		return ports, errors.New("First parameter cannot be equal or greater than second")
	}

	for i := n; i <= m; i++ {
		ports = append(ports, strconv.Itoa(i))
	}

	return ports, nil
}
