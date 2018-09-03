package spacefile

import "github.com/mittwald/spacectl/client/spaces"

type VirtualHostDef struct {
	Hostname string             `hcl:",key"`
	TLS      *TLSVirtualHostDef `hcl:"tls"`
}

type VirtualHostDefList []VirtualHostDef

type TLSVirtualHostDef struct {
	Type        string `hcl:"type"`
	Certificate string `hcl:"certificateID"`
}

func (v VirtualHostDef) ToDeclaration() spaces.VirtualHost {
	d := spaces.VirtualHost{
		Hostname: v.Hostname,
	}

	if v.TLS != nil {
		d.TLS.Type = v.TLS.Type
		d.TLS.Certificate.ID = v.TLS.Certificate
	} else {
		d.TLS.Type = "none"
	}
	return d
}
