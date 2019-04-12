// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package handler

import (
	"os"
	"context"

	yaml "gopkg.in/yaml.v2"
	"openpitrix.io/logger"
	"openpitrix.io/openpitrix/pkg/pi"

	"openpitrix.io/watcher/pkg/common"
)

func UpdateOpenPitrixEtcd() error {
	global := common.Global
	etcd := global.Etcd

	var content []byte
	var oldConfig []byte
	var newConfigMap common.AnyMap

	//read global_config file and convert to map
	content, newConfigMap, err := common.ReadYamlFile(global.WatchedFile)
	if err != nil {
		logger.Errorf(nil, "Failed to read yaml file %s: %+v", global.WatchedFile, err)
		return err //Do nothing if failed to read file
	}
	logger.Debugf(nil, "Global config yaml: %s", content)

	//get old config from etcd, and compare with global_config
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	err = etcd.Dlock(ctx, func() error {
		logger.Infof(nil, "Updating openpitrix etcd...")
		get, err := etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if err != nil {
			return err
		}
		logger.Debugf(nil, "Get count: %d", get.Count)
		logger.Debugf(nil, "Get: %+v", get.Kvs)
		var modified = new(bool)
		if get.Count == 0 {
			//init global_config if empty in etcd
			oldConfig = content
			*modified = true
		} else {
			//update old config from new config
			oldConfig = get.Kvs[0].Value
			oldConfigMap := make(common.AnyMap)
			err := yaml.Unmarshal(oldConfig, oldConfigMap)
			if err != nil {
				logger.Errorf(ctx, "Failed to unmarshal old config to map!")
				return err
			}

			ignoreKeyMap := make(common.AnyMap)
			err = yaml.Unmarshal([]byte(os.Getenv(common.IgnoreKeys)), ignoreKeyMap)
			if err != nil {
				logger.Errorf(ctx, "Failed to unmarshal ignore keys to map!")
				return err
			}

			common.CompareUpdateConfig(newConfigMap, oldConfigMap, ignoreKeyMap, modified)
			logger.Debugf(nil, "Modified: %t, config updated: %v", *modified, oldConfigMap)
			oldConfig, err = yaml.Marshal(oldConfigMap)
			if err != nil {
				logger.Errorf(nil, "Failed to convert old config from map to yaml: %+v", err)

			}
		}

		//put updated config to etcd if old config updated
		if *modified {
			_, err := etcd.Client.Put(ctx, pi.GlobalConfigKey, string(oldConfig))
			if err != nil {
				logger.Errorf(nil, "Failed to put data into etcd: %+v", err)
			}
		}
		return nil
	})

	if err != nil {
		logger.Errorf(nil, "Failed to update etcd: %+v", err)
	}
	return err
}
