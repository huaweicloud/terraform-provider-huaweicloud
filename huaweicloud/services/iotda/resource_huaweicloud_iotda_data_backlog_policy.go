package iotda

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/routing-rule/backlog-policy
// @API IoTDA GET /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}
func ResourceDataBacklogPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataBacklogPolicyCreate,
		ReadContext:   resourceDataBacklogPolicyRead,
		UpdateContext: resourceDataBacklogPolicyUpdate,
		DeleteContext: resourceDataBacklogPolicyDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backlog_size": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"backlog_time": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildBacklogRequestValue(value string) interface{} {
	r, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("[ERROR] error converting string value to int value: %s", err)
		return nil
	}

	return r
}

func buildDataBacklogPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name":  utils.ValueIgnoreEmpty(d.Get("name")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"backlog_size": buildBacklogRequestValue(d.Get("backlog_size").(string)),
		"backlog_time": buildBacklogRequestValue(d.Get("backlog_time").(string)),
	}
}

func resourceDataBacklogPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/routing-rule/backlog-policy"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDataBacklogPolicyBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA data backlog policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyID := utils.PathSearch("policy_id", respBody, "").(string)
	if policyID == "" {
		return diag.Errorf("error creating IoTDA data backlog policy: ID is not found in API response")
	}

	d.SetId(policyID)

	return resourceDataBacklogPolicyRead(ctx, d, meta)
}

func resourceDataBacklogPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	// When the resource does not exist, query API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data backlog policy")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	backlogSizeAttribute := utils.PathSearch("backlog_size", respBody, float64(0)).(float64)
	backlogTimeAttribute := utils.PathSearch("backlog_time", respBody, float64(0)).(float64)
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("policy_name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("backlog_size", flattenBackLogAttribute(backlogSizeAttribute)),
		d.Set("backlog_time", flattenBackLogAttribute(backlogTimeAttribute)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackLogAttribute(value float64) string {
	return strconv.FormatInt(int64(value), 10)
}

func resourceDataBacklogPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDataBacklogPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating IoTDA data backlog policy: %s", err)
	}

	return resourceDataBacklogPolicyRead(ctx, d, meta)
}

func resourceDataBacklogPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v5/iot/{project_id}/routing-rule/backlog-policy/{policy_id}"
		product = "iotda"
	)

	isDerived := WithDerivedAuth(cfg, region)
	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{policy_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	// When the resource does not exist, delete API will return `404` error code.
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA data backlog policy")
	}

	return nil
}
