package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/starrocks/databases/replication
func DataSourceTaurusDBHtapStarrocksReplications() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapStarrocksReplicationsRead,

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
			"replications": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     starrocksReplicationsSchema(),
			},
		},
	}
}

func starrocksReplicationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"source_database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"target_database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"task_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"percentage": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_need_repair": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"is_main_task": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapStarrocksReplicationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/starrocks/databases/replication"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		queryPath := fmt.Sprintf("%s?limit=%d&offset=%d", listPath, limit, offset)
		resp, err := client.Request("GET", queryPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP StarRocks replications: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		replications := utils.PathSearch("replications", respBody, make([]interface{}, 0)).([]interface{})
		if len(replications) == 0 {
			break
		}

		result = append(result, replications...)

		offset += len(replications)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("replications", flattenStarrocksReplications(result)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStarrocksReplications(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"source_database": utils.PathSearch("source_database", v, nil),
			"target_database": utils.PathSearch("target_database", v, nil),
			"task_name":       utils.PathSearch("task_name", v, nil),
			"status":          utils.PathSearch("status", v, nil),
			"stage":           utils.PathSearch("stage", v, nil),
			"percentage":      utils.PathSearch("percentage", v, nil),
			"is_need_repair":  utils.PathSearch("is_need_repair", v, nil),
			"is_main_task":    utils.PathSearch("is_main_task", v, nil),
		})
	}
	return res
}
