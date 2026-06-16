package secmaster

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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/shippers
func DataSourceSiemShippers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemShippersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shipper_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shipper_source_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shipper_source_strategy": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shipper_consumption_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"destination_shipper_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"shipper_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     siemShippersDataSchema(),
			},
		},
	}
}

func siemShippersDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"consumption_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dataspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pipe_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shipper_destination": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     siemShippersDestinationSchema(),
			},
			"shipper_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shipper_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shipper_source": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     siemShippersSourceSchema(),
			},
			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func siemShippersDestinationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_param": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dataspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"destination_info": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func siemShippersSourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"data_type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"dataspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dataspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"identity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pipe_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSiemShippersQueryParams(d *schema.ResourceData, offset int) string {
	// The `limit` and `offset` are required parameters.
	rst := fmt.Sprintf("?limit=1000&offset=%d", offset)

	if v, ok := d.GetOk("dataspace_id"); ok {
		rst += fmt.Sprintf("&dataspace_id=%v", v)
	}

	if v, ok := d.GetOk("pipe_id"); ok {
		rst += fmt.Sprintf("&pipe_id=%v", v)
	}

	if v, ok := d.GetOk("shipper_name"); ok {
		rst += fmt.Sprintf("&shipper_name=%v", v)
	}

	if v, ok := d.GetOk("shipper_source_region"); ok {
		rst += fmt.Sprintf("&shipper_source_region=%v", v)
	}

	if v, ok := d.GetOk("shipper_source_strategy"); ok {
		rst += fmt.Sprintf("&shipper_source_strategy=%v", v)
	}

	if v, ok := d.GetOk("shipper_consumption_type"); ok {
		rst += fmt.Sprintf("&shipper_consumption_type=%v", v)
	}

	if v, ok := d.GetOk("destination_shipper_type"); ok {
		rst += fmt.Sprintf("&destination_shipper_type=%v", v)
	}

	if v, ok := d.GetOk("shipper_status"); ok {
		rst += fmt.Sprintf("&shipper_status=%v", v)
	}

	if v, ok := d.GetOk("create_time"); ok {
		rst += fmt.Sprintf("&create_time=%v", v)
	}

	return rst
}

func dataSourceSiemShippersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/shippers"
		product = "secmaster"
		offset  = 0
		allData = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := requestPath + buildSiemShippersQueryParams(d, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SIEM shippers: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		dataResp := utils.PathSearch("data.data", respBody, make([]interface{}, 0)).([]interface{})
		if len(dataResp) == 0 {
			break
		}

		allData = append(allData, dataResp...)
		offset += len(dataResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSiemShippers(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSiemShippers(shippers []interface{}) []interface{} {
	if len(shippers) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(shippers))
	for _, v := range shippers {
		rst = append(rst, map[string]interface{}{
			"consumption_type": utils.PathSearch("consumption_type", v, nil),
			"create_time":      utils.PathSearch("create_time", v, nil),
			"dataspace_id":     utils.PathSearch("dataspace_id", v, nil),
			"dataspace_name":   utils.PathSearch("dataspace_name", v, nil),
			"domain_id":        utils.PathSearch("domain_id", v, nil),
			"id":               utils.PathSearch("id", v, nil),
			"pipe_id":          utils.PathSearch("pipe_id", v, nil),
			"pipe_name":        utils.PathSearch("pipe_name", v, nil),
			"project_id":       utils.PathSearch("project_id", v, nil),
			"shipper_destination": flattenSiemShipperDestination(
				utils.PathSearch("shipper_destination", v, nil)),
			"shipper_id":   utils.PathSearch("shipper_id", v, nil),
			"shipper_name": utils.PathSearch("shipper_name", v, nil),
			"shipper_source": flattenSiemShipperSource(
				utils.PathSearch("shipper_source", v, nil)),
			"status":       utils.PathSearch("status", v, nil),
			"update_time":  utils.PathSearch("update_time", v, nil),
			"version":      utils.PathSearch("version", v, nil),
			"workspace_id": utils.PathSearch("workspace_id", v, nil),
		})
	}

	return rst
}

func flattenSiemShipperDestination(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"data_param":       utils.PathSearch("data_param", raw, nil),
			"data_type":        utils.PathSearch("data_type", raw, nil),
			"dataspace":        utils.PathSearch("dataspace", raw, nil),
			"dataspace_name":   utils.PathSearch("dataspace_name", raw, nil),
			"destination_info": utils.PathSearch("destination_info", raw, nil),
			"id":               utils.PathSearch("id", raw, nil),
			"identity":         utils.PathSearch("identity", raw, nil),
			"pipe":             utils.PathSearch("pipe", raw, nil),
			"pipe_name":        utils.PathSearch("pipe_name", raw, nil),
			"region":           utils.PathSearch("region", raw, nil),
			"type":             utils.PathSearch("type", raw, nil),
			"workspace":        utils.PathSearch("workspace", raw, nil),
			"workspace_name":   utils.PathSearch("workspace_name", raw, nil),
		},
	}
}

func flattenSiemShipperSource(raw interface{}) []interface{} {
	if raw == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"data_type":      utils.PathSearch("data_type", raw, nil),
			"dataspace":      utils.PathSearch("dataspace", raw, nil),
			"dataspace_name": utils.PathSearch("dataspace_name", raw, nil),
			"id":             utils.PathSearch("id", raw, nil),
			"identity":       utils.PathSearch("identity", raw, nil),
			"pipe":           utils.PathSearch("pipe", raw, nil),
			"pipe_name":      utils.PathSearch("pipe_name", raw, nil),
			"region":         utils.PathSearch("region", raw, nil),
			"type":           utils.PathSearch("type", raw, nil),
			"workspace":      utils.PathSearch("workspace", raw, nil),
			"workspace_name": utils.PathSearch("workspace_name", raw, nil),
		},
	}
}
