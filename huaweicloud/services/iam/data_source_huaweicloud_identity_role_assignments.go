package iam

import (
	"context"
	"fmt"
	"strconv"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const perPage = 50

// @API IAM GET /v3.0/OS-PERMISSION/role-assignments
func DataSourceIdentityRoleAssignments() *schema.Resource {
	return &schema.Resource{
		ReadContext: DataSourceIdentityRoleAssignmentsRead,

		Schema: map[string]*schema.Schema{
			"role_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject_user_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject_group_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"subject_agency_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"scope_enterprise_projects_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_inherited": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, false),
			},
			"include_group": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					"true", "false",
				}, true),
			},
			"role_assignments": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"user_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"role_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"group_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"agency_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"is_inherited": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"scope": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
					},
				},
			},
		},
	}
}

func DataSourceIdentityRoleAssignmentsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	domainId := cfg.DomainID
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	listRecordPath := iamClient.Endpoint + "v3.0/OS-PERMISSION/role-assignments"
	listRecordPath += fmt.Sprintf("?per_page=%v", perPage)
	path := buildQueryRoleAssignmentsPath(d, listRecordPath, domainId)
	options := golangsdk.RequestOpts{KeepResponseBody: true}
	currentPage := 1
	var res []interface{}
	for {
		getPath := path + fmt.Sprintf("&page=%d", currentPage)
		response, err := iamClient.Request("GET", getPath, &options)
		if err != nil {
			return diag.Errorf("error listRoleAssignments: %s", err)
		}
		respBody, err := utils.FlattenResponse(response)
		if err != nil {
			return diag.FromErr(err)
		}
		result := flattenRoleAssignments(respBody)
		res = append(res, result...)
		if len(result) == 0 {
			break
		}
		currentPage++
	}
	id, _ := uuid.GenerateUUID()
	d.SetId(id)
	if err = d.Set("role_assignments", res); err != nil {
		return diag.Errorf("error setting role assignmnets: %s", err)
	}
	return nil
}

func flattenRoleAssignments(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("role_assignments", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	results := make([]interface{}, 0, len(curArray))
	for _, assignment := range curArray {
		results = append(results, map[string]interface{}{
			"role_id":      utils.PathSearch("role.id", assignment, nil),
			"is_inherited": utils.PathSearch("is_inherited", assignment, nil),
			"user_id":      utils.PathSearch("user.id", assignment, nil),
			"group_id":     utils.PathSearch("group.id", assignment, nil),
			"agency_id":    utils.PathSearch("agency.id", assignment, nil),
			"scope":        flattenAssignmentsScope(utils.PathSearch("scope", assignment, nil)),
		})
	}
	return results
}

func flattenAssignmentsScope(scopeBody interface{}) map[string]interface{} {
	if scopeBody == nil {
		return nil
	}
	scope := map[string]interface{}{
		"project_id":            utils.PathSearch("project.id", scopeBody, nil),
		"domain_id":             utils.PathSearch("domain.id", scopeBody, nil),
		"enterprise_project_id": utils.PathSearch("enterprise_project.id", scopeBody, nil),
	}
	return scope
}

func buildQueryRoleAssignmentsPath(d *schema.ResourceData, getPath string, domainId string) string {
	getPath += fmt.Sprintf("&domain_id=%s", domainId)
	if roleId := d.Get("role_id").(string); roleId != "" {
		getPath += fmt.Sprintf("&role_id=%s", roleId)
	}
	subject := d.Get("subject").(string)
	if subject != "" {
		getPath += fmt.Sprintf("&subject=%s", subject)
	}
	subjectUserId := d.Get("subject_user_id").(string)
	if subjectUserId != "" {
		getPath += fmt.Sprintf("&subject.user_id=%s", subjectUserId)
	}
	if subjectGroupId := d.Get("subject_group_id").(string); subjectGroupId != "" {
		getPath += fmt.Sprintf("&subject.group_id=%s", subjectGroupId)
	}
	if subjectAgencyId := d.Get("subject_agency_id").(string); subjectAgencyId != "" {
		getPath += fmt.Sprintf("&subject.agency_id=%s", subjectAgencyId)
	}
	scope := d.Get("scope").(string)
	if scope != "" {
		getPath += fmt.Sprintf("&scope=%s", scope)
	}
	if scopeProjectId := d.Get("scope_project_id").(string); scopeProjectId != "" {
		getPath += fmt.Sprintf("&scope.project_id=%s", scopeProjectId)
	}
	scopeDomainId := d.Get("scope_domain_id").(string)
	if scopeDomainId != "" {
		getPath += fmt.Sprintf("&scope.domain_id=%s", scopeDomainId)
	}
	if scopeEnterpriseProjectsId := d.Get("scope_enterprise_projects_id").(string); scopeEnterpriseProjectsId != "" {
		getPath += fmt.Sprintf("&scope.enterprise_projects_id=%s", scopeEnterpriseProjectsId)
	}
	if v, ok := d.GetOk("is_inherited"); ok {
		isInherited, _ := strconv.ParseBool(v.(string))
		getPath += fmt.Sprintf("&is_inherited=%v", isInherited)
	}
	if v, ok := d.GetOk("include_group"); ok {
		includeGroup, _ := strconv.ParseBool(v.(string))
		getPath += fmt.Sprintf("&include_group=%v", includeGroup)
	}
	return getPath
}
