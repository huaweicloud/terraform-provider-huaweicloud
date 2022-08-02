package dds

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/dds/v3/users"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceDatabaseUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabaseUserCreate,
		ReadContext:   resourceDatabaseUserRead,
		UpdateContext: resourceDatabaseUserUpdate,
		DeleteContext: resourceDatabaseUserDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(2 * time.Minute),
			Delete: schema.DefaultTimeout(2 * time.Minute),
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
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[\w-.]*$`),
						"The name can only contain letters, digits, underscores (_), hyphens (-) and dots (.)."),
					validation.StringDoesNotMatch(regexp.MustCompile(`^drsFull$|^drsIncremental$`),
						"Cannot use reserved names: 'drsFull' or 'drsIncremental'"),
					validation.StringLenBetween(1, 64),
				),
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`[A-Z]`),
						"Missing uppercase character, the name must contain uppercase and lowercase letters, digits"+
							" and special characters (~!@#%^*-_=+?)."),
					validation.StringMatch(regexp.MustCompile(`[a-z]`),
						"Missing lowercase character, the name must contain uppercase and lowercase letters, digits"+
							" and special characters (~!@#%^*-_=+?)."),
					validation.StringMatch(regexp.MustCompile(`[0-9]`),
						"Missing digit, the name must contain uppercase and lowercase letters, digits and special "+
							"characters (~!@#%^*-_=+?)."),
					validation.StringMatch(regexp.MustCompile(`[~!@#%^-_=+?]|\*`),
						"Missing special character, the name must contain uppercase and lowercase letters, digits and"+
							" special characters (~!@#%^*-_=+?)."),
					validation.StringLenBetween(8, 32),
				),
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^\w*$`),
						"The name can only contain letters, digits and underscores (_)."),
					validation.StringLenBetween(1, 64),
				),
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
							ValidateFunc: validation.All(
								validation.StringMatch(regexp.MustCompile(`^\w*$`),
									"The name can only contain letters, digits and underscores (_)."),
								validation.StringLenBetween(1, 64),
							),
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
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	// Before creating database user, we need to ensure that the database is not currently performing operations.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceActionsRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutCreate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the action of DDS instance to complete: %s", err)
	}

	dbName := d.Get("db_name").(string)
	userName := d.Get("name").(string)
	opts := users.CreateOpts{
		Password: d.Get("password").(string),
		Name:     userName,
		DbName:   dbName,
		Roles:    buildDatabaseRoles(d.Get("roles").([]interface{})),
	}
	err = users.Create(client, instanceId, opts)
	if err != nil {
		return diag.Errorf("error creating database user (%s): %v", userName, err)
	}
	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, dbName, userName))
	return resourceDatabaseUserRead(ctx, d, meta)
}

func resourceDatabaseUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to find user (%s) from DDS instance (%s)", name, instanceId))
	}
	user := resp[0]
	tflog.Debug(ctx, fmt.Sprintf("The user response is: %#v", user))

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
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	// Before updating database user, we need to ensure that the database is not currently performing operations.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceActionsRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutUpdate),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the action of DDS instance to complete: %s", err)
	}

	name := d.Get("name").(string)
	opts := users.PwdResetOpts{
		Name:     name,
		Password: d.Get("password").(string),
		DbName:   d.Get("db_name").(string),
	}
	err = users.ResetPassword(client, instanceId, opts)
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
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	// Before deleting database user, we need to ensure that the database is not currently performing operations.
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"ACTIVE"},
		Refresh:      instanceActionsRefreshFunc(client, instanceId),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        1 * time.Second,
		PollInterval: 5 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.Errorf("error waiting for the action of DDS instance to complete: %s", err)
	}

	dbName := d.Get("db_name").(string)
	userName := d.Get("name").(string)
	opts := users.DeleteOpts{
		Name:   userName,
		DbName: dbName,
	}
	err = users.Delete(client, instanceId, opts)
	if err != nil {
		return diag.Errorf("error deleting database user (%s) from instance (%s): %v", userName, instanceId, err)
	}

	stateConf = &resource.StateChangeConf{
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
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<db_name>/<name>")
	}

	d.Set("instance_id", parts[0])
	d.Set("db_name", parts[1])
	d.Set("name", parts[2])
	return []*schema.ResourceData{d}, nil
}
