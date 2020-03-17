package main

//calivan0423

import (
	"fmt"

	"github.com/calivan0423/Go/dictionary/mydict"
)

func main() {
	dictionary := mydict.Dictionary{"first": "first word", "second": "second word"}

	err := dictionary.Add("hello", "Greeting")
	if err != nil {
		fmt.Println(err)
	}

	definition, err2 := dictionary.Search("hello")
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(definition)

	baseWord := "bye"
	dictionary.Add(baseWord, "Fisrt")
	err3 := dictionary.Update(baseWord, "Second")
	if err3 != nil {
		fmt.Println(err3)
	}
	word, _ := dictionary.Search(baseWord)
	fmt.Println(word)

	dictionary.Delete(baseWord)
	word2, err4 := dictionary.Search(baseWord)
	if err4 != nil {
		fmt.Println(err4)
	}
	fmt.Println(word2)

}
