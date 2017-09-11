package main

import (
	"fmt"
	zhucron "mycron/zhucron"
	"time"

	"github.com/astaxie/beego"
)

// type FuncJob func()

// func (f FuncJob) Run() { f() }

type MyJob struct {
	Id      int
	Person  string
	RunFunc func(string, int) error
}

func NewMyJob(id int, name string) *MyJob {
	myjob := &MyJob{
		Id:     id,
		Person: name,
	}
	myjob.RunFunc = func(name string, id int) error {
		fmt.Println(name, id)
		return nil
	}
	return myjob
}

func (myjob *MyJob) Run() {
	myjob.RunFunc(myjob.Person, myjob.Id)
}

func t() {
	ticker := time.NewTicker(time.Second)

	for t := range ticker.C {
		fmt.Println(t)
	}

}

func zhuRemove(id int) (b zhucron.RemoveCheckFunc) {
	return func(e *zhucron.Entry) bool {
		switch e.MyType {
		case 0:
			mycron := e.Single.(*zhucron.MyCron)
			if v, ok := mycron.Job.(*MyJob); ok {
				if v.Id == id {
					beego.Error("remove", e)
					return true
				}
			} else {
				beego.Error("false")
			}
		case 1:
			oneoff := e.Single.(*zhucron.OneOff)
			if v, ok := oneoff.Job.(*MyJob); ok {
				if v.Id == id {
					beego.Error("remove", e)
					return true
				}
			} else {
				beego.Error("false")
			}
		}
		return false
	}
}

func main() {
	myjob1 := NewMyJob(1, "zhulilei0104")
	myjob2 := NewMyJob(2, "suhan0825")
	myjobcron := NewMyJob(3, "this is a cron")

	now := time.Now().Local()
	d, _ := time.ParseDuration("1m")
	endtime := now.Add(d)

	//cmd1 := func() { fmt.Println("Every hour on the one hour") }
	//cmd2 := func() { fmt.Println("Every hour on the two hour") }
	//cmd3 := func() { fmt.Println("Every hour on the three hour") }

	mycron1, err := zhucron.NewMyCronFrom("*/2 * * * * *", myjobcron)
	if err != nil {
		panic(err)
	}

	sub1, _ := time.ParseDuration("40s")
	oneoff1, err := zhucron.NewOneOffFrom(endtime, sub1, myjob1)
	if err != nil {
		panic(err)
	}

	sub2, _ := time.ParseDuration("50s")
	oneoff2, err := zhucron.NewOneOffFrom(endtime, sub2, myjob2)

	cron := zhucron.New()
	cron.Start()
	go cron.AddEntry(mycron1) //自己add的时候需要使用lock?用benchmark的时候没有遇到过,golang slice会有竞争么
	go cron.AddEntry(oneoff1)
	go cron.AddEntry(oneoff2)
	//cron.Start()

	go func() {
		timer := time.NewTimer(time.Second * 12)
		removeTimer := time.NewTimer(time.Second * 2)

		for {
			for _, e := range cron.Entries() {
				fmt.Printf("%+v\n", e)
				fmt.Printf("%+v\n", e.Single)
			}

			time.Sleep(1 * time.Second)
			select {
			case <-timer.C:
				//expire
				cron.Stop()
				time.Sleep(2 * time.Second)
				cron.Start()
			case <-removeTimer.C:
				beego.Error("here")
				cron.RemoveJob(zhuRemove(3))
			default:

			}
		}
	}()

	go t()

	select {} //阻塞主线程不退出

}
