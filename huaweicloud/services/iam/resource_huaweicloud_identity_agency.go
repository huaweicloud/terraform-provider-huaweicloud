package iam

import (
	"context"
	"errors"
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
	"github.com/chnsz/golangsdk/openstack/eps/v1/enterpriseprojects"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/eps_permissions"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var agencyNonUpdatableParams = []string{"name"}

// ResourceV3Agency
// @API IAM POST /v3.0/OS-AGENCY/agencies
// @API IAM GET /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM PUT /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM DELETE /v3.0/OS-AGENCY/agencies/{agency_id}
// @API IAM GET /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles
// @API IAM PUT /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-AGENCY/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM GET /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/inherited_to_projects
// @API IAM PUT /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects
// @API IAM DELETE /v3.0/OS-INHERIT/domains/{domain_id}/agencies/{agency_id}/roles/{role_id}/inherited_to_projects
// @API IAM GET /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles
// @API IAM PUT /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-AGENCY/projects/{project_id}/agencies/{agency_id}/roles/{role_id}
// @API IAM GET /v3.0/OS-PERMISSION/role-assignments
// @API IAM PUT /v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments
// @API IAM DELETE /v3.0/OS-PERMISSION/subjects/agency/scopes/enterprise-project/role-assignments
// @API IAM GET /v3/projects
// @API IAM GET /v3/roles
func ResourceV3Agency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3AgencyCreate,
		ReadContext:   resourceV3AgencyRead,
		UpdateContext: resourceV3AgencyUpdate,
		DeleteContext: resourceV3AgencyDelete,

		CustomizeDiff: config.FlexibleForceNew(agencyNonUpdatableParams),

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

func parseOneDayDuration(duration string) string {
	if duration == "ONEDAY" {
		return "1"
	}
	return duration
}

func getProjectIdByName(client *golangsdk.ServiceClient, domainId, name string) (string, error) {
	opts := projects.ListOpts{
		DomainID: domainId,
		Name:     name,
	}
	allPages, err := projects.List(client, &opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("failed to query projects: %s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return "", fmt.Errorf("failed to extract projects: %s", err)
	}

	if len(all) == 0 {
		return "", fmt.Errorf("can not find the ID of project %s", name)
	}

	item := all[0]
	return item.ID, nil
}

func getEnterpriseProjectByName(client *golangsdk.ServiceClient, name string) (string, error) {
	opts := enterpriseprojects.ListOpts{Name: name}
	enterpriseProjects, err := enterpriseprojects.List(client, opts).Extract()
	if err != nil {
		return "", fmt.Errorf("error retrieving enterprise projects: %s", err)
	}
	for _, v := range enterpriseProjects {
		if v.Name == name {
			// Use name to find the target enterprise project.
			return v.ID, nil
		}
	}
	return "", errors.New("your enterprise project doesn't exist")
}

func getAllProjectsByDomainId(client *golangsdk.ServiceClient, domainId string) (map[string]string, error) {
	opts := projects.ListOpts{
		DomainID: domainId,
	}
	allPages, err := projects.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("failed to query projects: %s", err)
	}

	allItems, err := projects.ExtractProjects(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to extract projects: %s", err)
	}

	all := make(map[string]string, len(allItems))
	for _, item := range allItems {
		all[item.Name] = item.ID
	}

	return all, nil
}

func listRolesByDomainId(client *golangsdk.ServiceClient, domainId string) (map[string]string, error) {
	opts := roles.ListOpts{
		DomainID: domainId,
	}
	allPages, err := roles.ListWithPages(client, opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %s", err)
	}

	allItems, err := roles.ExtractOffsetRoles(allPages)
	if err != nil {
		return nil, fmt.Errorf("failed to extract roles: %s", err)
	}
	if len(allItems) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(allItems))
	for _, item := range allItems {
		if name := item.DisplayName; name != "" {
			r[name] = item.ID
		} else {
			log.Printf("[WARN] role %s without display name", item.Name)
		}
	}

	return r, nil
}

