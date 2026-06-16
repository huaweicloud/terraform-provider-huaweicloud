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

// @API DRS GET /v5/{project_id}/links
func DataSourceJobLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceJobLinksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"job_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"job_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"job_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_types": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"job_direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"source_endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_endpoint_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildJobLinksQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?job_type=%s", d.Get("job_type").(string))
	return res
}

func dataSourceJobLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/links"
		result  = make([]interface{}, 0)
		limit   = 1000
		offset  = 0
		mErr    *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		queryParams := buildJobLinksQueryParams(d)
		currentListPath := fmt.Sprintf("%s%s&limit=%d&offset=%d", listPath, queryParams, limit, offset)
		listResp, err := client.Request("GET", currentListPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		jobLinks := utils.PathSearch("job_links", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, jobLinks...)
		if len(jobLinks) < limit {
			break
		}

		offset += len(jobLinks)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("job_links", flattenJobLinks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenJobLinks(jobLinks []interface{}) []interface{} {
	if len(jobLinks) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(jobLinks))
	for _, jobLink := range jobLinks {
		result = append(result, map[string]interface{}{
			"job_type":             utils.PathSearch("job_type", jobLink, nil),
			"engine_type":          utils.PathSearch("engine_type", jobLink, nil),
			"net_type":             utils.PathSearch("net_type", jobLink, nil),
			"task_types":           utils.PathSearch("task_types", jobLink, make([]interface{}, 0)),
			"job_direction":        utils.PathSearch("job_direction", jobLink, nil),
			"cluster_modes":        utils.PathSearch("cluster_modes", jobLink, make([]interface{}, 0)),
			"source_endpoint_type": utils.PathSearch("source_endpoint_type", jobLink, nil),
			"target_endpoint_type": utils.PathSearch("target_endpoint_type", jobLink, nil),
		})
	}

	return result
}
