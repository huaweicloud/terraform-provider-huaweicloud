package rds

import (
	"context"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/huaweicloud-sdk-go-v3/services/rds/v3/model"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

func ResourceRdsAccount() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRdsAccountCreate,
		UpdateContext: resourceRdsAccountUpdate,
		DeleteContext: resourceRdsAccountDelete,
		ReadContext:   resourceRdsAccountRead,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
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
		},
	}
}

func resourceRdsAccountCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating RDS client: %s", err)
	}

	dbUser := d.Get("name").(string)
	instanceId := d.Get("instance_id").(string)

	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	createOpts := &model.CreateDbUserRequest{
		InstanceId: instanceId,
		Body: &model.UserForCreation{
			Name:     dbUser,
			Password: d.Get("password").(string),
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		_, err = client.CreateDbUser(createOpts)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error creating RDS database user: %s", err)
	}

	id := instanceId + "/" + dbUser
	d.SetId(id)
	return resourceRdsAccountRead(ctx, d, meta)
}

func resourceRdsAccountRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating RDS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return fmtp.DiagErrorf("invalid id format, must be <instance_id>/<user>")
	}
	instanceId := parts[0]
	dbUser := parts[1]

	// items on every page, [1, 100]
	limit := int32(100)
	// List all db users
	request := &model.ListDbUsersRequest{
		InstanceId: instanceId,
		Limit:      limit,
		Page:       int32(1),
	}

	for {
		response, err := client.ListDbUsers(request)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error listing RDS db users")
		}
		users := *response.Users
		if len(users) == 0 {
			break
		}
		request.Page += 1
		for _, user := range users {
			if user.Name == dbUser {
				d.Set("instance_id", instanceId)
				d.Set("name", dbUser)
				return nil
			}
		}
	}

	// DB user deleted
	d.SetId("")
	log.Printf("[WARN] failed to fetch RDS db user %s: deleted", dbUser)

	return nil
}

func resourceRdsAccountUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	updateOpts := &model.SetDbUserPwdRequest{
		InstanceId: instanceId,
		Body: &model.DbUserPwdRequest{
			Name:     d.Get("name").(string),
			Password: d.Get("password").(string),
		},
	}

	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
		_, err = client.SetDbUserPwd(updateOpts)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error updating RDS database user: %s", err)
	}

	return resourceRdsAccountRead(ctx, d, meta)
}

func resourceRdsAccountDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	c := meta.(*config.Config)
	client, err := c.HcRdsV3Client(c.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("error creating RDS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	config.MutexKV.Lock(instanceId)
	defer config.MutexKV.Unlock(instanceId)

	deleteOpts := &model.DeleteDbUserRequest{
		InstanceId: instanceId,
		UserName:   d.Get("name").(string),
	}

	log.Printf("[DEBUG] Delete RDS db user options: %#v", deleteOpts)
	err = resource.RetryContext(ctx, d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		_, err = client.DeleteDbUser(deleteOpts)
		retryable, err := handleMultiOperationsError(err)
		if retryable {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return diag.Errorf("error deleting RDS database user: %s", err)
	}

	return nil
}
