// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product RMS
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

// @API Config PUT /v1/resource-manager/domains/{domain_id}/aggregators
// @API Config PUT /v1/resource-manager/domains/{domain_id}/aggregators/{id}
// @API Config GET /v1/resource-manager/domains/{domain_id}/aggregators/{id}
// @API Config DELETE /v1/resource-manager/domains/{domain_id}/aggregators/{id}
func ResourceAggregator() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAggregatorCreate,
		UpdateContext: resourceAggregatorUpdate,
		ReadContext:   resourceAggregatorRead,
		DeleteContext: resourceAggregatorDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the resource aggregator name.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the resource aggregator type, which can be **ACCOUNT** or **ORGANIZATION**.`,
			},
			"account_ids": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Computed:    true,
				Description: `Specifies the source account list.`,
			},
			"urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the resource aggregator identifier.`,
			},
		},
	}
}

func buildAggregatorBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"aggregator_name": d.Get("name"),
		"aggregator_type": d.Get("type"),
		"account_aggregation_sources": map[string]interface{}{
			"domain_ids": utils.ValueIgnoreEmpty(d.Get("account_ids").(*schema.Set).List()),
		},
	}
	return bodyParams
}

func resourceAggregatorCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		createAggregatorHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators"
		createAggregatorProduct = "rms"
	)

	cfg := meta.(*config.Config)
	createAggregatorClient, err := cfg.NewServiceClient(createAggregatorProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	createAggregatorPath := createAggregatorClient.Endpoint + createAggregatorHttpUrl
	createAggregatorPath = strings.ReplaceAll(createAggregatorPath, "{domain_id}", cfg.DomainID)

	createAggregatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	createAggregatorOpt.JSONBody = utils.RemoveNil(buildAggregatorBodyParams(d))
	createAggregatorResp, err := createAggregatorClient.Request("PUT", createAggregatorPath, &createAggregatorOpt)
	if err != nil {
		return diag.Errorf("error creating aggregator: %s", err)
	}

	createAggregatorRespBody, err := utils.FlattenResponse(createAggregatorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("aggregator_id", createAggregatorRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating aggregator: ID is not found in API response")
	}

	d.SetId(id)
	return resourceAggregatorRead(ctx, d, meta)
}

func resourceAggregatorRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		getAggregatorHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/{id}"
		getAggregatorProduct = "rms"
	)

	cfg := meta.(*config.Config)
	getAggregatorClient, err := cfg.NewServiceClient(getAggregatorProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	getAggregatorPath := getAggregatorClient.Endpoint + getAggregatorHttpUrl
	getAggregatorPath = strings.ReplaceAll(getAggregatorPath, "{domain_id}", cfg.DomainID)
	getAggregatorPath = strings.ReplaceAll(getAggregatorPath, "{id}", d.Id())

	getAggregatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAggregatorResp, err := getAggregatorClient.Request("GET", getAggregatorPath, &getAggregatorOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving aggregator")
	}

	getAggregatorRespBody, err := utils.FlattenResponse(getAggregatorResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("aggregator_name", getAggregatorRespBody, nil)),
		d.Set("type", utils.PathSearch("aggregator_type", getAggregatorRespBody, nil)),
		d.Set("account_ids", utils.PathSearch("account_aggregation_sources.domain_ids", getAggregatorRespBody, nil)),
		d.Set("urn", utils.PathSearch("aggregator_urn", getAggregatorRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAggregatorUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		updateAggregatorHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/{id}"
		updateAggregatorProduct = "rms"
	)

	cfg := meta.(*config.Config)
	updateAggregatorClient, err := cfg.NewServiceClient(updateAggregatorProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	updateAggregatorChanges := []string{
		"account_ids",
	}

	if d.HasChanges(updateAggregatorChanges...) {
		updateAggregatorPath := updateAggregatorClient.Endpoint + updateAggregatorHttpUrl
		updateAggregatorPath = strings.ReplaceAll(updateAggregatorPath, "{domain_id}", cfg.DomainID)
		updateAggregatorPath = strings.ReplaceAll(updateAggregatorPath, "{id}", d.Id())

		updateAggregatorOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateAggregatorOpt.JSONBody = utils.RemoveNil(buildAggregatorBodyParams(d))
		_, err = updateAggregatorClient.Request("PUT", updateAggregatorPath, &updateAggregatorOpt)
		if err != nil {
			return diag.Errorf("error updating aggregator: %s", err)
		}
	}
	return resourceAggregatorRead(ctx, d, meta)
}

func resourceAggregatorDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		deleteAggregatorHttpUrl = "v1/resource-manager/domains/{domain_id}/aggregators/{id}"
		deleteAggregatorProduct = "rms"
	)

	cfg := meta.(*config.Config)
	deleteAggregatorClient, err := cfg.NewServiceClient(deleteAggregatorProduct, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating RMS Client: %s", err)
	}

	deleteAggregatorPath := deleteAggregatorClient.Endpoint + deleteAggregatorHttpUrl
	deleteAggregatorPath = strings.ReplaceAll(deleteAggregatorPath, "{domain_id}", cfg.DomainID)
	deleteAggregatorPath = strings.ReplaceAll(deleteAggregatorPath, "{id}", d.Id())

	deleteAggregatorOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteAggregatorClient.Request("DELETE", deleteAggregatorPath, &deleteAggregatorOpt)
	if err != nil {
		return diag.Errorf("error deleting aggregator: %s", err)
	}

	return nil
}
