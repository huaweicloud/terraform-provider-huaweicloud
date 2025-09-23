// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CNAD
// ---------------------------------------------------------------

package cnad

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

// @API AAD POST /v1/cnad/policies
// @API AAD GET /v1/cnad/policies
// @API AAD DELETE /v1/cnad/policies/{policy_id}
// @API AAD GET /v1/cnad/policies/{policy_id}
// @API AAD PUT /v1/cnad/policies/{policy_id}
func ResourceCNADAdvancedPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCNADAdvancedPolicyCreate,
		UpdateContext: resourceCNADAdvancedPolicyUpdate,
		ReadContext:   resourceCNADAdvancedPolicyRead,
		DeleteContext: resourceCNADAdvancedPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the CNAD advanced instance ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy name, the maximum length is 255 characters.`,
			},
			"udp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies whether to block the UDP protocol. Valid values are **block** and **unblock**.`,
			},
			"threshold": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the cleaning threshold, the value ranges from 100 to 1000.`,
			},
			"block_location": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The location block list.`,
			},
			"block_protocol": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The protocol block list.`,
			},
			"connection_protection": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether enable connection protection.`,
			},
			"connection_protection_list": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The connection protection list.`,
			},
			"fingerprint_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The fingerprint count.`,
			},
			"port_block_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of port blockages.`,
			},
			"watermark_count": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of watermarks.`,
			},
		},
	}
}

func resourceCNADAdvancedPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createPolicyHttpUrl = "v1/cnad/policies"
		createPolicyProduct = "aad"
	)
	createPolicyClient, err := cfg.NewServiceClient(createPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	createPolicyPath := createPolicyClient.Endpoint + createPolicyHttpUrl
	createPolicyOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePolicyBodyParams(d)),
	}

	createPolicyResp, err := createPolicyClient.Request("POST", createPolicyPath, &createPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating CNAD advanced policy: %s", err)
	}

	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("id", createPolicyRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the CNAD advanced policy ID from the API response")
	}
	d.SetId(policyId)

	_, ok1 := d.GetOk("threshold")
	_, ok2 := d.GetOk("udp")
	if ok1 || ok2 {
		if err = updatePolicy(createPolicyClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceCNADAdvancedPolicyRead(ctx, d, meta)
}

func buildCreatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"package_id": utils.ValueIgnoreEmpty(d.Get("instance_id")),
		"name":       utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceCNADAdvancedPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	getPolicyClient, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	getPolicyDetailPath := getPolicyClient.Endpoint + "v1/cnad/policies/{policy_id}"
	getPolicyDetailPath = strings.ReplaceAll(getPolicyDetailPath, "{policy_id}", d.Id())
	getPolicyDetailOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}

	getPolicyDetailResp, err := getPolicyClient.Request("GET", getPolicyDetailPath, &getPolicyDetailOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CNAD advanced policy")
	}

	respBody, err := utils.FlattenResponse(getPolicyDetailResp)
	if err != nil {
		return diag.FromErr(err)
	}
	popPolicy := utils.PathSearch("pop_policy", respBody,
		make(map[string]interface{})).(map[string]interface{})

	mErr := multierror.Append(
		nil,
		d.Set("instance_id", utils.PathSearch("package_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("threshold", utils.PathSearch("clean_threshold", respBody, nil)),
		d.Set("block_location", popPolicy["block_location"]),
		d.Set("block_protocol", popPolicy["block_protocol"]),
		d.Set("connection_protection", popPolicy["connection_protection"]),
		d.Set("connection_protection_list", popPolicy["connection_protection_list"]),
		d.Set("fingerprint_count", popPolicy["fingerprint_count"]),
		d.Set("port_block_count", popPolicy["port_block_count"]),
		d.Set("watermark_count", popPolicy["watermark_count"]),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCNADAdvancedPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePolicyClient, err := cfg.NewServiceClient("aad", region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	if err = updatePolicy(updatePolicyClient, d); err != nil {
		return diag.FromErr(err)
	}
	return resourceCNADAdvancedPolicyRead(ctx, d, meta)
}

func updatePolicy(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updatePolicyHttpUrl := "v1/cnad/policies/{policy_id}"
	updatePolicyPath := client.Endpoint + updatePolicyHttpUrl
	updatePolicyPath = strings.ReplaceAll(updatePolicyPath, "{policy_id}", d.Id())

	updatePolicyOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdatePolicyBodyParams(d)),
	}
	_, err := client.Request("PUT", updatePolicyPath, &updatePolicyOpt)
	if err != nil {
		return fmt.Errorf("error updating CNAD advanced policy: %s", err)
	}
	return nil
}

func buildUpdatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":      utils.ValueIgnoreEmpty(d.Get("name")),
		"threshold": utils.ValueIgnoreEmpty(d.Get("threshold")),
		"udp":       utils.ValueIgnoreEmpty(d.Get("udp")),
	}
	return bodyParams
}

func resourceCNADAdvancedPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deletePolicyHttpUrl = "v1/cnad/policies/{policy_id}"
		deletePolicyProduct = "aad"
	)
	deletePolicyClient, err := cfg.NewServiceClient(deletePolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating CNAD Client: %s", err)
	}

	deletePolicyPath := deletePolicyClient.Endpoint + deletePolicyHttpUrl
	deletePolicyPath = strings.ReplaceAll(deletePolicyPath, "{policy_id}", d.Id())

	deletePolicyOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		KeepResponseBody: true,
	}
	_, err = deletePolicyClient.Request("DELETE", deletePolicyPath, &deletePolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting CNAD advanced policy: %s", err)
	}
	return nil
}
