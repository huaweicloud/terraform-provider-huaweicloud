package ddm

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

// @API DDM GET /v3/{project_id}/configurations/{config_id}
func DataSourceDdmConfigurationParameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDdmConfigurationParametersRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"config_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"configuration_parameters": {
				Type:     schema.TypeList,
				Elem:     parameterConfigurationParametersSchema(),
				Computed: true,
			},
		},
	}
}

func parameterConfigurationParametersSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"restart_required": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"readonly": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"value_range": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
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

func resourceDdmConfigurationParametersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/configurations/{config_id}"
		product = "ddm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DDM client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{config_id}", fmt.Sprintf("%v", d.Get("config_id")))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DDM configuration parameters")
	}

	getRespBody, err := utils.FlattenResponse(getResp)
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
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getRespBody, nil)),
		d.Set("datastore_name", utils.PathSearch("datastore_name", getRespBody, nil)),
		d.Set("created", utils.PathSearch("created", getRespBody, nil)),
		d.Set("updated", utils.PathSearch("updated", getRespBody, nil)),
		d.Set("configuration_parameters", flattenGetConfigurationParametersResponseBody(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetConfigurationParametersResponseBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("configuration_parameters", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"value":            utils.PathSearch("value", v, nil),
			"restart_required": utils.PathSearch("restart_required", v, nil),
			"readonly":         utils.PathSearch("readonly", v, nil),
			"value_range":      utils.PathSearch("value_range", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"description":      utils.PathSearch("description", v, nil),
		})
	}
	return rst
}
