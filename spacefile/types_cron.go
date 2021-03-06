package spacefile

import "github.com/mittwald/spacectl/client/spaces"

type CronjobDef struct {
	Identifier    string             `hcl:",key"`
	Schedule      string             `hcl:"schedule"`
	AllowParallel bool               `hcl:"allowParallel"`
	Command       *CommandCronjobDef `hcl:"command" hcle:"omitempty"`
	Timezone      string             `hcl:"timezone" hcle:"omitempty"`
}

type CommandCronjobDef struct {
	Command          string   `hcl:"command" hcle:"omitempty"`
	Arguments        []string `hcl:"arguments" hcle:"omitempty"`
	WorkingDirectory string   `hcl:"workingDirectory" hcle:"omitempty"`
}

type CronjobDefList []CronjobDef

func (c CronjobDef) ToDeclaration() spaces.Cronjob {
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

	return d
}

// CronjobFromDeclaration returns a CronjobDef for the spacefile
// from a API Cronjob Declaration
func CronjobFromDeclaration(decl *spaces.Cronjob) CronjobDef {
	def := CronjobDef{
		Identifier:    decl.ID,
		Schedule:      decl.Schedule,
		AllowParallel: decl.AllowParallel,
		Timezone:      decl.Timezone,
	}

	job := CommandCronjobDef{
		Command:          decl.Job.Command,
		Arguments:        decl.Job.Arguments,
		WorkingDirectory: decl.Job.WorkingDirectory,
	}

	def.Command = &job
	return def
}
