package workspace

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

// @API Workspace GET /v2/{project_id}/desktops/{desktop_id}/tags
func DataSourceDesktopTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDesktopTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the Workspace desktop is located.`,
			},
			"desktop_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the desktop to query tags.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        desktopTag(),
				Description: `The list of tags.`,
			},
		},
	}
}

func desktopTag() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The key of the tag.`,
			},
			"value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The value of the tag.`,
			},
		},
	}
}

func listDesktopTags(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl   = "v2/{project_id}/desktops/{desktop_id}/tags"
		desktopId = d.Get("desktop_id").(string)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{desktop_id}", desktopId)

	requestOpts := &golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, requestOpts)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	tags := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
	return tags, nil
}

func flattenDesktopTags(tags []interface{}) []interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, tag := range tags {
		tagMap := tag.(map[string]interface{})
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tagMap, nil),
			"value": utils.PathSearch("value", tagMap, nil),
		})
	}
	return result
}

func dataSourceDesktopTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating workspace client: %s", err)
	}

	tags, err := listDesktopTags(client, d)
	if err != nil {
		return diag.Errorf("error querying workspace desktop tags: %s", err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenDesktopTags(tags)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
