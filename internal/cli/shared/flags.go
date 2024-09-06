package shared

import (
	"github.com/urfave/cli/v2"

	"github.com/sdisaacson/TriggerLab/internal/env"
)

var PortNumberFlag = &cli.UintFlag{
	Name:    "port",
	Aliases: []string{"p"},
	Usage:   "Server TCP port number",
	Value:   8080,
	EnvVars: []string{env.ListenPort.String(), env.Port.String()},
}
