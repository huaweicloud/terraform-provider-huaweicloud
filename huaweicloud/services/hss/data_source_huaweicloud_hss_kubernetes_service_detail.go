package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS GET /v5/{project_id}/kubernetes/service/detail
func DataSourceKubernetesServiceDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceKubernetesServiceDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			// The field in the API document is `id`, here it is modified to `service_id`.
			// This field is required in the API documentation and optional in actual testing.
			// This is consistent with the API documentation.
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"endpoint_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creation_timestamp": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"labels": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"selector": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"session_affinity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_port_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"service_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protocol": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"target_port": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"node_port": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildKubernetesServiceDetailQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := fmt.Sprintf("?id=%v", d.Get("service_id"))

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}

	return queryParams
}

func dataSourceKubernetesServiceDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "hss"
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v5/{project_id}/kubernetes/service/detail"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildKubernetesServiceDetailQueryParams(d, epsId)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving HSS kubernetes service detail: %s", err)
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
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("endpoint_name", utils.PathSearch("endpoint_name", respBody, nil)),
		d.Set("namespace", utils.PathSearch("namespace", respBody, nil)),
		d.Set("creation_timestamp", utils.PathSearch("creation_timestamp", respBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", respBody, nil)),
		d.Set("labels", utils.PathSearch("labels", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("cluster_ip", utils.PathSearch("cluster_ip", respBody, nil)),
		d.Set("selector", utils.PathSearch("selector", respBody, nil)),
		d.Set("session_affinity", utils.PathSearch("session_affinity", respBody, nil)),
		d.Set("service_port_list",
			flattenServicePortList(utils.PathSearch("service_port_list", resp, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenServicePortList(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("id", v, nil),
			"service_id":  utils.PathSearch("service_id", v, nil),
			"name":        utils.PathSearch("name", v, nil),
			"protocol":    utils.PathSearch("protocol", v, nil),
			"port":        utils.PathSearch("port", v, nil),
			"target_port": utils.PathSearch("target_port", v, nil),
			"node_port":   utils.PathSearch("node_port", v, nil),
		})
	}
	return rst
}
