package runner

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Runner interface {
	Scheme() string
	Init() error
	Start() error
	Exit(context.Context) error
}

func Run(svr Runner) {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	// init signal
	if err := svr.Init(); err != nil {
		fmt.Printf("init server failed: %s", err.Error())
		return
	}
	Go(func() {
		if err := svr.Start(); err != nil {
			fmt.Printf("start server failed: %s", err.Error())
			os.Exit(1)
		}
	})

	s := <-sig
	switch s {
	case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
		ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
		defer cancel()
		if err := svr.Exit(ctx); err != nil {
			fmt.Printf("start server failed: %s", err.Error())
			return
		}
		time.Sleep(time.Second * 1)
		fmt.Println("server exit")
		return
	default:
		return
	}
}

func Go(fn func()) {
	var gw sync.WaitGroup
	gw.Add(1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				fmt.Println(err)
			}
		}()
		gw.Done()
		fn()
	}()
	gw.Wait()
}