func getAllRolesByDomain(client *golangsdk.ServiceClient, domainId string) (map[string]string, error) {
	systemRoles, err := listRolesByDomainId(client, "")
	if err != nil {
		return nil, fmt.Errorf("error listing system-defined roles: %s", err)
	}

	customRoles, err := listRolesByDomainId(client, domainId)
	if err != nil {
		return nil, fmt.Errorf("error listing custom roles: %s", err)
	}

	if systemRoles == nil {
		return customRoles, nil
	}

	if customRoles == nil {
		return systemRoles, nil
	}

	// merge customRoles into systemRoles
	for k, v := range customRoles {
		systemRoles[k] = v
	}
	return systemRoles, nil
}

func buildDelegatedDomain(d *schema.ResourceData) string {
	if domainName, ok := d.GetOk("delegated_domain_name"); ok {
		return domainName.(string)
	}
	return d.Get("delegated_service_name").(string)
}

// the type of duration can be string or int in Create and Update methods
func buildAgencyDuration(d *schema.ResourceData) interface{} {
	raw := d.Get("duration").(string)
	if raw == "" {
		return nil
	}

	// try to convert duration to int, if suceess, return the converted value
	if days, err := strconv.Atoi(raw); err == nil {
		return days
	}
	return raw
}

func buildCreateAgencyRequestBody(d *schema.ResourceData, domainId string) agency.CreateOpts {
	return agency.CreateOpts{
		DomainID:        domainId,
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Duration:        buildAgencyDuration(d),
		DelegatedDomain: buildDelegatedDomain(d),
	}
}

func changeToProjectRolePairs(projectRoles *schema.Set) map[string]bool {
	result := make(map[string]bool)
	for _, projectRole := range projectRoles.List() {
		projectRoleMap, ok := projectRole.(map[string]interface{})
		if !ok || len(projectRoleMap) < 1 {
			continue
		}

		projectName := projectRoleMap["project"].(string)
		roleNames := projectRoleMap["roles"].(*schema.Set)
		for _, roleName := range roleNames.List() {
			// The key is the combination of project name and role name, which formatted as "project_name|role_name"
			result[fmt.Sprintf("%s|%s", projectName, roleName.(string))] = true
		}
	}
	return result
}

func buildProjectRoles(projectRoles *schema.Set) []string {
	pairs := changeToProjectRolePairs(projectRoles)

	result := make([]string, 0, len(pairs))
	for projectRoleName := range pairs {
		result = append(result, projectRoleName)
	}
	return result
}

func changeToEnterpriseProjectRolePairs(enterpriseProjectRoles *schema.Set) map[string]bool {
	pairs := make(map[string]bool)
	for _, enterpriseProjectRole := range enterpriseProjectRoles.List() {
		enterpriseProjectRoleMap, ok := enterpriseProjectRole.(map[string]interface{})
		if !ok || len(enterpriseProjectRoleMap) < 1 {
			continue
		}

		enterpriseProjectName := enterpriseProjectRoleMap["enterprise_project"].(string)
		roleNames := enterpriseProjectRoleMap["roles"].(*schema.Set)
		for _, roleName := range roleNames.List() {
			// The key is the combination of enterprise project name and role name, which formatted as "enterprise_project_name|role_name"
			pairs[fmt.Sprintf("%s|%s", enterpriseProjectName, roleName.(string))] = true
		}
	}
	return pairs
}

func buildEnterpriseProjectRoles(enterpriseProjectRoles *schema.Set) []string {
	pairs := changeToEnterpriseProjectRolePairs(enterpriseProjectRoles)

	result := make([]string, 0, len(pairs))
	for enterpriseProjectRoleName := range pairs {
		result = append(result, enterpriseProjectRoleName)
	}
	return result
}

