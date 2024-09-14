package accessanalyzer

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AccessAnalyzer POST /v5/analyzers
// @API AccessAnalyzer GET /v5/analyzers/{analyzer_id}
// @API AccessAnalyzer DELETE /v5/analyzers/{analyzer_id}
// @API AccessAnalyzer POST /v5/{resource_type}/{resource_id}/tags/create
// @API AccessAnalyzer POST /v5/{resource_type}/{resource_id}/tags/delete
var nonUpdatableParams = []string{"name", "type"}

func ResourceAccessAnalyzer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessAnalyzerCreate,
		ReadContext:   resourceAccessAnalyzerRead,
		UpdateContext: resourceAccessAnalyzerUpdate,
		DeleteContext: resourceAccessAnalyzerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(nonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_analyzed_resource": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_resource_analyzed_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceAccessAnalyzerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	createAnalyzerHttpUrl := "v5/analyzers"
	createAnalyzerPath := client.Endpoint + createAnalyzerHttpUrl
	createAnalyzerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAnalyzerBodyParams(d)),
	}
	createAnalyzerResp, err := client.Request("POST", createAnalyzerPath, &createAnalyzerOpt)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer analyzer: %s", err)
	}
	createAnalyzerRespBody, err := utils.FlattenResponse(createAnalyzerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createAnalyzerRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating access analyzer: id is not found in API response")
	}
	d.SetId(id)

	return resourceAccessAnalyzerRead(ctx, d, meta)
}

func buildCreateAnalyzerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name": d.Get("name"),
		"type": d.Get("type").(string),
		"tags": utils.ExpandResourceTagsMap(d.Get("tags").(map[string]interface{})),
	}
	return bodyParams
}

func resourceAccessAnalyzerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	getAnalyzerHttpUrl := "v5/analyzers/{analyzer_id}"
	getAnalyzerPath := client.Endpoint + getAnalyzerHttpUrl
	getAnalyzerPath = strings.ReplaceAll(getAnalyzerPath, "{analyzer_id}", d.Id())
	getAnalyzerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAnalyzerResp, err := client.Request("GET", getAnalyzerPath, &getAnalyzerOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving access analyzer")
	}
	getAnalyzerRespBody, err := utils.FlattenResponse(getAnalyzerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	analyzer := utils.PathSearch("analyzer", getAnalyzerRespBody, nil)
	if analyzer == nil {
		return diag.Errorf("error getting access analyzer: analyzer is not found in API response")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("name", analyzer, nil)),
		d.Set("type", utils.PathSearch("type", analyzer, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", analyzer, make([]interface{}, 0)))),
		d.Set("status", utils.PathSearch("status", analyzer, nil)),
		d.Set("urn", utils.PathSearch("urn", analyzer, nil)),
		d.Set("last_analyzed_resource", utils.PathSearch("last_analyzed_resource", analyzer, nil)),
		d.Set("last_resource_analyzed_at", utils.PathSearch("last_resource_analyzed_at", analyzer, nil)),
		d.Set("created_at", utils.PathSearch("created_at", analyzer, nil)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting access analyzer fields: %s", err)
	}

	return nil
}

func resourceAccessAnalyzerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	oRaw, nRaw := d.GetChange("tags")
	oMap := oRaw.(map[string]interface{})
	nMap := nRaw.(map[string]interface{})

	// remove old tags
	if len(oMap) > 0 {
		if err = deleteTags(client, oMap, "analyzers", d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	// set new tags
	if len(nMap) > 0 {
		if err := createTags(client, nMap, "analyzers", d.Id()); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceAccessAnalyzerRead(ctx, d, meta)
}

func resourceAccessAnalyzerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("accessanalyzer", region)
	if err != nil {
		return diag.Errorf("error creating Access Analyzer client: %s", err)
	}

	deleteAnalyzerHttpUrl := "v5/analyzers/{analyzer_id}"
	deleteAnalyzerPath := client.Endpoint + deleteAnalyzerHttpUrl
	deleteAnalyzerPath = strings.ReplaceAll(deleteAnalyzerPath, "{analyzer_id}", d.Id())
	deleteAnalyzerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteAnalyzerPath, &deleteAnalyzerOpt)
	if err != nil {
		return diag.Errorf("error deleting access analyzer: %s", err)
	}

	return nil
}

func createTags(createTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, id string) error {
	if len(tags) > 0 {
		createTagsHttpUrl := "v5/{resource_type}/{resource_id}/tags/create"
		createTagsPath := createTagsClient.Endpoint + createTagsHttpUrl
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", resourceType)
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", id)
		createTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"tags": utils.ExpandResourceTagsMap(tags),
			},
		}

		_, err := createTagsClient.Request("POST", createTagsPath, &createTagsOpt)
		if err != nil {
			return fmt.Errorf("error creating tags: %s", err)
		}
	}
	return nil
}

func deleteTags(deleteTagsClient *golangsdk.ServiceClient, tags map[string]interface{}, resourceType, id string) error {
	if len(tags) > 0 {
		deleteTagsHttpUrl := "v5/{resource_type}/{resource_id}/tags/delete"
		deleteTagsPath := deleteTagsClient.Endpoint + deleteTagsHttpUrl
		deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_type}", resourceType)
		deleteTagsPath = strings.ReplaceAll(deleteTagsPath, "{resource_id}", id)
		deleteTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody: map[string]interface{}{
				"tag_keys": expandTagsKeyToStringList(tags),
			},
		}

		_, err := deleteTagsClient.Request("POST", deleteTagsPath, &deleteTagsOpt)
		if err != nil {
			return fmt.Errorf("error deleting tags: %s", err)
		}
	}
	return nil
}

func expandTagsKeyToStringList(tagmap map[string]interface{}) []string {
	var tagKeyList []string
	for k := range tagmap {
		tagKeyList = append(tagKeyList, k)
	}
	return tagKeyList
}
