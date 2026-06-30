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

// @API DRS POST /v3/{project_id}/jobs/batch-rpo-and-rto
func DataSourceDrsBatchRposAndRtos() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsBatchRposAndRtosRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"jobs": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"results": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"rpo_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     rpoRtoInfoSchema(),
						},
						"rto_info": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     rpoRtoInfoSchema(),
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_msg": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func rpoRtoInfoSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"check_point": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"delay": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"gtid_set": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildBatchRposAndRtosBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"jobs": utils.ExpandToStringList(d.Get("jobs").([]interface{})),
	}
}

func dataSourceDrsBatchRposAndRtosRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v3/{project_id}/jobs/batch-rpo-and-rto"
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
		JSONBody: buildBatchRposAndRtosBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS batch RPOs and RTOs: %s", err)
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
		d.Set("region", region),
		d.Set("results", flattenBatchRposAndRtosResults(respBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchRposAndRtosResults(respBody interface{}) []interface{} {
	resultsRaw := utils.PathSearch("results", respBody, make([]interface{}, 0)).([]interface{})
	if len(resultsRaw) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resultsRaw))
	for _, item := range resultsRaw {
		result = append(result, map[string]interface{}{
			"job_id":     utils.PathSearch("job_id", item, nil),
			"rpo_info":   flattenRpoRtoInfo(utils.PathSearch("rpo_info", item, nil)),
			"rto_info":   flattenRpoRtoInfo(utils.PathSearch("rto_info", item, nil)),
			"error_code": utils.PathSearch("error_code", item, nil),
			"error_msg":  utils.PathSearch("error_msg", item, nil),
		})
	}

	return result
}

func flattenRpoRtoInfo(info interface{}) []interface{} {
	if info == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"check_point": utils.PathSearch("check_point", info, nil),
			"delay":       utils.PathSearch("delay", info, nil),
			"gtid_set":    utils.PathSearch("gtid_set", info, nil),
			"time":        utils.PathSearch("time", info, nil),
		},
	}
}
