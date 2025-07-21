package iam

import (
	"bytes"
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk/openstack/identity/v3.0/acl"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/hashcode"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
func ResourceIdentityACL() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityACLCreate,
		ReadContext:   resourceIdentityACLRead,
		UpdateContext: resourceIdentityACLUpdate,
		DeleteContext: resourceIdentityACLDelete,

		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"console", "api",
				}, true),
			},
			"ip_cidrs": {
				Type:         schema.TypeSet,
				Optional:     true,
				AtLeastOneOf: []string{"ip_ranges"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateCIDR,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyCIDRHash,
			},
			"ip_ranges": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
				Set: resourceACLPolicyRangeHash,
			},
		},
	}
}

func resourceIdentityACLCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := meta.(*config.Config).DomainID
	if err := updateACLPolicy(d, meta, id); err != nil {
		return diag.Errorf("error creating identity ACL: %s", err)
	}

	d.SetId(id)
	return resourceIdentityACLRead(ctx, d, meta)
}

func resourceIdentityACLRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	mErr := &multierror.Error{}
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	var res *acl.ACLPolicy
	switch d.Get("type").(string) {
	case "console":
		res, err = acl.ConsoleACLPolicyGet(iamClient, d.Id()).ConsoleExtract()
		if err != nil {
			return diag.Errorf("error fetching identity ACL for console access")
		}
	default:
		res, err = acl.APIACLPolicyGet(iamClient, d.Id()).APIExtract()
		if err != nil {
			return diag.Errorf("error fetching identity ACL for API access")
		}
	}

	log.Printf("[DEBUG] Retrieved identity ACL: %#v", res)
	if len(res.AllowAddressNetmasks) > 0 {
		addressNetmasks := make([]map[string]string, 0, len(res.AllowAddressNetmasks))
		for _, v := range res.AllowAddressNetmasks {
			addressNetmask := map[string]string{
				"cidr":        v.AddressNetmask,
				"description": v.Description,
			}
			addressNetmasks = append(addressNetmasks, addressNetmask)
		}
		mErr = multierror.Append(mErr, d.Set("ip_cidrs", addressNetmasks))
	}
	if len(res.AllowIPRanges) > 0 {
		ipRanges := make([]map[string]string, 0, len(res.AllowIPRanges))
		for _, v := range res.AllowIPRanges {
			ipRange := map[string]string{
				"range":       v.IPRange,
				"description": v.Description,
			}
			ipRanges = append(ipRanges, ipRange)
		}
		mErr = multierror.Append(mErr, d.Set("ip_ranges", ipRanges))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity ACL fields: %s", err)
	}

	return nil
}

func resourceIdentityACLUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	id := d.Id()
	if d.HasChanges("ip_cidrs", "ip_ranges") {
		if err := updateACLPolicy(d, meta, id); err != nil {
			return diag.Errorf("error updating identity ACL: %s", err)
		}
	}

	return resourceIdentityACLRead(ctx, d, meta)
}

func resourceIdentityACLDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	netmasksList := make([]acl.AllowAddressNetmasks, 0, 1)
	netmask := acl.AllowAddressNetmasks{
		AddressNetmask: "0.0.0.0-255.255.255.255",
	}
	netmasksList = append(netmasksList, netmask)

	deleteOpts := &acl.ACLPolicy{
		AllowAddressNetmasks: netmasksList,
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(iamClient, deleteOpts, d.Id()).ConsoleExtract()
	default:
		_, err = acl.APIACLPolicyUpdate(iamClient, deleteOpts, d.Id()).APIExtract()
	}

	if err != nil {
		return diag.Errorf("error resetting identity ACL: %s", err)
	}

	return nil
}

func updateACLPolicy(d *schema.ResourceData, meta interface{}, id string) error {
	cfg := meta.(*config.Config)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return fmt.Errorf("error creating IAM client: %s", err)
	}

	updateOpts := &acl.ACLPolicy{}
	if addressNetmasks, ok := d.GetOk("ip_cidrs"); ok {
		netmasksList := make([]acl.AllowAddressNetmasks, 0, addressNetmasks.(*schema.Set).Len())
		for _, v := range addressNetmasks.(*schema.Set).List() {
			netmask := acl.AllowAddressNetmasks{
				AddressNetmask: v.(map[string]interface{})["cidr"].(string),
				Description:    v.(map[string]interface{})["description"].(string),
			}
			netmasksList = append(netmasksList, netmask)
		}
		updateOpts.AllowAddressNetmasks = netmasksList
	}
	if ipRanges, ok := d.GetOk("ip_ranges"); ok {
		rangeList := make([]acl.AllowIPRanges, 0, ipRanges.(*schema.Set).Len())
		for _, v := range ipRanges.(*schema.Set).List() {
			ipRange := acl.AllowIPRanges{
				IPRange:     v.(map[string]interface{})["range"].(string),
				Description: v.(map[string]interface{})["description"].(string),
			}
			rangeList = append(rangeList, ipRange)
		}
		updateOpts.AllowIPRanges = rangeList
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(iamClient, updateOpts, id).ConsoleExtract()
	case "api":
		_, err = acl.APIACLPolicyUpdate(iamClient, updateOpts, id).APIExtract()
	}

	return err
}

func resourceACLPolicyCIDRHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["cidr"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["cidr"].(string)))
	}

	return hashcode.String(buf.String())
}

func resourceACLPolicyRangeHash(v interface{}) int {
	var buf bytes.Buffer
	m := v.(map[string]interface{})

	if m["range"] != nil {
		buf.WriteString(fmt.Sprintf("%s-", m["range"].(string)))
	}

	return hashcode.String(buf.String())
}
