package cbh

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

// @API CBH GET /v2/{project_id}/cbs/instance/specification
func DataSourceCbhFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourceFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavor_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asset": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"max_connection": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Elem:     flavorSchema(),
				Computed: true,
			},
		},
	}
}

func flavorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ecs_system_data_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"asset": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"max_connection": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"data_disk_size": {
				Type:     schema.TypeFloat,
				Computed: true,
			},
		},
	}
	return &sc
}

func buildFlavorsReadPath(getCBHFlavorsClient *golangsdk.ServiceClient, getCBHFlavorsHttpUrl string, d *schema.ResourceData) string {
	getCBHFlavorsPath := getCBHFlavorsClient.Endpoint + getCBHFlavorsHttpUrl
	getCBHFlavorsPath = strings.ReplaceAll(getCBHFlavorsPath, "{project_id}",
		getCBHFlavorsClient.ProjectID)

	// If the `action` is set to **update**, the `spec_code` parameter is required.
	// If the `action` is not filled in, set it to **create** query.
	if action, ok := d.GetOk("action"); ok {
		getCBHFlavorsPath += "?action=" + fmt.Sprintf("%v", action)
		if action == "update" {
			getCBHFlavorsPath += "&spec_code=" + fmt.Sprintf("%v", d.Get("spec_code"))
		}
	} else {
		getCBHFlavorsPath += "?action=create"
	}

	return getCBHFlavorsPath
}

func datasourceFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// get CBH Flavors: Query CBH Flavors
	var (
		getCBHFlavorsHttpUrl = "v2/{project_id}/cbs/instance/specification"
		getCBHFlavorsProduct = "cbh"
	)
	getCBHFlavorsClient, err := cfg.NewServiceClient(getCBHFlavorsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CBH client: %s", err)
	}

	getCBHFlavorsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	getCBHFlavorsResp, err := getCBHFlavorsClient.Request("GET",
		buildFlavorsReadPath(getCBHFlavorsClient, getCBHFlavorsHttpUrl, d), &getCBHFlavorsOpt)

	if err != nil {
		return diag.Errorf("error retrieving CBH flavors, %s", err)
	}

	getCBHFlavorsRespBody, err := utils.FlattenResponse(getCBHFlavorsResp)
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
		d.Set("flavors", filterGetFlavorsResponseBodyFlavor(
			flattenListFlavorsBody(getCBHFlavorsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func filterGetFlavorsResponseBodyFlavor(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("flavor_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("type"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("type", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("asset"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("asset", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("memory"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("memory", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("vcpus"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("vcpus", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("max_connection"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("max_connection", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}
	return rst
}

func flattenListFlavorsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	// The responseBody is in the form of an array
	curArray := resp.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"id":                   utils.PathSearch("resource_spec_code", v, nil),
			"ecs_system_data_size": utils.PathSearch("ecs_system_data_size", v, nil),
			"vcpus":                utils.PathSearch("cpu", v, nil),
			"memory":               utils.PathSearch("ram", v, nil),
			"asset":                utils.PathSearch("asset", v, nil),
			"max_connection":       utils.PathSearch("connection", v, nil),
			"type":                 utils.PathSearch("type", v, nil),
			"data_disk_size":       utils.PathSearch("data_disk_size", v, nil),
		})
	}
	return rst
}
