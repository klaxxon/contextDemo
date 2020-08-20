package main

import (
	"context"
	"log"
	"sync"
	"time"
)

const (
	PARENTWAIT = 1
	CHILDWAIT  = 10
)

func grandChild(ctx context.Context, wg *sync.WaitGroup, id int) {
	log.Printf("Grandchild %d of %d started\n", id&0xff, id>>8)
	// Just wait for the call home
	<-ctx.Done()
	log.Printf("Grandchild %d of %d done\n", id&0xff, id>>8)
	// Tell parent we are home
	wg.Done()
}

func child(ctx context.Context, wg *sync.WaitGroup, id int) {
	log.Printf("Child %d started\n", id)
	// For this, we will have our own context so we can cancel just this
	// child and all its children.
	mctx, cancel := context.WithCancel(ctx)
	defer cancel()
	mwg := &sync.WaitGroup{}
	// We will spawn our own children
	for a := 0; a < 3; a++ {
		mwg.Add(1)
		go grandChild(mctx, mwg, a|(id<<8))
	}
	select {
	case <-ctx.Done(): // In case parent calls EVERYONE home before we do
		break
	case <-time.After(time.Duration(CHILDWAIT+id) * time.Second): // We call our own kids home
		log.Printf("Child %d cancelling kids\n", id)
		// Tell kids to come home
		cancel()
	}
	// Wait for all of them to come back
	mwg.Wait()
	// Tell parent I am home
	log.Printf("Child %d done\n", id)
	wg.Done()
}

func main() {
	// So we can make sure all kids return
	wg := &sync.WaitGroup{}
	// So we can signal everyone
	mainCtx, cancel := context.WithCancel(context.Background())

	// Kick off our children
	for a := 0; a < 3; a++ {
		wg.Add(1)
		go child(mainCtx, wg, a)
		time.Sleep(10 * time.Millisecond)
	}
	time.Sleep(PARENTWAIT * time.Second)
	log.Println("Main parent cancelling!")
	// Tell EVERYONE to come home
	cancel()
	// Wait for them
	wg.Wait()
	log.Println("Finished!")
}
