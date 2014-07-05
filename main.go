package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

const DIR = "/usr/local/"

var ver = flag.String("v", "", "the go version you want to download")
var url = flag.String("u", "", "the url of the go version you want to download")
var helpMsg = "Usage: go-update-install -v <version-number>\n" +
	"   or: go-update-install -u <url of binary to install"

func buildFilename(url string) string {
	tokens := strings.Split(url, "/")
	filename := DIR + tokens[len(tokens)-1]
	return filename
}

func download(url string) (string, error) {
	filename := buildFilename(url)
	f, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return filename, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return filename, err
}

func uninstallOld() error {
	dir := DIR + "go/"
	fmt.Printf("\nNow removing old installation deleting %v\n", dir)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("\nNo old installation found in %v\n", dir)
		} else {
			return err
		}
	} else {
		err = os.RemoveAll(dir)
		if err != nil {
			return fmt.Errorf("Error while deleting %v : %v", dir, err)
		}
	}
	return nil
}

func installNew(filename string) (string, error) {
	fmt.Printf("\nNow extracting archive %v\n", filename)
	// run the shell command: tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
	cmd := exec.Command("tar", "-C", DIR, "-xzf", filename)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func checkGoVersion() (string, error) {
	fmt.Println("\nNow checking go version\n")
	// run the shell command: go version
	cmd := exec.Command("go", "version")
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func main() {
	flag.Parse()
	if *url == "" {
		if *ver != "" {
			*url = "http://golang.org/dl/go" + *ver + ".linux-amd64.tar.gz"
		} else {
			fmt.Println(helpMsg)
			return
		}
	}
	filename, err := download(*url)
	if filename != "" {
		defer os.Remove(filename)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	err = uninstallOld()
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := installNew(buildFilename(*url))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)

	out, err = checkGoVersion()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(out)

}
