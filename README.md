# mycron

### this a project which is based github.com/robfig/cron, it increases the two main functions: crontab & one-off event。

# Example:

you can define your own Job like this:

```
type FuncJob func()

func (f FuncJob) Run() { f() }
```

it must has the Run method.

or Like this:

```
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
```

After defined your own Job,you can AddEntry Like this:

```
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
```

then just use it like this:

```
        cron := New()
        cron.Start()
```

if you want to remove one special Entry, you must first Yencapsulate a function like this:
```
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
```

then use:
```
 cron.RemoveJob(zhuRemove(3))
```

it will find the Entry which has Job that Id is 3 and remove the Entry.


Complete example, please see the test case or mian.go

