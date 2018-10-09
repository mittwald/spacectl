package spacefile

import (
	"fmt"

	"github.com/hashicorp/go-multierror"
)

type StageDef struct {
	Name         string             `hcl:",key"`
	Inherit      string             `hcl:"inherit" hcle:"omitempty"`
	Applications SoftwareDefList    `hcl:"application"`
	Databases    SoftwareDefList    `hcl:"database"`
	Cronjobs     CronjobDefList     `hcl:"cron"`
	VirtualHosts VirtualHostDefList `hcl:"virtualHost"`
	inheritStage *StageDef
}

func (d *StageDef) Validate(offline bool) error {
	var err *multierror.Error

	if len(d.Applications) > 1 {
		err = multierror.Append(err, fmt.Errorf("stage '%s' should not contain more than one application", d.Name))
	}

	if d.Application() != nil {
		err = multierror.Append(err, d.Application().Validate(offline))
	}

	for _, vhost := range d.VirtualHosts {
		if vhost.TLS == nil {
			continue
		}
		if vhost.TLS.Type == "certificate" {
			if vhost.TLS.Certificate == "" {
				err = multierror.Append(err, fmt.Errorf("no certificate given for vhost %s", vhost.Hostname))
			}
		}
	}

	return err.ErrorOrNil()
}

// Application returns a reference to the one application defined for this stage
func (d *StageDef) Application() *SoftwareDef {
	for i := range d.Applications {
		app := d.Applications[i]
		return &app
	}

	return nil
}

func (d *StageDef) resolveUserData() error {
	var mErr *multierror.Error
	var err error

	for i := range d.Applications {
		d.Applications[i].UserData, err = unfuckHCL(d.Applications[i].UserData, "")
		mErr = multierror.Append(mErr, err)

		if d.Applications[i].UserData == nil {
			d.Applications[i].UserData = map[string]string{}
		}
	}

	return mErr.ErrorOrNil()
}

func (d *StageDef) resolveInheritance(level int) error {
	if level > 4 {
		return fmt.Errorf("could not resolve stage dependencies after %d levels. Please check that there is no cyclic inheritance", level)
	}

	if d.inheritStage == nil {
		return nil
	}

	err := d.inheritStage.resolveInheritance(level + 1)
	if err != nil {
		return err
	}

	originalName := d.Name

	d.Applications, err = d.Applications.Merge(d.inheritStage.Applications)
	if err != nil {
		return err
	}

	d.Databases, err = d.Databases.Merge(d.inheritStage.Databases)
	if err != nil {
		return err
	}

	d.Name = originalName
	d.inheritStage = nil

	return nil
}
