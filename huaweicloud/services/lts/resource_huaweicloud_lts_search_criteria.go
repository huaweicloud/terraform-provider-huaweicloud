package lts

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

// @API LTS POST /v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias
// @API LTS GET /v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias
// @API LTS DELETE /v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias
func ResourceSearchCriteria() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSearchCriteriaCreate,
		ReadContext:   resourceSearchCriteriaRead,
		DeleteContext: resourceSearchCriteriaDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSearchCriteriaImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"log_group_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"log_stream_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"criteria": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
				Description: utils.SchemaDesc("The enterprise project ID.",
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func resourceSearchCriteriaCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSearchCriteria: Create an LTS search criteria.
	var (
		getSearchCriteriaHttpUrl = "v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias"
		getSearchCriteriaProduct = "lts"
	)
	createSearchCriteriaClient, err := cfg.NewServiceClient(getSearchCriteriaProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	createSearchCriteriaPath := createSearchCriteriaClient.Endpoint + getSearchCriteriaHttpUrl
	createSearchCriteriaPath = strings.ReplaceAll(createSearchCriteriaPath, "{project_id}", createSearchCriteriaClient.ProjectID)
	createSearchCriteriaPath = strings.ReplaceAll(createSearchCriteriaPath, "{group_id}", d.Get("log_group_id").(string))
	createSearchCriteriaPath = strings.ReplaceAll(createSearchCriteriaPath, "{topic_id}", d.Get("log_stream_id").(string))

	createSearchCriteriaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createSearchCriteriaOpt.JSONBody = utils.RemoveNil(buildCreateSearchCriteriaBodyParams(d, cfg))
	createSearchCriteriaResp, err := createSearchCriteriaClient.Request("POST", createSearchCriteriaPath, &createSearchCriteriaOpt)
	if err != nil {
		return diag.Errorf("error creating search criteria: %s", err)
	}

	createSearchCriteriaRespBody, err := utils.FlattenResponse(createSearchCriteriaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	criteriaId := utils.PathSearch("id", createSearchCriteriaRespBody, "").(string)
	if criteriaId == "" {
		return diag.Errorf("unable to find the LTS search criteria ID from the API response")
	}
	d.SetId(criteriaId)

	return resourceSearchCriteriaRead(ctx, d, meta)
}

func buildCreateSearchCriteriaBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"criteria":              d.Get("criteria"),
		"name":                  d.Get("name"),
		"search_type":           d.Get("type"),
		"enterprise_project_id": cfg.GetEnterpriseProjectID(d),
	}
	return bodyParams
}

func resourceSearchCriteriaRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSearchCriteria: Query the LTS search criteria detail.
	var (
		getSearchCriteriaHttpUrl = "v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias"
		getSearchCriteriaProduct = "lts"
	)
	getSearchCriteriaClient, err := cfg.NewServiceClient(getSearchCriteriaProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	getSearchCriteriaPath := getSearchCriteriaClient.Endpoint + getSearchCriteriaHttpUrl
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{project_id}", getSearchCriteriaClient.ProjectID)
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{group_id}", d.Get("log_group_id").(string))
	getSearchCriteriaPath = strings.ReplaceAll(getSearchCriteriaPath, "{topic_id}", d.Get("log_stream_id").(string))
	// If `type` parameter is not specified, the default query is "ORIGINALLOG" type, and the "VISUALIZATION" type cannot be queried.
	getSearchCriteriaPath = fmt.Sprintf("%s?search_type=%s", getSearchCriteriaPath, d.Get("type").(string))

	getSearchCriteriaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getSearchCriteriaOpt.JSONBody = d.Id()
	getSearchCriteriaResp, err := getSearchCriteriaClient.Request("GET", getSearchCriteriaPath, &getSearchCriteriaOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving search criteria")
	}

	getSearchCriteriaRespBody, err := utils.FlattenResponse(getSearchCriteriaResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("search_criterias[?id =='%s']|[0]", d.Id())
	getSearchCriteriaRespBody = utils.PathSearch(jsonPath, getSearchCriteriaRespBody, nil)
	if getSearchCriteriaRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("criteria", utils.PathSearch("criteria", getSearchCriteriaRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getSearchCriteriaRespBody, nil)),
		d.Set("type", utils.PathSearch("search_type", getSearchCriteriaRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getSearchCriteriaRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSearchCriteriaDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSearchCriteria: Delete an existing LTS search criteria.
	var (
		deleteSearchCriteriaHttpUrl = "v1.0/{project_id}/groups/{group_id}/topics/{topic_id}/search-criterias"
		deleteSearchCriteriaProduct = "lts"
	)
	deleteSearchCriteriaClient, err := cfg.NewServiceClient(deleteSearchCriteriaProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS Client: %s", err)
	}

	deleteSearchCriteriaPath := deleteSearchCriteriaClient.Endpoint + deleteSearchCriteriaHttpUrl
	deleteSearchCriteriaPath = strings.ReplaceAll(deleteSearchCriteriaPath, "{project_id}", deleteSearchCriteriaClient.ProjectID)
	deleteSearchCriteriaPath = strings.ReplaceAll(deleteSearchCriteriaPath, "{group_id}", d.Get("log_group_id").(string))
	deleteSearchCriteriaPath = strings.ReplaceAll(deleteSearchCriteriaPath, "{topic_id}", d.Get("log_stream_id").(string))

	deleteSearchCriteriaOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteSearchCriteriaOpt.JSONBody = utils.RemoveNil(BuildGetOrDeleteSearchCriteriaBodyParams(d.Id()))
	_, err = deleteSearchCriteriaClient.Request("DELETE", deleteSearchCriteriaPath, &deleteSearchCriteriaOpt)
	if err != nil {
		return diag.Errorf("error deleting search criteria: %s", err)
	}

	return nil
}

func BuildGetOrDeleteSearchCriteriaBodyParams(id string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id": id,
	}
	return bodyParams
}

func resourceSearchCriteriaImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid ID format, want '<log_group_id>/<log_stream_id>/<id>', but got '%s'", d.Id())
	}

	groupID := parts[0]
	streamID := parts[1]
	searchCriteriaID := parts[2]

	d.SetId(searchCriteriaID)
	mErr := multierror.Append(nil,
		d.Set("log_group_id", groupID),
		d.Set("log_stream_id", streamID),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
