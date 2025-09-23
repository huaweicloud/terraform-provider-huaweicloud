package iam

import (
	"bytes"
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

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

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
// @API IAM GET /v3.0/OS-AGENCY/projects/{projectID}/agencies/{agency_id}/roles
// @API IAM PUT /v3.0/OS-AGENCY/projects/{projectID}/agencies/{agency_id}/roles/{role_id}
// @API IAM DELETE /v3.0/OS-AGENCY/projects/{projectID}/agencies/{agency_id}/roles/{role_id}
// @API IAM GET /v3/projects
// @API IAM GET /v3/roles
func ResourceIAMAgencyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMAgencyV3Create,
		ReadContext:   resourceIAMAgencyV3Read,
		UpdateContext: resourceIAMAgencyV3Update,
		DeleteContext: resourceIAMAgencyV3Delete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delegated_domain_name": {
				Type:         schema.TypeString,
				Optional:     true,
				ExactlyOneOf: []string{"delegated_service_name"},
				Description:  "schema: Required",
			},
			"delegated_service_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Internal",
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:             schema.TypeString,
				Optional:         true,
				Default:          "FOREVER",
				DiffSuppressFunc: suppressDurationDiffs,
			},
			"project_role": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"project_role", "domain_roles", "all_resources_roles"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
					},
				},
				Set: resourceIAMAgencyProRoleHash,
			},
			"domain_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},
			"all_resources_roles": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
			},

			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// If `duration` is set to "ONEDAY" in the configuration, it will be set to "1" in the state.
// they have the same meaning, so we should suppress the difference.
//
//nolint:revive // keep same with the definition of SchemaDiffSuppressFunc
func suppressDurationDiffs(k, oldValue, newValue string, d *schema.ResourceData) bool {
	if oldValue == "ONEDAY" {
		oldValue = "1"
	}
	if newValue == "ONEDAY" {
		newValue = "1"
	}

	return oldValue == newValue
}

func resourceIAMAgencyProRoleHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})
	buf.WriteString(fmt.Sprintf("%s-", m["project"].(string)))

	r := m["roles"].(*schema.Set).List()
	s := make([]string, len(r))
	for i, item := range r {
		s[i] = item.(string)
	}
	buf.WriteString(strings.Join(s, "-"))

	return hashcode.String(buf.String())
}

