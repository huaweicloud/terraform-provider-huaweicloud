package cfw

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW GET /v1/{project_id}/dns/servers
// @API CFW PUT /v1/{project_id}/dns/servers
func ResourceDNSResolution() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDNSResolutionCreateOrUpdate,
		ReadContext:   resourceDNSResolutionRead,
		UpdateContext: resourceDNSResolutionCreateOrUpdate,
		DeleteContext: resourceDNSResolutionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"fw_instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the firewall.`,
			},
			"default_dns_servers": {
				Type: schema.TypeSet,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`),
						"the IP format should be xxx.xxx.xxx.xxx where xxx is a number between 0 and 255",
					),
				},
				Optional:     true,
				MinItems:     1,
				AtLeastOneOf: []string{"default_dns_servers", "custom_dns_servers"},
				Description:  `The default DNS servers.`,
			},
			"custom_dns_servers": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
					ValidateFunc: validation.StringMatch(
						regexp.MustCompile(`^((25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))\.){3}(25[0-5]|2[0-4]\d|(1\d{2}|[1-9]?\d))$`),
						"the IP format should be xxx.xxx.xxx.xxx where xxx is a number between 0 and 255",
					),
				},
				MinItems:    1,
				Description: `The custom DNS servers.`,
			},
			"health_check_domain_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The health check domain name.`,
			},
		},
	}
}

func resourceDNSResolutionCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl      = "v1/{project_id}/dns/servers"
		product      = "cfw"
		fwInstanceID = d.Get("fw_instance_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	respBody, err := readDNSServers(client, fwInstanceID)
	if err != nil {
		return diag.Errorf("error creating CFW DNS resolution configuration: %s", err)
	}

	defaultDNSServersExpr := "data[?is_customized==`0`].server_ip"
	defaultDNSServers := utils.PathSearch(defaultDNSServersExpr, respBody, make([]interface{}, 0)).([]interface{})
	if len(defaultDNSServers) == 0 {
		return diag.Errorf("error creating CFW DNS resolution configuration: the default DNS servers can not be found")
	}

	userDefaultDNSServers := d.Get("default_dns_servers").(*schema.Set).List()
	constains := utils.StrSliceContainsAnother(utils.ExpandToStringList(defaultDNSServers), utils.ExpandToStringList(userDefaultDNSServers))
	if !constains {
		return diag.Errorf("error creating CFW DNS resolution configuration: some non-default DNS servers have been specified")
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", fwInstanceID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = utils.RemoveNil(buildDNSResolutionConfigurationBodyParams(d, defaultDNSServers))
	resp, err := client.Request("PUT", path, &opt)
	if err != nil {
		return diag.Errorf("error creating CFW DNS resolution configuration: %s", err)
	}

	_, err = utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.IsNewResource() {
		d.SetId(fwInstanceID)
	}

	return resourceDNSResolutionRead(ctx, d, meta)
}

func readDNSServers(client *golangsdk.ServiceClient, id string) (interface{}, error) {
	httpUrl := "v1/{project_id}/dns/servers"
	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += "?offset=0&limit=100"
	path += fmt.Sprintf("&fw_instance_id=%s", id)

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		path,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return nil, err
	}

	respJson, err := json.Marshal(resp)
	if err != nil {
		return nil, fmt.Errorf("error marshaling CFW DNS resolution configuration: %s", err)
	}

	var respBody interface{}
	err = json.Unmarshal(respJson, &respBody)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling CFW DNS resolution configuration: %s", err)
	}

	return respBody, nil
}

func buildDNSResolutionConfigurationBodyParams(d *schema.ResourceData, defaultDNSServers []interface{}) map[string]interface{} {
	userDefaultServers := d.Get("default_dns_servers").(*schema.Set).List()
	userCustomServers := d.Get("custom_dns_servers").(*schema.Set).List()
	differenceServers := difference(userDefaultServers, defaultDNSServers)

	return map[string]interface{}{
		"dns_server":               buildDNSServersBodyParams(userDefaultServers, differenceServers, userCustomServers),
		"health_check_domain_name": utils.ValueIgnoreEmpty(d.Get("health_check_domain_name")),
	}
}

