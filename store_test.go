package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"testing"
)



func TestPathTransformFunc(t *testing.T){
	key := "momsbestpicture"
	pathKey := CASPathTransformFunc(key)
	fmt.Println(pathKey)

	expectedOriginalKey := "6804429f74181a63c50c3d81d733a12f14a353ff"
	expectedPathName := "68044/29f74/181a6/3c50c/3d81d/733a1/2f14a/353ff"
	if pathKey.Pathname != expectedPathName {
		t.Errorf("have %s want %s", pathKey.Pathname , expectedPathName)
	}

	if pathKey.Filename != expectedOriginalKey {
		t.Errorf("have %s want %s", pathKey.Filename , expectedOriginalKey)
	}

}

func TestStoreCASPathTransformFunc(t * testing.T){
	opts := StoreOpts{
		PathTransformFunc : CASPathTransformFunc,
	}

	s := NewStore(opts)
	key := "momspecial"
	data  := []byte("some jpg bytes")
	if err := s.writeStream(key, bytes.NewReader(data)); err != nil {
		t.Error(err)
	}
	r, err := s.Read(key)
	if err != nil {
		t.Error(err)
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(string(b))
	if string(b) != string(data) {
		t.Errorf("want %s , have %s", data, b)
	}

}

func TestStoreDefaultPathTransformFunc(t * testing.T){
	opts := StoreOpts{
		PathTransformFunc : DefaultPathTransformFunc,
	}

	s := NewStore(opts)

	data  := bytes.NewReader([]byte("some jpg bytes"))
	if err := s.writeStream("myspecialpicture", data); err != nil {
		t.Error(err)
	}
}