func resourceV3AgencyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamV3P0Client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}
	iamV3Client, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	domainId := cfg.DomainID
	if domainId == "" {
		return diag.Errorf("the parameter 'domain_id' in provider-level configuration must be specified")
	}

	agencyResp, err := agency.Create(iamV3P0Client, buildCreateAgencyRequestBody(d, domainId)).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM agency: %s", err)
	}

	agencyId := agencyResp.ID
	d.SetId(agencyId)

	// get all of the role IDs, include system-defined roles and custom roles
	allRoleIds, err := getAllRolesByDomain(iamV3Client, domainId)
	if err != nil {
		return diag.FromErr(err)
	}

	if rawRoles := d.Get("project_role").(*schema.Set); rawRoles.Len() > 0 {
		pRoles := buildProjectRoles(rawRoles)
		if err := attachProjectRoles(iamV3P0Client, iamV3Client, allRoleIds, pRoles, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if domainRoles := d.Get("domain_roles").(*schema.Set); domainRoles.Len() > 0 {
		if err := attachDomainRoles(iamV3P0Client, allRoleIds, utils.ExpandToStringListBySet(domainRoles), domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if inheritedRoles := d.Get("all_resources_roles").(*schema.Set); inheritedRoles.Len() > 0 {
		if err := attachAllResourcesRoles(iamV3P0Client, allRoleIds, utils.ExpandToStringListBySet(inheritedRoles), domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if epRawRoles := d.Get("enterprise_project_roles").(*schema.Set); epRawRoles.Len() > 0 {
		epRoles := buildEnterpriseProjectRoles(epRawRoles)
		epsClient, err := cfg.EnterpriseProjectClient(region)
		if err != nil {
			return diag.Errorf("error creating EPS client: %s", err)
		}
		if err = attachEnterpriseProjectRoles(iamV3P0Client, epsClient, allRoleIds, epRoles, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3AgencyRead(ctx, d, meta)
}

// The value can be "FOREVER" or the period in hours
// We should convert the period in days
func normalizeAgencyDuration(duration interface{}) string {
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

func getAgencyById(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration, agencyId string) (*agency.Agency, error) {
	var (
		agencyResp *agency.Agency
		err        error
	)

	// lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		agencyResp, err = agency.Get(client, agencyId).Extract()
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
	return agencyResp, err
}

func resourceV3AgencyRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		agencyId = d.Id()
	)

	iamV3P0Client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}
	iamV3Client, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	agencyResp, err := getAgencyById(ctx, iamV3P0Client, d.Timeout(schema.TimeoutRead), agencyId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving agency")
	}

	mErr := multierror.Append(nil,
		d.Set("name", agencyResp.Name),
		d.Set("description", agencyResp.Description),
		d.Set("expire_time", agencyResp.ExpireTime),
		d.Set("create_time", agencyResp.CreateTime),
		d.Set("duration", normalizeAgencyDuration(agencyResp.Duration)),
	)

	match, _ := regexp.MatchString("^op_svc_[A-Za-z]+$", agencyResp.DelegatedDomainName)
	if match {
		mErr = multierror.Append(mErr, d.Set("delegated_service_name", agencyResp.DelegatedDomainName))
	} else {
		mErr = multierror.Append(mErr, d.Set("delegated_domain_name", agencyResp.DelegatedDomainName))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity agency fields: %s", err)
	}

	allProjects, err := getAllProjectsByDomainId(iamV3Client, agencyResp.DomainID)
	if err != nil {
		return diag.Errorf("error querying the projects of domain: %s", err)
	}

	projectRoles := make([]interface{}, 0, len(allProjects))
	for projectName, projectId := range allProjects {
		// MOS is a special project, not visible to the user
		if projectId == "MOS" {
			continue
		}

		// the provider will query the roles in all projects, but the API rate limit threshold is 10 times per second.
		// so we should wait for some time to avoid exceeding the rate limit.
		// lintignore:R018
		time.Sleep(200 * time.Millisecond)

		allRoles, err := agency.ListRolesAttachedOnProject(iamV3P0Client, agencyId, projectId).ExtractRoles()
		if err != nil && !utils.IsResourceNotFound(err) {
			log.Printf("[ERROR] error querying the roles attached on project(%s): %s", projectName, err)
			continue
		}
		if len(allRoles) == 0 {
			continue
		}
		roleNamesUnderProject := make([]string, 0, len(allRoles))
		for _, role := range allRoles {
			roleNamesUnderProject = append(roleNamesUnderProject, role.DisplayName)
		}
		projectRoles = append(projectRoles, map[string]interface{}{
			"project": projectName,
			"roles":   roleNamesUnderProject,
		})
	}
	err = d.Set("project_role", projectRoles)
	if err != nil {
		log.Printf("[ERROR] unable to set the 'project_role' field: %s", err)
	}

	allDomainRoles, err := agency.ListRolesAttachedOnDomain(iamV3P0Client, agencyId, agencyResp.DomainID).ExtractRoles()
	if err != nil && !utils.IsResourceNotFound(err) {
		log.Printf("[ERROR] error querying the roles attached on domain: %s", err)
	}
	if len(allDomainRoles) != 0 {
		v := schema.Set{F: schema.HashString}
		for _, role := range allDomainRoles {
			v.Add(role.DisplayName)
		}
		err = d.Set("domain_roles", &v)
		if err != nil {
			log.Printf("[ERROR] unable to set the 'domain_roles' field: %s", err)
		}
	}

	// Unable to fetch all_resources_roles because the API response does not include `display_name` field
	// https://support.huaweicloud.com/api-iam/iam_12_0014.html

	// Unable to fetch enterprise_project_roles because the API response does not include `display_name` field
	// https://support.huaweicloud.com/api-iam/iam_10_0014.html
	return nil
}

func diffChangeOfEnterpriseProjectRole(oldVal, newVal *schema.Set) (remove, add []string) {
	remove = make([]string, 0)
	add = make([]string, 0)

	oldEnterpriseProjectRolePairs := changeToEnterpriseProjectRolePairs(oldVal)
	newEnterpriseProjectRolePairs := changeToEnterpriseProjectRolePairs(newVal)

	for k := range oldEnterpriseProjectRolePairs {
		if _, ok := newEnterpriseProjectRolePairs[k]; !ok {
			remove = append(remove, k)
		}
	}

	for k := range newEnterpriseProjectRolePairs {
		if _, ok := oldEnterpriseProjectRolePairs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func attachProjectRoles(iamClient, identityClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	projectRoles []string, domainId, agencyId string) error {
	if len(projectRoles) > 0 {
		log.Printf("[DEBUG] attaching roles %v in project scope to agency %s", projectRoles, agencyId)
	}

	for _, projectRole := range projectRoles {
		projectRolePair := strings.Split(projectRole, "|")
		if len(projectRolePair) != 2 {
			return fmt.Errorf("error parsing project role from %s: invalid format", projectRole)
		}

		projectId, err := getProjectIdByName(identityClient, domainId, projectRolePair[0])
		if err != nil {
			return fmt.Errorf("the project (%s) does not exist", projectRolePair[0])
		}
		roleId, ok := allRoleIds[projectRolePair[1]]
		if !ok {
			return fmt.Errorf("the role (%s) to be attached does not exist", projectRolePair[1])
		}

		err = agency.AttachRoleByProject(iamClient, agencyId, projectId, roleId).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role (%s) by project (%s) to agency (%s): %s",
				roleId, projectId, agencyId, err)
		}
	}

	return nil
}

func detachProjectRoles(iamClient, identityClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	projectRoles []string, domainId, agencyId string) error {
	if len(projectRoles) > 0 {
		log.Printf("[DEBUG] detaching roles %v in project scope from agency %s", projectRoles, agencyId)
	}

	for _, projectRole := range projectRoles {
		projectRolePair := strings.Split(projectRole, "|")
		if len(projectRolePair) != 2 {
			return fmt.Errorf("error parsing project role from %s: invalid format", projectRole)
		}

		projectId, err := getProjectIdByName(identityClient, domainId, projectRolePair[0])
		if err != nil {
			return fmt.Errorf("the project (%s) does not exist", projectRolePair[0])
		}

		roleId, ok := allRoleIds[projectRolePair[1]]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", projectRolePair[1])
			continue
		}

		err = agency.DetachRoleByProject(iamClient, agencyId, projectId, roleId).ExtractErr()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role (%s) by project (%s) from agency (%s): %s",
				roleId, projectId, agencyId, err)
		}
	}

	return nil
}

func attachDomainRoles(iamClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []string, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] attaching roles %v in domain scope to agency %s", roleNames, agencyId)
	}

	for _, roleName := range roleNames {
		roleId, ok := allRoleIds[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be attached does not exist", roleName)
			continue
		}

		err := agency.AttachRoleByDomain(iamClient, agencyId, domainId, roleId).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role (%s) by domain (%s) to agency (%s): %s",
				roleId, domainId, agencyId, err)
		}
	}

	return nil
}

