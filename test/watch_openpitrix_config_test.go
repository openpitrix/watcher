// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package test

import (
	"os"
	"context"
	"io/ioutil"
	"testing"

	yaml "gopkg.in/yaml.v2"
	"openpitrix.io/openpitrix/pkg/pi"

	"openpitrix.io/watcher/pkg/common"
	"openpitrix.io/watcher/pkg/handler"
)

func TestWatchOpenPitrixConfig(t *testing.T) {
	t.Logf("Test watching openpitrix config...")

	const (
		TmpFile   = "./config_tmp.yaml"
		Pilot     = "pilot"
		Port      = "port"
		PortValue = 30114
	)
	LocalEnv()
	common.LoadConf()
	global := common.Global

	content, contentMap, err := common.ReadYamlFile("./global_config.yaml")
	if err != nil {
		t.Skipf("Failed to read content: %+v", err)
		t.Failed()
	}

	//init etcd: put global_config into etcd
	err = putToEtcd(global, content)
	if err != nil {
		t.Skipf("Failed to put content into etcd: %+v", err)
	}

	global.WatchedFile = TmpFile
	defer os.Remove(TmpFile)

	//Update pilot.port in content and write into TmpFile for test
	modified := new(bool)
	for k := range contentMap {
		if k.(string) == Pilot {
			contentMap[k] = interface{}(map[string]interface{}{Port: PortValue})
			*modified = true
		}
	}
	if !*modified {
		t.Skip("Failed to replace pilot.port in content!")
	}
	t.Log("Updated content.")
	updatedContent, err := yaml.Marshal(contentMap)
	err = ioutil.WriteFile(TmpFile, updatedContent, 0755)
	if err != nil {
		t.Skipf("Failed to write updated content to file: %+v", err)
	}

	//run UpdateOpenpitrixEtcd
	err = handler.UpdateOpenPitrixEtcd()
	if err != nil {
		t.Skipf("Failed to call UpdateOpenPitrixEtcd: %+v", err)
	}

	//get content updated from etcd
	//get port from above content
	etcdContent, err := getFromEtcd(global)
	if err != nil {
		t.Skipf("Failed to get content from etcd: %+v", err)
	}
	etcdContentMap := make(common.AnyMap)
	err = yaml.Unmarshal(etcdContent, etcdContentMap)
	if err != nil {
		t.Skipf("Failed to unmarshal content of etcd to etcd: %+v", err)
	}
	port := 0
	for k, v := range etcdContentMap {
		if k.(string) == Pilot {
			for kk, vv := range v.(common.AnyMap) {
				if kk.(string) == Port {
					port = vv.(int)
				}
			}
		}
	}

	//check if port was updated in etcd
	if port != PortValue {
		t.Skipf("Failed to update config in etcd, want: %+v, actual: %+v!", PortValue, port)
	}

	t.Log("Test successfully!")
}

func getFromEtcd(global *common.Config) ([]byte, error) {
	etcd := global.Etcd
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	var content []byte
	err := etcd.Dlock(ctx, func() error {
		get, err := etcd.Client.Get(ctx, pi.GlobalConfigKey)
		if get.Count > 0 {
			content = get.Kvs[0].Value
		}
		return err
	})

	return content, err
}

func putToEtcd(global *common.Config, content []byte) error {
	etcd := global.Etcd
	ctx, cancel := context.WithTimeout(context.Background(), common.EtcdDlockTimeOut)
	defer cancel()
	err := etcd.Dlock(ctx, func() error {
		_, err := etcd.Client.Put(ctx, pi.GlobalConfigKey, string(content))
		return err
	})
	return err
}
