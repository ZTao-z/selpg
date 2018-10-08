package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

var (
	s, e, lines                                  int    = 1, 1, 72
	input_file, output_file, error_file, pipe    string = "", "", "", ""
	changeLine, isInFile, isOutFile, isErrorFile bool   = false, true, false, false
	err                                          error
	out                                          io.WriteCloser
)

func main() {
	err = getArgs()
	if err == nil {
		out = os.Stdout
		if pipe != "" {
			sendToPrinter()
		} else {
			err = readFile(input_file)
		}
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}

func getArgs() error {

	arg_len := len(os.Args)
	if arg_len < 3 {
		err = errors.New("fewer args")
		return err
	}

	for _, a := range os.Args[1:] {
		if a[0] == '-' {
			var pos int = 1
			if a[1] == '-' {
				pos = 2
			}
			switch a[pos] {
			case 's':
				s, err = strconv.Atoi(strings.Split(a, "-s")[1])
			case 'e':
				e, err = strconv.Atoi(strings.Split(a, "-e")[1])
			case 'l':
				lines, err = strconv.Atoi(strings.Split(a, "-l")[1])
			case 'f':
				changeLine = true
			case 'd':
				pipe = strings.Split(a, "-d")[1]
			case '-':
				err = errors.New("arguement wrong format")
			default:
				s, e, lines, changeLine = 1, 1, 72, false
			}
		} else {
			if a[0] == '<' {
				isInFile = true
				if isInFile && len(a) > 1 {
					os.Stdin, err = os.Open("data.in")
					isInFile = false
					input_file = ""
				}
			} else if a[0] == '>' {
				isOutFile = true
				if isOutFile && len(a) > 1 {
					os.Stdout, err = os.Create("data,out")
					isOutFile = false
					output_file = ""
				}
			} else if a[0] == '2' && a[1] == '>' {
				isErrorFile = true
				if isErrorFile && len(a) > 2 {
					error_file = a
				}
			} else if a[0] == '|' {

			} else {
				if isInFile {
					input_file = a
					isInFile = false
				}
			}
		}
		if err != nil {
			return err
		}
	}
	return nil
}

func readFile(filename string) error {
	f, err := os.Open(filename)
	if filename == "" {
		f = os.Stdin
	} else {
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	defer f.Close()
	buf := bufio.NewReader(f)

	var count, page int = 0, 1
	if !changeLine {
		// change page with line count
		for page = 1; page <= e; page += 1 {
			count = 1
			for {
				if count > lines {
					break
				}
				line, err := buf.ReadString('\n')
				count += 1
				if page < s {
					if err != nil {
						return err
					} else {
						continue
					}
				}
				//line = strings.TrimSpace(line)
				if err == nil {
					out.Write([]byte(line))
				} else if err == io.EOF {
					out.Write([]byte(line))
				} else {
					return err
				}
			}
		}
	} else {
		// change page with '\f'
		// echo -e '\fxxxx\n' >> test.txt
		page = 0
		for page < e {
			line, err := buf.ReadString('\f')
			page += 1
			if page < s {
				if err != nil {
					return err
				} else {
					continue
				}
			}
			if err == nil {
				out.Write([]byte(line))
			} else if err == io.EOF {
				out.Write([]byte(line))
			} else {
				return err
			}
		}
	}
	return nil
}

func sendToPrinter() {
	if pipe != "" {
		cmd := exec.Command("lp", "-d", pipe)
		out, err = cmd.StdinPipe()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		err = readFile(input_file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		err = readFile(input_file)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		out.Close()

		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		if err = cmd.Start(); err != nil {
			fmt.Fprintln(os.Stderr, err)
		}
		timer := time.AfterFunc(5*time.Second, func() {
			cmd.Process.Kill()
		})

		err = cmd.Wait()
		timer.Stop()
	}
}
