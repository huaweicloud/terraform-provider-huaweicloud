package dms

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/kafka/v2/model"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka PUT /v2/{engine}/{project_id}/instances/{instance_id}/users/{user_name}
// @API Kafka GET /v2/{project_id}/instances/{instance_id}/users
// @API Kafka POST /v2/{project_id}/instances/{instance_id}/users
// @API Kafka PUT /v2/{project_id}/instances/{instance_id}/users
func ResourceDmsKafkaUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsKafkaUserCreate,
		UpdateContext: resourceDmsKafkaUserUpdate,
		DeleteContext: resourceDmsKafkaUserDelete,
		ReadContext:   resourceDmsKafkaUserRead,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"default_app": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceDmsKafkaUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceUser := d.Get("name").(string)
	instanceId := d.Get("instance_id").(string)

	createOpts := &model.CreateInstanceUserRequest{
		InstanceId: instanceId,
		Body: &model.CreateInstanceUserReq{
			UserName:   utils.String(instanceUser),
			UserPasswd: utils.String(d.Get("password").(string)),
			UserDesc:   utils.String(d.Get("description").(string)),
		},
	}

	_, err = client.CreateInstanceUser(createOpts)
	if err != nil {
		return diag.Errorf("error creating DMS instance user: %s", err)
	}

	id := instanceId + "/" + instanceUser
	d.SetId(id)
	return resourceDmsKafkaUserRead(ctx, d, meta)
}

func resourceDmsKafkaUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return diag.Errorf("invalid id format, must be <instance_id>/<user>")
	}
	instanceId := parts[0]
	instanceUser := parts[1]

	// List all instance users
	request := &model.ShowInstanceUsersRequest{
		InstanceId: instanceId,
	}

	response, err := client.ShowInstanceUsers(request)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error listing DMS instance users")
	}

	if response.Users != nil && len(*response.Users) != 0 {
		users := *response.Users
		for _, user := range users {
			if *user.UserName == instanceUser {
				d.Set("instance_id", instanceId)
				d.Set("name", instanceUser)
				d.Set("description", user.UserDesc)
				d.Set("default_app", user.DefaultApp)
				d.Set("role", user.Role)
				d.Set("created_at", utils.FormatTimeStampRFC3339(*user.CreatedTime/1000, false))
				return nil
			}
		}
	}

	// DB user deleted
	d.SetId("")
	log.Printf("[WARN] failed to fetch DMS instance user %s: deleted", instanceUser)

	return nil
}

func resourceDmsKafkaUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)

	updateOpts := &model.UpdateInstanceUserRequest{
		Engine:     "kafka",
		InstanceId: instanceId,
		UserName:   d.Get("name").(string),
		Body: &model.UpdateUserReq{
			UserName:    utils.String(d.Get("name").(string)),
			NewPassword: utils.String(d.Get("password").(string)),
			UserDesc:    utils.String(d.Get("description").(string)),
		},
	}

	_, err = client.UpdateInstanceUser(updateOpts)
	if err != nil {
		return diag.Errorf("error updating DMS instance user: %s", err)
	}

	return resourceDmsKafkaUserRead(ctx, d, meta)
}

func resourceDmsKafkaUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcDmsV2Client(c.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	users := []string{d.Get("name").(string)}
	action := model.GetBatchDeleteInstanceUsersReqActionEnum().DELETE

	deleteOpts := &model.BatchDeleteInstanceUsersRequest{
		InstanceId: instanceId,
		Body: &model.BatchDeleteInstanceUsersReq{
			Action: &action,
			Users:  &users,
		},
	}

	log.Printf("[DEBUG] Delete DMS instance user options: %#v", deleteOpts)
	_, err = client.BatchDeleteInstanceUsers(deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting DMS instance user: %s", err)
	}

	return nil
}