func detachDomainRoles(iamClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []string, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in domain scope from agency %s", roleNames, agencyId)
	}

	for _, roleName := range roleNames {
		roleId, ok := allRoleIds[roleName]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", roleName)
			continue
		}

		err := agency.DetachRoleByDomain(iamClient, agencyId, domainId, roleId).ExtractErr()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role (%s) by domain (%s) from agency (%s): %s",
				roleId, domainId, agencyId, err)
		}
	}

	return nil
}

func attachAllResourcesRoles(iamClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []string, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] attaching roles %v in all resources to agency %s", roleNames, agencyId)
	}

	for _, roleName := range roleNames {
		roleId, ok := allRoleIds[roleName]
		if !ok {
			return fmt.Errorf("the role (%s) to be attached does not exist", roleName)
		}

		err := agency.AttachAllResources(iamClient, agencyId, domainId, roleId).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role (%s) in all resources to agency (%s): %s",
				roleId, agencyId, err)
		}
	}

	return nil
}

func detachAllResourcesRoles(iamClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	roleNames []string, domainId, agencyId string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in all resources from agency %s", roleNames, agencyId)
	}

	for _, roleName := range roleNames {
		roleId, ok := allRoleIds[roleName]
		if !ok {
			return fmt.Errorf("the role (%s) to be detached does not exist", roleName)
		}

		err := agency.DetachAllResources(iamClient, agencyId, domainId, roleId).ExtractErr()
		if err != nil {
			return fmt.Errorf("error detaching role (%s) in all resources from agency (%s): %s",
				roleId, agencyId, err)
		}
	}

	return nil
}

