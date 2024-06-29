package lts

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

// @API LTS GET /v2/{project_id}/lts/aom-mapping
func DataSourceAOMAccesses() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceAOMAccessesRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"log_group_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the ID of the log group name to be queried.`,
			},
			"log_stream_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the log stream name to be queried.`,
			},
			"accesses": {
				Type:        schema.TypeList,
				Elem:        AOMAccesseschema(),
				Computed:    true,
				Description: `All AOM access rules that match the filter parameters.`,
			},
		},
	}
}

func AOMAccesseschema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the AOM access rule.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the AOM access rule.`,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster ID corresponding to the AOM access rule.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster name corresponding to the AOM access rule.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The namespace corresponding to the AOM access rule.`,
			},
			"workloads": {
				Type:        schema.TypeList,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Computed:    true,
				Description: `The list of the workloads corresponding to AOM access rule.`,
			},
			"container_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the container corresponding to AOM access rule.`,
			},
			"access_rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The AOM access log details.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"file_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log path.`,
						},
						"log_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log group.`,
						},
						"log_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the log group.`,
						},
						"log_stream_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the log stream.`,
						},
						"log_stream_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the stream.`,
						},
					},
				},
			},
		},
	}
}

func resourceAOMAccessesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr    *multierror.Error
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/lts/aom-mapping"
		product = "lts"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating LTS client: %s", err)
	}

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath += buildListAOMAccessesQueryParams(d)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return diag.Errorf("error retrieving AOM accesses: %s", err)
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("accesses", flattenAndFilterAOMAccesses(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenAndFilterAOMAccesses(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curArray := resp.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"id":             utils.PathSearch("rule_id", v, nil),
			"name":           utils.PathSearch("rule_name", v, nil),
			"cluster_id":     utils.PathSearch("rule_info.cluster_id", v, nil),
			"cluster_name":   utils.PathSearch("rule_info.cluster_name", v, nil),
			"namespace":      utils.PathSearch("rule_info.namespace", v, nil),
			"container_name": utils.PathSearch("rule_info.container_name", v, nil),
			"workloads":      utils.PathSearch("rule_info.deployments", v, nil),
			"access_rules":   flattenAccessRules(v),
		}
	}
	return rst
}

func buildListAOMAccessesQueryParams(d *schema.ResourceData) string {
	queryParam := ""
	if v, ok := d.GetOk("log_group_name"); ok {
		queryParam = fmt.Sprintf("%s&log_group_name=%v", queryParam, v)
	}

	if v, ok := d.GetOk("log_stream_name"); ok {
		queryParam = fmt.Sprintf("%s&log_stream_name=%v", queryParam, v)
	}

	if queryParam != "" {
		queryParam = "?" + queryParam[1:]
	}
	return queryParam
}
