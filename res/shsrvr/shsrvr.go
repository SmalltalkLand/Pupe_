package shsrvr

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func run(cmd *exec.Cmd, pid string) {
	cmd.Run()
	cmd2 := exec.Command("/bin/kill", "-9", pid)
	cmd2.Run()
}
func server() {
	file, err := os.Open(os.Args[2])
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		var line string
		var nline string
		line, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}
		nline, err = reader.ReadString('\n')
		if err != nil && err != io.EOF {
			break
		}

		// Process the line here.
		cmd := exec.Command(os.Args[3], append(os.Args[4:], strings.Split(nline, " ")...)...)
		in, err := os.Open(fmt.Sprintf("/proc/%s/fd/0", line))
		fileOut, err := os.Open(fmt.Sprintf("/proc/%s/fd/1", line))
		fileError, err := os.Open(fmt.Sprintf("/proc/%s/fd/1", line))
		cmd.Stdin = in
		cmd.Stdout = fileOut
		cmd.Stderr = fileError
		go run(cmd, line)
		if err != nil {
			break
		}
	}
}
func client() {
	file, err := os.Open(os.Args[2])
	defer file.Close()
	if err != nil {
		return
	}
	w := bufio.NewWriter(file)
	n, err := w.WriteString(fmt.Sprintf("%d\n%s", os.Getpid(), strings.Join(os.Args[3:], " ")))
	if n == 0 || err != nil {
		return
	}
	for {

	}
}
func main() {
	if os.Args[1] == "server" {
		server()
	}
	if os.Args[1] == "client" {
		client()
	}
}
