package main

import (
  "flag"
  "fmt"
  "os"
  "syscall"
)

func usage() {
    fmt.Fprintf(os.Stderr, "usage: wrapper binary [args...]\n")
    flag.PrintDefaults()
    os.Exit(2)
}

func main() {
    flag.Usage = usage
    flag.Parse()
    args := flag.Args()

    if len(args) < 1 {
       fmt.Println("Input file is missing.");
       os.Exit(1);
    }

    binary := "/usr/bin/qemu-arm-static"
    args = append([]string{binary, "/bin/busybox", "sh", "-c"}, args...)
    fmt.Println(args)
    env := os.Environ()
    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}