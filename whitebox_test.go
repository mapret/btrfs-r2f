package main

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

/*const testFolder = "TestFilesystemRoot"

func cleanTestFolder() {
	_, err := os.Stat(testFolder)
	if !os.IsNotExist(err) {
		err = os.RemoveAll(testFolder)
		if err != nil {
			panic(err)
		}
	}
	err = os.Mkdir(testFolder, 0700)
	if err != nil {
		panic(err)
	}
}

func TestMain(m *testing.M) {
	cleanTestFolder()
	code := m.Run()
	os.Exit(code)
}*/

func testCommandAndExpectStdout(t *testing.T, command func(io.Reader, Config) bool, data []byte, expectedResult string) {
	stdout := bytes.NewBufferString("")
	dataReader := bytes.NewReader(data)
	config := Config{
		DryRun:  true,
		Verbose: true,
		Stdout:  stdout,
	}

	returnValue := command(dataReader, config)
	if !returnValue {
		t.Errorf("Command returned false")
	}
	if dataReader.Len() != 0 {
		t.Errorf("Data stream not empty")
	}
	stdoutValue := stdout.String()
	if expectedResult != stdoutValue {
		extraNewline := "\n"
		if strings.HasSuffix(expectedResult, "\n") && strings.HasSuffix(stdoutValue, "\n") {
			extraNewline = ""
		}
		t.Errorf("command stdout mismatch: expected\n\t%s%sbut was\n\t%s%s", expectedResult, extraNewline, stdoutValue, extraNewline)
	}
}

func TestMkdirCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1\x03\x00\x08\x00\x01\x23\x45\x67\x89\xab\xcd\xef")
	testCommandAndExpectStdout(t, mkdirCommand, data, "mkdir dir1 (17279655951921914625)\n")
}

func TestLinkCommand(t *testing.T) {
	data := []byte("\x0f\x00\x04\x00dir1\x11\x00\x04\x00dir2")
	testCommandAndExpectStdout(t, linkCommand, data, "link dir1 to dir2\n")
}
