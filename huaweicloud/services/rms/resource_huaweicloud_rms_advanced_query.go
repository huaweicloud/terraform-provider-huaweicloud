// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product Config
// ---------------------------------------------------------------

package rms

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Config POST /v1/resource-manager/domains/{domain_id}/stored-queries
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}
// @API Config GET /v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}
// @API Config PUT /v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}
func ResourceAdvancedQuery() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAdvancedQueryCreate,
		UpdateContext: resourceAdvancedQueryUpdate,
		ReadContext:   resourceAdvancedQueryRead,
		DeleteContext: resourceAdvancedQueryDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ResourceQL name.`,
			},
			"expression": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ResourceQL expression.`,
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the ResourceQL type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ResourceQL description.`,
			},
		},
	}
}

func resourceAdvancedQueryCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAdvancedQuery: Create a RMS advanced query.
	var (
		createAdvancedQueryHttpUrl = "v1/resource-manager/domains/{domain_id}/stored-queries"
		createAdvancedQueryProduct = "rms"
	)
	createAdvancedQueryClient, err := cfg.NewServiceClient(createAdvancedQueryProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	createAdvancedQueryPath := createAdvancedQueryClient.Endpoint + createAdvancedQueryHttpUrl
	createAdvancedQueryPath = strings.ReplaceAll(createAdvancedQueryPath, "{domain_id}", cfg.DomainID)

	createAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createAdvancedQueryOpt.JSONBody = utils.RemoveNil(buildCreateAdvancedQueryBodyParams(d))
	createAdvancedQueryResp, err := createAdvancedQueryClient.Request("POST", createAdvancedQueryPath,
		&createAdvancedQueryOpt)
	if err != nil {
		return diag.Errorf("error creating RMS advanced query: %s", err)
	}

	createAdvancedQueryRespBody, err := utils.FlattenResponse(createAdvancedQueryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createAdvancedQueryRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating RMS advanced query: ID is not found in API response")
	}
	d.SetId(id)

	return resourceAdvancedQueryRead(ctx, d, meta)
}

func buildCreateAdvancedQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"expression":  d.Get("expression"),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceAdvancedQueryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAdvancedQuery: Query the RMS advanced query
	var (
		getAdvancedQueryHttpUrl = "v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}"
		getAdvancedQueryProduct = "rms"
	)
	getAdvancedQueryClient, err := cfg.NewServiceClient(getAdvancedQueryProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	getAdvancedQueryPath := getAdvancedQueryClient.Endpoint + getAdvancedQueryHttpUrl
	getAdvancedQueryPath = strings.ReplaceAll(getAdvancedQueryPath, "{domain_id}", cfg.DomainID)
	getAdvancedQueryPath = strings.ReplaceAll(getAdvancedQueryPath, "{query_id}", d.Id())

	getAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAdvancedQueryResp, err := getAdvancedQueryClient.Request("GET", getAdvancedQueryPath, &getAdvancedQueryOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving RMS advanced query")
	}

	getAdvancedQueryRespBody, err := utils.FlattenResponse(getAdvancedQueryResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("name", utils.PathSearch("name", getAdvancedQueryRespBody, nil)),
		d.Set("expression", utils.PathSearch("expression", getAdvancedQueryRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getAdvancedQueryRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getAdvancedQueryRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAdvancedQueryUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAdvancedQueryChanges := []string{
		"name",
		"expression",
		"type",
		"description",
	}

	if d.HasChanges(updateAdvancedQueryChanges...) {
		// updateAdvancedQuery: Update a RMS advanced query.
		var (
			updateAdvancedQueryHttpUrl = "v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}"
			updateAdvancedQueryProduct = "rms"
		)
		updateAdvancedQueryClient, err := cfg.NewServiceClient(updateAdvancedQueryProduct, region)
		if err != nil {
			return diag.Errorf("error creating RMS client: %s", err)
		}

		updateAdvancedQueryPath := updateAdvancedQueryClient.Endpoint + updateAdvancedQueryHttpUrl
		updateAdvancedQueryPath = strings.ReplaceAll(updateAdvancedQueryPath, "{domain_id}", cfg.DomainID)
		updateAdvancedQueryPath = strings.ReplaceAll(updateAdvancedQueryPath, "{query_id}", d.Id())

		updateAdvancedQueryOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}

		updateAdvancedQueryOpt.JSONBody = utils.RemoveNil(buildUpdateAdvancedQueryBodyParams(d))
		_, err = updateAdvancedQueryClient.Request("PUT", updateAdvancedQueryPath, &updateAdvancedQueryOpt)
		if err != nil {
			return diag.Errorf("error updating RMS advanced query: %s", err)
		}
	}
	return resourceAdvancedQueryRead(ctx, d, meta)
}

func buildUpdateAdvancedQueryBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"expression":  d.Get("expression"),
		"type":        utils.ValueIgnoreEmpty(d.Get("type")),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceAdvancedQueryDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAdvancedQuery: Delete an existing RMS advanced query
	var (
		deleteAdvancedQueryHttpUrl = "v1/resource-manager/domains/{domain_id}/stored-queries/{query_id}"
		deleteAdvancedQueryProduct = "rms"
	)
	deleteAdvancedQueryClient, err := cfg.NewServiceClient(deleteAdvancedQueryProduct, region)
	if err != nil {
		return diag.Errorf("error creating RMS client: %s", err)
	}

	deleteAdvancedQueryPath := deleteAdvancedQueryClient.Endpoint + deleteAdvancedQueryHttpUrl
	deleteAdvancedQueryPath = strings.ReplaceAll(deleteAdvancedQueryPath, "{domain_id}", cfg.DomainID)
	deleteAdvancedQueryPath = strings.ReplaceAll(deleteAdvancedQueryPath, "{query_id}", d.Id())

	deleteAdvancedQueryOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = deleteAdvancedQueryClient.Request("DELETE", deleteAdvancedQueryPath, &deleteAdvancedQueryOpt)
	if err != nil {
		return diag.Errorf("error deleting RMS advanced query: %s", err)
	}

	return nil
}
