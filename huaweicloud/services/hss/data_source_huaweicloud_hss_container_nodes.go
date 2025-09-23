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

// @API HSS GET /v5/{project_id}/container/nodes
func DataSourceContainerNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceContainerNodesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"agent_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"container_tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"agent_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"host_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agent_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"protect_interrupt": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"protect_degradation": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"degradation_reason": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container_tags": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"private_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"public_ip": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enterprise_project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"detect_result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"asset": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"vulnerability": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"intrusion": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"policy_group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"policy_group_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildContainerNodesQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := "?limit=200"
	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("host_name"); ok {
		queryParams = fmt.Sprintf("%s&host_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("agent_status"); ok {
		queryParams = fmt.Sprintf("%s&agent_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("protect_status"); ok {
		queryParams = fmt.Sprintf("%s&protect_status=%v", queryParams, v)
	}
	if v, ok := d.GetOk("container_tags"); ok {
		queryParams = fmt.Sprintf("%s&container_tags=%v", queryParams, v)
	}

	return queryParams
}

func dataSourceContainerNodesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		product = "hss"
		httpUrl = "v5/{project_id}/container/nodes"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildContainerNodesQueryParams(d, epsId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving HSS container nodes: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data_list", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		result = append(result, dataResp...)
		offset += len(dataResp)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data_list", flattenContainerDataList(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenContainerDataList(dataResp []interface{}) []interface{} {
	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"agent_id":                utils.PathSearch("agent_id", v, nil),
			"host_id":                 utils.PathSearch("host_id", v, nil),
			"host_name":               utils.PathSearch("host_name", v, nil),
			"host_status":             utils.PathSearch("host_status", v, nil),
			"agent_status":            utils.PathSearch("agent_status", v, nil),
			"protect_status":          utils.PathSearch("protect_status", v, nil),
			"protect_interrupt":       utils.PathSearch("protect_interrupt", v, nil),
			"protect_degradation":     utils.PathSearch("protect_degradation", v, nil),
			"degradation_reason":      utils.PathSearch("degradation_reason", v, nil),
			"container_tags":          utils.PathSearch("container_tags", v, nil),
			"private_ip":              utils.PathSearch("private_ip", v, nil),
			"public_ip":               utils.PathSearch("public_ip", v, nil),
			"resource_id":             utils.PathSearch("resource_id", v, nil),
			"group_name":              utils.PathSearch("group_name", v, nil),
			"enterprise_project_name": utils.PathSearch("enterprise_project_name", v, nil),
			"detect_result":           utils.PathSearch("detect_result", v, nil),
			"asset":                   utils.PathSearch("asset", v, nil),
			"vulnerability":           utils.PathSearch("vulnerability", v, nil),
			"intrusion":               utils.PathSearch("intrusion", v, nil),
			"policy_group_id":         utils.PathSearch("policy_group_id", v, nil),
			"policy_group_name":       utils.PathSearch("policy_group_name", v, nil),
		})
	}

	return result
}
