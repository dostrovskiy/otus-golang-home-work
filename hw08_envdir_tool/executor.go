package main

import (
	"log"
	"os"
	"os/exec"
	"strings"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	execCmd := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	execCmd.Stdout = os.Stdout
	execCmd.Env = getEnv(execCmd.Environ(), env)
	err := execCmd.Run()
	if err != nil {
		log.Fatalf("execution error: %+v\n", err)
		return 1
	}
	return execCmd.ProcessState.ExitCode()
}

func getEnv(oldEnv []string, addEnv Environment) []string {
	newEnv := []string{}
	// take all old env variables excluding those that need to be removed
	for _, env := range oldEnv {
		idx := strings.Index(env, "=")
		if idx == -1 {
			continue
		}
		if v, ok := addEnv[env[:idx]]; ok && v.NeedRemove {
			continue
		}
		newEnv = append(newEnv, env)
	}
	// adding new env variables
	for k, v := range addEnv {
		if !v.NeedRemove {
			newEnv = append(newEnv, k+"="+v.Value)
		}
	}
	return newEnv
}
