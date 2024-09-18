// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

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

// @API CC POST /v3/{domain_id}/ccaas/inter-region-bandwidths
// @API CC DELETE /v3/{domain_id}/ccaas/inter-region-bandwidths/{id}
// @API CC GET /v3/{domain_id}/ccaas/inter-region-bandwidths/{id}
// @API CC PUT /v3/{domain_id}/ccaas/inter-region-bandwidths/{id}
func ResourceInterRegionBandwidth() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInterRegionBandwidthCreate,
		UpdateContext: resourceInterRegionBandwidthUpdate,
		ReadContext:   resourceInterRegionBandwidthRead,
		DeleteContext: resourceInterRegionBandwidthDelete,
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
			"cloud_connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Cloud connection ID.`,
			},
			"bandwidth_package_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Bandwidth package ID.`,
			},
			"bandwidth": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Inter-region bandwidth.`,
			},
			"inter_region_ids": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				ForceNew:    true,
				Description: `Two regions to which bandwidth is allocated.`,
			},
			"inter_regions": {
				Type:     schema.TypeList,
				Elem:     InterRegionBandwidthInterRegionsSchema(),
				Computed: true,
			},
		},
	}
}

func InterRegionBandwidthInterRegionsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Inter-region bandwidth ID.`,
			},
			"project_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Project ID of a region where the inter-region bandwidth is used.`,
			},
			"local_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `ID of the local region where the inter-region bandwidth is used.`,
			},
			"remote_region_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `ID of the remote region where the inter-region bandwidth is used.`,
			},
		},
	}
	return &sc
}

func resourceInterRegionBandwidthCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createInterRegionBandwidthHttpUrl = "v3/{domain_id}/ccaas/inter-region-bandwidths"
		createInterRegionBandwidthProduct = "cc"
	)
	createInterRegionBandwidthClient, err := cfg.NewServiceClient(createInterRegionBandwidthProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	createInterRegionBandwidthPath := createInterRegionBandwidthClient.Endpoint + createInterRegionBandwidthHttpUrl
	createInterRegionBandwidthPath = strings.ReplaceAll(createInterRegionBandwidthPath, "{domain_id}", cfg.DomainID)

	createInterRegionBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}

	createInterRegionBandwidthOpt.JSONBody = utils.RemoveNil(buildCreateInterRegionBandwidthBodyParams(d))
	createInterRegionBandwidthResp, err := createInterRegionBandwidthClient.Request("POST", createInterRegionBandwidthPath,
		&createInterRegionBandwidthOpt)
	if err != nil {
		return diag.Errorf("error creating inter-region bandwidth: %s", err)
	}

	createInterRegionBandwidthRespBody, err := utils.FlattenResponse(createInterRegionBandwidthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("inter_region_bandwidth.id", createInterRegionBandwidthRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating inter-region bandwidth: ID is not found in API response")
	}
	d.SetId(id)

	return resourceInterRegionBandwidthRead(ctx, d, meta)
}

func buildCreateInterRegionBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"inter_region_bandwidth": map[string]interface{}{
			"cloud_connection_id":  utils.ValueIgnoreEmpty(d.Get("cloud_connection_id")),
			"bandwidth_package_id": utils.ValueIgnoreEmpty(d.Get("bandwidth_package_id")),
			"bandwidth":            utils.ValueIgnoreEmpty(d.Get("bandwidth")),
			"inter_region_ids":     utils.ValueIgnoreEmpty(d.Get("inter_region_ids")),
		},
	}
	return bodyParams
}

func resourceInterRegionBandwidthRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getInterRegionBandwidthHttpUrl = "v3/{domain_id}/ccaas/inter-region-bandwidths/{id}"
		getInterRegionBandwidthProduct = "cc"
	)
	getInterRegionBandwidthClient, err := cfg.NewServiceClient(getInterRegionBandwidthProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	getInterRegionBandwidthPath := getInterRegionBandwidthClient.Endpoint + getInterRegionBandwidthHttpUrl
	getInterRegionBandwidthPath = strings.ReplaceAll(getInterRegionBandwidthPath, "{domain_id}", cfg.DomainID)
	getInterRegionBandwidthPath = strings.ReplaceAll(getInterRegionBandwidthPath, "{id}", d.Id())

	getInterRegionBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}

	getInterRegionBandwidthResp, err := getInterRegionBandwidthClient.Request("GET", getInterRegionBandwidthPath, &getInterRegionBandwidthOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving inter-region bandwidth")
	}

	getInterRegionBandwidthRespBody, err := utils.FlattenResponse(getInterRegionBandwidthResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("bandwidth_package_id", utils.PathSearch("inter_region_bandwidth.bandwidth_package_id", getInterRegionBandwidthRespBody, nil)),
		d.Set("cloud_connection_id", utils.PathSearch("inter_region_bandwidth.cloud_connection_id", getInterRegionBandwidthRespBody, nil)),
		d.Set("inter_region_ids", utils.PathSearch("inter_region_bandwidth.inter_regions[*].local_region_id", getInterRegionBandwidthRespBody, nil)),
		d.Set("bandwidth", utils.PathSearch("inter_region_bandwidth.bandwidth", getInterRegionBandwidthRespBody, nil)),
		d.Set("inter_regions", flattenGetInterRegionBandwidthResponseBodyInterRegions(getInterRegionBandwidthRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetInterRegionBandwidthResponseBodyInterRegions(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("inter_region_bandwidth.inter_regions", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"project_id":       utils.PathSearch("project_id", v, nil),
			"local_region_id":  utils.PathSearch("local_region_id", v, nil),
			"remote_region_id": utils.PathSearch("remote_region_id", v, nil),
		})
	}
	return rst
}

func resourceInterRegionBandwidthUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateInterRegionBandwidthChanges := []string{
		"bandwidth",
	}

	if d.HasChanges(updateInterRegionBandwidthChanges...) {
		var (
			updateInterRegionBandwidthHttpUrl = "v3/{domain_id}/ccaas/inter-region-bandwidths/{id}"
			updateInterRegionBandwidthProduct = "cc"
		)
		updateInterRegionBandwidthClient, err := cfg.NewServiceClient(updateInterRegionBandwidthProduct, region)
		if err != nil {
			return diag.Errorf("error creating CC Client: %s", err)
		}

		updateInterRegionBandwidthPath := updateInterRegionBandwidthClient.Endpoint + updateInterRegionBandwidthHttpUrl
		updateInterRegionBandwidthPath = strings.ReplaceAll(updateInterRegionBandwidthPath, "{domain_id}", cfg.DomainID)
		updateInterRegionBandwidthPath = strings.ReplaceAll(updateInterRegionBandwidthPath, "{id}", d.Id())

		updateInterRegionBandwidthOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}

		updateInterRegionBandwidthOpt.JSONBody = utils.RemoveNil(buildUpdateInterRegionBandwidthBodyParams(d))
		_, err = updateInterRegionBandwidthClient.Request("PUT", updateInterRegionBandwidthPath, &updateInterRegionBandwidthOpt)
		if err != nil {
			return diag.Errorf("error updating inter-region bandwidth: %s", err)
		}
	}
	return resourceInterRegionBandwidthRead(ctx, d, meta)
}

func buildUpdateInterRegionBandwidthBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"inter_region_bandwidth": map[string]interface{}{
			"bandwidth": utils.ValueIgnoreEmpty(d.Get("bandwidth")),
		},
	}
	return bodyParams
}

func resourceInterRegionBandwidthDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteInterRegionBandwidthHttpUrl = "v3/{domain_id}/ccaas/inter-region-bandwidths/{id}"
		deleteInterRegionBandwidthProduct = "cc"
	)
	deleteInterRegionBandwidthClient, err := cfg.NewServiceClient(deleteInterRegionBandwidthProduct, region)
	if err != nil {
		return diag.Errorf("error creating CC Client: %s", err)
	}

	deleteInterRegionBandwidthPath := deleteInterRegionBandwidthClient.Endpoint + deleteInterRegionBandwidthHttpUrl
	deleteInterRegionBandwidthPath = strings.ReplaceAll(deleteInterRegionBandwidthPath, "{domain_id}", cfg.DomainID)
	deleteInterRegionBandwidthPath = strings.ReplaceAll(deleteInterRegionBandwidthPath, "{id}", d.Id())

	deleteInterRegionBandwidthOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}

	_, err = deleteInterRegionBandwidthClient.Request("DELETE", deleteInterRegionBandwidthPath, &deleteInterRegionBandwidthOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting inter-region bandwidth")
	}

	return nil
}
