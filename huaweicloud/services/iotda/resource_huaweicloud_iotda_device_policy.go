package iotda

import (
	"context"
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

var devicePolicyNonUpdatableParams = []string{"space_id"}

// @API IoTDA POST /v5/iot/{project_id}/device-policies
// @API IoTDA GET /v5/iot/{project_id}/device-policies/{policy_id}
// @API IoTDA PUT /v5/iot/{project_id}/device-policies/{policy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/device-policies/{policy_id}
func ResourceDevicePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDevicePolicyCreate,
		ReadContext:   resourceDevicePolicyRead,
		UpdateContext: resourceDevicePolicyUpdate,
		DeleteContext: resourceDevicePolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(devicePolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"policy_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"statement": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"effect": {
							Type:     schema.TypeString,
							Required: true,
						},
						"actions": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"resources": {
							Type:     schema.TypeList,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateOrUpdateStatementBodyParams(d *schema.ResourceData) []map[string]interface{} {
	statementList := d.Get("statement").([]interface{})
	if len(statementList) < 1 {
		return nil
	}

	bodyParams := make([]map[string]interface{}, 0, len(statementList))
	for _, v := range statementList {
		statementMap := v.(map[string]interface{})
		statementParams := map[string]interface{}{
			"effect":    statementMap["effect"],
			"actions":   utils.ExpandToStringList(statementMap["actions"].([]interface{})),
			"resources": utils.ExpandToStringList(statementMap["resources"].([]interface{})),
		}

		bodyParams = append(bodyParams, statementParams)
	}

	return bodyParams
}

func buildCreateDevicePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name": d.Get("policy_name"),
		"app_id":      utils.ValueIgnoreEmpty(d.Get("space_id")),
		"statement":   buildCreateOrUpdateStatementBodyParams(d),
	}

	return bodyParams
}

func resourceDevicePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-policies"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDevicePolicyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device policy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", createRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating IoTDA device policy: ID is not found in API response")
	}

	d.SetId(policyId)

	return resourceDevicePolicyRead(ctx, d, meta)
}

func resourceDevicePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-policies/{policy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", d.Id())
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		// When the resource does not exist, query API will return `404` error code.
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("policy_name", utils.PathSearch("policy_name", getRespBody, nil)),
		d.Set("statement", flattenStatement(utils.PathSearch("statement", getRespBody, nil))),
		d.Set("space_id", utils.PathSearch("app_id", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("update_time", utils.PathSearch("update_time", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStatement(statementResp interface{}) []map[string]interface{} {
	if statementResp == nil {
		return nil
	}

	statementList := statementResp.([]interface{})
	if len(statementList) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(statementList))
	for _, v := range statementList {
		statement := map[string]interface{}{
			"effect":    utils.PathSearch("effect", v, nil),
			"actions":   utils.ExpandToStringList(utils.PathSearch("actions", v, make([]interface{}, 0)).([]interface{})),
			"resources": utils.ExpandToStringList(utils.PathSearch("resources", v, make([]interface{}, 0)).([]interface{})),
		}

		result = append(result, statement)
	}

	return result
}

func buildUpdateDevicePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"policy_name": d.Get("policy_name"),
		"statement":   buildCreateOrUpdateStatementBodyParams(d),
	}

	return bodyParams
}

func resourceDevicePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-policies/{policy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Id())
	updateOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildUpdateDevicePolicyBodyParams(d),
	}

	_, err = client.Request("PUT", updatePath, &updateOpts)
	if err != nil {
		return diag.Errorf("error updating IoTDA device policy: %s", err)
	}

	return resourceDevicePolicyRead(ctx, d, meta)
}

func resourceDevicePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/device-policies/{policy_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		return diag.Errorf("error deleting IoTDA device policy: %s", err)
	}

	return nil
}
