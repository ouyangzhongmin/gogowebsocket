package timingwheel

import (
	"fmt"
	"testing"
	"time"
)

func TestTimingWheel(t *testing.T) {
	w := NewTimingWheel(1*time.Second, 10)
	ch := w.After(2 * time.Second)
	for {
		select {
		case <-ch:
			fmt.Println("到时间了")
			ch = w.After(2 * time.Second)
		}
	}
	time.Sleep(time.Minute * 5)
}

func TestTimingWheel2(t *testing.T) {
	w := NewTimingWheel(100*time.Millisecond, 20)
	for i := 0; i < 1000; i++ {
		go func() {
			for {
				select {
				case <-w.After(500 * time.Millisecond):
					fmt.Println("到时间了")
					return
				}
			}
		}()
	}

	time.Sleep(time.Minute * 5)

}
