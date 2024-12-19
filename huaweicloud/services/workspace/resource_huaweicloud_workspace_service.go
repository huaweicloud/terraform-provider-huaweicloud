package workspace

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/workspace/v2/services"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

type AuthenticationType string

const (
	// The authentication type of Workspace service.
	AuthTypeLocal    AuthenticationType = "LITE_AS"
	AuthTypeAdDomain AuthenticationType = "LOCAL_AD"
)

func adDomainSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"admin_account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"active_domain_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"active_domain_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"standby_domain_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standby_domain_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"active_dns_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"standby_dns_ip": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"delete_computer_object": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func securityGroupSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func assistAuthSchemaResource() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"enable": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"receive_mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"auth_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"app_secret": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"auth_server_access_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cert_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"rule_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

// @API Workspace GET /v2/{project_id}/assist-auth-config/method-config
// @API Workspace PUT /v2/{project_id}/assist-auth-config/method-config
// @API Workspace DELETE /v2/{project_id}/workspaces
// @API Workspace GET /v2/{project_id}/workspaces
// @API Workspace POST /v2/{project_id}/workspaces
// @API Workspace PUT /v2/{project_id}/workspaces
// @API Workspace GET /v2/{project_id}/workspaces/lock-status
// @API Workspace PUT /v2/{project_id}/workspaces/lock-status
func ResourceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceCreate,
		ReadContext:   resourceServiceRead,
		UpdateContext: resourceServiceUpdate,
		DeleteContext: resourceServiceDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Update: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Importer: &schema.ResourceImporter{
			StateContext: resourceServiceImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpc_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"network_ids": {
				Type:     schema.TypeList,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"access_mode": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"INTERNET", "DEDICATED", "BOTH",
				}, false),
			},
			"auth_type": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(AuthTypeLocal), string(AuthTypeAdDomain),
				}, false),
				Default: string(AuthTypeLocal),
			},
			"ad_domain": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem:     adDomainSchemaResource(),
			},
			"enterprise_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"internet_access_port": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"dedicated_subnets": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 5,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"management_subnet_cidr": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"lock_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"internet_access_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"infrastructure_security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityGroupSchemaResource(),
			},
			"desktop_security_group": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     securityGroupSchemaResource(),
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"otp_config_info": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     assistAuthSchemaResource(),
			},
			"is_locked": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"lock_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"lock_reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func isDeleteObject(isDelete bool) int {
	if isDelete {
		return 1
	}
	return 0
}

func buildServiceAuthConfiguration(d *schema.ResourceData) *services.Domain {
	authType := d.Get("auth_type").(string)
	if authType == string(AuthTypeLocal) {
		return &services.Domain{
			Type: authType,
		}
	}
	return buildServiceAdDomain(d.Get("ad_domain").([]interface{}))
}

func buildServiceAdDomain(adDomains []interface{}) *services.Domain {
	if len(adDomains) < 1 {
		return nil
	}

	domain := adDomains[0].(map[string]interface{})
	result := services.Domain{
		Type:                 string(AuthTypeAdDomain),
		Name:                 domain["name"].(string),
		AdminAccount:         domain["admin_account"].(string),
		Password:             domain["password"].(string),
		ActiveDomainIp:       domain["active_domain_ip"].(string),
		AcitveDomainName:     domain["active_domain_name"].(string),
		StandyDomainIp:       domain["standby_domain_ip"].(string),
		StandyDomainName:     domain["standby_domain_name"].(string),
		ActiveDnsIp:          domain["active_dns_ip"].(string),
		StandyDnsIp:          domain["standby_dns_ip"].(string),
		DeleteComputerObject: utils.Int(isDeleteObject(domain["delete_computer_object"].(bool))),
	}

	return &result
}

func buildServiceNetworkIds(networkIds []interface{}) []services.Subnet {
	if len(networkIds) < 1 {
		return nil
	}

	result := make([]services.Subnet, len(networkIds))
	for i, networkId := range networkIds {
		result[i] = services.Subnet{
			NetworkId: networkId.(string),
		}
	}

	return result
}