func attachEnterpriseProjectRoles(iamClient, epsClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	enterpriseProjectRoles []string, agencyId string) error {
	if len(enterpriseProjectRoles) == 0 {
		return nil
	}

	roleAssignments := make([]eps_permissions.RoleAssignment, 0, len(enterpriseProjectRoles))
	for _, enterpriseProjectRole := range enterpriseProjectRoles {
		enterpriseProjectRolePair := strings.Split(enterpriseProjectRole, "|")
		if len(enterpriseProjectRolePair) != 2 {
			return fmt.Errorf("error parsing enterprise project role from %s: invalid format", enterpriseProjectRole)
		}
		epid, err := getEnterpriseProjectByName(epsClient, enterpriseProjectRolePair[0])
		if err != nil {
			return fmt.Errorf("the enterprise project (%s) does not exist", enterpriseProjectRolePair[0])
		}
		roleId, ok := allRoleIds[enterpriseProjectRolePair[1]]
		if !ok {
			return fmt.Errorf("the role (%s) to be attached does not exist", enterpriseProjectRolePair[1])
		}
		roleAssignments = append(roleAssignments, eps_permissions.RoleAssignment{
			AgencyID:            agencyId,
			EnterprisePorjectID: epid,
			RoleID:              roleId,
		})
	}
	opt := eps_permissions.AgencyPermissionsOpts{RoleAssignments: roleAssignments}
	err := eps_permissions.AgencyPermissionsCreate(iamClient, &opt).ExtractErr()
	if err != nil {
		return fmt.Errorf("error attaching roles by enterprise project to agency (%s), "+
			"roleAssignments[{agency_id, enterprise_project_id, role_id}]: %v, error: %s",
			agencyId, enterpriseProjectRoles, err)
	}
	return nil
}

