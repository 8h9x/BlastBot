package cache

import (
	"os"
	"path"
)

const CACHE_FOLDER_LINUX = "/var/cache/blast"
const CACHE_FOLDER_DARWIN = "/Library/Caches/blast"
var CACHE_FOLDER_WINDOWS = path.Join(os.Getenv("PROGRAMDATA"), "blast", "cache")