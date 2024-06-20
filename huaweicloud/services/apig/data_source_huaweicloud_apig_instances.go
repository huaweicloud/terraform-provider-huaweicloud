package apig

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

// @API APIG GET /v2/{project_id}/apigw/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the instance to be queried.",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the name of the instance to be queried.",
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the status of the instance to be queried.",
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Specifies the ID of the enterprise project to which the instances belong.",
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The ID of instance.",
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The name of instance.",
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of instance.",
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The status of instance.",
						},
						"edition": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The edition of instance.",
						},
						"enterprise_project_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The enterprise project ID of the instance.",
						},
						"eip_address": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The elastic IP address of instance binding.",
						},
						"loadbalancer_provider": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The type of load balancer used by the instance.",
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The creation time of the instance, in RFC3339 format.",
						},
					},
				},
			},
		},
	}
}

func buildListInstancesParams(d *schema.ResourceData) string {
	res := ""
	if instanceId, ok := d.GetOk("instance_id"); ok {
		res = fmt.Sprintf("%s&instance_id=%v", res, instanceId)
	}
	if name, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&instance_name=%v", res, name)
	}
	if status, ok := d.GetOk("status"); ok {
		res = fmt.Sprintf("%s&status=%v", res, status)
	}
	return res
}

func queryInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/instances?limit=200"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	queryParams := buildListInstancesParams(d)
	listPath += queryParams

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, fmt.Errorf("error retrieving instances: %s", err)
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		instances := utils.PathSearch("instances", respBody, make([]interface{}, 0)).([]interface{})
		if len(instances) < 1 {
			break
		}
		result = append(result, instances...)
		offset += len(instances)
	}
	return result, nil
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}
	instances, err := queryInstances(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", filterListInstanceBody(flattenInstances(instances), d, cfg.GetEnterpriseProjectID(d))),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstances(instances []interface{}) []interface{} {
	if len(instances) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(instances))
	for _, authorizer := range instances {
		createTime := utils.PathSearch("create_time", authorizer, float64(0)).(float64)
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", authorizer, nil),
			"name":                  utils.PathSearch("instance_name", authorizer, nil),
			"type":                  utils.PathSearch("type", authorizer, nil),
			"status":                utils.PathSearch("status", authorizer, nil),
			"edition":               utils.PathSearch("spec", authorizer, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", authorizer, nil),
			"eip_address":           utils.PathSearch("eip_address", authorizer, nil),
			"loadbalancer_provider": utils.PathSearch("loadbalancer_provider", authorizer, nil),
			"created_at":            utils.FormatTimeStampRFC3339(int64(createTime)/1000, true),
		})
	}
	return result
}

func filterListInstanceBody(all []interface{}, _ *schema.ResourceData, enterpriseProjectId string) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if enterpriseProjectId != "" && enterpriseProjectId != utils.PathSearch("enterprise_project_id", v, "").(string) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}
