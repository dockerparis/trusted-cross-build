package main

import (
  "fmt"
  "os"
  "strings"
  "syscall"
)


func exists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
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

    binary := "/usr/bin/qemu-arm-static"
    fmt.Printf("Input args:  %s\n", args)

    // prevent double proxy
    args[0] = strings.Replace(args[0], "/binproxy/", "/bin/", -1)

    // "resolve" binary full path
    if args[0][0] != '/' {
       for _, pathdir := range []string{"/bin", "/usr/bin"} {
           // FIXME: resolve binaries using $PATH
           path := pathdir + "/" + args[0]
           _exists, _ := exists(path)
           if _exists {
               args[0] = path
               break
           }
       }
    }

    // If target is sh, replace with bash (sh is the current wrapper)
    if args[0] == "/bin/sh" {
       args[0] = "/bin/bash"
    }

    if len(args) > 1 {
       // if we are in a "sh -c" context, the binary is args[2]
       if args[0] == "/bin/bash" && args[1] == "-c" {
          args[2] = strings.Replace(args[2], "/bin/", "/binproxy/", -1)
       } else {
          args[1] = strings.Replace(args[1], "/bin/sh", "/binproxy/sh", -1)
       }
    }

    args = append([]string{binary}, args...)
    fmt.Printf("Output args: %s\n", args)

    //os.Setenv("LD_PRELOAD", "/bin/ld_wrapper.so")
    env := os.Environ()

    execErr := syscall.Exec(binary, args, env)
    if execErr != nil {
        panic(execErr)
    }
}