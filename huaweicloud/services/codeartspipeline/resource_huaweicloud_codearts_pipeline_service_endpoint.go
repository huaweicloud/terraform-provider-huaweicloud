package codeartspipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var serviceEnopointNonUpdatableParams = []string{
	"project_id", "module_id", "url", "name", "authorization", "authorization.0.parameters", "authorization.0.scheme",
}

// @API CodeArtsPipeline POST /v1/serviceconnection/endpoints
// @API CodeArtsPipeline GET /v1/serviceconnection/endpoints
func ResourceCodeArtsPipelineServiceEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePipelineServiceEndpointCreate,
		ReadContext:   resourcePipelineServiceEndpointRead,
		UpdateContext: resourcePipelineServiceEndpointUpdate,
		DeleteContext: resourcePipelineServiceEndpointDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceImportStateFuncWithProjectIdAndId,
		},

		CustomizeDiff: config.FlexibleForceNew(serviceEnopointNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"module_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the module ID.`,
			},
			"url": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the URL.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the endpoint name.`,
			},
			"authorization": {
				Type:        schema.TypeList,
				Optional:    true,
				MaxItems:    1,
				Description: `Specifies the permission information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"scheme": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the authentication mode.`,
						},
						"parameters": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `Specifies the authentication parameter.`,
						},
					},
				},
			},
			"created_by": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the permission information.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user name.`,
						},
					},
				},
			},
		},
	}
}

func resourcePipelineServiceEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/serviceconnection/endpoints"
	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreatePipelineServiceEndpointBodyParams(d, region)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline service endpoint: %s", err)
	}
	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := checkResponseError(createRespBody, ""); err != nil {
		return diag.Errorf("error creating CodeArts Pipeline service endpoint: %s", err)
	}

	id := utils.PathSearch("result.uuid", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the CodeArts Pipeline service endpoint ID from the API response")
	}

	d.SetId(id)

	return resourcePipelineServiceEndpointRead(ctx, d, meta)
}

func buildCreatePipelineServiceEndpointBodyParams(d *schema.ResourceData, region string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"region_name":   region,
		"project_uuid":  d.Get("project_id"),
		"module_id":     utils.ValueIgnoreEmpty(d.Get("module_id")),
		"url":           utils.ValueIgnoreEmpty(d.Get("url")),
		"name":          utils.ValueIgnoreEmpty(d.Get("name")),
		"authorization": buildPipelineServiceEndpointAuthorization(d),
	}

	return bodyParams
}

func buildPipelineServiceEndpointAuthorization(d *schema.ResourceData) interface{} {
	rawParams := d.Get("authorization").([]interface{})
	if len(rawParams) == 0 {
		return nil
	}

	if param, ok := rawParams[0].(map[string]interface{}); ok {
		rst := map[string]interface{}{
			// `parameters` is an object in request
			"parameters": parseJson(param["parameters"].(string)),
			"scheme":     utils.ValueIgnoreEmpty(param["scheme"]),
		}
		return rst
	}

	return nil
}

func resourcePipelineServiceEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getRespBody, err := GetPipelineServiceEndpoint(client, region, d.Get("project_id").(string), d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CodeArts Pipeline service endpoint")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("project_id", utils.PathSearch("project_uuid", getRespBody, nil)),
		d.Set("module_id", utils.PathSearch("module_id", getRespBody, nil)),
		d.Set("url", utils.PathSearch("url", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("created_by", flattenPipelineServiceEndpointCreatedBy(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetPipelineServiceEndpoint(client *golangsdk.ServiceClient, region, projectId, id string) (interface{}, error) {
	httpUrl := "v1/serviceconnection/endpoints"
	listPath := client.Endpoint + httpUrl
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	pageSize := 100
	listPath += fmt.Sprintf("?region_name=%s&project_uuid=%s&page_size=%d", region, projectId, pageSize)
	pageIndex := 1
	for {
		currentPath := listPath + fmt.Sprintf("&page_index=%d", pageIndex)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}
		if err := checkResponseError(listRespBody, projectNotFoundError); err != nil {
			return nil, err
		}

		endpoints := utils.PathSearch("result.endpoints", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(endpoints) == 0 {
			return nil, golangsdk.ErrDefault404{}
		}

		searchPath := fmt.Sprintf("result.endpoints[?uuid=='%s']|[0]", id)
		result := utils.PathSearch(searchPath, listRespBody, nil)
		if result != nil {
			return result, nil
		}

		pageIndex++
	}
}

func flattenPipelineServiceEndpointCreatedBy(resp interface{}) []interface{} {
	createdBy := utils.PathSearch("created_by", resp, nil)
	if createdBy == nil {
		return nil
	}

	rst := map[string]interface{}{
		"user_id":   utils.PathSearch("user_id", createdBy, nil),
		"user_name": utils.PathSearch("username", createdBy, nil),
	}

	return []interface{}{rst}
}

func resourcePipelineServiceEndpointUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourcePipelineServiceEndpointDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting pipeline service endpoint resource is not supported. The resource is only removed from the" +
		"state, the service endpoint remains in the cloud."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
