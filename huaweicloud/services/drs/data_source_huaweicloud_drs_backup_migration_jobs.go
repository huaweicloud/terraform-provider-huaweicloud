package drs

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

// @API DRS GET /v5/{project_id}/backup-migration-jobs
func DataSourceBackupMigrationJobs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBackupMigrationJobsRead,

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
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dbs_instance_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"create_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"completed_at": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     backupMigrationJobsSchema(),
			},
		},
	}
}

func backupMigrationJobsSchema() *schema.Resource {
	return &schema.Resource{
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
			"engine_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"error_log": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"finish_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBackupMigrationJobsQueryParams(d *schema.ResourceData, offset int, epsId string) string {
	queryParams := "?limit=2000"

	if v, ok := d.GetOk("name"); ok {
		queryParams += fmt.Sprintf("&name=%s", v.(string))
	}
	if v, ok := d.GetOk("status"); ok {
		queryParams += fmt.Sprintf("&status=%s", v.(string))
	}
	if v := d.Get("dbs_instance_ids").([]interface{}); len(v) > 0 {
		for _, instanceId := range v {
			queryParams += fmt.Sprintf("&dbs_instance_ids=%s", instanceId)
		}
	}
	if v, ok := d.GetOk("description"); ok {
		queryParams += fmt.Sprintf("&description=%s", v.(string))
	}
	if v, ok := d.GetOk("create_at"); ok {
		queryParams += fmt.Sprintf("&create_at=%s", v.(string))
	}
	if v, ok := d.GetOk("completed_at"); ok {
		queryParams += fmt.Sprintf("&completed_at=%s", v.(string))
	}
	if epsId != "" {
		queryParams += fmt.Sprintf("&enterprise_project_id=%s", epsId)
	}
	if v, ok := d.GetOk("tags"); ok {
		queryParams += fmt.Sprintf("&tags=%s", v.(string))
	}
	if v, ok := d.GetOk("sort_key"); ok {
		queryParams += fmt.Sprintf("&sort_key=%s", v.(string))
	}
	if v, ok := d.GetOk("sort_dir"); ok {
		queryParams += fmt.Sprintf("&sort_dir=%s", v.(string))
	}
	if offset > 0 {
		queryParams += fmt.Sprintf("&offset=%d", offset)
	}

	return queryParams
}

func dataSourceBackupMigrationJobsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/backup-migration-jobs"
		offset  = 0
		epsId   = cfg.GetEnterpriseProjectID(d)
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		requestPathWithQuery := requestPath + buildBackupMigrationJobsQueryParams(d, offset, epsId)
		resp, err := client.Request("GET", requestPathWithQuery, &requestOpt)
		if err != nil {
			return diag.Errorf("error retrieving DRS backup migration jobs: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobsResp := utils.PathSearch("jobs", respBody, make([]interface{}, 0)).([]interface{})
		if len(jobsResp) == 0 {
			break
		}

		result = append(result, jobsResp...)
		offset += len(jobsResp)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("jobs", flattenBackupMigrationJobs(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBackupMigrationJobs(respArray []interface{}) []interface{} {
	if len(respArray) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(respArray))
	for _, item := range respArray {
		result = append(result, map[string]interface{}{
			"id":                    utils.PathSearch("id", item, nil),
			"name":                  utils.PathSearch("name", item, nil),
			"status":                utils.PathSearch("status", item, nil),
			"engine_type":           utils.PathSearch("engine_type", item, nil),
			"error_log":             utils.PathSearch("error_log", item, nil),
			"description":           utils.PathSearch("description", item, nil),
			"create_time":           utils.PathSearch("create_time", item, nil),
			"finish_time":           utils.PathSearch("finish_time", item, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", item, nil),
		})
	}

	return result
}
