package main

import (
  "fmt"
  "os"
  "strings"
  "syscall"
)

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

    binary := "/usr/bin/qemu-arm-static"
    fmt.Println(args)
    args[0] = strings.Replace(args[0], "/binproxy/", "/bin/", -1)
    args_0 := args[0]
    if args_0 == "/bin/sh" {
       args_0 = "/bin/bash"
    }
    args_0 = strings.Replace(args_0, "/binproxy/", "/bin/", -1)
    if args_0[0] != '/' {
       args_0 = "/bin/" + args_0
    }
    if args[0] == "/bin/sh" && args[1] == "-c" {
       args[2] = strings.Replace(args[2], "/bin/", "/binproxy/", -1)
       args[2] = strings.Replace(args[2], "/bin/proxysh", "/bin/bash", -1)
    }
    args = append([]string{binary, args_0}, args[1:]...)
    fmt.Println(args)
    //os.Setenv("LD_PRELOAD", "/bin/ld_wrapper.so")
    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}