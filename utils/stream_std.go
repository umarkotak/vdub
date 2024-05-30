package utils

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

func StreamStd(std io.ReadCloser) {
	scanner := bufio.NewScanner(std)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func StreamCmd(stdout, stderr io.ReadCloser) {
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
}

func StreamCmdTranscript(stdout, stderr io.ReadCloser) {
	scanner := bufio.NewScanner(io.MultiReader(stdout, stderr))
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()

		prefix := ""
		if strings.Contains(m, "[") {
			prefix = "\n"
		} else {
			prefix = " "
		}

		fmt.Printf("%s%s", prefix, m)
	}
}
