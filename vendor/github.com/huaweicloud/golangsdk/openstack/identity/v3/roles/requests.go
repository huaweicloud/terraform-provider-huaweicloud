package roles

import (
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/pagination"
)

// ListOptsBuilder allows extensions to add additional parameters to
// the List request
type ListOptsBuilder interface {
	ToRoleListQuery() (string, error)
}

// ListOpts provides options to filter the List results.
type ListOpts struct {
	// DomainID filters the response by a domain ID.
	DomainID string `q:"domain_id"`

	// Name filters the response by role name.
	Name string `q:"name"`
}

// ToRoleListQuery formats a ListOpts into a query string.
func (opts ListOpts) ToRoleListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	return q.String(), err
}

// List enumerates the roles to which the current token has access.
func List(client *golangsdk.ServiceClient, opts ListOptsBuilder) pagination.Pager {
	url := listURL(client)
	if opts != nil {
		query, err := opts.ToRoleListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RolePage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// Get retrieves details on a single role, by ID.
func Get(client *golangsdk.ServiceClient, id string) (r GetResult) {
	_, r.Err = client.Get(getURL(client, id), &r.Body, nil)
	return
}

// CreateOptsBuilder allows extensions to add additional parameters to
// the Create request.
type CreateOptsBuilder interface {
	ToRoleCreateMap() (map[string]interface{}, error)
}

// CreateOpts provides options used to create a role.
type CreateOpts struct {
	// Name is the name of the new role.
	Name string `json:"name" required:"true"`

	// DomainID is the ID of the domain the role belongs to.
	DomainID string `json:"domain_id,omitempty"`

	// Extra is free-form extra key/value pairs to describe the role.
	Extra map[string]interface{} `json:"-"`
}

// ToRoleCreateMap formats a CreateOpts into a create request.
func (opts CreateOpts) ToRoleCreateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "role")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["role"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Create creates a new Role.
func Create(client *golangsdk.ServiceClient, opts CreateOptsBuilder) (r CreateResult) {
	b, err := opts.ToRoleCreateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(createURL(client), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{201},
	})
	return
}

// UpdateOptsBuilder allows extensions to add additional parameters to
// the Update request.
type UpdateOptsBuilder interface {
	ToRoleUpdateMap() (map[string]interface{}, error)
}

// UpdateOpts provides options for updating a role.
type UpdateOpts struct {
	// Name is the name of the new role.
	Name string `json:"name,omitempty"`

	// Extra is free-form extra key/value pairs to describe the role.
	Extra map[string]interface{} `json:"-"`
}

// ToRoleUpdateMap formats a UpdateOpts into an update request.
func (opts UpdateOpts) ToRoleUpdateMap() (map[string]interface{}, error) {
	b, err := golangsdk.BuildRequestBody(opts, "role")
	if err != nil {
		return nil, err
	}

	if opts.Extra != nil {
		if v, ok := b["role"].(map[string]interface{}); ok {
			for key, value := range opts.Extra {
				v[key] = value
			}
		}
	}

	return b, nil
}

// Update updates an existing Role.
func Update(client *golangsdk.ServiceClient, roleID string, opts UpdateOptsBuilder) (r UpdateResult) {
	b, err := opts.ToRoleUpdateMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Patch(updateURL(client, roleID), &b, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// Delete deletes a role.
func Delete(client *golangsdk.ServiceClient, roleID string) (r DeleteResult) {
	_, r.Err = client.Delete(deleteURL(client, roleID), nil)
	return
}

// ListAssignmentsOptsBuilder allows extensions to add additional parameters to
// the ListAssignments request.
type ListAssignmentsOptsBuilder interface {
	extractAssignment() (string, string, string, string, error)
}

// ListAssignmentsOpts allows you to query the ListAssignments method.
// Specify one of or a combination of GroupId, RoleId, ScopeDomainId,
// ScopeProjectId, and/or UserId to search for roles assigned to corresponding
// entities.
type ListAssignmentsOpts struct {
	// GroupID is the group ID to query.
	GroupID string `q:"group.id"`

	// ScopeDomainID filters the results by the given domain ID.
	ScopeDomainID string `q:"scope.domain.id"`

	// ScopeProjectID filters the results by the given Project ID.
	ScopeProjectID string `q:"scope.project.id"`

	// UserID filterst he results by the given User ID.
	UserID string `q:"user.id"`
}

// ToRolesListAssignmentsQuery formats a ListAssignmentsOpts into a query string.
func (opts ListAssignmentsOpts) extractAssignment() (string, string, string, string, error) {
	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ScopeProjectID != "" {
		targetID = opts.ScopeProjectID
		targetType = "projects"
	} else {
		targetID = opts.ScopeDomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	return targetType, targetID, actorType, actorID, nil
}

// ListAssignments enumerates the roles assigned to a specified resource.
func ListAssignments(client *golangsdk.ServiceClient, opts ListAssignmentsOptsBuilder) pagination.Pager {
	targetType, targetID, actorType, actorID, _ := opts.extractAssignment()

	url := listAssignmentsURL(client, targetType, targetID, actorType, actorID)
	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return RoleAssignmentPage{pagination.LinkedPageBase{PageResult: r}}
	})
}

// AssignOpts provides options to assign a role
type AssignOpts struct {
	// UserID is the ID of a user to assign a role
	// Note: exactly one of UserID or GroupID must be provided
	UserID string `xor:"GroupID"`

	// GroupID is the ID of a group to assign a role
	// Note: exactly one of UserID or GroupID must be provided
	GroupID string `xor:"UserID"`

	// ProjectID is the ID of a project to assign a role on
	// Note: exactly one of ProjectID or DomainID must be provided
	ProjectID string `xor:"DomainID"`

	// DomainID is the ID of a domain to assign a role on
	// Note: exactly one of ProjectID or DomainID must be provided
	DomainID string `xor:"ProjectID"`
}

// UnassignOpts provides options to unassign a role
type UnassignOpts struct {
	// UserID is the ID of a user to unassign a role
	// Note: exactly one of UserID or GroupID must be provided
	UserID string `xor:"GroupID"`

	// GroupID is the ID of a group to unassign a role
	// Note: exactly one of UserID or GroupID must be provided
	GroupID string `xor:"UserID"`

	// ProjectID is the ID of a project to unassign a role on
	// Note: exactly one of ProjectID or DomainID must be provided
	ProjectID string `xor:"DomainID"`

	// DomainID is the ID of a domain to unassign a role on
	// Note: exactly one of ProjectID or DomainID must be provided
	DomainID string `xor:"ProjectID"`
}

// Assign is the operation responsible for assigning a role
// to a user/group on a project/domain.
func Assign(client *golangsdk.ServiceClient, roleID string, opts AssignOpts) (r AssignmentResult) {
	// Check xor conditions
	_, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ProjectID != "" {
		targetID = opts.ProjectID
		targetType = "projects"
	} else {
		targetID = opts.DomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	_, r.Err = client.Put(assignURL(client, targetType, targetID, actorType, actorID, roleID), nil, nil, &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}

// Unassign is the operation responsible for unassigning a role
// from a user/group on a project/domain.
func Unassign(client *golangsdk.ServiceClient, roleID string, opts UnassignOpts) (r UnassignmentResult) {
	// Check xor conditions
	_, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		r.Err = err
		return
	}

	// Get corresponding URL
	var targetID string
	var targetType string
	if opts.ProjectID != "" {
		targetID = opts.ProjectID
		targetType = "projects"
	} else {
		targetID = opts.DomainID
		targetType = "domains"
	}

	var actorID string
	var actorType string
	if opts.UserID != "" {
		actorID = opts.UserID
		actorType = "users"
	} else {
		actorID = opts.GroupID
		actorType = "groups"
	}

	_, r.Err = client.Delete(assignURL(client, targetType, targetID, actorType, actorID, roleID), &golangsdk.RequestOpts{
		OkCodes: []int{204},
	})
	return
}
