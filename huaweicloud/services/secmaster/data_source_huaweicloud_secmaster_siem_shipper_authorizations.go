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

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/siem/shippers/authorizations
func DataSourceSiemShipperAuthorizations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSiemShipperAuthorizationsRead,

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
			"source_region": {
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
			"auth_status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     siemShipperAuthorizationsDataSchema(),
			},
		},
	}
}

func siemShipperAuthorizationsDataSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"authorize_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"data_source_query": {
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
			"id": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"pipe": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"request_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"handler_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"run_status": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"shipper_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"shipper_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildSiemShipperAuthorizationsQueryParams(d *schema.ResourceData, offset int) string {
	// The `limit` and `offset` are required parameters.
	rst := fmt.Sprintf("?limit=1000&offset=%d", offset)

	if v, ok := d.GetOk("source_region"); ok {
		rst += fmt.Sprintf("&source_region=%v", v)
	}

	if v, ok := d.GetOk("destination_shipper_type"); ok {
		rst += fmt.Sprintf("&destination_shipper_type=%v", v)
	}

	if v, ok := d.GetOk("shipper_status"); ok {
		rst += fmt.Sprintf("&shipper_status=%v", v)
	}

	if v, ok := d.GetOk("auth_status"); ok {
		rst += fmt.Sprintf("&auth_status=%v", v)
	}

	return rst
}

func dataSourceSiemShipperAuthorizationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/workspaces/{workspace_id}/siem/shippers/authorizations"
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
		currentPath := requestPath + buildSiemShipperAuthorizationsQueryParams(d, offset)
		resp, err := client.Request("GET", currentPath, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster SIEM shipper authorizations: %s", err)
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

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenSiemShipperAuthorizations(allData)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSiemShipperAuthorizations(authorizations []interface{}) []interface{} {
	if len(authorizations) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(authorizations))
	for _, v := range authorizations {
		rst = append(rst, map[string]interface{}{
			"authorize_status":  utils.PathSearch("authorize_status", v, nil),
			"data_source_query": utils.PathSearch("data_source_query", v, nil),
			"data_type":         utils.PathSearch("data_type", v, nil),
			"dataspace":         utils.PathSearch("dataspace", v, nil),
			"id":                utils.PathSearch("id", v, nil),
			"pipe":              utils.PathSearch("pipe", v, nil),
			"region":            utils.PathSearch("region", v, nil),
			"request_time":      utils.PathSearch("request_time", v, nil),
			"handler_time":      utils.PathSearch("handler_time", v, nil),
			"run_status":        utils.PathSearch("run_status", v, nil),
			"shipper_id":        utils.PathSearch("shipper_id", v, nil),
			"shipper_name":      utils.PathSearch("shipper_name", v, nil),
			"workspace":         utils.PathSearch("workspace", v, nil),
		})
	}

	return rst
}
