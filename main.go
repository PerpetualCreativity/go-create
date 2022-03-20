package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/PerpetualCreativity/fancyChecks"
	"github.com/schollz/progressbar/v3"
)

var name, file, output *string
var compress, exclude, cgoInc, firstClassInc *bool
var osList []string

var fc = fancyChecks.New("gocreate", "", "Status", "Error")
var pb = progressbar.New(0)
var er []string

func main() {
	listOut, err := exec.Command("go", "tool", "dist", "list", "-json").Output()
	fc.ErrCheck(err, "Failed to get list of Go-supported build targets")

	type OsArch struct {
		GOOS         string
		GOARCH       string
		CgoSupported bool
		FirstClass   bool
	}
	var osarch []OsArch

	err = json.Unmarshal(listOut, &osarch)
	fc.ErrCheck(err, "Failed to parse list of Go-supported build targets")

	wg := new(sync.WaitGroup)
	count := 0
	for _, osa := range osarch {
		inc := *exclude

		if !inc {
			if (*cgoInc && osa.CgoSupported) || (*firstClassInc && osa.FirstClass) {
				inc = true
			} else {
				for _, e := range osList {
					if e == osa.GOOS {
						inc = !*exclude
						break
					}
				}
			}
		}

		if inc {
			wg.Add(1)
			count++
			pb.ChangeMax(count)
			go build(wg, osa.GOOS, osa.GOARCH)
		}
	}

	wg.Wait()

	if _, err := os.Stat(*file); err == nil {
		os.Remove(strings.Replace(*file, ".go", "", 1))
		os.Remove(strings.Replace(*file, ".go", ".exe", 1))
	}

	fmt.Println("")
	if len(er) != 0 {
		for _, e := range er {
			fmt.Println(e)
		}
		fc.Success("Successfully built all other executables.")
	} else {
		fc.Success("Successfully built all executables.")
	}
}

func init() {
	name = flag.String("name", "", "specify name of output")
	file = flag.String("file", "main.go", "specify file to build")
	compress = flag.Bool("compress", false, "compress output")
	output = flag.String("output", "", "directory to place outputs")
	exclude = flag.Bool("exclude", false, "exclude specified OSs instead of including")
	cgoInc = flag.Bool("cgo", false, "build on all platforms that support Cgo")
	firstClassInc = flag.Bool("first-class", false, "build on all first-class platforms")
	flag.Parse()
	osList = flag.Args()

	fc.ErrNComp(*name, "", "Name must be defined.")

	_, err := os.Stat(*file)
	fc.ErrCheck(err, fmt.Sprintf("File %s does not exist.", *file))

	fc.ErrNComp(*output, "", "Output must be defined.")

	if _, err := os.Stat(*output); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(*output, os.ModePerm)
		fc.ErrCheck(err, fmt.Sprintf("Output folder %s does not exist. When trying to create it, the following error(s) occured", *output))
	}
}

func build(wg *sync.WaitGroup, ops string, arch string) {
	defer wg.Done()

	var cmd *exec.Cmd
	if *compress {
		cmd = exec.Command("go", "build", "-ldflags=-s", "-ldflags=-w", *file)
	} else {
		cmd = exec.Command("go", "build", *file)
	}
	cmd.Env = append(
		os.Environ(),
		fmt.Sprintf("GOOS=%s", ops),
		fmt.Sprintf("GOARCH=%s", arch),
	)
	bin, err := cmd.CombinedOutput()
	if err != nil {
		er = append(er, fmt.Sprintf("Failed to build on OS %s, arch %s", ops, arch))
	}

	err = os.WriteFile(fmt.Sprintf("%s/%s-%s-%s", *output, *name, ops, arch), bin, 0666)
	fc.ErrCheck(err, "Could not write output")
	pb.Add(1)
}

