package rds

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RDS GET /v3/{project_id}/jobs?id={job_id}
func DataSourceRdsListJobInfo() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsListJobInfoRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     jobSchema(),
			},
		},
	}
}

func jobSchema() *schema.Resource {
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
			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"ended": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"process": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fail_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entities": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instances": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceSchema(),
			},
		},
	}
}

func instanceSchema() *schema.Resource {
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
		},
	}
}

func dataSourceRdsListJobInfoRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	jobID := d.Get("job_id").(string)

	client, err := cfg.NewServiceClient("rds", region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	url := client.Endpoint + "v3/{project_id}/jobs?id={job_id}"
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{job_id}", jobID)

	opts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", url, &opts)
	if err != nil {
		return diag.Errorf("error retrieving RDS job %q: %s", jobID, err)
	}

	body, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating resource ID: %s", err)
	}
	d.SetId(id)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("job", flattenJob(body)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJob(resp interface{}) []interface{} {
	raw := utils.PathSearch("job", resp, nil)
	if raw == nil {
		return nil
	}

	job, ok := raw.(map[string]interface{})
	if !ok {
		return nil
	}

	instanceID := utils.PathSearch("instance.id", job, nil)
	instanceName := utils.PathSearch("instance.name", job, nil)

	instanceList := []interface{}{}
	if instanceID != nil || instanceName != nil {
		instanceList = append(instanceList, map[string]interface{}{
			"id":   instanceID,
			"name": instanceName,
		})
	}

	entitiesRaw := utils.PathSearch("entities", job, nil)
	var entitiesStr string
	if entitiesRaw != nil {
		entitiesJSON, err := json.Marshal(entitiesRaw)
		if err != nil {
			entitiesStr = "error marshaling entities"
		} else {
			entitiesStr = string(entitiesJSON)
		}
	}

	return []interface{}{
		map[string]interface{}{
			"id":          utils.PathSearch("id", job, nil),
			"name":        utils.PathSearch("name", job, nil),
			"status":      utils.PathSearch("status", job, nil),
			"created":     utils.PathSearch("created", job, nil),
			"ended":       utils.PathSearch("ended", job, nil),
			"process":     utils.PathSearch("process", job, nil),
			"fail_reason": utils.PathSearch("fail_reason", job, nil),
			"entities":    entitiesStr,
			"instances":   instanceList,
		},
	}
}
