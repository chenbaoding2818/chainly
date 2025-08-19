package main

import (
	"github.com/chenbaoding2818/chainly/core"
)

func main() {
	// Initialize the server
	server := core.NewServer()
	// Register default components
	server.RegisterDefaultComponents()
	// Register custom components (注册自定义的组件)
	server.RegisterComponent(nil)
	// Run the server
	server.Run()
}
