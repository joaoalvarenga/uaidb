package main

import (
	"log"
	"strings"
)

const ORDER = 4
const MIN_ORDER = 2

type IndexNodeKey uint64

type IndexNode struct {
	Value    IndexNodeKey
	Position uint64
}

type IndexPage struct {
	N int
	R [ORDER + 1]IndexNode
	P [ORDER + 1]*IndexPage
}

type Index struct {
	Root *IndexPage
}

func CreateIndex() *Index {
	return &Index{
		Root: nil,
	}
}

func InsertInPage(ap *IndexPage, node IndexNode, apDir *IndexPage) {
	k := ap.N
	positionNotFound := k > 0
	for positionNotFound {
		if node.Value >= ap.R[k-1].Value {
			positionNotFound = false
			break
		}
		ap.R[k] = ap.R[k-1]
		ap.P[k+1] = ap.P[k]
		k--
		if k < 1 {
			positionNotFound = false
		}
	}
	ap.R[k] = node
	ap.P[k+1] = apDir
	ap.N++
}

func Insert(node IndexNode, ap *IndexPage, grown *bool, outputNode *IndexNode, outputPage **IndexPage) {
	i := 1
	if ap == nil {
		*grown = true
		*outputNode = node
		*outputPage = nil
		return
	}
	for i < ap.N && node.Value > ap.R[i-1].Value {
		i++
	}
	if node.Value == ap.R[i-1].Value {
		log.Fatal("Key already exists")
		*grown = false
		return
	}

	if node.Value < ap.R[i-1].Value {
		i--
	}
	Insert(node, ap.P[i], grown, outputNode, outputPage)
	if !*grown {
		return
	}
	if ap.N < ORDER {
		InsertInPage(ap, *outputNode, *outputPage)
		*grown = false
		return
	}
	apTemp := &IndexPage{
		N: 0,
	}
	if i < MIN_ORDER+1 {
		InsertInPage(apTemp, ap.R[ORDER-1], ap.P[ORDER])
		ap.N--
		InsertInPage(ap, *outputNode, *outputPage)
	} else {
		InsertInPage(apTemp, *outputNode, *outputPage)
	}
	for j := MIN_ORDER + 2; j <= ORDER; j++ {
		InsertInPage(apTemp, ap.R[j-1], ap.P[j])
	}
	ap.N = MIN_ORDER
	ap.P[0] = ap.P[MIN_ORDER+1]
	*outputNode = ap.R[MIN_ORDER]
	*outputPage = apTemp
}

func (t *Index) Insert(value IndexNodeKey, position uint64) {
	var grown bool
	var outputNode IndexNode
	var outputPage *IndexPage
	var apTemp *IndexPage
	node := IndexNode{
		Value:    value,
		Position: position,
	}
	Insert(node, t.Root, &grown, &outputNode, &outputPage)
	if grown {
		apTemp = &IndexPage{
			N: 1,
		}
		apTemp.R[0] = outputNode
		apTemp.P[1] = outputPage
		apTemp.P[0] = t.Root
		t.Root = apTemp
	}
}

func PrintNode(page *IndexPage, level int) {
	if page == nil {
		return
	}
	for i := 0; i < page.N; i++ {
		PrintNode(page.P[i], level+1)
		log.Printf("%s├── %d\n", strings.Repeat(" ", 4*level), page.R[i].Value)

	}
	if page.P[page.N] != nil {
		PrintNode(page.P[page.N], level+1)
	}

}

func Search(key IndexNodeKey, page *IndexPage) *IndexNode {
	i := 1
	if page == nil {
		return nil
	}
	for i < page.N && key > page.R[i-1].Value {
		i++
	}
	if key == page.R[i-1].Value {
		return &page.R[i-1]
	}
	if key < page.R[i-1].Value {
		return Search(key, page.P[i-1])
	}
	return Search(key, page.P[i])
}

func (t *Index) Search(key IndexNodeKey) *IndexNode {
	return Search(key, t.Root)
}

func (t *Index) Print() {
	if t.Root == nil {
		log.Printf("Empty index")
	}
	PrintNode(t.Root, 1)
}
