package spaces

import "time"

type Cronjob struct {
	ID            string     `json:"id,omitempty"`
	Schedule      string     `json:"schedule"`
	AllowParallel bool       `json:"allowParallel"`
	Job           CronjobJob `json:"job"`
	ReadOnly      bool `json:"readonly"`
	Timezone      string     `json:"timezone"`
}

type CronjobJob struct {
	Type             string                   `json:"type"`
	Command          string                   `json:"command,omitempty"`
	Arguments        []string                 `json:"arguments,omitempty"`
	WorkingDirectory string                   `json:"workingDirectory,omitempty"`
	NextExecution    *CommandCronjobExecution `json:"nextExecution,omitempty"`
	LastExecution    *CommandCronjobExecution `json:"lastExecution,omitempty"`
}

type CommandCronjobExecution struct {
	Date     time.Time `json:"date"`
	ExitCode int       `json:"exitCode,omitempty"`
}