func getProjectIDByName(client *golangsdk.ServiceClient, domainID, name string) (string, error) {
	opts := projects.ListOpts{
		DomainID: domainID,
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

func getAllProjectsOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	opts := projects.ListOpts{
		DomainID: domainID,
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

func listRolesOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	opts := roles.ListOpts{
		DomainID: domainID,
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
			log.Printf("[WARN] role %s without displayname", item.Name)
		}
	}

	return r, nil
}

func getAllRolesOfDomain(client *golangsdk.ServiceClient, domainID string) (map[string]string, error) {
	systemRoles, err := listRolesOfDomain(client, "")
	if err != nil {
		return nil, fmt.Errorf("error listing system-defined roles: %s", err)
	}

	customRoles, err := listRolesOfDomain(client, domainID)
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
	if v, ok := d.GetOk("delegated_domain_name"); ok {
		return v.(string)
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

func resourceIAMAgencyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	identityClient, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating identity client: %s", err)
	}

	domainID := cfg.DomainID
	if domainID == "" {
		return diag.Errorf("the domain_id must be specified in the provider configuration")
	}

	opts := agency.CreateOpts{
		DomainID:        domainID,
		Name:            d.Get("name").(string),
		Description:     d.Get("description").(string),
		Duration:        buildAgencyDuration(d),
		DelegatedDomain: buildDelegatedDomain(d),
	}

	log.Printf("[DEBUG] create IAM agency Options: %#v", opts)
	a, err := agency.Create(iamClient, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating IAM agency: %s", err)
	}

	agencyID := a.ID
	d.SetId(agencyID)

	// get all of the role IDs, include system-defined roles and custom roles
	allRoleIDs, err := getAllRolesOfDomain(identityClient, domainID)
	if err != nil {
		return diag.FromErr(err)
	}

	rawRoles := d.Get("project_role").(*schema.Set)
	pRoles := buildProjectRoles(rawRoles)
	if err := attachProjectRoles(iamClient, identityClient, allRoleIDs, pRoles, domainID, agencyID); err != nil {
		return diag.FromErr(err)
	}

	domainRoles := utils.ExpandToStringListBySet(d.Get("domain_roles").(*schema.Set))
	if err := attachDomainRoles(iamClient, allRoleIDs, domainRoles, domainID, agencyID); err != nil {
		return diag.FromErr(err)
	}

	inheritedRoles := utils.ExpandToStringListBySet(d.Get("all_resources_roles").(*schema.Set))
	if err := attachAllResourcesRoles(iamClient, allRoleIDs, inheritedRoles, domainID, agencyID); err != nil {
		return diag.FromErr(err)
	}

	return resourceIAMAgencyV3Read(ctx, d, meta)
}

// the value can be "FOREVER" or the period in hour
// we should convert the period in day
func normalizeAgencyDuration(dura interface{}) string {
	var result string
	switch v := dura.(type) {
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

func resourceIAMAgencyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	identityClient, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating identity client: %s", err)
	}

	agencyID := d.Id()
	var a *agency.Agency

	a, err = agency.Get(iamClient, agencyID).Extract()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); !ok || !d.IsNewResource() {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM agency")
		}

		// if got 404 error in new resource, wait 10 seconds and try again
		// lintignore:R018
		time.Sleep(10 * time.Second)
		a, err = agency.Get(iamClient, agencyID).Extract()
		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving IAM agency")
		}
	}

	log.Printf("[DEBUG] retrieved IAM agency %s: %#v", agencyID, a)
	mErr := multierror.Append(nil,
		d.Set("name", a.Name),
		d.Set("description", a.Description),
		d.Set("expire_time", a.ExpireTime),
		d.Set("create_time", a.CreateTime),
		d.Set("duration", normalizeAgencyDuration(a.Duration)),
	)

	match, _ := regexp.MatchString("^op_svc_[A-Za-z]+$", a.DelegatedDomainName)
	if match {
		mErr = multierror.Append(mErr, d.Set("delegated_service_name", a.DelegatedDomainName))
	} else {
		mErr = multierror.Append(mErr, d.Set("delegated_domain_name", a.DelegatedDomainName))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity agency fields: %s", err)
	}

	allProjects, err := getAllProjectsOfDomain(identityClient, a.DomainID)
	if err != nil {
		return diag.Errorf("error querying the projects of domain: %s", err)
	}

	prs := schema.Set{F: resourceIAMAgencyProRoleHash}
	for pn, pid := range allProjects {
		// MOS is a special project, not visible to the user
		if pn == "MOS" {
			continue
		}

		// the provider will query the roles in all projects, but the API rate limit threshold is 10 times per second.
		// so we should wait for some time to avoid exceeding the rate limit.
		// lintignore:R018
		time.Sleep(200 * time.Millisecond)

		allRoles, err := agency.ListRolesAttachedOnProject(iamClient, agencyID, pid).ExtractRoles()
		if err != nil && !utils.IsResourceNotFound(err) {
			log.Printf("[ERROR] error querying the roles attached on project(%s): %s", pn, err)
			continue
		}
		if len(allRoles) == 0 {
			continue
		}
		v := schema.Set{F: schema.HashString}
		for _, role := range allRoles {
			v.Add(role.DisplayName)
		}
		prs.Add(map[string]interface{}{
			"project": pn,
			"roles":   &v,
		})
	}
	err = d.Set("project_role", &prs)
	if err != nil {
		log.Printf("[ERROR] Set project_role failed: %s", err)
	}

	allDomainRoles, err := agency.ListRolesAttachedOnDomain(iamClient, agencyID, a.DomainID).ExtractRoles()
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
			log.Printf("[ERROR] set domain_roles failed: %s", err)
		}
	}

	// Unable to fetch all_resources_roles because the API response does not include `display_name` field
	// https://support.huaweicloud.com/api-iam/iam_12_0014.html
	return nil
}

func buildProjectRoles(prs *schema.Set) []string {
	addprs := changeToPRPair(prs)
	pRoles := make([]string, 0, len(addprs))
	for key := range addprs {
		pRoles = append(pRoles, key)
	}
	return pRoles
}

func changeToPRPair(prs *schema.Set) map[string]bool {
	r := make(map[string]bool)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})

		pn := pr["project"].(string)
		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			key := fmt.Sprintf("%s|%s", pn, role.(string))
			r[key] = true
		}
	}
	return r
}

