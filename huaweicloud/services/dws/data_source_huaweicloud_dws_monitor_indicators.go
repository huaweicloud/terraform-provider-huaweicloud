package dws

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

// @API DWS GET /v1.0/{project_id}/dms/metric-data/indicators
func DataSourceMonitorIndicators() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceMonitorIndicatorsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the monitor indicators are located.`,
			},

			// Attributes.
			"indicators": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        monitorIndicatorsSchema(),
				Description: `The list of monitor indicators.`,
			},
		},
	}
}

func monitorIndicatorsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"indicator_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The monitor indicator name.`,
			},
			"plugin_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The collection plugin name.`,
			},
			"default_collect_rate": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The default collection rate.`,
			},
			"support_datastore_version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The supported datastore version.`,
			},
		},
	}
}

func listMonitorIndicators(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v1.0/{project_id}/dms/metric-data/indicators"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	indicators, ok := respBody.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format, expected list but got %T", respBody)
	}

	return indicators, nil
}

func flattenMonitorIndicators(indicators []interface{}) []map[string]interface{} {
	if len(indicators) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(indicators))
	for _, indicator := range indicators {
		result = append(result, map[string]interface{}{
			"indicator_name":            utils.PathSearch("indicator_name", indicator, nil),
			"plugin_name":               utils.PathSearch("plugin_name", indicator, nil),
			"default_collect_rate":      utils.PathSearch("default_collect_rate", indicator, nil),
			"support_datastore_version": utils.PathSearch("support_datastore_version", indicator, nil),
		})
	}

	return result
}

func dataSourceMonitorIndicatorsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	indicators, err := listMonitorIndicators(client)
	if err != nil {
		return diag.Errorf("error retrieving monitor indicators: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("indicators", flattenMonitorIndicators(indicators)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
