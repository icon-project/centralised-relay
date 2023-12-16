package lvldb

import (
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LVLDB struct {
	db *leveldb.DB
	sync.RWMutex
}

func NewLvlDB(path string, readonly bool) (*LVLDB, error) {
	opts := &opt.Options{
		ReadOnly:               readonly,
		OpenFilesCacheCapacity: 5000,
	}

	ldb, err := leveldb.OpenFile(path, opts)
	if err != nil {
		return nil, errors.Wrap(err, "levelDB.OpenFile fail")
	}

	return &LVLDB{db: ldb}, nil
}

func (db *LVLDB) GetByKey(key []byte) ([]byte, error) {
	db.RLock()
	defer db.RUnlock()
	return db.db.Get(key, nil)
}

func (db *LVLDB) SetByKey(key []byte, value []byte) error {
	db.Lock()
	defer db.Unlock()
	return db.db.Put(key, value, nil)
}

func (db *LVLDB) DeleteByKey(key []byte) error {
	db.Lock()
	defer db.Unlock()
	return db.db.Delete(key, nil)
}

func (db *LVLDB) NewIterator(prefix []byte) iterator.Iterator {
	db.RLock()
	defer db.RUnlock()
	return db.db.NewIterator(util.BytesPrefix(prefix), nil)
}

func (db *LVLDB) RemoveDbFile(filepath string) error {
	if err := os.Remove(filepath); err != nil {
		return errors.Wrapf(err, "unable to remove db file")
	}
	return nil
}

func (db *LVLDB) ClearStore() error {
	db.Lock()
	defer db.Unlock()

	iter := db.db.NewIterator(nil, nil)
	batch := new(leveldb.Batch)

	for iter.Next() {
		key := iter.Key()
		batch.Delete(key)
	}
	iter.Release()
	err := iter.Error()
	if err != nil {
		return nil
	}

	return db.db.Write(batch, nil)
}

// SnapShot snaphots the current state of the database
func (db *LVLDB) SnapShot() (*leveldb.Snapshot, error) {
	db.RLock()
	defer db.RUnlock()
	return db.db.GetSnapshot()
}

func (db *LVLDB) Close() error {
	db.Lock()
	defer db.Unlock()
	return db.db.Close()
}
