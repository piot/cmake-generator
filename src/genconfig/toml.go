/*---------------------------------------------------------------------------------------------
 *  Copyright (c) Peter Bjorklund. All rights reserved.
 *  Licensed under the MIT License. See LICENSE in the project root for license information.
 *--------------------------------------------------------------------------------------------*/

package genconfig

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pelletier/go-toml/v2"
)

type ArtifactType uint

const (
	Library ArtifactType = iota
	Executable
)

func (e ArtifactType) String() string {
	switch e {
	case Library:
		return "library"
	case Executable:
		return "executable"
	default:
		return fmt.Sprintf("%v", e)
	}
}

type GenConfig struct {
	CmakeGenVersion string
	Name            string
	ArtifactType    string `toml:"artifact_type"`
	SourceDirs      []string
	Dependencies    []string
}

func ReadGenConfigFromReader(reader io.Reader) (*GenConfig, error) {
	tomlString, tomlParseErr := io.ReadAll(reader)
	if tomlParseErr != nil {
		return nil, tomlParseErr
	}
	config := &GenConfig{}
	unmarshalErr := toml.Unmarshal(tomlString, config)
	if unmarshalErr != nil {
		log.Printf("unmarshal: %v", unmarshalErr)
		return nil, unmarshalErr
	}
	if config.CmakeGenVersion != "0.0.0" {
		return nil, fmt.Errorf("wrong deps file format version '%v'", config.CmakeGenVersion)
	}

	if config.ArtifactType == "" {
		config.ArtifactType = "library"
	}

	return config, unmarshalErr
}

func ReadGenConfigFromFilename(filename string) (*GenConfig, error) {
	reader, openErr := os.Open(filename)
	if openErr != nil {
		return nil, openErr
	}
	return ReadGenConfigFromReader(reader)
}
