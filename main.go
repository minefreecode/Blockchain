package main

import (
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
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

// Создание новой книги
func newBook(w http.ResponseWriter, r *http.Request) {
	var book Book //Выделение памяти

	if err := json.NewDecoder(r.Body).Decode(&book); err != nil { //Если не удалось декодировать книгу из тела запроса
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Невозможно создать: %v\n", err)
		w.Write([]byte("Невозможно создать новый блок"))
		return
	}

	h := md5.New()                                     //Создание объекта хеша
	io.WriteString(h, book.ISBN+book.PublishData)      //Вывод данных книги
	book.ID = fmt.Sprintf("%x", h.Sum(nil))            //Получение идентификатора книги
	response, err := json.MarshalIndent(book, "", " ") //Подготовка ответа в виде сериализации JSON с отступами
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Не удалось распределить полезную нагрузку: %v\n", err)
		w.Write([]byte("Невозможно сохранить данные"))
		return
	}
	w.WriteHeader(http.StatusOK) //Запись статуса в ответ
	w.Write(response)            //Запись ответа
}

// Запись блока
func writeBlock(w http.ResponseWriter, r *http.Request) {
	var checkoutItem BookCheckout //Выделение памяти для проверок книги

	if err := json.NewDecoder(r.Body).Decode(&checkoutItem); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Невозможно записать блок: %v", err)
		w.Write([]byte("Невозможно записать блок"))
		return
	}
}
