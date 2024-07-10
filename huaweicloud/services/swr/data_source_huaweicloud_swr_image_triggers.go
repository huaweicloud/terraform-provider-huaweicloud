package swr

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/triggers
func DataSourceImageTriggers() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceImageTriggersRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"organization": {
				Type:     schema.TypeString,
				Required: true,
			},
			"repository": {
				Type:     schema.TypeString,
				Required: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"condition_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"triggers": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workload_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"workload_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"condition_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"condition_value": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"container": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_at": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"created_by": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"histories": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"tag": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"result": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"detail": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func dataSourceImageTriggersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		listImageTriggersHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/triggers"
		listImageTriggersProduct = "swr"
	)

	listImageTriggersClient, err := cfg.NewServiceClient(listImageTriggersProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	listImageTriggersPath := listImageTriggersClient.Endpoint + listImageTriggersHttpUrl
	listImageTriggersPath = strings.ReplaceAll(listImageTriggersPath, "{namespace}", organization)
	listImageTriggersPath = strings.ReplaceAll(listImageTriggersPath, "{repository}", repository)

	listImageTriggersOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listImageTriggersResp, err := listImageTriggersClient.Request("GET", listImageTriggersPath, &listImageTriggersOpt)
	if err != nil {
		return diag.Errorf("error querying SWR image triggers: %v", err)
	}

	listImageTriggersRespBody, err := utils.FlattenResponse(listImageTriggersResp)
	if err != nil {
		return diag.Errorf("error retrieving SWR image triggers: %s", err)
	}

	results := make([]map[string]interface{}, 0)

	imageTriggers := listImageTriggersRespBody.([]interface{})
	for _, imageTrigger := range imageTriggers {
		name := utils.PathSearch("name", imageTrigger, "").(string)
		enabled := utils.PathSearch("enable", imageTrigger, "").(string)
		conditionType := utils.PathSearch("trigger_type", imageTrigger, "").(string)
		clusterName := utils.PathSearch("cluster_name", imageTrigger, "").(string)
		if val, ok := d.GetOk("name"); ok && name != val {
			continue
		}
		if val, ok := d.GetOk("enabled"); ok && enabled != val {
			continue
		}
		if val, ok := d.GetOk("condition_type"); ok && conditionType != val {
			continue
		}
		if val, ok := d.GetOk("cluster_name"); ok && clusterName != val {
			continue
		}

		results = append(results, map[string]interface{}{
			"action":          utils.PathSearch("action", imageTrigger, nil),
			"workload_name":   utils.PathSearch("application", imageTrigger, nil),
			"workload_type":   utils.PathSearch("app_type", imageTrigger, nil),
			"cluster_id":      utils.PathSearch("cluster_id", imageTrigger, nil),
			"cluster_name":    clusterName,
			"namespace":       utils.PathSearch("cluster_ns", imageTrigger, nil),
			"name":            name,
			"type":            utils.PathSearch("trigger_mode", imageTrigger, nil),
			"condition_type":  conditionType,
			"condition_value": utils.PathSearch("condition", imageTrigger, nil),
			"container":       utils.PathSearch("container", imageTrigger, nil),
			"enabled":         enabled,
			"created_at":      utils.PathSearch("created_at", imageTrigger, nil),
			"created_by":      utils.PathSearch("creator_name", imageTrigger, nil),
			"histories":       flattenImageTriggerHistories(imageTrigger),
		})
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("triggers", results),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenImageTriggerHistories(imageTrigger interface{}) []interface{} {
	histories := utils.PathSearch("trigger_history", imageTrigger, make([]interface{}, 0)).([]interface{})
	rst := make([]interface{}, 0, len(histories))
	for _, history := range histories {
		rst = append(rst, map[string]interface{}{
			"tag":    utils.PathSearch("tag", history, nil),
			"result": utils.PathSearch("result", history, nil),
			"detail": utils.PathSearch("detail", history, nil),
		})
	}
	return rst
}
