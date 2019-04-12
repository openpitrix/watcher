// Copyright 2019 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package common

import (
	"time"
	"context"
	"strings"

	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"go.etcd.io/etcd/clientv3/namespace"
	"openpitrix.io/openpitrix/pkg/pi"
	"openpitrix.io/logger"
)

const (
	EtcdDlockTimeOut = time.Second * 60
)

type Etcd struct {
	Endpoints string `default:"openpitrix-etcd:2379"` // Example: "localhost:2379,localhost:22379,localhost:32379"
	Client    *clientv3.Client
}

func (etcd *Etcd) NewEtcdClient(prefix string) {
	logger.Infof(nil, "New client of etcd: %+v...", etcd)
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   strings.Split(etcd.Endpoints, ","),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		logger.Errorf(nil, "Failed to connect etcd.")
		panic(err)
	}
	cli.KV = namespace.NewKV(cli.KV, prefix)
	etcd.Client = cli
}

type Mutex struct {
	*concurrency.Mutex
}

func (etcd *Etcd) NewMutex(ctx context.Context, dLockKey string) (*Mutex, error) {
	if etcd.Client == nil {
		return nil, NilError{"The etcd client is nil!"}
	}
	logger.Debugf(nil, "Create new session of etcd...")
	session, err := concurrency.NewSession(etcd.Client, concurrency.WithContext(ctx))
	if err != nil {
		return nil, err
	}
	return &Mutex{concurrency.NewMutex(session, dLockKey)}, nil
}

// Lock locks the mutex with a cancelable context. If the context is canceled
// while trying to acquire the lock, the mutex tries to clean its stale lock entry.
func (m *Mutex) Lock(ctx context.Context) error {
	logger.Debugf(ctx, "Lock etcd client...")
	return m.Mutex.Lock(ctx)
}

func (m *Mutex) Unlock(ctx context.Context) error {
	logger.Debugf(ctx, "Unlock etcd client...")
	return m.Mutex.Unlock(ctx)
}

type callback func() error

func (etcd *Etcd) Dlock(ctx context.Context, cb callback) error {
	etcd.NewEtcdClient(pi.EtcdPrefix)
	defer etcd.Client.Close()
	logger.Infof(ctx, "Create dlock with key [%s]", pi.DlockKey)
	mutex, err := etcd.NewMutex(ctx, pi.DlockKey)
	if err != nil {
		logger.Errorf(ctx, "Dlock lock error, failed to create mutex: %+v", err)
		panic(err)
	}
	err = mutex.Lock(ctx)
	if err != nil {
		logger.Errorf(ctx, "Dlock lock error, failed to lock mutex: %+v", err)
		return err
	}
	defer mutex.Unlock(ctx)
	logger.Debugf(ctx, "Running callback...")
	err = cb()
	return err
}
