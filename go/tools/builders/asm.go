// Copyright 2017 The Bazel Authors. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// asm builds a single .s file with "go tool asm". It is invoked by the
// Go rules as an action.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
)

func run(args []string) error {
	// process the args
	if len(args) < 2 {
		return fmt.Errorf("Usage: asm -go gotool source.s -- <extra options>")
	}
	flags := flag.NewFlagSet("asm", flag.ExitOnError)
	goenv := envFlags(flags)
	if err := flags.Parse(args); err != nil {
		return err
	}
	if len(flags.Args()) < 1 {
		return fmt.Errorf("Missing source file to asm")
	}
	source := flags.Args()[0]
	remains := flags.Args()[1:]

	// filter our input file list
	bctx := goenv.BuildContext()
	matched, err := matchFile(bctx, source)
	if err != nil {
		return err
	}
	if !matched {
		source = os.DevNull
	}
	goargs := []string{"tool", "asm"}
	goargs = append(goargs, remains...)
	goargs = append(goargs, source)
	env := os.Environ()
	env = append(env, goenv.Env()...)
	cmd := exec.Command(goenv.Go, goargs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running assembler: %v", err)
	}
	return nil
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}
