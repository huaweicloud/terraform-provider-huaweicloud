// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Organizations
// ---------------------------------------------------------------

package organizations

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

// @API Organizations POST /v1/organizations/policies
// @API Organizations DELETE /v1/organizations/policies/{policy_id}
// @API Organizations GET /v1/organizations/policies/{policy_id}
// @API Organizations PATCH /v1/organizations/policies/{policy_id}
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		UpdateContext: resourcePolicyUpdate,
		ReadContext:   resourcePolicyRead,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name to be assigned to the policy.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the policy text content to be added to the new policy.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of the policy to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description to be assigned to the policy.`,
			},
			"tags": common.TagsSchema(),
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the uniform resource name of the policy.`,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPolicy: create Organizations policy
	var (
		createPolicyHttpUrl = "v1/organizations/policies"
		createPolicyProduct = "organizations"
	)
	createPolicyClient, err := cfg.NewServiceClient(createPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPolicyPath := createPolicyClient.Endpoint + createPolicyHttpUrl

	createPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createPolicyOpt.JSONBody = utils.RemoveNil(buildCreatePolicyBodyParams(d))
	createPolicyResp, err := createPolicyClient.Request("POST", createPolicyPath, &createPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating Organizations policy: %s", err)
	}

	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy.policy_summary.id", createPolicyRespBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the Organizations policy ID from the API response")
	}
	d.SetId(policyId)

	return resourcePolicyRead(ctx, d, meta)
}

func buildCreatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"content":     utils.ValueIgnoreEmpty(d.Get("content")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":        utils.ExpandResourceTags(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPolicy: Query Organizations policy
	var (
		getPolicyHttpUrl = "v1/organizations/policies/{policy_id}"
		getPolicyProduct = "organizations"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl
	getPolicyPath = strings.ReplaceAll(getPolicyPath, "{policy_id}", d.Id())

	getPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getPolicyResp, err := getPolicyClient.Request("GET", getPolicyPath, &getPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Organizations policy")
	}

	getPolicyRespBody, err := utils.FlattenResponse(getPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("policy.policy_summary.name", getPolicyRespBody, nil)),
		d.Set("content", utils.PathSearch("policy.content", getPolicyRespBody, nil)),
		d.Set("type", utils.PathSearch("policy.policy_summary.type", getPolicyRespBody, nil)),
		d.Set("description", utils.PathSearch("policy.policy_summary.description", getPolicyRespBody, nil)),
		d.Set("urn", utils.PathSearch("policy.policy_summary.urn", getPolicyRespBody, nil)),
	)

	tagMap, err := getTags(getPolicyClient, policiesType, d.Id())
	if err != nil {
		log.Printf("[WARN] error fetching tags of Organizations policy (%s): %s", d.Id(), err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePolicyChanges := []string{
		"name",
		"content",
		"description",
	}

	// updatePolicy: update Organizations policy
	var (
		updatePolicyHttpUrl = "v1/organizations/policies/{policy_id}"
		updatePolicyProduct = "organizations"
	)
	updatePolicyClient, err := cfg.NewServiceClient(updatePolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.HasChanges(updatePolicyChanges...) {
		updatePolicyPath := updatePolicyClient.Endpoint + updatePolicyHttpUrl
		updatePolicyPath = strings.ReplaceAll(updatePolicyPath, "{policy_id}", d.Id())

		updatePolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updatePolicyOpt.JSONBody = utils.RemoveNil(buildUpdatePolicyBodyParams(d))
		_, err = updatePolicyClient.Request("PATCH", updatePolicyPath, &updatePolicyOpt)
		if err != nil {
			return diag.Errorf("error updating Organizations policy: %s", err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, updatePolicyClient, policiesType, d.Id(), "tags")
		if err != nil {
			return diag.Errorf("error updating tags of Organizations policy %s: %s", d.Id(), err)
		}
	}

	return resourcePolicyRead(ctx, d, meta)
}

func buildUpdatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"content":     utils.ValueIgnoreEmpty(d.Get("content")),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourcePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePolicy: Delete Organizations policy
	var (
		deletePolicyHttpUrl = "v1/organizations/policies/{policy_id}"
		deletePolicyProduct = "organizations"
	)
	deletePolicyClient, err := cfg.NewServiceClient(deletePolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePolicyPath := deletePolicyClient.Endpoint + deletePolicyHttpUrl
	deletePolicyPath = strings.ReplaceAll(deletePolicyPath, "{policy_id}", d.Id())

	deletePolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deletePolicyClient.Request("DELETE", deletePolicyPath, &deletePolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting Organizations policy: %s", err)
	}

	return nil
}
