package zhucron

import (
	"log"
	"runtime"
	"sort"
	"time"

	"github.com/astaxie/beego"
)

func (c *Cron) setAllNext() {
	for _, entry := range c.entries {
		entry.setNext()
	}
}

func New() *Cron {
	return NewWithLocation(time.Now().Location())
}

func NewWithLocation(location *time.Location) *Cron {
	return &Cron{
		entries:  nil,
		add:      make(chan *Entry),
		remove:   make(chan RemoveCheckFunc),
		stop:     make(chan struct{}),
		snapshot: make(chan []*Entry),
		running:  false,
		ErrorLog: nil,
	}
}

func (c *Cron) AddEntry(entry *Entry) {
	if !c.running {
		c.entries = append(c.entries, entry)
		return
	}
	c.add <- entry
}

func (c *Cron) RemoveJob(cb RemoveCheckFunc) {
	c.remove <- cb
}

func (c *Cron) Entries() []*Entry {
	if c.running {
		c.snapshot <- nil
		x := <-c.snapshot
		return x
	}
	return c.entrySnapshot()
}

func (c *Cron) entrySnapshot() []*Entry {
	entries := []*Entry{}
	for _, e := range c.entries {
		entries = append(entries, e)
	}
	return entries
}

func (c *Cron) Start() {
	if c.running {
		return
	}
	c.running = true
	go c.run()
}

func (c *Cron) Run() {
	if c.running {
		return
	}
	c.running = true
	c.run()
}

func (c *Cron) logf(format string, args ...interface{}) {
	if c.ErrorLog != nil {
		c.ErrorLog.Printf(format, args...)
	} else {
		log.Printf(format, args...)
	}
}

func (c *Cron) runWithRecovery(e *Entry) {
	defer func() {
		if r := recover(); r != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			buf = buf[:runtime.Stack(buf, false)]
			c.logf("cron: panic running job: %v\n%s", r, buf)
		}
	}()
	e.Single.Run()
}

func (c *Cron) Stop() {
	if !c.running {
		return
	}
	c.stop <- struct{}{}
	c.running = false
}

func (c *Cron) run() {
	now := time.Now().Local()
	//beego.Error("now is", now)
	c.setAllNext()

	for {
		sort.Sort(byTime(c.entries))

		var timer *time.Timer
		if len(c.entries) > 0 {
			//	beego.Error("0 is", c.entries[0].getNext())
		}

		if len(c.entries) == 0 || c.entries[0].getNext().IsZero() {
			timer = time.NewTimer(100000 * time.Hour)
			//beego.Error("here")
		} else {
			timer = time.NewTimer(c.entries[0].getNext().Sub(now))
			//beego.Error("duration is", c.entries[0].getNext().Sub(now))
		}

		//beego.Error(timer)

		for {
			select {
			case now = <-timer.C:
				now := time.Now().Local()

				for k, e := range c.entries {
					if e.getNext().After(now) || e.getNext().IsZero() {
						break
					}
					beego.Error(e)
					go c.runWithRecovery(e)
					switch e.MyType {
					case 0:
						mycron := e.Single.(*MyCron)
						mycron.Prev = mycron.Next
						mycron.Next = mycron.Schedule.Next(now)
					case 1:
						c.entries = append(c.entries[:k], c.entries[k+1:]...)
						//beego.Error("length is", len(c.entries))
					}
				}

			case newEntry := <-c.add:
				timer.Stop()
				newEntry.setNext()
				c.entries = append(c.entries, newEntry)

			case <-c.snapshot:
				c.snapshot <- c.entrySnapshot()
				continue

			case <-c.stop:
				timer.Stop()
				return

			case cb := <-c.remove:
				newEntries := make([]*Entry, 0)
				for _, e := range c.entries {
					if !cb(e) {
						newEntries = append(newEntries, e)
					}
				}
				c.entries = newEntries
			}
			break
		}

	}
}