func diffChangeOfProjectRole(oldVal, newVal *schema.Set) (remove, add []string) {
	remove = make([]string, 0)
	add = make([]string, 0)

	oldprs := changeToPRPair(oldVal)
	newprs := changeToPRPair(newVal)

	for k := range oldprs {
		if _, ok := newprs[k]; !ok {
			remove = append(remove, k)
		}
	}

	for k := range newprs {
		if _, ok := oldprs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func attachProjectRoles(iamClient, identityClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	pRoles []string, domainID, agencyID string) error {
	if len(pRoles) > 0 {
		log.Printf("[DEBUG] attaching roles %v in project scope to agency %s", pRoles, agencyID)
	}

	for _, v := range pRoles {
		pr := strings.Split(v, "|")
		if len(pr) != 2 {
			return fmt.Errorf("error parsing project role from %s: invalid format", v)
		}

		pid, err := getProjectIDByName(identityClient, domainID, pr[0])
		if err != nil {
			return fmt.Errorf("the project(%s) is not exist", pr[0])
		}
		rid, ok := allRoleIDs[pr[1]]
		if !ok {
			return fmt.Errorf("the role(%s) to be attached is not exist", pr[1])
		}

		err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role(%s) by project(%s) to agency(%s): %s",
				rid, pid, agencyID, err)
		}
	}

	return nil
}

func detachProjectRoles(iamClient, identityClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	pRoles []string, domainID, agencyID string) error {
	if len(pRoles) > 0 {
		log.Printf("[DEBUG] detaching roles %v in project scope from agency %s", pRoles, agencyID)
	}

	for _, v := range pRoles {
		pr := strings.Split(v, "|")
		if len(pr) != 2 {
			return fmt.Errorf("error parsing project role from %s: invalid format", v)
		}

		pid, err := getProjectIDByName(identityClient, domainID, pr[0])
		if err != nil {
			return fmt.Errorf("the project(%s) is not exist", pr[0])
		}

		rid, ok := allRoleIDs[pr[1]]
		if !ok {
			log.Printf("[WARN] the role(%s) to be detached is not exist", pr[1])
			continue
		}

		err = agency.DetachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role(%s) by project(%s) from agency(%s): %s",
				rid, pid, agencyID, err)
		}
	}

	return nil
}

func attachDomainRoles(iamClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	roleNames []string, domainID, agencyID string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] attaching roles %v in domain scope to agency %s", roleNames, agencyID)
	}

	for _, r := range roleNames {
		rid, ok := allRoleIDs[r]
		if !ok {
			return fmt.Errorf("the role(%s) to be attached is not exist", r)
		}

		err := agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role(%s) by domain(%s) to agency(%s): %s",
				rid, domainID, agencyID, err)
		}
	}

	return nil
}

func detachDomainRoles(iamClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	roleNames []string, domainID, agencyID string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in domain scope from agency %s", roleNames, agencyID)
	}

	for _, r := range roleNames {
		rid, ok := allRoleIDs[r]
		if !ok {
			log.Printf("[WARN] the role(%s) to be detached is not exist", r)
			continue
		}

		err := agency.DetachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role(%s) by domain(%s) from agency(%s): %s",
				rid, domainID, agencyID, err)
		}
	}

	return nil
}

func attachAllResourcesRoles(iamClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	roleNames []string, domainID, agencyID string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] attaching roles %v in all resources to agency %s", roleNames, agencyID)
	}

	for _, r := range roleNames {
		rid, ok := allRoleIDs[r]
		if !ok {
			return fmt.Errorf("the role(%s) to be attached is not exist", r)
		}

		err := agency.AttachAllResources(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role(%s) in all resources to agency(%s): %s",
				r, agencyID, err)
		}
	}

	return nil
}

func detachAllResourcesRoles(iamClient *golangsdk.ServiceClient, allRoleIDs map[string]string,
	roleNames []string, domainID, agencyID string) error {
	if len(roleNames) > 0 {
		log.Printf("[DEBUG] detaching roles %v in all resources from agency %s", roleNames, agencyID)
	}

	for _, r := range roleNames {
		rid, ok := allRoleIDs[r]
		if !ok {
			return fmt.Errorf("the role(%s) to be detached is not exist", r)
		}

		err := agency.DetachAllResources(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("error detaching role(%s) in all resources from agency(%s): %s",
				r, agencyID, err)
		}
	}

	return nil
}

