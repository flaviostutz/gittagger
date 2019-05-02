package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
)

const VERSION = "0.0.1"

func main() {
	logLevel := flag.String("loglevel", "debug", "debug, info, warning, error")
	gitRepoURL := flag.String("git-repo-url", "", "Git repository URL that will be tagged. ex.: https://github.com/flaviostutz/gittagger-test.git")
	gitUserName0 := flag.String("git-username", "", "Git username")
	gitEmail0 := flag.String("git-email", "", "Git email")
	flag.Parse()

	switch *logLevel {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
		break
	case "warning":
		logrus.SetLevel(logrus.WarnLevel)
		break
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
		break
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Infof("====Starting GitTagger %s====", VERSION)

	gitUsername := *gitUserName0
	gitEmail := *gitEmail0
	// if gitUsername == "" {
	// 	panic("'git-username' is mandatory")
	// }
	if gitEmail == "" {
		panic("'git-email' is mandatory")
	}

	logrus.Infof("Git repo: %s", *gitRepoURL)
	logrus.Infof("Git username: %s", gitUsername)
	logrus.Infof("Git email: %s", gitEmail)

	logrus.Infof("Initializing git workspace")
	ExecShellf("rm -rf /opt/repo")
	ExecShellf("mkdir -p /opt/repo")
	ExecShellf("cd /opt && git clone %s repo", *gitRepoURL)

	_, err := ExecShellf("cd /opt/repo && git config --global user.email %s", gitEmail)
	if err != nil {
		panic(err)
	}

	logrus.Infof("Listening on port 50000")

	router := mux.NewRouter()
	router.HandleFunc("/tag/{name}", handlerTag).Methods("POST")
	router.HandleFunc("/files/{name}", handlerFile).Methods("POST")
	err = http.ListenAndServe("0.0.0.0:50000", router)
	if err != nil {
		logrus.Errorf("Error while listening requests: %s", err)
		os.Exit(1)
	}
}

func handlerTag(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tagname := vars["name"]

	logrus.Debugf("Adding tag %s", tagname)
	out, err := ExecShellf("cd /opt/repo && git tag %s", tagname)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, out)
		return
	}

	logrus.Debugf("Adding tag %s", tagname)
	out, err = ExecShellf("cd /opt/repo && git push origin %s", tagname)
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, out)
		return
	}

	writeResponse(w, http.StatusCreated, fmt.Sprintf("Tag %s pushed successfully to git repository", tagname))
}

func handlerFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filename := vars["name"]
	logrus.Debugf("Updating file contents. name=%s", filename)

	//read body contents
	body, err := ioutil.ReadAll(r.Body)
	b1 := []byte(body)

	//update file contents
	err = ioutil.WriteFile(fmt.Sprintf("/opt/repo/%s", filename), b1, 0644)
	if err != nil {
		logrus.Errorf("failed writing to file: %s", err)
		writeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	out, err := ExecShellf("cd /opt/repo && git add . && git commit -m \"updating file contents using gittagger\" && git push origin master")
	if err != nil {
		writeResponse(w, http.StatusInternalServerError, out)
		return
	}

	writeResponse(w, http.StatusCreated, fmt.Sprintf("File '%s' updated and pushed to git repo successfully", filename))
}

func writeResponse(w http.ResponseWriter, statusCode int, message string) {
	msg := make(map[string]string)
	msg["message"] = message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(msg)
}
