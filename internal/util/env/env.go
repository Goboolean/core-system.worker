package env

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

// Just import this package to get all the env variables at the root of the project
// Import this package anonymously as shown below:
// import _ "github.com/Goboolean/core-system.worker/internal/util/env"

const baseDIR = "core-system.worker"

func init() {
	path, err := filepath.Abs(".")
	if err != nil {
		panic(err)
	}

	for base := filepath.Base(path); base != "core-system.worker" && base != "app"; {
		path = filepath.Dir(path)
		base = filepath.Base(path)

		if base == "." || base == "/" {
			panic(errRootNotFound)
		}
	}

	if err := os.Chdir(path); err != nil {
		panic(err)
	}

	if err := godotenv.Load(); err != nil {
		path, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		panic(fmt.Errorf("%v, working directory: %s", err, path))
	}
}

var errRootNotFound = errors.New("could not find root directory, be sure to set root of the project as fetch-server")
