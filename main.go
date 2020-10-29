package main

import "C"
import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	"openOEP/desktopCapture"
	"openOEP/singleton"
)

const (
	AV_PIX_FMT_YUV420P = 0
	NALUTypeVPS        = 32 // VPS
	NALUTypeSPS        = 33 // SPS
	NALUTypePPS        = 34 // PPS
)

var (
	vps []byte // VPS
	sps []byte // SPS
	pps []byte // PPS
)

func main() {
	// watch os signal
	sigCh := make(chan os.Signal, 16)
	signal.Notify(sigCh)

	// wait group
	var wg sync.WaitGroup

	// define and start workers
	workers := []func(){
		func() {
			// capture desktop
			defer wg.Done()
			ec := desktopCapture.StartCapture(25,
				0,
				0,
				1920,
				1080,
				960,
				540,
				AV_PIX_FMT_YUV420P)
			fmt.Printf("startCapture exited: %d\n", int(ec))
			close(singleton.X265Queue)
		},
		func() {
			// push stream
			defer wg.Done()
			for nal := range singleton.X265Queue {
				switch nal.Type {
				case NALUTypeVPS:
					vps = nal.Payload
					fmt.Println("--VPS--")
				case NALUTypeSPS:
					sps = nal.Payload
					fmt.Println("--SPS--")
				case NALUTypePPS:
					pps = nal.Payload
					fmt.Println("--PPS--")
				default:
					//fmt.Printf("%d,%d,%p\n", nal.Type, nal.Size, nal.Payload)
				}
				// push push push
			}
			sigCh <- os.Interrupt
		},
	}
	wg.Add(len(workers))
	for _, worker := range workers {
		go worker()
	}

	// wait for os signal
sigLoop:
	for sig := range sigCh {
		fmt.Println(sig)
		switch sig {
		case os.Kill, os.Interrupt:
			desktopCapture.StopCapture()
			signal.Stop(sigCh)
			break sigLoop
		}
	}

	// wait all workers done
	wg.Wait()
}
