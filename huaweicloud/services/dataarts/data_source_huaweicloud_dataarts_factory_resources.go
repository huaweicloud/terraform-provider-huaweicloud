package dataarts

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DataArtsStudio GET /v1/{project_id}/resources
func DataSourceFactoryResources() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceFactoryResourcesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the resources are located.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the resources belong.`,
			},

			// Attributes.
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the resource.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the resource.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The type of the resource.`,
						},
						"location": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The OBS path of the resource file.`,
						},
						"directory": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The directory where the resource is located.`,
						},
						"description": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The description of the resource.`,
						},
						"depend_files": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of dependent file paths.`,
						},
					},
				},
				Description: `The list of resources that match the filter parameters.`,
			},
		},
	}
}

func buildFactoryMoreHeaders(workspaceId string, isExclusive ...bool) map[string]string {
	result := map[string]string{
		"Content-Type": "application/json",
		"Dlm-Type":     "EXCLUSIVE",
	}

	// By default, the DLM type is set to 'EXCLUSIVE'. If the isExclusive parameter is set to 'false', the DLM type is
	// not set.
	if len(isExclusive) > 0 && !isExclusive[0] {
		delete(result, "Dlm-Type")
	}

	if workspaceId != "" {
		result["workspace"] = workspaceId
	}

	return result
}

func listFactoryResources(client *golangsdk.ServiceClient, workspaceId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/resources?limit={limit}"
		offset  = 0
		limit   = 100
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildFactoryMoreHeaders(workspaceId),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		resources := utils.PathSearch("resources", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, resources...)
		if len(resources) < limit {
			break
		}
		offset += len(resources)
	}

	return result, nil
}

func flattenFactoryResources(resources []interface{}) []map[string]interface{} {
	if len(resources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(resources))
	for _, resource := range resources {
		result = append(result, map[string]interface{}{
			"id":           utils.PathSearch("resourceId", resource, nil),
			"name":         utils.PathSearch("name", resource, nil),
			"type":         utils.PathSearch("type", resource, nil),
			"location":     utils.PathSearch("location", resource, nil),
			"directory":    utils.PathSearch("directory", resource, nil),
			"description":  utils.PathSearch("desc", resource, nil),
			"depend_files": utils.PathSearch("dependFiles", resource, make([]interface{}, 0)).([]interface{}),
		})
	}

	return result
}

func dataSourceFactoryResourcesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		workspaceId = d.Get("workspace_id").(string)
	)

	client, err := cfg.NewServiceClient("dataarts-dlf", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	resources, err := listFactoryResources(client, workspaceId)
	if err != nil {
		return diag.Errorf("error querying DataArts Factory resources: %s", err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("resources", flattenFactoryResources(resources)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
