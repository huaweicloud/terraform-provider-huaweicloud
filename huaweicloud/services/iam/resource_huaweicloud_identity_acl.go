package iam

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v3AclNonUpdatableParams = []string{"type"}

// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy
// @API IAM GET /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
// @API IAM PUT /v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy
func ResourceV3Acl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3AclCreate,
		ReadContext:   resourceV3AclRead,
		UpdateContext: resourceV3AclUpdate,
		DeleteContext: resourceV3AclDelete,

		CustomizeDiff: config.FlexibleForceNew(v3AclNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the ACL policy.`,
			},
			"ip_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
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
				DiffSuppressFunc: func(key, oldVal, newVal string, d *schema.ResourceData) bool {
					ipCidrs := d.Get("ip_cidrs").([]interface{})
					if key == "ip_ranges.#" {
						return len(ipCidrs) == 0 && oldVal == "1" && newVal == "0"
					}
					if len(ipCidrs) == 0 && newVal == "" &&
						oldVal == "0.0.0.0-255.255.255.255" {
						return true
					}
					return false
				},
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
			"ipv6_cidrs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Required: true,
						},
						"description": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"ipv6_ranges": {
				Type:     schema.TypeList,
				Optional: true,
				DiffSuppressFunc: func(key, oldVal, newVal string, d *schema.ResourceData) bool {
					ipv6Cidrs := d.Get("ipv6_cidrs").([]interface{})
					if key == "ipv6_ranges.#" {
						return len(ipv6Cidrs) == 0 && oldVal == "1" && newVal == "0"
					}
					if len(ipv6Cidrs) == 0 && newVal == "" &&
						oldVal == "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF" {
						return true
					}
					return false
				},
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
			"ipv6_ciders_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr": {
							Type:     schema.TypeString,
							Computed: true,
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
			"ipv6_ranges_order": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				DiffSuppressFunc: utils.SuppressDiffAll,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"range": {
							Type:     schema.TypeString,
							Computed: true,
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

func buildV3AclIpCidersOrder(d *schema.ResourceData) []interface{} {
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

func buildV3AclIpRangesOrder(d *schema.ResourceData) []interface{} {
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

func buildV3AclIpv6CidersOrder(d *schema.ResourceData) []interface{} {
	ipv6Ciders, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ipv6_cidrs").([]interface{})
	if !ok || ipv6Ciders == nil {
		return nil
	}

	result := make([]interface{}, 0, len(ipv6Ciders))
	for _, ipv6Cider := range ipv6Ciders {
		result = append(result, map[string]interface{}{
			"cidr": utils.PathSearch("cidr", ipv6Cider, "").(string),
		})
	}

	return result
}

func buildV3AclIpv6RangesOrder(d *schema.ResourceData) []interface{} {
	ipv6Ranges, ok := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ipv6_ranges").([]interface{})
	if !ok || ipv6Ranges == nil {
		return nil
	}

	result := make([]interface{}, 0, len(ipv6Ranges))
	for _, ipv6Range := range ipv6Ranges {
		result = append(result, map[string]interface{}{
			"range": utils.PathSearch("range", ipv6Range, "").(string),
		})
	}

	return result
}

func resourceV3AclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = updateV3AclPolicy(client, d, domainId, utils.RemoveNil(buildUpdateV3AclPolicyBodyParams(d)))
	if err != nil {
		return diag.Errorf("error creating identity ACL: %s", err)
	}
	d.SetId(fmt.Sprintf("%s/%s", domainId, d.Get("type").(string)))

	if err = d.Set("ip_ciders_order", buildV3AclIpCidersOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ciders_order field after creating ACL: %s", err)
	}
	if err = d.Set("ip_ranges_order", buildV3AclIpRangesOrder(d)); err != nil {
		log.Printf("[ERROR] error setting the ip_ranges_order field after creating ACL: %s", err)
	}
	if d.Get("type").(string) == "console" {
		if err = d.Set("ipv6_ciders_order", buildV3AclIpv6CidersOrder(d)); err != nil {
			log.Printf("[ERROR] error setting the ipv6_ciders_order field after creating ACL: %s", err)
		}
		if err = d.Set("ipv6_ranges_order", buildV3AclIpv6RangesOrder(d)); err != nil {
			log.Printf("[ERROR] error setting the ipv6_ranges_order field after creating ACL: %s", err)
		}
	}

	return resourceV3AclRead(ctx, d, meta)
}

func orderV3AclIpCidersByIpCidersOrderOrigin(ipCiders, ipCidersOrigin []interface{}) []interface{} {
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
			break
		}
	}
	// Add any remaining unsorted ip ciders to the end of the sorted list.
	sortedIpCiders = append(sortedIpCiders, ipCidersCopy...)
	return sortedIpCiders
}

func flattenV3AclIpCiders(ipCiders, ipCidersOrigin []interface{}) []interface{} {
	if len(ipCiders) < 1 {
		return nil
	}

	sortedIpCiders := orderV3AclIpCidersByIpCidersOrderOrigin(ipCiders, ipCidersOrigin)
	result := make([]interface{}, 0, len(sortedIpCiders))
	for _, ipCider := range sortedIpCiders {
		result = append(result, map[string]interface{}{
			"cidr":        utils.PathSearch("cidr", ipCider, nil),
			"description": utils.PathSearch("description", ipCider, nil),
		})
	}

	return result
}

func orderV3AclIpRangesByIpRangesOrderOrigin(ipRanges, ipRangesOrigin []interface{}) []interface{} {
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
			break
		}
	}
	// Add any remaining unsorted ip ranges to the end of the sorted list.
	sortedIpRanges = append(sortedIpRanges, ipRangesCopy...)
	return sortedIpRanges
}

func flattenV3AclIpRanges(ipRanges, ipRangesOrigin []interface{}) []interface{} {
	if len(ipRanges) < 1 {
		return nil
	}

	sortedIpRanges := orderV3AclIpRangesByIpRangesOrderOrigin(ipRanges, ipRangesOrigin)
	result := make([]interface{}, 0, len(sortedIpRanges))
	for _, ipRange := range sortedIpRanges {
		result = append(result, map[string]interface{}{
			"range":       utils.PathSearch("range", ipRange, nil),
			"description": utils.PathSearch("description", ipRange, nil),
		})
	}

	return result
}

func GetV3AclByDomainId(client *golangsdk.ServiceClient, aclType, domainId string) (interface{}, error) {
	httpUrl := ""
	switch aclType {
	case "console":
		httpUrl = "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy"
	case "api":
		httpUrl = "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy"
	default:
		return nil, fmt.Errorf("invalid acl type: %s, should be console or api", aclType)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{domain_id}", domainId)
	getOpts := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpts)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	var respBody interface{}
	if aclType == "console" {
		respBody = utils.PathSearch("console_acl_policy", getRespBody, nil)
		ipCidrs := utils.PathSearch("allow_address_netmasks", respBody, make([]interface{}, 0)).([]interface{})
		ipRanges := utils.PathSearch("allow_ip_ranges", respBody, make([]interface{}, 0)).([]interface{})
		ipRange := utils.PathSearch("allow_ip_ranges[0].ip_range", respBody, "").(string)
		ipv6Cidrs := utils.PathSearch("allow_address_netmasks_ipv6", respBody, make([]interface{}, 0)).([]interface{})
		ipv6Ranges := utils.PathSearch("allow_ip_ranges_ipv6", respBody, make([]interface{}, 0)).([]interface{})
		ipv6Range := utils.PathSearch("allow_ip_ranges_ipv6[0].ip_range", respBody, "").(string)
		if len(ipCidrs) == 0 && len(ipRanges) == 1 && ipRange == "0.0.0.0-255.255.255.255" &&
			len(ipv6Cidrs) == 0 && len(ipv6Ranges) == 1 &&
			ipv6Range == "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF" {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy",
					RequestId: "NONE",
					Body:      fmt.Appendf(nil, "identity ACL for console access <%s> has been reverted", domainId),
				},
			}
		}
	} else {
		respBody = utils.PathSearch("api_acl_policy", getRespBody, nil)
		ipCidrs := utils.PathSearch("allow_address_netmasks", respBody, make([]interface{}, 0)).([]interface{})
		ipRanges := utils.PathSearch("allow_ip_ranges", respBody, make([]interface{}, 0)).([]interface{})
		ipRange := utils.PathSearch("allow_ip_ranges[0].ip_range", respBody, "").(string)
		if len(ipCidrs) == 0 && len(ipRanges) == 1 && ipRange == "0.0.0.0-255.255.255.255" {
			return nil, golangsdk.ErrDefault404{
				ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
					Method:    "GET",
					URL:       "/v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy",
					RequestId: "NONE",
					Body:      fmt.Appendf(nil, "identity ACL for API access (domain: %s) has been reverted", domainId),
				},
			}
		}
	}

	return respBody, nil
}

func resourceV3AclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		mErr     = &multierror.Error{}
		cfg      = meta.(*config.Config)
		domainId = cfg.DomainID
	)
	iamClient, err := cfg.IAMV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating IAM client: %s", err)
	}

	aclPolicy, err := GetV3AclByDomainId(iamClient, d.Get("type").(string), domainId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error fetching identity ACL")
	}

	ipCidrs := utils.PathSearch("allow_address_netmasks", aclPolicy, make([]interface{}, 0)).([]interface{})
	if len(ipCidrs) > 0 {
		addressNetmasks := make([]interface{}, 0, len(ipCidrs))
		for _, v := range ipCidrs {
			addressNetmasks = append(addressNetmasks, map[string]interface{}{
				"cidr":        utils.PathSearch("address_netmask", v, nil),
				"description": utils.PathSearch("description", v, nil),
			})
		}
		mErr = multierror.Append(mErr, d.Set("ip_cidrs", flattenV3AclIpCiders(addressNetmasks, d.Get("ip_ciders_order").([]interface{}))))
	}
	allowIpRanges := utils.PathSearch("allow_ip_ranges", aclPolicy, make([]interface{}, 0)).([]interface{})
	if len(allowIpRanges) > 0 {
		ipRanges := make([]interface{}, 0, len(allowIpRanges))
		for _, v := range allowIpRanges {
			ipRanges = append(ipRanges, map[string]interface{}{
				"range":       utils.PathSearch("ip_range", v, nil),
				"description": utils.PathSearch("description", v, nil),
			})
		}
		mErr = multierror.Append(mErr, d.Set("ip_ranges", flattenV3AclIpRanges(ipRanges, d.Get("ip_ranges_order").([]interface{}))))
	}
	ipv6Cidrs := utils.PathSearch("allow_address_netmasks_ipv6", aclPolicy, make([]interface{}, 0)).([]interface{})
	if len(ipv6Cidrs) > 0 {
		addressNetmasksv6 := make([]interface{}, 0, len(ipv6Cidrs))
		for _, v := range ipv6Cidrs {
			addressNetmasksv6 = append(addressNetmasksv6, map[string]interface{}{
				"cidr":        utils.PathSearch("address_netmask", v, nil),
				"description": utils.PathSearch("description", v, nil),
			})
		}
		mErr = multierror.Append(mErr, d.Set("ipv6_cidrs",
			flattenV3AclIpCiders(addressNetmasksv6, d.Get("ipv6_ciders_order").([]interface{}))))
	}
	allowIpv6Ranges := utils.PathSearch("allow_ip_ranges_ipv6", aclPolicy, make([]interface{}, 0)).([]interface{})
	if len(allowIpv6Ranges) > 0 {
		ipv6Ranges := make([]interface{}, 0, len(allowIpv6Ranges))
		for _, v := range allowIpv6Ranges {
			ipv6Ranges = append(ipv6Ranges, map[string]interface{}{
				"range":       utils.PathSearch("ip_range", v, nil),
				"description": utils.PathSearch("description", v, nil),
			})
		}
		mErr = multierror.Append(mErr, d.Set("ipv6_ranges",
			flattenV3AclIpRanges(ipv6Ranges, d.Get("ipv6_ranges_order").([]interface{}))))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting identity ACL fields: %s", err)
	}
	return nil
}

func resourceV3AclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	if d.HasChangeExcept("enable_force_new") {
		if err = updateV3AclPolicy(client, d, domainId, utils.RemoveNil(buildUpdateV3AclPolicyBodyParams(d))); err != nil {
			return diag.Errorf("error updating identity ACL: %s", err)
		}

		if err = d.Set("ip_ciders_order", buildV3AclIpCidersOrder(d)); err != nil {
			log.Printf("[ERROR] error setting the ip_ciders_order field after updating ACL: %s", err)
		}
		if err = d.Set("ip_ranges_order", buildV3AclIpRangesOrder(d)); err != nil {
			log.Printf("[ERROR] error setting the ip_ranges_order field after updating ACL: %s", err)
		}
		if d.Get("type").(string) == "console" {
			if err = d.Set("ipv6_ciders_order", buildV3AclIpv6CidersOrder(d)); err != nil {
				log.Printf("[ERROR] error setting the ipv6_ciders_order field after creating ACL: %s", err)
			}
			if err = d.Set("ipv6_ranges_order", buildV3AclIpv6RangesOrder(d)); err != nil {
				log.Printf("[ERROR] error setting the ipv6_ranges_order field after creating ACL: %s", err)
			}
		}
	}

	return resourceV3AclRead(ctx, d, meta)
}

func resourceV3AclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	err = updateV3AclPolicy(iamClient, d, domainId, buildDeleteV3AclPolicyBodyParams(d))
	if err != nil {
		return diag.Errorf("error resetting identity ACL: %s", err)
	}

	return nil
}

func buildDeleteV3AclPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"allow_address_netmasks": []map[string]interface{}{
			{
				"address_netmask": "0.0.0.0-255.255.255.255",
			},
		},
	}
	if d.Get("type") == "console" {
		bodyParams["allow_address_netmasks_ipv6"] = []map[string]interface{}{
			{
				"address_netmask": "0000:0000:0000:0000:0000:0000:0000:0000-FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF:FFFF",
			},
		}
		return map[string]interface{}{
			"console_acl_policy": bodyParams,
		}
	}

	return map[string]interface{}{
		"api_acl_policy": bodyParams,
	}
}

func updateV3AclPolicy(client *golangsdk.ServiceClient, d *schema.ResourceData, domainId string,
	bodyParams interface{}) error {
	httpUrl := ""
	policyType := d.Get("type").(string)
	switch policyType {
	case "console":
		httpUrl = "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/console-acl-policy"
	case "api":
		httpUrl = "v3.0/OS-SECURITYPOLICY/domains/{domain_id}/api-acl-policy"
	default:
		return fmt.Errorf("invalid policy type: %s, should be console or api", policyType)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{domain_id}", domainId)
	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		KeepResponseBody: true,
		JSONBody:         bodyParams,
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)

	return err
}

func buildUpdateV3AclPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"allow_address_netmasks": buildUpdateV3AclPolicyIpCidrsBodyParams(d),
		"allow_ip_ranges":        buildUpdateV3AclPolicyIpRangesBodyParams(d),
	}
	if d.Get("type") == "console" {
		bodyParams["allow_address_netmasks_ipv6"] = buildUpdateV3AclPolicyIpv6CidrsBodyParams(d)
		bodyParams["allow_ip_ranges_ipv6"] = buildUpdateV3AclPolicyIpv6RangesBodyParams(d)
		return map[string]interface{}{
			"console_acl_policy": bodyParams,
		}
	}
	return map[string]interface{}{
		"api_acl_policy": bodyParams,
	}
}

func buildUpdateV3AclPolicyIpCidrsBodyParams(d *schema.ResourceData) []interface{} {
	ipCidrs := d.Get("ip_cidrs").([]interface{})
	if len(ipCidrs) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(ipCidrs))
	for _, ipCidr := range ipCidrs {
		if v, ok := ipCidr.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"address_netmask": v["cidr"],
				"description":     utils.ValueIgnoreEmpty(v["description"]),
			})
		}
	}
	return res
}

func buildUpdateV3AclPolicyIpRangesBodyParams(d *schema.ResourceData) []interface{} {
	ipRanges := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ip_ranges").([]interface{})
	if len(ipRanges) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(ipRanges))
	for _, ipRange := range ipRanges {
		if v, ok := ipRange.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"ip_range":    v["range"],
				"description": utils.ValueIgnoreEmpty(v["description"]),
			})
		}
	}
	return res
}

func buildUpdateV3AclPolicyIpv6CidrsBodyParams(d *schema.ResourceData) []interface{} {
	ipv6Cidrs := d.Get("ipv6_cidrs").([]interface{})
	if len(ipv6Cidrs) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(ipv6Cidrs))
	for _, ipv6Cidr := range ipv6Cidrs {
		if v, ok := ipv6Cidr.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"address_netmask": v["cidr"],
				"description":     utils.ValueIgnoreEmpty(v["description"]),
			})
		}
	}
	return res
}

func buildUpdateV3AclPolicyIpv6RangesBodyParams(d *schema.ResourceData) []interface{} {
	ipv6Ranges := utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "ipv6_ranges").([]interface{})
	if len(ipv6Ranges) == 0 {
		return nil
	}
	res := make([]interface{}, 0, len(ipv6Ranges))
	for _, ipv6Range := range ipv6Ranges {
		if v, ok := ipv6Range.(map[string]interface{}); ok {
			res = append(res, map[string]interface{}{
				"ip_range":    v["range"],
				"description": utils.ValueIgnoreEmpty(v["description"]),
			})
		}
	}
	return res
}
