package iam

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/identity/v3.0/acl"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceAclNonUpdatableParams = []string{"type"}

// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
func ResourceAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAclCreate,
		ReadContext:   resourceAclRead,
		UpdateContext: resourceAclUpdate,
		DeleteContext: resourceAclDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceAclNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the ACL policy.`,
			},
			"ip_cidrs": {
				Type:         schema.TypeList,
				Optional:     true,
				AtLeastOneOf: []string{"ip_ranges"},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:         schema.TypeString,
							Required:     true,
							ValidateFunc: utils.ValidateCIDR,
							Description:  `The IPv4 CIDR block which allow access through console or API.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the IPv4 CIDR block.`,
						},
					},
				},
				Description: `The list of IPv4 CIDR blocks from which console access or API access is allowed.`,
			},
			"ip_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The IPv4 address range which allow access through console or API.`,
						},
						"description": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The description of the IP address range.`,
						},
					},
				},
				Description: `The list of IPv4 address ranges from which console access or API access is allowed.`,
			},

			// Internal attributes.
			"ip_ciders_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The origin IPv4 CIDR block.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of IPv4 CIDR blocks that used to reorder the 'ip_cidrs' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"ip_ranges_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The origin IPv4 range.`,
						},
					},
				},
				Description: utils.SchemaDesc(
					`The origin list of IPv4 ranges that used to reorder the 'ip_ranges' parameter.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},

			// Internal parameters.
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

func buildIpCidersOrder(d *schema.ResourceData) []interface{} {
	ipCiders, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ip_cidrs").([]interface{})
	if !ok || ipCiders == nil {
		return nil
	}

	result := make([]interface{}, 0, len(ipCiders))
	for _, ipCider := range ipCiders {
		result = append(result, map[string]interface{}{
			"cidr": utils.PathSearch("cidr", ipCider, "").(string),
		})
	}

	return result
}

func buildIpRangesOrder(d *schema.ResourceData) []interface{} {
	ipRanges, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ip_ranges").([]interface{})
	if !ok || ipRanges == nil {
		return nil
	}

	result := make([]interface{}, 0, len(ipRanges))
	for _, ipRange := range ipRanges {
		result = append(result, map[string]interface{}{
			"range": utils.PathSearch("range", ipRange, "").(string),
		})
	}

	return result
}

func resourceAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		err      error
		domainId = cfg.DomainID
	)

	// ACL policy change operations may encounter concurrency issues (causing other ACL policy changes to fail),
	// so, it is necessary to lock the domain ID to prevent concurrent changes.
	config.MutexKV.Lock(domainId)
	defer config.MutexKV.Unlock(domainId)

	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", domainId, d.Get("type").(string)))
	err = updateAclPolicy(client, d, domainId)
	if err != nil {
		return diag.Errorf("error creating identity ACL: %s", err)
	}

	if err = d.Set("ip_ciders_order", buildIpCidersOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ciders_order field after creating ACL: %s", err)
	}
	if err = d.Set("ip_ranges_order", buildIpRangesOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ranges_order field after creating ACL: %s", err)
	}

	return resourceAclRead(ctx, d, meta)
}

func orderIpCidersByIpCidersOrderOrigin(ipCiders, ipCidersOrigin []interface{}) []interface{} {
	if len(ipCidersOrigin) == 0 {
		return ipCiders
	}

	sortedIpCiders := make([]interface{}, 0)
	ipCidersCopy := ipCiders
	for _, ipCiderOrigin := range ipCidersOrigin {
		cidrOrigin := utils.PathSearch("cidr", ipCiderOrigin, "").(string)
		for index, ipCider := range ipCidersCopy {
			if utils.PathSearch("cidr", ipCider, "").(string) != cidrOrigin {
				continue
			}
			// Add the found ip cidr to the sorted ip ciders list.
			sortedIpCiders = append(sortedIpCiders, ipCidersCopy[index])
			// Remove the processed ip cidr from the original ip ciders array.
			ipCidersCopy = append(ipCidersCopy[:index], ipCidersCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted ip ciders to the end of the sorted list.
	sortedIpCiders = append(sortedIpCiders, ipCidersCopy...)
	return sortedIpCiders
}

func flattenIpCiders(ipCiders, ipCidersOrigin []interface{}) []interface{} {
	if len(ipCiders) < 1 {
		return nil
	}

	sortedIpCiders := orderIpCidersByIpCidersOrderOrigin(ipCiders, ipCidersOrigin)
	result := make([]interface{}, 0, len(sortedIpCiders))
	for _, ipCider := range sortedIpCiders {
		result = append(result, map[string]interface{}{
			"cidr":        utils.PathSearch("cidr", ipCider, nil),
			"description": utils.PathSearch("description", ipCider, nil),
		})
	}

	return result
}

func orderIpRangesByIpRangesOrderOrigin(ipRanges, ipRangesOrigin []interface{}) []interface{} {
	if len(ipRangesOrigin) == 0 {
		return ipRanges
	}

	sortedIpRanges := make([]interface{}, 0)
	ipRangesCopy := ipRanges
	for _, ipRangeOrigin := range ipRangesOrigin {
		rangeOrigin := utils.PathSearch("range", ipRangeOrigin, "").(string)
		for index, ipRange := range ipRangesCopy {
			if utils.PathSearch("range", ipRange, "").(string) != rangeOrigin {
				continue
			}
			// Add the found ip cidr to the sorted ip ciders list.
			sortedIpRanges = append(sortedIpRanges, ipRangesCopy[index])
			// Remove the processed ip cidr from the original ip ciders array.
			ipRangesCopy = append(ipRangesCopy[:index], ipRangesCopy[index+1:]...)
		}
	}
	// Add any remaining unsorted ip ranges to the end of the sorted list.
	sortedIpRanges = append(sortedIpRanges, ipRangesCopy...)
	return sortedIpRanges
}

func flattenIpRanges(ipRanges, ipRangesOrigin []interface{}) []interface{} {
	if len(ipRanges) < 1 {
		return nil
	}

	sortedIpRanges := orderIpRangesByIpRangesOrderOrigin(ipRanges, ipRangesOrigin)
	result := make([]interface{}, 0, len(sortedIpRanges))
	for _, ipRange := range sortedIpRanges {
		result = append(result, map[string]interface{}{
			"range":       utils.PathSearch("range", ipRange, nil),
			"description": utils.PathSearch("description", ipRange, nil),
		})
	}

	return result
}

func GetAclByDomainId(client *golangsdk.ServiceClient, aclType, domainId string) (*acl.ACLPolicy, error) {
	var (
		result *acl.ACLPolicy
		err    error
	)
	switch aclType {
	case "console":
		result, err = acl.ConsoleACLPolicyGet(client, domainId).ConsoleExtract()
		if err != nil {
			return nil, err
		}
		if len(result.AllowAddressNetmasks) == 0 && len(result.AllowIPRanges) == 1 &&
			result.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("identity ACL for console access <%s> has been reverted", domainId)),
				},
			}
		}
	case "api":
		result, err = acl.APIACLPolicyGet(client, domainId).APIExtract()
		if err != nil {
			return nil, err
		}
		if len(result.AllowAddressNetmasks) == 0 && len(result.AllowIPRanges) == 1 &&
			result.AllowIPRanges[0].IPRange == "0.0.0.0-255.255.255.255" {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy",
					RequestId: "NONE",
					Body:      []byte(fmt.Sprintf("identity ACL for API access (domain: %s) has been reverted", domainId)),
				},
			}
		}
	default:
		return nil, fmt.Errorf("invalid ACL type: %s", aclType)
	}
	return result, nil
}

func resourceAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr     = &multierror.Error{}
		cfg      = meta.(*config.Config)
		domainId = cfg.DomainID
	)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	aclPolicy, err := GetAclByDomainId(iamClient, d.Get("type").(string), domainId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching identity ACL")
	}

	if len(aclPolicy.AllowAddressNetmasks) > 0 {
		addressNetmasks := make([]interface{}, 0, len(aclPolicy.AllowAddressNetmasks))
		for _, v := range aclPolicy.AllowAddressNetmasks {
			addressNetmasks = append(addressNetmasks, map[string]interface{}{
				"cidr":        v.AddressNetmask,
				"description": v.Description,
			})
		}
		mErr = multierror.Append(mErr, d.Set("ip_cidrs", flattenIpCiders(addressNetmasks, d.Get("ip_ciders_order").([]interface{}))))
	}
	if len(aclPolicy.AllowIPRanges) > 0 {
		ipRanges := make([]interface{}, 0, len(aclPolicy.AllowIPRanges))
		for _, v := range aclPolicy.AllowIPRanges {
			ipRanges = append(ipRanges, map[string]interface{}{
				"range":       v.IPRange,
				"description": v.Description,
			})
		}
		mErr = multierror.Append(mErr, d.Set("ip_ranges", flattenIpRanges(ipRanges, d.Get("ip_ranges_order").([]interface{}))))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity ACL fields: %s", err)
	}
	return nil
}

func resourceAclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		domainId = cfg.DomainID
	)

	// ACL policy change operations may encounter concurrency issues (causing other ACL policy changes to fail),
	// so, it is necessary to lock the domain ID to prevent concurrent changes.
	config.MutexKV.Lock(domainId)
	defer config.MutexKV.Unlock(domainId)

	client, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	if err := updateAclPolicy(client, d, domainId); err != nil {
		return diag.Errorf("error updating identity ACL: %s", err)
	}

	if err = d.Set("ip_ciders_order", buildIpCidersOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ciders_order field after updating ACL: %s", err)
	}
	if err = d.Set("ip_ranges_order", buildIpRangesOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ranges_order field after updating ACL: %s", err)
	}

	return resourceAclRead(ctx, d, meta)
}

func resourceAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		domainId = cfg.DomainID
	)

	// ACL policy change operations may encounter concurrency issues (causing other ACL policy changes to fail),
	// so, it is necessary to lock the domain ID to prevent concurrent changes.
	config.MutexKV.Lock(domainId)
	defer config.MutexKV.Unlock(domainId)

	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	deleteOpts := &acl.ACLPolicy{
		AllowAddressNetmasks: []acl.AllowAddressNetmasks{
			{
				AddressNetmask: "0.0.0.0-255.255.255.255",
			},
		},
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(iamClient, deleteOpts, domainId).ConsoleExtract()
	default:
		_, err = acl.APIACLPolicyUpdate(iamClient, deleteOpts, domainId).APIExtract()
	}

	if err != nil {
		return diag.Errorf("error resetting identity ACL: %s", err)
	}

	return nil
}

func updateAclPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string) error {
	var (
		updateOpts = &acl.ACLPolicy{}
		err        error
	)

	if addressNetmasks, ok := d.Get("ip_cidrs").([]interface{}); ok && len(addressNetmasks) > 0 {
		netmasksList := make([]acl.AllowAddressNetmasks, 0, len(addressNetmasks))
		for _, v := range addressNetmasks {
			netmasksList = append(netmasksList, acl.AllowAddressNetmasks{
				AddressNetmask: v.(map[string]interface{})["cidr"].(string),
				Description:    v.(map[string]interface{})["description"].(string),
			})
		}
		updateOpts.AllowAddressNetmasks = netmasksList
	}

	if ipRanges, ok := d.Get("ip_ranges").([]interface{}); ok && len(ipRanges) > 0 {
		ipRangesList := make([]acl.AllowIPRanges, 0, len(ipRanges))
		for _, v := range ipRanges {
			ipRangesList = append(ipRangesList, acl.AllowIPRanges{
				IPRange:     v.(map[string]interface{})["range"].(string),
				Description: v.(map[string]interface{})["description"].(string),
			})
		}
		updateOpts.AllowIPRanges = ipRangesList
	}

	switch d.Get("type").(string) {
	case "console":
		_, err = acl.ConsoleACLPolicyUpdate(client, updateOpts, domainId).ConsoleExtract()
	case "api":
		_, err = acl.APIACLPolicyUpdate(client, updateOpts, domainId).APIExtract()
	}

	return err
}
