package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/dgryski/go-farm"

	"github.com/dgraph-io/dgraph/loader"
	"github.com/dgraph-io/dgraph/posting"
	"github.com/dgraph-io/dgraph/rdb"
	"github.com/dgraph-io/dgraph/store"
	"github.com/dgraph-io/dgraph/uid"
)

func TestMergeFolders(t *testing.T) {
	logrus.SetLevel(logrus.DebugLevel)
	rootDir, err := ioutil.TempDir("", "storetest_")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(rootDir)

	dir1, err := ioutil.TempDir(rootDir, "dir_")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(dir1)

	dir2, err := ioutil.TempDir(rootDir, "dir_")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(dir2)

	destDir, err := ioutil.TempDir("", "dest_")
	if err != nil {
		t.Error(err)
		return
	}
	defer os.RemoveAll(destDir)

	ps1, err := store.NewStore(dir1)
	if err != nil {
		t.Error(err)
		return
	}

	ps2, err := store.NewStore(dir2)
	if err != nil {
		t.Error(err)
		return
	}

	list := []string{"alice", "bob", "mallory", "ash", "man", "dgraph",
		"ash", "alice"}
	var numInstances uint64 = 2
	posting.Init()
	uid.Init(ps1)
	loader.Init(nil, ps1)
	for _, str := range list {
		if farm.Fingerprint64([]byte(str))%numInstances == 0 {
			_, err := uid.GetOrAssign(str, 0, numInstances)
			if err != nil {
				fmt.Errorf("error while assigning uid")
			}
		}
	}
	uid.Init(ps2)
	loader.Init(nil, ps2)
	for _, str := range list {
		if farm.Fingerprint64([]byte(str))%numInstances == 1 {
			uid.Init(ps2)
			_, err := uid.GetOrAssign(str, 1, numInstances)
			if err != nil {
				fmt.Errorf("error while assigning uid")
			}

		}
	}
	posting.MergeLists(100)
	ps1.Close()
	ps2.Close()

	mergeFolders(rootDir, destDir)

	var opt *rdb.Options
	var ropt *rdb.ReadOptions
	opt = rdb.NewDefaultOptions()
	ropt = rdb.NewDefaultReadOptions()
	db, err := rdb.OpenDb(opt, destDir)
	it := db.NewIterator(ropt)

	count := 0
	for it.SeekToFirst(); it.Valid(); it.Next() {
		count++
	}

	if count != 6 { // There are totally 6 unique strings
		fmt.Errorf("Not all the items have been assigned uid")
	}
}
