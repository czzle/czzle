package main

import (
	"fmt"
	"time"

	"github.com/czzle/czzle/internal/pkg/appd"
)

type logger struct{}

func (logger) Log(kv ...interface{}) error {
	_, err := fmt.Println(kv...)
	return err
}

func main() {
	app := appd.New(
		appd.ProbePort("8889"),
		appd.WithLogger(new(logger)),
		appd.BeforeStart(
			func() error {
				time.Sleep(time.Second * 5)
				return nil
			},
		),
		appd.BeforeStop(
			func() error {
				time.Sleep(time.Second * 5)
				return nil
			},
		),
		appd.Start(
			func(ch chan<- error) {
				time.Sleep(time.Second * 10)
				ch <- nil
			},
		),
	)
	app.Run()
}