func detachEnterpriseProjectRoles(iamClient, epsClient *golangsdk.ServiceClient, allRoleIds map[string]string,
	enterpriseProjectRoles []string, agencyId string) error {
	if len(enterpriseProjectRoles) == 0 {
		return nil
	}

	roleAssignments := make([]eps_permissions.RoleAssignment, 0, len(enterpriseProjectRoles))
	for _, enterpriseProjectRole := range enterpriseProjectRoles {
		enterpriseProjectRolePair := strings.Split(enterpriseProjectRole, "|")
		if len(enterpriseProjectRolePair) != 2 {
			return fmt.Errorf("error parsing enterprise project role from %s: invalid format", enterpriseProjectRole)
		}

		epid, err := getEnterpriseProjectByName(epsClient, enterpriseProjectRolePair[0])
		if err != nil {
			return fmt.Errorf("the enterprise project (%s) does not exist", enterpriseProjectRolePair[0])
		}
		roleId, ok := allRoleIds[enterpriseProjectRolePair[1]]
		if !ok {
			log.Printf("[WARN] the role (%s) to be detached does not exist", enterpriseProjectRolePair[1])
			continue
		}
		roleAssignments = append(roleAssignments, eps_permissions.RoleAssignment{
			AgencyID:            agencyId,
			EnterprisePorjectID: epid,
			RoleID:              roleId,
		})
	}
	opt := eps_permissions.AgencyPermissionsOpts{RoleAssignments: roleAssignments}
	err := eps_permissions.AgencyPermissionsDelete(iamClient, &opt).ExtractErr()
	if err != nil {
		return fmt.Errorf("error detaching roles by enterprise project from agency (%s), "+
			"roleAssignments[{agency_id, enterprise_project_id, role_id}]: %v, error: %s",
			agencyId, roleAssignments, err)
	}
	return nil
}

func diffChangeOfProjectRole(oldVal, newVal *schema.Set) (remove, add []string) {
	remove = make([]string, 0)
	add = make([]string, 0)

	oldProjectRolePairs := changeToProjectRolePairs(oldVal)
	newProjectRolePairs := changeToProjectRolePairs(newVal)

	for k := range oldProjectRolePairs {
		if _, ok := newProjectRolePairs[k]; !ok {
			remove = append(remove, k)
		}
	}

	for k := range newProjectRolePairs {
		if _, ok := oldProjectRolePairs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func updateProjectRoles(d *schema.ResourceData, iamClient, identityClient *golangsdk.ServiceClient,
	allRoleIds map[string]string, domainId, agencyId string) error {
	o, n := d.GetChange("project_role")
	deleteProjectRoles, addProjectRoles := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))

	if err := detachProjectRoles(iamClient, identityClient, allRoleIds, deleteProjectRoles, domainId, agencyId); err != nil {
		return err
	}

	//nolint:revive
	if err := attachProjectRoles(iamClient, identityClient, allRoleIds, addProjectRoles, domainId, agencyId); err != nil {
		return err
	}

	return nil
}

func updateDomainRoles(d *schema.ResourceData, iamClient *golangsdk.ServiceClient,
	allRoleIds map[string]string, domainId, agencyId string) error {
	o, n := d.GetChange("domain_roles")
	oldr := o.(*schema.Set)
	newr := n.(*schema.Set)

	detachRoles := utils.ExpandToStringListBySet(oldr.Difference(newr))
	if err := detachDomainRoles(iamClient, allRoleIds, detachRoles, domainId, agencyId); err != nil {
		return err
	}

	attachRoles := utils.ExpandToStringListBySet(newr.Difference(oldr))
	//nolint:revive
	if err := attachDomainRoles(iamClient, allRoleIds, attachRoles, domainId, agencyId); err != nil {
		return err
	}

	return nil
}

