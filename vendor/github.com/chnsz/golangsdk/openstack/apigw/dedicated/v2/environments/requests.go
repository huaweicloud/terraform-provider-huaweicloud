package environments

import (
	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"
)

// EnvironmentOpts allows to create a new environment or update an existing environment using given parameters.
type EnvironmentOpts struct {
	// Environment name, which can contain 3 to 64 characters, starting with a letter.
	// Only letters, digits and underscores (_) are allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Name string `json:"name" required:"true"`
	// Description of the environment, which can contain a maximum of 255 characters,
	// and the angle brackets (< and >) are not allowed.
	// Chinese characters must be in UTF-8 or Unicode format.
	Description *string `json:"remark,omitempty"`
}

type EnvironmentOptsBuilder interface {
	ToEnvironmentOptsMap() (map[string]interface{}, error)
}

func (opts EnvironmentOpts) ToEnvironmentOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// Create is a method by which to create function that create a new environment.
func Create(client *golangsdk.ServiceClient, instanceId string, opts EnvironmentOptsBuilder) (r CreateResult) {
	reqBody, err := opts.ToEnvironmentOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId, "envs"), reqBody, &r.Body, nil)
	return
}

// Update is a method by which to udpate an existing environment.
func Update(client *golangsdk.ServiceClient, instanceId, envId string, opts EnvironmentOptsBuilder) (r UpdateResult) {
	reqBody, err := opts.ToEnvironmentOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Put(resourceURL(client, instanceId, "envs", envId), reqBody, &r.Body, &golangsdk.RequestOpts{
		OkCodes: []int{200},
	})
	return
}

// ListOpts allows to filter list data using given parameters.
type ListOpts struct {
	// Environment name.
	Name string `q:"name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
}

type ListOptsBuilder interface {
	ToListQuery() (string, error)
}

func (opts ListOpts) ToListQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// List is a method to obtain an array of one or more groups according to the query parameters.
func List(client *golangsdk.ServiceClient, instanceId string, opts ListOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId, "envs")
	if opts != nil {
		query, err := opts.ToListQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return EnvironmentPage{pagination.SinglePageBase(r)}
	})
}

// Delete is a method to delete an existing group.
func Delete(client *golangsdk.ServiceClient, instanceId, envId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, "envs", envId), nil)
	return
}

// CreateVariableOpts allows to create a new envirable variable for an existing group using given parameters.
type CreateVariableOpts struct {
	// Variable name, which can contain 3 to 32 characters, starting with a letter.
	// Only letters, digits, hyphens (-), and underscores (_) are allowed.
	// In the definition of an API, #Name# (case-sensitive) indicates a variable.
	// It is replaced by the actual value when the API is published in an environment.
	Name string `json:"variable_name" required:"true"`
	// Variable value, which can contain 1 to 255 characters.
	// Only letters, digits, and special characters (_-/.:) are allowed.
	Value string `json:"variable_value" required:"true"`
	// Environment ID.
	EnvId string `json:"env_id" required:"true"`
	// Group ID.
	GroupId string `json:"group_id" required:"true"`
}

type CreateVariableOptsBuilder interface {
	ToCreateVariableOptsMap() (map[string]interface{}, error)
}

func (opts CreateVariableOpts) ToCreateVariableOptsMap() (map[string]interface{}, error) {
	return golangsdk.BuildRequestBody(opts, "")
}

// CreateVariable is a method by which to create function that create a new environment variable.
func CreateVariable(client *golangsdk.ServiceClient, instanceId string,
	opts CreateVariableOptsBuilder) (r VariableCreateResult) {
	reqBody, err := opts.ToCreateVariableOptsMap()
	if err != nil {
		r.Err = err
		return
	}
	_, r.Err = client.Post(rootURL(client, instanceId, "env-variables"), reqBody, &r.Body, nil)
	return
}

// GetVariable is a method to obtain the specified environment variable according to the instance id and variable id.
func GetVariable(client *golangsdk.ServiceClient, instanceId, varId string) (r VariableGetResult) {
	_, r.Err = client.Get(resourceURL(client, instanceId, "env-variables", varId), &r.Body, nil)
	return
}

// ListVariablesOpts allows to filter list data using given parameters.
type ListVariablesOpts struct {
	// API group ID.
	GroupId string `q:"group_id"`
	// Environment ID.
	EnvId string `q:"env_id"`
	// Variable name.
	Name string `q:"variable_name"`
	// Offset from which the query starts.
	// If the offset is less than 0, the value is automatically converted to 0. Default to 0.
	Offset int `q:"offset"`
	// Number of items displayed on each page. The valid values are range form 1 to 500, default to 20.
	Limit int `q:"limit"`
	// Parameter name for exact matching. Only API variable names are supported.
	PreciseSearch string `q:"precise_search"`
}

type ListVariablesOptsBuilder interface {
	ToListVariablesQuery() (string, error)
}

func (opts ListVariablesOpts) ToListVariablesQuery() (string, error) {
	q, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return "", err
	}
	return q.String(), err
}

// ListVariables is a method to obtain an array of one or more variables according to the query parameters.
func ListVariables(client *golangsdk.ServiceClient, instanceId string, opts ListVariablesOptsBuilder) pagination.Pager {
	url := rootURL(client, instanceId, "env-variables")
	if opts != nil {
		query, err := opts.ToListVariablesQuery()
		if err != nil {
			return pagination.Pager{Err: err}
		}
		url += query
	}

	return pagination.NewPager(client, url, func(r pagination.PageResult) pagination.Page {
		return VariablePage{pagination.SinglePageBase(r)}
	})
}

// DeleteVariable is a method to delete an existing variable.
func DeleteVariable(client *golangsdk.ServiceClient, instanceId, varId string) (r DeleteResult) {
	_, r.Err = client.Delete(resourceURL(client, instanceId, "env-variables", varId), nil)
	return
}
