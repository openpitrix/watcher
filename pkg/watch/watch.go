// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package watch

import (
	"fmt"
	"time"

	"github.com/radovskyb/watcher"
	"openpitrix.io/logger"

	"openpitrix.io/watcher/pkg/common"
	"openpitrix.io/watcher/pkg/handler"
)

const (
	UpdateOpenPitrixEtcd = "UpdateOpenPitrixEtcd"
)

func Watch() {
	handle() //init global_config

	w := watcher.New()
	w.SetMaxEvents(1)
	w.FilterOps(watcher.Write, watcher.Create)

	global := common.Global
	w.Add(global.WatchedFile)
	logger.Infof(nil, "Watching file: %s", global.WatchedFile)

	go func() {
		for {
			select {
			case event := <-w.Event:
				logger.Infof(nil, "%+v", event)
				handle()
			case err := <-w.Error:
				panic(err)
			case <-w.Closed:
				logger.Errorf(nil, "The watcher is closed!")
				return
			}
		}
	}()

	duration := time.Duration(global.Duration) * time.Second

	err := w.Start(duration)

	if err != nil {
		logger.Errorf(nil, "Failed to start watching: %+v", err)
		panic(err)
	}
}

func handle() {
	switch common.Global.Handler {
	case UpdateOpenPitrixEtcd:
		handler.UpdateOpenPitrixEtcd()
	default:
		msg := fmt.Sprintf("The func %s not exist!", common.Global.Handler)
		panic(msg)
	}
}
