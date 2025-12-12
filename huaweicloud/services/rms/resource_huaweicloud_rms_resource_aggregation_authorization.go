// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RMS
// ---------------------------------------------------------------

package rms

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

// @API Config PUT /v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization
// @API Config GET /v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization/{id}
func ResourceAggregationAuthorization() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAggregationAuthCreate,
		ReadContext:   resourceAggregationAuthRead,
		DeleteContext: resourceAggregationAuthDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"account_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the resource aggregation account to be authorized.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the authorization identifier of the resource aggregation account.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the resource aggregation account was authorized.`,
			},
			"tags": common.TagsSchema(),
		},
	}
}

func buildCreateAggregationAuthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"authorized_account_id": d.Get("account_id"),
	}
	if tagMap := d.Get("tags").(map[string]interface{}); len(tagMap) > 0 {
		bodyParams["tags"] = utils.ExpandResourceTags(tagMap)
	}
	return bodyParams
}

func resourceAggregationAuthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		createAggregationAuthHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization"
		createAggregationAuthProduct = "rms"
	)

	cfg := meta.(*config.Config)
	createAggregationAuthClient, err := cfg.NewServiceClient(createAggregationAuthProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	createAggregationAuthPath := createAggregationAuthClient.Endpoint + createAggregationAuthHttpUrl
	createAggregationAuthPath = strings.ReplaceAll(createAggregationAuthPath, "{domain_id}", cfg.DomainID)

	createAggregationAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAggregationAuthOpt.JSONBody = buildCreateAggregationAuthBodyParams(d)
	createAggregationAuthResp, err := createAggregationAuthClient.Request("PUT", createAggregationAuthPath, &createAggregationAuthOpt)
	if err != nil {
		return diag.Errorf("error creating aggregation authorization: %s", err)
	}

	createAggregationAuthRespBody, err := utils.FlattenResponse(createAggregationAuthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("authorized_account_id", createAggregationAuthRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating aggregation authorization: ID is not found in API response")
	}

	d.SetId(id)
	return resourceAggregationAuthRead(ctx, d, meta)
}

func resourceAggregationAuthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		getAggregationAuthHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization"
		getAggregationAuthProduct = "rms"
	)

	cfg := meta.(*config.Config)
	getAggregationAuthClient, err := cfg.NewServiceClient(getAggregationAuthProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	getAggregationAuthPath := getAggregationAuthClient.Endpoint + getAggregationAuthHttpUrl
	getAggregationAuthPath = strings.ReplaceAll(getAggregationAuthPath, "{domain_id}", cfg.DomainID)
	getAggregationAuthPath += fmt.Sprintf("?account_id=%s", d.Id())

	getAggregationAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAggregationAuthResp, err := getAggregationAuthClient.Request("GET", getAggregationAuthPath, &getAggregationAuthOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving aggregation authorization")
	}

	respBody, err := utils.FlattenResponse(getAggregationAuthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	item := utils.PathSearch("aggregation_authorizations[0]", respBody, nil)
	if item == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving aggregation authorization")
	}

	mErr := multierror.Append(nil,
		d.Set("urn", utils.PathSearch("aggregation_authorization_urn", item, nil)),
		d.Set("account_id", utils.PathSearch("authorized_account_id", item, nil)),
		d.Set("created_at", utils.PathSearch("created_at", item, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", item, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAggregationAuthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		deleteAggregationAuthHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/aggregation-authorization/{id}"
		deleteAggregationAuthProduct = "rms"
	)

	cfg := meta.(*config.Config)
	deleteAggregationAuthClient, err := cfg.NewServiceClient(deleteAggregationAuthProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	deleteAggregationAuthPath := deleteAggregationAuthClient.Endpoint + deleteAggregationAuthHttpUrl
	deleteAggregationAuthPath = strings.ReplaceAll(deleteAggregationAuthPath, "{domain_id}", cfg.DomainID)
	deleteAggregationAuthPath = strings.ReplaceAll(deleteAggregationAuthPath, "{id}", d.Id())

	deleteAggregationAuthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteAggregationAuthClient.Request("DELETE", deleteAggregationAuthPath, &deleteAggregationAuthOpt)
	if err != nil {
		return diag.Errorf("error deleting aggregation authorization: %s", err)
	}

	return nil
}
