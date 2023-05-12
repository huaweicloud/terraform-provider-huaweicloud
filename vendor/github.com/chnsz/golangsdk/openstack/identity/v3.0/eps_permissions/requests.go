package eps_permissions

import (
	"github.com/chnsz/golangsdk"
)

// UserGroupPermissionsCreate assigns enterprise project, user group and role.
func UserGroupPermissionsCreate(client *golangsdk.ServiceClient, enterpriseProjectID, groupID, roleID string) (r CommonResult) {
	_, r.Err = client.Put(userGroupPermissionsURL(client, enterpriseProjectID, groupID, roleID), nil, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// UserGroupPermissionsCreate unassigns enterprise project, user group and role.
func UserGroupPermissionsDelete(client *golangsdk.ServiceClient, enterpriseProjectID, groupID, roleID string) (r CommonResult) {
	_, r.Err = client.Delete(userGroupPermissionsURL(client, enterpriseProjectID, groupID, roleID), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// UserGroupPermissionsGet gets roles assigned to specified enterprise project and user group.
func UserGroupPermissionsGet(client *golangsdk.ServiceClient, enterpriseProjectID, groupID string) (r RoleResult) {
	_, r.Err = client.Get(userGroupPermissionsGetURL(client, enterpriseProjectID, groupID), &r.Body, nil)
	return
}

// UserPermissionsCreate assigns enterprise project, user and role.
func UserPermissionsCreate(client *golangsdk.ServiceClient, enterpriseProjectID, userID, roleID string) (r CommonResult) {
	_, r.Err = client.Put(userPermissionsURL(client, enterpriseProjectID, userID, roleID), nil, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// UserPermissionsDelete unassigns enterprise project, group and role.
func UserPermissionsDelete(client *golangsdk.ServiceClient, enterpriseProjectID, userID, roleID string) (r CommonResult) {
	_, r.Err = client.Delete(userPermissionsURL(client, enterpriseProjectID, userID, roleID), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// UserPermissionsGet gets roles assigned to specified enterprise project and user.
func UserPermissionsGet(client *golangsdk.ServiceClient, enterpriseProjectID, userID string) (r RoleResult) {
	_, r.Err = client.Get(userPermissionsGetURL(client, enterpriseProjectID, userID), &r.Body, nil)
	return
}

type AgencyPermissionsOpts struct {
	RoleAssignments []RoleAssignment `json:"role_assignments" required:"true"`
}

type RoleAssignment struct {
	AgencyID            string `json:"agency_id" required:"true"`
	EnterprisePorjectID string `json:"enterprise_project_id" required:"true"`
	RoleID              string `json:"role_id" required:"true"`
}

// AgencyPermissionsCreate assigns enterprise project, agency and role.
func AgencyPermissionsCreate(client *golangsdk.ServiceClient, opts *AgencyPermissionsOpts) (r CommonResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Put(agencyPermissionsURL(client), &b, nil, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// AgencyPermissionsDelete unassigns enterprise project, agency and role.
func AgencyPermissionsDelete(client *golangsdk.ServiceClient, opts *AgencyPermissionsOpts) (r CommonResult) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	_, r.Err = client.Delete(agencyPermissionsURL(client), &golangsdk.RequestOpts{
		JSONBody: b,
		OkCodes:  []int{204},
	})
	return
}
