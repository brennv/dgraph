package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/dgraph-io/dgraph/store"
	"github.com/dgraph-io/dgraph/x"
)

const (
	dbPath1 = "/tmp/testdbn"
	dbPath2 = "/tmp/testdbn2"
)

func randStr() []byte {
	return []byte(strconv.Itoa(rand.Int()))
}

func main() {
	st, err := store.NewStore(dbPath1)
	x.Check(err)
	fmt.Println("Initial population")
	for i := 0; i < 1000000; i++ {
		st.SetOne(randStr(), randStr())
	}

	fmt.Println("NewCheckpoint")
	checkpoint, err := st.NewCheckpoint()
	x.Check(err)

	go func() {
		// Make sure we start after we save checkpoint.
		time.Sleep(10 * time.Millisecond)
		for i := 0; i < 123456; i++ {
			if (i % 10000) == 0 {
				fmt.Printf("SetOne %d\n", i)
			}
			if i >= 110000 {
				st.SetOne([]byte("aaaaa"), []byte("old"))
			} else {
				st.SetOne([]byte("aaaaa"), []byte("new"))
			}
		}
	}()
	fmt.Println("Start saving checkpoint")
	checkpoint.Save(dbPath2)
	fmt.Println("Done saving checkpoint")

	checkpoint.Destroy()

	st2, err := store.NewStore(dbPath2)
	result, err := st2.Get([]byte("aaaaa"))
	if result != nil {
		fmt.Printf("Result: [%s]\n", string(result))
	} else {
		fmt.Println("Not found")
	}
}
