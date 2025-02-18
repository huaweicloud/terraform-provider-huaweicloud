package iotda

import (
	"context"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IoTDA POST /v5/iot/{project_id}/routing-rule/flowcontrol-policy
// @API IoTDA GET /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
// @API IoTDA PUT /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
// @API IoTDA DELETE /v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}
func ResourceDataFlowControlPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataFlowControlPolicyCreate,
		ReadContext:   resourceDataFlowControlPolicyRead,
		UpdateContext: resourceDataFlowControlPolicyUpdate,
		DeleteContext: resourceDataFlowControlPolicyDelete,

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
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"scope_value": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func buildCreateDataFlowControlPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"scope":       utils.ValueIgnoreEmpty(d.Get("scope")),
		"scope_value": utils.ValueIgnoreEmpty(d.Get("scope_value")),
		"limit":       utils.ValueIgnoreEmpty(d.Get("limit")),
	}
}

func resourceDataFlowControlPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	log.Printf("1111111111: %v", d.Get("limit"))
	createPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/flowcontrol-policy"
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDataFlowControlPolicyBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA data flow control policy: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy_id", createRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("error creating IoTDA data flow control policy: ID is not found in API response")
	}

	d.SetId(policyId)

	return resourceDataFlowControlPolicyRead(ctx, d, meta)
}

func resourceDataFlowControlPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	getPath := client.Endpoint + "v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{policy_id}", d.Id())
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the resource does not exist, query API will return `404` error code.
	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA data flow control policy")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("policy_name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("scope", utils.PathSearch("scope", getRespBody, nil)),
		d.Set("scope_value", utils.PathSearch("scope_value", getRespBody, nil)),
		d.Set("limit", utils.PathSearch("limit", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDataFlowControlPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"policy_name": utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"limit":       utils.ValueIgnoreEmpty(d.Get("limit")),
	}
}

func resourceDataFlowControlPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	updatePath := client.Endpoint + "v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}"
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDataFlowControlPolicyBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating IoTDA data flow control policy: %s", err)
	}

	return resourceDataFlowControlPolicyRead(ctx, d, meta)
}

func resourceDataFlowControlPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		product   = "iotda"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth(product, region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + "v5/iot/{project_id}/routing-rule/flowcontrol-policy/{policy_id}"
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// When the resource does not exist, delete API will return `404` error code.
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA data flow control policy")
	}

	return nil
}
