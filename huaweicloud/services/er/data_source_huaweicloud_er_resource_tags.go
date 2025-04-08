package er

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

// @API ER GET /v3/{project_id}/{resource_type}/{resource_id}/tags
func DataSourceResourceTags() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceResourceTagsRead,

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
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource ID to which the tags belong that to be queried.`,
			},
			"tags": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The tags of a specified resource.`,
			},
		},
	}
}

func listResourceTags(client *golangsdk.ServiceClient, resourceType, resourceId string) ([]interface{}, error) {
	httpUrl := "v3/{project_id}/{resource_type}/{resource_id}/tags"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{resource_type}", resourceType)
	listPath = strings.ReplaceAll(listPath, "{resource_id}", resourceId)

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

func dataSourceResourceTagsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg          = meta.(*config.Config)
		region       = cfg.GetRegion(d)
		resourceType = d.Get("resource_type").(string)
		resourceId   = d.Get("resource_id").(string)
	)
	client, err := cfg.NewServiceClient("er", region)
	if err != nil {
		return diag.Errorf("error creating ER client: %s", err)
	}

	tagList, err := listResourceTags(client, resourceType, resourceId)
	if err != nil {
		return diag.Errorf("error querying tags for a specified resource (%s) of the ER service: %s", resourceId, err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("tags", utils.FlattenTagsToMap(tagList)),
	)

	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving data source fields of the ER resource tags: %s", mErr)
	}
	return nil
}
