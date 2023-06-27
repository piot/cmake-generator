/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/piot/cmake-generator/src/command"
)

var version string

// SharedOptions are command line shared options.
type SharedOptions struct {
}

// BuildCmd is the options for a build.
type BuildCmd struct {
	//	Template string   `type:"path" required:""`
	Config string `arg:"" type:"path" default:"./cmake_gen.toml" help:"the cmake gen configuration file"`
}

// Options are all the command line options.
type Options struct {
	Build BuildCmd `cmd:""`
}

// Run is called if a build command was issued.
func (o *BuildCmd) Run() error {
	if o.Config == "" {
		return fmt.Errorf("missing config file")
	}
	return command.Build(o.Config)
}

func main() {
	ctx := kong.Parse(&Options{})

	err := ctx.Run()
	if err != nil {
		fmt.Printf("ERROR:%v\n", err)
		os.Exit(-1)
	}
}
