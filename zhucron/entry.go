package zhucron

import "time"

func (entry *Entry) getNext() time.Time {
	switch entry.MyType {
	case 0:
		return entry.Single.(*MyCron).Next
	case 1:
		return entry.Single.(*OneOff).Next
	}
	neverTime, _ := time.Parse("2016/09/28-06:40", "1970/01/01-12:00")
	return neverTime
}

func (entry *Entry) setNext() {
	now := time.Now().Local()
	switch entry.MyType {
	case 0:
		entry.Single.(*MyCron).Next = entry.Single.(*MyCron).Schedule.Next(now)
	case 1:
		oneoff := entry.Single.(*OneOff)
		oneoff.Next = oneoff.EndTime.Add(-oneoff.Sub)
	}
}
