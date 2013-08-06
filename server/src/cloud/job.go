package cloud

import (
	"fmt"
	"log"
	"math/rand"
	_ "os/exec"
	"strings"
	"time"
)

type JobID string
type Result bool

var jobs = make(map[JobID]Result)

func init() {
	rand.Seed(int64(time.Now().Nanosecond()))
}

func DoJob(file string) JobID {
	id := JobID(fmt.Sprintf("[%d]", rand.Int()))
	log.Printf("Rendering job %s started for %s", id, file)
	jobs[id] = false
	go func() {
		if strings.Contains(file, "scene.txt") {
			time.Sleep(time.Second)
		} else {
			log.Print("10 seconds...")
			time.Sleep(10 * time.Second)
		}
		jobs[id] = true
		log.Printf("Rendering job %s for %s completed", id, file)
	}()
	return id
}

func JobDone(id JobID) (*Result, error) {
	if r, ok := jobs[id]; ok {
		return &r, nil
	} else {
		return nil, &CloudError{"Unknown job"}
	}
}
