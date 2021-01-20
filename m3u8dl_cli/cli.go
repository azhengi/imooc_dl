package m3u8dl_cli

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
)

const CLI_PATH = "./N_m3u8DL-CLI_v2.9.1.exe"

func checkCLIExists() (bool, error) {
	_, err := os.Stat(CLI_PATH)
	if err == nil {
		return true, nil
	}
	return false, err
}

func Run(filename string) {
	exist, _ := checkCLIExists()
	if exist {
		line := CLI_PATH + " " + filename
		var stdoutBuf, stderrBuf bytes.Buffer
		cmd := exec.Command("cmd", "/c", "start "+line)
		// cmd := exec.Command("cmd")

		// stdin, _ := cmd.StdinPipe()
		stdoutIn, _ := cmd.StdoutPipe()
		stderrIn, _ := cmd.StderrPipe()

		var errStdout, errStderr error
		stdout := io.MultiWriter(os.Stdout, &stdoutBuf)
		stderr := io.MultiWriter(os.Stderr, &stderrBuf)

		err := cmd.Start()
		if err != nil {
			log.Fatalf("cmd.Start() failed with '%s'\n", err)
		}

		go func() {
			_, errStdout = io.Copy(stdout, stdoutIn)
		}()
		go func() {
			_, errStderr = io.Copy(stderr, stderrIn)
		}()

		// io.WriteString(stdin, "index.m3u8")
		// stdin.Close()

		defer func() {
			err = cmd.Wait()
			if err != nil {
				log.Fatalf("cmd.Run() failed with %s\n", err)
			}
		}()

		if errStdout != nil || errStderr != nil {
			log.Fatal("failed to capture stdout or stderr\n")
		}

		outStr, errStr := string(stdoutBuf.Bytes()), string(stderrBuf.Bytes())
		// fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
		fmt.Println("\n========================= print =========================")
		fmt.Printf("standard output:\n%s", outStr)
		fmt.Println()
		fmt.Printf("standard error:\n%s", errStr)
	}

}
