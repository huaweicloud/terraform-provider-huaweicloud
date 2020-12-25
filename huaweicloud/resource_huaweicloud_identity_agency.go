package huaweicloud

import (
	"bytes"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/hashcode"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/agency"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/domains"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/projects"
	"github.com/huaweicloud/golangsdk/openstack/identity/v3/roles"
)

func resourceIAMAgencyV3() *schema.Resource {
	return &schema.Resource{
		Create: resourceIAMAgencyV3Create,
		Read:   resourceIAMAgencyV3Read,
		Update: resourceIAMAgencyV3Update,
		Delete: resourceIAMAgencyV3Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"delegated_domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expire_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
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
							MinItems: 1,
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
				MinItems: 1,
				MaxItems: 25,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Set:      schema.HashString,
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

func listProjectsOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	opts := projects.ListOpts{
		DomainID: domainID,
	}
	allPages, err := projects.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("List projects failed, err=%s", err)
	}

	all, err := projects.ExtractProjects(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract projects failed, err=%s", err)
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		r[item.Name] = item.ID
	}
	log.Printf("[TRACE] projects = %#v\n", r)
	return r, nil
}

func listRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	opts := roles.ListOpts{
		DomainID: domainID,
	}
	allPages, err := roles.List(client, &opts).AllPages()
	if err != nil {
		return nil, fmt.Errorf("List roles failed, err=%s", err)
	}

	all, err := roles.ExtractRoles(allPages)
	if err != nil {
		return nil, fmt.Errorf("Extract roles failed, err=%s", err)
	}
	if len(all) == 0 {
		return nil, nil
	}

	r := make(map[string]string, len(all))
	for _, item := range all {
		dn, ok := item.Extra["display_name"].(string)
		if ok {
			r[dn] = item.ID
		} else {
			log.Printf("[DEBUG] Can not retrieve role:%#v", item)
		}
	}
	log.Printf("[TRACE] list roles = %#v, len=%d\n", r, len(r))
	return r, nil
}

func allRolesOfDomain(domainID string, client *golangsdk.ServiceClient) (map[string]string, error) {
	roles, err := listRolesOfDomain("", client)
	if err != nil {
		return nil, fmt.Errorf("Error listing global roles, err=%s", err)
	}

	customRoles, err := listRolesOfDomain(domainID, client)
	if err != nil {
		return nil, fmt.Errorf("Error listing domain's custom roles, err=%s", err)
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

func getDomainID(config *Config, client *golangsdk.ServiceClient) (string, error) {
	if config.DomainID != "" {
		return config.DomainID, nil
	}

	name := config.DomainName
	if name == "" {
		return "", fmt.Errorf("The required domain name was missed")
	}

	old := client.ResourceBase
	defer func() { client.ResourceBase = old }()
	// ResourceBase: https://iam.{CLOUD}/v3/auth/
	client.ResourceBase = old + "auth/"

	opts := domains.ListOpts{
		Name: name,
	}
	allPages, err := domains.List(client, &opts).AllPages()
	if err != nil {
		return "", fmt.Errorf("List domains failed, err=%s", err)
	}

	all, err := domains.ExtractDomains(allPages)
	if err != nil {
		return "", fmt.Errorf("Extract domains failed, err=%s", err)
	}

	count := len(all)
	switch count {
	case 0:
		err := &golangsdk.ErrResourceNotFound{}
		err.ResourceType = "iam"
		err.Name = name
		return "", err
	case 1:
		domainID := all[0].ID
		config.DomainID = domainID
		return domainID, nil
	default:
		err := &golangsdk.ErrMultipleResourcesFound{}
		err.ResourceType = "iam"
		err.Name = name
		err.Count = count
		return "", err
	}
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

func resourceIAMAgencyV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud IAM client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	domainID, err := getDomainID(config, identityClient)
	if err != nil {
		return fmt.Errorf("Error getting the domain id, err=%s", err)
	}

	opts := agency.CreateOpts{
		Name:            d.Get("name").(string),
		DomainID:        domainID,
		DelegatedDomain: d.Get("delegated_domain_name").(string),
		Description:     d.Get("description").(string),
	}
	log.Printf("[DEBUG] Create IAM-Agency Options: %#v", opts)
	a, err := agency.Create(iamClient, opts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating IAM-Agency: %s", err)
	}

	agencyID := a.ID
	d.SetId(agencyID)

	projects, err := listProjectsOfDomain(domainID, identityClient)
	if err != nil {
		return fmt.Errorf("Error querying the projects, err=%s", err)
	}

	roles, err := allRolesOfDomain(domainID, identityClient)
	if err != nil {
		return fmt.Errorf("Error querying the roles, err=%s", err)
	}

	prs := d.Get("project_role").(*schema.Set)
	for _, v := range prs.List() {
		pr := v.(map[string]interface{})
		pn := pr["project"].(string)
		pid, ok := projects[pn]
		if !ok {
			return fmt.Errorf("The project(%s) is not exist", pn)
		}

		rs := pr["roles"].(*schema.Set)
		for _, role := range rs.List() {
			r := role.(string)
			rid, ok := roles[r]
			if !ok {
				return fmt.Errorf("The project role(%s) is not exist", r)
			}

			err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by project(%s) to agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}
	}

	drs := d.Get("domain_roles").(*schema.Set)
	for _, role := range drs.List() {
		r := role.(string)
		rid, ok := roles[r]
		if !ok {
			return fmt.Errorf("The domain role(%s) is not exist", r)
		}

		err = agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
		if err != nil {
			return fmt.Errorf("Error attaching role(%s) by domain(%s) to agency(%s), err=%s",
				rid, domainID, agencyID, err)
		}
	}

	return resourceIAMAgencyV3Read(d, meta)
}

func resourceIAMAgencyV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	a, err := agency.Get(iamClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "IAM-Agency")
	}
	log.Printf("[DEBUG] Retrieved IAM-Agency %s: %#v", d.Id(), a)

	d.Set("region", GetRegion(d, config))
	d.Set("name", a.Name)
	d.Set("delegated_domain_name", a.DelegatedDomainName)
	d.Set("description", a.Description)
	d.Set("duration", a.Duration)
	d.Set("expire_time", a.ExpireTime)
	d.Set("create_time", a.CreateTime)

	projects, err := listProjectsOfDomain(a.DomainID, identityClient)
	if err != nil {
		return fmt.Errorf("Error querying the projects, err=%s", err)
	}
	agencyID := d.Id()
	prs := schema.Set{F: resourceIAMAgencyProRoleHash}
	for pn, pid := range projects {
		roles, err := agency.ListRolesAttachedOnProject(iamClient, agencyID, pid).ExtractRoles()
		if err != nil && !isResourceNotFound(err) {
			return fmt.Errorf("Error querying the roles attached on project(%s), err=%s", pn, err)
		}
		if len(roles) == 0 {
			continue
		}
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.Extra["display_name"])
		}
		prs.Add(map[string]interface{}{
			"project": pn,
			"roles":   &v,
		})
	}
	err = d.Set("project_role", &prs)
	if err != nil {
		log.Printf("[ERROR]Set project_role failed, err=%s", err)
	}

	roles, err := agency.ListRolesAttachedOnDomain(iamClient, agencyID, a.DomainID).ExtractRoles()
	if err != nil && !isResourceNotFound(err) {
		return fmt.Errorf("Error querying the roles attached on domain, err=%s", err)
	}
	if len(roles) != 0 {
		v := schema.Set{F: schema.HashString}
		for _, role := range roles {
			v.Add(role.Extra["display_name"])
		}
		err = d.Set("domain_roles", &v)
		if err != nil {
			log.Printf("[ERROR]Set domain_roles failed, err=%s", err)
		}
	}

	return nil
}

