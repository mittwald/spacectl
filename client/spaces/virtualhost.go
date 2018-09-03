package spaces

import (
	"fmt"

	"github.com/mittwald/spacectl/client/errors"
)

// ListVirtualHostsByStage returns an array of VirtualHosts for the stage
func (c *spacesClient) ListVirtualHostsByStage(spaceID, stage string) (VirtualHostList, error) {
	var existingVHosts VirtualHostList

	err := c.client.Get("/spaces/"+spaceID+"/stages/"+stage+"/virtualhosts", &existingVHosts)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not access virtualhosts for space: %s, stage: %s", spaceID, stage)}
	}

	return existingVHosts, nil
}

// UpdateVirtualHost updates or creates the VirtualHost from the given model
func (c *spacesClient) UpdateVirtualHost(spaceID, stage string, vhost VirtualHost) (*VirtualHost, error) {
	hosts, err := c.ListVirtualHostsByStage(spaceID, stage)
	if err != nil {
		return nil, err
	}

	var newVirtualHost VirtualHost
	if hosts.Exists(vhost.Hostname) {
		err = c.client.Put("/spaces/"+spaceID+"/stages/"+stage+"/virtualhosts/"+vhost.Hostname, vhost, &newVirtualHost)
	} else {
		err = c.client.Post("/spaces/"+spaceID+"/stages/"+stage+"/virtualhosts", vhost, &newVirtualHost)
	}
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not create/update virtualHost %s", vhost.Hostname)}
	}

	return &newVirtualHost, err
}

// DeleteVirtualHost deletes the VirtualHost with the given hostname
func (c *spacesClient) DeleteVirtualHost(spaceID, stage, hostname string) error {
	var target Empty
	err := c.client.Delete("/spaces/"+spaceID+"/stages/"+stage+"/virtualhosts/"+hostname, &target)
	if err != nil {
		return errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not delete virtualHost %s", hostname)}
	}
	return nil
}