func updateProjectRoles(d *schema.ResourceData, iamClient, identityClient *golangsdk.ServiceClient,
	allRoleIDs map[string]string, domainID, agencyID string) error {
	o, n := d.GetChange("project_role")
	deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))

	if err := detachProjectRoles(iamClient, identityClient, allRoleIDs, deleteprs, domainID, agencyID); err != nil {
		return err
	}

	//nolint:revive
	if err := attachProjectRoles(iamClient, identityClient, allRoleIDs, addprs, domainID, agencyID); err != nil {
		return err
	}

	return nil
}

func updateDomainRoles(d *schema.ResourceData, iamClient *golangsdk.ServiceClient,
	allRoleIDs map[string]string, domainID, agencyID string) error {
	o, n := d.GetChange("domain_roles")
	oldr := o.(*schema.Set)
	newr := n.(*schema.Set)

	detachRoles := utils.ExpandToStringListBySet(oldr.Difference(newr))
	if err := detachDomainRoles(iamClient, allRoleIDs, detachRoles, domainID, agencyID); err != nil {
		return err
	}

	attachRoles := utils.ExpandToStringListBySet(newr.Difference(oldr))
	//nolint:revive
	if err := attachDomainRoles(iamClient, allRoleIDs, attachRoles, domainID, agencyID); err != nil {
		return err
	}

	return nil
}

func updateAllResourcesRoles(d *schema.ResourceData, iamClient *golangsdk.ServiceClient,
	allRoleIDs map[string]string, domainID, agencyID string) error {
	o, n := d.GetChange("all_resources_roles")
	oldr := o.(*schema.Set)
	newr := n.(*schema.Set)

	detachRoles := utils.ExpandToStringListBySet(oldr.Difference(newr))
	if err := detachAllResourcesRoles(iamClient, allRoleIDs, detachRoles, domainID, agencyID); err != nil {
		return err
	}

	attachRoles := utils.ExpandToStringListBySet(newr.Difference(oldr))
	//nolint:revive
	if err := attachAllResourcesRoles(iamClient, allRoleIDs, attachRoles, domainID, agencyID); err != nil {
		return err
	}

	return nil
}

func resourceIAMAgencyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	iamClient, err := cfg.IAMV3Client(region)
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}
	identityClient, err := cfg.IdentityV3Client(region)
	if err != nil {
		return diag.Errorf("error creating identity client: %s", err)
	}

	agencyID := d.Id()
	domainID := cfg.DomainID
	if domainID == "" {
		return diag.Errorf("the domain_id must be specified in the provider configuration")
	}

	if d.HasChanges("delegated_domain_name", "delegated_service_name", "description", "duration") {
		updateOpts := agency.UpdateOpts{
			Description:     d.Get("description").(string),
			Duration:        buildAgencyDuration(d),
			DelegatedDomain: buildDelegatedDomain(d),
		}

		log.Printf("[DEBUG] updating IAM agency %s with options: %#v", agencyID, updateOpts)
		timeout := d.Timeout(schema.TimeoutUpdate)
		//lintignore:R006
		err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
			_, err := agency.Update(iamClient, agencyID, updateOpts).Extract()
			if err != nil {
				return common.CheckForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return diag.Errorf("error updating IAM agency %s: %s", agencyID, err)
		}
	}

	var allRoles map[string]string
	if d.HasChanges("project_role", "domain_roles", "all_resources_roles") {
		allRoles, err = getAllRolesOfDomain(identityClient, domainID)
		if err != nil {
			return diag.Errorf("error querying the roles: %s", err)
		}
	}

	if d.HasChange("project_role") {
		if err = updateProjectRoles(d, iamClient, identityClient, allRoles, domainID, agencyID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("domain_roles") {
		if err = updateDomainRoles(d, iamClient, allRoles, domainID, agencyID); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("all_resources_roles") {
		if err = updateAllResourcesRoles(d, iamClient, allRoles, domainID, agencyID); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceIAMAgencyV3Read(ctx, d, meta)
}

func resourceIAMAgencyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	rID := d.Id()
	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.RetryContext(ctx, timeout, func() *resource.RetryError {
		err := agency.Delete(iamClient, rID).ExtractErr()
		if err != nil {
			return common.CheckForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if utils.IsResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable IAM agency: %s", rID)
			return nil
		}
		return diag.Errorf("error deleting IAM agency %s: %s", rID, err)
	}

	return nil
}
