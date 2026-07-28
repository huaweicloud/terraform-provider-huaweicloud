package elb

import (
	"context"
	"fmt"
	"log"
	"strconv"
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

var domainResolutionNonUpdatableParams = []string{"loadbalancer_id"}

// @API ELB POST /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/user-defined-config
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips
// @API ELB GET /v3/{project_id}/elb/loadbalancers/{loadbalancer_id}
func ResourceDomainResolution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainResolutionCreate,
		ReadContext:   resourceDomainResolutionRead,
		UpdateContext: resourceDomainResolutionUpdate,
		DeleteContext: resourceDomainResolutionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(domainResolutionNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"public_domain_name_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"private_domain_name_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"public_dns_zone_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_dns_record_set_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"private_dns_zone_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_dns_zone_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"private_dns_record_set_ttl": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"ips": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     domainResolutionSchema(),
			},
		},
	}
}

func domainResolutionSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDomainResolutionParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"loadbalancer": map[string]interface{}{
			"public_domain_name_enable":  utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "public_domain_name_enable"),
			"public_dns_zone_name":       utils.ValueIgnoreEmpty(d.Get("public_dns_zone_name")),
			"public_dns_record_set_ttl":  utils.ValueIgnoreEmpty(d.Get("public_dns_record_set_ttl")),
			"private_domain_name_enable": utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "private_domain_name_enable"),
			"private_dns_zone_name":      utils.ValueIgnoreEmpty(d.Get("private_dns_zone_name")),
			"private_dns_zone_type":      utils.ValueIgnoreEmpty(d.Get("private_dns_zone_type")),
			"private_dns_record_set_ttl": utils.ValueIgnoreEmpty(d.Get("private_dns_record_set_ttl")),
		},
	}

	return params
}

func resourceDomainResolutionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		loadbalancerId = d.Get("loadbalancer_id").(string)
		httpUrl        = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/user-defined-config"
		publicEnabled  = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "public_domain_name_enable")
		privateEnabled = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "private_domain_name_enable")
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	if !domainResolutionEnabled(publicEnabled, privateEnabled) {
		return diag.Errorf("at least one of `public_domain_name_enable` and `private_domain_name_enable` value needs to be true")
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{loadbalancer_id}", loadbalancerId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         utils.RemoveNil(buildDomainResolutionParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating domain name resolution: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(loadbalancerId)

	return resourceDomainResolutionRead(ctx, d, meta)
}

func domainResolutionEnabled(publicEnabled, privateEnabled interface{}) bool {
	var (
		publicSwitch  bool
		privateSwitch bool
	)

	if publicEnabled != nil {
		publicSwitch = publicEnabled.(bool)
	}
	if privateEnabled != nil {
		privateSwitch = privateEnabled.(bool)
	}

	return publicSwitch || privateSwitch
}

func resourceDomainResolutionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	domainInformation, ips, err := getDomainResolution(client, d.Id())
	if err != nil {
		// If the load balancer not exsit, the query API return code is `404`.
		return common.CheckDeletedDiag(d, err, "error retrieving domain name resolution information")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("loadbalancer_id", d.Id()),
		d.Set("public_domain_name_enable", utils.PathSearch("public_domain_name_enable", domainInformation, nil)),
		d.Set("public_dns_zone_name", utils.PathSearch("public_dns_zone_name", domainInformation, nil)),
		d.Set("public_dns_record_set_ttl", utils.PathSearch("public_dns_record_set_ttl", domainInformation, nil)),
		d.Set("private_domain_name_enable", utils.PathSearch("private_domain_name_enable", domainInformation, nil)),
		d.Set("private_dns_zone_name", utils.PathSearch("private_dns_zone_name", domainInformation, nil)),
		d.Set("private_dns_zone_type", utils.PathSearch("private_dns_zone_type", domainInformation, nil)),
		d.Set("private_dns_record_set_ttl", utils.PathSearch("private_dns_record_set_ttl", domainInformation, nil)),
		d.Set("ips", flattenDomainResolutionIpAddresses(ips)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getDomainResolution(client *golangsdk.ServiceClient, loadbalancerId string) (interface{}, []interface{}, error) {
	respBody, err := GetDomainResolutionInformation(client, loadbalancerId)
	if err != nil {
		return nil, nil, err
	}

	ips, err := getDomainResolutionIpAddresses(client, loadbalancerId)
	if err != nil {
		log.Printf("error retrieving domain name resolution IP address information")
	}

	return respBody, ips, nil
}

func GetDomainResolutionInformation(client *golangsdk.ServiceClient, loadbalancerId string) (interface{}, error) {
	httpUrl := "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{loadbalancer_id}", loadbalancerId)
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	resp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	domainInformation := utils.PathSearch("loadbalancer.user_defined_dns", respBody, nil)
	if domainInformation == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return domainInformation, nil
}

func getDomainResolutionIpAddresses(client *golangsdk.ServiceClient, loadbalancerId string) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/ips?limit={limit}"
		result  = make([]interface{}, 0)
		limit   = 2000
		marker  = ""
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{loadbalancer_id}", loadbalancerId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
	}

	for {
		listPathWithMarker := listPath
		if marker != "" {
			listPathWithMarker = fmt.Sprintf("%s&marker=%s", listPathWithMarker, marker)
		}

		resp, err := client.Request("GET", listPathWithMarker, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		ipAddreesses := utils.PathSearch("ips", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, ipAddreesses...)
		if len(ipAddreesses) < limit {
			break
		}

		marker = utils.PathSearch("page_info.next_marker", respBody, "").(string)
		if marker == "" {
			break
		}
	}

	return result, nil
}

func flattenDomainResolutionIpAddresses(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	res := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"enable":      utils.PathSearch("enable", v, nil),
			"ip_address":  utils.PathSearch("ip_address", v, nil),
			"type":        utils.PathSearch("type", v, nil),
			"domain_name": utils.PathSearch("domain_name", v, nil),
			"created_at":  utils.PathSearch("created_at", v, nil),
			"updated_at":  utils.PathSearch("updated_at", v, nil),
		})
	}
	return res
}

func resourceDomainResolutionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg            = meta.(*config.Config)
		region         = cfg.GetRegion(d)
		httpUrl        = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/user-defined-config"
		publicEnabled  = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "public_domain_name_enable")
		privateEnabled = utils.GetNestedObjectFromRawConfig(d.GetRawConfig(), "private_domain_name_enable")
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	if d.HasChangeExcept("enable_force_new") {
		if !domainResolutionEnabled(publicEnabled, privateEnabled) {
			return diag.Errorf("at least one of `public_domain_name_enable` and `private_domain_name_enable` value needs to be true")
		}

		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{loadbalancer_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
			JSONBody:         utils.RemoveNil(buildDomainResolutionParams(d)),
		}

		_, err = client.Request("POST", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating domain name resolution: %s", err)
		}
	}

	return resourceDomainResolutionRead(ctx, d, meta)
}

func buildDeleteDomainResolutionParams() map[string]interface{} {
	params := map[string]interface{}{
		"loadbalancer": map[string]interface{}{
			"public_domain_name_enable":  false,
			"private_domain_name_enable": false,
		},
	}

	return params
}

func resourceDomainResolutionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/elb/loadbalancers/{loadbalancer_id}/dns/user-defined-config"
	)

	client, err := cfg.NewServiceClient("elb", region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{loadbalancer_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json;charset=UTF-8"},
		JSONBody:         buildDeleteDomainResolutionParams(),
	}

	_, err = client.Request("POST", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting domain name resolution")
	}

	return nil
}
