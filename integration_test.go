package btrfsr2f

import (
	"bytes"
	"crypto/md5"
	"encoding/binary"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"testing"
)

const testFolder = "test_filesystem_root"
const quickBrownFox = "The quick brown fox jumps over the lazy dog.\n"

func prepareTestFolder(t *testing.T, emptyDirectory bool) {
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

	if !emptyDirectory {
		err = ioutil.WriteFile(path.Join(testFolder, "file1"), []byte(quickBrownFox), 0700)
		if err != nil {
			t.Fatal("Failed to create file1")
		}
		err = os.Mkdir(path.Join(testFolder, "dir1"), 0700)
		if err != nil {
			t.Fatal("Failed to create dir1")
		}
	}
}

func runProgram(t *testing.T, data []byte) {
	stdout := bytes.NewBufferString("")
	dataReader := bytes.NewReader(data)
	config := Config{
		DryRun:  false,
		Verbose: false,
		Stdout:  stdout,
		Root:    "test_filesystem_root",
	}

	ExecuteProgram(dataReader, config)
	if dataReader.Len() != 0 {
		t.Errorf("Data stream not empty")
	}
}

func makeCommandStream(command sendCommand, data []byte) []byte {
	// Stream header
	header := make([]byte, 0)
	header = append(header, []byte("btrfs-stream\x00")...) // Magic string
	header = append(header, []byte("\x01\x00\x00\x00")...) // Send version

	// Command header
	commandBytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(commandBytes, uint16(command))
	header = append(header, []byte("\x01\x02\x03\x04")...) // Command size
	header = append(header, commandBytes...)               // Command type
	header = append(header, []byte("\x01\x02\x03\x04")...) // CRC32

	// End command
	data = append(data, []byte("\x00\x00\x00\x00")...)
	data = append(data, []byte("\x15\x00")...)
	data = append(data, []byte("\x00\x00\x00\x00")...)

	return append(header, data...)
}

func simpleResolveLnkFile(pathToCheck string) string {
	if runtime.GOOS != "windows" || !strings.HasSuffix(pathToCheck, ".lnk") {
		return pathToCheck
	}
	// TODO: Make this better
	data, _ := ioutil.ReadFile(pathToCheck)
	s := string(data)
	start := strings.Index(s, "Data2\x00")
	if start == -1 {
		start = strings.Index(s, "Data3\x00")
	}
	start += 6
	if start == 5 { // TODO: This reeeeeeeaaally needs to change
		start = strings.Index(s, "\x00\x43\x3a") + 1
	}
	end := strings.Index(s[start:], "\x00") + start
	pathToCheck = s[start:end]

	newPath, _ := filepath.Abs(testFolder)
	newPath, _ = filepath.Rel(newPath, pathToCheck)
	return path.Join(testFolder, newPath)
}

func compareHashes(t *testing.T, hashSource string) {
	hashlist := make([]string, 0)
	err := filepath.Walk(testFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			realPath := simpleResolveLnkFile(path)
			data, err := ioutil.ReadFile(realPath)
			if err != nil {
				t.Fatal(err)
			}
			hash := md5.Sum(data)
			path = strings.ReplaceAll(path, testFolder, ".")
			path = strings.ReplaceAll(path, "\\", "/")
			path = strings.ReplaceAll(path, ".lnk", "")
			hashlist = append(hashlist, hex.EncodeToString(hash[:])+"  "+path)
		}
		return nil
	})
	sort.Strings(hashlist)
	hashesActual := strings.Join(hashlist, "\n") + "\n"
	hashesExpectationBytes, _ := ioutil.ReadFile(hashSource)
	hashesExpectation := string(hashesExpectationBytes)
	hashesExpectation = strings.ReplaceAll(hashesExpectation, "\r", "") // Git in Windows Docker image inserts \r
	if hashesActual != hashesExpectation {
		t.Fatalf("Hash mismatch: Expected\n%s  but was\n%s", hashesExpectation, hashesActual)
	}

	if err != nil {
		t.Fatal("Directory listing failed")
	}
}

func TestStream(t *testing.T) {
	prepareTestFolder(t, true)
	data, _ := ioutil.ReadFile("data/stream01.bin")
	runProgram(t, data)
	compareHashes(t, "data/stream01_files.txt")

	data, _ = ioutil.ReadFile("data/stream02.bin")
	runProgram(t, data)
	compareHashes(t, "data/stream02_files.txt")
}

func TestStreamWithCommandline(t *testing.T) {
	prepareTestFolder(t, true)
	args := [...]string{programName, "-o", testFolder, "-i", "data/stream01.bin"}
	Main(args[:])
	compareHashes(t, "data/stream01_files.txt")
}
