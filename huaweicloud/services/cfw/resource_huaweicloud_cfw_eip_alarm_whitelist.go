package cfw

import (
	"context"
	"errors"
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

var eipAlarmWhitelistNonUpdatableParams = []string{"fw_instance_id", "eip_id", "public_ip", "object_id", "public_ipv6",
	"enterprise_project_id"}

// @API CFW POST /v1/{project_id}/eip/alarm-whitelist
// @API CFW GET /v1/{project_id}/eip/alarm-whitelist/{fw_instance_id}
func ResourceEipAlarmWhitelist() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipAlarmWhitelistCreate,
		ReadContext:   resourceEipAlarmWhitelistRead,
		UpdateContext: resourceEipAlarmWhitelistUpdate,
		DeleteContext: resourceEipAlarmWhitelistDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEipAlarmWhitelistImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(eipAlarmWhitelistNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			// The API does not return this field, so the `Computed` attribute is not added.
			"object_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"public_ipv6": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The API does not return this field, so the `Computed` attribute is not added.
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"device_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildAddEipAlarmWhitelistBodyParam(d *schema.ResourceData) map[string]interface{} {
	eipInfo := utils.RemoveNil(map[string]interface{}{
		"eip_id":         d.Get("eip_id"),
		"public_ip":      d.Get("public_ip"),
		"fw_instance_id": d.Get("fw_instance_id").(string),
		"object_id":      utils.ValueIgnoreEmpty(d.Get("object_id")),
		"public_ipv6":    utils.ValueIgnoreEmpty(d.Get("public_ipv6")),
	})

	return map[string]interface{}{
		"fw_instance_id": d.Get("fw_instance_id"),
		"eip_infos":      []interface{}{eipInfo},
	}
}

func resourceEipAlarmWhitelistCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/eip/alarm-whitelist"
		product = "cfw"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildAddEipAlarmWhitelistBodyParam(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error adding CFW EIP alarm whitelist: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", respBody, "").(string)
	if id == "" {
		return diag.Errorf("error adding CFW EIP alarm whitelist: ID is not found in API response")
	}

	d.SetId(id)

	return resourceEipAlarmWhitelistRead(ctx, d, meta)
}

func BuildEipAlarmWhitelistQueryParams(epsId, publicIP string) string {
	queryParams := fmt.Sprintf("?ip_address=%s", publicIP)
	// This query does not require pagination, but the `limit` and `offset` parameters in the API are required.
	queryParams += "&limit=1024&offset=0"
	if epsId != "" {
		queryParams += fmt.Sprintf("&eps_id=%s", epsId)
	}

	return queryParams
}
func resourceEipAlarmWhitelistRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/eip/alarm-whitelist/{fw_instance_id}"
		product  = "cfw"
		epsId    = cfg.GetEnterpriseProjectID(d)
		publicIP = d.Get("public_ip").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{fw_instance_id}", d.Id())
	requestPath += BuildEipAlarmWhitelistQueryParams(epsId, publicIP)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW EIP alarm whitelist: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	listResp := utils.PathSearch("data.list", respBody, make([]interface{}, 0)).([]interface{})
	if len(listResp) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	typeResp := utils.PathSearch("type", listResp[0], float64(0)).(float64)
	if typeResp == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("eip_id", utils.PathSearch("eip_id", listResp[0], nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", listResp[0], nil)),
		d.Set("public_ipv6", utils.PathSearch("public_ipv6", listResp[0], nil)),
		d.Set("device_name", utils.PathSearch("device_name", listResp[0], nil)),
		d.Set("type", typeResp),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEipAlarmWhitelistUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEipAlarmWhitelistDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `Currently, this resource does not support the delete function. Deleting this resource will not clear
    the corresponding request record, but will only remove the resource information from the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}

func resourceEipAlarmWhitelistImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <id>/<public_ip>")
	}

	d.SetId(parts[0])

	mErr := multierror.Append(nil,
		d.Set("fw_instance_id", parts[0]),
		d.Set("public_ip", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
