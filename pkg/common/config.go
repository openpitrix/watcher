// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package common

import (
	"fmt"

	"github.com/koding/multiconfig"
	"openpitrix.io/logger"
)

const (
	ConfigPrefix = "WATCHER"
	//The env name of ignore keys that for update openpitrix etcd;
	//You can set env for that, like(string):
	// {
	//   "runtime": true,
	//   "cluster": {
	//     "registry_mirror": true
	//   }
	// }
	IgnoreKeys = "IGNORE_KEYS"
)

type Config struct {
	WatchedFile string `default:"/opt/global_config.yaml"` //The file that need to be watched
	Duration    int64  `default:"10"`                      //The duration for watcher polling cycle which repeats
	Handler     string `default:"UpdateOpenPitrixEtcd"`    //The action func name to run when files change
	LogLevel    string `default:"info"`
	Etcd        *Etcd
}

var Global = new(Config)

func LoadConf() {
	loader := multiconfig.MultiLoader(
		&multiconfig.TagLoader{},
		&multiconfig.EnvironmentLoader{Prefix: ConfigPrefix, CamelCase: true},
	)
	//get config from env
	Global.Etcd = &Etcd{}
	err := loader.Load(Global.Etcd)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load etcd config: %+v", err)
		panic(errMsg)
	}
	err = loader.Load(Global)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to load global config: %+v", err)
		panic(errMsg)
	}

	logger.SetLevelByString(Global.LogLevel)
	logger.Debugf(nil, "Etcd config: %+v", Global.Etcd)
	logger.Debugf(nil, "Global config: %+v", Global)
}

type NilError struct {
	msg string
}

func (e NilError) Error() string {
	return e.msg
}
