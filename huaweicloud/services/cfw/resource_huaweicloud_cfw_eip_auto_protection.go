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

var eipAutoProtectionNonUpdatableParams = []string{"fw_instance_id", "object_id", "status", "enterprise_project_id"}

// @API CFW POST /v1/{project_id}/eip/auto-protect-status/switch
// @API CFW GET /v1/{project_id}/eip/auto-protect-status/{object_id}
func ResourceEipAutoProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEipAutoProtectionCreate,
		ReadContext:   resourceEipAutoProtectionRead,
		UpdateContext: resourceEipAutoProtectionUpdate,
		DeleteContext: resourceEipAutoProtectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceEipAutoProtectionImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(eipAutoProtectionNonUpdatableParams),

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
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeInt,
				Required: true,
			},
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
			"available_eip_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"beyond_max_count": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"eip_protected_self": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"eip_total": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"eip_un_protected": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildEipAutoProtectionSwitchQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	return ""
}

func buildEipAutoProtectionSwitchBodyParam(fwInstanceID, objectID string, status int) map[string]interface{} {
	return map[string]interface{}{
		"fw_instance_id": fwInstanceID,
		"object_id":      objectID,
		"status":         status,
	}
}
func resourceEipAutoProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/eip/auto-protect-status/switch"
		product      = "cfw"
		fwInstanceID = d.Get("fw_instance_id").(string)
		objectID     = d.Get("object_id").(string)
		status       = d.Get("status").(int)
		epsId        = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEipAutoProtectionSwitchQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildEipAutoProtectionSwitchBodyParam(fwInstanceID, objectID, status),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error enabling CFW EIP auto protection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// The value of the `data` field is the same as the `object_id`, use it as the resource ID.
	data := utils.PathSearch("data", respBody, "").(string)
	if data == "" {
		return diag.Errorf("error enabling CFW EIP auto protection: data is not found in API response")
	}

	d.SetId(data)

	return resourceEipAutoProtectionRead(ctx, d, meta)
}

func resourceEipAutoProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/eip/auto-protect-status/{object_id}"
		product = "cfw"
		epsId   = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{object_id}", d.Id())
	requestPath += buildEipAutoProtectionSwitchQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving CFW EIP auto protection: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	// When status = 0, that is, when protection is turned off, perform checkDeleted `404` processing.
	statusResp := utils.PathSearch("data.status", respBody, float64(0)).(float64)
	if statusResp == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CFW EIP auto protection")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("object_id", utils.PathSearch("data.object_id", respBody, nil)),
		d.Set("status", statusResp),
		d.Set("available_eip_count", utils.PathSearch("data.available_eip_count", respBody, nil)),
		d.Set("beyond_max_count", utils.PathSearch("data.beyond_max_count", respBody, nil)),
		d.Set("eip_protected_self", utils.PathSearch("data.eip_protected_self", respBody, nil)),
		d.Set("eip_total", utils.PathSearch("data.eip_total", respBody, nil)),
		d.Set("eip_un_protected", utils.PathSearch("data.eip_un_protected", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceEipAutoProtectionUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEipAutoProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		httpUrl      = "v1/{project_id}/eip/auto-protect-status/switch"
		product      = "cfw"
		fwInstanceID = d.Get("fw_instance_id").(string)
		objectID     = d.Get("object_id").(string)
		epsId        = cfg.GetEnterpriseProjectID(d)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildEipAutoProtectionSwitchQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildEipAutoProtectionSwitchBodyParam(fwInstanceID, objectID, 0),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error disabling CFW EIP auto protection: %s", err)
	}

	return nil
}

func resourceEipAutoProtectionImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <fw_instance_id>/<id>")
	}

	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("fw_instance_id", parts[0])
}
