package main

import "fmt"
type coder interface {
	code()
	debug()
}
type Gopher struct {
	language string
}
func (p Gopher) code() {
	fmt.Printf("I am coding %s language\n", p.language)
}
func (p *Gopher) debug() {
	fmt.Printf("I am debuging %s language\n", p.language)
}
func main() {
	var c coder = Gopher{"Go"}
	c.code()
	//c.debug()
}
