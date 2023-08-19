/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package sub_dir

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/piot/cmake-generator/src/command"
	"github.com/piot/deps/src/depslib"
)

func getProjectNameFromRepo(repo string) string {
	projectName := strings.Split(repo, "/")[1]

	return projectName
}

func stringSliceContains(slice []string, searchFor string) bool {
	for _, str := range slice {
		if str == searchFor {
			return true
		}
	}

	return false
}

func SubDir() error {
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

	fmt.Printf("cmake_minimum_required(VERSION 3.16.3)\nproject(%v C)\n\n",
		command.GetLibraryNameFromRepo(depsConfig.Name))

	firstConfig := configFiles[0]
	depsRootDir := filepath.Dir(filepath.Dir(filepath.Dir(firstConfig)))
	depsPath := filepath.Dir(firstConfig)
	log.Printf("depsRoot: %v, depsPath: %v", depsRootDir, depsPath)

	cache, _, err := depslib.CalculateTotalDependencies(depsRootDir, depsPath, depsConfig, depslib.ReadLocal, false)
	if err != nil {
		return err
	}

	var dependencyPackageNames []string
	for _, packageDependency := range cache.Nodes {
		if packageDependency.Name() == depsConfig.Name {
			continue
		}
		dependencyPackageNames = append(dependencyPackageNames, packageDependency.Name())
	}

	sort.Strings(dependencyPackageNames)

	for _, name := range dependencyPackageNames {
		fmt.Printf("add_subdirectory(deps/%v/src/lib)\n", name)
	}

	fmt.Printf("\n\n")

	normalSubDirectories := []string{"lib", "tests", "examples"}
	for _, subDirectory := range normalSubDirectories {
		if _, err := os.Stat(subDirectory); os.IsNotExist(err) {
			continue
		}
		fmt.Printf("add_subdirectory(%v)\n", subDirectory)
	}

	fmt.Printf("\n\n")

	var developmentPackageNames []string
	for _, packageDevDependency := range depsConfig.Development {
		if stringSliceContains(dependencyPackageNames, packageDevDependency.Name) {
			continue
		}
		developmentPackageNames = append(developmentPackageNames, packageDevDependency.Name)
	}

	sort.Strings(developmentPackageNames)
	for _, name := range developmentPackageNames {
		fmt.Printf("add_subdirectory(deps/%v/src/lib)\n", name)
	}
	return nil
}
