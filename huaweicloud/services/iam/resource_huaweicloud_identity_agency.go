package iam

import (
	"bytes"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3/agency"
	"github.com/chnsz/golangsdk/openstack/identity/v3/projects"
	"github.com/chnsz/golangsdk/openstack/identity/v3/roles"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
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
			},
			"delegated_service_name": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^op_svc_[A-Za-z]+"),
					"the value must start with op_svc_."),
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"project_role": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"project_role", "domain_roles"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"roles": {
							Type:     schema.TypeSet,
							Required: true,
							MaxItems: 25,
							Elem:     &schema.Schema{Type: schema.TypeString},
							Set:      schema.HashString,
						},
						"project": {
							Type:     schema.TypeString,
							Required: true,
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

			"duration": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "FOREVER",
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

func getProjectIDOfDomain(client *golangsdk.ServiceClient, domainID, name string) (string, error) {
	opts := projects.ListOpts{
		DomainID: domainID,
		Name:     name,
	}
	allPages, err := projects.List(client, &opts).AllPages()
	if err != nil {
		return "", fmtp.Errorf("List projects failed, err=%s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return "", fmtp.Errorf("Extract projects failed, err=%s", err)
	}

	if len(all) == 0 {
		return "", fmtp.Errorf("Wrong name or no access to the project: %s", name)
	}

	item := all[0]
	return item.ID, nil
}

func listProjectsOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	opts := projects.ListOpts{
		DomainID: domainID,
	}
	allPages, err := projects.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmtp.Errorf("List projects failed, err=%s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return nil, fmtp.Errorf("Extract projects failed, err=%s", err)
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		r[item.Name] = item.ID
	}
	logp.Printf("[TRACE] projects = %#v\n", r)
	return r, nil
}

func listRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	opts := roles.ListOpts{
		DomainID: domainID,
	}
	allPages, err := roles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmtp.Errorf("List roles failed, err=%s", err)
	}

	all, err := roles.ExtractRoles(allPages)
	if err != nil {
		return nil, fmtp.Errorf("Extract roles failed, err=%s", err)
	}
	if len(all) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		if name := item.DisplayName; name != "" {
			r[name] = item.ID
		} else {
			logp.Printf("[WARN] role %s without displayname", item.Name)
		}
	}
	logp.Printf("[TRACE] list roles = %#v, len=%d\n", r, len(r))
	return r, nil
}

func getAllRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	roles, err := listRolesOfDomain("", client)
	if err != nil {
		return nil, fmtp.Errorf("Error listing system-defined roles, err=%s", err)
	}

	customRoles, err := listRolesOfDomain(domainID, client)
	if err != nil {
		return nil, fmtp.Errorf("Error listing domain's custom roles, err=%s", err)
	}

	if roles == nil {
		return customRoles, nil
	}

	if customRoles == nil {
		return roles, nil
	}

	for k, v := range customRoles {
		roles[k] = v
	}
	return roles, nil
}

func changeToPRPair(prs *schema.Set) (r map[string]bool) {
	r = make(map[string]bool)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})

		pn := pr["project"].(string)
		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			r[pn+"|"+role.(string)] = true
		}
	}
	return
}

func diffChangeOfProjectRole(old, newv *schema.Set) (delete, add []string) {
	delete = make([]string, 0)
	add = make([]string, 0)

	oldprs := changeToPRPair(old)
	newprs := changeToPRPair(newv)

	for k := range oldprs {
		if _, ok := newprs[k]; !ok {
			delete = append(delete, k)
		}
	}

	for k := range newprs {
		if _, ok := oldprs[k]; !ok {
			add = append(add, k)
		}
	}
	return
}

func resourceIAMAgencyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud IAM client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	domainID := config.DomainID
	if domainID == "" {
		return fmtp.DiagErrorf("the domain_id must be specified in the provider configuration")
	}

	opts := agency.CreateOpts{
		Name:        d.Get("name").(string),
		DomainID:    domainID,
		Description: d.Get("description").(string),
		Duration:    d.Get("duration").(string),
	}
	if v, ok := d.GetOk("delegated_domain_name"); ok {
		opts.DelegatedDomain = v.(string)
	} else {
		opts.DelegatedDomain = d.Get("delegated_service_name").(string)
	}
	logp.Printf("[DEBUG] Create IAM-Agency Options: %#v", opts)

	a, err := agency.Create(iamClient, opts).Extract()
	if err != nil {
		return fmtp.DiagErrorf("Error creating IAM-Agency: %s", err)
	}

	agencyID := a.ID
	d.SetId(agencyID)

	roles, err := getAllRolesOfDomain(domainID, identityClient)
	if err != nil {
		return fmtp.DiagErrorf("Error querying the roles, err=%s", err)
	}

	prs := d.Get("project_role").(*schema.Set)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})
		pname := pr["project"].(string)
		pid, err := getProjectIDOfDomain(identityClient, domainID, pname)
		if err != nil {
			return fmtp.DiagErrorf("The project(%s) is not exist", pname)
		}

		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			r := role.(string)
			rid, ok := roles[r]
			if !ok {
				return fmtp.DiagErrorf("The project role(%s) is not exist", r)
			}

			err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return fmtp.DiagErrorf("Error attaching role(%s) by project(%s) to agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}
	}

	drs := d.Get("domain_roles").(*schema.Set)
	for _, role := range drs.List() {
		r := role.(string)
		rid, ok := roles[r]
		if !ok {
			return fmtp.DiagErrorf("The domain role(%s) is not exist", r)
		}

		err = agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmtp.DiagErrorf("Error attaching role(%s) by domain(%s) to agency(%s), err=%s",
				rid, domainID, agencyID, err)
		}
	}

	return resourceIAMAgencyV3Read(ctx, d, meta)
}

func resourceIAMAgencyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	a, err := agency.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "IAM-Agency")
	}
	logp.Printf("[DEBUG] Retrieved IAM-Agency %s: %#v", d.Id(), a)

	mErr := multierror.Append(nil,
		d.Set("name", a.Name),
		d.Set("description", a.Description),
		d.Set("expire_time", a.ExpireTime),
		d.Set("create_time", a.CreateTime),
	)

	if a.Duration != "" {
		mErr = multierror.Append(mErr, d.Set("duration", a.Duration))
	} else {
		mErr = multierror.Append(mErr, d.Set("duration", "FOREVER"))
	}

	if ok, err := regexp.MatchString("^op_svc_[A-Za-z]+$", a.DelegatedDomainName); err != nil {
		logp.Printf("[ERROR] Regexp error, err= %s", err)
	} else if ok {
		mErr = multierror.Append(mErr, d.Set("delegated_service_name", a.DelegatedDomainName))
	} else {
		mErr = multierror.Append(mErr, d.Set("delegated_domain_name", a.DelegatedDomainName))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("error setting identity agency fields: %s", err)
	}

	projects, err := listProjectsOfDomain(a.DomainID, identityClient)
	if err != nil {
		return fmtp.DiagErrorf("Error querying the projects, err=%s", err)
	}
	agencyID := d.Id()
	prs := schema.Set{F: resourceIAMAgencyProRoleHash}
	for pn, pid := range projects {
		roles, err := agency.ListRolesAttachedOnProject(iamClient, agencyID, pid).ExtractRoles()
		if err != nil && !utils.IsResourceNotFound(err) {
			return fmtp.DiagErrorf("Error querying the roles attached on project(%s), err=%s", pn, err)
		}
		if len(roles) == 0 {
			continue
		}
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.DisplayName)
		}
		prs.Add(map[string]interface{}{
			"project": pn,
			"roles":   &v,
		})
	}
	err = d.Set("project_role", &prs)
	if err != nil {
		logp.Printf("[ERROR]Set project_role failed, err=%s", err)
	}

	roles, err := agency.ListRolesAttachedOnDomain(iamClient, agencyID, a.DomainID).ExtractRoles()
	if err != nil && !utils.IsResourceNotFound(err) {
		return fmtp.DiagErrorf("Error querying the roles attached on domain, err=%s", err)
	}
	if len(roles) != 0 {
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.DisplayName)
		}
		err = d.Set("domain_roles", &v)
		if err != nil {
			logp.Printf("[ERROR]Set domain_roles failed, err=%s", err)
		}
	}

	return nil
}

func resourceIAMAgencyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud identity client: %s", err)
	}

	agencyID := d.Id()
	domainID := config.DomainID
	if domainID == "" {
		return fmtp.DiagErrorf("the domain_id must be specified in the provider configuration")
	}

	if d.HasChanges("delegated_domain_name", "delegated_service_name", "description", "duration") {
		updateOpts := agency.UpdateOpts{
			Description: d.Get("description").(string),
			Duration:    d.Get("duration").(string),
		}

		if v, ok := d.GetOk("delegated_domain_name"); ok {
			updateOpts.DelegatedDomain = v.(string)
		} else {
			updateOpts.DelegatedDomain = d.Get("delegated_service_name").(string)
		}

		logp.Printf("[DEBUG] Updating IAM-Agency %s with options: %#v", agencyID, updateOpts)
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
			return fmtp.DiagErrorf("Error updating IAM-Agency %s: %s", agencyID, err)
		}
	}

	var roles map[string]string
	if d.HasChanges("project_role", "domain_roles") {
		roles, err = getAllRolesOfDomain(domainID, identityClient)
		if err != nil {
			return fmtp.DiagErrorf("Error querying the roles, err=%s", err)
		}
	}

	if d.HasChange("project_role") {
		o, n := d.GetChange("project_role")
		deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))
		for _, v := range deleteprs {
			pr := strings.Split(v, "|")
			pid, err := getProjectIDOfDomain(identityClient, domainID, pr[0])
			if err != nil {
				return fmtp.DiagErrorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmtp.DiagErrorf("The role(%s) is not exist", pr[1])
			}

			err = agency.DetachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil && !utils.IsResourceNotFound(err) {
				return fmtp.DiagErrorf("Error detaching role(%s) by project{%s} from agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}

		for _, v := range addprs {
			pr := strings.Split(v, "|")
			pid, err := getProjectIDOfDomain(identityClient, domainID, pr[0])
			if err != nil {
				return fmtp.DiagErrorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmtp.DiagErrorf("The role(%s) is not exist", pr[1])
			}

			err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return fmtp.DiagErrorf("Error attaching role(%s) by project{%s} to agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}
	}

	if d.HasChange("domain_roles") {
		o, n := d.GetChange("domain_roles")
		oldr := o.(*schema.Set)
		newr := n.(*schema.Set)

		for _, r := range oldr.Difference(newr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmtp.DiagErrorf("The role(%s) is not exist", r.(string))
			}

			err = agency.DetachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
			if err != nil && !utils.IsResourceNotFound(err) {
				return fmtp.DiagErrorf("Error detaching role(%s) by domain{%s} from agency(%s), err=%s",
					rid, domainID, agencyID, err)
			}
		}

		for _, r := range newr.Difference(oldr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmtp.DiagErrorf("The role(%s) is not exist", r.(string))
			}

			err = agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
			if err != nil {
				return fmtp.DiagErrorf("Error attaching role(%s) by domain{%s} to agency(%s), err=%s",
					rid, domainID, agencyID, err)
			}
		}
	}
	return resourceIAMAgencyV3Read(ctx, d, meta)
}

func resourceIAMAgencyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	iamClient, err := config.IAMV3Client(config.GetRegion(d))
	if err != nil {
		return fmtp.DiagErrorf("Error creating HuaweiCloud client: %s", err)
	}

	rID := d.Id()
	logp.Printf("[DEBUG] Deleting IAM-Agency %s", rID)

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
			logp.Printf("[INFO] deleting an unavailable IAM-Agency: %s", rID)
			return nil
		}
		return fmtp.DiagErrorf("Error deleting IAM-Agency %s: %s", rID, err)
	}

	return nil
}
