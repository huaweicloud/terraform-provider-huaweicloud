package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventSystemUserWhiteListNonUpdatableParams = []string{"host_id", "enterprise_project_id"}

// @API HSS POST /v5/{project_id}/event/white-list/userlist
// @API HSS GET /v5/{project_id}/event/white-list/userlist
// @API HSS PUT /v5/{project_id}/event/white-list/userlist
// @API HSS DELETE /v5/{project_id}/event/white-list/userlist
func ResourceEventSystemUserWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventSystemUserWhiteListCreate,
		ReadContext:   resourceEventSystemUserWhiteListRead,
		UpdateContext: resourceEventSystemUserWhiteListUpdate,
		DeleteContext: resourceEventSystemUserWhiteListDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(eventSystemUserWhiteListNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"system_user_name_list": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"remarks": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_all": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"enterprise_project_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildEventSystemUserWhiteListQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildCreateOrUpdateEventSystemUserWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"host_id":               d.Get("host_id"),
		"system_user_name_list": utils.ExpandToStringList(d.Get("system_user_name_list").([]interface{})),
		"remarks":               utils.ValueIgnoreEmpty(d.Get("remarks")),
	}

	return bodyParams
}

func resourceEventSystemUserWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		hostID = d.Get("host_id").(string)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildEventSystemUserWhiteListQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateEventSystemUserWhiteListBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating HSS event system user white list: %s", err)
	}

	d.SetId(hostID)

	return resourceEventSystemUserWhiteListRead(ctx, d, meta)
}

func buildReadEventSystemUserWhiteListQueryParams(hostID, epsId string) string {
	queryParams := fmt.Sprintf("?host_id=%s", hostID)

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%s", queryParams, epsId)
	}

	return queryParams
}

func resourceEventSystemUserWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
		hostID = d.Id()
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	queryPath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath += buildReadEventSystemUserWhiteListQueryParams(hostID, epsId)
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS event system user white list: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataList := utils.PathSearch("data_list", respBody, make([]interface{}, 0)).([]interface{})
	// The List API always returns a status code of `200`.
	// If the list API returns an empty `data_list`, the resource is considered not found.
	if len(dataList) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_id", utils.PathSearch("host_id", dataList[0], nil)),
		d.Set("system_user_name_list", utils.ExpandToStringList(
			utils.PathSearch("system_user_name_list", dataList[0], make([]interface{}, 0)).([]interface{}))),
		d.Set("remarks", utils.PathSearch("remarks", dataList[0], nil)),
		d.Set("enterprise_project_name", utils.PathSearch("enterprise_project_name", dataList[0], nil)),
		d.Set("host_name", utils.PathSearch("host_name", dataList[0], nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", dataList[0], nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", dataList[0], nil)),
		d.Set("update_time", utils.PathSearch("update_time", dataList[0], nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEventSystemUserWhiteListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	updatePath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath += buildEventSystemUserWhiteListQueryParams(epsId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateOrUpdateEventSystemUserWhiteListBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating HSS event system user white list: %s", err)
	}

	return resourceEventSystemUserWhiteListRead(ctx, d, meta)
}

func buildDeleteEventSystemUserWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	// When `delete_all = true` is specified, the body of delete only needs to pass the `delete_all` parameter.
	if d.Get("delete_all").(bool) {
		return map[string]interface{}{
			"delete_all": true,
		}
	}

	deleteBodyParams := map[string]interface{}{
		"host_id":               d.Get("host_id"),
		"system_user_name_list": utils.ExpandToStringList(d.Get("system_user_name_list").([]interface{})),
	}

	return map[string]interface{}{
		"data_list": []interface{}{
			deleteBodyParams,
		},
	}
}

func resourceEventSystemUserWhiteListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + "v5/{project_id}/event/white-list/userlist"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildEventSystemUserWhiteListQueryParams(epsId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteEventSystemUserWhiteListBodyParams(d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS event system user white list: %s", err)
	}

	return nil
}
