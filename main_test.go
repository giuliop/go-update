package main

import (
	"os"
	"os/exec"
	"testing"
)

const testUrl = "http://pendlay.3dcartstores.com/assets/images/1web/pendlaylogo.jpg"
const filename = "pendlaylogo.jpg"

func TestDownload(t *testing.T) {
	saveDir := "./"
	f, err := download(testUrl, saveDir)
	defer func() {
		os.Remove(saveDir + filename)
	}()
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(f)
	if err != nil {
		t.Fatal(err)
	}
	if f != saveDir+filename {
		t.Fatalf("Filepath is %v, expecting %v", f, saveDir+filename)
	}
}

func TestUninstallOld(t *testing.T) {
	// run the shell command: mkdir -p ./one/two ; touch ./one/two/test
	cmd := exec.Command("mkdir", "-p", "./one/two")
	cmd2 := exec.Command("touch", "./one/two/test")
	_, err := cmd.Output()
	_, err2 := cmd2.Output()
	_, err3 := os.Stat("./one")
	if err != nil || err2 != nil || err3 != nil {
		t.Fatal(err, err2, err3)
	}
	err = uninstallOld("./one")
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat("./one")
	if !os.IsNotExist(err) {
		t.Fatal("Directory ./one should have been deleted")
	}
}
