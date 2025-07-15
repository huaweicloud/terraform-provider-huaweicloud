package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var eventLoginWhiteListNonUpdatableParams = []string{"private_ip", "login_ip", "login_user_name", "remarks",
	"handle_event", "enterprise_project_id"}

// @API HSS POST /v5/{project_id}/event/white-list/login
// @API HSS GET /v5/{project_id}/event/white-list/login
// @API HSS DELETE /v5/{project_id}/event/white-list/login
func ResourceEventLoginWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEventLoginWhiteListCreate,
		ReadContext:   resourceEventLoginWhiteListRead,
		UpdateContext: resourceEventLoginWhiteListUpdate,
		DeleteContext: resourceEventLoginWhiteListDelete,

		// Due to the special nature of the query API, this resource currently does not support the import function.

		CustomizeDiff: config.FlexibleForceNew(eventLoginWhiteListNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"remarks": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"handle_event": {
				Type:     schema.TypeBool,
				Optional: true,
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
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildEventLoginWhiteListQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildCreateEventLoginWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"private_ip":      d.Get("private_ip"),
		"login_ip":        d.Get("login_ip"),
		"login_user_name": d.Get("login_user_name"),
		"remarks":         utils.ValueIgnoreEmpty(d.Get("remarks")),
	}

	if d.Get("handle_event").(bool) {
		bodyParams["handle_event"] = true
	}

	return bodyParams
}

func resourceEventLoginWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	createPath := client.Endpoint + "v5/{project_id}/event/white-list/login"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath += buildEventLoginWhiteListQueryParams(epsId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateEventLoginWhiteListBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating HSS event login white list: %s", err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	return resourceEventLoginWhiteListRead(ctx, d, meta)
}

func buildReadEventLoginWhiteListQueryParams(d *schema.ResourceData, epsId string) string {
	// In the list API, `limit` and `offset` are required parameters.
	queryParams := "?limit=200&offset=0"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%s", queryParams, epsId)
	}

	queryParams = fmt.Sprintf("%s&private_ip=%s", queryParams, d.Get("private_ip"))
	queryParams = fmt.Sprintf("%s&login_ip=%s", queryParams, d.Get("login_ip"))
	queryParams = fmt.Sprintf("%s&login_user_name=%s", queryParams, d.Get("login_user_name"))

	return queryParams
}

func resourceEventLoginWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	queryPath := client.Endpoint + "v5/{project_id}/event/white-list/login"
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath += buildReadEventLoginWhiteListQueryParams(d, epsId)
	queryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", queryPath, &queryOpt)
	if err != nil {
		return diag.Errorf("error retrieving HSS event login white list: %s", err)
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
		d.Set("private_ip", utils.PathSearch("private_ip", dataList[0], nil)),
		d.Set("login_ip", utils.PathSearch("login_ip", dataList[0], nil)),
		d.Set("login_user_name", utils.PathSearch("login_user_name", dataList[0], nil)),
		d.Set("remarks", utils.PathSearch("remarks", dataList[0], nil)),
		d.Set("enterprise_project_name", utils.PathSearch("enterprise_project_name", dataList[0], nil)),
		d.Set("update_time", utils.PathSearch("update_time", dataList[0], nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEventLoginWhiteListUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// This resource doesn't support update operation.
	return nil
}

func buildDeleteEventLoginWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	// When `delete_all = true` is specified, the body of delete only needs to pass the `delete_all` parameter.
	if d.Get("delete_all").(bool) {
		return map[string]interface{}{
			"delete_all": true,
		}
	}

	return map[string]interface{}{
		"data_list": []interface{}{
			map[string]interface{}{
				"private_ip":      d.Get("private_ip"),
				"login_ip":        d.Get("login_ip"),
				"login_user_name": d.Get("login_user_name"),
			},
		},
	}
}

func resourceEventLoginWhiteListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient("hss", region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	deletePath := client.Endpoint + "v5/{project_id}/event/white-list/login"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath += buildEventLoginWhiteListQueryParams(epsId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildDeleteEventLoginWhiteListBodyParams(d),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting HSS event login white list: %s", err)
	}

	return nil
}
