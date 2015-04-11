package main

import (
	"bytes"
	"compress/zlib"
	"io"
	"io/ioutil"
	"log"
	"os"
)

var stream = []byte{0x73, 0x74, 0x72, 0x65, 0x61, 0x6d}
var endstream = []byte{0x65, 0x6e, 0x64, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d}

func main() {
	path := "/home/oneplus/Desktop/gopdf/cc5.pdf"
	/*content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(2)
	}
	fmt.Printf("%s\n", content)
	*/
	fd, err := os.Open(path)
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(2)
	}
	defer fd.Close()

	binary, err := ioutil.ReadAll(fd)
	if err != nil {
		log.Printf("%s", err.Error())
		os.Exit(2)
	}

	indexs := FindBytes(binary, stream)
	indexends := FindBytes(binary, endstream)

	indexs = CleanStreamIndexs(indexs, indexends)
	//fmt.Printf("\n%#v \n %#v", indexs, indexends)

	for i, val := range indexs {
		if i == 3 {
			zipcon := binary[val+len(stream)+1 : indexends[i]]
			Print(zipcon)
		}
		//fmt.Printf("%v\n\n", zipcon)
	}

}

func Print(buff []byte) {

	b := bytes.NewReader(buff)

	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, r)

	r.Close()
}

func CleanStreamIndexs(indexs []int, indexends []int) []int {
	i := 0
	indexsmax := len(indexs)
	for i < indexsmax {
		indexs[i] = indexs[i] + len(stream)
		i++
	}

	i = 0
	indexendsmax := len(indexends)
	for i < indexendsmax {
		indexends[i] = indexends[i] + len(endstream)
		i++
	}
	i = 0
	var realindexs []int
	for i < indexsmax {
		j := 0
		dup := false
		for j < indexendsmax {
			if indexs[i] == indexends[j] {
				dup = true
				break
			}
			j++
		}

		if !dup {
			realindexs = append(realindexs, indexs[i]-len(stream))
		}

		i++
	}

	return realindexs
}

func FindBytes(binary []byte, patten []byte) []int {

	var indexs []int
	binarylen := len(binary)
	pattenlen := len(patten)
	i := 0
	for index, _ := range binary {
		if index+pattenlen > binarylen {
			break
		}
		i = 0
		for binary[index+i] == patten[i] {
			if i == pattenlen-1 {
				indexs = append(indexs, index)
				break
			}
			i++
		}

	}

	return indexs
}
