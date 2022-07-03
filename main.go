package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"tcp/utils"
)

const pathConfig = "configs/config.json"

func main() {
	// создаём конфиг
	config := utils.NewConfig()
	// читаем файл 
	byteValue, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		log.Fatal("failed to read from config.json")
	}
	// записываем в нашу структуру данные с файла
	err = json.Unmarshal(byteValue, &config)
	// создаём обьект нашей основной структуры
	srv := &utils.Server{}
	// получаем порты вёденные пользователем из терминала если не ввёл то оставляем стандарьные
	port, err := GetPort(config)
	// настраиваем наш сервер введёнными значениями
	err = srv.ConstructorSrv(port, config.MaxConnections)

	if err != nil {
		log.Printf("ERROR -> main: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Listening on the port %v\n", port)
	for {
		// получаем коннект из функции
		conn, err := srv.Server.Accept()
		if err != nil {
			log.Fatal(err)
			break
		}
		// запускаем горутину которая добавляет нового пользователя в чат 
		go srv.ConnectMessenger(conn)
	}
	

}

// GetPort функция которая берёт значение с  консоли 
func GetPort(Port *utils.Config) (string, error) {
	args := os.Args
	if len(args) < 2 {
		return Port.BindAddr, nil
	}else if len(args) > 2 {
		return "", errors.New("User inputs more than 2 arguments")
	}
	return ":" + os.Args[1], nil
}