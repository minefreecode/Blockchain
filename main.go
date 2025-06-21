package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strconv"
	"time"
)

// Block Блок для хеширования данных
type Block struct {
	Pos       int // Номер позиции в блокчейне
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

// Проверка, являет ил хеш правильным
func (b *Block) validateHash(hash string) bool {
	b.generateHash()      //Сгенерировать новый хеш
	return b.Hash == hash //Проверка совпадения хешей
}

func (bc *Blockchain) AddBlock(data BookCheckout) {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	block := CreateBlock(prevBlock, data)
	if validBlock(block, prevBlock) {
		bc.blocks = append(bc.blocks, block)
	}
}

// CreateBlock Создание нового блока
func CreateBlock(prevBlock *Block, data BookCheckout) *Block {
	block := new(Block)                   // Выделение памяти для структуры
	block.Pos = prevBlock.Pos + 1         // Установка позиции из предыдущего блока
	block.Data = data                     //Установка данных о проверке книги
	block.Timestamp = time.Now().String() // Время
	block.PrevHash = prevBlock.Hash       //Сохранение хеша предыдущего блока
	block.generateHash()                  // Генерация нового хеша для блока

	return block
}

// Валидация блока
// Проверки:
// 1)Валиден ли хег предыдущего блока, есть ли связь со следущим по хеши;
// 2)проверка хеша самого блока
// 3) Проверка позиции Pos предыдущего блока
func validBlock(block *Block, prevBlock *Block) bool {
	if block.PrevHash != prevBlock.Hash {
		return false
	}
	if !block.validateHash(block.Hash) {
		return false
	}
	if prevBlock.Pos+1 != block.Pos {
		return false
	}
	return true
}
