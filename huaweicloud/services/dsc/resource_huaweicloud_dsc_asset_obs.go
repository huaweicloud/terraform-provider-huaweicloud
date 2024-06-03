// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DSC
// ---------------------------------------------------------------

package dsc

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

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

func resourceAssetObsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAssetObs: create an OBS asset.
	var (
		createAssetObsHttpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		createAssetObsProduct = "dsc"
	)
	createAssetObsClient, err := cfg.NewServiceClient(createAssetObsProduct, region)
	if err != nil {
		return diag.Errorf("error creating AssetObs Client: %s", err)
	}

	createAssetObsPath := createAssetObsClient.Endpoint + createAssetObsHttpUrl
	createAssetObsPath = strings.ReplaceAll(createAssetObsPath, "{project_id}", createAssetObsClient.ProjectID)

	createAssetObsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAssetObsOpt.JSONBody = utils.RemoveNil(buildCreateAssetObsBodyParams(d, cfg))
	_, err = createAssetObsClient.Request("POST", createAssetObsPath, &createAssetObsOpt)
	if err != nil {
		return diag.Errorf("error creating AssetObs: %s", err)
	}

	// getAssetObs: Query the asset OBS
	var (
		getAssetObsHttpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		getAssetObsProduct = "dsc"
	)
	getAssetObsClient, err := cfg.NewServiceClient(getAssetObsProduct, region)
	if err != nil {
		return diag.Errorf("error creating AssetObs Client: %s", err)
	}

	getAssetObsPath := getAssetObsClient.Endpoint + getAssetObsHttpUrl
	getAssetObsPath = strings.ReplaceAll(getAssetObsPath, "{project_id}", getAssetObsClient.ProjectID)
	getAssetObsPath += "?added=true"

	getAssetObsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAssetObsResp, err := getAssetObsClient.Request("GET", getAssetObsPath, &getAssetObsOpt)
	if err != nil {
		return diag.Errorf("error creating AssetObs: %s", err)
	}

	getAssetObsRespBody, err := utils.FlattenResponse(getAssetObsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	name := d.Get("name").(string)
	id, err := jmespath.Search("id", FilterAssetObs(getAssetObsRespBody, name, "asset_name"))
	if err != nil {
		return diag.Errorf("error creating AssetObs: ID is not found in API response")
	}
	d.SetId(id.(string))

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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAssetObs: Query the asset OBS
	var (
		getAssetObsHttpUrl = "v1/{project_id}/sdg/asset/obs/buckets"
		getAssetObsProduct = "dsc"
	)
	getAssetObsClient, err := cfg.NewServiceClient(getAssetObsProduct, region)
	if err != nil {
		return diag.Errorf("error creating AssetObs Client: %s", err)
	}

	getAssetObsPath := getAssetObsClient.Endpoint + getAssetObsHttpUrl
	getAssetObsPath = strings.ReplaceAll(getAssetObsPath, "{project_id}", getAssetObsClient.ProjectID)
	getAssetObsPath += "?added=true"

	getAssetObsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAssetObsResp, err := getAssetObsClient.Request("GET", getAssetObsPath, &getAssetObsOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AssetObs")
	}

	getAssetObsRespBody, err := utils.FlattenResponse(getAssetObsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	assetObs := FilterAssetObs(getAssetObsRespBody, d.Id(), "id")

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("asset_name", assetObs, nil)),
		d.Set("bucket_name", utils.PathSearch("bucket_name", assetObs, nil)),
		d.Set("bucket_policy", utils.PathSearch("bucket_policy", assetObs, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func FilterAssetObs(resp interface{}, val string, path string) interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("buckets", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	for _, v := range curArray {
		if val == utils.PathSearch(path, v, nil) {
			return v
		}
	}
	return nil
}

func resourceAssetObsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAssetObshasChanges := []string{
		"name",
	}

	if d.HasChanges(updateAssetObshasChanges...) {
		// updateAssetObs: update the OBS asset name.
		var (
			updateAssetObsHttpUrl = "v1/{project_id}/sdg/asset/{id}/name"
			updateAssetObsProduct = "dsc"
		)
		updateAssetObsClient, err := cfg.NewServiceClient(updateAssetObsProduct, region)
		if err != nil {
			return diag.Errorf("error creating AssetObs Client: %s", err)
		}

		updateAssetObsPath := updateAssetObsClient.Endpoint + updateAssetObsHttpUrl
		updateAssetObsPath = strings.ReplaceAll(updateAssetObsPath, "{project_id}", updateAssetObsClient.ProjectID)
		updateAssetObsPath = strings.ReplaceAll(updateAssetObsPath, "{id}", d.Id())

		updateAssetObsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateAssetObsOpt.JSONBody = utils.RemoveNil(buildUpdateAssetObsBodyParams(d, cfg))
		_, err = updateAssetObsClient.Request("PUT", updateAssetObsPath, &updateAssetObsOpt)
		if err != nil {
			return diag.Errorf("error updating AssetObs: %s", err)
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
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteAssetObsHttpUrl = "v1/{project_id}/sdg/asset/obs/bucket/{id}"
		deleteAssetObsProduct = "dsc"
	)
	deleteAssetObsClient, err := cfg.NewServiceClient(deleteAssetObsProduct, region)
	if err != nil {
		return diag.Errorf("error creating AssetObs Client: %s", err)
	}

	deleteAssetObsPath := deleteAssetObsClient.Endpoint + deleteAssetObsHttpUrl
	deleteAssetObsPath = strings.ReplaceAll(deleteAssetObsPath, "{project_id}", deleteAssetObsClient.ProjectID)
	deleteAssetObsPath = strings.ReplaceAll(deleteAssetObsPath, "{id}", d.Id())

	deleteAssetObsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteAssetObsClient.Request("DELETE", deleteAssetObsPath, &deleteAssetObsOpt)
	if err != nil {
		return diag.Errorf("error deleting AssetObs: %s", err)
	}

	return nil
}
