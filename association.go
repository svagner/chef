package chef

// import "fmt"
import "errors"

type AssociationService struct {
	client *Client
}

// Chef API docs: https://docs.chef.io/api_chef_server.html#association-requests
// https://github.com/chef/chef-server/blob/master/src/oc_erchef/apps/oc_chef_wm/src/oc_chef_wm_org_invites.erl  Invitation implementation
// https://github.com/chef/chef-server/blob/master/src/oc_erchef/apps/oc_chef_wm/src/oc_chef_wm_org_associations.erl user org associations
type Association struct {
	Uri              string `json:"uri"` // the last part of the uri is the invitation id
	OrganizationUser struct {
		UserName string `json:"username,omitempty"`
	} `json:"organization_user"`
	Organization struct {
		Name string `json:"name,omitempty"`
	} `json:"organization"`
	User struct {
		Email     string `json:"email,omitempty"`
		FirstName string `json:"first_name,omitempty"`
	} `json:"user"`
}

type OrgUserListEntry struct {
	User struct {
		Username string `json:"username,omitempty"`
	} `json:"user,omitempty"`
}

type OrgUser struct {
	Username    string `json:"username,omitempty"`
	Email       string `json:"email,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	PublicKey   string `json:"public_key,omitempty"`
}

type RescindInvite struct {
	Id       string `json:"id,omitempty"`
	Orgname  string `json:"orgname,omitempty"`
	Username string `json:"username,omitempty"`
}

type Invite struct {
	Id       string `json:"id,omitempty"`
	UserName string `json:"username,omitempty"`
}

type Request struct {
	User string `json:"user"`
}

type AddNow struct {
	Username string `json:"username"`
}

// Get gets a list of the pending invitations for an organization.
func (e *AssociationService) ListInvites() (invitelist []Invite, err error) {
	err = e.client.magicRequestDecoder("GET", "association_requests", nil, &invitelist)
	return
}

// Creates an invitation for a user to join an organization on the chef server
func (e *AssociationService) Invite(invite Request) (data Association, err error) {
	body, err := JSONReader(invite)
	if err != nil {
		return
	}
	err = e.client.magicRequestDecoder("POST", "association_requests/", body, &data)
	return
}

// Delete removes a pending invitation to an organization
func (e *AssociationService) DeleteInvite(id string) (rescind RescindInvite, err error) {
	err = e.client.magicRequestDecoder("DELETE", "association_requests/"+id, nil, &rescind)
	return
}

// InviteID Finds an invitation id for a user
func (e *AssociationService) InviteId(user string) (id string, err error) {
	var invitelist []Invite
	err = e.client.magicRequestDecoder("GET", "association_requests", nil, &invitelist)
	if err != nil {
		return
	}
	// Find an invite for the user or return err
	for _, in := range invitelist {
		if in.UserName == user {
			id = in.Id
		}
	}
	if id == "" {
		err = errors.New("User request not found")
	}
	return
}

// AcceptInvite Accepts an invitation
// TODO: Gets a 405, code is in knife is it part of erchef?
func (e *AssociationService) AcceptInvite(id string) (data string, err error) {
	body, err := JSONReader("{ \"accept\" }")
	if err != nil {
		return
	}
	err = e.client.magicRequestDecoder("PUT", "association_requests/"+id, body, &data)
	return
}

// Get a list of the users in an organization
func (e *AssociationService) List() (data []OrgUserListEntry, err error) {
	err = e.client.magicRequestDecoder("GET", "users", nil, &data)
	return
}

// Add a user immediately
func (e *AssociationService) Add(addme AddNow) (err error) {
	body, err := JSONReader(addme)
	if err != nil {
		return
	}
	err = e.client.magicRequestDecoder("POST", "users", body, nil)
	return
}

// Get the details of a user in an organization
func (e *AssociationService) Get(name string) (data OrgUser, err error) {
	err = e.client.magicRequestDecoder("GET", "users/"+name, nil, &data)
	return
}

// Delete removes a user from an organization
func (e *AssociationService) Delete(name string) (data OrgUser, err error) {
	err = e.client.magicRequestDecoder("DELETE", "users/"+name, nil, &data)
	return
}
