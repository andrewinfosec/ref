package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
)

const HELP = `
ref -- Manage references for large writing projects.

ref add       - create reference number, copy to clipboard, and open corresponding directory
ref <number>  - open files associated with reference <number>
ref loc       - set location of database to current directory (stored in ~/.ref)
`

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func dbPath() string {
	user, err := user.Current()
	checkErr(err)

	if _, err := os.Stat(user.HomeDir + "/.ref"); os.IsNotExist(err) {
		fmt.Println("~/.ref does not exist. Use `ref loc` to create." + "\n\n" + HELP)
		os.Exit(0)
	}

	data, err := ioutil.ReadFile(user.HomeDir + "/.ref")
	checkErr(err)
	return strings.TrimSpace(string(data))
}

func cmdAdd() {
	// identify reference number
	db := dbPath()

	files, err := ioutil.ReadDir(db)
	checkErr(err)

	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}

	var highest int = 0
	for _, v := range fileNames {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
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

	_, err = exec.Command(bash, "-c", echo+" -n "+number+"| "+pbcopy).Output()
	checkErr(err)

	// create directory
	err = os.Mkdir(db+"/"+number, 0700)
	checkErr(err)

	// open directory
	binary, err := exec.LookPath("open")
	checkErr(err)
	err = syscall.Exec(binary, []string{"open", db + "/" + number}, os.Environ())
	checkErr(err)
}

func cmdLoc() {
	loc, err := os.Getwd()
	checkErr(err)

	user, err := user.Current()
	checkErr(err)

	err = ioutil.WriteFile(user.HomeDir+"/.ref", []byte(loc+"\n"), 0600)
	checkErr(err)

	fmt.Println("Location of database set to " + loc + " in " + user.HomeDir + "/.ref")
}

func openRef(n string) {
	db := filepath.Clean(dbPath())

	if _, err := os.Stat(db + "/" + n); os.IsNotExist(err) {
		fmt.Println("No such reference number:", n)
		os.Exit(0)
	}

	bash, err := exec.LookPath("bash")
	checkErr(err)

	open, err := exec.LookPath("open")
	checkErr(err)

	exec.Command(bash, "-c", open+" "+"\""+db+"\""+"/"+n+"/*.html").Output()
	exec.Command(bash, "-c", open+" "+"\""+db+"\""+"/"+n+"/*.pdf").Output()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(HELP)
		os.Exit(0)
	}

	switch os.Args[1] {
	case "add":
		cmdAdd()
	case "loc":
		cmdLoc()
	default:
		number := regexp.MustCompile(`\d+$`).FindString(os.Args[1])
		if number == "" {
			fmt.Println("Invalid argument: " + os.Args[1] + "\n" + HELP)
			os.Exit(1)
		}
		openRef(number)
	}
}
