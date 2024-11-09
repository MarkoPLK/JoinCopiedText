package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
)

var (
    xclip = "xclip"
    xclipPasteArgs = []string{xclip,"-o", "-selection", "clipboard"}
    xclipCopyArgs = []string{xclip, "-in", "-selection", "clipboard"}
)

func main() {
    fmt.Printf("Limpiando saltos de línea del texto copiado...\n\n")

    if err := checkXclip(); err != nil {
        fmt.Println("checkXclip: ", err.Error())
        os.Exit(1)
    }


    bufferCopy, err := readClipboard()
    if err != nil {
        fmt.Println("readClipboard: ", err.Error())
        os.Exit(1)
    }

    fmt.Printf("[INFO] Text copied:\n")
    color.Red("%s\n", bufferCopy)
    fmt.Println()

    content := cleanNewLines(bufferCopy)

    fmt.Printf("[INFO] Text free of line breaks:\n")
    color.Green("%s\n", content)

    err = writeClipboard(content)
    if err != nil {
        fmt.Println("writeClipboard : ", err.Error())
        os.Exit(1)
    }

}

func checkXclip() error {
    if _, err := exec.LookPath(xclip); err != nil {
        return errors.New("Xclip not installed")
    }
    return nil
}

func readClipboard() ([]byte, error) {
    
    copiedClipboardCmd := exec.Command(xclipPasteArgs[0], xclipPasteArgs[1:]...)

    bufferCopyClipboard, err := copiedClipboardCmd.Output()
    if err != nil {
        return nil, err
    }

    return bufferCopyClipboard, nil

}

func writeClipboard(content []byte) (error) {
    
    copiedClipboardCmd := exec.Command(xclipCopyArgs[0], xclipCopyArgs[1:]...)

    stdin, err := copiedClipboardCmd.StdinPipe()
    if err != nil {
        return err
    }

    go func() {
        defer stdin.Close()
        stdin.Write(content)
    }()

    if err := copiedClipboardCmd.Run(); err != nil {
        return err
    }

    return nil
}

func cleanNewLines(content []byte) ([]byte) {
    var res []byte = make([]byte, 0)

    var lastChar byte = byte(0)
    // var lastChar int = -1
    for _, c := range content {
        // Check for LF chars and if is not post "."
        if c == byte(10) && lastChar != byte(46) { 
            if lastChar == byte(45) { // if it is a word splitter "-" remove it
                res = res[:len(res)-1]
            } else {
                res = append(res, byte(32)) // Switch "\n" for a space
            }
        } else {
            res = append(res, c) 
        }
        lastChar = c
    }

    return res
}
