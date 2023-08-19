/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package command

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
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
	SourceFilesString       string
	DependencyLibraryNames  []string
	DependencyLibraryString string
	ArtifactType            string
}

func GetLibraryNameFromRepo(repo string) string {
	projectName := strings.Split(repo, "/")[1]
	if strings.HasSuffix(projectName, "-c") {
		return projectName[:len(projectName)-2]
	}

	return projectName
}

func Build(genConfigFilename string) error {
	absoluteConfigPath, absErr := filepath.Abs(genConfigFilename)
	if absErr != nil {
		return absErr
	}

	absoluteConfigDirectory := filepath.Dir(absoluteConfigPath)

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

	if len(config.Dependencies) == 0 {
		for _, dep := range depsConfig.Dependencies {
			dependencyLibraryNames = append(dependencyLibraryNames, GetLibraryNameFromRepo(dep.Name))
		}
	} else {
		dependencyLibraryNames = config.Dependencies
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

	cmakeName := config.Name
	if config.Name == "" {
		cmakeName = GetLibraryNameFromRepo(depsConfig.Name)

	}

	publicDependencyLibraryNames := ""

	if len(dependencyLibraryNames) > 0 {
		const libraryNamesSeparator = "\n  "
		publicDependencyLibraryNames = libraryNamesSeparator + strings.Join(dependencyLibraryNames,
			libraryNamesSeparator)
	}

	var sourceFiles []string
	for _, sourceDir := range config.SourceDirs {
		globWildcard := filepath.Join(absoluteConfigDirectory, sourceDir, "*.c")
		matches, err := filepath.Glob(globWildcard)
		if err != nil {
			return err
		}

		for _, match := range matches {
			relativePath, relativeErr := filepath.Rel(absoluteConfigDirectory, match)
			if relativeErr != nil {
				return relativeErr
			}
			sourceFiles = append(sourceFiles, relativePath)
		}

	}

	sort.Strings(sourceFiles)
	const sourceFileSeparator = "\n  "
	sourceFilesString := sourceFileSeparator + strings.Join(sourceFiles, sourceFileSeparator)

	data := Data{
		Name:                    cmakeName,
		SourceDirs:              config.SourceDirs,
		DependencyLibraryNames:  dependencyLibraryNames,
		SourceFilesString:       sourceFilesString,
		ArtifactType:            config.ArtifactType,
		DependencyLibraryString: publicDependencyLibraryNames}

	return result.Execute(os.Stdout, data)
}
