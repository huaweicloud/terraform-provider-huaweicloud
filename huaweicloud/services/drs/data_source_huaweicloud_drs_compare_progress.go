package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/compare-progress/{compare_job_id}
func DataSourceDrsCompareProgress() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsCompareProgressRead,

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
			"compare_job_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"full_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"progress": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"src_speed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recheck_entities": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"incre_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"delay": {
							Type:     schema.TypeFloat,
							Computed: true,
						},
						"src_speed": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rps": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"log_point": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"recheck_entities": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
			"global_info": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"src_speed": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDrsCompareProgressRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		product      = "drs"
		httpUrl      = "v5/{project_id}/jobs/{job_id}/compare-progress/{compare_job_id}"
		jobId        = d.Get("job_id").(string)
		compareJobId = d.Get("compare_job_id").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{job_id}", jobId)
	requestPath = strings.ReplaceAll(requestPath, "{compare_job_id}", compareJobId)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &reqOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS compare progress: %s", err)
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

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("full_info", flattenCompareProgressFullInfo(utils.PathSearch("full_info", respBody, nil))),
		d.Set("incre_info", flattenCompareProgressIncreInfo(utils.PathSearch("incre_info", respBody, nil))),
		d.Set("global_info", flattenCompareProgressGlobalInfo(utils.PathSearch("global_info", respBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenCompareProgressFullInfo(fullInfo interface{}) []interface{} {
	if fullInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"progress":         utils.PathSearch("progress", fullInfo, nil),
			"src_speed":        utils.PathSearch("src_speed", fullInfo, nil),
			"recheck_entities": utils.PathSearch("recheck_entities", fullInfo, nil),
		},
	}
}

func flattenCompareProgressIncreInfo(increInfo interface{}) []interface{} {
	if increInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"delay":            utils.PathSearch("delay", increInfo, nil),
			"src_speed":        utils.PathSearch("src_speed", increInfo, nil),
			"rps":              utils.PathSearch("rps", increInfo, nil),
			"log_point":        utils.PathSearch("log_point", increInfo, nil),
			"recheck_entities": utils.PathSearch("recheck_entities", increInfo, nil),
		},
	}
}

func flattenCompareProgressGlobalInfo(globalInfo interface{}) []interface{} {
	if globalInfo == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"src_speed": utils.PathSearch("src_speed", globalInfo, nil),
		},
	}
}
