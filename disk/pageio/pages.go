package main

import (
	"os"
)

// get disk page size
func getPageSize() int {
	return os.Getpagesize()
}

type page struct {
	data   []byte
	offset int
	dirty  bool
}

func main() {
	pageSize := getPageSize()
	println(pageSize/1024, "KB")

	filename := "test.db"
	stat, err := os.Stat(filename)
	if err != nil {
		println("file not found")
		return
	}

	ps := make([]page, stat.Size()/int64(pageSize)+1)
	handle, err := os.OpenFile(filename, os.O_RDWR, 0644)
	if err != nil {
		println("file open error")
		return
	}

	for i := range ps {
		ps[i].data = make([]byte, pageSize)
		ps[i].offset = i * pageSize
		handle.ReadAt(ps[i].data, int64(i*pageSize))
	}

	// let's assum some pages are modified
	for i := range ps {
		if i%2 == 0 {
			ps[i].data[0] = 0
		}
		ps[i].dirty = true
	}

	// write back to disk
	for i := range ps {
		if ps[i].dirty {
			handle.WriteAt(ps[i].data, int64(ps[i].offset))
		}
	}

	// let's assum some pages are removed
	for i := range ps {
		if i%3 == 0 {
			ps[i].data = nil
			ps[i].dirty = true
		}
	}

	nps := make([]page, 0)
	for i := range ps {
		if ps[i].data != nil {
			nps = append(nps, ps[i])
		}
	}

	// adjust the offsets
	for i := range nps {
		nps[i].offset = i * pageSize
	}

	// write back to disk
	for i := range nps {
		if nps[i].dirty {
			handle.WriteAt(nps[i].data, int64(nps[i].offset))
		}
	}

	// set file size
	handle.Truncate(int64(len(nps) * pageSize))

	handle.Close()
}
