package gaussdb

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/asp
func DataSourceGaussDbAspCollectionResults() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceGaussDbAspCollectionResultsRead,

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
			"start_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"job_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"asp": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbAspCollectionResultsAspSchema(),
			},
		},
	}
}

func gaussDbAspCollectionResultsAspSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"job_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"file_path": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"file_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"start_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"end_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"download_url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"obs_bucket": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     gaussDbAspCollectionResultsAspObsBucketSchema(),
			},
		},
	}
}

func gaussDbAspCollectionResultsAspObsBucketSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"url": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"domain_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGaussDbAspCollectionResultsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/asp"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildGetAspCollectionResultsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving GaussDB ASP collection results: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("asp", flattenGetAspCollectionResultsAspResponseBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAspCollectionResultsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("start_time"); ok {
		res = fmt.Sprintf("%s&start_time=%v", res, strings.ReplaceAll(v.(string), "+", "%2B"))
	}
	if v, ok := d.GetOk("end_time"); ok {
		res = fmt.Sprintf("%s&end_time=%v", res, strings.ReplaceAll(v.(string), "+", "%2B"))
	}
	if v, ok := d.GetOk("job_id"); ok {
		res = fmt.Sprintf("%s&job_id=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}

func flattenGetAspCollectionResultsAspResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("asp", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"job_id":       utils.PathSearch("job_id", v, nil),
			"file_size":    utils.PathSearch("file_size", v, nil),
			"file_path":    utils.PathSearch("file_path", v, nil),
			"file_name":    utils.PathSearch("file_name", v, nil),
			"start_time":   utils.PathSearch("start_time", v, nil),
			"end_time":     utils.PathSearch("end_time", v, nil),
			"download_url": utils.PathSearch("download_url", v, nil),
			"status":       utils.PathSearch("status", v, nil),
			"obs_bucket":   flattenGetAspCollectionResultsAspObsBucketResponseBody(v),
		})
	}
	return res
}

func flattenGetAspCollectionResultsAspObsBucketResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("obs_bucket", resp, nil)
	if curJson == nil {
		return nil
	}

	res := []interface{}{
		map[string]interface{}{
			"name":      utils.PathSearch("name", curJson, nil),
			"type":      utils.PathSearch("type", curJson, nil),
			"url":       utils.PathSearch("url", curJson, nil),
			"port":      utils.PathSearch("port", curJson, nil),
			"domain_id": utils.PathSearch("domain_id", curJson, nil),
		},
	}
	return res
}
