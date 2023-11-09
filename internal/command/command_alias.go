package command

import "fmt"

// commandAlias handles listing the available aliases
func (ctx *Context) commandAlias(args []string) error {
	fmt.Println("the available aliases are")
	for key, value := range aliasMap {
		fmt.Printf("%s -> %s\n", key, value)
	}
	return nil
}
