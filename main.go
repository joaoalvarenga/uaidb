package main

import "fmt"

func main() {

	collection := CreateCollection("test")
	collection.CreateIndex("count")
	collection.AddDocument(Document{
		ID:  0,
		Raw: "{\"count\":10}",
	})
	collection.AddDocument(Document{
		ID:  0,
		Raw: "{\"count\":10}",
	})
	doc := collection.GetByID(0)
	fmt.Println(doc)

}
