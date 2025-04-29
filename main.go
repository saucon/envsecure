/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/saucon/envsecure/cmd"
	"log"
	"os"
)

func main() {
	err := cmd.Execute()
	if err != nil {
		log.Printf("log error: %s", err.Error())
		os.Exit(1)
		return
	}
}
