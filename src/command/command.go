/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package command

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/piot/cmake-generator/src/genconfig"
	"github.com/piot/deps/src/depslib"
)

type Options struct {
}

type Data struct {
	Name                    string
	SourceDirs              []string
	DependencyLibraryNames  []string
	DependencyLibraryString string
}

func Build(genConfigFilename string) error {
	config, err := genconfig.ReadGenConfigFromFilename(genConfigFilename)
	if err != nil {
		return err
	}

	configFiles, findErr := depslib.FindClosestConfigurationFiles(".")
	if findErr != nil {
		return findErr
	}
	if len(configFiles) == 0 {
		return fmt.Errorf("didn't find config files")
	}

	depsConfig, configErr := depslib.ReadConfigFromFilename(configFiles[0])
	if configErr != nil {
		return configErr
	}

	var dependencyLibraryNames []string
	for _, dep := range depsConfig.Dependencies {
		projectNames := strings.Split(dep.Name, "/")
		dependencyLibraryNames = append(dependencyLibraryNames, projectNames[1])
	}

	reader := bufio.NewReader(os.Stdin)

	octets, readErr := io.ReadAll(reader)
	if readErr != nil {
		return readErr
	}

	result, err := template.New("test").Parse(string(octets))
	if err != nil {
		return err
	}

	nameSplit := strings.Split(depsConfig.Name, "/")

	data := Data{Name: nameSplit[1], SourceDirs: config.SourceDirs, DependencyLibraryNames: dependencyLibraryNames,
		DependencyLibraryString: strings.Join(dependencyLibraryNames, " ")}

	log.Printf("data: %v", data)
	result.Execute(os.Stdout, data)

	return nil
}