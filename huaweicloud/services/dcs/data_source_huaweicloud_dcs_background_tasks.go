package dcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/tasks
func DataSourceDcsBackgroundTasks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsBackgroundTasksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tasks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsCenterTaskSchema(),
			},
		},
	}
}

func dataSourceDcsBackgroundTasksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/tasks"
		product = "dcs"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	basePath = strings.ReplaceAll(basePath, "{instance_id}", d.Get("instance_id").(string))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	var offset int
	result := make([]interface{}, 0)

	for {
		path := basePath + buildQueryParams(d, offset)
		resp, err := client.Request("GET", path, &opt)
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving background tasks")
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		tasks := utils.PathSearch("tasks", respBody, make([]interface{}, 0)).([]interface{})

		if len(tasks) == 0 {
			break
		}

		result = append(result, tasks...)
		offset += len(tasks)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("tasks", flattenListCenterTasksBody(map[string]interface{}{
			"tasks": result,
		})),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildQueryParams(d *schema.ResourceData, offset int) string {
	res := fmt.Sprintf("?limit=10&offset=%d", offset)
	if v, ok := d.GetOk("begin_time"); ok {
		res += "&begin_time=" + v.(string)
	}
	if v, ok := d.GetOk("end_time"); ok {
		res += "&end_time=" + v.(string)
	}
	return res
}
