package rds

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

// @API RDS GET /v3/{project_id}/instances/{instance_id}/configurations
func DataSourceRdsInstanceConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsInstanceConfigurationsRead,
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
			"configuration_parameters": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     configurationParametersSchema(),
			},
		},
	}
}

func configurationParametersSchema() *schema.Resource {
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

func dataSourceRdsInstanceConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/configurations"
	getUrl := client.Endpoint + httpUrl
	getUrl = strings.ReplaceAll(getUrl, "{project_id}", client.ProjectID)
	getUrl = strings.ReplaceAll(getUrl, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
			"Accept":       "application/json",
		},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getUrl, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving RDS instance configuration list: %s", err)
	}

	body, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("datastore_version_name", utils.PathSearch("datastore_version_name", body, nil)),
		d.Set("datastore_name", utils.PathSearch("datastore_name", body, nil)),
		d.Set("created", utils.PathSearch("created", body, nil)),
		d.Set("updated", utils.PathSearch("updated", body, nil)),
		d.Set("configuration_parameters", flattenInstanceConfigurationParameters(body)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceConfigurationParameters(cp interface{}) []interface{} {
	if cp == nil {
		return nil
	}

	raw := utils.PathSearch("configuration_parameters", cp, nil)

	if raw == nil {
		return nil
	}

	params, ok := raw.([]interface{})
	if !ok || len(params) == 0 {
		return nil
	}

	out := make([]interface{}, 0, len(params))
	for _, p := range params {
		out = append(out, map[string]interface{}{
			"name":             utils.PathSearch("name", p, nil),
			"value":            utils.PathSearch("value", p, nil),
			"restart_required": utils.PathSearch("restart_required", p, nil),
			"readonly":         utils.PathSearch("readonly", p, nil),
			"value_range":      utils.PathSearch("value_range", p, nil),
			"type":             utils.PathSearch("type", p, nil),
			"description":      utils.PathSearch("description", p, nil),
		})
	}
	return out
}
