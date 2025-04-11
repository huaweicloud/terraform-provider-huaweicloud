package cbr

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

// @API CBR GET /v3/{project_id}/vault/tags
func DataSourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region in which to query the resource tags.`,
			},
			"resource_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource type to which the tags belong that to be queried.`,
			},
			"tags": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of all tags for resources of the same type.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key of the resource tag.`,
						},
						"values": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `All values corresponding to the key.`,
						},
					},
				},
			},
		},
	}
}

func listProjectTags(client *golangsdk.ServiceClient, resourceType string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/{resource_type}/tags"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", resourceType)

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
	return utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func flattenProjectTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":    utils.PathSearch("key", tag, nil),
			"values": utils.PathSearch("values", tag, make([]interface{}, 0)),
		})
	}

	return result
}

func dataSourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		resourceType = d.Get("resource_type").(string)
	)
	client, err := cfg.NewServiceClient("cbr", region)
	if err != nil {
		return diag.Errorf("error creating CBR client: %s", err)
	}

	tagList, err := listProjectTags(client, resourceType)
	if err != nil {
		return diag.Errorf("error querying tags of the CBR service: %s", err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", flattenProjectTags(tagList)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the CBR project tags: %s", mErr)
	}
	return nil
}
