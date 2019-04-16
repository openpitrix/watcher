// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test

import (
	"strconv"
	"testing"

	"openpitrix.io/watcher/pkg/common"
)

func TestLoadConf(t *testing.T) {
	t.Logf("Test loading configurations...")

	LocalEnv()
	common.LoadConf()
	global := common.Global

	if global.Handler != Envs[Handler] {
		t.Errorf("WATCHER_HANDLER value is wrong: %s, want: %s", global.Handler, Envs[Handler])
	}

	if global.WatchedFile != Envs[WatchedFile] {
		t.Errorf("WATCHER_WATCHED_FILE value is wrong: %s, want: %s", global.WatchedFile, Envs[WatchedFile])
	}

	d, _ := strconv.ParseInt(Envs[Duration], 10, 64)
	if global.Duration != d {
		t.Errorf("WATCHER_DURATION value is wrong: %d, want: %s", global.Duration, Envs[Duration])
	}

	if global.LogLevel != Envs[LogLevel] {
		t.Errorf("WATCHER_LOG_LEVEL value is wrong: %s, want: %s", global.LogLevel, Envs[LogLevel])
	}

	if global.Etcd.Endpoints != Envs[EtcdEndpoints] {
		t.Errorf("WATCHER_ETCD_ENDPOINTS value is wrong: %s, want: %s", global.Etcd.Endpoints, Envs[EtcdEndpoints])
	}

	t.Log("Test successfully!")
}
