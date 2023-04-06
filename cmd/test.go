package main

import "log"

func main() {
	s := "Hello"
	s_link := &s
	s_double := *s_link
	log.Printf("s_double %s", s_double)

}
