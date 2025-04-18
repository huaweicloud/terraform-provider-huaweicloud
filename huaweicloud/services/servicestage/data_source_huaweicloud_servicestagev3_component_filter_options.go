package servicestage

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ServiceStage GET /v3/{project_id}/cas/components/filterOptions
func DataSourceV3ComponentUsedResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3ComponentUsedResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the components are located.`,
			},
			"applications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the application that component used.`,
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the application that component used.`,
						},
					},
				},
				Description: `The list of applications that component used.`,
			},
			"enterprise_projects": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID list of the enterprise projects that component used.`,
			},
			"environments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the environment that component used.`,
						},
						"label": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the environment that component used.`,
						},
					},
				},
				Description: `The list of environments that component used.`,
			},
		},
	}
}

func listV3ComponentUsedResources(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v3/{project_id}/cas/components/filterOptions"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenV3ComponentUsedResourcesInfo(objInfos []interface{}) []map[string]interface{} {
	if len(objInfos) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(objInfos))
	for _, obj := range objInfos {
		result = append(result, map[string]interface{}{
			"id":    utils.PathSearch("id", obj, nil),
			"label": utils.PathSearch("label", obj, nil),
		})
	}

	return result
}

func dataSourceV3ComponentUsedResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	filterOptions, err := listV3ComponentUsedResources(client)
	if err != nil {
		return diag.Errorf("error getting component used resources: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("applications", flattenV3ComponentUsedResourcesInfo(utils.PathSearch("filter_options.applications",
			filterOptions, make([]interface{}, 0)).([]interface{}))),
		d.Set("enterprise_projects", utils.PathSearch("filter_options.enterprise_projects",
			filterOptions, make([]interface{}, 0))),
		d.Set("environments", flattenV3ComponentUsedResourcesInfo(utils.PathSearch("filter_options.environments",
			filterOptions, make([]interface{}, 0)).([]interface{}))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
