package cce

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

// @API CCE GET /cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}/history
func DataSourceCCEReleaseHistory() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEReleaseHistoryRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the CCE cluster.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the name of the chart release.",
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the namespace to which the chart release belongs.",
			},
			"releases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"chart_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the chart.",
						},
						"chart_public": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: "Whether the chart is public.",
						},
						"chart_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The version of the chart.",
						},
						"cluster_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of the CCE cluster.",
						},
						"cluster_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the CCE cluster.",
						},
						"create_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the chart release.",
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The description of the chart release.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of the chart release.",
						},
						"namespace": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The namespace to which the chart release belongs.",
						},
						"parameters": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The parameters of the chart release in JSON format.",
						},
						"resources": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The resources required by the chart release in JSON format.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of the chart release.",
						},
						"status_description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status description of the chart release.",
						},
						"update_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The update time of the chart release.",
						},
						"values": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The values of the chart release in JSON format.",
						},
						"version": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: "The version number of the chart release.",
						},
					},
				},
			},
		},
	}
}

func dataSourceCCEReleaseHistoryRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "cce"
		httpUrl   = "cce/cam/v3/clusters/{cluster_id}/namespace/{namespace}/releases/{name}/history"
		clusterID = d.Get("cluster_id").(string)
		namespace = d.Get("namespace").(string)
		name      = d.Get("name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{cluster_id}", clusterID)
	requestPath = strings.ReplaceAll(requestPath, "{namespace}", namespace)
	requestPath = strings.ReplaceAll(requestPath, "{name}", name)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving CCE chart release history: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("releases", flattenCCEReleaseHistory(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCCEReleaseHistory(resp interface{}) []interface{} {
	respList, ok := resp.([]interface{})
	if !ok {
		return nil
	}
	rst := make([]interface{}, 0, len(respList))
	for _, v := range respList {
		rst = append(rst, map[string]interface{}{
			"chart_name":         utils.PathSearch("chart_name", v, nil),
			"chart_public":       utils.PathSearch("chart_public", v, nil),
			"chart_version":      utils.PathSearch("chart_version", v, nil),
			"cluster_id":         utils.PathSearch("cluster_id", v, nil),
			"cluster_name":       utils.PathSearch("cluster_name", v, nil),
			"create_at":          utils.PathSearch("create_at", v, nil),
			"description":        utils.PathSearch("description", v, nil),
			"name":               utils.PathSearch("name", v, nil),
			"namespace":          utils.PathSearch("namespace", v, nil),
			"parameters":         utils.PathSearch("parameters", v, nil),
			"resources":          utils.PathSearch("resources", v, nil),
			"status":             utils.PathSearch("status", v, nil),
			"status_description": utils.PathSearch("status_description", v, nil),
			"update_at":          utils.PathSearch("update_at", v, nil),
			"values":             utils.PathSearch("values", v, nil),
			"version":            utils.PathSearch("version", v, nil),
		})
	}

	return rst
}
