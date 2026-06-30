package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/instances
func DataSourceInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceInstancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the instances are located.`,
			},

			// Required parameters.
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the database.`,
			},

			// Attributes.
			"instances": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of instances that matched filter parameters.`,
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
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version of the instance.`,
						},
						"engine_type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The engine type of the instance.`,
						},
						"ip": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The IP address of the instance.`,
						},
						"port": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The port of the instance.`,
						},
						"cpu": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CPU cores of the instance.`,
						},
						"mem": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The memory size of the instance, in GB.`,
						},
						"login_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether login is enabled.`,
						},
						"slow_sql_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether slow SQL analysis is enabled.`,
						},
						"deadlock_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether deadlock analysis is enabled.`,
						},
						"lock_blocking_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether lock blocking analysis is enabled.`,
						},
						"charge_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether the instance is charged.`,
						},
						"full_sql_flag": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Whether full SQL is enabled.`,
						},
					},
				},
			},
		},
	}
}

func listInstances(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances?limit={limit}&datastore_type={datastore_type}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listPath = strings.ReplaceAll(listPath, "{datastore_type}", d.Get("datastore_type").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		instances := utils.PathSearch("instance_list", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, instances...)
		if len(instances) < limit {
			break
		}

		offset += len(instances)
	}

	return result, nil
}

func flattenInstances(instances []interface{}) []map[string]interface{} {
	if len(instances) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(instances))
	for _, instance := range instances {
		result = append(result, map[string]interface{}{
			"id":                 utils.PathSearch("instance_id", instance, nil),
			"name":               utils.PathSearch("instance_name", instance, nil),
			"status":             utils.PathSearch("instance_status", instance, nil),
			"version":            utils.PathSearch("version", instance, nil),
			"engine_type":        utils.PathSearch("engine_type", instance, nil),
			"ip":                 utils.PathSearch("ip", instance, nil),
			"port":               utils.PathSearch("port", instance, nil),
			"cpu":                utils.PathSearch("cpu", instance, nil),
			"mem":                utils.PathSearch("mem", instance, nil),
			"login_flag":         utils.PathSearch("login_flag", instance, nil),
			"slow_sql_flag":      utils.PathSearch("slow_sql_flag", instance, nil),
			"deadlock_flag":      utils.PathSearch("deadlock_flag", instance, nil),
			"lock_blocking_flag": utils.PathSearch("lock_blocking_flag", instance, nil),
			"charge_flag":        utils.PathSearch("charge_flag", instance, nil),
			"full_sql_flag":      utils.PathSearch("full_sql_flag", instance, nil),
		})
	}
	return result
}

func dataSourceInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	instances, err := listInstances(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS instances: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instances", flattenInstances(instances)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
