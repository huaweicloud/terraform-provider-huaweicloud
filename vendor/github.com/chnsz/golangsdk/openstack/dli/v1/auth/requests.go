package auth

import (
	"github.com/chnsz/golangsdk"
)

type GrantDataPermissionOpts struct {
	UserName string `json:"user_name" required:"true"`
	// Grants or revokes the permission. The parameter value can be grant, revoke, or update.
	// grant: Indicates to grant users with permissions.
	// revoke: Indicates to revoke permissions.
	// update: Indicates to clear all the original permissions and assign the permissions in the provided
	// permission array.
	// NOTE:
	// Users can perform the update operation only when they have been granted with the grant and revoke permissions.
	Action string `json:"action" required:"true"`
	// Permission granting information. For details, see Table 3.
	Privileges []DataPermission `json:"privileges" required:"true"`
}

type DataPermission struct {
	// Data objects to be assigned. If they are named:
	// databases.Database_name, data in the entire database will be shared.
	// databases.Database_name.tables.Table_name, data in the specified table will be shared.
	// databases.Database_name.tables.Table_name.columns.Column_name, data in the specified column will be shared.
	// jobs.flink.Flink job ID, data the specified job will be shared.
	// groups. Package group name, data in the specified package group will be shared.
	// resources. Package name, data in the specified package will be shared.
	Object string `json:"object" required:"true"`
	// List of permissions to be granted, revoked, or updated.
	// NOTE:
	// If Action is Update and the update list is empty, all permissions of the user in the database or table
	// are revoked.
	Privileges []string `json:"privileges" required:"true"`
}

type ListDataPermissionOpts struct {
	// Data object to be assigned, which corresponds to the object in API permission assignment.
	// "jobs.flink.job_ID", data in the specified job will be queried.
	// "groups.Package_group_name", data in the specified package group will be queried.
	// "resources.Package_ame", data in the specified package will be queried.
	// NOTE:
	// When you view the packages in a group, the object format is resources.package group name/package name.
	Object string `q:"object"`
	Limit  int    `q:"limit"`
	Offset int    `q:"offset"`
}

type GrantQueuePermissionOpts struct {
	QueueName string `json:"queue_name" required:"true"`
	UserName  string `json:"user_name" required:"true"`
	// Grants or revokes the permission. The parameter value can be grant, revoke, or update.
	// grant: Indicates to grant users with permissions.
	// revoke: Indicates to revoke permissions.
	// update: Indicates to clear all the original permissions and assign the permissions in the provided
	// permission array.
	Action string `json:"action" required:"true"`
	// List of permissions to be granted, revoked, or updated. The following permissions are supported:
	// SUBMIT_JOB: indicates to submit a job.
	// CANCEL_JOB: indicates to cancel a job.
	// DROP_QUEUE: indicates to a delete a queue.
	// GRANT_PRIVILEGE: indicates to assign a permission.
	// REVOKE_PRIVILEGE: indicates to revoke a permission.
	// SHOW_PRIVILEGE: indicates to view other user's permissions.
	// RESTART: indicates to restart the queue.
	// SCALE_QUEUE: indicates to change the queue specifications.
	// NOTE:
	// If the update list is empty, all permissions of the queue granted to the user are revoked.
	Privileges []string `json:"privileges" required:"true"`
}

var RequestOpts = golangsdk.RequestOpts{
	MoreHeaders: map[string]string{"Content-Type": "application/json", "X-Language": "en-us"},
}

func GrantDataPermission(c *golangsdk.ServiceClient, opts GrantDataPermissionOpts) (*CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CommonResp
	_, err = c.Put(grantDataPermissionURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ListDataPermission(c *golangsdk.ServiceClient, opts ListDataPermissionOpts) (*DataPermissions, error) {
	url := ListDataPermissionUrl(c)
	query, err := golangsdk.BuildQueryString(opts)
	if err != nil {
		return nil, err
	}
	url += query.String()

	var rst DataPermissions
	_, err = c.Get(url, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func GrantQueuePermission(c *golangsdk.ServiceClient, opts GrantQueuePermissionOpts) (*CommonResp, error) {
	b, err := golangsdk.BuildRequestBody(opts, "")
	if err != nil {
		return nil, err
	}

	var rst CommonResp
	_, err = c.Put(grantQueuePermissionURL(c), b, &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ListQueuePermission(c *golangsdk.ServiceClient, queueName string) (*QueuePermissions, error) {
	var rst QueuePermissions
	_, err := c.Get(listQueuePermissionURL(c, queueName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ListDatabasePermission(c *golangsdk.ServiceClient, databaseName string) (*DatabasePermissions, error) {
	var rst DatabasePermissions
	_, err := c.Get(listDatabasePermissionURL(c, databaseName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func ListTablePermission(c *golangsdk.ServiceClient, databaseName string, tableName string) (*TablePermissions, error) {
	var rst TablePermissions
	_, err := c.Get(listTablePermissionURL(c, databaseName, tableName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}

func GetTablePermissionOfUser(c *golangsdk.ServiceClient, dbName string, tableName string, userName string) (
	*TablePermissionsOfUser, error) {
	var rst TablePermissionsOfUser
	_, err := c.Get(getTablePermissionOfUserURL(c, dbName, tableName, userName), &rst, &golangsdk.RequestOpts{
		MoreHeaders: RequestOpts.MoreHeaders,
	})
	return &rst, err
}
