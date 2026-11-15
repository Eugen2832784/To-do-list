package main

import (
	"fmt"
	"os"

	"todo/internal/cli"
	"todo/internal/storage"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}
}

func run() error {
	path, err := storage.DefaultPath()
	if err != nil {
		return err
	}

	app := cli.New(storage.New(path), os.Stdout)
	return app.Run(os.Args[1:])
}
