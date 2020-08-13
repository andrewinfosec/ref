package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

const HELP = `
ref -- Manage references for large writing projects.

ref add       - create reference number, copy to clipboard, and open corresponding directory
ref <number>  - open file associated with reference <number>
ref loc       - set location of database to current directory
`

func checkErr(err error) {
	// {{{
	if err != nil {
		panic(err)
	}
	// }}}
}

func cmdAdd() {
	// {{{
	// identify location of db
	user, err := user.Current()
	checkErr(err)

	if _, err := os.Stat(user.HomeDir + "/.ref"); os.IsNotExist(err) {
		fmt.Println("Error: ~/.ref does not exist. Use REF LOC to create." + "\n\n" + HELP)
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(user.HomeDir + "/.ref")
	checkErr(err)
	dbpath := strings.TrimSpace(string(data))

	// identify reference number
	var dirNames []string
	perFile := func(fullPath string, f os.FileInfo, err error) error {
		if fullPath != dbpath && f.IsDir() {
			dirNames = append(dirNames, f.Name())
		}
		return nil
	}
	err2 := filepath.Walk(dbpath, perFile)
	checkErr(err2)

	var highest int = 0
	for _, v := range dirNames {
		i, err := strconv.Atoi(v)
		checkErr(err)
		if i > highest {
			highest = i
		}
	}
	number := strconv.Itoa(highest + 1)
	fmt.Println("Created reference " + number + ".")

	// copy number to clipboard
	bash, err := exec.LookPath("bash")
	checkErr(err)

	echo, err := exec.LookPath("echo")
	checkErr(err)

	pbcopy, err := exec.LookPath("pbcopy")
	checkErr(err)

	_, err3 := exec.Command(bash, "-c", echo+" -n "+number+"| "+pbcopy).Output()
	checkErr(err3)

	// create directory
	err4 := os.Mkdir(dbpath+"/"+number, 0600)
	checkErr(err4)

	// open directory
	binary, err := exec.LookPath("open")
	checkErr(err)
	err5 := syscall.Exec(binary, []string{"open", dbpath + "/" + number}, os.Environ())
	checkErr(err5)

	os.Exit(0)
	// }}}
}

func cmdLoc() {
	// {{{
	loc, err := os.Getwd()
	checkErr(err)

	user, err := user.Current()
	checkErr(err)

	err2 := ioutil.WriteFile(user.HomeDir+"/.ref", []byte(loc+"\n"), 0600)
	checkErr(err2)

	fmt.Println("Location of database set to " + loc + " in " + user.HomeDir + "/.ref")

	os.Exit(0)
	// }}}
}

func main() {
	// {{{
	for _, v := range os.Args[1:] {
		switch v {
		case "add":
			cmdAdd()
		case "loc":
			cmdLoc()
		default:
			// XXX

			fmt.Println("Unknown argument: " + v + "\n\n" + HELP)
			os.Exit(0)
		}
	}

	fmt.Println(HELP)
	// }}}
}
