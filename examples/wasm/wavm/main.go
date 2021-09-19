package main

import (
	"log"
	"time"
)

func main() {
	log.Print("Initializing WAVM")

	workers := make([]*WorkerT, 0, 128)

	//in := bufio.NewScanner(os.Stdin)
	for i := 0; i < 100; i++ {
		worker := Load(true, true)
		workers = append(workers, worker)

		//if i % 10 == 0 {
		//	in.Text()
		//}
	}

	//for _, worker := range workers {
	//	worker.Close()
	//}
	//worker.Close()

	println("done")
	time.Sleep(time.Hour)
}
