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

const defaultRootFolderName = "ggnetwork"

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
		Filename: hashStr,
	}
	 
}

type PathTransformFunc func(string) PathKey

type PathKey struct {
	Pathname string
	Filename string
}

func (p PathKey) FullPath(root string) string {

	if len(root) == 0 {
		return fmt.Sprintf("%s/%s",p.Pathname,p.Filename)
	} else {
		return fmt.Sprintf("%s/%s/%s",root, p.Pathname,p.Filename)
	}
}

type StoreOpts struct {
	Root			  string // Root is the folder name of the root, containing all the files/folder of the system
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) PathKey {
	return PathKey{
		Pathname: key,
		Filename: key}
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store{
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}

	if len(opts.Root) == 0 {
		opts.Root = defaultRootFolderName
	}
	return &Store{
		StoreOpts: opts,
	}
}

func ( s*Store) Has (key string) bool {
	PathKey:=s.PathTransformFunc(key)

	_, err := os.Stat(PathKey.FullPath(s.Root))
	if err != nil{
		return false
	}

	return true
}

func (p PathKey) FirstPathName() string{
	paths := strings.Split(p.Pathname, "/")

	if len(paths) == 0 {
		return ""
	}
	return paths[0]
}

func (s *Store) Delete(key string) error{
	pathKey := s.PathTransformFunc(key)

	defer func(){
		log.Printf("delete [%s] from disk", pathKey.Filename)
	}()
	
	return os.RemoveAll(s.Root+"/"+pathKey.FirstPathName())
}

func (s *Store) Read(key string)(io.Reader, error){
	f,err := s.readStream(key)
	if err!= nil{
		return nil, err
	}

	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, f)
	
	return buf, err
}

func (s *Store) readStream(key string) (io.ReadCloser, error) {
	pathKey := s.PathTransformFunc(key)

	return os.Open(pathKey.FullPath(s.Root))
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathKey := s.PathTransformFunc(key)
	pathName := s.Root + "/" + pathKey.Pathname
	
	if err := os.MkdirAll(pathName, os.ModePerm); err!=nil {
		return err
	}

	buf := new (bytes.Buffer)
	io.Copy(buf, r)

	// filenameBytes := md5.Sum(buf.Bytes())
	// filename := hex.EncodeToString((filenameBytes[:]))
	// pathAndFilename := pathName + "/" + filename
	fullPath := pathKey.FullPath(s.Root)

	f, err := os.Create(fullPath)
	if err != nil {
		return err
	}

	n, err := io.Copy(f, buf)
	if err != nil {
		return err
	}

	log.Printf("written (%d) bytes to disk: %s", n, fullPath)
	return nil
}