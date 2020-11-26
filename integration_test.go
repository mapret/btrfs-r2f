package main

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

func prepareTestFolder(t *testing.T) {
	prepareTestFolder2(t, false)
}

func prepareTestFolder2(t *testing.T, emptyDirectory bool) {
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

func TestLinkStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file2\x11\x00\x05\x00file1")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_LINK, data))

	stat1, _ := os.Stat(path.Join(testFolder, "file1"))
	stat2, _ := os.Stat(path.Join(testFolder, "file2"))
	if !os.SameFile(stat1, stat2) {
		t.Fatal("file1 and file2 are not the same")
	}
}

func TestMkdirStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x04\x00dir2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_MKDIR, data))

	stat, err := os.Stat(path.Join(testFolder, "dir2"))
	if os.IsNotExist(err) || !stat.IsDir() {
		t.Fatal("dir2 was not created")
	}
}

func TestMkfileStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_MKFILE, data))

	stat, err := os.Stat(path.Join(testFolder, "file2"))
	if os.IsNotExist(err) || stat.IsDir() {
		t.Fatal("file2 was not created")
	}
}

func TestRenameStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file1\x10\x00\x07\x00newname")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_RENAME, data))

	stat, err := os.Stat(path.Join(testFolder, "newname"))
	if os.IsNotExist(err) || stat.IsDir() {
		t.Fatal("file1 was not renamed to newname")
	}
}

func TestRmdirStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x04\x00dir1")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_RMDIR, data))

	_, err := os.Stat(path.Join(testFolder, "newname"))
	if !os.IsNotExist(err) {
		t.Fatal("dir1 was not deleted")
	}
}

func TestSymlinkStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file2\x03\x00\x08\x00\xef\xcd\xab\x89\x67\x45\x23\x01\x11\x00\x05\x00file1")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_SYMLINK, data))

	if runtime.GOOS == "windows" {
		_, err := os.Stat(path.Join(testFolder, "file2.lnk"))
		if os.IsNotExist(err) {
			t.Fatal("file2 was not created")
		}
	} else {
		stat, err := os.Lstat(path.Join(testFolder, "file2"))
		if os.IsNotExist(err) {
			t.Fatal("file2 was not created")
		}
		if stat.Mode()&os.ModeSymlink == 0 {
			t.Fatal("file2 is not a symlink")
		}
	}
}

func TestTruncateStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file1\x04\x00\x08\x00\x10\x00\x00\x00\x00\x00\x00\x00")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_TRUNCATE, data))

	stat, _ := os.Stat(path.Join(testFolder, "file1"))
	if stat.Size() != 16 {
		t.Fatal("file1 was not truncated")
	}
	file1bytes, _ := ioutil.ReadFile(path.Join(testFolder, "file1"))
	if string(file1bytes) != "The quick brown " {
		t.Fatal("file1 was not truncated correctly")
	}
}

func TestUnlinkStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file1")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_UNLINK, data))

	_, err := os.Stat(path.Join(testFolder, "file1"))
	if !os.IsNotExist(err) {
		t.Fatal("file1 was not deleted")
	}
}

func TestWriteStream(t *testing.T) {
	prepareTestFolder(t)
	data := []byte("\x0f\x00\x05\x00file1\x12\x00\x08\x00\x04\x00\x00\x00\x00\x00\x00\x00\x13\x00\x05\x00QUICK")
	runProgram(t, makeCommandStream(BTRFS_SEND_C_WRITE, data))

	stat, _ := os.Stat(path.Join(testFolder, "file1"))
	if stat.Size() != int64(len(quickBrownFox)) {
		t.Fatal("size of file1 changed")
	}
	file1bytes, _ := ioutil.ReadFile(path.Join(testFolder, "file1"))
	if string(file1bytes) != strings.ReplaceAll(quickBrownFox, "quick", "QUICK") {
		t.Fatal("file1 was not written to correctly")
	}
}

func simpleResolveLnkFile(pathToCheck string) string {
	if runtime.GOOS != "windows" || !strings.HasSuffix(pathToCheck, ".lnk") {
		return pathToCheck
	}
	// TODO: Make this better
	data, _ := ioutil.ReadFile(pathToCheck)
	s := string(data)
	start := strings.Index(s, "Data2\x00") + 6
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
	prepareTestFolder2(t, true)
	data, _ := ioutil.ReadFile("data/stream01.bin")
	runProgram(t, data)
	compareHashes(t, "data/stream01_files.txt")

	data, _ = ioutil.ReadFile("data/stream02.bin")
	runProgram(t, data)
	compareHashes(t, "data/stream02_files.txt")
}

func TestStreamWithCommandline(t *testing.T) {
	prepareTestFolder2(t, true)
	args := [...]string{programName, "-o", testFolder, "-i", "data/stream01.bin"}
	Main(args[:])
	compareHashes(t, "data/stream01_files.txt")
}
