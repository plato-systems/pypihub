package redirect

import (
	"container/list"
	"fmt"
	"math/rand"
	"sync"
)

const maxSize = 1024

var (
	table = map[string]*entry{}
	fifo  = list.New()
	lock  = sync.Mutex{}
)

type entry struct {
	dest string
	user string
	pass string
}

func Register(src, dest, user, pass string) string {
	lock.Lock()
	defer lock.Unlock()

	if len(table) == maxSize {
		k := fifo.Remove(fifo.Front())
		delete(table, k.(string))
	}

	key := fmt.Sprintf("%x/%s", rand.Int(), src)
	table[key] = &entry{dest, user, pass}
	fifo.PushBack(key)

	return BaseURLPath + key
}
