package redirect

const maxSize = 1024

var table = map[string]string{}

func Register(src, dest string) {
	if len(table) == maxSize && table[src] == "" {
		for k := range table {
			delete(table, k)
			return
		}
	}
	table[src] = dest
}
