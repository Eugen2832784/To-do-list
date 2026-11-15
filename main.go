package main

import (
	"fmt"
	"os"

	"todo/internal/cli"
	"todo/internal/storage"
)



func main() {
	app := cli.New(storage.New("tasks.json"), os.Stdout)

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}
