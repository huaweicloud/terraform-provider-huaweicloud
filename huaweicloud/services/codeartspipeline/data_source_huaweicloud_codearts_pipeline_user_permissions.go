package codeartspipeline

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

// @API CodeArtsPipeline GET /v5/{project_id}/api/pipeline-permissions/{pipeline_id}/user-permission
func DataSourceCodeArtsPipelineUserPermissions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCodeArtsPipelineUserPermissionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the CodeArts project ID.`,
			},
			"pipeline_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the pipeline ID.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the user name.`,
			},
			"users": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `Indicates the template list.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user ID.`,
						},
						"user_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the user name.`,
						},
						"operation_query": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has the permission to query.`,
						},
						"operation_execute": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has the permission to execute.`,
						},
						"operation_update": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has the permission to update.`,
						},
						"operation_delete": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has the permission to delete.`,
						},
						"operation_authorize": {
							Type:        schema.TypeBool,
							Computed:    true,
							Description: `Indicates whether the user has the permission to authorize.`,
						},
						"role_id": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `Indicates the role ID.`,
						},
						"role_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Indicates the role name.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceCodeArtsPipelineUserPermissionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("codearts_pipeline", region)
	if err != nil {
		return diag.Errorf("error creating CodeArts Pipeline client: %s", err)
	}

	getHttpUrl := "v5/{project_id}/api/pipeline-permissions/{pipeline_id}/user-permission?limit=10"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", d.Get("project_id").(string))
	getPath = strings.ReplaceAll(getPath, "{pipeline_id}", d.Get("pipeline_id").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	// have to send
	getPath += fmt.Sprintf("&subject=%v", d.Get("user_name"))

	offset := 0
	rst := make([]map[string]interface{}, 0)
	for {
		currentPath := getPath + fmt.Sprintf("&offset=%d", offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving pipeline users: %s", err)
		}
		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.Errorf("error flatten response: %s", err)
		}

		if err := checkResponseError(getRespBody, ""); err != nil {
			return diag.Errorf("error retrieving pipeline users: %s", err)
		}

		users := utils.PathSearch("users", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(users) == 0 {
			break
		}

		for _, user := range users {
			rst = append(rst, map[string]interface{}{
				"user_id":             utils.PathSearch("user_id", user, nil),
				"user_name":           utils.PathSearch("user_name", user, nil),
				"role_id":             utils.PathSearch("role_id", user, nil),
				"role_name":           utils.PathSearch("role_name", user, nil),
				"operation_query":     utils.PathSearch("operation_query", user, nil),
				"operation_execute":   utils.PathSearch("operation_execute", user, nil),
				"operation_update":    utils.PathSearch("operation_update", user, nil),
				"operation_delete":    utils.PathSearch("operation_delete", user, nil),
				"operation_authorize": utils.PathSearch("operation_authorize", user, nil),
			})
		}

		offset += 10
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("users", rst),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
