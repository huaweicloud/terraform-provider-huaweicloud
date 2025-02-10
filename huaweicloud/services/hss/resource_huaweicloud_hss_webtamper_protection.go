package hss

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API HSS POST /v5/{project_id}/webtamper/static/status
// @API HSS POST /v5/{project_id}/webtamper/rasp/status
// @API HSS GET /v5/{project_id}/webtamper/hosts
// @API HSS GET /v5/{project_id}/host-management/hosts
func ResourceWebTamperProtection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWebTamperProtectionCreate,
		ReadContext:   resourceWebTamperProtectionRead,
		UpdateContext: resourceWebTamperProtectionUpdate,
		DeleteContext: resourceWebTamperProtectionDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceWebTamperProtectionImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"host_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"quota_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"is_dynamics_protect": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"host_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"group_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_bit": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rasp_protect_status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"anti_tampering_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"detect_tampering_times": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildOpenWebTamperProtectionQueryParams(epsId string) string {
	if epsId != "" {
		return fmt.Sprintf("?enterprise_project_id=%v", epsId)
	}

	return ""
}

func buildOpenWebTamperProtectionBodyParams(d *schema.ResourceData, hostId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// **true** means enable, **false** means disabled.
		"status":        true,
		"host_id_list":  []string{hostId},
		"resource_id":   utils.StringIgnoreEmpty(d.Get("quota_id").(string)),
		"charging_mode": chargingModePacketCycle,
	}

	return bodyParams
}

func openWebTamperProtection(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData, hostId string) error {
	var (
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	requestPath := client.Endpoint + "v5/{project_id}/webtamper/static/status"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOpenWebTamperProtectionQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         utils.RemoveNil(buildOpenWebTamperProtectionBodyParams(d, hostId)),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)

	return err
}

func buildOpenOrCloseDynamicProtectionBodyParams(d *schema.ResourceData, hostId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status":       d.Get("is_dynamics_protect").(bool),
		"host_id_list": []string{hostId},
	}

	return bodyParams
}

func openOrCloseDynamicProtection(client *golangsdk.ServiceClient, cfg *config.Config, d *schema.ResourceData, hostId string) error {
	var (
		region = cfg.GetRegion(d)
		epsId  = cfg.GetEnterpriseProjectID(d)
	)

	requestPath := client.Endpoint + "v5/{project_id}/webtamper/rasp/status"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOpenWebTamperProtectionQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildOpenOrCloseDynamicProtectionBodyParams(d, hostId),
	}

	_, err := client.Request("POST", requestPath, &requestOpt)

	return err
}

func resourceWebTamperProtectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		hostId  = d.Get("host_id").(string)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, epsId, hostId)
	if checkHostAvailableErr != nil {
		return diag.FromErr(checkHostAvailableErr)
	}

	err = openWebTamperProtection(client, cfg, d, hostId)
	if err != nil {
		return diag.Errorf("error opening HSS web tamper protection: %s", err)
	}

	d.SetId(hostId)

	if d.Get("is_dynamics_protect").(bool) {
		err = openOrCloseDynamicProtection(client, cfg, d, hostId)
		if err != nil {
			return diag.Errorf("error opening HSS dynamic web tamper protection: %s", err)
		}
	}

	return resourceWebTamperProtectionRead(ctx, d, meta)
}

func buildWebTamperProtectionQueryParams(epsId, hostId string) string {
	return fmt.Sprintf("?enterprise_project_id=%v&host_id=%v&protect_status=%v", epsId, hostId, ProtectStatusOpened)
}

func GetWebTamperProtectionHost(client *golangsdk.ServiceClient, region, epsId, hostId string) (interface{}, error) {
	getPath := client.Endpoint + "v5/{project_id}/webtamper/hosts"
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath += buildWebTamperProtectionQueryParams(epsId, hostId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving HSS web tamper protection host: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	hostResp := utils.PathSearch("data_list[0]", getRespBody, nil)
	if hostResp == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return hostResp, nil
}

func resourceWebTamperProtectionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		hostId  = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	hostResp, err := GetWebTamperProtectionHost(client, region, epsId, hostId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "HSS web tamper protection")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("host_id", utils.PathSearch("host_id", hostResp, nil)),
		d.Set("host_name", utils.PathSearch("host_name", hostResp, nil)),
		d.Set("public_ip", utils.PathSearch("public_ip", hostResp, nil)),
		d.Set("private_ip", utils.PathSearch("private_ip", hostResp, nil)),
		d.Set("group_name", utils.PathSearch("group_name", hostResp, nil)),
		d.Set("os_bit", utils.PathSearch("os_bit", hostResp, nil)),
		d.Set("os_type", utils.PathSearch("os_type", hostResp, nil)),
		d.Set("protect_status", utils.PathSearch("protect_status", hostResp, nil)),
		d.Set("rasp_protect_status", utils.PathSearch("rasp_protect_status", hostResp, nil)),
		d.Set("anti_tampering_times", utils.PathSearch("anti_tampering_times", hostResp, nil)),
		d.Set("detect_tampering_times", utils.PathSearch("detect_tampering_times", hostResp, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceWebTamperProtectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		hostId  = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, epsId, hostId)
	if checkHostAvailableErr != nil {
		return diag.FromErr(checkHostAvailableErr)
	}

	if d.HasChange("is_dynamics_protect") {
		err = openOrCloseDynamicProtection(client, cfg, d, hostId)
		if err != nil {
			return diag.Errorf("error updating HSS dynamic web tamper protection: %s", err)
		}
	}

	return resourceWebTamperProtectionRead(ctx, d, meta)
}

func buildCloseWebTamperProtectionBodyParams(hostId string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// **true** means enable, **false** means disabled.
		"status":       false,
		"host_id_list": []string{hostId},
	}

	return bodyParams
}

func resourceWebTamperProtectionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		epsId   = cfg.GetEnterpriseProjectID(d)
		hostId  = d.Id()
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating HSS client: %s", err)
	}

	requestPath := client.Endpoint + "v5/{project_id}/webtamper/static/status"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildOpenWebTamperProtectionQueryParams(epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"region": region},
		JSONBody:         buildCloseWebTamperProtectionBodyParams(hostId),
	}

	_, err = client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		// Repeatedly closing web tamper protection, API will not report errors.
		// If the host does not exist, closing web tamper protection will result in an error as follows:
		// {"error_code": "00000010","error_description": "拒绝访问"}
		// The API documentation does not provide any explanatory information about this error,
		// so the logic of checkDeleted is not added.
		return diag.Errorf("error closing HSS web tamper protection: %s", err)
	}

	return nil
}

func resourceWebTamperProtectionImportState(_ context.Context, d *schema.ResourceData, meta interface{}) (
	[]*schema.ResourceData, error) {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		hostId  = d.Id()
		epsId   = cfg.GetEnterpriseProjectID(d, QueryAllEpsValue)
		product = "hss"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return []*schema.ResourceData{d}, fmt.Errorf("error creating HSS client: %s", err)
	}

	checkHostAvailableErr := checkHostAvailable(client, epsId, hostId)

	return []*schema.ResourceData{d}, checkHostAvailableErr
}
