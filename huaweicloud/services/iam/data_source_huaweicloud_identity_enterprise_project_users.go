package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-PERMISSION/users/{user_id}/enterprise-projects
// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/users
func DataSourceIdentityEnterpriseProjectUsers() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityEnterpriseProjectUsersRead,
		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"users": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"policy_num": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"description": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
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
						"lastest_policy_time": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"enterprise_projects": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"project_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityEnterpriseProjectUsersRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	enterpriseProjectId := d.Get("enterprise_project_id").(string)
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/users"
	path = strings.ReplaceAll(path, "{enterprise_project_id}", enterpriseProjectId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return diag.Errorf("error getting identity enterprise-projects users: %s", err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return diag.FromErr(err)
	}
	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id)
	usersList := flattenEnterpriseProjectUsersList(utils.PathSearch("users", respBody, make([]interface{}, 0)).([]interface{}), iamClient)
	if e := d.Set("users", usersList); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenEnterpriseProjectUsersList(list []interface{}, iamClient *golangsdk.ServiceClient) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, user := range list {
		userID := utils.PathSearch("id", user, "").(string)
		projects, err := GetUserEnterpriseProjectsRead(iamClient, userID)
		if err != nil {
			log.Printf("[DEBUG] Failed to get enterprise projects for user ID %s: %s", userID, err)
			return nil
		}
		res[i] = map[string]interface{}{
			"policy_num":          utils.PathSearch("policy_num", user, nil),
			"description":         utils.PathSearch("description", user, ""),
			"domain_id":           utils.PathSearch("domain_id", user, ""),
			"name":                utils.PathSearch("name", user, ""),
			"enabled":             utils.PathSearch("enabled", user, nil),
			"lastest_policy_time": utils.PathSearch("lastest_policy_time", user, nil),
			"id":                  userID,
			"enterprise_projects": projects,
		}
	}
	return res
}

func GetUserEnterpriseProjectsRead(iamClient *golangsdk.ServiceClient, userId string) ([]map[string]interface{}, error) {
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/users/{user_id}/enterprise-projects"
	path = strings.ReplaceAll(path, "{user_id}", userId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting identity enterprise-projects for user %s: %w", userId, err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error flattening response for user %s: %w", userId, err)
	}
	enterpriseProjectList := flattenUserEnterpriseProjectList(utils.PathSearch("\"enterprise-projects\"",
		respBody, make([]interface{}, 0)).([]interface{}))
	return enterpriseProjectList, nil
}

func flattenUserEnterpriseProjectList(list []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, project := range list {
		res[i] = map[string]interface{}{
			"project_id": utils.PathSearch("projectId", project, ""),
		}
	}
	return res
}
