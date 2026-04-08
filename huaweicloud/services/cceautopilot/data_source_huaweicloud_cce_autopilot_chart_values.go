package cceautopilot

import (
	"context"
	"encoding/json"
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

// @API CCE GET /autopilot/v2/charts/{chart_id}/values
func DataSourceCceAutopilotShowChartValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCceAutopilotShowChartValuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				Description: "Specifies the region in which to query the resource. " +
					"If omitted, the provider-level region will be used.",
			},
			"chart_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of CCE autopilot chart.",
			},
			"values": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The values of CCE autopilot chart template.",
			},
		},
	}
}

func dataSourceCceAutopilotShowChartValuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "cce"
		httpUrl = "autopilot/v2/charts/{chart_id}/values"
		chartID = d.Get("chart_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{chart_id}", chartID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving CCE chart values: %s", err)
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
		d.Set("values", flattenCCEAutopilotChartValues(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCCEAutopilotChartValues(resp interface{}) map[string]interface{} {
	if resp == nil {
		return nil
	}
	values, ok := resp.(map[string]interface{})
	if !ok {
		return nil
	}
	result := make(map[string]interface{})
	for k, v := range values {
		jsonBytes, err := json.Marshal(v)
		if err != nil {
			// If marshal fails, use string representation as fallback
			result[k] = fmt.Sprintf("%v", v)
		} else {
			result[k] = string(jsonBytes)
		}
	}
	return result
}
