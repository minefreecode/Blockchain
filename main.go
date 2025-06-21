package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
)

// Block Блок для хеширования данных
type Block struct {
	Pos       int
	Data      BookCheckout
	Timestamp string
	Hash      string // Хеш-значение в блоке. Получается из Data + Время + Pos
	PrevHash  string
}

// BookCheckout Данные проверки книги
type BookCheckout struct {
	BookID       string `json:"book_id"`
	User         string `json:"user"`
	CheckoutDate string `json:"checkout_date"`
	IsGenesis    bool   `json:"is_genesis"`
}

type Book struct {
	ID          string `json:"id"`
	ISBN        string `json:"isbn"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	PublishData string `json:"publish_date"`
}

type Blockchain struct {
	blocks []*Block
}

func (b *Block) generateHash() {
	bytes, _ := json.Marshal(b.Data) //Перевод в байты

	data := strconv.Itoa(b.Pos) + b.Timestamp + string(bytes) //Соединение данных для хеширования в строку

	hash := sha256.New()                       // Объект хеша
	hash.Write([]byte(data))                   //Запись данных в хеш
	b.Hash = hex.EncodeToString(hash.Sum(nil)) //Получение строки хеша в види строки и заненсение его в блок
}

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)

}

func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := new(Block)

	return block
}
