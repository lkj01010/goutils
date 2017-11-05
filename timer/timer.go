package timer
import "time"

type Timer struct {
	*time.Timer
}

func NewTimer() Timer {
	t := Timer{time.NewTimer(0)}
	return t
}

func (t Timer)SetInterval(interval time.Duration, callback func()) {
	go func() {
		for {
			t.Reset(interval)
			<- t.C
			callback()
		}
	}()
}

