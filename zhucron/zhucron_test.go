package zhucron

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

var n sync.WaitGroup

type MyJob struct {
	Id      int
	Person  string
	RunFunc func(string, int) error
}

func (myjob *MyJob) Run() {
	myjob.RunFunc(myjob.Person, myjob.Id)
}

func NewMyJob(id int, name string) *MyJob {
	myjob := &MyJob{
		Id:     id,
		Person: name,
	}
	myjob.RunFunc = func(name string, id int) error {
		fmt.Println(name, id)
		n.Done()
		return nil
	}
	return myjob
}

func tt(t *testing.T) {
	ticker := time.NewTicker(time.Second)

	for t1 := range ticker.C {
		fmt.Println(t1)
	}

}

func TestConcurrent(t *testing.T) {

	now := time.Now().Local()
	d, _ := time.ParseDuration("20s")
	endtime := now.Add(d)

	cron := New()
	cron.Start()

	i := 0
	for {
		i++

		myjob := NewMyJob(i, strconv.Itoa(i))
		sub, _ := time.ParseDuration(strconv.Itoa(i*2) + "s")
		oneoff, _ := NewOneOffFrom(endtime, sub, myjob)
		go cron.AddEntry(oneoff)
		n.Add(1)
		if i > 5 {
			break
		}
	}

	// myjob1 := NewMyJob(1, "zhulilei0104")
	// n.Add(1)
	// myjob2 := NewMyJob(2, "suhan0825")
	// n.Add(1)
	// myjobcron := NewMyJob(3, "this is a cron")
	// n.Add(1)

	//mycron1, err := NewMyCronFrom("*/2 * * * * *", myjobcron)
	/*
		if err != nil {
			panic(err)
		}

		sub1, _ := time.ParseDuration("40s")
		oneoff1, err := NewOneOffFrom(endtime, sub1, myjob1)
		if err != nil {
			panic(err)
		}

		sub2, _ := time.ParseDuration("50s")
		oneoff2, err := NewOneOffFrom(endtime, sub2, myjob2)

		go cron.AddEntry(mycron1) //自己add的时候需要使用lock
		go cron.AddEntry(oneoff1)
		go cron.AddEntry(oneoff2)
	*/

	go tt(t)

	n.Wait()

	//select {} //阻塞主线程不退出
}
