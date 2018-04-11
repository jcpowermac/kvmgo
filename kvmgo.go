package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	libvirt "github.com/libvirt/libvirt-go"
)

type GuestExecStatusArguments struct {
	Pid int `json:"pid"`
}

type GuestExecArguments struct {
	Path          string   `json:"path"`
	Arg           []string `json:"arg,omitempty"`
	Env           []string `json:"env,omitempty"`
	InputData     string   `json:"input-data,omitempty"`
	CaptureOutput bool     `json:"capture-output,omitempty"`
}
type ExecAgentCommand struct {
	Execute   string             `json:"execute"`
	Arguments GuestExecArguments `json:"arguments"`
}
type ExecStatusAgentCommand struct {
	Execute   string                   `json:"execute"`
	Arguments GuestExecStatusArguments `json:"arguments"`
}

type ExecAgentCommandReturn struct {
	Pid int `json:"pid"`
}

type ExecAgentCommandOutput struct {
	Return ExecAgentCommandReturn `json:"return"`
}

func main() {
	connect, err := libvirt.NewConnect("qemu+ssh://root@172.31.52.110/system")
	if err != nil {
		log.Fatal(err)
	}

	defer connect.Close()
	domain, err := connect.LookupDomainByName("win2k16")
	if err != nil {
		log.Fatal(err)
	}

	var command ExecAgentCommand
	command.Execute = "guest-exec"
	command.Arguments.Path = "C:/Windows/System32/setx.exe"
	command.Arguments.Arg = []string{"TESTVAR", "TESTVALUE", "/M"}
	command.Arguments.CaptureOutput = true

	c, err := json.Marshal(command)

	output, err := domain.QemuAgentCommand(string(c), 1000, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("output: %s", output)
	co := new(ExecAgentCommandOutput)
	err = json.Unmarshal([]byte(output), co)

	if err != nil {
		log.Fatal(err)
	}

	// yes I would change this to check the output...
	time.Sleep(time.Second * 10)
	var statusCommand ExecStatusAgentCommand
	statusCommand.Execute = "guest-exec-status"
	statusCommand.Arguments.Pid = (*co).Return.Pid

	sc, err := json.Marshal(statusCommand)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("sc: %s", sc)

	statusOutput, err := domain.QemuAgentCommand(string(sc), 1000, 0)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("statusOutput: %s", statusOutput)

}
