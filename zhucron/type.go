package zhucron

import (
	"log"
	"time"
)

// 删除job的检查函数，返回true则删除
type RemoveCheckFunc func(e *Entry) bool

type Cron struct {
	entries  []*Entry
	stop     chan struct{}
	add      chan *Entry
	remove   chan RemoveCheckFunc
	snapshot chan []*Entry
	running  bool
	ErrorLog *log.Logger
}

type Entry struct {
	MyType int    //0代表MyCron;1代表OneOff
	Single Single //Single is interface which has Run method
}

type Single interface {
	Run()
}
type MyCron struct {
	Schedule Schedule
	Next     time.Time
	Prev     time.Time
	Job      Job
}

type Schedule interface {
	Next(time.Time) time.Time
}

func (c *MyCron) Run() {
	c.Job.Run()
}

type OneOff struct {
	EndTime time.Time
	Sub     time.Duration
	Next    time.Time
	Process bool //true mean already processed,false mean has next
	Job     Job
}

func (o *OneOff) Run() {
	o.Job.Run()
}

type Job interface {
	Run()
}

type FuncJob func()

func (f FuncJob) Run() { f() }

type byTime []*Entry

func (s byTime) Len() int      { return len(s) }
func (s byTime) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
func (s byTime) Less(i, j int) bool {
	if s[i].getNext().IsZero() {
		return false
	}

	if s[j].getNext().IsZero() {
		return false
	}

	return s[i].getNext().Before(s[j].getNext())
}

/*
func main() {
	cmd1 := func() { fmt.Println("Every hour on the one hour") }
	funcjob1 := FuncJob(cmd1)

	cmd2 := func() { fmt.Println("Every hour on the two hour") }
	funcjob2 := FuncJob(cmd2)

	cron1 := &Cron{
		Next: "cron",
		Job:  funcjob1,
	}

	entry1 := &Entry{
		MyType: 1,
		Single: cron1,
	}

	oneoff1 := &OneOff{
		Next: "oneoff",
		Job:  funcjob2,
	}

	entry2 := &Entry{
		MyType: 2,
		Single: oneoff1,
	}

	fmt.Printf("%+v\n", entry1)
	fmt.Printf("%+v\n", entry2)

	entry1.Single.Run()

	entry2.Single.Run()

}
*/
