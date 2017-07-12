/*
http://www.apache.org/licenses/LICENSE-2.0.txt

Copyright 2017 Intel Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

Author: Saurabh Deochake
*/

// A wrapper around fio and nvme-cli master

package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

//verify environment variable
func checkEnvironment() bool {
    res := verifyGoPath()
    if res == true {
        output, err := exec.Command("fio").CombinedOutput()
        if err != nil {
            os.Stderr.WriteString(err.Error())
        }
        fmt.Println(string(output))
    }else{
        res := setGoPath()
        if res == false{
            return false
        }
    }
}
// check if GOPATH is set
func verifyGoPath() bool{
    fmt.Println("Verifying GOPATH...")
    _, err := exec.Command("echo","$GOPATH").CombinedOutput()
    if err != nil{
        return false
    }
    return true
}

// if GOPATH is not set, then set it
func setGoPath() bool{
    fmt.Println("Setting GOPATH environment variable...")
    output, err := exec.Command("whereis","go").CombinedOutput()
    if err != nil {
        os.Stderr.WriteString(err.Error())
        return false
    }else{
        environment := "GOPATH="+
        strings.Replace(strings.Split(string(output),":")[1]," ","", -1)
        fmt.Println("Env:%s", environment)
        _, err := exec.Command("export",environment).CombinedOutput()
        if err != nil{
            return false
        }
    }
    return true
}


func main(){
    res := checkEnvironment()
    if res == false{
        fmt.Println("Could not set GOPATH. Please set it manually.")
        os.exit(1)
    }else{
        fmt.Println("GOPATH already set up...")
    }
}
