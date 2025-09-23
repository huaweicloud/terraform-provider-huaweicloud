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

// @API ServiceStage GET /v3/{project_id}/cas/innerimages
func DataSourceV3InnerRuntimeStacks() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV3InnerRuntimeStacksRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the inner runtime stacks are located.`,
			},
			"runtime_stacks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the inner runtime stack.`,
						},
						"url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The image URL of the inner runtime stack.`,
						},
					},
				},
				Description: "All inner runtime stack details.",
			},
		},
	}
}

func listV3InnerRuntimeStacks(client *golangsdk.ServiceClient) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/cas/innerimages"

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

func flattenV3InnerRuntimeStacks(runtimeStacks []interface{}) []map[string]interface{} {
	if len(runtimeStacks) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(runtimeStacks))
	for _, runtimeStack := range runtimeStacks {
		result = append(result, map[string]interface{}{
			"type": utils.PathSearch("type", runtimeStack, nil),
			"url":  utils.PathSearch("url", runtimeStack, nil),
		})
	}
	return result
}

func dataSourceV3InnerRuntimeStacksRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	innerRuntimeStacks, err := listV3InnerRuntimeStacks(client)
	if err != nil {
		return diag.Errorf("error getting inner runtime stacks: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("runtime_stacks", flattenV3InnerRuntimeStacks(innerRuntimeStacks)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
