package cdn

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CDN GET /v1.0/cdn/configuration/tags
func DataSourceDomainTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDomainTagsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the domain tags are located.`,
			},

			// Required parameters.
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the domain to query tags.`,
			},

			// Attributes.
			"tags": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
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
				},
				Description: `The list of domain tags that matched filter parameters.`,
			},
		},
	}
}

func listTags(client *golangsdk.ServiceClient, resourceId string) ([]interface{}, error) {
	httpUrl := "v1.0/cdn/configuration/tags"
	listPath := client.Endpoint + httpUrl
	listPath += fmt.Sprintf("?resource_id=%s", resourceId)

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

func flattenTags(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}

	return result
}

func dataSourceDomainTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		resourceId = d.Get("resource_id").(string)
	)

	client, err := cfg.NewServiceClient("cdn", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	tags, err := listTags(client, resourceId)
	if err != nil {
		return diag.Errorf("error querying domain tags: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("tags", flattenTags(tags)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
