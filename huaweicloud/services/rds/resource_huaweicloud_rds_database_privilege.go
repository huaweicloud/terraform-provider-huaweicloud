package rds

import (
	"context"
	"log"
	"strings"

	"github.com/chnsz/golangsdk"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	v3 "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3"
	rds "github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func ResourceRdsDatabasePrivilege() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsDatabasePrivilegeCreate,
		DeleteContext: resourceRdsDatabasePrivilegeDelete,
		ReadContext:   resourceRdsDatabasePrivilegeRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"db_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"users": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 50,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
						},
						"readonly": {
							Type:     schema.TypeBool,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func buildUserOpts(rawUsers []interface{}) []rds.UserWithPrivilege {
	if len(rawUsers) < 1 {
		return nil
	}

	usersOpts := make([]rds.UserWithPrivilege, len(rawUsers))
	for i, v := range rawUsers {
		user := v.(map[string]interface{})
		usersOpts[i] = rds.UserWithPrivilege{
			Name:     user["name"].(string),
			Readonly: user["readonly"].(bool),
		}
	}
	return usersOpts
}

func buildRevokeUserOpts(rawUsers []interface{}) []rds.RevokeRequestBodyUsers {
	if len(rawUsers) < 1 {
		return nil
	}

	usersOpts := make([]rds.RevokeRequestBodyUsers, len(rawUsers))
	for i, v := range rawUsers {
		user := v.(map[string]interface{})
		usersOpts[i] = rds.RevokeRequestBodyUsers{
			Name: user["name"].(string),
		}
	}
	return usersOpts
}

func resourceRdsDatabasePrivilegeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	dbName := d.Get("db_name").(string)
	createOpts := rds.GrantRequest{
		DbName: d.Get("db_name").(string),
		Users:  buildUserOpts(d.Get("users").([]interface{})),
	}
	log.Printf("[DEBUG] Create RDS database privilege options: %#v", createOpts)

	privilegeReq := rds.AllowDbUserPrivilegeRequest{
		InstanceId: instanceId,
		Body:       &createOpts,
	}

	_, err = client.AllowDbUserPrivilege(&privilegeReq)
	if err != nil {
		return diag.Errorf("error creating RDS database privilege: %s", err)
	}

	id := instanceId + "/" + dbName
	d.SetId(id)
	return resourceRdsDatabasePrivilegeRead(ctx, d, meta)
}

func resourceRdsDatabasePrivilegeRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	// Split instance_id and database from resource id
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<database_name>")
	}
	instanceId := parts[0]
	dbName := parts[1]

	users, err := QueryDatabaseUsers(client, instanceId, dbName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error listing RDS db authorized users")
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", instanceId),
		d.Set("db_name", dbName),
		d.Set("users", flattenUsers(users)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting RDS db privilege fields: %s", err)
	}

	return nil
}

func resourceRdsDatabasePrivilegeDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	deleteOpts := rds.RevokeRequestBody{
		DbName: d.Get("db_name").(string),
		Users:  buildRevokeUserOpts(d.Get("users").([]interface{})),
	}
	log.Printf("[DEBUG] Delete RDS database privilege options: %#v", deleteOpts)

	deleteReq := rds.RevokeRequest{
		InstanceId: d.Get("instance_id").(string),
		Body:       &deleteOpts,
	}

	_, err = client.Revoke(&deleteReq)
	if err != nil {
		return diag.Errorf("error deleting RDS database privilege: %s", err)
	}

	return nil
}

func flattenUsers(users []rds.UserWithPrivilege) []map[string]interface{} {
	if len(users) < 1 {
		return nil
	}

	usersToSet := make([]map[string]interface{}, len(users))
	for i, v := range users {
		usersToSet[i] = map[string]interface{}{
			"name":     v.Name,
			"readonly": v.Readonly,
		}
	}
	return usersToSet
}

func QueryDatabaseUsers(client *v3.RdsClient, instanceId, dbName string) ([]rds.UserWithPrivilege, error) {
	request := rds.ListAuthorizedDbUsersRequest{
		InstanceId: instanceId,
		DbName:     dbName,
		Limit:      int32(100),
		Page:       int32(1),
	}

	// List all databases
	allUsers := []rds.UserWithPrivilege{}
	for {
		response, err := client.ListAuthorizedDbUsers(&request)
		if err != nil {
			return nil, err
		}
		if response.Users == nil || len(*response.Users) == 0 {
			break
		}

		users := *response.Users
		allUsers = append(allUsers, users...)
		request.Page += 1
	}

	if len(allUsers) < 1 {
		return nil, golangsdk.ErrDefault404{}
	}

	return allUsers, nil
}
