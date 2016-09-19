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
	dbPath1 = "/tmp/testdba"
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

	go func() {
		fmt.Println("~~~Start updating")
		for i := 0; i < 5000000; i++ {
			if (i % 10000) == 0 {
				fmt.Printf("SetOne %d\n", i)
			}
			if i < 230000 {
				st.SetOne([]byte("aaaaa"), []byte("old"))
			} else {
				st.SetOne([]byte("aaaaa"), []byte("new"))
			}
		}
	}()
	fmt.Println("Start saving snapshot")
	start := time.Now()
	snapshot := st.NewSnapshot()
	fmt.Printf("Done saving snapshot; time elapsed %v\n", time.Since(start))
	time.Sleep(10 * time.Millisecond)

	fmt.Println("Before using snapshot")
	result, err := st.Get([]byte("aaaaa"))
	if result != nil {
		fmt.Printf("Result: [%s]\n", string(result))
	} else {
		fmt.Println("Not found")
	}

	fmt.Println("After using snapshot")
	st.SetSnapshot(snapshot)
	result, err = st.Get([]byte("aaaaa"))
	if result != nil {
		fmt.Printf("Result: [%s]\n", string(result))
	} else {
		fmt.Println("Not found")
	}
}