func updateAllResourcesRoles(d *schema.ResourceData, iamClient *golangsdk.ServiceClient,
	allRoleIds map[string]string, domainId, agencyId string) error {
	o, n := d.GetChange("all_resources_roles")
	oldr := o.(*schema.Set)
	newr := n.(*schema.Set)

	detachRoles := utils.ExpandToStringListBySet(oldr.Difference(newr))
	if err := detachAllResourcesRoles(iamClient, allRoleIds, detachRoles, domainId, agencyId); err != nil {
		return err
	}

	attachRoles := utils.ExpandToStringListBySet(newr.Difference(oldr))
	return attachAllResourcesRoles(iamClient, allRoleIds, attachRoles, domainId, agencyId)
}

func updateEnterpriseProjectRoles(d *schema.ResourceData, iamClient, epsClient *golangsdk.ServiceClient,
	allRoleIds map[string]string, agencyId string) error {
	o, n := d.GetChange("enterprise_project_roles")
	deleteEnterpriseProjectRoles, addEnterpriseProjectRoles := diffChangeOfEnterpriseProjectRole(o.(*schema.Set), n.(*schema.Set))
	if err := detachEnterpriseProjectRoles(iamClient, epsClient, allRoleIds, deleteEnterpriseProjectRoles, agencyId); err != nil {
		return err
	}

	return attachEnterpriseProjectRoles(iamClient, epsClient, allRoleIds, addEnterpriseProjectRoles, agencyId)
}

func resourceV3AgencyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamV3P0Client, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}
	iamV3Client, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM v3 client: %s", err)
	}

	agencyId := d.Id()
	domainId := cfg.DomainID
	if domainId == "" {
		return diag.Errorf("the parameter 'domain_id' in provider-level configuration must be specified")
	}

	if d.HasChanges("delegated_domain_name", "delegated_service_name", "description", "duration") {
		updateOpts := agency.UpdateOpts{
			Description:     d.Get("description").(string),
			Duration:        buildAgencyDuration(d),
			DelegatedDomain: buildDelegatedDomain(d),
		}

		timeout := d.Timeout(schema.TimeoutUpdate)
		// lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err := agency.Update(iamV3P0Client, agencyId, updateOpts).Extract()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("error updating agency (%s): %s", agencyId, err)
		}
	}

	var allRoles map[string]string
	if d.HasChanges("project_role", "domain_roles", "all_resources_roles", "enterprise_project_roles") {
		allRoles, err = getAllRolesByDomain(iamV3Client, domainId)
		if err != nil {
			return diag.Errorf("error querying the roles: %s", err)
		}
	}

	if d.HasChange("project_role") {
		if err = updateProjectRoles(d, iamV3P0Client, iamV3Client, allRoles, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("domain_roles") {
		if err = updateDomainRoles(d, iamV3P0Client, allRoles, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("all_resources_roles") {
		if err = updateAllResourcesRoles(d, iamV3P0Client, allRoles, domainId, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("enterprise_project_roles") {
		epsClient, err := cfg.EnterpriseProjectClient(region)
		if err != nil {
			return diag.Errorf("error creating EPS client: %s", err)
		}
		if err = updateEnterpriseProjectRoles(d, iamV3P0Client, epsClient, allRoles, agencyId); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceV3AgencyRead(ctx, d, meta)
}

func resourceV3AgencyDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM v3.0 client: %s", err)
	}

	rID := d.Id()
	timeout := d.Timeout(schema.TimeoutDelete)
	// lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := agency.Delete(client, rID).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting agency")
	}

	return nil
}
