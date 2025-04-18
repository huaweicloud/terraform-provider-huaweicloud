package servicestage

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

// @API ServiceStage GET /v3/{project_id}/cas/runtimestacks
func DataSourceV3RuntimeStacks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3RuntimeStacksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the runtime stacks are located.`,
			},
			"runtime_stacks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the runtime stack.`,
						},
						"deploy_mode": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The deploy mode.`,
						},
						"version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The version number.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The URL of the runtime stack.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the runtime stack.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The status of the runtime stack.`,
						},
					},
				},
				Description: "All runtime stack details.",
			},
		},
	}
}

func listV3RuntimeStacks(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/cas/runtimestacks"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("runtime_stacks", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenV3RuntimeStacks(runtimeStacks []interface{}) []map[string]interface{} {
	if len(runtimeStacks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(runtimeStacks))
	for _, runtimeStack := range runtimeStacks {
		result = append(result, map[string]interface{}{
			"name":        utils.PathSearch("name", runtimeStack, nil),
			"deploy_mode": utils.PathSearch("deploy_mode", runtimeStack, nil),
			"version":     utils.PathSearch("version", runtimeStack, nil),
			"url":         utils.PathSearch("url", runtimeStack, nil),
			"type":        utils.PathSearch("type", runtimeStack, nil),
			"status":      utils.PathSearch("status", runtimeStack, nil),
		})
	}
	return result
}

func dataSourceV3RuntimeStacksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	runtimeStacks, err := listV3RuntimeStacks(client)
	if err != nil {
		return diag.Errorf("error getting runtime stacks: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("runtime_stacks", flattenV3RuntimeStacks(runtimeStacks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
