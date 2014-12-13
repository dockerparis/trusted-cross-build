package main

import (
  "fmt"
  "os"
  "syscall"
)

type Utsname syscall.Utsname

func uname() (*syscall.Utsname, error) {
     uts := &syscall.Utsname{}

     if err := syscall.Uname(uts); err != nil {
     	return nil, err
	}
	return uts, nil
}

func main() {
    args := os.Args

    if len(args) < 1 {
       fmt.Println("Input file is missing.")
       os.Exit(1);
    }

    uts, _ := uname()
    Machine := ""
    for _, c := range uts.Machine { 
        if c == 0 { 
            break 
        } 
        Machine += string(byte(c)) 
    } 
    fmt.Println(Machine)

    binary := "/usr/bin/qemu-arm-static"
    fmt.Println(args)
    //args = append([]string{binary, "/bin/busybox"}, args[2:]...)
    args = append([]string{binary, "/bin/bash", "-c"}, args[2:]...)
    //args = []string{binary, "/bin/busybox", "ls"}
    fmt.Println(args)
    env := os.Environ()
    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}