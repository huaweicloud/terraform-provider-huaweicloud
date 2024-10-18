// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product GES
// ---------------------------------------------------------------

package ges

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

// @API GES DELETE /v2/{project_id}/graphs/metadatas/{id}
// @API GES GET /v1.0/{project_id}/graphs/metadatas/{id}
// @API GES GET /v2/{project_id}/graphs/metadatas
// @API GES POST /v2/{project_id}/graphs/metadatas
func ResourceGesMetadata() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGesMetadataCreate,
		UpdateContext: resourceGesMetadataUpdate,
		ReadContext:   resourceGesMetadataRead,
		DeleteContext: resourceGesMetadataDelete,
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
				ForceNew:    true,
				Description: `Metadata name.`,
			},
			"metadata_path": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `OBS Path for storing the metadata.`,
			},
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Metadata description.`,
			},
			"ges_metadata": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     metadataSchema(),
				Required: true,
			},
			"encryption": {
				Type:     schema.TypeList,
				MaxItems: 1,
				Elem:     metadataEncryptionSchema(),
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Status of a metadata. **200** is available.`,
			},
		},
	}
}

func metadataSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"labels": {
				Type:        schema.TypeList,
				Elem:        metadataLabelsSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Label list`,
			},
		},
	}
	return &sc
}

func metadataLabelsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Name of a label.`,
			},
			"properties": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
				},
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func metadataEncryptionSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:        schema.TypeBool,
				Optional:    true,
				Computed:    true,
				Description: `Whether to enable data encryption The value can be true or false. The default value is false.`,
			},
			"master_key_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `ID of the customer master key created by DEW in the project where the graph is created.`,
			},
		},
	}
	return &sc
}

func resourceGesMetadataCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createMetadata: create a GES metadata.
	var (
		createMetadataHttpUrl = "v2/{project_id}/graphs/metadatas"
		createMetadataProduct = "ges"
	)
	createMetadataClient, err := cfg.NewServiceClient(createMetadataProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	createMetadataPath := createMetadataClient.Endpoint + createMetadataHttpUrl
	createMetadataPath = strings.ReplaceAll(createMetadataPath, "{project_id}", createMetadataClient.ProjectID)

	createMetadataOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	createMetadataOpt.JSONBody = utils.RemoveNil(buildMetadataBodyParams(d))
	createMetadataResp, err := createMetadataClient.Request("POST", createMetadataPath, &createMetadataOpt)
	if err != nil {
		return diag.Errorf("error creating GES metadata: %s", err)
	}

	createMetadataRespBody, err := utils.FlattenResponse(createMetadataResp)
	if err != nil {
		return diag.FromErr(err)
	}

	metadataId := utils.PathSearch("id", createMetadataRespBody, "").(string)
	if metadataId == "" {
		return diag.Errorf("unable to find the GES metadata ID from the API response")
	}
	d.SetId(metadataId)

	return resourceGesMetadataRead(ctx, d, meta)
}

func buildMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"metadata_path": utils.ValueIgnoreEmpty(d.Get("metadata_path")),
		"description":   utils.ValueIgnoreEmpty(d.Get("description")),
		"is_overwrite":  false,
		"ges_metadata":  buildMetadataReqBodyMetadata(d.Get("ges_metadata")),
		"encryption":    buildMetadataReqBodyEncryption(d.Get("encryption")),
	}
	return bodyParams
}

func buildMetadataReqBodyMetadata(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"labels": buildMetadataLabels(raw["labels"]),
		}
		return params
	}
	return nil
}

func buildMetadataReqBodyEncryption(rawParams interface{}) map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}
		raw := rawArray[0].(map[string]interface{})
		params := map[string]interface{}{
			"enable":        utils.ValueIgnoreEmpty(raw["enable"]),
			"master_key_id": utils.ValueIgnoreEmpty(raw["master_key_id"]),
		}
		return params
	}
	return nil
}

func buildMetadataLabels(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		rst := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			rst[i] = map[string]interface{}{
				"name":       utils.ValueIgnoreEmpty(raw["name"]),
				"properties": raw["properties"],
			}
		}
		return rst
	}
	return nil
}

func resourceGesMetadataRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getMetadataDetail: Query the GES metadata detail.
	var (
		getMetadataDetailHttpUrl = "v1.0/{project_id}/graphs/metadatas/{id}"
		getMetadataDetailProduct = "ges"
	)
	getMetadataDetailClient, err := cfg.NewServiceClient(getMetadataDetailProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	getMetadataDetailPath := getMetadataDetailClient.Endpoint + getMetadataDetailHttpUrl
	getMetadataDetailPath = strings.ReplaceAll(getMetadataDetailPath, "{project_id}", getMetadataDetailClient.ProjectID)
	getMetadataDetailPath = strings.ReplaceAll(getMetadataDetailPath, "{id}", d.Id())

	getMetadataDetailOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getMetadataDetailResp, err := getMetadataDetailClient.Request("GET", getMetadataDetailPath, &getMetadataDetailOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errCode", "GES.2067"),
			"error retrieving GES metadata")
	}

	getMetadataDetailRespBody, err := utils.FlattenResponse(getMetadataDetailResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("ges_metadata", flattenGetMetadataRespBodyMetadata(getMetadataDetailRespBody)),
	)

	// getMetadataList: Query the GES metadata List.
	var (
		getMetadataListHttpUrl = "v2/{project_id}/graphs/metadatas"
		getMetadataListProduct = "ges"
	)
	getMetadataListClient, err := cfg.NewServiceClient(getMetadataListProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	getMetadataListPath := getMetadataListClient.Endpoint + getMetadataListHttpUrl
	getMetadataListPath = strings.ReplaceAll(getMetadataListPath, "{project_id}", getMetadataListClient.ProjectID)

	getMetadataListPath += "?limit=100" // user only can add 50 metadatas.

	getMetadataListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	getMetadataListResp, err := getMetadataListClient.Request("GET", getMetadataListPath, &getMetadataListOpt)
	if err != nil {
		return diag.FromErr(err)
	}

	getMetadataListRespBody, err := utils.FlattenResponse(getMetadataListResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("schema_list[?id=='%s']|[0]", d.Id())
	getMetadataListRespBody = utils.PathSearch(jsonPath, getMetadataListRespBody, nil)
	if getMetadataListRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getMetadataListRespBody, nil)),
		d.Set("metadata_path", utils.PathSearch("metadata_path", getMetadataListRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getMetadataListRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getMetadataListRespBody, nil)),
		d.Set("encryption", flattenGetMetadatasRespBodyEncryption(getMetadataListRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetMetadatasRespBodyEncryption(resp interface{}) []interface{} {
	rst := []interface{}{
		map[string]interface{}{
			"enable":        utils.PathSearch("encrypted", resp, nil),
			"master_key_id": utils.PathSearch("master_key_id", resp, nil),
		},
	}
	return rst
}

func flattenGetMetadataRespBodyMetadata(resp interface{}) []interface{} {
	var rst []interface{}
	curJson := utils.PathSearch("gesMetadata", resp, nil)
	if curJson == nil {
		return rst
	}

	rst = []interface{}{
		map[string]interface{}{
			"labels": flattenMetadataLabels(curJson),
		},
	}
	return rst
}

func flattenMetadataLabels(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("labels", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":       utils.PathSearch("name", v, nil),
			"properties": utils.PathSearch("properties", v, nil),
		})
	}
	return rst
}

func resourceGesMetadataUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateMetadataChanges := []string{
		"metadata_path",
		"ges_metadata",
		"encryption",
	}

	if d.HasChanges(updateMetadataChanges...) {
		// updateMetadata: update metadata
		var (
			updateMetadataHttpUrl = "v2/{project_id}/graphs/metadatas"
			updateMetadataProduct = "ges"
		)
		updateMetadataClient, err := cfg.NewServiceClient(updateMetadataProduct, region)
		if err != nil {
			return diag.Errorf("error creating GES Client: %s", err)
		}

		updateMetadataPath := updateMetadataClient.Endpoint + updateMetadataHttpUrl
		updateMetadataPath = strings.ReplaceAll(updateMetadataPath, "{project_id}", updateMetadataClient.ProjectID)

		updateMetadataOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		}

		updateMetadataOpt.JSONBody = utils.RemoveNil(buildUpdateMetadataBodyParams(d))
		_, err = updateMetadataClient.Request("POST", updateMetadataPath, &updateMetadataOpt)
		if err != nil {
			return diag.Errorf("error updating GES metadata: %s", err)
		}
	}
	return resourceGesMetadataRead(ctx, d, meta)
}

func buildUpdateMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"metadata_path": utils.ValueIgnoreEmpty(d.Get("metadata_path")),
		"description":   utils.ValueIgnoreEmpty(d.Get("description")),
		"is_overwrite":  true,
		"ges_metadata":  buildMetadataReqBodyMetadata(d.Get("ges_metadata")),
		"encryption":    buildMetadataReqBodyEncryption(d.Get("encryption")),
	}
	return bodyParams
}

func resourceGesMetadataDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteMetadata: delete GES metadata
	var (
		deleteMetadataHttpUrl = "v2/{project_id}/graphs/metadatas/{id}"
		deleteMetadataProduct = "ges"
	)
	deleteMetadataClient, err := cfg.NewServiceClient(deleteMetadataProduct, region)
	if err != nil {
		return diag.Errorf("error creating GES Client: %s", err)
	}

	deleteMetadataPath := deleteMetadataClient.Endpoint + deleteMetadataHttpUrl
	deleteMetadataPath = strings.ReplaceAll(deleteMetadataPath, "{project_id}", deleteMetadataClient.ProjectID)
	deleteMetadataPath = strings.ReplaceAll(deleteMetadataPath, "{id}", d.Id())

	deleteMetadataOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	_, err = deleteMetadataClient.Request("DELETE", deleteMetadataPath, &deleteMetadataOpt)
	if err != nil {
		return diag.Errorf("error deleting GES metadata: %s", err)
	}

	return nil
}
