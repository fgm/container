package main

import (
	"fmt"
	"io"
	"time"

	"github.com/fgm/container"
	"github.com/fgm/container/queue"
)

type Element int

func consumer(w io.Writer, wq container.WaitableQueue[Element]) {
	lq := wq.(container.Countable) // This implementation provides Countable, so the assertion cannot fail.

	for {
		_, ok := <-wq.WaitChan()
		if !ok {
			fmt.Fprintf(w, "Consumer exiting on WaitableQueue closure, %d remaining in queue\n", lq.Len())
			return
		}
		fmt.Fprintln(w, "Queue might not be empty, looking for items")
		var (
			item Element
			wqs  container.WaitableQueueState
		)
		for ok {
			item, ok, wqs = wq.Dequeue()
			if !ok {
				fmt.Fprintln(w, "Back to wait")
				break
			}
			fmt.Fprintf(w, "Received: %v, %d in queue, status %s\n", item, lq.Len(), wqs)
			time.Sleep(30 * time.Millisecond) // Consume more slowly than producer
		}
	}
}

func producer(w io.Writer, wq container.WaitableQueue[Element]) {
	delay := time.After(3 * time.Second)
	for i := range 60 {
		wqs := wq.Enqueue(Element(i))
		fmt.Fprintf(w, "Sent: %v, status %s\n", i, wqs)
		time.Sleep(2 * time.Millisecond)
	}
	<-delay
	wq.Close()
}

func realMain(w io.Writer, sizeHint, low, high int) int {
	wq, err := queue.NewWaitableQueue[Element](sizeHint, low, high)
	if err != nil {
		fmt.Fprintf(w, "Failed to create WaitableQueue: %v\n", err)
		return 1
	}

	go consumer(w, wq)
	producer(w, wq)
	time.Sleep(100 * time.Millisecond) // Leave time for consumer() to display its exit message.
	return 0
}
