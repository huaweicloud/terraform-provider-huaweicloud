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

// @API IAM GET /v3.0/OS-PERMISSION/groups/{group_id}/enterprise-projects
// @API IAM GET /v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/groups
func DataSourceIdentityEnterpriseProjectGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityEnterpriseProjectGroupsRead,
		Schema: map[string]*schema.Schema{
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"groups": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"create_time": {
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

func DataSourceIdentityEnterpriseProjectGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IdentityV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	enterpriseProjectId := d.Get("enterprise_project_id").(string)
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/enterprise-projects/{enterprise_project_id}/groups"
	path = strings.ReplaceAll(path, "{enterprise_project_id}", enterpriseProjectId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return diag.Errorf("error getting identity enterprise-projects groups: %s", err)
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
	groupsList := flattenEnterpriseProjectGroupsList(utils.PathSearch("groups", respBody, make([]interface{}, 0)).([]interface{}), iamClient)
	if e := d.Set("groups", groupsList); e != nil {
		return diag.FromErr(e)
	}
	return nil
}

func flattenEnterpriseProjectGroupsList(list []interface{}, iamClient *golangsdk.ServiceClient) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, group := range list {
		groupID := utils.PathSearch("id", group, "").(string)
		projects, err := GetGroupEnterpriseProjectsRead(iamClient, groupID)
		if err != nil {
			log.Printf("[DEBUG] Failed to get enterprise projects for group ID %s: %s", groupID, err)
			return nil
		}
		res[i] = map[string]interface{}{
			"create_time":         utils.PathSearch("createTime", group, nil),
			"description":         utils.PathSearch("description", group, ""),
			"domain_id":           utils.PathSearch("domain_id", group, ""),
			"name":                utils.PathSearch("name", group, ""),
			"id":                  groupID,
			"enterprise_projects": projects,
		}
	}
	return res
}

func GetGroupEnterpriseProjectsRead(iamClient *golangsdk.ServiceClient, groupId string) ([]map[string]interface{}, error) {
	path := iamClient.Endpoint + "v3.0/OS-PERMISSION/groups/{group_id}/enterprise-projects"
	path = strings.ReplaceAll(path, "{group_id}", groupId)
	identityOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	response, err := iamClient.Request("GET", path, &identityOpt)
	if err != nil {
		return nil, fmt.Errorf("error getting identity enterprise-projects for group %s: %w", groupId, err)
	}
	respBody, err := utils.FlattenResponse(response)
	if err != nil {
		return nil, fmt.Errorf("error flattening response for group %s: %w", groupId, err)
	}
	enterpriseProjectList := flattenEnterpriseProjectList(utils.PathSearch("\"enterprise-projects\"",
		respBody, make([]interface{}, 0)).([]interface{}))
	return enterpriseProjectList, nil
}

func flattenEnterpriseProjectList(list []interface{}) []map[string]interface{} {
	res := make([]map[string]interface{}, len(list))
	for i, project := range list {
		res[i] = map[string]interface{}{
			"project_id": utils.PathSearch("projectId", project, ""),
		}
	}
	return res
}
