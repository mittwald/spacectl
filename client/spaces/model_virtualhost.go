package spaces

type VirtualHost struct {
	Hostname string         `json:"hostname"`
	TLS      VirtualHostTLS `json:"tls"`
}

type VirtualHostTLS struct {
	Type        string                 `json:"type"`
	Certificate VirtualHostCertificate `json:"certificate"`
}

type VirtualHostCertificate struct {
	ID string `json:"id"`
}

type VirtualHostList []VirtualHost

func (l VirtualHostList) Exists(hostname string) bool {
	for _, host := range l {
		if host.Hostname == hostname {
			return true
		}
	}
	return false
}
