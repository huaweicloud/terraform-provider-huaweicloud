package rocketmq

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

// @API RocketMQ GET /v2/{project_id}/recycle
func DataSourceRocketmqRecycleInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRocketmqRecycleInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the RocketMQ recycle instances are located.`,
			},

			// Attributes.
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        recycleInstanceSchema(),
				Description: `The list of the recycle instances.`,
			},
		},
	}
}

func recycleInstanceSchema() *schema.Resource {
	return &schema.Resource{
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
				Description: `The message engine type.`,
			},
			"in_recycle_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the instance was recycled, in RFC3339 format.`,
			},
			"save_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of days the instance is retained in the recycle bin.`,
			},
			"auto_delete_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the instance will be automatically deleted, in RFC3339 format.`,
			},
		},
	}
}

func flattenRecycleInstances(recycleInstances []interface{}) []map[string]interface{} {
	if len(recycleInstances) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(recycleInstances))
	for _, instance := range recycleInstances {
		result = append(result, map[string]interface{}{
			"id":               utils.PathSearch("instance_id", instance, nil),
			"name":             utils.PathSearch("name", instance, nil),
			"status":           utils.PathSearch("status", instance, nil),
			"engine":           utils.PathSearch("engine", instance, nil),
			"in_recycle_time":  utils.FormatTimeStampRFC3339(int64(utils.PathSearch("in_recycle_time", instance, float64(0)).(float64))/1000, false),
			"save_time":        utils.PathSearch("save_time", instance, nil),
			"auto_delete_time": utils.FormatTimeStampRFC3339(int64(utils.PathSearch("auto_delete_time", instance, float64(0)).(float64))/1000, false),
		})
	}

	return result
}

func dataSourceRocketmqRecycleInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/recycle"
	)

	client, err := cfg.NewServiceClient("dmsv2", region)
	if err != nil {
		return diag.Errorf("error creating DMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving recycle instances: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenRecycleInstances(utils.PathSearch("recycle_instances", respBody, make([]interface{}, 0)).([]interface{}))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
