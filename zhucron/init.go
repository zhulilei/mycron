package zhucron

import "time"

func NewMyCronFrom(spec string, job Job) (*Entry, error) {
	//func NewMyCronFrom(spec string, cmd func()) (*Entry, error) {
	schedule, err := Parse(spec)
	if err != nil {
		return nil, err
	}
	mycron := &MyCron{
		Schedule: schedule,
		//Job:      FuncJob(cmd),
		Job: job,
	}

	entry := &Entry{
		MyType: 0,
		Single: mycron,
	}
	return entry, nil
}

//func NewOneOffFrom(endTime time.Time, sub time.Duration, cmd func()) (*Entry, error) {
func NewOneOffFrom(endTime time.Time, sub time.Duration, job Job) (*Entry, error) {
	oneoff := &OneOff{
		EndTime: endTime,
		Sub:     sub,
		//Job:     FuncJob(cmd),
		Job: job,
	}
	entry := &Entry{
		MyType: 1,
		Single: oneoff,
	}
	return entry, nil
}
