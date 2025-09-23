// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product UCS
// ---------------------------------------------------------------

package ucs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API UCS GET /v1/permissions/rules
// @API UCS POST /v1/permissions/rules
// @API UCS DELETE /v1/permissions/rules/{id}
// @API UCS PUT /v1/permissions/rules/{id}
func ResourcePolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePolicyCreate,
		UpdateContext: resourcePolicyUpdate,
		ReadContext:   resourcePolicyRead,
		DeleteContext: resourcePolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the UCS policy.`,
			},
			"iam_user_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `Specifies the list of iam user IDs to associate to the policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the type of the UCS policy.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the UCS policy.`,
			},
			"details": {
				Type:        schema.TypeList,
				Elem:        PolicyContentSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the details of the UCS policy.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The created time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The updated time.`,
			},
		},
	}
}

func PolicyContentSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"operations": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the list of operations.`,
			},
			"resources": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the list of resources.`,
			},
		},
	}
	return &sc
}

func resourcePolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createPolicy: Create a UCS Policy.
	var (
		createPolicyHttpUrl = "v1/permissions/rules"
		createPolicyProduct = "ucs"
	)
	createPolicyClient, err := cfg.NewServiceClient(createPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	createPolicyPath := createPolicyClient.Endpoint + createPolicyHttpUrl

	createPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createPolicyOpt.JSONBody = utils.RemoveNil(buildCreatePolicyBodyParams(d))
	createPolicyResp, err := createPolicyClient.Request("POST", createPolicyPath, &createPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating Policy: %s", err)
	}

	createPolicyRespBody, err := utils.FlattenResponse(createPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("uid", createPolicyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating Policy: ID is not found in API response")
	}
	d.SetId(id)

	return resourcePolicyRead(ctx, d, meta)
}

func buildCreatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": buildPolicyMetadataBodyParams(d),
		"spec":     buildPolicySpecBodyParams(d),
	}
	return bodyParams
}

func buildPolicyMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func buildPolicySpecBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
		"iamuserids":  utils.ValueIgnoreEmpty(d.Get("iam_user_ids")),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
	}

	// when the type is admin, develop or readonly, the contents must be empty
	if d.Get("type") == "custom" {
		bodyParams["contents"] = buildPolicyRequestBodyContent(d.Get("details"))
	}
	return bodyParams
}

func buildPolicyRequestBodyContent(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"verbs":     utils.ValueIgnoreEmpty(raw["operations"]),
				"resources": utils.ValueIgnoreEmpty(raw["resources"]),
			}
		}
		return rst
	}
	return nil
}

func resourcePolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getPolicy: Query the UCS Policy detail
	var (
		getPolicyHttpUrl = "v1/permissions/rules"
		getPolicyProduct = "ucs"
	)
	getPolicyClient, err := cfg.NewServiceClient(getPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	getPolicyPath := getPolicyClient.Endpoint + getPolicyHttpUrl

	getPolicyResp, err := pagination.ListAllItems(
		getPolicyClient,
		"offset",
		getPolicyPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Policy")
	}

	getPolicyRespJson, err := json.Marshal(getPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getPolicyRespBody interface{}
	err = json.Unmarshal(getPolicyRespJson, &getPolicyRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("items[?metadata.uid=='%s']|[0]", d.Id())
	getPolicyRespBody = utils.PathSearch(jsonPath, getPolicyRespBody, nil)
	if getPolicyRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("metadata.name", getPolicyRespBody, nil)),
		d.Set("description", utils.PathSearch("spec.description", getPolicyRespBody, nil)),
		d.Set("iam_user_ids", utils.PathSearch("spec.iamUserIDs", getPolicyRespBody, nil)),
		d.Set("type", utils.PathSearch("spec.type", getPolicyRespBody, nil)),
		d.Set("created_at", utils.PathSearch("metadata.creationTimestamp", getPolicyRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("metadata.updateTimestamp", getPolicyRespBody, nil)),
		d.Set("details", flattenGetPolicyResponseBodyContent(getPolicyRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetPolicyResponseBodyContent(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("spec.contents", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"operations": utils.PathSearch("verbs", v, nil),
			"resources":  utils.PathSearch("resources", v, nil),
		})
	}
	return rst
}

func resourcePolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updatePolicyChanges := []string{
		"description",
		"iam_user_ids",
		"type",
		"details",
	}

	if d.HasChanges(updatePolicyChanges...) {
		// updatePolicy: Update the UCS Policy
		var (
			updatePolicyHttpUrl = "v1/permissions/rules/{id}"
			updatePolicyProduct = "ucs"
		)
		updatePolicyClient, err := cfg.NewServiceClient(updatePolicyProduct, region)
		if err != nil {
			return diag.Errorf("error creating UCS Client: %s", err)
		}

		updatePolicyPath := updatePolicyClient.Endpoint + updatePolicyHttpUrl
		updatePolicyPath = strings.ReplaceAll(updatePolicyPath, "{id}", d.Id())

		updatePolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updatePolicyOpt.JSONBody = utils.RemoveNil(buildUpdatePolicyBodyParams(d))
		_, err = updatePolicyClient.Request("PUT", updatePolicyPath, &updatePolicyOpt)
		if err != nil {
			return diag.Errorf("error updating Policy: %s", err)
		}
	}
	return resourcePolicyRead(ctx, d, meta)
}

func buildUpdatePolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"spec": buildPolicySpecBodyParams(d),
	}
	return bodyParams
}

func resourcePolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deletePolicy: Delete an existing UCS Policy
	var (
		deletePolicyHttpUrl = "v1/permissions/rules/{id}"
		deletePolicyProduct = "ucs"
	)
	deletePolicyClient, err := cfg.NewServiceClient(deletePolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating UCS Client: %s", err)
	}

	deletePolicyPath := deletePolicyClient.Endpoint + deletePolicyHttpUrl
	deletePolicyPath = strings.ReplaceAll(deletePolicyPath, "{id}", d.Id())

	deletePolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	_, err = deletePolicyClient.Request("DELETE", deletePolicyPath, &deletePolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting Policy: %s", err)
	}

	return nil
}
