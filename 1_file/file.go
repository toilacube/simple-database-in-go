package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

func randomInt() int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Int()
}

func SaveData(path string, data []byte) error {
	fp, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
	return err
	}
	defer fp.Close()
	_, err = fp.Write(data)
	return err
}

func SaveData2(path string, data []byte) error {
	temp_file := fmt.Sprintf("%s.%d", path, randomInt())
	file, err := os.OpenFile(temp_file, os.O_WRONLY | os.O_CREATE | os.O_EXCL, 0664)
	if err != nil {
		return err
	}

	defer func() {
		file.Close()
		if err != nil {
			os.Remove(temp_file)
		}
	}()

	// Save data to temp file
	if _, err = file.Write(data); err != nil {
		return err
	}

	// fsync: flush written data to disk
	if err = file.Sync(); err != nil {
		return err
	}

	// rename temp file to target path
	err = os.Rename(temp_file, path)
	return err

}

func main(){
	// write a file
	err := SaveData("test.txt", []byte("Hello, World zzzzz!"))
	fmt.Println(err)
}