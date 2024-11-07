package main

import "log"

func ExitOnError(err error) {
	if err == nil {
		return
	}

	log.Fatalln(err)
}
