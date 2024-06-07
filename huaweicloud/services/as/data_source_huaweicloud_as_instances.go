// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product AS
// ---------------------------------------------------------------

package as

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

// @API AS GET /autoscaling-api/v1/{project_id}/scaling_group_instance/{scaling_group_id}/list
func DataSourceASInstances() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceASInstancesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the AS group ID.`,
			},
			"life_cycle_state": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance lifecycle status in the AS group.`,
			},
			"health_status": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance health status.`,
			},
			"protect_from_scaling_down": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the instance protection status.`,
			},
			"instances": {
				Type:        schema.TypeList,
				Elem:        dataSourceInstancesSchema(),
				Computed:    true,
				Description: `The details about the instances in the AS group.`,
			},
		},
	}
}

func dataSourceInstancesSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance ID.`,
			},
			"instance_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance name.`,
			},
			"scaling_group_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the AS group to which the instance belongs.`,
			},
			"scaling_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the AS group to which the instance belongs.`,
			},
			"life_cycle_state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance lifecycle status in the AS group.`,
			},
			"health_status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the instance health status.`,
			},
			"scaling_configuration_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the AS configuration name.`,
			},
			"scaling_configuration_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the AS configuration ID.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the time when the instance is added to the AS group. The time format complies with UTC.`,
			},
			"protect_from_scaling_down": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Indicates the instance protection status.`,
			},
		},
	}
	return &sc
}

func resourceASInstancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		mErr         *multierror.Error
		httpUrl      = "autoscaling-api/v1/{project_id}/scaling_group_instance/{scaling_group_id}/list"
		product      = "autoscaling"
		startNumber  = 0
		allInstances []interface{}
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating AS client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{scaling_group_id}", d.Get("scaling_group_id").(string))
	basePath += buildGetASInstancesQueryParams(d)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		getPath := fmt.Sprintf("%s&start_number=%d", basePath, startNumber)
		getResp, err := client.Request("GET", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving AS instances: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		instances := utils.PathSearch("scaling_group_instances", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(instances) == 0 {
			break
		}

		allInstances = append(allInstances, instances...)
		startNumber += len(instances)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instances", flattenDataSourceInstances(allInstances)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDataSourceInstances(allInstances []interface{}) []interface{} {
	rst := make([]interface{}, 0, len(allInstances))
	for _, v := range allInstances {
		rst = append(rst, map[string]interface{}{
			"instance_id":                utils.PathSearch("instance_id", v, nil),
			"instance_name":              utils.PathSearch("instance_name", v, nil),
			"scaling_group_id":           utils.PathSearch("scaling_group_id", v, nil),
			"scaling_group_name":         utils.PathSearch("scaling_group_name", v, nil),
			"life_cycle_state":           utils.PathSearch("life_cycle_state", v, nil),
			"health_status":              utils.PathSearch("health_status", v, nil),
			"scaling_configuration_name": utils.PathSearch("scaling_configuration_name", v, nil),
			"scaling_configuration_id":   utils.PathSearch("scaling_configuration_id", v, nil),
			"created_at":                 utils.PathSearch("create_time", v, nil),
			"protect_from_scaling_down":  utils.PathSearch("protect_from_scaling_down", v, nil),
		})
	}
	return rst
}

func buildGetASInstancesQueryParams(d *schema.ResourceData) string {
	res := "?limit=100"
	if v, ok := d.GetOk("life_cycle_state"); ok {
		res = fmt.Sprintf("%s&life_cycle_state=%v", res, v)
	}

	if v, ok := d.GetOk("health_status"); ok {
		res = fmt.Sprintf("%s&health_status=%v", res, v)
	}

	if v, ok := d.GetOk("protect_from_scaling_down"); ok {
		res = fmt.Sprintf("%s&protect_from_scaling_down=%v", res, v)
	}
	return res
}
