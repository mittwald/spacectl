package spaces

import (
	"fmt"
	"net/url"

	"github.com/mittwald/spacectl/client/errors"
)

// ListVirtualHostsByStage returns an array of VirtualHosts for the stage
func (c *spacesClient) ListVirtualHostsByStage(spaceID, stage string) (VirtualHostList, error) {
	var existingVHosts VirtualHostList

	listPath := fmt.Sprintf("/spaces/%s/stages/%s/virtualhosts", url.PathEscape(spaceID), url.PathEscape(stage))
	err := c.client.Get(listPath, &existingVHosts)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not access virtualhosts for space: %s, stage: %s", spaceID, stage)}
	}

	return existingVHosts, nil
}

// UpdateVirtualHost updates or creates the VirtualHost from the given model
func (c *spacesClient) UpdateVirtualHost(spaceID, stage string, vhost VirtualHost) (*VirtualHost, error) {
	var newVirtualHost VirtualHost
	createPath := fmt.Sprintf("/spaces/%s/stages/%s/virtualhosts", url.PathEscape(spaceID), url.PathEscape(stage))
	err := c.client.Put(createPath+"/"+url.PathEscape(vhost.Hostname), vhost, &newVirtualHost)
	if err != nil {
		return nil, errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not create/update virtualHost %s", vhost.Hostname)}
	}

	return &newVirtualHost, err
}

// DeleteVirtualHost deletes the VirtualHost with the given hostname
func (c *spacesClient) DeleteVirtualHost(spaceID, stage, hostname string) error {
	var target Empty
	deletePath := fmt.Sprintf("/spaces/%s/stages/%s/virtualhosts/%s", url.PathEscape(spaceID), url.PathEscape(stage), url.PathEscape(hostname))
	err := c.client.Delete(deletePath, &target)
	if err != nil {
		return errors.ErrNested{Inner: err, Msg: fmt.Sprintf("could not delete virtualHost %s", hostname)}
	}
	return nil
}
