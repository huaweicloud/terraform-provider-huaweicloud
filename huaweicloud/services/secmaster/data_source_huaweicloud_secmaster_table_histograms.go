package secmaster

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

// @API Secmaster POST /v2/{project_id}/workspaces/{workspace_id}/siem/tables/{table_id}/histograms
func DataSourceTableHistograms() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTableHistogramsRead,

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
			"table_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"from": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"to": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"histograms": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"count": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"from": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"to": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildTableHistogramsBodyParams(d *schema.ResourceData) map[string]interface{} {
	rst := map[string]interface{}{
		"query": d.Get("query"),
	}

	if v, ok := d.GetOk("from"); ok {
		rst["from"] = v
	}

	if v, ok := d.GetOk("to"); ok {
		rst["to"] = v
	}

	return rst
}

func dataSourceTableHistogramsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v2/{project_id}/workspaces/{workspace_id}/siem/tables/{table_id}/histograms"
		workspaceID = d.Get("workspace_id").(string)
		tableID     = d.Get("table_id").(string)
	)

	client, err := cfg.NewServiceClient("secmaster", region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", workspaceID)
	requestPath = strings.ReplaceAll(requestPath, "{table_id}", tableID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildTableHistogramsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving SecMaster table histograms: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(dataSourceId)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("histograms", flattenHistograms(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHistograms(respBody interface{}) []interface{} {
	if respBody == nil {
		return nil
	}

	respArray := utils.PathSearch("histograms", respBody, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(respArray))
	for _, v := range respArray {
		rst = append(rst, map[string]interface{}{
			"count": utils.PathSearch("count", v, nil),
			"from":  utils.PathSearch("from", v, nil),
			"to":    utils.PathSearch("to", v, nil),
		})
	}

	return rst
}
