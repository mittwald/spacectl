package teams

type Team struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	DNSName string `json:"dnsLabel"`
}

type Invite struct {
	ID      string `json:"id"`
	Message string `json:"message"`
	State   string `json:"state"`
}

type MembershipUser struct {
	ID              string `json:"id"`
	FirstName       string `json:"firstName"`
	LastName        string `json:"lastName"`
	Email           string `json:"email"`
	ProfileImageURI string `json:"profileImageURI"`
}

type Membership struct {
	User MembershipUser `json:"user"`
	Role string         `json:"role"`
}
