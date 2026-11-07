package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
)

type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Grade int    `json:"grade"`
}

type Book struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Author string `json:"author"`
}

var (
	students = []Student{
		{ID: 1, Name: "Анна", Grade: 5},
		{ID: 2, Name: "Иван", Grade: 4},
		{ID: 3, Name: "Мария", Grade: 5},
	}
	books = []Book{
		{ID: 1, Title: "Язык программирования Go", Author: "Алан Донован"},
		{ID: 2, Title: "Чистый код", Author: "Роберт Мартин"},
	}
	mu sync.RWMutex
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Запуск: go run main.go <номер задачи 1–12 или hw>")
		return
	}

	mux := http.NewServeMux()
	switch os.Args[1] {
	case "1":
		mux.HandleFunc("GET /students", task1)
	case "2":
		mux.HandleFunc("GET /students", task2WithHeader)
		mux.HandleFunc("GET /students/no-header", task2WithoutHeader)
	case "3":
		mux.HandleFunc("GET /students", task1)
		mux.HandleFunc("GET /books", task3)
	case "4":
		mux.HandleFunc("GET /stats", task4)
	case "5":
		mux.HandleFunc("GET /students/{id}", task5)
	case "6":
		mux.HandleFunc("GET /students/{id}", task6)
	case "7":
		mux.HandleFunc("GET /students/top", task7)
		mux.HandleFunc("GET /students/{id}", task5)
	case "8":
		mux.HandleFunc("GET /students", task8Students)
		mux.HandleFunc("GET /books", task8Books)
		mux.HandleFunc("GET /stats", task8Stats)
		mux.HandleFunc("GET /students/{id}", task8Student)
		mux.HandleFunc("GET /students/top", task8Top)
	case "9":
		mux.HandleFunc("POST /students", task9)
	case "10":
		mux.HandleFunc("POST /students", task10)
	case "11":
		mux.HandleFunc("POST /students", task11)
	case "12":
		mux.HandleFunc("GET /students", task12List)
		mux.HandleFunc("GET /students/{id}", task12Get)
		mux.HandleFunc("POST /students", task12Create)
	case "hw":
		mux.HandleFunc("GET /students", homeworkStudents)
		mux.HandleFunc("GET /students/{id}", task12Get)
		mux.HandleFunc("POST /students", task12Create)
		mux.HandleFunc("GET /students/count", homeworkCount)
		mux.HandleFunc("GET /books", homeworkBooks)
		mux.HandleFunc("GET /books/{id}", homeworkBook)
		mux.HandleFunc("POST /books", homeworkCreateBook)
	default:
		fmt.Println("Допустимые задачи: 1–12, hw")
		return
	}

	log.Printf("Задача %s. Сервер: http://localhost:8080", os.Args[1])
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// Задание 1

func task1(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// Задание 2

func task2WithHeader(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

func task2WithoutHeader(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(students)
}

// Задание 3
func task3(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// Задание 4

func task4(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"count": len(students)})
}

// Задание 5

func task5(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	for _, s := range students {
		if s.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	http.Error(w, "не найдено", http.StatusNotFound)
}

// Задание 6

func task6(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	for _, s := range students {
		if s.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(s)
			return
		}
	}
	http.Error(w, "не найдено", http.StatusNotFound)
}

// Задание 7

func task7(w http.ResponseWriter, r *http.Request) {
	result := make([]Student, 0)
	for _, s := range students {
		if s.Grade == 5 {
			result = append(result, s)
		}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// Задание 8

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("Ошибка записи JSON: %v", err)
	}
}

func task8Students(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, students)
}

func task8Books(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, books)
}

func task8Stats(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, map[string]int{"count": len(students)})
}

func task8Student(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	for _, s := range students {
		if s.ID == id {
			writeJSON(w, s)
			return
		}
	}
	http.Error(w, "не найдено", http.StatusNotFound)
}

func task8Top(w http.ResponseWriter, r *http.Request) {
	result := make([]Student, 0)
	for _, s := range students {
		if s.Grade == 5 {
			result = append(result, s)
		}
	}
	writeJSON(w, result)
}

// Задание 9

func task9(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	mu.Lock()
	s.ID = len(students) + 1
	students = append(students, s)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, s)
}

// Задание 10

func task10(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	mu.Lock()
	s.ID = len(students) + 1
	students = append(students, s)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, s)
}

// Задание 11

func task11(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(s.Name) == "" {
		http.Error(w, "имя обязательно", http.StatusBadRequest)
		return
	}
	mu.Lock()
	s.ID = len(students) + 1
	students = append(students, s)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, s)
}

// Задание 12

func task12List(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	result := append([]Student{}, students...)
	mu.RUnlock()
	writeJSON(w, result)
}

func task12Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	mu.RLock()
	for _, s := range students {
		if s.ID == id {
			mu.RUnlock()
			writeJSON(w, s)
			return
		}
	}
	mu.RUnlock()
	http.Error(w, "не найдено", http.StatusNotFound)
}

// ДЗ

// Пункт 1
func homeworkBooks(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	result := append([]Book{}, books...)
	mu.RUnlock()
	writeJSON(w, result)
}

func homeworkBook(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "некорректный id", http.StatusBadRequest)
		return
	}
	mu.RLock()
	for _, b := range books {
		if b.ID == id {
			mu.RUnlock()
			writeJSON(w, b)
			return
		}
	}
	mu.RUnlock()
	http.Error(w, "не найдено", http.StatusNotFound)
}

// Пункт 2
func homeworkCreateBook(w http.ResponseWriter, r *http.Request) {
	var b Book
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(b.Title) == "" {
		http.Error(w, "заголовок обязателен", http.StatusBadRequest)
		return
	}
	mu.Lock()
	b.ID = len(books) + 1
	books = append(books, b)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, b)
}

// Пункт 3
func homeworkCount(w http.ResponseWriter, r *http.Request) {
	mu.RLock()
	count := len(students)
	mu.RUnlock()
	writeJSON(w, map[string]int{"count": count})
}

// Пункт 4
func homeworkStudents(w http.ResponseWriter, r *http.Request) {
	var grade int
	filter := r.URL.Query().Has("grade")
	if filter {
		var err error
		grade, err = strconv.Atoi(r.URL.Query().Get("grade"))
		if err != nil {
			http.Error(w, "оценка должна быть целым числом", http.StatusBadRequest)
			return
		}
	}
	mu.RLock()
	result := make([]Student, 0)
	for _, s := range students {
		if !filter || s.Grade == grade {
			result = append(result, s)
		}
	}
	mu.RUnlock()
	writeJSON(w, result)
}

func task12Create(w http.ResponseWriter, r *http.Request) {
	var s Student
	if err := json.NewDecoder(r.Body).Decode(&s); err != nil {
		http.Error(w, "некорректный JSON", http.StatusBadRequest)
		return
	}
	if strings.TrimSpace(s.Name) == "" {
		http.Error(w, "имя обязательно", http.StatusBadRequest)
		return
	}
	mu.Lock()
	s.ID = len(students) + 1
	students = append(students, s)
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, s)
}
