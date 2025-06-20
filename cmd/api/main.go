package main

import "github.com/Ablebil/sea-catering-be/internal/bootstrap"

func main() {
	if err := bootstrap.Start(); err != nil {
		panic(err)
	}
}
