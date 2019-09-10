package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	testfile        = "QmZHeNJTU4jFzgBAouHSqbT2tyYJxgk6i15e7x5pudBune"
	testfilecontent = "Hello BTFS!"

	routeError = "/error"
)

type errorData struct {
	HVal        string `json:"h_val"`
	PeerId      string `json:"peer_id"`
	ErrorStatus string `json:"error_status"`
}

// we need to delete the file for get test from last run
func prepare_test(btfsBinaryPath, statusServerDomain, peerId string) bool {
	cmd := exec.Command(btfsBinaryPath, "rm", testfile)
	err := cmd.Start()

	if err != nil {
		fmt.Printf("btfs rm failed with message: [%v]\n", err)
		errMsg := fmt.Sprintf("btfs rm failed with message: [%v]", err)
		log.Info(errMsg)
		send_error(errMsg, statusServerDomain, peerId)
		return false
	} else {
		log.Info("btfs test preparation succeed\n")
	}
	return true
}

func get_functest(btfsBinaryPath string) error {
	// btfs get file saved to current working directory
	dir, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("get working directory failed: [%v]", err))
	}

	cmd := exec.Command(btfsBinaryPath, "get", "-o", dir, testfile)
	out, err := cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs get test failed: [%v], Out[%s]", err, string(out)))
	}

	data, err := ioutil.ReadFile(dir + "/" + testfile)
	if err != nil {
		fmt.Sprintf("btfs get test: read file failed: [%v]\n", err)
		return errors.New(fmt.Sprintf("btfs get test: read file failed: [%v]", err))
	}

	// remote last "\n" before compare
	if string(data[:len(data)-1]) != testfilecontent {
		fmt.Sprintf("btfs get test: get different content[%s]\n", string(data))
		return errors.New(fmt.Sprintf("btfs get test: get different content[%s]", string(data)))
	}

	return nil
}

func add_functest(btfsBinaryPath string) error {
	// write btfs id command output to a file in current working directory
	// then btfs add that file for test
	dir, err := os.Getwd()
	if err != nil {
		return errors.New(fmt.Sprintf("get working directory failed: [%v]", err))
	}

	cmd := exec.Command(btfsBinaryPath, "id")
	out, err := cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test: btfs id failed: [%v], Out[%s]", err, string(out)))
	}

	// add current time stamp to file content so every time adding-file hash is different
	currentTime := time.Now().String()
	out = append(out, currentTime...)

	origin := out
	filename := dir + "/btfstest.txt"
	err = ioutil.WriteFile(filename, out, 0644)
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test: write file failed: [%v]", err))
	}

	cmd = exec.Command(btfsBinaryPath, "add", filename)
	out, err = cmd.Output()
	if err != nil {
		return errors.New(fmt.Sprintf("btfs add test failed: [%v]", err))
	}

	s := strings.Split(string(out), " ")
	if len(s) < 2 {
		return errors.New(fmt.Sprintf("btfs add test failed: invalid add result[%s]", string(out)))
	}

	addfilehash := s[1]
	cmd = exec.Command(btfsBinaryPath, "cat", addfilehash)
	out, err = cmd.Output()

	if string(out) != string(origin) {
		return errors.New(fmt.Sprintf("btfs add test failed: cat different content, btfs add file:[%s], btfs cat file:[%s]",
			string(origin), string(out)))
	}

	return nil
}

// function to send error message to status server
func send_error(errMsg, statusServerDomain, peerId string) {
	errData := new(errorData)
	errData.ErrorStatus = errMsg
	errData.PeerId = peerId
	errDataMarshaled, err := json.Marshal(errData)

	// reports to status server by making HTTP request
	req, err := http.NewRequest("POST", fmt.Sprintf("%s%s", statusServerDomain, routeError), bytes.NewReader(errDataMarshaled))
	if err != nil {
		log.Info(fmt.Sprintf("failed to make new http request: %s", err.Error()))
		return
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Info(fmt.Sprintf("failed to perform http.DefaultClient.Do(): %s", err.Error()))
		return
	}
	defer res.Body.Close()
}