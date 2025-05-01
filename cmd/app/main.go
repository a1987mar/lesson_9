package main

import (
	"fmt"
	"lesson4/pkg/users"
	"log/slog"

	"lesson4/pkg/documentstore"
)

func main() {

	//marshalExample()
	//unmarshalExample()
	slog.Info("App started")
	// lesson7()
	lesson9()
}

func marshalExample() {
	s := &documentstore.MyStruct{X: 15}
	doc, err := documentstore.MarshalDocument(s)
	if err != nil {
		fmt.Printf("failed to marshal document: %+v\n", err)
		return
	}
	fmt.Printf("marshaled document: %+v\n", doc)
}

func unmarshalExample() {
	doc := &documentstore.Document{Fields: map[string]documentstore.DocumentField{}}
	doc.Fields["X"] = documentstore.DocumentField{
		Type:  documentstore.DocumentFieldTypeNumber,
		Value: int(32),
	}

	s := &documentstore.MyStruct{}
	err := documentstore.UnmarshalDocument(doc, s)
	if err != nil {
		fmt.Printf("failed to unmarshal document: %+v\n", err)
		return
	}
	fmt.Printf("unmarshaled document: %+v\n", s)
}

func lesson7() {

	slog.Info("start app")
	doc1 := []documentstore.Document{
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u1"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Andrii"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u2"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Lubov"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u4"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Taras"},
			},
		},
		{
			Fields: map[string]documentstore.DocumentField{
				"id":   {Type: documentstore.DocumentFieldTypeString, Value: "u3"},
				"name": {Type: documentstore.DocumentFieldTypeString, Value: "Roman"},
			},
		},
	}

	usersCreated := make([]users.User, 0)
	slog.Info("add new store")
	st := documentstore.NewStore()
	ser := users.NewService(st)
	for i, doc := range doc1 {
		u1, err := ser.CreateUser(doc.Fields["id"].Value.(string), doc.Fields["name"].Value.(string), &doc)
		if err != nil {
			fmt.Println("", i+1, err)
		}
		usersCreated = append(usersCreated, *u1)
	}
	//
	fmt.Printf("%+v USERS \n", usersCreated)
	getU, err := ser.GetUser("u2")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%+v знайдено \n", getU)
	//
	delUser := ser.DeleteUser("u1")
	fmt.Println(delUser)

	uList, err := ser.ListUsers()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(uList)
	slog.Info("App done")
}

func lesson9() {
	// Створюємо нову колекцію
	store := documentstore.NewStore()
	_, users := store.CreateCollection("users", "id")

	// Додаємо документи
	users.Put(documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "1"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Andrii"},
		},
	})
	users.Put(documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "2"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Taras"},
		},
	})
	users.Put(documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "3"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Roman"},
		},
	})

	users.Put(documentstore.Document{
		Fields: map[string]documentstore.DocumentField{
			"id":   {Type: documentstore.DocumentFieldTypeString, Value: "4"},
			"name": {Type: documentstore.DocumentFieldTypeString, Value: "Stepan"},
		},
	})

	// Створюємо індекс по полю "name"
	if err := users.CreateIndex("name"); err != nil {
		fmt.Printf("Failed to create index: %v\n", err)
		return
	}

	// Виконуємо запит за допомогою індексу
	min := "Roman"
	max := "Stepan"
	results, err := users.Query("name", documentstore.QueryParams{
		MinValue: &min,
		MaxValue: &max,
		Desc:     false,
	})
	if err != nil {
		fmt.Printf("Query failed: %v\n", err)
		return
	}

	fmt.Println("Query results:")
	for _, doc := range results {
		id := doc.Fields["id"].Value.(string)
		name := doc.Fields["name"].Value.(string)
		fmt.Printf("User ID: %s, Name: %s\n", id, name)
	}

	// Видаляємо індекс
	if err := users.DeleteIndex("name"); err != nil {
		fmt.Printf("Failed to delete index: %v\n", err)
	}
}