func resourceIAMAgencyV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud client: %s", err)
	}
	identityClient, err := config.IdentityV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud identity client: %s", err)
	}

	agencyID := d.Id()

	if d.HasChanges("delegated_domain_name", "description") {
		updateOpts := agency.UpdateOpts{
			DelegatedDomain: d.Get("delegated_domain_name").(string),
			Description:     d.Get("description").(string),
		}
		log.Printf("[DEBUG] Updating IAM-Agency %s with options: %#v", agencyID, updateOpts)
		timeout := d.Timeout(schema.TimeoutUpdate)
		//lintignore:R006
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err := agency.Update(iamClient, agencyID, updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("Error updating IAM-Agency %s: %s", agencyID, err)
		}
	}

	domainID := ""
	var roles map[string]string
	if d.HasChanges("project_role", "domain_roles") {
		domainID, err = getDomainID(config, identityClient)
		if err != nil {
			return fmt.Errorf("Error getting the domain id, err=%s", err)
		}

		roles, err = allRolesOfDomain(domainID, identityClient)
		if err != nil {
			return fmt.Errorf("Error querying the roles, err=%s", err)
		}
	}

	if d.HasChange("project_role") {
		projects, err := listProjectsOfDomain(domainID, identityClient)
		if err != nil {
			return fmt.Errorf("Error querying the projects, err=%s", err)
		}

		o, n := d.GetChange("project_role")
		deleteprs, addprs := diffChangeOfProjectRole(o.(*schema.Set), n.(*schema.Set))
		for _, v := range deleteprs {
			pr := strings.Split(v, "|")
			pid, ok := projects[pr[0]]
			if !ok {
				return fmt.Errorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", pr[1])
			}

			err = agency.DetachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil && !isResourceNotFound(err) {
				return fmt.Errorf("Error detaching role(%s) by project{%s} from agency(%s), err=%s",
					rid, pid, agencyID, err)
			}
		}

		for _, v := range addprs {
			pr := strings.Split(v, "|")
			pid, ok := projects[pr[0]]
			if !ok {
				return fmt.Errorf("The project(%s) is not exist", pr[0])
			}
			rid, ok := roles[pr[1]]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", pr[1])
			}

			err = agency.AttachRoleByProject(iamClient, agencyID, pid, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by project{%s} to agency(%s), err=%s",
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
				return fmt.Errorf("The role(%s) is not exist", r.(string))
			}

			err = agency.DetachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
			if err != nil && !isResourceNotFound(err) {
				return fmt.Errorf("Error detaching role(%s) by domain{%s} from agency(%s), err=%s",
					rid, domainID, agencyID, err)
			}
		}

		for _, r := range newr.Difference(oldr).List() {
			rid, ok := roles[r.(string)]
			if !ok {
				return fmt.Errorf("The role(%s) is not exist", r.(string))
			}

			err = agency.AttachRoleByDomain(iamClient, agencyID, domainID, rid).ExtractErr()
			if err != nil {
				return fmt.Errorf("Error attaching role(%s) by domain{%s} to agency(%s), err=%s",
					rid, domainID, agencyID, err)
			}
		}
	}
	return resourceIAMAgencyV3Read(d, meta)
}

func resourceIAMAgencyV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	iamClient, err := config.IAMV3Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud client: %s", err)
	}

	rID := d.Id()
	log.Printf("[DEBUG] Deleting IAM-Agency %s", rID)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := agency.Delete(iamClient, rID).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable IAM-Agency: %s", rID)
			return nil
		}
		return fmt.Errorf("Error deleting IAM-Agency %s: %s", rID, err)
	}

	return nil
}
