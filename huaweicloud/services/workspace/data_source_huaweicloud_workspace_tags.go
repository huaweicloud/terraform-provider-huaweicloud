package workspace

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

// @API Workspace GET /v2/{project_id}/desktops/tags
func DataSourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The region in which to obtain the desktop tags.",
			},
			"key": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The key of the tag to be queried.",
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        tag(),
				Description: "The list of workspace tags.",
			},
		},
	}
}

func tag() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The key of the tag.",
			},
			"values": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The values of the tag.",
			},
		},
	}
}

func buildListTagsParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("key"); ok {
		res = fmt.Sprintf("%s&key=%v", res, v)
	}
	return res
}

func listTags(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/desktops/tags"
		offset  = 0
		limit   = 200
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s?limit=%d&offset=%d", listPath, limit, offset)
		listPathWithOffset += buildListTagsParams(d)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		tags := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, tags...)
		if len(tags) < limit {
			break
		}
		offset += len(tags)
	}
	return result, nil
}

func flattenTags(tags []interface{}) []interface{} {
	if len(tags) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(tags))
	for _, item := range tags {
		tag := map[string]interface{}{
			"key":    utils.PathSearch("key", item, nil),
			"values": utils.PathSearch("values", item, nil),
		}
		result = append(result, tag)
	}

	return result
}

func dataSourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating WorkSpace client: %s", err)
	}

	tags, err := listTags(client, d)
	if err != nil {
		return diag.Errorf("error querying WorkSpace tags: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("tags", flattenTags(tags)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
