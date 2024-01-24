package dds

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dds/v3/users"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API DDS POST /v3/{project_id}/instances/{instance_id}/db-user
// @API DDS GET /v3/{project_id}/instances
// @API DDS GET /v3/{project_id}/instances/{instance_id}/db-user/detail
// @API DDS PUT /v3/{project_id}/instances/{instance_id}/reset-password
// @API DDS DELETE /v3/{project_id}/instances/{instance_id}/db-user
func ResourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseUserCreate,
		ReadContext:   resourceDatabaseUserRead,
		UpdateContext: resourceDatabaseUserUpdate,
		DeleteContext: resourceDatabaseUserDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(60 * time.Minute),
			Update: schema.DefaultTimeout(60 * time.Minute),
			Delete: schema.DefaultTimeout(60 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceDatabaseUserImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"roles": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"db_name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
					},
				},
			},
			"privileges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privilegeSchemaResource(),
			},
			"inherited_privileges": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     privilegeSchemaResource(),
			},
		},
	}
}

func resourceDatabaseUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	userName := d.Get("name").(string)
	opts := users.CreateOpts{
		Password: d.Get("password").(string),
		Name:     userName,
		DbName:   dbName,
		Roles:    buildDatabaseRoles(d.Get("roles").([]interface{})),
	}
	retryFunc := func() (interface{}, bool, error) {
		err = users.Create(client, instanceId, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceActionsRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutCreate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error creating database user (%s): %v", userName, err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, dbName, userName))
	return resourceDatabaseUserRead(ctx, d, meta)
}

func resourceDatabaseUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.DdsV3Client(region)
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	opts := users.ListOpts{
		DbName: d.Get("db_name").(string),
		Name:   name,
	}
	resp, err := users.List(client, instanceId, opts)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting user (%s) from DDS instance (%s)", name, instanceId))
	}
	if len(resp) < 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}
	user := resp[0]
	log.Printf("[DEBUG] The user response is: %#v", user)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", user.Name),
		d.Set("db_name", user.DbName),
		d.Set("roles", flattenDatabaseRoles(user.Roles)),
		d.Set("privileges", flattenDatabasePrivileges(user.Privileges)),
		d.Set("inherited_privileges", flattenDatabasePrivileges(user.InheritedPrivileges)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDatabaseUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	name := d.Get("name").(string)
	opts := users.PwdResetOpts{
		Name:     name,
		Password: d.Get("password").(string),
		DbName:   d.Get("db_name").(string),
	}
	retryFunc := func() (interface{}, bool, error) {
		err = users.ResetPassword(client, instanceId, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceActionsRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error resetting the password of the database user (%s): %v", name, err)
	}

	return resourceDatabaseUserRead(ctx, d, meta)
}

func databaseUserRefreshFunc(client *golangsdk.ServiceClient, instanceId, dbName,
	userName string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		opts := users.ListOpts{
			DbName: dbName,
			Name:   userName,
		}
		resp, err := users.List(client, instanceId, opts)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error getting database user (%s) from DDS instance (%s)", userName,
				instanceId)
		}
		if len(resp) < 1 {
			return resp, "DELETED", nil
		}
		return resp, "ACTIVE", nil
	}
}

func resourceDatabaseUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.DdsV3Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DDS v3 client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	userName := d.Get("name").(string)
	opts := users.DeleteOpts{
		Name:   userName,
		DbName: dbName,
	}
	retryFunc := func() (interface{}, bool, error) {
		err = users.Delete(client, instanceId, opts)
		retry, err := handleMultiOperationsError(err)
		return nil, retry, err
	}
	_, err = common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    retryFunc,
		WaitFunc:     instanceActionsRefreshFunc(client, instanceId),
		WaitTarget:   []string{"ACTIVE"},
		WaitPending:  []string{"PENDING"},
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 1 * time.Second,
		PollInterval: 10 * time.Second,
	})
	if err != nil {
		return diag.Errorf("error deleting database user (%s) from instance (%s): %v", userName, instanceId, err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"ACTIVE"},
		Target:       []string{"DELETED"},
		Refresh:      databaseUserRefreshFunc(client, instanceId, dbName, userName),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        2 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for database user deletion to complete: %s", err)
	}

	return nil
}

func resourceDatabaseUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<db_name>/<name>")
	}

	d.Set("instance_id", parts[0])
	d.Set("db_name", parts[1])
	d.Set("name", parts[2])
	return []*schema.ResourceData{d}, nil
}
