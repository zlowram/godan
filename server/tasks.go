package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/jroimartin/rpcmq"
)

const (
	taskLimit = 0
)

type taskManager struct {
	ctrl     *ctrlTable
	wg       sync.WaitGroup
	client   *rpcmq.Client
	database *sql.DB
}

func newTaskManager(c *rpcmq.Client, d *sql.DB) *taskManager {
	tm := &taskManager{
		ctrl:     &ctrlTable{m: make(map[string]bool)},
		client:   c,
		database: d,
	}
	return tm
}

func (tm *taskManager) runTasks(tasks []string) {
	tm.wg.Add(1)
	go func() {
		for _, task := range tasks {
			if taskLimit > 0 {
				for {
					sent, _ := tm.ctrl.status()
					if sent < taskLimit {
						break
					}
					<-time.After(2 * time.Second)
				}
			}
			uuid, err := tm.client.Call("getService", []byte(task), 0)
			if err != nil {
				log.Println("Call:", err)
			}

			tm.ctrl.insert(uuid, false)
		}

		tm.ctrl.setFull(true)
		tm.wg.Done()
	}()

	tm.wg.Add(1)
	go func() {
		for !tm.ctrl.done() {
			r := <-tm.client.Results()
			banners := make([]banner, 100)
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
				_, err := tm.database.Exec("INSERT INTO banners (ip, port, service, content) VALUES (INET_ATON(?), ?, ?, ?)", b.Ip, b.Port, b.Service, b.Content)
				if err != nil {
					log.Fatal(err)
				}
			}

			tm.ctrl.insert(r.UUID, true)
		}
		tm.wg.Done()
	}()
	tm.wg.Wait()
}

type banner struct {
	Ip      string
	Port    string
	Service string
	Content string
	Error   string
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
