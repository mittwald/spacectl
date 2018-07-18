package spaces

import "fmt"

func (c *spacesClient) ListCaughtEmails(spaceID string) (CaughtEmailList, error) {
	var space Space
	var emails []CaughtEmail

	err := c.client.Get("/spaces/"+spaceID, &space)
	if err != nil {
		return nil, fmt.Errorf("could not load stage: %s", err)
	}

	link, err := space.Links.GetLinkByRel("caughtEmails")
	if err != nil {
		return nil, fmt.Errorf("could not access caught emails: %s", err)
	}

	err = link.Get(c.client, &emails)
	if err != nil {
		return nil, fmt.Errorf("could not load caught emails: %s", err)
	}

	return emails, nil
}
