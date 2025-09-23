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

// @API LTS DELETE /v3/{project_id}/lts/access-config
// @API LTS POST /v2.0/{project_id}/lts/createAgencyAccess
// @API LTS POST /v1/{project_id}/{resource_type}/{resource_id}/tags/action
// @API LTS POST /v3/{project_id}/lts/access-config-list
func ResourceCrossAccountAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCrossAccountAccessCreate,
		UpdateContext: resourceCrossAccountAccessUpdate,
		ReadContext:   resourceCrossAccountAccessRead,
		DeleteContext: resourceHostAccessConfigDelete,

		CustomizeDiff: config.MergeDefaultTags(),

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
				Description: `Specifies the name of the cross account access.`,
			},
			"log_agencystream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log stream ID that already exists in the delegated account.`,
			},
			"log_agencystream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log stream name that already exists in the delegated account.`,
			},
			"log_agencygroup_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specify the log group ID that already exists in the delegated account.`,
			},
			"log_agencygroup_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specify the log group name that already exists in the delegated account.`,
			},
			"log_stream_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log stream ID that already exists in the delegatee account.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the log stream name that already exists in the delegatee account.`,
			},
			"log_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specify the log group ID that already exists in the delegatee account.`,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specify the log group name that already exists in the delegatee account.`,
			},
			"agency_project_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the delegator project ID.`,
			},
			"agency_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the agency created in IAM by the delegator.`,
			},
			"agency_domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the delegator account to verify the delegation.`,
			},
			"tags": common.TagsSchema(),
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the cross account access, in RFC3339 format.`,
			},
			"access_config_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The log access configuration type.`,
			},
		},
	}
}

func resourceCrossAccountAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v2.0/{project_id}/lts/createAgencyAccess"
		product       = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	createPath := ltsClient.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", ltsClient.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createOpt.JSONBody = utils.RemoveNil(buildCreateBodyParams(d, ltsClient.ProjectID))
	createResp, err := ltsClient.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating LTS cross account access: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	accessId := utils.PathSearch("[0].access_config_id", createRespBody, "").(string)
	if accessId == "" {
		return diag.Errorf("unable to find the access ID of the LTS cross account from the API response")
	}
	d.SetId(accessId)

	// deal tags
	if tags, ok := d.GetOk("tags"); ok {
		createTagsHttpUrl := "v1/{project_id}/{resource_type}/{resource_id}/tags/action"
		createTagsPath := ltsClient.Endpoint + createTagsHttpUrl
		createTagsPath = strings.ReplaceAll(createTagsPath, "{project_id}", ltsClient.ProjectID)
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_type}", "ltsAccessConfig")
		createTagsPath = strings.ReplaceAll(createTagsPath, "{resource_id}", d.Id())

		createTagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		createTagsOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tags.(map[string]interface{}), "create"))
		_, err := ltsClient.Request("POST", createTagsPath, &createTagsOpt)

		if err != nil {
			return diag.Errorf("error creating LTS cross account access tags: %s", err)
		}
	}

	return resourceCrossAccountAccessRead(ctx, d, meta)
}

func buildCreateBodyParams(d *schema.ResourceData, projectId string) map[string]interface{} {
	bodyParams := make([]map[string]interface{}, 1)
	bodyParams[0] = map[string]interface{}{
		"agency_access_type":      "AGENCYACCESS",
		"agency_log_access":       d.Get("name"),
		"log_agencyStream_id":     d.Get("log_agencystream_id"),
		"log_agencyStream_name":   d.Get("log_agencystream_name"),
		"log_agencyGroup_id":      d.Get("log_agencygroup_id"),
		"log_agencyGroup_name":    d.Get("log_agencygroup_name"),
		"log_beAgencystream_id":   d.Get("log_stream_id"),
		"log_beAgencystream_name": d.Get("log_stream_name"),
		"log_beAgencygroup_id":    d.Get("log_group_id"),
		"log_beAgencygroup_name":  d.Get("log_group_name"),
		"be_agency_project_id":    projectId,
		"agency_project_id":       d.Get("agency_project_id"),
		"agency_name":             d.Get("agency_name"),
		"agency_domain_name":      d.Get("agency_domain_name"),
	}
	bodyParam := map[string]interface{}{
		"preview_agency_list": bodyParams,
	}
	return bodyParam
}

func buildTagsBodyParams(tags map[string]interface{}, action string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"action":  action,
		"is_open": false,
		"tags":    utils.ExpandResourceTags(tags),
	}
	return bodyParams
}
func resourceCrossAccountAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listHttpUrl = "v3/{project_id}/lts/access-config-list"
		listProduct = "lts"
	)
	ltsClient, err := cfg.NewServiceClient(listProduct, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listPath := ltsClient.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", ltsClient.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	name := d.Get("name").(string)
	listOpt.JSONBody = map[string]interface{}{
		"access_config_name_list": []string{name},
	}
	listResp, err := ltsClient.Request("POST", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving LTS cross account access")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("result[?access_config_name=='%s']|[0]", name)
	listRespBody = utils.PathSearch(jsonPath, listRespBody, nil)
	if listRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "Cross account access")
	}

	created := utils.PathSearch("create_time", listRespBody, float64(0)).(float64)
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("access_config_name", listRespBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(created)/1000, false)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("access_config_tag", listRespBody, nil))),
		d.Set("access_config_type", utils.PathSearch("access_config_type", listRespBody, nil)),
		d.Set("log_group_id", utils.PathSearch("log_info.log_group_id", listRespBody, nil)),
		d.Set("log_group_name", utils.PathSearch("log_info.log_group_name", listRespBody, nil)),
		d.Set("log_stream_id", utils.PathSearch("log_info.log_stream_id", listRespBody, nil)),
		d.Set("log_stream_name", utils.PathSearch("log_info.log_stream_name", listRespBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCrossAccountAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	product := "lts"
	tagsHttpUrl := "v1/{project_id}/{resource_type}/{resource_id}/tags/action"

	ltsClient, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	if d.HasChanges("tags") {
		tagsPath := ltsClient.Endpoint + tagsHttpUrl
		tagsPath = strings.ReplaceAll(tagsPath, "{project_id}", ltsClient.ProjectID)
		tagsPath = strings.ReplaceAll(tagsPath, "{resource_type}", "ltsAccessConfig")
		tagsPath = strings.ReplaceAll(tagsPath, "{resource_id}", d.Id())

		tags := d.Get("tags").(map[string]interface{})
		tagsOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		// Creation, update and deletion options can all be implemented through action value: "create",
		// and the API logic is overwriting update (delete with empty array).
		tagsOpt.JSONBody = utils.RemoveNil(buildTagsBodyParams(tags, "create"))
		_, err = ltsClient.Request("POST", tagsPath, &tagsOpt)
		if err != nil {
			return diag.Errorf("error updating cross account access tags: %s", err)
		}
	}

	return resourceCrossAccountAccessRead(ctx, d, meta)
}
