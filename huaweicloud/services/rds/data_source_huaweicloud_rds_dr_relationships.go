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

// @API RDS POST /v3/{project_id}/instances/disaster-recovery-infos
func DataSourceRdsDrRelationships() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsDrRelationshipsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"relationship_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"master_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slave_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slave_region": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_at_start": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"create_at_end": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_field": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_dr_infos": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"failed_message": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"master_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_instance_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"slave_region": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"build_process": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"replica_state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wal_write_receive_delay_in_mb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wal_write_replay_delay_in_mb": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"wal_receive_replay_delay_in_ms": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRdsDrRelationshipsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/disaster-recovery-infos"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	offset := 0
	limit := 100
	res := make([]map[string]interface{}, 0)
	for {
		getOpt.JSONBody = utils.RemoveNil(buildGetDrRelationshipsParams(d, limit, offset))
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS DR relationships: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		drRelationships := flattenRdsGetDrRelationships(getRespBody)
		if len(drRelationships) == 0 {
			break
		}
		res = append(res, drRelationships...)
		offset += limit
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("instance_dr_infos", res),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetDrRelationshipsParams(d *schema.ResourceData, limit, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"id":                 utils.ValueIgnoreEmpty(d.Get("relationship_id").(string)),
		"status":             utils.ValueIgnoreEmpty(d.Get("status").(string)),
		"master_instance_id": utils.ValueIgnoreEmpty(d.Get("master_instance_id").(string)),
		"master_region":      utils.ValueIgnoreEmpty(d.Get("master_region").(string)),
		"slave_instance_id":  utils.ValueIgnoreEmpty(d.Get("slave_instance_id").(string)),
		"slave_region":       utils.ValueIgnoreEmpty(d.Get("slave_region").(string)),
		"create_at_start":    utils.ValueIgnoreEmpty(d.Get("create_at_start").(int)),
		"create_at_end":      utils.ValueIgnoreEmpty(d.Get("create_at_end").(int)),
		"order":              utils.ValueIgnoreEmpty(d.Get("order").(string)),
		"sort_field":         utils.ValueIgnoreEmpty(d.Get("sort_field").(string)),
		"limit":              limit,
		"offset":             offset,
	}
	return bodyParams
}

func flattenRdsGetDrRelationships(resp interface{}) []map[string]interface{} {
	curJson := utils.PathSearch("instance_dr_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	if len(curArray) == 0 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(curArray))
	for _, v := range curArray {
		result = append(result, map[string]interface{}{
			"id":                             utils.PathSearch("id", v, nil),
			"status":                         utils.PathSearch("status", v, nil),
			"failed_message":                 utils.PathSearch("failed_message", v, nil),
			"master_instance_id":             utils.PathSearch("master_instance_id", v, nil),
			"master_region":                  utils.PathSearch("master_region", v, nil),
			"slave_instance_id":              utils.PathSearch("slave_instance_id", v, nil),
			"slave_region":                   utils.PathSearch("slave_region", v, nil),
			"build_process":                  utils.PathSearch("build_process", v, nil),
			"time":                           utils.PathSearch("time", v, nil),
			"replica_state":                  utils.PathSearch("replica_state", v, nil),
			"wal_write_receive_delay_in_mb":  utils.PathSearch("wal_write_receive_delay_in_mb", v, nil),
			"wal_write_replay_delay_in_mb":   utils.PathSearch("wal_write_replay_delay_in_mb", v, nil),
			"wal_receive_replay_delay_in_ms": utils.PathSearch("wal_receive_replay_delay_in_ms", v, nil),
		})
	}
	return result
}
