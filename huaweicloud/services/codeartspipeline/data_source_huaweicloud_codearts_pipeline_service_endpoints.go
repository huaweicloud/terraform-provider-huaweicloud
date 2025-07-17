package codeartspipeline

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CodeArtsPipeline GET /v1/serviceconnection/endpoints
func DataSourceCodeArtsPipelineServiceEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineServiceEndpointsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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
			"endpoints": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the endpoint list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the endpoint ID.`,
						},
						"module_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the module ID.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the URL.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the endpoint name.`,
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
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineServiceEndpointsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	httpUrl := "v1/serviceconnection/endpoints"
	listPath := client.Endpoint + httpUrl
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	pageSize := 100
	listPath += fmt.Sprintf("?region_name=%s&project_uuid=%s&page_size=%d", region, d.Get("project_id").(string), pageSize)
	if v, ok := d.GetOk("module_id"); ok {
		listPath += fmt.Sprintf("&module_id=%v", v)
	}
	pageIndex := 1
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := listPath + fmt.Sprintf("&page_index=%d", pageIndex)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return diag.Errorf("error getting service endpoints: %s", err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.Errorf("error flattening response: %s", err)
		}
		if err := checkResponseError(listRespBody, projectNotFoundError); err != nil {
			return diag.Errorf("error getting service endpoints: %s", err)
		}

		endpoints := utils.PathSearch("result.endpoints", listRespBody, make([]interface{}, 0)).([]interface{})
		if len(endpoints) == 0 {
			break
		}

		for _, endpoint := range endpoints {
			rst = append(rst, map[string]interface{}{
				"id":         utils.PathSearch("uuid", endpoint, nil),
				"module_id":  utils.PathSearch("module_id", endpoint, nil),
				"url":        utils.PathSearch("url", endpoint, nil),
				"name":       utils.PathSearch("name", endpoint, nil),
				"created_by": flattenPipelineServiceEndpointCreatedBy(endpoint),
			})
		}

		pageIndex++
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("endpoints", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
