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

// @API DataArtsStudio GET /v2/{project_id}/design/directorys
func DataSourceArchitectureDirectories() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceArchitectureDirectoriesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the directories are located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the directories belong.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of directory to be queried.`,
			},

			// Attributes.
			"directories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataArchitectureDirectory(),
				Description: `The list of directories that matched filter parameters.`,
			},
		},
	}
}

func dataArchitectureDirectory() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the directory.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the directory.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the directory.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The description of the directory.`,
			},
			"parent_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The parent ID of the directory.`,
			},
			"root_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The root directory ID.`,
			},
			"qualified_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The qualified name of the directory.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the directory, in RFC3339 format.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The latest update time of the directory, in RFC3339 format.`,
			},
			"created_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who created the directory.`,
			},
			"updated_by": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who updated the directory.`,
			},
			"children": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of sub-directories.`,
			},
		},
	}
}

func listArchitectureDirectories(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		//nolint:misspell
		httpUrl = "v2/{project_id}/design/directorys?type={type}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{type}", d.Get("type").(string))
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildArchitectureMoreHeaders(d.Get("workspace_id").(string)),
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		directories := utils.PathSearch("data.value", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, directories...)

		if len(directories) < limit {
			break
		}
		offset += len(directories)
	}

	return result, nil
}

func flattenArchitectureDirectoryNode(directory interface{}) []map[string]interface{} {
	result := []map[string]interface{}{
		{
			"id":             utils.PathSearch("id", directory, nil),
			"name":           utils.PathSearch("name", directory, nil),
			"type":           utils.PathSearch("type", directory, nil),
			"description":    utils.PathSearch("description", directory, nil),
			"parent_id":      utils.PathSearch("parent_id", directory, nil),
			"root_id":        utils.PathSearch("root_id", directory, nil),
			"qualified_name": utils.PathSearch("qualified_name", directory, nil),
			"created_at":     utils.PathSearch("create_time", directory, nil),
			"updated_at":     utils.PathSearch("update_time", directory, nil),
			"created_by":     utils.PathSearch("create_by", directory, nil),
			"updated_by":     utils.PathSearch("update_by", directory, nil),
			"children":       utils.PathSearch(`children[*].name`, directory, make([]interface{}, 0)),
		},
	}

	children := utils.PathSearch("children", directory, make([]interface{}, 0)).([]interface{})
	for _, child := range children {
		result = append(result, flattenArchitectureDirectoryNode(child)...)
	}

	return result
}

func flattenArchitectureDirectories(directories []interface{}) []map[string]interface{} {
	if len(directories) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0)
	for _, directory := range directories {
		result = append(result, flattenArchitectureDirectoryNode(directory)...)
	}
	return result
}

func dataSourceArchitectureDirectoriesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts client: %s", err)
	}

	directories, err := listArchitectureDirectories(client, d)
	if err != nil {
		return diag.Errorf("error querying DataArts Architecture directories: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("directories", flattenArchitectureDirectories(directories)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
