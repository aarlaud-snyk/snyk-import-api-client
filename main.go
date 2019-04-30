package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

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

	bodyJSON, _ := json.MarshalIndent(body, "", " ")
	url := fmt.Sprintf("https://snyk.io/api/v1/org/%s/integrations/%s/import", orgId, integrationId)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(bodyJSON))

	req.Header.Add("Content-Type", "application/json; charset=utf-8")
	tokenString := fmt.Sprintf("token %s", token)
	req.Header.Add("Authorization", tokenString)

	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Errored when sending request to the server")
		return
	}

	defer resp.Body.Close()
	respBody, _ := ioutil.ReadAll(resp.Body)

	fmt.Println(resp.Status)
	fmt.Println(string(respBody))
}

func main() {
	tokenPtr := flag.String("token", "", "your Snyk token")
	orgIdPtr := flag.String("orgId", "", "your Organization ID")
	integrationIdPtr := flag.String("intId", "", "your Integration ID")
	projectKeyPtr := flag.String("projectkey", "", "your Integration ID")
	repoNamePtr := flag.String("repo", "", "your Integration ID")
	repoSlugPtr := flag.String("reposlug", "", "your Integration ID")

	flag.Parse()

	orgId := *orgIdPtr
	integrationId := *integrationIdPtr
	token := *tokenPtr
	projectKey := *projectKeyPtr
	repoName := *repoNamePtr
	repoSlug := *repoSlugPtr

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}
	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("Usage: find <repoFolder> -name pom.xml | ./snykimport")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	var paths []string

	for {
		input, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			break
		}
		paths = append(paths, strings.TrimSuffix(input, "\n"))

	}
	importProject(token, orgId, integrationId, projectKey, repoName, repoSlug, paths)

}
