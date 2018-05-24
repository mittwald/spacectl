package spacefile

import "github.com/mittwald/spacectl/client/spaces"

type CronjobDef struct {
	Identifier    string             `hcl:",key"`
	Schedule      string             `hcl:"schedule"`
	AllowParallel bool               `hcl:"allowParallel"`
	Command       *CommandCronjobDef `hcl:"command"`
	Timezone      string             `hcl:"timezone"`
}

type CommandCronjobDef struct {
	Command          string   `hcl:"command"`
	Arguments        []string `hcl:"arguments"`
	WorkingDirectory string   `hcl:"workingDirectory"`
}

type CronjobDefList []CronjobDef

func (c CronjobDef) ToDeclaration() (spaces.Cronjob, error) {
	d := spaces.Cronjob{
		ID:            c.Identifier,
		AllowParallel: c.AllowParallel,
		Schedule:      c.Schedule,
		Timezone:      c.Timezone,
	}

	if c.Command != nil {
		d.Job.Type = "command"
		d.Job.Command = c.Command.Command
		d.Job.WorkingDirectory = c.Command.WorkingDirectory
		d.Job.Arguments = c.Command.Arguments
	}

	return d, nil
}
