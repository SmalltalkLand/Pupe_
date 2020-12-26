package av

import (
	"bufio"
	"io"
	"math/rand"
	"os"
	"strconv"
	"syscall"
	"time"
)

import "C"
import "github.com/fatih/color"
import "golang.org/x/sys/unix"
import "github.com/fsnotify/fsnotify"

func containssi(s []string, e string) bool {
    for _, a := range s {
        if strings.Contains(a,e){
            return true
        }
    }
    return false
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
func PidFdOpen(Pid int, Flags uint32) (*os.File, error) {
	pidFd, errno := C.pidfd_open(C.int(Pid), C.uint32_t(Flags))
	if errno != nil {
		return nil, errno
	}

	errno = unCloexec(int(pidFd))
	if errno != nil {
		return nil, errno
	}

	return os.NewFile(uintptr(pidFd), fmt.Sprintf("%d", Pid)), nil
}
func PidEnter(Pid int,flags int) error{
	fd, err := PidFdOpen(pid,0)
	if err != nil{return err}
	err = unix.Setns(fd.Fd(),flags)
	return err
}
func RsString(length int) string {
	return RsStringWithCharset(length, rscharset)
}
func getCmdLine(pid int) (b5 []string, err error) {
	f, err := os.Open("/proc/" + strconv.Itoa(pid) + "/cmdline")
	if err != nil {
		return nil, err
	}
	r4 := bufio.NewReader(f)
	b5 = make([]string, 0)
	for {
		b4, err := r4.ReadBytes(0)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		b5 = append(b5, string(b4))
	}
	return b5, nil
}
func getExePath(pid int) (exe string, err error) {
	exe, err = os.Readlink("/proc/" + strconv.Itoa(pid) + "/exe")
	return exe, err
}
func isSelfPid(pid int) bool {
	return pid == os.Getpid()
}
func error_(pid int) {
	for {
		e := syscall.Kill(pid, 9)
		if e == nil {
			return
		}
	}
}
func getRoot(pid int,procfs string) (string, error){
	root, err := os.Readlink(procfs + "/" + strconv.Itoa(pid) + "/root")
	return root, err
}
func watch(x func(y fsnotify.Watcher, z chan bool),path string,done chan bool){
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		x(watcher,done)
	}()

	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}
}
func main2(pid int) {
	oldSelf, err := syscall.Getpriority(syscall.PRIO_PROCESS, os.Getpid())
	if err != nil {
		error_(pid)
		return
	}
	err = syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), -20)
	if err != nil {
		error_(pid)
		return
	}
	oldOther, err := syscall.Getpriority(syscall.PRIO_PROCESS, pid)
	if err != nil {
		error_(pid)
		return
	}
	err = syscall.Setpriority(syscall.PRIO_PROCESS, pid, 19)
	if err != nil {
		error_(pid)
		return
	}
	defer syscall.Setpriority(syscall.PRIO_PROCESS, pid, oldOther)
	defer syscall.Setpriority(syscall.PRIO_PROCESS, os.Getpid(), oldSelf)
	pidFlags := syscall.CLONE_NEWCGROUP | syscall.CLONE_NEWIPC | syscall.CLONE_NEWNET | syscall.CLONE_NEWNS | syscall.CLONE_NEWTIME | syscall.CLONE_NEWUSER | syscall.CLONE_NEWUTS
	Ppid := os.getppid()
	defer PidEnter(Ppid,pidFlags)
	err = PidEnter(pid,pidFlags)
	if err != nil {
		error_(pid)
		return
	}
	procfs := "/" + RsString(100)
	err = syscall.Mount("x",procfs,"procfs",0,"")
	defer syscall.Unmount(procfs,0)
	if err != nil {
		error_(pid)
		return
	}
	rootfs := "/" + RsString(100)
	childRoot, err := getRoot(pid,procfs)
	if err != nil {
		error_(pid)
		return
	}
	err = syscall.Mount("/",childRoot + rootfs,"bind",0,"")
	if err != nil {
		error_(pid)
		return
	}
	err = syscall.Chroot(childRoot)
	if err != nil {
		error_(pid)
		return
	}
	defer  syscall.Chroot(rootfs)
	done := make(chan bool)
	oldGor := runtime.NumGoroutine()
	badFiles := []string{"/etc/ld.so.preload"}
watch(func(y  fsnotify.Watcher, z chan bool){
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				link, err := os.Readlink(event.Name)
				if containssi(badFiles,link){
					error_(pid)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	},procfs + "/" + strconv.Itoa(pid) + "/fd",done)

	gor := runtime.NumGoroutine() - oldGor
	for i := 1; i < gor; i++ {
		_ = <- done
	}
//	files, err := ioutil.ReadDir(procfs + "/" + strconv.Itoa(pid) + "/fd")
//	if err != nil {
//		error_(pid)
//		return
//	}
//	for _, file := range files {
	//	filename := procfs + "/" + strconv.Itoa(pid) + "/fd/" + file.Name()
	//	
	//	if err != nil {
	//		error_(pid)
	//		return
	//	}

	//}
}
func main() {
	if os.Args[2] != "-s"{
		cmd := exec.Command("/proc/self/exe",append([]string{os.Args[1],"-s"},os.Args[2:])...)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Run()
		return
	}
	pid, err := strconv.Atoi(os.Args[1])
	if err != nil {
		return
	}
	syscall.PtraceAttach(pid)
	if isSelfPid(pid) {
		syscall.PtraceDetach(pid)
		return
	}
	main2(pid)
	syscall.PtraceDetach(pid)
}
