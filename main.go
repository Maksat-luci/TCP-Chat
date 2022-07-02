package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"tcp/utils"
)

const pathConfig = "configs/config.json"

func main() {
	config := utils.NewConfig()
	byteValue, err := ioutil.ReadFile(pathConfig)
	if err != nil {
		log.Fatal("failed to read from config.json")
	}
	
	err = json.Unmarshal(byteValue, &config)
	fmt.Println(config)
}
