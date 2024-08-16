package main

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

func CASPathTransformFunc(key string) PathKey{
	hash := sha1.Sum([]byte(key)) //[20]byte => []byte=>[:]
	hashStr := hex.EncodeToString(hash[:])

	blockSize := 5
	sliceLen := len(hashStr) / blockSize

	paths := make([]string, sliceLen)

	for i := 0; i < sliceLen; i++ {
		from, to := i*blockSize, (i*blockSize)+blockSize
		paths[i] = hashStr[from:to]
	}

	return PathKey{
		Pathname: strings.Join(paths,"/"),
		Original: hashStr,
	}
	 
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Original string
}

func (p PathKey) Filename() string {
	return fmt.Sprintf("%s/%s",p.Pathname,p.Original)
}

type StoreOpts struct {
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{Pathname: key,Original: key}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store{
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathName := pathKey.Pathname
	
	if err := os.MkdirAll(pathName, os.ModePerm); err!=nil {
		return err
	}

	buf := new (bytes.Buffer)
	io.Copy(buf, r)

	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString((filenameBytes[:]))
	// pathAndFilename := pathName + "/" + filename
	pathAndFilename := pathKey.Filename()

	f, err := os.Create(pathAndFilename)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, pathAndFilename)
	return nil
}