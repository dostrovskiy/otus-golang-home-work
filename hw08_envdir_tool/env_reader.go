package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// ErrInvalidEnvName is returned when file name contains '='.
var ErrInvalidEnvName = errors.New("invalid env name")

// Environment represents a set of environment variables.
type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	envs := make(Environment, len(files))
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if strings.Contains(file.Name(), "=") {
			return nil, ErrInvalidEnvName
		}
		eval, err := readEnv(dir, file)
		if err != nil {
			return nil, err
		}
		envs[file.Name()] = *eval
	}
	return envs, nil
}

func readEnv(dir string, file os.DirEntry) (eval *EnvValue, err error) {
	finfo, err := file.Info()
	if err != nil {
		return nil, err
	}
	if finfo.Size() == 0 {
		return &EnvValue{NeedRemove: true}, nil
	}
	fi, err := os.Open(filepath.Join(dir, file.Name()))
	if err != nil {
		return nil, err
	}
	defer fi.Close()
	sc := bufio.NewScanner(fi)
	if !sc.Scan() {
		if err := sc.Err(); err != nil {
			return nil, err
		}
		return &EnvValue{Value: ""}, nil
	}
	bs := bytes.ReplaceAll(sc.Bytes(), []byte{'\x00'}, []byte{'\n'})
	return &EnvValue{Value: strings.TrimRight(string(bs), "\t ")}, nil
}
