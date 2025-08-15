package coc

import (
	"context"
	"fmt"
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

var documentNonUpdatableParams = []string{"enterprise_project_id", "tags"}

// @API COC POST /v1/documents
// @API COC PUT /v1/documents/{document_id}
// @API COC GET /v1/documents/{document_id}
// @API COC DELETE /v1/documents/{document_id}
func ResourceDocument() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDocumentCreate,
		ReadContext:   resourceDocumentRead,
		UpdateContext: resourceDocumentUpdate,
		DeleteContext: resourceDocumentDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: customdiff.All(
			config.FlexibleForceNew(documentNonUpdatableParams),
			config.MergeDefaultTags(),
		),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"risk_level": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"modifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"version_uuid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func resourceDocumentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	createHttpUrl := "v1/documents"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateDocumentBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC document: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening COC document: %s", err)
	}

	id := utils.PathSearch("data", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC document ID from the API response")
	}

	d.SetId(id)

	return resourceDocumentRead(ctx, d, meta)
}

func buildCreateDocumentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":                  d.Get("name"),
		"content":               d.Get("content"),
		"enterprise_project_id": d.Get("enterprise_project_id"),
		"risk_level":            d.Get("risk_level"),
		"description":           utils.ValueIgnoreEmpty(d.Get("description")),
		"tags":                  utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}

	return bodyParams
}

func resourceDocumentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	document, err := GetDocument(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00042101"),
			"error retrieving document")
	}

	mErr = multierror.Append(mErr,
		d.Set("name", utils.PathSearch("data.name", document, nil)),
		d.Set("content", utils.PathSearch("data.content", document, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("data.enterprise_project_id", document, nil)),
		d.Set("risk_level", utils.PathSearch("data.risk_level", document, nil)),
		d.Set("description", utils.PathSearch("data.description", document, nil)),
		d.Set("create_time", utils.PathSearch("data.create_time", document, nil)),
		d.Set("update_time", utils.PathSearch("data.update_time", document, nil)),
		d.Set("version", utils.PathSearch("data.version", document, nil)),
		d.Set("creator", utils.PathSearch("data.creator", document, nil)),
		d.Set("modifier", utils.PathSearch("data.modifier", document, nil)),
		d.Set("versions", flattenDocumentVersions(
			utils.PathSearch("data.versions", document, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDocumentVersions(rawParams interface{}) []map[string]interface{} {
	if paramsList, ok := rawParams.([]interface{}); ok {
		if len(paramsList) < 1 {
			return nil
		}
		configurations := make([]map[string]interface{}, len(paramsList))
		for i, params := range paramsList {
			raw := params.(map[string]interface{})
			configurations[i] = map[string]interface{}{
				"version":      utils.PathSearch("version", raw, nil),
				"version_uuid": utils.PathSearch("version_uuid", raw, nil),
				"create_time":  utils.PathSearch("create_time", raw, nil),
			}
		}
		return configurations
	}

	return nil
}

func GetDocument(client *golangsdk.ServiceClient, documentID string) (interface{}, error) {
	getHttpUrl := "v1/documents/{document_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{document_id}", documentID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening document: %s", err)
	}

	return getRespBody, nil
}

func resourceDocumentUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	updateHttpUrl := "v1/documents/{document_id}"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{document_id}", d.Id())
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildUpdateDocumentBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating document: %s", err)
	}

	return resourceDocumentRead(ctx, d, meta)
}

func buildUpdateDocumentBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"content":     d.Get("content"),
		"description": d.Get("description"),
		"risk_level":  utils.ValueIgnoreEmpty(d.Get("risk_level")),
	}

	return bodyParams
}

func resourceDocumentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	deleteHttpUrl := "v1/documents/{document_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{document_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "COC.00042101"),
			"error deleting document")
	}

	return nil
}
