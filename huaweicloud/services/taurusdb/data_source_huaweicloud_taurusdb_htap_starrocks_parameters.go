package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/configurations
func DataSourceTaurusDBHtapStarrocksParameters() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksParametersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"node_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configurations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksParametersConfigurationsSchema(),
			},
			"parameter_values": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksParametersParameterValuesSchema(),
			},
		},
	}
}

func starrocksParametersConfigurationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"configuration_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"datastore_version_name": {
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
		},
	}
}

func starrocksParametersParameterValuesSchema() *schema.Resource {
	return &schema.Resource{
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
}

func dataSourceTaurusDBHtapStarrocksParametersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		nodeType   = d.Get("node_type").(string)
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	configuration, parameterValues, err := getStarrocksParameters(client, instanceId, nodeType)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("configurations", flattenStarrocksParametersConfigurationsBody(configuration)),
		d.Set("parameter_values", flattenStarrocksParametersParameterValuesBody(parameterValues)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func getStarrocksParameters(client *golangsdk.ServiceClient, instanceId string, nodeType string) (interface{}, []interface{}, error) {
	var (
		offset         = 0
		result         = make([]interface{}, 0)
		configurations interface{}
	)
	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/configurations"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = fmt.Sprintf("%s?node_type=%s", listPath, nodeType)

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	for {
		currentQueryPath := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", currentQueryPath, &listOpts)
		if err != nil {
			return nil, nil, fmt.Errorf("error retrieving TaurusDB HTAP StarRocks (%s) Node (%s) Parameters: %s", instanceId, nodeType, err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, nil, err
		}

		// Save configurations from the first page response
		if configurations == nil {
			configurations = utils.PathSearch("configurations", respBody, nil)
		}

		parameterValues := utils.PathSearch("parameter_values", respBody, make([]interface{}, 0)).([]interface{})
		if len(parameterValues) == 0 {
			break
		}

		result = append(result, parameterValues...)

		if len(parameterValues) < 100 {
			break
		}

		offset += len(parameterValues)
	}
	return configurations, result, nil
}

func flattenStarrocksParametersConfigurationsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"configuration_id":       utils.PathSearch("configuration_id", resp, nil),
			"datastore_version_name": utils.PathSearch("datastore_version_name", resp, nil),
			"datastore_name":         utils.PathSearch("datastore_name", resp, nil),
			"created":                utils.PathSearch("created", resp, nil),
			"updated":                utils.PathSearch("updated", resp, nil),
		},
	}
}

func flattenStarrocksParametersParameterValuesBody(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":             utils.PathSearch("name", v, nil),
			"value":            utils.PathSearch("value", v, nil),
			"restart_required": utils.PathSearch("restart_required", v, nil),
			"readonly":         utils.PathSearch("readonly", v, nil),
			"value_range":      utils.PathSearch("value_range", v, nil),
			"type":             utils.PathSearch("type", v, nil),
			"description":      utils.PathSearch("description", v, nil),
		})
	}
	return res
}
