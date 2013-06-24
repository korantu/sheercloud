package cloud

import (
	"fmt"
	"log"
	"math/rand"
	_ "os/exec"
	"time"
)

type JobID string
type Result bool

var jobs = make(map[JobID]Result)

func DoJob(file string) JobID {
	id := JobID(fmt.Sprintf("[%d]", rand.Int()))
	log.Printf("Rendering job %s started for %s", id, file)
	jobs[id] = false
	go func() {
		time.Sleep(time.Second)
		jobs[id] = true
	}()
	return id
}

func JobDone(id JobID) (*Result, error) {
	if r, ok := jobs[id]; ok {
		return &r, nil
	} else {
		problem := CloudError("Unknown job")
		return nil, &problem
	}
}
