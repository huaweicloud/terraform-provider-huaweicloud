package workspace

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API Workspace GET /v2/{project_id}/workspaces
// @API Workspace GET /v2/{project_id}/assist-auth-config/method-config
// @API Workspace GET /v2/{project_id}/workspaces/lock-status
func DataSourceService() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceServiceRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the environments are located.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The VPC ID to which the service belongs.`,
			},
			"network_ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The network ID list of subnets that the service have.`,
			},
			"access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The access mode of Workspace service.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication type of Workspace service.`,
			},
			"ad_domain": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAdDomainSchema(),
				Description: `The configuration of AD domain.`,
			},
			"enterprise_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The enterprise ID.`,
			},
			"internet_access_port": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The internet access port.`,
			},
			"internet_access_address": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The internet access address.`,
			},
			"dedicated_subnets": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The subnet segments of the dedicated access.`,
			},
			"management_subnet_cidr": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The subnet segment of the management component.`,
			},
			"infrastructure_security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataServiceSecurityGroupSchema(),
				Description: `The management component security group automatically created under the specified VPC
after service is registered.`,
			},
			"desktop_security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dataServiceSecurityGroupSchema(),
				Description: `The desktop security group automatically created under the specified VPC after the service
is registered.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The current status of the Workspace service.`,
			},
			"otp_config_info": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        dataServiceAssistAuthSchema(),
				Description: `The configuration of auxiliary authentication.`,
			},
			"is_locked": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `Whether the service is locked.`,
			},
			"lock_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the service is locked.`,
			},
			"lock_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the service is locked.`,
			},
		},
	}
}

func dataServiceAdDomainSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name.`,
			},
			"admin_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain administrator account.`,
			},
			"active_domain_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of primary domain controller.`,
			},
			"active_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the primary domain controller.`,
			},
			"active_dns_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The primary DNS IP address.`,
			},
			"standby_domain_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The IP address of standby domain controller.`,
			},
			"standby_domain_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the standby domain controller.`,
			},
			"standby_dns_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The standby DNS IP address.`,
			},
			"delete_computer_object": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to delete the corresponding computer object on AD while deleting the desktop.`,
			},
		},
	}
}

func dataServiceSecurityGroupSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The security group ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The security group name.`,
			},
		},
	}
}

func dataServiceAssistAuthSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether to enable auxiliary authentication.`,
			},
			"receive_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The verification code receiving mode.`,
			},
			"auth_url": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The auxiliary authentication server address.`,
			},
			"app_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The auxiliary authentication server access account.`,
			},
			"app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication service access password.`,
			},
			"auth_server_access_mode": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication service access mode.`,
			},
			"cert_content": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate content, in PEM format.`,
			},
			"rule_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication application object type.`,
			},
			"rule": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication application object.`,
			},
		},
	}
}

func showService(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/workspaces"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenDataServiceAdDomain(adDomain interface{}) []map[string]interface{} {
	if adDomain == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":                   utils.PathSearch("domain_name", adDomain, nil),
			"admin_account":          utils.PathSearch("domain_admin_account", adDomain, nil),
			"active_domain_ip":       utils.PathSearch("active_domain_ip", adDomain, nil),
			"active_domain_name":     utils.PathSearch("active_domain_name", adDomain, nil),
			"active_dns_ip":          utils.PathSearch("active_dns_ip", adDomain, nil),
			"standby_domain_ip":      utils.PathSearch("standby_domain_ip", adDomain, nil),
			"standby_domain_name":    utils.PathSearch("standby_domain_name", adDomain, nil),
			"standby_dns_ip":         utils.PathSearch("standby_dns_ip", adDomain, nil),
			"delete_computer_object": parseIsDeleteObject(utils.PathSearch("delete_computer_object", adDomain, "").(string)),
		},
	}
}

func parseDataServiceInternetAccessPort(portStr string) int {
	var (
		portNum int
		err     error
	)
	if portStr == "" {
		return portNum
	}

	portNum, err = strconv.Atoi(portStr)
	if err != nil {
		log.Printf("[ERROR] error converting internet access port (value: %s) from string to integer: %s", portStr, err)
	}
	return portNum
}

func flattenDataServiceSecurityGroup(securityGroup interface{}) []map[string]interface{} {
	if securityGroup == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":   utils.PathSearch("id", securityGroup, nil),
			"name": utils.PathSearch("name", securityGroup, nil),
		},
	}
}

func showAssistAuthConfig(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/assist-auth-config/method-config"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func flattenDataServiceAssistAuthConfig(assistAuthConfig interface{}) []map[string]interface{} {
	if assistAuthConfig == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"enable":                  utils.PathSearch("enable", assistAuthConfig, nil),
			"receive_mode":            utils.PathSearch("receive_mode", assistAuthConfig, nil),
			"auth_url":                utils.PathSearch("auth_url", assistAuthConfig, nil),
			"app_id":                  utils.PathSearch("app_id", assistAuthConfig, nil),
			"app_secret":              utils.PathSearch("app_secret", assistAuthConfig, nil),
			"auth_server_access_mode": utils.PathSearch("auth_server_access_mode", assistAuthConfig, nil),
			"cert_content":            utils.PathSearch("cert_content", assistAuthConfig, nil),
			"rule":                    utils.PathSearch("apply_rule.rule", assistAuthConfig, nil),
			"rule_type":               utils.PathSearch("apply_rule.rule_type", assistAuthConfig, nil),
		},
	}
}

func showServiceLockStatus(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v2/{project_id}/workspaces/lock-status"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}
	return utils.FlattenResponse(requestResp)
}

func dataSourceServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	service, err := showService(client)
	if err != nil {
		return diag.Errorf("error querying service: %s", err)
	}

	if serviceId := utils.PathSearch("id", service, "").(string); serviceId != "" {
		d.SetId(serviceId)
	} else {
		uuid, err := uuid.GenerateUUID()
		if err != nil {
			return diag.Errorf("unable to generate ID: %s", err)
		}
		d.SetId(uuid)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("vpc_id", utils.PathSearch("vpc_id", service, nil)),
		d.Set("network_ids", utils.PathSearch("subnet_ids[*].subnet_id", service, make([]interface{}, 0))),
		d.Set("access_mode", utils.PathSearch("access_mode", service, nil)),
		d.Set("auth_type", utils.PathSearch("ad_domains.domain_type", service, nil)),
		d.Set("ad_domain", flattenDataServiceAdDomain(utils.PathSearch("ad_domains", service, nil))),
		d.Set("enterprise_id", utils.PathSearch("enterprise_id", service, nil)),
		d.Set("internet_access_address", utils.PathSearch("internet_access_address", service, nil)),
		d.Set("internet_access_port", parseDataServiceInternetAccessPort(utils.PathSearch("internet_access_port", service, "").(string))),
		d.Set("dedicated_subnets", strings.Split(utils.PathSearch("dedicated_subnets", service, "").(string), ";")),
		d.Set("management_subnet_cidr", utils.PathSearch("management_subnet_cidr", service, nil)),
		d.Set("infrastructure_security_group", flattenDataServiceSecurityGroup(utils.PathSearch("infrastructure_security_group", service, nil))),
		d.Set("desktop_security_group", flattenDataServiceSecurityGroup(utils.PathSearch("desktop_security_group", service, nil))),
		d.Set("status", utils.PathSearch("status", service, nil)),
	)

	assistAuthConfig, err := showAssistAuthConfig(client)
	if err != nil {
		log.Printf("[ERROR] error querying assist auth configuration of the Workspace service: %s", err)
	} else {
		mErr = multierror.Append(mErr, d.Set("otp_config_info", flattenDataServiceAssistAuthConfig(assistAuthConfig)))
	}

	lockStatus, err := showServiceLockStatus(client)
	if err != nil {
		log.Printf("[ERROR] error querying assist auth configuration of the Workspace service: %s", err)
	} else {
		mErr = multierror.Append(mErr,
			d.Set("is_locked", utils.PathSearch("is_locked", lockStatus, nil)),
			d.Set("lock_time", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(
				utils.PathSearch("lock_time", lockStatus, "").(string),
				"2006-01-02 15:04:05",
			)/1000, false)),
			d.Set("lock_reason", utils.PathSearch("lock_reason", lockStatus, nil)),
		)
	}

	return diag.FromErr(mErr.ErrorOrNil())
}
