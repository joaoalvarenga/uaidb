package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type Collection struct {
	name    string
	indexes map[string]*Index
}

func CreateCollection(name string) *Collection {

	os.MkdirAll(name, os.ModePerm)

	out := &Collection{
		name:    name,
		indexes: make(map[string]*Index),
	}
	out.indexes["_ID"] = &Index{}
	return out
}

func (c *Collection) CreateIndex(fieldName string) {
	c.indexes[fieldName] = &Index{}
}

func (c *Collection) IndexDocument(document Document) {
	c.indexes["_ID"].Insert(IndexNodeKey(document.ID), document.ID)
	documentData := make(map[string]interface{})
	json.Unmarshal([]byte(document.Raw), &documentData)
	hasher := sha256.New()
	for indexName := range c.indexes {
		if indexName == "_ID" {
			continue
		}
		value, ok := documentData[indexName]
		if !ok {
			continue
		}
		var buf bytes.Buffer
		valueEncoder := gob.NewEncoder(&buf)
		valueEncoder.Encode(value)
		data := buf.Bytes()
		hasher.Write(data)
		bs := hasher.Sum(nil)
		num, _ := binary.Uvarint(bs)
		c.indexes[indexName].Insert(IndexNodeKey(num), document.ID)
	}

}

func (c *Collection) AddDocument(document Document) {
	filename := filepath.Join(c.name, strconv.FormatUint(document.ID, 10))
	f, err := os.Create(filename)
	if err != nil {
		log.Fatal("Couldn't open file")
	}

	enc := gob.NewEncoder(f)
	enc.Encode(document)
	f.Close()

	c.IndexDocument(document)
}

func (c *Collection) GetByID(id int) *Document {
	filename := filepath.Join(c.name, strconv.Itoa(id))
	f, err := os.Open(filename)
	if err != nil {
		log.Fatal("Couldn't open file")
	}

	dec := gob.NewDecoder(f)
	document := Document{}
	dec.Decode(&document)
	f.Close()
	return &document
}
