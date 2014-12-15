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

    filename := "/bin/wrapper"
    if _, err := os.Stat(filename); os.IsNotExist(err) {
       fmt.Printf("no such file or directory: %s\n", filename)

       err :=  os.Rename("/bin/sh", filename)
       if err != nil {
           fmt.Println(err)
           return
       }

       err2 := os.Symlink("/bin/wrapper", "/bin/sh")
       if err2 != nil {
           fmt.Println(err)
	   return
       }
       return
    }

    //os.Remove("/bin/sh")
    //os.Symlink("/bin/bash", "/bin/sh")

    /*
    uts, _ := uname()
    Machine := ""
    for _, c := range uts.Machine {
        if c == 0 {
            break
        }
        Machine += string(byte(c))
    }
    fmt.Println(Machine)
    */

    binary := "/usr/bin/qemu-arm-static"
    fmt.Println(args)
    //args = append([]string{binary, "/bin/busybox"}, args[2:]...)
    //args = append([]string{binary, "/bin/bash", "-c"}, args[2:]...)
    args_0 := args[0]
    if args_0 == "/bin/sh" {
       args_0 = "/bin/bash"
    }
    if args_0[0] != '/' {
       args_0 = "/bin/" + args_0
    }
    args = append([]string{binary, args_0}, args[1:]...)
    //args = []string{binary, "/bin/busybox", "ls"}
    fmt.Println(args)
    //os.Setenv("LD_PRELOAD", "/bin/ld_wrapper.so")
    env := os.Environ()

    fmt.Println("exec")

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        fmt.Println("Failed")
        panic(execErr)
    }

    fmt.Println("Cleaning")

    //os.Remove("/bin/sh")
    //os.Symlink("/bin/wrapper", "/bin/sh")
}