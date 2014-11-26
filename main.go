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
var helpMsg = "Usage: go-update -v <version-number>\n" +
	"   or: go-update -u <url of binary to install>"

// download downloads the url to saveDir, overwriting the file if present
func download(url string, saveDir string) (string, error) {
	tokens := strings.Split(url, "/")
	filepath := saveDir + tokens[len(tokens)-1]
	f, err := os.Create(filepath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	resp, err := http.Get(url)
	if err != nil {
		return filepath, err
	}
	defer resp.Body.Close()
	_, err = io.Copy(f, resp.Body)
	return filepath, err
}

// uninstallOld remove the old Go installation in dir
func uninstallOld(dir string) error {
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

// installNew install Go from filepath to installDir
func installNew(filepath, installDir string) (string, error) {
	fmt.Printf("\nNow extracting archive %v\n", filepath)
	// run the shell command: tar -C /usr/local -xzf go$VERSION.$OS-$ARCH.tar.gz
	cmd := exec.Command("tar", "-C", installDir, "-xzf", filepath)
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

func checkGoVersion() (string, error) {
	fmt.Println("\nNow checking go version\n")
	// run the shell command: go version
	cmd := exec.Command("/usr/local/go/bin/go", "version")
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
	filepath, err := download(*url, DIR)
	if filepath != "" {
		defer os.Remove(filepath)
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	err = uninstallOld(DIR + "go/")
	if err != nil {
		fmt.Println(err)
		return
	}
	out, err := installNew(filepath, DIR)
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
