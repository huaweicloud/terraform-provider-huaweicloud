package organizations

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var dryRunPolicyNonUpdatableParams = []string{"type"}

// @API Organizations POST /v1/organizations/dry-run-policies
// @API Organizations GET /v1/organizations/dry-run-policies/{policy_id}
// @API Organizations GET /v1/organizations/{resource_type}/{resource_id}/tags
// @API Organizations PATCH /v1/organizations/dry-run-policies/{policy_id}
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/create
// @API Organizations POST /v1/organizations/{resource_type}/{resource_id}/tags/delete
// @API Organizations DELETE /v1/organizations/dry-run-policies/{policy_id}
func ResourceDryRunPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDryRunPolicyCreate,
		ReadContext:   resourceDryRunPolicyRead,
		UpdateContext: resourceDryRunPolicyUpdate,
		DeleteContext: resourceDryRunPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(dryRunPolicyNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the dry-run policy.`,
			},
			"content": {
				Type:         schema.TypeString,
				Required:     true,
				Description:  `The content of the dry-run policy, in JSON format.`,
				ValidateFunc: validation.StringIsJSON,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the dry-run policy.`,
			},
			// Optional parameters.
			"tags": common.TagsSchema(`The key/value pairs associated with the dry-run policy.`),
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the dry-run policy.`,
			},
			// Attributes.
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The uniform resource name of the dry-run policy.`,
			},
			"is_builtin": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the dry-run policy is a built-in policy.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceDryRunPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/organizations/dry-run-policies"
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildCreateDryRunPolicyBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error creating dry-run policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening dry-run policy create response: %s", err)
	}

	policyId := utils.PathSearch("policy.policy_summary.id", respBody, "").(string)
	if policyId == "" {
		return diag.Errorf("unable to find the dry-run policy ID from the API response")
	}

	d.SetId(policyId)

	return resourceDryRunPolicyRead(ctx, d, meta)
}

func buildCreateDryRunPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":    d.Get("name"),
		"content": d.Get("content"),
		"type":    d.Get("type"),
		// The API requires description to be specified, so when not specified, pass through as empty string.
		"description": d.Get("description"),
		"tags":        utils.ValueIgnoreEmpty(utils.ExpandResourceTags(d.Get("tags").(map[string]interface{}))),
	}
}

// GetDryRunPolicyById queries a dry-run policy by ID.
func GetDryRunPolicyById(client *golangsdk.ServiceClient, policyId string) (interface{}, error) {
	httpUrl := "v1/organizations/dry-run-policies/{policy_id}"
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

func resourceDryRunPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", region)
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	respBody, err := GetDryRunPolicyById(client, policyId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", organizationNotFoundErrCodes...),
			fmt.Sprintf("error retrieving dry-run policy (%s)", policyId),
		)
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("policy.policy_summary.name", respBody, nil)),
		d.Set("content", utils.PathSearch("policy.content", respBody, nil)),
		d.Set("type", utils.PathSearch("policy.policy_summary.type", respBody, nil)),
		d.Set("description", utils.PathSearch("policy.policy_summary.description", respBody, nil)),
		d.Set("urn", utils.PathSearch("policy.policy_summary.urn", respBody, nil)),
		d.Set("is_builtin", utils.PathSearch("policy.policy_summary.is_builtin", respBody, nil)),
	)

	tags, err := getTags(client, policiesType, policyId)
	if err != nil {
		log.Printf("[WARN] error getting tags of dry-run policy (%s): %s", policyId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("tags", tags))
	}

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateDryRunPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"content":     d.Get("content"),
		"description": d.Get("description"),
	}
}

func updateDryRunPolicy(client *golangsdk.ServiceClient, policyId string, d *schema.ResourceData) error {
	httpUrl := "v1/organizations/dry-run-policies/{policy_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{policy_id}", policyId)
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         utils.RemoveNil(buildUpdateDryRunPolicyBodyParams(d)),
	}

	_, err := client.Request("PATCH", updatePath, &updateOpt)
	return err
}

func resourceDryRunPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		policyId = d.Id()
	)
	client, err := cfg.NewServiceClient("organizations", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Organizations client: %s", err)
	}

	if d.HasChanges("name", "content", "description") {
		if err = updateDryRunPolicy(client, policyId, d); err != nil {
			return diag.Errorf("error updating dry-run policy (%s): %s", policyId, err)
		}
	}

	if d.HasChange("tags") {
		err = updateTags(d, client, policiesType, policyId, "tags")
		if err != nil {
			return diag.Errorf("error updating tags of dry-run policy (%s): %s", policyId, err)
		}
	}

	return resourceDryRunPolicyRead(ctx, d, meta)
}

func resourceDryRunPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		httpUrl  = "v1/organizations/dry-run-policies/{policy_id}"
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
			fmt.Sprintf("error deleting dry-run policy (%s)", policyId),
		)
	}

	return nil
}
