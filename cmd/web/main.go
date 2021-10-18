package main

import (
	"database/sql"
	"flag"
	"html/template" // Новый импорт
	"log"
	"net/http"
	"os"

	"SewingWorkshop/pkg/models/mysql"

	_ "github.com/lib/pq"
)

// Добавляем поле templateCache в структуру зависимостей. Это позволит
// получить доступ к кэшу во всех обработчиках.
type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	products      *mysql.ProductModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":8080", "Сетевой адрес веб-сервера")
	dsn := flag.String("dsn", "host=localhost port=5432 user=evgenijvolkov password=25459198 dbname=ksp sslmode=disable", "Название MySQL источника данных")
	flag.Parse()

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	db, err := openDB(*dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	// Инициализируем новый кэш шаблона...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// И добавляем его в зависимостях нашего
	// веб-приложения.
	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		products:      &mysql.ProductModel{DB: db},
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Запуск сервера на http://127.0.0.1%s", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
