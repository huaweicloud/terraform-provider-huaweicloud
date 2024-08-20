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

// @API DWS GET /v2/{project_id}/disaster-recoveries
func DataSourceDisasterRecoveryTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDisasterRecoveryTasksRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"primary_cluster_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standby_cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standby_cluster_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Elem:     disasterClusterSchema(),
				Computed: true,
			},
		},
	}
}

func disasterClusterSchema() *schema.Resource {
	nodeResource := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"dr_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"primary_cluster_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_role": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_region": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"standby_cluster_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"last_disaster_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}

	return &nodeResource
}

func resourceDisasterRecoveryTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dws"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}
	httpUrl := "v2/{project_id}/disaster-recoveries"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DWS disaster recoveries: %s", err)
	}
	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing DWS disaster recoveries: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	disasterList := utils.PathSearch("disaster_recovery", respBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tasks", filterDisasterRecoveryTasks(flattenDisasterRecoveryTasks(disasterList), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDisasterRecoveryTasks(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                         utils.PathSearch("id", v, nil),
			"name":                       utils.PathSearch("name", v, nil),
			"status":                     utils.PathSearch("status", v, nil),
			"dr_type":                    utils.PathSearch("dr_type", v, nil),
			"primary_cluster_id":         utils.PathSearch("primary_cluster_id", v, nil),
			"primary_cluster_name":       utils.PathSearch("primary_cluster_name", v, nil),
			"primary_cluster_role":       utils.PathSearch("primary_cluster_role", v, nil),
			"primary_cluster_status":     utils.PathSearch("primary_cluster_status", v, nil),
			"primary_cluster_region":     utils.PathSearch("primary_cluster_region", v, nil),
			"primary_cluster_project_id": utils.PathSearch("primary_cluster_project_id", v, nil),
			"standby_cluster_id":         utils.PathSearch("standby_cluster_id", v, nil),
			"standby_cluster_name":       utils.PathSearch("standby_cluster_name", v, nil),
			"standby_cluster_role":       utils.PathSearch("standby_cluster_role", v, nil),
			"standby_cluster_status":     utils.PathSearch("standby_cluster_status", v, nil),
			"standby_cluster_region":     utils.PathSearch("standby_cluster_region", v, nil),
			"standby_cluster_project_id": utils.PathSearch("standby_cluster_project_id", v, nil),
			"last_disaster_time":         utils.PathSearch("last_disaster_time", v, nil),
			"start_at":                   utils.PathSearch("start_time", v, nil),
			"create_at":                  utils.PathSearch("create_time", v, nil),
		})
	}
	return rst
}

func filterDisasterRecoveryTasks(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))
	for _, v := range all {
		if param, ok := d.GetOk("name"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("dr_type"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("dr_type", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("status"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("primary_cluster_name"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("primary_cluster_name", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("primary_cluster_region"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("primary_cluster_region", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("standby_cluster_name"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("standby_cluster_name", v, nil)) {
				continue
			}
		}
		if param, ok := d.GetOk("standby_cluster_region"); ok {
			if fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("standby_cluster_region", v, nil)) {
				continue
			}
		}

		rst = append(rst, v)
	}
	return rst
}