func difference(userDefaultServers, defaultDNSServers []interface{}) []interface{} {
	userDefaultServersSet := make(map[string]struct{})
	for _, userServer := range userDefaultServers {
		userDefaultServersSet[userServer.(string)] = struct{}{}
	}

	result := make([]interface{}, 0, len(userDefaultServers))
	for _, defaultServer := range defaultDNSServers {
		if _, ok := userDefaultServersSet[defaultServer.(string)]; !ok {
			result = append(result, defaultServer)
		}
	}
	return result
}

func buildDNSServersBodyParams(userDefaultServers, differenceServers, userCustomServers []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(userCustomServers)+len(differenceServers)+len(userCustomServers))
	for _, server := range userDefaultServers {
		result = append(result, map[string]interface{}{
			"server_ip":     server,
			"is_applied":    1,
			"is_customized": 0,
		})
	}
	for _, server := range differenceServers {
		result = append(result, map[string]interface{}{
			"server_ip":     server,
			"is_applied":    0,
			"is_customized": 0,
		})
	}
	for _, server := range userCustomServers {
		result = append(result, map[string]interface{}{
			"server_ip":     server,
			"is_applied":    1,
			"is_customized": 1,
		})
	}
	return result
}

func resourceDNSResolutionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	respBody, err := readDNSServers(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving CFW DNS resolution configuration",
		)
	}

	dnsServers := utils.PathSearch("data[?is_applied==`1`]", respBody, make([]interface{}, 0)).([]interface{})
	if len(dnsServers) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CFW DNS resolution configuration")
	}

	defaultDNSServers := utils.PathSearch("[?is_customized==`0`].server_ip", dnsServers, make([]interface{}, 0)).([]interface{})
	customDNSServers := utils.PathSearch("[?is_customized==`1`].server_ip", dnsServers, make([]interface{}, 0)).([]interface{})

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("fw_instance_id", d.Id()),
		d.Set("default_dns_servers", defaultDNSServers),
		d.Set("custom_dns_servers", customDNSServers),
		d.Set("health_check_domain_name", utils.PathSearch("[0].health_check_domain_name", dnsServers, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDNSResolutionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl      = "v1/{project_id}/dns/servers"
		product      = "cfw"
		fwInstanceID = d.Get("fw_instance_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	respBody, err := readDNSServers(client, fwInstanceID)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW DNS resolution configuration",
		)
	}

	expression := "data[?is_customized==`0`].server_ip"
	servers := utils.PathSearch(expression, respBody, make([]interface{}, 0)).([]interface{})
	if len(servers) == 0 {
		return diag.Errorf("error deleting CFW DNS resolution configuration: the default DNS servers can not be found")
	}

	appliedServers := utils.PathSearch("data[?is_applied==`1`]", respBody, make([]interface{}, 0)).([]interface{})
	if len(appliedServers) == 0 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving CFW DNS resolution configuration")
	}

	path := client.Endpoint + httpUrl
	path = strings.ReplaceAll(path, "{project_id}", client.ProjectID)
	path += fmt.Sprintf("?fw_instance_id=%s", fwInstanceID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	opt.JSONBody = buildDeleteDNSConfigurationBodyParams(servers)
	_, err = client.Request("PUT", path, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting CFW DNS resolution configuration",
		)
	}

	return nil
}

func buildDeleteDNSConfigurationBodyParams(servers []interface{}) map[string]interface{} {
	return map[string]interface{}{
		"dns_server":               buildDeleteDNSServersBodyParams(servers),
		"health_check_domain_name": "www.huaweicloud.com",
	}
}

func buildDeleteDNSServersBodyParams(servers []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(servers))
	for _, server := range servers {
		result = append(result, map[string]interface{}{
			"server_ip":     server.(string),
			"is_applied":    0,
			"is_customized": 0,
		})
	}
	return result
}
