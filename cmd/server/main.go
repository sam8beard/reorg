package main

import (
	"embed"
	vite "github.com/torenware/vite-go"
	"log"
	"os"
)

// For development
const vueDir = "../../frontend"
const entryPoint = "../../frontend/src/main.js"

var frontendFS = os.DirFS(vueDir)

// For production
//
//go:embed "dist"
var dist embed.FS

func main() {

	config := vite.ViteConfig{
		Environment: "development",
		AssetsPath:  vueDir,
		EntryPoint:  entryPoint,
		Platform:    "",
		FS:          frontendFS,
	}

	// initialize vite library
	glue, err := vite.NewVueGlue(&config)
	if err != nil {
		log.Panicf("could not initialize vite library: %v", err)
	}
}
