package main

import (
  "fmt"
  "os"
  "syscall"
)

func main() {
    args := os.Args

    if len(args) < 1 {
       fmt.Println("Input file is missing.");
       os.Exit(1);
    }

    binary := "/usr/bin/qemu-arm-static"
    fmt.Println(args)
    //args = append([]string{binary, "/bin/busybox"}, args[2:]...)
    args = []string{binary, "/bin/busybox", "ls"}
    fmt.Println(args)
    env := os.Environ()
    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}