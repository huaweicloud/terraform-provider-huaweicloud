// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DSC
// ---------------------------------------------------------------

package dsc

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

// @API DSC PUT /v1/{project_id}/sdg/asset/{id}/name
// @API DSC DELETE /v1/{project_id}/sdg/asset/obs/bucket/{id}
// @API DSC GET /v1/{project_id}/sdg/asset/obs/buckets
// @API DSC POST /v1/{project_id}/sdg/asset/obs/buckets
func ResourceAssetObs() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAssetObsCreate,
		UpdateContext: resourceAssetObsUpdate,
		ReadContext:   resourceAssetObsRead,
		DeleteContext: resourceAssetObsDelete,
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
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of asset.`,
			},
			"bucket_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The bucket name.`,
			},
			"bucket_policy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The bucket policy.`,
			},
		},
	}
}

func queryAssetObsBucket(client *golangsdk.ServiceClient) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/sdg/asset/obs/buckets"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += "?added=true"
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceAssetObsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		product = "dsc"
		name    = d.Get("name").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAssetObsBodyParams(d, cfg)),
	}
	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error adding DSC asset OBS bucket: %s", err)
	}

	respBody, err := queryAssetObsBucket(client)
	if err != nil {
		return diag.Errorf("error retrieving DSC asset OBS buckets: %s", err)
	}

	expression := fmt.Sprintf("buckets[?asset_name=='%s']|[0].id", name)
	assetObsId := utils.PathSearch(expression, respBody, "").(string)
	if assetObsId == "" {
		return diag.Errorf("error adding DSC asset OBS bucket: ID is not found in API response")
	}
	d.SetId(assetObsId)

	return resourceAssetObsRead(ctx, d, meta)
}

func buildCreateAssetObsBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"buckets": []map[string]interface{}{{
			"asset_name":    utils.ValueIgnoreEmpty(d.Get("name")),
			"location":      cfg.GetRegion(d),
			"bucket_name":   utils.ValueIgnoreEmpty(d.Get("bucket_name")),
			"bucket_policy": utils.ValueIgnoreEmpty(d.Get("bucket_policy")),
		},
		},
	}
	return bodyParams
}

func resourceAssetObsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	respBody, err := queryAssetObsBucket(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DSC asset OBS buckets")
	}

	expression := fmt.Sprintf("buckets[?id=='%s']|[0]", d.Id())
	assetObs := utils.PathSearch(expression, respBody, nil)
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("asset_name", assetObs, nil)),
		d.Set("bucket_name", utils.PathSearch("bucket_name", assetObs, nil)),
		d.Set("bucket_policy", utils.PathSearch("bucket_policy", assetObs, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAssetObsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	if d.HasChange("name") {
		var (
			httpUrl = "v1/{project_id}/sdg/asset/{id}/name"
			product = "dsc"
		)
		client, err := cfg.NewServiceClient(product, region)
		if err != nil {
			return diag.Errorf("error creating DSC client: %s", err)
		}

		requestPath := client.Endpoint + httpUrl
		requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
		requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
		requestOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateAssetObsBodyParams(d, cfg)),
		}
		_, err = client.Request("PUT", requestPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error updating DSC asset name: %s", err)
		}
	}
	return resourceAssetObsRead(ctx, d, meta)
}

func buildUpdateAssetObsBodyParams(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": utils.ValueIgnoreEmpty(d.Get("name")),
	}
	return bodyParams
}

func resourceAssetObsDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/sdg/asset/obs/bucket/{id}"
		product = "dsc"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting DSC asset OBS bucket: %s", err)
	}

	return nil
}
