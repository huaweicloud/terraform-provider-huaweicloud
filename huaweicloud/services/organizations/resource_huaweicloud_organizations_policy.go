package organizations

import (
	"context"
	"fmt"
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
// @API Organizations GET /v1/organizations/policies/{policy_id}
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations PATCH /v1/organizations/policies/{policy_id}
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations DELETE /v1/organizations/policies/{policy_id}
func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		ReadContext:   resourcePolicyRead,
		UpdateContext: resourcePolicyUpdate,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.MergeDefaultTags(),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name to be assigned to the policy.`,
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The policy text content to be added to the new policy.`,
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					equal, _ := utils.CompareJsonTemplateAreEquivalent(old, new)
					return equal
				},
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the policy to be created.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description to be assigned to the policy.`,
			},
			"tags": common.TagsSchema(`The key/value pairs associated with the policy.`),
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the policy.`,
			},
		},
	}
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations/policies"
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePolicyBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	policyId := utils.PathSearch("policy.policy_summary.id", respBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the policy ID from the API response")
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
		"tags":        utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
	return bodyParams
}

func GetPolicyById(client *golangsdk.ServiceClient, policyId string) (interface{}, error) {
	httpUrl := "v1/organizations/policies/{policy_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{policy_id}", policyId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
		mErr     *multierror.Error
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	respBody, err := GetPolicyById(client, policyId)
	if err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error retrieving policy (%s)", policyId),
		)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("policy.policy_summary.name", respBody, nil)),
		d.Set("content", utils.PathSearch("policy.content", respBody, nil)),
		d.Set("type", utils.PathSearch("policy.policy_summary.type", respBody, nil)),
		d.Set("description", utils.PathSearch("policy.policy_summary.description", respBody, nil)),
		d.Set("urn", utils.PathSearch("policy.policy_summary.urn", respBody, nil)),
	)

	tagMap, err := getTags(client, policiesType, policyId)
	if err != nil {
		log.Printf("[WARN] error fetching tags of policy (%s): %s", policyId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tagMap))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/policies/{policy_id}"
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.HasChanges("name", "content", "description") {
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{policy_id}", policyId)
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdatePolicyBodyParams(d)),
		}

		_, err = client.Request("PATCH", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating policy (%s): %s", policyId, err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, client, policiesType, policyId, "tags")
		if err != nil {
			return diag.Errorf("error updating tags of policy (%s): %s", policyId, err)
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
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/policies/{policy_id}"
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{policy_id}", policyId)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error deleting policy (%s)", policyId),
		)
	}

	return nil
}
