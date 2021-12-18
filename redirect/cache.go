package redirect

import (
	"fmt"
	"time"
)

const maxSize = 1024

var table = map[string]*entry{}

type entry struct {
	dest string
	user string
	pass string
}

func Register(src, dest, user, pass string) string {
	if len(table) == maxSize {
		for k := range table {
			delete(table, k)
			break
		}
	}

	key := fmt.Sprintf("%x/%s", time.Now().UnixMilli(), src)
	table[key] = &entry{dest, user, pass}

	return BaseURLPath + key
}
