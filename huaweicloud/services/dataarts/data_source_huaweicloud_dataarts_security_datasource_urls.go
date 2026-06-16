package dataarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/security/permission-sets/datasource/urls
func DataSourceSecurityDatasourceUrls() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSecurityDatasourceUrlsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the datasource urls are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the datasource urls belong.`,
			},

			// Optional parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the cluster.`,
			},

			"datasource_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The type of the datasource.`,
			},

			"parent_permission_set_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the parent permission set.`,
			},

			// Attributes.
			"urls": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of datasource urls that matched filter parameters.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the url path.`,
						},
						"contains": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the parent permission set contains this permission.`,
						},
					},
				},
			},
		},
	}
}

func buildSecurityDatasourceUrlsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("cluster_id"); ok {
		res = fmt.Sprintf("%s&cluster_id=%v", res, v)
	}
	if v, ok := d.GetOk("datasource_type"); ok {
		res = fmt.Sprintf("%s&datasource_type=%v", res, v)
	}
	if v, ok := d.GetOk("parent_permission_set_id"); ok {
		res = fmt.Sprintf("%s&parent_permission_set_id=%v", res, v)
	}

	if len(res) < 1 {
		return res
	}
	return "?" + res[1:]
}

func listSecurityDatasourceUrls(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/security/permission-sets/datasource/urls"
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildSecurityDatasourceUrlsQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(d.Get("workspace_id").(string)),
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func flattenSecurityDatasourceUrls(urls []interface{}) []map[string]interface{} {
	if len(urls) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(urls))
	for _, url := range urls {
		result = append(result, map[string]interface{}{
			"name":     utils.PathSearch("name", url, nil),
			"contains": utils.PathSearch("contains", url, nil),
		})
	}

	return result
}

func dataSourceSecurityDatasourceUrlsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := listSecurityDatasourceUrls(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Security datasource urls: %s", err)
	}

	randomUUID, err := uuid.NewUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("urls", flattenSecurityDatasourceUrls(
			utils.PathSearch("urls", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
