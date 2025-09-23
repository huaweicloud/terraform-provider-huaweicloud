package kafka

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

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
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	instanceUser := d.Get("name").(string)

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/users"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildCreateDmsKafkaUserBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DMS instance user: %s", err)
	}

	id := instanceId + "/" + instanceUser
	d.SetId(id)
	return resourceDmsKafkaUserRead(ctx, d, meta)
}

func buildCreateDmsKafkaUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name":   d.Get("name"),
		"user_desc":   d.Get("description"),
		"user_passwd": d.Get("password"),
	}
	return bodyParams
}

func resourceDmsKafkaUserRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	// Split instance_id and user from resource id
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<user>")
	}
	instanceId := parts[0]
	instanceUser := parts[1]

	user, err := GetDmsKafkaUser(client, instanceId, instanceUser)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error listing DMS instance users")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("instance_id", instanceId),
		d.Set("name", instanceUser),
		d.Set("description", utils.PathSearch("user_desc", user, nil)),
		d.Set("default_app", utils.PathSearch("default_app", user, nil)),
		d.Set("role", utils.PathSearch("role", user, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("created_time", user, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetDmsKafkaUser(client *golangsdk.ServiceClient, instId, name string) (interface{}, error) {
	getHttpUrl := "v2/{project_id}/instances/{instance_id}/users"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	searchPath := fmt.Sprintf("users[?user_name=='%s']|[0]", name)
	user := utils.PathSearch(searchPath, getRespBody, nil)
	if user == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return user, nil
}

func resourceDmsKafkaUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	updateHttpUrl := "v2/{engine}/{project_id}/instances/{instance_id}/users/{user_name}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{engine}", "kafka")
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", d.Get("instance_id").(string))
	updatePath = strings.ReplaceAll(updatePath, "{user_name}", d.Get("name").(string))

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: utils.RemoveNil(buildUpdateDmsKafkaUserBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DMS instance user: %s", err)
	}

	return resourceDmsKafkaUserRead(ctx, d, meta)
}

func buildUpdateDmsKafkaUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"user_name":    d.Get("name"),
		"user_desc":    d.Get("description"),
		"new_password": d.Get("password"),
	}
	return bodyParams
}

func resourceDmsKafkaUserDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.DmsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/users"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", d.Get("instance_id").(string))

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		JSONBody: map[string]interface{}{
			"action": "delete",
			"users":  []interface{}{d.Get("name")},
		},
	}

	_, err = client.Request("PUT", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting DMS instance user")
	}

	return nil
}
