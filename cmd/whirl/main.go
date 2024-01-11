package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"

	"github.com/whirl-lang/whirl/pkg/pipeline"
)

type Args struct {
	args []string
}

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("Usage: whirl <filename> [args]")
		return
	}

	dir, err := os.Getwd()

	if err != nil {
		panic(err)
	}

	file, err := os.Create("out.c")

	if err != nil {
		panic(err)
	}

	path := path.Join(dir, args[1])

	ParseFile(path, file)

	out := ExecuteFile(path)
	ParseArgs(Args{args: args}, file)

	fmt.Println(string(out))
}

func ParseArgs(args Args, file *os.File) {
	//TODO: Add proper arg parsing and flags
	deletingFile := true

	if len(args.args) > 2 {
		for _, arg := range args.args[2:] {
			switch arg {
			case "-c":
				deletingFile = false
			}
		}
	}

	file.Close()

	if deletingFile {
		err := os.Remove("out.c")

		if err != nil {
			panic(err)
		}
	}
}

func ExecuteFile(path string) []byte {
	out, _ := exec.Command("tcc", "-run", "out.c").CombinedOutput()

	return out
}

func ParseFile(filename string, out io.Writer) {
	file, err := os.Open(filename)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	stat, err := file.Stat()

	if err != nil {
		panic(err)
	}

	bytes := make([]byte, stat.Size())
	file.Read(bytes)

	pipeline.TranspileC(bytes, path.Dir(filename), out)
}
