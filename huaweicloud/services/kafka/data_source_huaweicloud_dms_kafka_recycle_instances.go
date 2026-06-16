package kafka

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Kafka GET /v2/{project_id}/recycle
func DataSourceRecycleInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecycleInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the recycle bin instances are located.`,
			},

			// Attributes.
			"retention_days": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The retention days of the recycle bin.`,
			},
			"default_use_recycle": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the recycle bin is enabled.`,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the instance.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the instance.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the instance.`,
						},
						"engine": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The engine of the instance.`,
						},
						"in_recycle_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the instance was placed in the recycle bin, in RFC3339 format.`,
						},
						"save_time": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The time when the instance was saved, in days.`,
						},
						"auto_delete_time": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The time when the instance was automatically deleted, in RFC3339 format.`,
						},
						"cost_per_hour": {
							Type:        schema.TypeFloat,
							Computed:    true,
							Description: `The cost per hour of the instance.`,
						},
						"error_message": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The error message.`,
						},
						"product_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the flavor of the instance.`,
						},
					},
				},
				Description: `The list of recycle instances.`,
			},
		},
	}
}

func listRecycleInstances(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/recycle"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func dataSourceRecycleInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dms", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	respBody, err := listRecycleInstances(client)
	if err != nil {
		return diag.Errorf("error retrieving recycle bin instance list: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate data source ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("retention_days", utils.PathSearch("retention_days", respBody, nil)),
		d.Set("default_use_recycle", utils.PathSearch("default_use_recycle", respBody, nil)),
		d.Set("instances", flattenRecycleInstances(utils.PathSearch("recycle_instances", respBody,
			make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenRecycleInstances(recycleInstances []interface{}) []map[string]interface{} {
	if len(recycleInstances) == 0 {
		return nil
	}

	results := make([]map[string]interface{}, 0, len(recycleInstances))
	for _, item := range recycleInstances {
		results = append(results, map[string]interface{}{
			"id":     utils.PathSearch("instance_id", item, nil),
			"name":   utils.PathSearch("name", item, nil),
			"status": utils.PathSearch("status", item, nil),
			"engine": utils.PathSearch("engine", item, nil),
			"in_recycle_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("in_recycle_time",
				item, float64(0)).(float64))/1000, false),
			"save_time": utils.PathSearch("save_time", item, nil),
			"auto_delete_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("auto_delete_time",
				item, float64(0)).(float64))/1000, false),
			"cost_per_hour": utils.PathSearch("cost_per_hour", item, nil),
			"error_message": utils.PathSearch("error_message", item, nil),
			"product_id":    utils.PathSearch("product_id", item, nil),
		})
	}
	return results
}
