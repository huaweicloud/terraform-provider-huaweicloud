package rabbitmq

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

// @API RabbitMQ PUT /v2/{project_id}/instances/{instance_id}/users/{user_name}
// @API RabbitMQ POST /v2/{project_id}/instances/{instance_id}/users
// @API RabbitMQ GET /v2/{project_id}/instances/{instance_id}/users
// @API RabbitMQ DELETE /v2/{project_id}/instances/{instance_id}/users/{user_name}
func ResourceDmsRabbitmqUser() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDmsRabbitmqUserCreate,
		ReadContext:   resourceDmsRabbitmqUserRead,
		UpdateContext: resourceDmsRabbitmqUserUpdate,
		DeleteContext: resourceDmsRabbitmqUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceRabbitmqUserImportState,
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
			"access_key": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"secret_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"vhosts": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"vhost": {
							Type:     schema.TypeString,
							Required: true,
						},
						"conf": {
							Type:     schema.TypeString,
							Required: true,
						},
						"write": {
							Type:     schema.TypeString,
							Required: true,
						},
						"read": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func resourceDmsRabbitmqUserCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	accessKey := d.Get("access_key").(string)

	createHttpUrl := "v2/{project_id}/instances/{instance_id}/users"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRabbitmqUserBodyParams(d),
	}
	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating user: %s", err)
	}

	id := fmt.Sprintf("%s/%s", instanceID, accessKey)
	d.SetId(id)

	return resourceDmsRabbitmqUserRead(ctx, d, cfg)
}

func buildRabbitmqUserBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_key": utils.ValueIgnoreEmpty(d.Get("access_key")),
		"secret_key": utils.ValueIgnoreEmpty(d.Get("secret_key")),
		"vhosts":     buildRabbitmqUserBodyParamsVhosts(d),
	}
	return bodyParams
}

func buildRabbitmqUserBodyParamsVhosts(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("vhosts").(*schema.Set).List()
	if len(rawParams) == 0 {
		return nil
	}
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, v := range rawParams {
		params := v.(map[string]interface{})
		m := map[string]interface{}{
			"vhost": params["vhost"],
			"conf":  params["conf"],
			"write": params["write"],
			"read":  params["read"],
		}
		rst = append(rst, m)
	}
	return rst
}

func resourceDmsRabbitmqUserRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	accessKey := d.Get("access_key").(string)

	user, err := GetRabbitmqUser(client, instanceID, accessKey)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving the user")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("secret_key", utils.PathSearch("secret_key", user, nil)),
		d.Set("vhosts", flattenRabbitmqUserVhosts(utils.PathSearch("vhosts", user, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetRabbitmqUser(client *golangsdk.ServiceClient, instanceID, accessKey string) (interface{}, error) {
	listHttpUrl := "v2/{project_id}/instances/{instance_id}/users"
	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// pageLimit is `10`
	listPath += fmt.Sprintf("?limit=%d", pageLimit)
	offset := 0
	for {
		currentPath := listPath + fmt.Sprintf("&offset=%d", offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		searchPath := fmt.Sprintf("users[?access_key=='%s']|[0]", accessKey)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		// `total` means the number of all users, and type is float64.
		offset += pageLimit
		total := utils.PathSearch("total", listRespBody, float64(0))
		if int(total.(float64)) <= offset {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v2/rabbitmq/{project_id}/instances/{instance_id}/users",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("the user (%s) does not exist", accessKey)),
				},
			}
		}
	}
}

func flattenRabbitmqUserVhosts(rawParams []interface{}) interface{} {
	if len(rawParams) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(rawParams))
	for _, params := range rawParams {
		rst = append(rst, map[string]interface{}{
			"vhost": utils.PathSearch("vhost", params, nil),
			"conf":  utils.PathSearch("conf", params, nil),
			"write": utils.PathSearch("write", params, nil),
			"read":  utils.PathSearch("read", params, nil),
		})
	}
	return rst
}

func resourceDmsRabbitmqUserUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	accessKey := d.Get("access_key").(string)

	updateHttpUrl := "v2/{project_id}/instances/{instance_id}/users/{user_name}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{instance_id}", instanceID)
	updatePath = strings.ReplaceAll(updatePath, "{user_name}", accessKey)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildRabbitmqUserBodyParams(d),
	}
	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating user: %s", err)
	}

	return resourceDmsRabbitmqUserRead(ctx, d, cfg)
}

func resourceDmsRabbitmqUserDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	accessKey := d.Get("access_key").(string)

	deleteHttpUrl := "v2/{project_id}/instances/{instance_id}/users/{user_name}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", instanceID)
	deletePath = strings.ReplaceAll(deletePath, "{user_name}", accessKey)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DMS.00500972"),
			"error deleting user")
	}

	return nil
}

func resourceRabbitmqUserImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <instance_id>/<access_key>")
	}

	d.Set("instance_id", parts[0])
	d.Set("access_key", parts[1])

	return []*schema.ResourceData{d}, nil
}
