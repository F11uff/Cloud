package main

import (
	"cloud/Balancer/config"
	"fmt"
	"log"
	"os"
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.LstdFlags | log.Lshortfile) // Логирование в удобном формате(Добавляет в логирование имя файла и время в формате YYYY:MM:DD HH:MM:SS)

	cfg, err := config.InitConfig()

	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	fmt.Printf("%+v\n", cfg)

}
