package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"taskTest/internal/config"
	"taskTest/internal/entity"
	"taskTest/internal/storage/psg"
)

func main() {
	//считываем .env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Ошибка загрузки файла .env", err)
	}

	//парсим переменные окружения в conf
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	//подключение к бд
	db, err := psg.NewPostgresStorage(conf)

	if err != nil {
		log.Fatal("Ошибка подключения к бд:", err)
	}

	log.Println("База данных запущена")

	ctx := context.Background()

	user1 := entity.User{
		Name:  "Ivan",
		Email: "ffgg@gsds.ru",
	}

	id, err := db.Create(ctx, &user1)
	if err != nil {
		log.Println(err)
	}

	us, err := db.Read(ctx, id)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(us.Email == user1.Email)

}
