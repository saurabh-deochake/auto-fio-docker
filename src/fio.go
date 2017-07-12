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
    if res == false {
        res := setGoPath()
        if res == false{
            return false
        }
    }
    result := verifyDocker()
    if result == true{
        return true
    }else{
        return false
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

//Create FIO Docker container
func createBenchmarkContainer() bool {
    fmt.Println("Creating Docker Container for benchmarking...")
    cmd :=
    "docker run --cap-add=SYS_ADMIN -it --device=/dev/nvme0n1:/dev/xvda:rw saurabhd04/docker_fio"
    // -------------------------- DO SOMETHING WITH ERR --------------------
    _, err :=
        exec.Command("bash", "-c", cmd).CombinedOutput()
    if err != nil {
        os.Stderr.WriteString(err.Error())
        return false
    }
    return true
}

//Verify if Docker is installed and running
func verifyDocker() bool{
    fmt.Println("Verifying Docker environment...")
    output, err :=
     exec.Command("bash", "-c", "rpm -qa | grep docker").CombinedOutput()
    if err != nil {
        fmt.Println("Inside error...")
        os.Stderr.WriteString(err.Error())
        return false
    }
    //fmt.Println(string(output))
    //return true

    if string(output[:]) == ""{
        fmt.Println("Is Docker installed? Please check and install Docker...")
        return false
    }else{
        // check if docker is running
        fmt.Println("Inside checking output slice")
        output, err =
            exec.Command("bash", "-c", "ps -ef | grep \"[d]ockerd\"").CombinedOutput()
        if err != nil {
            os.Stderr.WriteString(err.Error())
            return false
        }
        fmt.Println("Above dockerd check")

        if string(output[:]) == ""{
            fmt.Println("Dockerd is not running... Exiting!")
            return false
        }
        // Check if saurabhd04/docker_fio container is running
        fmt.Println("check docker_fio")

        output, err =
            exec.Command("bash", "-c", "docker ps | grep docker_fio").CombinedOutput()
            //s :=  string(output[:])
            //fmt.Println("This is %s",string(err.Error()))
        /*
        if err != nil {
            fmt.Println("Here?")
            os.Stderr.WriteString(err.Error())
            fmt.Println("Gela?")

            return false
        }*/

        fmt.Println("running?")

        if string(output[:]) == ""{
            fmt.Println("Docker container containing FIO not running...")
            res := createBenchmarkContainer()
            if res == true{
                fmt.Println("Docker container created...")
                return true
            }else{
                return false
            }
        }

    }
    return true
}

func main(){

    res := checkEnvironment()
    if res == false{
        fmt.Println("Could not set GOPATH. Please set it manually.")
        os.Exit(1)
    }

}
