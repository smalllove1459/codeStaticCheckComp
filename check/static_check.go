package check

import (
	"bytes"
	"log"
	"os/exec"
)

//first git clone ur project to ur disk
func gitCloneCode(url string) (string, error) {
	cmd := exec.Command("git", "clone", url)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	log.Println("git clone project completed.......")
	return result.String(), nil
}

//locate ur project folder name
func getGitpath(url string) string {
	urlList := strings.Split(url, "/")
	log.Println(urlList)
	last := urlList[len(urlList)-1:]
	return strings.Split(last[0], ".")[0]
}

//second go build ur project

func listFiles(path string, f os.FileInfo, err error) error {
	var strRet string
	// strRet, _ = os.Getwd()

	osType := os.Getenv("GOOS")
	if osType == "windows" {
		strRet += "\\"
	} else if osType == "linux" {
		strRet += "/"
	}
	strRet += path

	ok := strings.HasSuffix(strRet, ".go")
	if ok {
		// lintFiles = append(lintFiles, strRet)
		lintFiles = append(lintFiles, strRet)
	}

	return nil
}

func getFolderInfo(path string) string {
	err := filepath.Walk(path, listFiles)
	if err != nil {
		log.Println(err.Error())
	}
	return ""
}

func lintCode(fileName string) (string, error) {
	cmd := exec.Command("golint", fileName)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return result.String(), nil
}

func clocAllFiles(dirName string) string {
	cmd := exec.Command("cloc", dirName)
	var result bytes.Buffer
	cmd.Stdout = &result
	err := cmd.Run()
	if err != nil {
		return err.Error()
	}
	return result.String()
}
