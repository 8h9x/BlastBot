package main

import (
	"flag"
	"github.com/8h9x/BlastBot/database/pkg/fortniteapicom"
	"log"
	"runtime"

	"github.com/8h9x/BlastBot/database/internal/cache"
)

func main() {
	var purgeCache bool
	flag.BoolVar(&purgeCache, "purge", false, "Delete existing cache and redownload all images from server")
	flag.Parse()

	log.Println(purgeCache)

	cosmetics, err := fortniteapicom.FetchCosmetics()
	if err != nil {
		log.Fatal(err)
	}

	log.Println(cosmetics)

	switch runtime.GOOS {
    case "linux":
    case "freebsd":
        log.Println("Linux/FreeBSD", cache.CACHE_FOLDER_LINUX)
    case "darwin":
		log.Println("macOS", cache.CACHE_FOLDER_DARWIN)
    case "windows":
		log.Println("Windows", cache.CACHE_FOLDER_WINDOWS)
    default:
		log.Println("unknown OS:", runtime.GOOS)
    }
}