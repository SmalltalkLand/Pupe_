package shsrvr

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
	"time"
)

type Pipe struct {
	reader io.Reader
	writer io.Writer
}

func add(ports map[string]Pipe, plugins *[]Pipe, outside Pipe, prefix []string, target Pipe) {
	scanner := bufio.NewScanner(target.reader)
	for scanner.Scan() {
		t := scanner.Text()
		plugs := parse(strings.Split(t, " "))
		go plugy(plugs, ports, plugins, outside, prefix)
	}
}

var files map[string]os.File

func open(path string) *os.File {
	val, ok := files[path]
	if ok {
		return &val
	}
	x, _ := os.Open(path)
	files[path] = *x
	return x
}

const rscharset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var rsseededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func RsStringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rsseededRand.Intn(len(charset))]
	}
	return string(b)
}

func RsString(length int) string {
	return RsStringWithCharset(length, rscharset)
}
func get(ports map[string]Pipe, key string, other string, plugins *[]Pipe, outside Pipe, prefix []string) Pipe {
	if key == "other" {
		return get(ports, other, key, plugins, outside, prefix)
	}
	if strings.HasPrefix(key, "open:") {
		x := strings.Split(key, ":")
		var read, write *os.File
		read = open(x[1])
		write = open(x[2])
		return Pipe{
			reader: read,
			writer: write,
		}
	}
	if strings.HasPrefix(key, "sub:") {
		x := strings.Split(key, ":")
		y, _ := base64.StdEncoding.DecodeString(x[1])
		z := string(y)
		in1, out1 := io.Pipe()
		in2, out2 := io.Pipe()
		ports2 := make(map[string]Pipe)
		for k, v := range ports {
			ports2[k] = v
		}
		go plugy(parse(strings.Split(z, " ")), ports2, plugins, Pipe{
			reader: in1,
			writer: out2,
		}, prefix)
		return Pipe{
			reader: in2,
			writer: out1,
		}
	}
	return ports[key]
}
func cache(p func(n int) Pipe, pp []string, target []int) []Pipe {
	if pp[0] == pp[1] {
		pp := p(0)
		x := make([]Pipe, 2)
		x[0] = pp
		x[1] = pp
		return x
	}
	x := make([]Pipe, 2)
	x[0] = p(0)
	x[1] = p(1)
	return x
}
func other(x int) int {
	if x == 0 {
		return 1
	}
	return 0
}
func plugy(plugs map[string]string, ports map[string]Pipe, plugins *[]Pipe, outside Pipe, prefix []string) {
	keys := make([]string, 0, len(plugs))
	zeroone := make([]int, 2)
	zeroone[0] = 0
	zeroone[1] = 1
	for k := range plugs {
		keys = append(keys, k)
	}
	for _, k := range keys {
		reader, writer := io.Pipe()
		ports[k] = Pipe{
			reader: reader,
			writer: writer,
		}
	}
	for _, v := range plugs {
		split := strings.Split(v, " ")
		for i := 1; i < 2; i++ {
			k := split[i]
			_, ok := ports[k]
			if !ok {
				reader, writer := io.Pipe()
				ports[k] = Pipe{
					reader: reader,
					writer: writer,
				}
			}
		}
	}
	ports["outside"] = outside
	for k, v := range plugs {
		split := strings.Split(v, " ")
		for _, plugin := range *plugins {
			scanner := bufio.NewScanner(plugin.reader)
			writer := bufio.NewWriter(plugin.writer)
			s := base64.StdEncoding.EncodeToString([]byte(v + ";" + k + ";"))
			fmt.Fprintf(writer, "init %s", s)
			for {
				line := scanner.Text()
				split2 := strings.Split(line, " ")
				if split2[0] == "assign" && split2[1] == s {
					split[0] = split2[3]
					split[1] = split2[4]
					k = split2[2]
				}
				if split2[0] == "done" && split2[1] == s {
					break
				}
			}
		}
		if strings.HasPrefix(k, "plugin:") {
			x := cache(func(i int) Pipe { return get(ports, split[i], split[other(i)], plugins, outside, prefix) }, split, zeroone)
			plugin := Pipe{
				reader: x[0].reader,
				writer: x[1].writer,
			}
			*plugins = append(*plugins, plugin)
			scanner := bufio.NewScanner(plugin.reader)
			writer := bufio.NewWriter(plugin.writer)
			go func() {
				writer.WriteString("")
				for {
					line := scanner.Text()
					split2 := strings.Split(line, " ")
					if split2[0] == "swap" {
						temp := ports[split2[1]]
						ports[split2[1]] = ports[split2[2]]
						ports[split2[2]] = temp
					}
				}
			}()
		}
		if k == "adder" {
			x := cache(func(i int) Pipe { return get(ports, split[i], split[other(i)], plugins, outside, prefix) }, split, zeroone)
			go add(ports, plugins, outside, prefix, Pipe{
				reader: x[0].reader,
				writer: x[1].writer,
			})
		} else if k == "sandbox" {
			ports2 := make(map[string]Pipe)
			p3 := make([]string, 0)
			mode := 0
			for i, x := range split {
				if i > 2 {
					if x == "-p" {
						mode = 1
						continue
					}
					if x == "-o" {
						mode = 0
						continue
					}
					if mode == 0 {
						split2 := strings.Split(x, "->")
						ports2[strings.Join(split2[1:], "->")] = ports[split2[0]]
					}
					if mode == 1 {
						p3 = append(p3, x)
					}
				}
			}
			p2 := make([]Pipe, 0)
			for _, plugin := range *plugins {
				p2 = append(p2, plugin)
				scanner := bufio.NewScanner(plugin.reader)
				writer := bufio.NewWriter(plugin.writer)
				s := base64.StdEncoding.EncodeToString([]byte(v + ";" + k + ";"))
				fmt.Fprintf(writer, "sand %s", s)
				for {
					line := scanner.Text()
					split2 := strings.Split(line, " ")
					if split2[0] == "inject" && split2[1] == s {
						ports2[split2[2]] = ports[split2[3]]
					}
					if split2[0] == "done" && split2[1] == s {
						break
					}
				}
			}
			onetwo := make([]int, 2)
			onetwo[0] = 1
			onetwo[1] = 2
			x := cache(func(i int) Pipe { return get(ports, split[i], split[other(i)], plugins, outside, prefix) }, split[1:], onetwo)
			go add(ports, &p2, Pipe{
				reader: x[0].reader,
				writer: x[1].writer,
			}, p3, ports[split[0]])
		} else {
			split2 := split[2:]
			extra := make(map[int]Pipe)
			rw := make(map[int]int)
			for {
				if split2[0] == "-+r" {
					target, err := strconv.Atoi(split2[1])
					if err == nil {
						rw[target] = 0
						extra[target] = get(ports, split2[2], "", plugins, outside, prefix)
						split2 = split2[3:]
						continue
					}
				}
				if split2[0] == "-+rw" {
					target, err := strconv.Atoi(split2[1])
					if err == nil {
						rw[target] = 2
						extra[target] = get(ports, split2[2], "", plugins, outside, prefix)
						split2 = split2[3:]
						continue
					}
				}
				if split2[0] == "-+w" {
					target, err := strconv.Atoi(split2[1])
					if err == nil {
						rw[target] = 1
						extra[target] = get(ports, split2[2], "", plugins, outside, prefix)
						split2 = split2[3:]
						continue
					}
				}
				break
			}
			files := make([]*os.File, 0)
			for k, v := range rw {
				w := extra[k]
				var xr io.Reader
				var xw io.Writer
				if v == 0 {
					xr = w.reader
				} else {
					xw = w.writer
					if v == 2 {
						xr = w.reader
					}
				}
				var filer *os.File
				var filew *os.File
				if v == 0 || v == 2 {
					keyr := "/tmp/plgy/" + base64.StdEncoding.EncodeToString([]byte(RsString(25)))
					syscall.Mkfifo(keyr, 0700)
					filer, _ = os.Open(keyr)
				}
				if v == 1 {
					keyw := "/tmp/plgy/" + base64.StdEncoding.EncodeToString([]byte(RsString(25)))
					syscall.Mkfifo(keyw, 0700)
					filew, _ = os.Open(keyw)
				}
				if v == 0 {
					go io.Copy(filer, xr)
				} else {
					go io.Copy(xw, filew)
					if v == 2 {
						go io.Copy(filer, xr)
					}
				}
				if v == 0 {
					files = append(files, filer)
				} else {
					files = append(files, filew)
					if v == 2 {
						files = append(files, filer)
					}
				}
			}
			x := cache(func(i int) Pipe { return get(ports, split[i], split[other(i)], plugins, outside, prefix) }, split, zeroone)
			cmd := exec.Command("/usr/bin/env", append(prefix, split2...)...)
			cmd.Stdin = x[0].reader
			cmd.Stdout = x[1].writer
			cmd.ExtraFiles = files
			go cmd.Run()
		}
	}
}
func parse(x []string) map[string]string {
	plugs := make(map[string]string)
	i := 0
	vv := ""
	vvv := ""
	j := 0
	for _, v := range x {
		if i == 0 {
			i = 1
			vv = v
		} else {
			if v == "--" && j == 0 {
				plugs[vv] = vvv
				i = 0
				j = 0
				vvv = ""
			} else if v == "/--" {
				j = 1
			} else {
				vvv += " "
				vvv += v
				j = 0
			}
		}
	}
	return plugs
}
func main() {
	plugs := parse(os.Args)
	plugins := make([]Pipe, 0)
	plugy(plugs, make(map[string]Pipe), &plugins, Pipe{
		reader: os.Stdin,
		writer: os.Stdout,
	}, make([]string, 0))
}
