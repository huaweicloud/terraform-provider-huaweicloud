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

// @API DRS GET /v5/{project_id}/support-links
func DataSourceSupportLinks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSupportLinksRead,

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
			"support_links": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"engine_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"net_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"task_modes": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"job_direction": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_mode": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"job_instance_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_support_bind_eip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildSupportLinksQueryParams(d *schema.ResourceData) string {
	return fmt.Sprintf("?job_type=%s", d.Get("job_type").(string))
}

func dataSourceSupportLinksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "drs"
		httpUrl = "v5/{project_id}/support-links"
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
		queryParams := buildSupportLinksQueryParams(d)
		currentListPath := fmt.Sprintf("%s%s&limit=%d&offset=%d", listPath, queryParams, limit, offset)
		listResp, err := client.Request("GET", currentListPath, &reqOpt)
		if err != nil {
			return diag.FromErr(err)
		}
		listRespBody, err := utils.FlattenResponse(listResp)
		if err != nil {
			return diag.FromErr(err)
		}

		supportLinks := utils.PathSearch("support_links", listRespBody, make([]interface{}, 0)).([]interface{})
		result = append(result, supportLinks...)

		if len(supportLinks) < limit {
			break
		}

		offset += len(supportLinks)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("support_links", flattenSupportLinks(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSupportLinks(supportLinks []interface{}) []interface{} {
	if len(supportLinks) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(supportLinks))
	for _, supportLink := range supportLinks {
		result = append(result, map[string]interface{}{
			"engine_type":         utils.PathSearch("engine_type", supportLink, nil),
			"net_type":            utils.PathSearch("net_type", supportLink, nil),
			"task_modes":          utils.PathSearch("task_modes", supportLink, make([]interface{}, 0)),
			"job_direction":       utils.PathSearch("job_direction", supportLink, nil),
			"cluster_mode":        utils.PathSearch("cluster_mode", supportLink, nil),
			"job_instance_type":   utils.PathSearch("job_instance_type", supportLink, nil),
			"is_support_bind_eip": utils.PathSearch("is_support_bind_eip", supportLink, false),
		})
	}

	return result
}