func buildServiceCreateOpts(d *schema.ResourceData) services.CreateOpts {
	return services.CreateOpts{
		AdDomain:             buildServiceAuthConfiguration(d),
		VpcId:                d.Get("vpc_id").(string),
		Subnets:              buildServiceNetworkIds(d.Get("network_ids").([]interface{})),
		AccessMode:           d.Get("access_mode").(string),
		EnterpriseId:         d.Get("enterprise_id").(string),
		DedicatedSubnets:     strings.Join(utils.ExpandToStringList(d.Get("dedicated_subnets").([]interface{})), ";"),
		ManagementSubnetCidr: d.Get("management_subnet_cidr").(string),
	}
}

func refreshServiceStatusFunc(client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := services.Get(client)
		if err != nil {
			return resp, "", err
		}

		return resp, resp.Status, nil
	}
}

func waitForServiceCreateCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration) (string,
	error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PREPARING", "SUBSCRIBING"},
		Target:       []string{"SUBSCRIBED"},
		Refresh:      refreshServiceStatusFunc(client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 15 * time.Second,
	}

	resp, err := stateConf.WaitForStateContext(ctx)
	return resp.(*services.Service).ID, err
}

func resourceServiceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	createOpts := buildServiceCreateOpts(d)
	_, err = services.Create(client, createOpts)
	if err != nil {
		return diag.Errorf("error creating Workspace service: %s", err)
	}
	serviceId, err := waitForServiceCreateCompleted(ctx, client, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("an error occurred while registering the service: %s", err)
	}

	d.SetId(serviceId)

	if _, ok := d.GetOk("internet_access_port"); ok {
		err = updateServiceSubnetIds(ctx, client, d)
		if err != nil {
			return diag.Errorf("error configuring access port: %s", err)
		}
	}

	return resourceServiceRead(ctx, d, meta)
}

func parseIsDeleteObject(isDelete string) bool {
	return isDelete == "1"
}

func flattenServiceAdDomain(d *schema.ResourceData, domain services.DomainResp) []map[string]interface{} {
	if domain == (services.DomainResp{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"name":                   domain.Name,
			"admin_account":          domain.AdminAccount,
			"password":               d.Get("ad_domain.0.password"),
			"active_domain_ip":       domain.ActiveDomainIp,
			"active_domain_name":     domain.AcitveDomainName,
			"standby_domain_ip":      domain.StandyDomainIp,
			"standby_domain_name":    domain.StandyDomainName,
			"active_dns_ip":          domain.ActiveDnsIp,
			"standby_dns_ip":         domain.StandyDnsIp,
			"delete_computer_object": parseIsDeleteObject(domain.DeleteComputerObject),
		},
	}
}

func flattenServiceNetworkIds(networks []services.Subnet) []interface{} {
	if len(networks) < 1 {
		return nil
	}

	result := make([]interface{}, len(networks))
	for i, subnet := range networks {
		result[i] = subnet.NetworkId
	}
	return result
}

func flattenServiceServiceGroup(secgroup services.SecurityGroup) []map[string]interface{} {
	if secgroup == (services.SecurityGroup{}) {
		return nil
	}

	return []map[string]interface{}{
		{
			"id":   secgroup.ID,
			"name": secgroup.Name,
		},
	}
}

func getAssistAuthConfig(client *golangsdk.ServiceClient) ([]map[string]interface{}, error) {
	resp, err := services.GetAuthConfig(client)
	if err != nil {
		return nil, fmt.Errorf("error getting the auxiliary authentication configuration details: %s", err)
	}

	authConfig := resp.OptConfigInfo
	if !authConfig.Enable {
		return nil, nil
	}

	result := []map[string]interface{}{
		{
			"enable":                  authConfig.Enable,
			"receive_mode":            authConfig.ReceiveMode,
			"auth_url":                authConfig.AuthUrl,
			"app_id":                  authConfig.AppId,
			"app_secret":              authConfig.AppSecrte,
			"auth_server_access_mode": authConfig.AuthServerAccessMode,
			"cert_content":            authConfig.CertContent,
			"rule_type":               authConfig.ApplyRule.RuleType,
			"rule":                    authConfig.ApplyRule.Rule,
		},
	}
	return result, err
}

func getLockStatus(client *golangsdk.ServiceClient) (*services.LockStatusResp, error) {
	resp, err := services.GetLockStatus(client)
	if err != nil {
		return nil, fmt.Errorf("error retrieving the lock status of Workspace service: %s", err)
	}

	return resp, nil
}

func resourceServiceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	client, err := conf.WorkspaceV2Client(region)
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	resp, err := services.Get(client)
	if err != nil {
		return diag.Errorf("error retrieving resource details of Workspace service: %s", err)
	}
	if resp.Status == "CLOSED" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("ad_domain", flattenServiceAdDomain(d, resp.AdDomain)),
		d.Set("auth_type", resp.AdDomain.Type),
		d.Set("vpc_id", resp.VpcId),
		d.Set("network_ids", flattenServiceNetworkIds(resp.SubnetIds)),
		d.Set("access_mode", resp.AccessMode),
		d.Set("enterprise_id", resp.EnterpriseId),
		d.Set("dedicated_subnets", strings.Split(resp.DedicatedSubnets, ";")),
		d.Set("management_subnet_cidr", resp.ManagementSubentCidr),
		d.Set("infrastructure_security_group", flattenServiceServiceGroup(resp.InfrastructureSecurityGroup)),
		d.Set("desktop_security_group", flattenServiceServiceGroup(resp.DesktopSecurityGroup)),
		d.Set("status", resp.Status),
	)

	if resp.InternetAccessPort != "" {
		if portNum, err := strconv.Atoi(resp.InternetAccessPort); err == nil {
			mErr = multierror.Append(mErr,
				d.Set("internet_access_port", portNum),
				d.Set("internet_access_address", resp.InternetAccessAddress),
			)
		} else {
			log.Printf("[WARN] the internet access port cannot convert to number")
		}
	}

	configInfo, err := getAssistAuthConfig(client)
	if err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		mErr = multierror.Append(mErr, d.Set("otp_config_info", configInfo))
	}

	lockResp, err := getLockStatus(client)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(mErr,
		d.Set("is_locked", lockResp.IsLocked),
		d.Set("lock_time", lockResp.LockTime),
		d.Set("lock_reason", lockResp.LockReason),
	)

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting service fields: %s", err)
	}
	return nil
}

func doingUpdate(ctx context.Context, client *golangsdk.ServiceClient, opts services.UpdateOpts,
	timeout time.Duration) error {
	resp, err := services.Update(client, opts)
	if err != nil {
		return err
	}

	_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, timeout)
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)

	return nil
}

func updateServiceConnection(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	log.Printf("[DEBUG] start to update the service connection")
	opts := services.UpdateOpts{
		AdDomain:         buildServiceAdDomain(d.Get("ad_domain").([]interface{})),
		AccessMode:       d.Get("access_mode").(string),
		DedicatedSubnets: strings.Join(utils.ExpandToStringList(d.Get("dedicated_subnets").([]interface{})), ";"),
	}
	return doingUpdate(ctx, client, opts, d.Timeout(schema.TimeoutUpdate))
}

func updateServiceSubnetIds(_ context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	log.Printf("[DEBUG] start updating the network ID list of service")
	opts := services.UpdateOpts{
		Subnets: utils.ExpandToStringList(d.Get("network_ids").([]interface{})),
	}
	// Updating subnet configuration will not return job ID.
	_, err := services.Update(client, opts)
	if err != nil {
		return err
	}
	return nil
}

func updateServiceInternetAccess(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	log.Printf("[DEBUG] start to update the internet access port of service")
	opts := services.UpdateOpts{
		InternetAccessPort: strconv.Itoa(d.Get("internet_access_port").(int)),
	}
	return doingUpdate(ctx, client, opts, d.Timeout(schema.TimeoutUpdate))
}

func updateServiceEnterpriseId(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	log.Printf("[DEBUG] start to update the enterprise ID of service")
	opts := services.UpdateOpts{
		EnterpriseId: d.Get("enterprise_id").(string),
	}
	return doingUpdate(ctx, client, opts, d.Timeout(schema.TimeoutUpdate))
}

func buildAuthConfig(configInfo map[string]interface{}) *services.OtpConfigInfo {
	result := services.OtpConfigInfo{
		Enable:               utils.Bool(configInfo["enable"].(bool)),
		ReceiveMode:          configInfo["receive_mode"].(string),
		AuthUrl:              configInfo["auth_url"].(string),
		AppId:                configInfo["app_id"].(string),
		AppSecrte:            configInfo["app_secret"].(string),
		AuthServerAccessMode: configInfo["auth_server_access_mode"].(string),
		CertContent:          configInfo["cert_content"].(string),
		ApplyRule: &services.ApplyRule{
			RuleType: configInfo["rule_type"].(string),
			Rule:     configInfo["rule"].(string),
		},
	}

	return &result
}

func updateAssistAuthConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	authConfig := d.Get("otp_config_info").([]interface{})
	if len(authConfig) < 1 {
		return nil
	}

	opts := services.UpdateAuthConfigOpts{
		AuthType:      "OTP",
		OptConfigInfo: buildAuthConfig(authConfig[0].(map[string]interface{})),
	}
	return services.UpdateAssistAuthConfig(client, opts)
}

func unlockService(ctx context.Context, client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	opts := services.UnlockOpts{
		OperateType: "unlock",
	}
	resp, err := services.UnlockService(client, opts)
	if err != nil {
		return fmt.Errorf("error unLocking of the Workspace service: %s", err)
	}

	_, err = waitForWorkspaceJobCompleted(ctx, client, resp.JobId, d.Timeout(schema.TimeoutUpdate))
	if err != nil {
		return fmt.Errorf("error waiting for the job (%s) completed: %s", resp.JobId, err)
	}
	log.Printf("[DEBUG] The job (%s) has been completed", resp.JobId)

	return err
}

func resourceServiceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	if d.HasChanges("ad_domain", "access_mode", "dedicated_subnets") {
		if err = updateServiceConnection(ctx, client, d); err != nil {
			return diag.Errorf("error updating connection parameters of service: %s", err)
		}
	}
	if d.HasChange("network_ids") {
		if err = updateServiceSubnetIds(ctx, client, d); err != nil {
			return diag.Errorf("error updating subnet list of service: %s", err)
		}
	}
	if d.HasChange("internet_access_port") {
		if err = updateServiceInternetAccess(ctx, client, d); err != nil {
			return diag.Errorf("error updating internet access port of service: %s", err)
		}
	}
	if d.HasChange("enterprise_id") {
		if err = updateServiceEnterpriseId(ctx, client, d); err != nil {
			return diag.Errorf("error updating enterprise ID of service: %s", err)
		}
	}

	if d.HasChanges("otp_config_info") {
		if err = updateAssistAuthConfig(client, d); err != nil {
			return diag.Errorf("error updating authentication config parameters of service: %s", err)
		}
	}

	if d.HasChanges("lock_enabled") {
		lockEnabled := d.Get("lock_enabled").(bool)
		// If the current service is not locked, this action is not required.
		if !lockEnabled {
			return nil
		}

		if err = unlockService(ctx, client, d); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceServiceRead(ctx, d, meta)
}

func refreshServiceClosableStatusFunc(client *golangsdk.ServiceClient) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := services.Get(client)
		if err != nil {
			return resp, "ERROR", err
		}

		if !resp.Closable {
			return resp, "PENDING", nil
		}
		return resp, "COMPLETE", nil
	}
}

func waitForServiceClosableReady(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETE"},
		Refresh:      refreshServiceClosableStatusFunc(client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func waitForServiceDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"DEREGISTERING"},
		Target:       []string{"CLOSED"},
		Refresh:      refreshServiceStatusFunc(client),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceServiceDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace v2 client: %s", err)
	}

	err = waitForServiceClosableReady(ctx, client, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("The current service is not allowed to be deleted (the value of the closable attribute is false): %s", err)
	}

	_, err = services.Delete(client)
	if err != nil {
		return diag.Errorf("error unregistring service (%s): %s", d.Id(), err)
	}
	err = waitForServiceDeleteCompleted(ctx, client, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return diag.Errorf("an error occurred while unregistering the service: %s", err)
	}

	return nil
}

func resourceServiceImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if !utils.IsUUID(d.Id()) {
		conf := meta.(*config.Config)
		client, err := conf.WorkspaceV2Client(conf.GetRegion(d))
		if err != nil {
			return nil, fmt.Errorf("error creating Workspace v2 client: %s", err)
		}

		resp, err := services.Get(client)
		if err != nil {
			return nil, fmt.Errorf("error retrieving Workspace service detail: %s", err)
		}

		// Refresh the service ID.
		d.SetId(resp.ID)
	}

	return []*schema.ResourceData{d}, nil
}
