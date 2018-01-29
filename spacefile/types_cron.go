package spacefile

type CronjobDef struct {
	Identifier    string             `hcl:",key"`
	Schedule      string             `hcl:"schedule"`
	AllowParallel bool               `hcl:"allowParallel"`
	Command       *CommandCronjobDef `hcl:"command"`
}

type CommandCronjobDef struct {
	Command          string   `hcl:"command"`
	Arguments        []string `hcl:"arguments"`
	WorkingDirectory string   `hcl:"workingDirectory"`
}

type CronjobDefList []CronjobDef
