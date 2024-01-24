package vpcep

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API VPCEP GET /v1/{project_id}/vpc-endpoints
func DataSourceVPCEPEndpoints() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceVpcepEndpointsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"endpoints": {
				Type:     schema.TypeList,
				Elem:     endpointSchema(),
				Computed: true,
			},
		},
	}
}

func endpointSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"packet_id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"enable_dns": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"subnet_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsComputedSchema(),
			"whitelist": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_whitelist": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func dataSourceVpcepEndpointsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	// get VPC endpoints: Query VPC endpoints
	var (
		getVPCEPEndpointsHttpUrl = "v1/{project_id}/vpc-endpoints"
		getVPCEPEndpointsProduct = "vpcep"
	)
	getVPCEPEndpointsClient, err := conf.NewServiceClient(getVPCEPEndpointsProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPCEP client: %s", err)
	}

	getVPCEPEndpointsPath := getVPCEPEndpointsClient.Endpoint + getVPCEPEndpointsHttpUrl
	getVPCEPEndpointsPath = strings.ReplaceAll(getVPCEPEndpointsPath, "{project_id}",
		getVPCEPEndpointsClient.ProjectID)
	getVPCEPEndpointsPath += buildVPCEPEndpointsQueryParams(d, conf)

	getVPCEPEndpointsResp, err := pagination.ListAllItems(
		getVPCEPEndpointsClient,
		"offset",
		getVPCEPEndpointsPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return diag.Errorf("error retrieving endpoints, %s", err)
	}

	listVPCEPEndpointsRespJson, err := json.Marshal(getVPCEPEndpointsResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listVPCEPEndpointsRespBody interface{}
	err = json.Unmarshal(listVPCEPEndpointsRespJson, &listVPCEPEndpointsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}
	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("endpoints", flattenListEndpointsBody(listVPCEPEndpointsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildVPCEPEndpointsQueryParams(d *schema.ResourceData, _ *config.Config) string {
	res := ""
	if v, ok := d.GetOk("service_name"); ok {
		res = fmt.Sprintf("%s&endpoint_service_name=%v", res, v)
	}
	if v, ok := d.GetOk("vpc_id"); ok {
		res = fmt.Sprintf("%s&vpc_id=%v", res, v)
	}
	if v, ok := d.GetOk("endpoint_id"); ok {
		res = fmt.Sprintf("%s&id=%v", res, v)
	}
	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenListEndpointsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("endpoints", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":               utils.PathSearch("id", v, nil),
			"service_type":     utils.PathSearch("service_type", v, nil),
			"status":           utils.PathSearch("status", v, nil),
			"service_name":     utils.PathSearch("endpoint_service_name", v, nil),
			"service_id":       utils.PathSearch("endpoint_service_id", v, nil),
			"packet_id":        utils.PathSearch("marker_id", v, nil),
			"ip_address":       utils.PathSearch("ip", v, nil),
			"vpc_id":           utils.PathSearch("vpc_id", v, nil),
			"subnet_id":        utils.PathSearch("subnet_id", v, nil),
			"created_at":       utils.PathSearch("created_at", v, nil),
			"updated_at":       utils.PathSearch("updated_at", v, nil),
			"description":      utils.PathSearch("description", v, nil),
			"enable_dns":       utils.StringToBool(utils.PathSearch("enable_dns", v, "")),
			"enable_whitelist": utils.StringToBool(utils.PathSearch("enable_whitelist", v, "")),
			"tags":             utils.FlattenTagsToMap(utils.PathSearch("tags", v, make([]interface{}, 0))),
			"whitelist":        utils.ExpandToStringList(utils.PathSearch("whitelist", v, make([]interface{}, 0)).([]interface{})),
		})
	}
	return rst
}
