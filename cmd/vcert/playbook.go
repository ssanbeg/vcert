/*
 * Copyright 2023 Venafi, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *  http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"

	"github.com/Venafi/vcert/v4/pkg/playbook/app/domain"
	"github.com/Venafi/vcert/v4/pkg/playbook/app/parser"
	"github.com/Venafi/vcert/v4/pkg/playbook/app/service"
	"github.com/Venafi/vcert/v4/pkg/playbook/util"
)

const (
	commandRunPlaybookName = "run"
)

var commandRunPlaybook = &cli.Command{
	Name: commandRunPlaybookName,
	Usage: `Enables users to request and retrieve one or more certificates, 
	install them as either CAPI, JKS, PEM, or PKCS#12, run after-install operations 
	(script, command-line instruction, etc.), and monitor certificate(s) for renewal 
	on subsequent runs.`,
	UsageText: `vcert run
   vcert run -f /path/to/my/file.yml
   vcert run -f ./myFile.yaml --force-renew
   vcert run -f ./myFile.yaml --debug`,
	Action: doRunPlaybook,
	Flags:  playbookFlags,
}

type runOptions struct {
	debug    bool
	filepath string
	force    bool
}

var (
	playbookOptions = runOptions{}

	PBFlagDebug = &cli.BoolFlag{
		Name:        "debug",
		Aliases:     []string{"d"},
		Usage:       "Enables debug log messages",
		Required:    false,
		Value:       false,
		Destination: &playbookOptions.debug,
	}

	PBFlagFilepath = &cli.StringFlag{
		Name:        "file",
		Aliases:     []string{"f"},
		Usage:       "the path to the playbook file to be run",
		Required:    true,
		Value:       domain.DefaultFilepath,
		Destination: &playbookOptions.filepath,
	}

	PBFlagForce = &cli.BoolFlag{
		Name:        "force-renew",
		Aliases:     nil,
		Usage:       "forces certificate renewal regardless of expiration date or renew window",
		Required:    false,
		Value:       false,
		Destination: &playbookOptions.force,
	}

	playbookFlags = flagsApppend(
		PBFlagDebug,
		PBFlagFilepath,
		PBFlagForce,
	)
)

func doRunPlaybook(_ *cli.Context) error {
	err := util.ConfigureLogger(playbookOptions.debug)
	if err != nil {
		return err
	}
	zap.L().Info("running playbook file", zap.String("file", playbookOptions.filepath))
	zap.L().Debug("debug is enabled")

	playbook, err := parser.ReadPlaybook(playbookOptions.filepath)
	if err != nil {
		zap.L().Error(fmt.Errorf("%w", err).Error())
		os.Exit(1)
	}

	_, err = playbook.IsValid()
	if err != nil {
		zap.L().Error("invalid playbook file", zap.String("file", playbookOptions.filepath), zap.Error(err))
		os.Exit(1)
	}

	//Set the forceRenew variable
	playbook.Config.ForceRenew = playbookOptions.force

	if len(playbook.CertificateTasks) == 0 {
		zap.L().Info("no tasks in the playbook. Nothing to do")
		return nil
	}

	// emulate the setTLSConfig from vcert
	err = setTLSConfig()
	if err != nil {
		zap.L().Error("tls config error", zap.Error(err))
		os.Exit(1)
	}

	if playbook.Config.Connection.Type == domain.CTypeTPP {
		err = service.ValidateTPPCredentials(&playbook)
		if err != nil {
			zap.L().Error("invalid tpp credentials", zap.Error(err))
			os.Exit(1)
		}
	}

	for _, certTask := range playbook.CertificateTasks {
		zap.L().Info("running playbook task", zap.String("task", certTask.Name))
		errors := service.Execute(playbook.Config, certTask)
		if len(errors) > 0 {
			for _, err2 := range errors {
				zap.L().Error("error running task", zap.String("task", certTask.Name), zap.Error(err2))
			}
			os.Exit(1)
		}
	}

	zap.L().Info("playbook run finished")
	return nil
}