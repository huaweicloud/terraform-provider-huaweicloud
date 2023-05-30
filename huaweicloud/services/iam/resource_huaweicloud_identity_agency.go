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
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func ResourceIAMAgencyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIAMAgencyV3Create,
		ReadContext:   resourceIAMAgencyV3Read,
		UpdateContext: resourceIAMAgencyV3Update,
		DeleteContext: resourceIAMAgencyV3Delete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
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
				ValidateFunc: validation.StringDoesNotMatch(regexp.MustCompile("^op_svc_[A-Za-z]+"),
					"the value can not start with op_svc_, use `delegated_service_name` instead"),
				Description: "schema: Required",
			},
			"delegated_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^op_svc_[A-Za-z]+"),
					"the value must start with op_svc_."),
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
				AtLeastOneOf: []string{"project_role", "domain_roles"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"project": {
							Type:     schema.TypeString,
							Required: true,
						},
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 25,
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
				MaxItems: 25,
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
	allPages, err := roles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("failed to query roles: %s", err)
	}

	allItems, err := roles.ExtractRoles(allPages)
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

	prs := d.Get("project_role").(*schema.Set).List()
	for _, v := range prs {
		pr := v.(map[string]interface{})
		pname := pr["project"].(string)
		pid, err := getProjectIDByName(identityClient, domainID, pname)
		if err != nil {
			return diag.FromErr(err)
		}

		rs := pr["roles"].(*schema.Set).List()
		for _, role := range rs {
			r := role.(string)
			rid, ok := allRoleIDs[r]
			if !ok {
				return diag.Errorf("the project role(%s) is not exist", r)
			}

			err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return diag.Errorf("error attaching role(%s) by project(%s) to agency(%s): %s",
					rid, pid, agencyID, err)
			}
		}
	}

	drs := d.Get("domain_roles").(*schema.Set)
	for _, role := range drs.List() {
		r := role.(string)
		rid, ok := allRoleIDs[r]
		if !ok {
			return diag.Errorf("the domain role(%s) is not exist", r)
		}

		err = agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return diag.Errorf("error attaching role(%s) by domain(%s) to agency(%s): %s",
				rid, domainID, agencyID, err)
		}
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
	a, err := agency.Get(iamClient, agencyID).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM agency")
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
			log.Printf("[ERROR] Set domain_roles failed: %s", err)
		}
	}

	return nil
}

func changeToPRPair(prs *schema.Set) (r map[string]bool) {
	r = make(map[string]bool)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})

		pn := pr["project"].(string)
		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			key := fmt.Sprintf("%s|%s", pn, role.(string))
			r[key] = true
		}
	}
	return
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

func updateProjectRoles(d *schema.ResourceData, iamClient, identityClient *golangsdk.ServiceClient,
	allRoleIDs map[string]string, domainID, agencyID string) error {
	o, n := d.GetChange("project_role")
	deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))
	for _, v := range deleteprs {
		pr := strings.Split(v, "|")

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

	for _, v := range addprs {
		pr := strings.Split(v, "|")
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

func updateDomainRoles(d *schema.ResourceData, iamClient *golangsdk.ServiceClient,
	allRoleIDs map[string]string, domainID, agencyID string) error {
	o, n := d.GetChange("domain_roles")
	oldr := o.(*schema.Set)
	newr := n.(*schema.Set)

	for _, r := range oldr.Difference(newr).List() {
		rid, ok := allRoleIDs[r.(string)]
		if !ok {
			log.Printf("[WARN] the role(%s) to be detached is not exist", r.(string))
			continue
		}

		err := agency.DetachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmt.Errorf("error detaching role(%s) by domain(%s) from agency(%s): %s",
				rid, domainID, agencyID, err)
		}
	}

	for _, r := range newr.Difference(oldr).List() {
		rid, ok := allRoleIDs[r.(string)]
		if !ok {
			return fmt.Errorf("the role(%s) to be attached is not exist", r.(string))
		}

		err := agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("error attaching role(%s) by domain(%s) to agency(%s): %s",
				rid, domainID, agencyID, err)
		}
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
	if d.HasChanges("project_role", "domain_roles") {
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
