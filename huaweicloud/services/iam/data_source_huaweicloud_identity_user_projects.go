package iam

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// DataSourceIdentityUserProjects
// @API IAM GET /v3/users/{user_id}/projects
// @API IAM GET /v5/caller-identity
func DataSourceIdentityUserProjects() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityUserProjectsRead,

		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The IAM user id. The default value is the current user id.",
			},

			"projects": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"enabled": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityUserProjectsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	userProjectsPath := iamClient.Endpoint + "v3/users/{user_id}/projects"
	userId := d.Get("user_id").(string)
	if userId == "" {
		userId, err = getUserId(cfg)
		if err != nil {
			return diag.Errorf("error getUserId: %s", err)
		}
	}
	userProjectsPath = strings.ReplaceAll(userProjectsPath, "{user_id}", userId)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := iamClient.Request("GET", userProjectsPath, &options)
	if err != nil {
		return diag.Errorf("error listProjectsForUser: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)
	projectsBody := utils.PathSearch("projects", respBody, make([]interface{}, 0)).([]interface{})
	projects := flattenProjects(projectsBody)
	if err = d.Set("projects", projects); err != nil {
		return diag.Errorf("error setting projects fields: %s", err)
	}
	return nil
}

func getUserId(cfg *config.Config) (string, error) {
	if cfg.UserID != "" {
		return cfg.UserID, nil
	}
	callerIdentity, err := queryCallerIdentity(cfg)
	if err != nil {
		return "", err
	}
	return utils.PathSearch("principal_id", callerIdentity, "").(string), nil
}

func queryCallerIdentity(cfg *config.Config) (interface{}, error) {
	client, err := cfg.StsClient(cfg.Region)
	if err != nil {
		return nil, err
	}
	callerIdentityPath := client.Endpoint + "v5/caller-identity"
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	response, err := client.Request("GET", callerIdentityPath, &options)
	if err != nil {
		return nil, err
	}
	callerIdentity, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, err
	}
	return callerIdentity, nil
}
