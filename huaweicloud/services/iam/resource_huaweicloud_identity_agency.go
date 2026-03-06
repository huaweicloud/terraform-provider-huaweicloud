package iam

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3AgencyNonUpdatableParams = []string{"name"}

// @API IAM GET /v3/projects
// @API EPS GET /v1.0/enterprise-projects
// @API IAM GET /v3/roles
// @API IAM POST /v3.0/OS-AGENCY/agencies
// @API IAM PUT /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM PUT /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM PUT /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects
// @API IAM PUT /v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments
// @API IAM GET /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM GET /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles
// @API IAM GET /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles
// @API IAM PUT /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM DELETE /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects
// @API IAM DELETE /v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments
// @API IAM GET /v5/agencies/{agency_id}/attached-policies
// @API IAM POST /v5/policies/{policy_id}/detach-agency
// @API IAM DELETE /v3.0/OS-AGENCY/agencies/{agency_id}
func ResourceV3Agency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3AgencyCreate,
		ReadContext:   resourceV3AgencyRead,
		UpdateContext: resourceV3AgencyUpdate,
		DeleteContext: resourceV3AgencyDelete,

		CustomizeDiff: config.FlexibleForceNew(v3AgencyNonUpdatableParams),

		Timeouts: &schema.ResourceTimeout{
			Read:   schema.DefaultTimeout(2 * time.Minute),
			Update: schema.DefaultTimeout(1 * time.Minute),
			Delete: schema.DefaultTimeout(1 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			// Required parameters.
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the agency.`,
			},
			"delegated_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"delegated_service_name"},
				Description: utils.SchemaDesc(
					`The name of the delegated user domain.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},

			// Optional parameters.
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description (supplementary information) of the agency.`,
			},
			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FOREVER",
				DiffSuppressFunc: func(_, oldValue, newValue string, _ *schema.ResourceData) bool {
					return parseOneDayDuration(oldValue) == parseOneDayDuration(newValue)
				},
				Description: `The validity period of the agency.`,
			},
			"project_role": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the project.`,
						},
						"roles": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of role names used for assignment in a specified project.`,
						},
					},
				},
				Description: `The roles assignment for the agency which the projects are used to grant.`,
			},
			"domain_roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Set:         schema.HashString,
				Description: `The roles assignment for the agency which the domain are used to grant.`,
			},
			"all_resources_roles": {
				Type:        schema.TypeSet,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The roles assignment for the agency which the all resources are used to grant.`,
			},
			"enterprise_project_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"enterprise_project": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The name of the enterprise project.`,
						},
						"roles": {
							Type:        schema.TypeSet,
							Required:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `The list of role names used for assignment in a specified enterprise project.`,
						},
					},
				},
				Description: `The roles assignment for the agency which the enterprise projects are used to grant.`,
			},
			"force_dissociate_v5_policies": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: `Whether to force dissociate the associated v5 policies when deleting the agency.`,
			},

			// Attributes.
			"expire_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The expiration time of the agency.`,
			},
			"create_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the agency.`,
			},

			// Internal parameters.
			"delegated_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				Description: utils.SchemaDesc(
					`The name of the delegated service.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
		},
	}
}

func parseOneDayDuration(duration string) interface{} {
	if duration == "ONEDAY" {
		return "1"
	}

	return duration
}

func buildCreateV3AgencyDelegatedDomain(d *schema.ResourceData) interface{} {
	if domainName, ok := d.GetOk("delegated_domain_name"); ok {
		return domainName.(string)
	}

	return d.Get("delegated_service_name").(string)
}

// The type of duration can be string or int in Create and Update methods
func buildCreateV3AgencyDuration(d *schema.ResourceData) interface{} {
	raw := d.Get("duration").(string)
	if raw == "" {
		return nil
	}

	// Try to convert duration to int, if success, return the converted value
	if days, err := strconv.Atoi(raw); err == nil {
		return days
	}

	return raw
}

func buildCreateV3AgencyBodyParams(d *schema.ResourceData, domainId string) map[string]interface{} {
	return map[string]interface{}{
		"agency": map[string]interface{}{
			// Required parameters.
			"domain_id":         domainId,
			"name":              d.Get("name").(string),
			"trust_domain_name": buildCreateV3AgencyDelegatedDomain(d),
			// Optional parameters.
			"description": utils.ValueIgnoreEmpty(d.Get("description").(string)),
			"duration":    buildCreateV3AgencyDuration(d),
		},
	}
}

func createV3Agency(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) (interface{}, error) {
	httpUrl := "v3.0/OS-AGENCY/agencies"

	createPath := client.Endpoint + httpUrl
	createAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateV3AgencyBodyParams(d, domainId)),
	}

	requestResp, err := client.Request("POST", createPath, &createAgencyOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

// If the domainId is omitted, the system's built-in roles will be queried.
// Otherwise, custom roles under the specified domain will be queried.
func listRoles(client *golangsdk.ServiceClient, domainId ...string) ([]interface{}, error) {
	var (
		httpUrl = "v3/roles?per_page={per_page}"
		perPage = 300
		page    = 1
		result  = make([]interface{}, 0)
	)

	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{per_page}", strconv.Itoa(perPage))
	if len(domainId) > 0 {
		httpUrl = fmt.Sprintf("%s&domain_id=%s", httpUrl, domainId[0])
	}

	getRoleOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		httpUrlWithPage := fmt.Sprintf("%s&page=%d", httpUrl, page)
		requestResp, err := client.Request("GET", httpUrlWithPage, &getRoleOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		roles := utils.PathSearch("roles", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, roles...)
		if len(roles) < perPage {
			break
		}
		page++
	}

	return result, nil
}

func listAllRoles(client *golangsdk.ServiceClient, domainId string) ([]interface{}, error) {
	systemRoles, err := listRoles(client)
	if err != nil {
		return nil, fmt.Errorf("error listing system-defined roles: %s", err)
	}
	customRoles, err := listRoles(client, domainId)
	if err != nil {
		return nil, fmt.Errorf("error listing custom roles: %s", err)
	}

	return append(systemRoles, customRoles...), nil
}

// parseRolesToPairs parses the roles to pairs, the key is the role name, the value is the role ID.
// This method used to convert the list of characters into a mapping from character names to character IDs, which
// facilitates fast indexing later. The tree structure of the map is superior to slice in terms of performance.
func parseRolesToPairs(roles []interface{}) map[string]string {
	result := make(map[string]string)

	for _, role := range roles {
		roleName := utils.PathSearch("display_name", role, "").(string)
		roleId := utils.PathSearch("id", role, "").(string)
		if roleName == "" || roleId == "" {
			log.Printf("[WARN] invalid role name (%s) or ID (%s)", roleName, roleId)
			continue
		}
		result[roleName] = roleId
	}

	return result
}

func listProjects(client *golangsdk.ServiceClient, domainId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/projects?domain_id={domain_id}&per_page={per_page}"
		perPage = 5000
		page    = 1
		result  = make([]interface{}, 0)
	)

	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{domain_id}", domainId)
	httpUrl = strings.ReplaceAll(httpUrl, "{per_page}", strconv.Itoa(perPage))

	getProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		httpUrlWithPage := fmt.Sprintf("%s&page=%d", httpUrl, page)
		requestResp, err := client.Request("GET", httpUrlWithPage, &getProjectOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		projects := utils.PathSearch("projects", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, projects...)
		if len(projects) < perPage {
			break
		}
		page++
	}

	return result, nil
}

func parseProjectRolesToPairs(projectRoles *schema.Set) map[string]interface{} {
	// The key is the combination of project name and role name, which formatted as "project_name|role_name".
	// Map structure can ensure de-duplication, so there is no need to worry about duplication.
	result := make(map[string]interface{})

	for _, projectRole := range projectRoles.List() {
		projectName := utils.PathSearch("project", projectRole, "").(string)
		roleNames := utils.PathSearch("roles", projectRole, schema.NewSet(schema.HashString, nil)).(*schema.Set)

		// The project name and role names cannot be empty
		if projectName == "" || roleNames.Len() < 1 {
			log.Printf("[WARN] invalid project name (%s) or role names (%v)", projectName, roleNames.List())
			continue
		}

		for _, roleName := range roleNames.List() {
			// The key is the combination of project name and role name, which formatted as "project_name|role_name"
			result[fmt.Sprintf("%s|%v", projectName, roleName)] = true
		}
	}

	return result
}

// buildProjectRoles builds the project roles from the project roles set, the format of each element is "project_name|role_name"
func buildProjectRoles(projectRoles *schema.Set) []interface{} {
	return utils.PathSearch("keys(@)", parseProjectRolesToPairs(projectRoles), make([]interface{}, 0)).([]interface{})
}

func attachProjectRoleToV3Agency(client *golangsdk.ServiceClient, agencyId, projectId, roleId string) error {
	httpUrl := "v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}"

	attachPath := client.Endpoint + httpUrl
	attachPath = strings.ReplaceAll(attachPath, "{project_id}", projectId)
	attachPath = strings.ReplaceAll(attachPath, "{agency_id}", agencyId)
	attachPath = strings.ReplaceAll(attachPath, "{role_id}", roleId)

	attachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("PUT", attachPath, &attachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func attachProjectRolesToV3Agency(client *golangsdk.ServiceClient, parsedRolePairs map[string]string,
	projectRoles []interface{}, domainId, agencyId string) error {
	if len(projectRoles) < 1 {
		return nil
	}

	log.Printf("[DEBUG] attaching roles in project scope to agency (%s), the roles are %v", agencyId, projectRoles)

	projects, err := listProjects(client, domainId)
	if err != nil {
		return fmt.Errorf("error querying the projects of domain (%s): %s", domainId, err)
	}

	for _, projectRoleVal := range projectRoles {
		var projectName, roleName string
		projectRole, ok := projectRoleVal.(string)
		if !ok {
			log.Printf("[WARN] invalid type of project role (%[1]v), want string, but got %[1]T", projectRole)
			continue
		}
		projectRoleParts := strings.Split(projectRole, "|")
		if len(projectRoleParts) != 2 {
			log.Printf("[WARN] invalid project role (%v): invalid format", projectRole)
			continue
		}
		projectName = projectRoleParts[0]
		roleName = projectRoleParts[1]

		projectId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", projectName), projects, "").(string)
		roleId, ok := parsedRolePairs[roleName]
		if !ok {
			log.Printf("[ERROR] the role (%s) to be attached does not exist", roleName)
			continue
		}

		err = attachProjectRoleToV3Agency(client, agencyId, projectId, roleId)
		if err != nil {
			return fmt.Errorf("error attaching project role (%s) to agency (%s): %s",
				roleId, agencyId, err)
		}
	}

	return nil
}

func listEnterpriseProjects(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl = "v1.0/enterprise-projects?limit={limit}"
		limit   = 1000
		offset  = 0
		result  = make([]interface{}, 0)
	)

	httpUrl = client.Endpoint + httpUrl
	httpUrl = strings.ReplaceAll(httpUrl, "{limit}", strconv.Itoa(limit))

	getEnterpriseProjectOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		httpUrlWithOffset := fmt.Sprintf("%s&offset=%d", httpUrl, offset)
		requestResp, err := client.Request("GET", httpUrlWithOffset, &getEnterpriseProjectOpt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}
		enterpriseProjects := utils.PathSearch("enterprise_projects", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, enterpriseProjects...)
		if len(enterpriseProjects) < limit {
			break
		}
		offset += len(enterpriseProjects)
	}

	return result, nil
}

func attachDomainRoleToV3Agency(client *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	httpUrl := "v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}"

	attachPath := client.Endpoint + httpUrl
	attachPath = strings.ReplaceAll(attachPath, "{domain_id}", domainId)
	attachPath = strings.ReplaceAll(attachPath, "{agency_id}", agencyId)
	attachPath = strings.ReplaceAll(attachPath, "{role_id}", roleId)

	attachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("PUT", attachPath, &attachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func attachDomainRolesToV3Agency(client *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []interface{}, domainId, agencyId string) error {
	if len(roleNames) < 1 {
		return nil
	}

	log.Printf("[DEBUG] attaching roles %v in domain scope to agency %s", roleNames, agencyId)

	for _, roleName := range roleNames {
		roleNameStr, ok := roleName.(string)
		if !ok {
			log.Printf("[WARN] invalid type of role name (%[1]v), want string, but got %[1]T", roleName)
			continue
		}
		roleId, ok := allRoleIds[roleNameStr]
		if !ok {
			log.Printf("[WARN] the role (%s) to be attached does not exist", roleNameStr)
			continue
		}

		err := attachDomainRoleToV3Agency(client, agencyId, domainId, roleId)
		if err != nil {
			return fmt.Errorf("error attaching role (%s) to agency (%s) by domain (%s): %s",
				roleId, agencyId, domainId, err)
		}
	}

	return nil
}

func attachAllResourcesRoleToV3Agency(client *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	httpUrl := "v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects"

	attachPath := client.Endpoint + httpUrl
	attachPath = strings.ReplaceAll(attachPath, "{domain_id}", domainId)
	attachPath = strings.ReplaceAll(attachPath, "{agency_id}", agencyId)
	attachPath = strings.ReplaceAll(attachPath, "{role_id}", roleId)

	attachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err := client.Request("PUT", attachPath, &attachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func attachAllResourcesRolesToV3Agency(client *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []interface{}, domainId, agencyId string) error {
	if len(roleNames) < 1 {
		return nil
	}

	log.Printf("[DEBUG] attaching roles %v in all resources to agency %s", roleNames, agencyId)

	for _, roleName := range roleNames {
		roleNameStr, ok := roleName.(string)
		if !ok {
			log.Printf("[WARN] invalid type of role name (%[1]v), want string, but got %[1]T", roleName)
			continue
		}
		roleId, ok := allRoleIds[roleNameStr]
		if !ok {
			log.Printf("[WARN] the role (%s) to be attached does not exist", roleNameStr)
			continue
		}

		err := attachAllResourcesRoleToV3Agency(client, agencyId, domainId, roleId)
		if err != nil {
			return fmt.Errorf("error attaching role (%s) in all resources to agency (%s): %s",
				roleId, agencyId, err)
		}
	}

	return nil
}

func parseEnterpriseProjectRolesToPairs(enterpriseProjectRoles *schema.Set) map[string]bool {
	// The key is the combination of enterprise project name and role name, which formatted as "enterprise_project_name|role_name".
	// Map structure can ensure de-duplication, so there is no need to worry about duplication.
	result := make(map[string]bool)

	for _, enterpriseProjectRole := range enterpriseProjectRoles.List() {
		enterpriseProjectName := utils.PathSearch("enterprise_project", enterpriseProjectRole, "").(string)
		roleNames := utils.PathSearch("roles", enterpriseProjectRole, schema.NewSet(schema.HashString, nil)).(*schema.Set)

		// The enterprise project name and role names cannot be empty
		if enterpriseProjectName == "" || roleNames.Len() < 1 {
			log.Printf("[WARN] invalid enterprise project name (%s) or role names (%v)", enterpriseProjectName, roleNames.List())
			continue
		}

		for _, roleName := range roleNames.List() {
			// The key is the combination of enterprise project name and role name, which formatted as "enterprise_project_name|role_name"
			result[fmt.Sprintf("%s|%s", enterpriseProjectName, roleName.(string))] = true
		}
	}

	return result
}

func buildEnterpriseProjectRoles(enterpriseProjectRoles *schema.Set) []interface{} {
	return utils.PathSearch("keys(@)", parseEnterpriseProjectRolesToPairs(enterpriseProjectRoles), make([]interface{}, 0)).([]interface{})
}

func attachEnterpriseProjectRolesToV3Agency(iamClient, epsClient *golangsdk.ServiceClient, parsedRolePairs map[string]string,
	enterpriseProjectRoles []interface{}, agencyId string) error {
	if len(enterpriseProjectRoles) < 1 {
		return nil
	}

	log.Printf("[DEBUG] attaching roles %v in enterprise project scope to agency %s", enterpriseProjectRoles, agencyId)

	enterpriseProjects, err := listEnterpriseProjects(epsClient)
	if err != nil {
		return fmt.Errorf("error querying the enterprise projects: %s", err)
	}

	var (
		httpUrl     = "v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments"
		attachRoles = make([]interface{}, 0, len(enterpriseProjectRoles))
	)

	for _, enterpriseProjectRoleVal := range enterpriseProjectRoles {
		var enterpriseProjectName, roleName, roleId string
		enterpriseProjectRole, ok := enterpriseProjectRoleVal.(string)
		if !ok {
			log.Printf("[WARN] invalid type of enterprise project role (%[1]v), want string, but got %[1]T", enterpriseProjectRole)
			continue
		}
		enterpriseProjectRolePair := strings.Split(enterpriseProjectRole, "|")
		if len(enterpriseProjectRolePair) != 2 {
			log.Printf("[WARN] invalid enterprise project role (%v): invalid format", enterpriseProjectRole)
			continue
		}
		enterpriseProjectName = enterpriseProjectRolePair[0]
		roleName = enterpriseProjectRolePair[1]

		enterpriseProjectId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", enterpriseProjectName), enterpriseProjects, "").(string)
		roleId, ok = parsedRolePairs[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be attached does not exist", roleName)
			continue
		}

		attachRoles = append(attachRoles, map[string]interface{}{
			"agency_id":             agencyId,
			"enterprise_project_id": enterpriseProjectId,
			"role_id":               roleId,
		})
	}

	if len(attachRoles) < 1 {
		log.Printf("[DEBUG] no roles to attach by enterprise project to agency %s", agencyId)
		return nil
	}

	attachPath := iamClient.Endpoint + httpUrl
	attachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"role_assignments": attachRoles,
		},
	}

	_, err = iamClient.Request("PUT", attachPath, &attachOpt)
	if err != nil {
		return fmt.Errorf("error attaching roles by enterprise project to agency (%s): %s", agencyId, err)
	}

	return nil
}

func resourceV3AgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	iamClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	domainId := cfg.DomainID
	if domainId == "" {
		return diag.Errorf("the parameter 'domain_id' in provider-level configuration must be specified before creating agency")
	}

	respBody, err := createV3Agency(iamClient, d, domainId)
	if err != nil {
		return diag.Errorf("error creating agency: %s", err)
	}

	agencyId := utils.PathSearch("agency.id", respBody, "").(string)
	if agencyId == "" {
		return diag.Errorf("unable to find the agency ID from the API response")
	}
	d.SetId(agencyId)

	// get all of the role IDs, include system-defined roles and custom roles
	allRoles, err := listAllRoles(iamClient, domainId)
	if err != nil {
		return diag.FromErr(err)
	}
	parsedRolePairs := parseRolesToPairs(allRoles)

	if projectRoles := d.Get("project_role").(*schema.Set); projectRoles.Len() > 0 {
		if err := attachProjectRolesToV3Agency(iamClient, parsedRolePairs, buildProjectRoles(projectRoles), domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if domainRoles := d.Get("domain_roles").(*schema.Set); domainRoles.Len() > 0 {
		if err := attachDomainRolesToV3Agency(iamClient, parsedRolePairs, domainRoles.List(), domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if allResourceRoles := d.Get("all_resources_roles").(*schema.Set); allResourceRoles.Len() > 0 {
		if err := attachAllResourcesRolesToV3Agency(iamClient, parsedRolePairs, allResourceRoles.List(), domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if epRawRoles := d.Get("enterprise_project_roles").(*schema.Set); epRawRoles.Len() > 0 {
		epRoles := buildEnterpriseProjectRoles(epRawRoles)
		epsClient, err := cfg.EnterpriseProjectClient(region)
		if err != nil {
			return diag.Errorf("error creating EPS client: %s", err)
		}
		if err = attachEnterpriseProjectRolesToV3Agency(iamClient, epsClient, parsedRolePairs, epRoles, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3AgencyRead(ctx, d, meta)
}

// The value can be "FOREVER" or the period in hours
// We should convert the period in days
func normalizeAgencyDuration(duration interface{}) interface{} {
	var result string

	switch v := duration.(type) {
	case string:
		if hours, err := strconv.Atoi(v); err == nil {
			days := hours / 24
			result = strconv.Itoa(days)
		} else {
			result = v
		}
	case int:
		days := v / 24
		result = strconv.Itoa(days)
	default:
		result = "FOREVER"
	}

	return result
}

func getV3AgencyById(client *golangsdk.ServiceClient, agencyId string) (interface{}, error) {
	httpUrl := "v3.0/OS-AGENCY/agencies/{agency_id}"

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{agency_id}", agencyId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	respBody, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(respBody)
}

func GetV3AgencyByIdWithRetry(ctx context.Context, client *golangsdk.ServiceClient, agencyId string, timeout ...time.Duration) (interface{}, error) {
	var (
		respBody   interface{}
		err        error
		timeoutVal time.Duration
	)

	if len(timeout) < 1 || timeout[0] <= time.Duration(0) {
		return getV3AgencyById(client, agencyId)
	}
	timeoutVal = timeout[0]

	// lintignore:R006
	err = resource.RetryContext(ctx, timeoutVal, func() *resource.RetryError {
		respBody, err = getV3AgencyById(client, agencyId)
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			// Retrieving agency details may result in a 404 error, requiring appropriate retries.
			// If the details are not retrieved within the timeout period, an error will be returned.
			// lintignore:R018
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return nil
	})

	return respBody, err
}

func listAttachedProjectRolesForV3AgencyByProjectId(client *golangsdk.ServiceClient, agencyId, projectId string) ([]interface{}, error) {
	var httpUrl = "v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", projectId)
	listPath = strings.ReplaceAll(listPath, "{agency_id}", agencyId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("roles", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func listAttachedProjectRolesForV3Agency(client *golangsdk.ServiceClient, domainId, agencyId string) ([]interface{}, error) {
	projects, err := listProjects(client, domainId)
	if err != nil {
		return nil, fmt.Errorf("error querying the projects of domain (%s): %s", domainId, err)
	}

	result := make([]interface{}, 0, len(projects))
	for _, project := range projects {
		projectId := utils.PathSearch("id", project, "").(string)
		projectName := utils.PathSearch("name", project, "").(string)

		// MOS is a special project for CBC service which is used for billing, not visible to the user
		if projectId == "MOS" {
			continue
		}

		// the provider will query the roles in all projects, but the API rate limit threshold is 10 times per second.
		// so we should wait for some time to avoid exceeding the rate limit.
		// lintignore:R018
		time.Sleep(200 * time.Millisecond)

		attachedProjectRoles, err := listAttachedProjectRolesForV3AgencyByProjectId(client, agencyId, projectId)
		if err != nil && !utils.IsResourceNotFound(err) {
			log.Printf("[ERROR] error querying the roles attached on project(%s): %s", projectName, err)
			continue
		}
		if len(attachedProjectRoles) < 1 {
			continue
		}

		result = append(result, map[string]interface{}{
			"project": projectName,
			"roles":   utils.PathSearch("[*].display_name", attachedProjectRoles, make([]interface{}, 0)).([]interface{}),
		})
	}

	return result, nil
}

func listAttachedDomainRolesForV3AgencyByDomainId(client *golangsdk.ServiceClient, domainId, agencyId string) ([]interface{}, error) {
	var httpUrl = "v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles"

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{domain_id}", domainId)
	listPath = strings.ReplaceAll(listPath, "{agency_id}", agencyId)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch("roles", respBody, make([]interface{}, 0)).([]interface{}), nil
}

func resourceV3AgencyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Id()
		timeout  time.Duration
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.IsNewResource() {
		timeout = d.Timeout(schema.TimeoutRead)
	}
	agency, err := GetV3AgencyByIdWithRetry(ctx, client, agencyId, timeout)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving agency")
	}

	mErr := multierror.Append(nil,
		d.Set("name", utils.PathSearch("agency.name", agency, "")),
		d.Set("description", utils.PathSearch("agency.description", agency, "")),
		d.Set("expire_time", utils.PathSearch("agency.expire_time", agency, "")),
		d.Set("create_time", utils.PathSearch("agency.create_time", agency, "")),
		d.Set("duration", normalizeAgencyDuration(utils.PathSearch("agency.duration", agency, ""))),
	)

	delegatedDomainName := utils.PathSearch("agency.trust_domain_name", agency, "").(string)
	match, _ := regexp.MatchString("^op_svc_[A-Za-z]+$", delegatedDomainName)
	if match {
		mErr = multierror.Append(mErr, d.Set("delegated_service_name", delegatedDomainName))
	} else {
		mErr = multierror.Append(mErr, d.Set("delegated_domain_name", delegatedDomainName))
	}

	domainId := utils.PathSearch("agency.domain_id", agency, "").(string)
	projectRoles, err := listAttachedProjectRolesForV3Agency(client, domainId, agencyId)
	if err != nil {
		log.Printf("[ERROR] error querying the roles attached on project for agency (%s): %s", agencyId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("project_role", projectRoles))
	}

	domainRoles, err := listAttachedDomainRolesForV3AgencyByDomainId(client, domainId, agencyId)
	if err != nil {
		log.Printf("[ERROR] error querying the roles attached on domain for agency (%s): %s", agencyId, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("domain_roles", utils.PathSearch("[*].display_name", domainRoles,
			make([]interface{}, 0)).([]interface{})))
	}

	// Unable to fetch all_resources_roles because the API response does not include `display_name` field
	// https://support.huaweicloud.com/api-iam/iam_12_0014.html

	// Unable to fetch enterprise_project_roles because the API response does not include `display_name` field
	// https://support.huaweicloud.com/api-iam/iam_10_0014.html
	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateV3AgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"agency": map[string]interface{}{
			"description":      d.Get("description").(string),
			"duration":         buildCreateV3AgencyDuration(d),
			"delegated_domain": buildCreateV3AgencyDelegatedDomain(d),
		},
	}
}

func updateV3Agency(ctx context.Context, client *golangsdk.ServiceClient, agencyId string, d *schema.ResourceData,
	timeout time.Duration) error {
	httpUrl := "v3.0/OS-AGENCY/agencies/{agency_id}"

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{agency_id}", agencyId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildUpdateV3AgencyBodyParams(d)),
	}

	// lintignore:R006
	err := resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		_, retryErr := client.Request("PUT", updatePath, &updateOpt)
		if retryErr != nil {
			return common.CheckForRetryableError(retryErr)
		}
		// lintignore:R018
		time.Sleep(10 * time.Second)
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func diffChangesOfProjectRolesForV3Agency(oldVal, newVal *schema.Set) (removeProjectRoles, addProjectRoles []interface{}) {
	removeProjectRoles = make([]interface{}, 0)
	addProjectRoles = make([]interface{}, 0)

	oldProjectRolePairs := parseProjectRolesToPairs(oldVal)
	newProjectRolePairs := parseProjectRolesToPairs(newVal)

	for k := range oldProjectRolePairs {
		if _, ok := newProjectRolePairs[k]; !ok {
			removeProjectRoles = append(removeProjectRoles, k)
		}
	}

	for k := range newProjectRolePairs {
		if _, ok := oldProjectRolePairs[k]; !ok {
			addProjectRoles = append(addProjectRoles, k)
		}
	}

	return removeProjectRoles, addProjectRoles
}

func detachProjectRoleFromV3Agency(client *golangsdk.ServiceClient, agencyId, projectId, roleId string) error {
	httpUrl := "v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}"

	detachPath := client.Endpoint + httpUrl
	detachPath = strings.ReplaceAll(detachPath, "{project_id}", projectId)
	detachPath = strings.ReplaceAll(detachPath, "{agency_id}", agencyId)
	detachPath = strings.ReplaceAll(detachPath, "{role_id}", roleId)

	detachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", detachPath, &detachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func detachProjectRolesFromV3Agency(client *golangsdk.ServiceClient, parsedRolePairs map[string]string, projectRoles []interface{},
	domainId, agencyId string) error {
	if len(projectRoles) < 1 {
		return nil
	}

	log.Printf("[DEBUG] detaching roles %v in project scope from agency %s, the roles are %v", projectRoles, agencyId, projectRoles)

	projects, err := listProjects(client, domainId)
	if err != nil {
		return fmt.Errorf("error querying the projects of domain (%s): %s", domainId, err)
	}

	for _, projectRoleVal := range projectRoles {
		var projectName, roleName string
		projectRole, ok := projectRoleVal.(string)
		if !ok {
			log.Printf("[WARN] invalid type of project role (%[1]v), want string, but got %[1]T", projectRole)
			continue
		}
		projectRolePair := strings.Split(projectRole, "|")
		if len(projectRolePair) != 2 {
			log.Printf("[WARN] invalid project role (%v): invalid format", projectRole)
			continue
		}
		projectName = projectRolePair[0]
		roleName = projectRolePair[1]

		projectId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", projectName), projects, "").(string)
		roleId, ok := parsedRolePairs[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", roleName)
			continue
		}

		err = detachProjectRoleFromV3Agency(client, agencyId, projectId, roleId)
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role (%s) by project (%s) from agency (%s): %s",
				roleId, projectId, agencyId, err)
		}
	}

	return nil
}

func updateProjectRolesForV3Agency(client *golangsdk.ServiceClient, d *schema.ResourceData, parsedRolePairs map[string]string,
	domainId, agencyId string) error {
	var (
		oldRaw, newRaw                      = d.GetChange("project_role")
		removeProjectRoles, addProjectRoles = diffChangesOfProjectRolesForV3Agency(oldRaw.(*schema.Set), newRaw.(*schema.Set))
	)

	if len(removeProjectRoles) > 0 {
		if err := detachProjectRolesFromV3Agency(client, parsedRolePairs, removeProjectRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	if len(addProjectRoles) > 0 {
		if err := attachProjectRolesToV3Agency(client, parsedRolePairs, addProjectRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	return nil
}

func detachDomainRoleFromV3Agency(client *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	httpUrl := "v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}"

	detachPath := client.Endpoint + httpUrl
	detachPath = strings.ReplaceAll(detachPath, "{domain_id}", domainId)
	detachPath = strings.ReplaceAll(detachPath, "{agency_id}", agencyId)
	detachPath = strings.ReplaceAll(detachPath, "{role_id}", roleId)

	detachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", detachPath, &detachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func detachDomainRolesFromV3Agency(client *golangsdk.ServiceClient, parsedRolePairs map[string]string,
	roleNames []interface{}, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in domain scope from agency %s", roleNames, agencyId)
	}

	for _, roleNameVal := range roleNames {
		roleName, ok := roleNameVal.(string)
		if !ok {
			log.Printf("[WARN] invalid type of role name (%[1]v), want string, but got %[1]T", roleName)
			continue
		}
		roleId, ok := parsedRolePairs[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", roleName)
			continue
		}

		err := detachDomainRoleFromV3Agency(client, agencyId, domainId, roleId)
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role (%s) by domain (%s) from agency (%s): %s",
				roleId, domainId, agencyId, err)
		}
	}

	return nil
}

func updateDomainRolesForV3Agency(client *golangsdk.ServiceClient, d *schema.ResourceData, parsedRolePairs map[string]string,
	domainId, agencyId string) error {
	var (
		oldRaw, newRaw    = d.GetChange("domain_roles")
		deleteDomainRoles = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set)).List()
		addDomainRoles    = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set)).List()
	)

	if len(deleteDomainRoles) > 0 {
		if err := detachDomainRolesFromV3Agency(client, parsedRolePairs, deleteDomainRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	if len(addDomainRoles) > 0 {
		if err := attachDomainRolesToV3Agency(client, parsedRolePairs, addDomainRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	return nil
}

func detachAllResourcesRoleFromV3Agency(client *golangsdk.ServiceClient, agencyId, domainId, roleId string) error {
	httpUrl := "v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects"

	detachPath := client.Endpoint + httpUrl
	detachPath = strings.ReplaceAll(detachPath, "{domain_id}", domainId)
	detachPath = strings.ReplaceAll(detachPath, "{agency_id}", agencyId)
	detachPath = strings.ReplaceAll(detachPath, "{role_id}", roleId)

	detachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", detachPath, &detachOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func detachAllResourcesRolesFromV3Agency(client *golangsdk.ServiceClient, parsedRolePairs map[string]string,
	roleNames []interface{}, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in all resources from agency %s", roleNames, agencyId)
	}

	for _, roleName := range roleNames {
		roleNameStr, ok := roleName.(string)
		if !ok {
			log.Printf("[WARN] invalid type of role name (%[1]v), want string, but got %[1]T", roleName)
			continue
		}
		roleId, ok := parsedRolePairs[roleNameStr]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", roleNameStr)
			continue
		}

		err := detachAllResourcesRoleFromV3Agency(client, agencyId, domainId, roleId)
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role (%s) in all resources from agency (%s): %s",
				roleId, agencyId, err)
		}
	}

	return nil
}

func updateAllResourcesRolesForV3Agency(client *golangsdk.ServiceClient, d *schema.ResourceData,
	parsedRolePairs map[string]string, domainId, agencyId string) error {
	var (
		oldRaw, newRaw          = d.GetChange("all_resources_roles")
		deleteAllResourcesRoles = oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set)).List()
		addAllResourcesRoles    = newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set)).List()
	)

	if len(deleteAllResourcesRoles) > 0 {
		if err := detachAllResourcesRolesFromV3Agency(client, parsedRolePairs, deleteAllResourcesRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	if len(addAllResourcesRoles) > 0 {
		if err := attachAllResourcesRolesToV3Agency(client, parsedRolePairs, addAllResourcesRoles, domainId, agencyId); err != nil {
			return err
		}
	}

	return nil
}

func diffChangesOfEnterpriseProjectRolesForV3Agency(oldVal, newVal *schema.Set) (deleteEnterpriseProjectRoles,
	addEnterpriseProjectRoles []interface{}) {
	deleteEnterpriseProjectRoles = make([]interface{}, 0)
	addEnterpriseProjectRoles = make([]interface{}, 0)

	oldEnterpriseProjectRolePairs := parseEnterpriseProjectRolesToPairs(oldVal)
	newEnterpriseProjectRolePairs := parseEnterpriseProjectRolesToPairs(newVal)

	for k := range oldEnterpriseProjectRolePairs {
		if _, ok := newEnterpriseProjectRolePairs[k]; !ok {
			deleteEnterpriseProjectRoles = append(deleteEnterpriseProjectRoles, k)
		}
	}

	for k := range newEnterpriseProjectRolePairs {
		if _, ok := oldEnterpriseProjectRolePairs[k]; !ok {
			addEnterpriseProjectRoles = append(addEnterpriseProjectRoles, k)
		}
	}

	return deleteEnterpriseProjectRoles, addEnterpriseProjectRoles
}

func detachEnterpriseProjectRolesFromV3Agency(iamClient, epsClient *golangsdk.ServiceClient, parsedRolePairs map[string]string,
	enterpriseProjectRoles []interface{}, agencyId string) error {
	if len(enterpriseProjectRoles) < 1 {
		return nil
	}

	log.Printf("[DEBUG] detaching roles %v in enterprise project scope from agency %s", enterpriseProjectRoles, agencyId)

	enterpriseProjects, err := listEnterpriseProjects(epsClient)
	if err != nil {
		return fmt.Errorf("error querying the enterprise projects: %s", err)
	}

	var (
		httpUrl     = "v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments"
		detachRoles = make([]interface{}, 0, len(enterpriseProjectRoles))
	)

	for _, enterpriseProjectRoleVal := range enterpriseProjectRoles {
		var enterpriseProjectName, roleName string
		enterpriseProjectRole, ok := enterpriseProjectRoleVal.(string)
		if !ok {
			log.Printf("[WARN] invalid type of enterprise project role (%[1]v), want string, but got %[1]T", enterpriseProjectRole)
			continue
		}
		enterpriseProjectRolePair := strings.Split(enterpriseProjectRole, "|")
		if len(enterpriseProjectRolePair) != 2 {
			log.Printf("[WARN] invalid enterprise project role (%v): invalid format", enterpriseProjectRole)
			continue
		}
		enterpriseProjectName = enterpriseProjectRolePair[0]
		roleName = enterpriseProjectRolePair[1]

		enterpriseProjectId := utils.PathSearch(fmt.Sprintf("[?name=='%s']|[0].id", enterpriseProjectName), enterpriseProjects, "").(string)
		roleId, ok := parsedRolePairs[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", roleName)
			continue
		}

		detachRoles = append(detachRoles, map[string]interface{}{
			"agency_id":             agencyId,
			"enterprise_project_id": enterpriseProjectId,
			"role_id":               roleId,
		})
	}

	if len(detachRoles) < 1 {
		log.Printf("[DEBUG] no roles to detach by enterprise project from agency %s", agencyId)
		return nil
	}

	detachPath := iamClient.Endpoint + httpUrl
	detachOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"role_assignments": detachRoles,
		},
	}

	_, err = iamClient.Request("DELETE", detachPath, &detachOpt)
	if err != nil {
		return fmt.Errorf("error detaching roles by enterprise project from agency (%s): %s", agencyId, err)
	}

	return nil
}

func updateEnterpriseProjectRoles(iamClient, epsClient *golangsdk.ServiceClient, d *schema.ResourceData,
	parsedRolePairs map[string]string, agencyId string) error {
	var (
		oldRaw, newRaw                                          = d.GetChange("enterprise_project_roles")
		deleteEnterpriseProjectRoles, addEnterpriseProjectRoles = diffChangesOfEnterpriseProjectRolesForV3Agency(oldRaw.(*schema.Set),
			newRaw.(*schema.Set))
	)

	if len(deleteEnterpriseProjectRoles) > 0 {
		if err := detachEnterpriseProjectRolesFromV3Agency(iamClient, epsClient, parsedRolePairs, deleteEnterpriseProjectRoles,
			agencyId); err != nil {
			return err
		}
	}

	if len(addEnterpriseProjectRoles) > 0 {
		if err := attachEnterpriseProjectRolesToV3Agency(iamClient, epsClient, parsedRolePairs, addEnterpriseProjectRoles, agencyId); err != nil {
			return err
		}
	}

	return nil
}

func resourceV3AgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Id()
		domainId = cfg.DomainID
	)

	iamClient, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if domainId == "" {
		return diag.Errorf("the parameter 'domain_id' in provider-level configuration must be specified")
	}

	if d.HasChanges("delegated_domain_name", "delegated_service_name", "description", "duration") {
		if err = updateV3Agency(ctx, iamClient, agencyId, d, d.Timeout(schema.TimeoutUpdate)); err != nil {
			return diag.Errorf("error updating agency (%s): %s", agencyId, err)
		}
	}

	var parsedRolePairs map[string]string
	if d.HasChanges("project_role", "domain_roles", "all_resources_roles", "enterprise_project_roles") {
		// get all of the role IDs, include system-defined roles and custom roles
		allRoles, err := listAllRoles(iamClient, domainId)
		if err != nil {
			return diag.FromErr(err)
		}
		parsedRolePairs = parseRolesToPairs(allRoles)
	}

	if d.HasChange("project_role") {
		if err = updateProjectRolesForV3Agency(iamClient, d, parsedRolePairs, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("domain_roles") {
		if err = updateDomainRolesForV3Agency(iamClient, d, parsedRolePairs, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("all_resources_roles") {
		if err = updateAllResourcesRolesForV3Agency(iamClient, d, parsedRolePairs, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_roles") {
		epsClient, err := cfg.EnterpriseProjectClient(region)
		if err != nil {
			return diag.Errorf("error creating EPS client: %s", err)
		}
		if err = updateEnterpriseProjectRoles(iamClient, epsClient, d, parsedRolePairs, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3AgencyRead(ctx, d, meta)
}

func listV5AgencyAssociatedPolicies(client *golangsdk.ServiceClient, agencyId string) ([]interface{}, error) {
	var (
		httpUrl = "v5/agencies/{agency_id}/attached-policies?limit={limit}"
		limit   = 100
		marker  = ""
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{agency_id}", agencyId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := listPath
		if marker != "" {
			currentPath = fmt.Sprintf("%s&marker=%s", currentPath, marker)
		}
		resp, err := client.Request("GET", currentPath, &opt)
		if err != nil {
			return nil, err
		}
		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}
		policies := utils.PathSearch("attached_policies", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, policies...)
		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func deleteV5PoliciesFromAgency(client *golangsdk.ServiceClient, agencyId string, policies []interface{}) error {
	httpUrl := "v5/policies/{policy_id}/detach-agency"

	for _, policy := range policies {
		policyId := utils.PathSearch("id", policy, "").(string)

		detachPath := client.Endpoint + httpUrl
		detachPath = strings.ReplaceAll(detachPath, "{policy_id}", policyId)

		detachOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: map[string]interface{}{
				"agency_id": agencyId,
			},
		}

		_, err := client.Request("POST", detachPath, &detachOpt)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				log.Printf("[WARN] the policy (%s) was already detached from the agency (%s)", policyId, agencyId)
				continue
			}
			// Note: here we cannot format the error, otherwise the original status code will be lost
			return err
		}
	}
	return nil
}

func deleteV3Agency(client *golangsdk.ServiceClient, agencyId string) error {
	httpUrl := "v3.0/OS-AGENCY/agencies/{agency_id}"

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{agency_id}", agencyId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	// Note: here we cannot format the error, otherwise the original status code will be lost
	return err
}

func deleteV3AgencyWithRetry(ctx context.Context, client *golangsdk.ServiceClient, agencyId string, timeout time.Duration) error {
	// lintignore:R006
	return resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := deleteV3Agency(client, agencyId)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return resource.NonRetryableError(err)
			}
			// Retrieving agency details may result in a 404 error, requiring appropriate retries.
			// If the details are not retrieved within the timeout period, an error will be returned.
			// lintignore:R018
			time.Sleep(10 * time.Second)
			return resource.RetryableError(err)
		}
		return nil
	})
}

func resourceV3AgencyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Id()
		timeout  = d.Timeout(schema.TimeoutDelete)
	)

	client, err := cfg.NewServiceClient("iam", region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if d.Get("force_dissociate_v5_policies").(bool) {
		policies, err := listV5AgencyAssociatedPolicies(client, agencyId)
		if err != nil {
			log.Printf("[WARN] error listing associated v5 policies with the agency (%s): %s", agencyId, err)
		} else {
			err = deleteV5PoliciesFromAgency(client, agencyId, policies)
			if err != nil {
				return diag.Errorf("error dissociating v5 policies from the agency (%s): %s", agencyId, err)
			}
		}
	}

	err = deleteV3AgencyWithRetry(ctx, client, agencyId, timeout)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting agency (%s)", agencyId))
	}

	return nil
}
