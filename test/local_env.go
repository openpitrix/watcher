// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test

import (
	"os"
)

const (
	Handler       = "WATCHER_HANDLER"
	WatchedFile   = "WATCHER_WATCHED_FILE"
	Duration      = "WATCHER_DURATION"
	LogLevel      = "WATCHER_LOG_LEVEL"
	EtcdEndpoints = "WATCHER_ETCD_ENDPOINTS"
)

var Envs = map[string]string{
	Handler:       "UpdateOpenPitrixEtcd",
	WatchedFile:   "./test/global_config.yaml",
	Duration:      "5",
	LogLevel:      "debug",
	EtcdEndpoints: "127.0.0.1:2379",
}

func LocalEnv() {
	for k, v := range Envs {
		os.Setenv(k, v)
	}
}
