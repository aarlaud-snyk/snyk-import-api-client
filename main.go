package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

var (
	targetFolder string
	targetFile   string
	searchResult []string
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

var manifestFileTypePatterns = [15]string{
	"Gopkg.*",
	"vendor/vendor.json",
	"build.gradle",
	"pom.xml",
	"build.sbt",
	"requirements.txt",
	"Gemfile.lock",
	"package.json",
	"package-lock.json",
	"yarn.lock",
	"project.json",
	"*.csproj",
	"*.vbproj",
	"*.fsproj",
	"packages.config"}

var ignorePathPatterns []string
var repoBasePath string

func Init(
	traceHandle io.Writer,
	infoHandle io.Writer,
	warningHandle io.Writer,
	errorHandle io.Writer) {

	Trace = log.New(traceHandle,
		"",
		0)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func importProject(token string, orgId string, integrationId string, projectKey string, repoName string, repoSlug string, paths []string) {

	type File struct {
		Path string `json:"path"`
	}

	type Target struct {
		ProjectKey string `json:"projectKey"`
		Name       string `json:"name"`
		Branch     string `json:"branch"`
	}

	type BodyStruct struct {
		Target `json:"target"`
		Files  []File `json:"files"`
	}

	client := &http.Client{}

	files := make([]File, 0)
	for _, file := range paths {
		files = append(files, File{Path: file})
	}

	body := BodyStruct{
		Target: Target{
			ProjectKey: projectKey,
			Name:       repoName,
			Branch:     repoSlug,
		},
		Files: files,
	}
	Trace.Println("Programmatically Importing into Snyk")

	bodyJSON, _ := json.MarshalIndent(body, "", " ")
	url := fmt.Sprintf("https://snyk.io/api/v1/org/%s/integrations/%s/import", orgId, integrationId)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	tokenString := fmt.Sprintf("token %s", token)
	req.Header.Add("Authorization", tokenString)

	resp, err := client.Do(req)

	if err != nil {
		Error.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	Info.Println(resp.Status)
	Info.Println(string(respBody))
}

func findFile(path string, fileInfo os.FileInfo, err error) error {

	if err != nil {
		Error.Println(err)
		return nil
	}

	// get absolute path of the folder that we are searching
	absolute, err := filepath.Abs(path)

	if err != nil {
		Error.Println(err)
		return nil
	}

	if fileInfo.IsDir() {
		//fmt.Println("Searching directory ... ", absolute)

		// correct permission to scan folder?
		testDir, err := os.Open(absolute)

		if err != nil {
			if os.IsPermission(err) {
				Warning.Println("No permission to scan ... ", absolute)
				Error.Println(err)
			}
		}
		testDir.Close()
		return nil
	} else {

		// ok, we are dealing with a file
		// is this the target file?

		// yes, need to support wildcard search as well
		// https://www.socketloop.com/tutorials/golang-match-strings-by-wildcard-patterns-with-filepath-match-function
		//matched, err := filepath.Match("main.go", fileInfo.Name())
		for _, pattern := range manifestFileTypePatterns {
			matched, err := filepath.Match(pattern, fileInfo.Name())
			if err != nil {
				Error.Println(err)
			}

			if matched {
				ignored := false
				for _, pathPattern := range ignorePathPatterns {
					if strings.Contains(absolute, pathPattern) {
						ignored = true
						Info.Printf("Ignoring %s\n", absolute)
					}
				}
				if !ignored {
					// yes, add into our search result

					add := "/" + strings.Replace(absolute, repoBasePath, "", 1)
					searchResult = append(searchResult, add)
				}

			}
		}

	}

	return nil
}

func main() {

	rootPathPtr := flag.String("path", "", "your repo root path")
	ignorePathPtr := flag.String("excludeFile", "", "your ignore file path")
	tokenPtr := flag.String("token", "", "your Snyk token")
	orgIdPtr := flag.String("orgId", "", "your Organization ID")
	integrationIdPtr := flag.String("intId", "", "your Integration ID")
	projectKeyPtr := flag.String("projectkey", "", "your Integration ID")
	repoNamePtr := flag.String("repo", "", "your Integration ID")
	repoSlugPtr := flag.String("reposlug", "", "your Integration ID")
	debugPtr := flag.Bool("d", false, "Debug")

	flag.Parse()

	targetFolder := *rootPathPtr
	ignorePath := *ignorePathPtr
	orgId := *orgIdPtr
	integrationId := *integrationIdPtr
	token := *tokenPtr
	projectKey := *projectKeyPtr
	repoName := *repoNamePtr
	repoSlug := *repoSlugPtr
	repoBasePath = targetFolder
	debug := *debugPtr

	if debug {
		Init(os.Stdout, os.Stdout, os.Stdout, os.Stderr)
	} else {
		Init(os.Stdout, ioutil.Discard, ioutil.Discard, os.Stderr)
	}

	var paths []string
	if ignorePath != "" {
		file, err := os.Open(ignorePath)
		if err != nil {
			Error.Println(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {

			ignorePathPatterns = append(ignorePathPatterns, scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			Error.Println(err)
		}
		Info.Println(ignorePathPatterns)
	}
	if targetFolder != "" {
		// sanity check
		testFile, err := os.Open(targetFolder)
		if err != nil {
			Error.Println(err)
			os.Exit(-1)
		}
		defer testFile.Close()

		testFileInfo, _ := testFile.Stat()
		if !testFileInfo.IsDir() {
			Info.Println(targetFolder, " is not a directory!")
			os.Exit(-1)
		}

		err = filepath.Walk(targetFolder, findFile)

		if err != nil {
			Error.Println(err)
			os.Exit(-1)
		}

		// display our search result

		Trace.Println("\n\nFound ", len(searchResult), " hits!")
		Info.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")

		for _, v := range searchResult {
			Info.Println(v)
		}
		Info.Println("@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	} else {
		reader := bufio.NewReader(os.Stdin)

		for {
			input, err := reader.ReadString('\n')
			if err != nil && err == io.EOF {
				break
			}
			paths = append(paths, strings.TrimSuffix(input, "\n"))

		}
	}

	importProject(token, orgId, integrationId, projectKey, repoName, repoSlug, paths)

}
