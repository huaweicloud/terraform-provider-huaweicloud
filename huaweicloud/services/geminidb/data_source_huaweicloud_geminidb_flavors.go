package geminidb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GeminiDB GET /v3.1/{project_id}/flavors
func DataSourceGeminiDBFlavors() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGeminiDBFlavorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"engine_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"product_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"flavors": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     geminiDBFlavorSchema(),
			},
		},
	}
}

func geminiDBFlavorSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"engine_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"engine_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"vcpus": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ram": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"spec_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"az_status": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceGeminiDBFlavorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	httpUrl := "v3.1/{project_id}/flavors"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListGeminiDBFlavorsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB flavors: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	flavors := flattenListGeminiDBFlavors(listRespBody)
	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("flavors", flavors),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListGeminiDBFlavorsQueryParams(d *schema.ResourceData) string {
	queryParams := ""

	if v, ok := d.GetOk("engine_name"); ok {
		queryParams = fmt.Sprintf("%s&engine_name=%v", queryParams, v)
	}
	if v, ok := d.GetOk("mode"); ok {
		queryParams = fmt.Sprintf("%s&mode=%v", queryParams, v)
	}
	if v, ok := d.GetOk("product_type"); ok {
		queryParams = fmt.Sprintf("%s&product_type=%v", queryParams, v)
	}

	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func flattenListGeminiDBFlavors(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("flavors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"engine_name":    utils.PathSearch("engine_name", v, nil),
			"engine_version": utils.PathSearch("engine_version", v, nil),
			"vcpus":          utils.PathSearch("vcpus", v, nil),
			"ram":            utils.PathSearch("ram", v, nil),
			"spec_code":      utils.PathSearch("spec_code", v, nil),
			"az_status":      utils.PathSearch("az_status", v, nil),
		})
	}
	return rst
}
