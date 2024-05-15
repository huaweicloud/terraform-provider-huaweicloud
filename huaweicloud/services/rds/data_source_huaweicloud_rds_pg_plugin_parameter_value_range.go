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

var (
	defaultValues = map[string][]string{
		"shared_preload_libraries": {"passwordcheck.so", "pg_stat_statements", "pg_sql_history", "pgaudit"},
	}
)

// @API RDS GET /v3/{project_id}/instances/{instance_id}/parameter/{name}
func DataSourceRdsPgPluginParameterValues() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourcePgPluginParameterValuesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the RDS instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the parameter name.`,
			},
			"restart_required": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates whether a reboot is required.`,
			},
			"default_values": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of default parameter values.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the list of parameter values.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourcePgPluginParameterValuesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/parameter/{name}"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{name}", d.Get("name").(string))

	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return diag.Errorf("error retrieving RDS PostgreSQL plugin parameter values: %s", err)
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}

	valueRangeRaw := utils.PathSearch("value_range", respBody, nil)
	if valueRangeRaw == nil {
		return diag.Errorf("error getting RDS PostgreSQL plugin parameter values, it is empty")
	}

	valueRange := strings.Split(valueRangeRaw.(string), ",")
	defaults := defaultValues[d.Get("name").(string)]
	values := make([]string, 0, len(valueRange)-len(defaults))
	for _, value := range valueRange {
		if !utils.StrSliceContains(defaults, value) {
			values = append(values, value)
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("restart_required", utils.PathSearch("restart_required", respBody, nil)),
		d.Set("default_values", defaults),
		d.Set("values", values),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
