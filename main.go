package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"taskTest/internal/config"
	"taskTest/internal/entity"
	"taskTest/internal/storage/psg"
	"taskTest/internal/usecase"
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
	db, err := psg.InitStorage(conf)
	if err != nil {
		log.Fatal("Ошибка подключения к бд:", err)
	}

	ucDB := usecase.NewUseCase(db)
	log.Println("База данных запущена")

	user1 := entity.User{
		ID:    1557,
		Name:  "Ivan",
		Email: "ffgg@gsds.ru",
	}

	err = ucDB.CreateUser(&user1)
	if err != nil {
		log.Println(err)
	}

	us, err := ucDB.GetUserByID(1557)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(us.Email == user1.Email)
	user1.Name = "Xaxapuz"

	err = ucDB.UpdateUser(&user1)
	if err != nil {
		log.Fatal(err)
	}
	us, err = ucDB.GetUserByID(1557)
	fmt.Println(us.Name == user1.Name)

}
