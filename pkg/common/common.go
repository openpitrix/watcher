// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package common

import (
	"io/ioutil"
	"reflect"

	yaml "gopkg.in/yaml.v2"
	"openpitrix.io/logger"
)

type AnyMap map[interface{}]interface{}

//return: content contentMap error
func ReadYamlFile(f string) ([]byte, AnyMap, error) {
	//read global_config file and convert to map
	content, err := ioutil.ReadFile(f)
	if err != nil {
		logger.Errorf(nil, "Failed to read file %s!", f)
		return nil, nil, err
	}

	contentMap := make(AnyMap)
	err = yaml.Unmarshal(content, contentMap)
	if err != nil {
		logger.Errorf(nil, "Failed to Unmarshal yaml to map!")
	}
	return content, contentMap, err
}

//Base old config, update that from new config.
func CompareUpdateConfig(new, old AnyMap, ignoreKeys AnyMap, modified *bool) {
	for k, v := range old {
		kStr := k.(string)

		//check if k is in ignore keys
		var subIgnoreKeys AnyMap
		var t interface{}
		if ignoreKeys == nil || ignoreKeys[kStr] == nil {
			t = nil
		} else {
			t = reflect.TypeOf(ignoreKeys[kStr]).Kind()
		}

		if t == reflect.Bool && ignoreKeys[kStr].(bool) {
			logger.Infof(nil, "Ignore to update config: %s", kStr)
			continue //only in this condition, ignore update old config
		} else if t == reflect.Map {
			//get sub-ignore-keys
			subIgnoreKeys = ignoreKeys[kStr].(AnyMap)
		}

		if v == nil { //check if old value and new value are nil
			if new == nil || new[k] == nil {
				continue
			} else {
				logger.Infof(nil, "Updating, key: %s, old value: %v, new value: %v", k, v, new[k])
				//update old config from new config
				old[k] = new[k]
				continue
			}
		}

		switch reflect.TypeOf(v).Kind() {
		case reflect.Map:
			logger.Debugf(nil, "Key: %+v", k)
			CompareUpdateConfig(new[k].(AnyMap), v.(AnyMap), subIgnoreKeys, modified)
		default:
			if new[k] != v { //update old config from new config
				logger.Infof(nil, "Updating, key: %s, old value: %v, new value: %v", k, v, new[k])
				old[k] = new[k]
				*modified = true
			}
		}
	}
}
