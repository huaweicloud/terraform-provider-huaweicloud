package vpn

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var accessPolicyNonUpdatableParams = []string{"vpn_server_id"}

// @API VPN POST /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies
// @API VPN GET /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}
// @API VPN PUT /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}
// @API VPN DELETE /v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}
func ResourceAccessPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAccessPolicyCreate,
		UpdateContext: resourceAccessPolicyUpdate,
		ReadContext:   resourceAccessPolicyRead,
		DeleteContext: resourceAccessPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAccessPolicyImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(accessPolicyNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"vpn_server_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The VPN server ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The access policy name.`,
			},
			"user_group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The user group ID.`,
			},
			"dest_ip_cidrs": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The list of destination IP CIDRs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the access policy.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"user_group_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user group name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The update time.`,
			},
		},
	}
}

func resourceAccessPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		createAccessPolicyHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies"
		createAccessPolicyProduct = "vpn"
	)
	createAccessPolicyClient, err := conf.NewServiceClient(createAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	createAccessPolicyPath := createAccessPolicyClient.Endpoint + createAccessPolicyHttpUrl
	createAccessPolicyPath = strings.ReplaceAll(createAccessPolicyPath, "{project_id}", createAccessPolicyClient.ProjectID)
	createAccessPolicyPath = strings.ReplaceAll(createAccessPolicyPath, "{vpn_server_id}", d.Get("vpn_server_id").(string))

	createAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createAccessPolicyOpt.JSONBody = utils.RemoveNil(buildCreateAccessPolicyBodyParams(d))
	createAccessPolicyResp, err := createAccessPolicyClient.Request("POST", createAccessPolicyPath, &createAccessPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating VPN access policy: %s", err)
	}

	createAccessPolicyRespBody, err := utils.FlattenResponse(createAccessPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("access_policy.id", createAccessPolicyRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find access policy ID in API response")
	}
	d.SetId(id)
	// The creation interface is asynchronous.
	// If the access policy information disappears, then the creation fails.
	// Wait for a while to check if the creation is successful.
	// lintignore:R018
	time.Sleep(30 * time.Second)

	return resourceAccessPolicyRead(ctx, d, meta)
}

func buildCreateAccessPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_policy": map[string]interface{}{
			"name":          d.Get("name"),
			"user_group_id": d.Get("user_group_id"),
			"dest_ip_cidrs": d.Get("dest_ip_cidrs").(*schema.Set).List(),
			"description":   utils.ValueIgnoreEmpty(d.Get("description")),
		},
	}
	return bodyParams
}

func resourceAccessPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var mErr *multierror.Error

	getAccessPolicyProduct := "vpn"
	getAccessPolicyClient, err := conf.NewServiceClient(getAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	getAccessPolicyBody, err := GetAccessPolicy(getAccessPolicyClient, serverId, id)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPN access policy")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("access_policy.name", getAccessPolicyBody, nil)),
		d.Set("description", utils.PathSearch("access_policy.description", getAccessPolicyBody, nil)),
		d.Set("user_group_id", utils.PathSearch("access_policy.user_group_id", getAccessPolicyBody, nil)),
		d.Set("dest_ip_cidrs", utils.PathSearch("access_policy.dest_ip_cidrs", getAccessPolicyBody, nil)),
		d.Set("user_group_name", utils.PathSearch("access_policy.user_group_name", getAccessPolicyBody, nil)),
		d.Set("created_at", utils.PathSearch("access_policy.created_at", getAccessPolicyBody, nil)),
		d.Set("updated_at", utils.PathSearch("access_policy.updated_at", getAccessPolicyBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func GetAccessPolicy(client *golangsdk.ServiceClient, serverId, id string) (interface{}, error) {
	getAccessPolicyHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}"
	getAccessPolicyPath := buildAccessPolicyURL(client, getAccessPolicyHttpUrl, serverId, id)

	getAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getAccessPolicyResp, err := client.Request("GET", getAccessPolicyPath, &getAccessPolicyOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getAccessPolicyResp)
}

func resourceAccessPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	updateAccessPolicyClient, err := conf.NewServiceClient("vpn", region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	updateAccessPolicyHttpUrl := "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}"
	updateAccessPolicyPath := buildAccessPolicyURL(updateAccessPolicyClient, updateAccessPolicyHttpUrl, serverId, id)

	updateAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	updateAccessPolicyOpt.JSONBody = buildUpdateAccessPolicyBodyParams(d)
	_, err = updateAccessPolicyClient.Request("PUT", updateAccessPolicyPath, &updateAccessPolicyOpt)
	if err != nil {
		return diag.Errorf("error updating VPN access policy: %s", err)
	}

	return resourceAccessPolicyRead(ctx, d, meta)
}

func buildAccessPolicyURL(client *golangsdk.ServiceClient, urlTemplate, serverId, id string) string {
	url := client.Endpoint + urlTemplate
	url = strings.ReplaceAll(url, "{project_id}", client.ProjectID)
	url = strings.ReplaceAll(url, "{vpn_server_id}", serverId)
	url = strings.ReplaceAll(url, "{policy_id}", id)
	return url
}

func buildUpdateAccessPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"access_policy": map[string]interface{}{
			"name":          d.Get("name"),
			"user_group_id": d.Get("user_group_id"),
			"dest_ip_cidrs": d.Get("dest_ip_cidrs").(*schema.Set).List(),
			"description":   d.Get("description"),
		},
	}
	return bodyParams
}

func resourceAccessPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)

	var (
		deleteAccessPolicyHttpUrl = "v5/{project_id}/p2c-vpn-gateways/vpn-servers/{vpn_server_id}/access-policies/{policy_id}"
		deleteAccessPolicyProduct = "vpn"
	)
	deleteAccessPolicyClient, err := conf.NewServiceClient(deleteAccessPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating VPN client: %s", err)
	}

	serverId := d.Get("vpn_server_id").(string)
	id := d.Id()
	deleteAccessPolicyPath := buildAccessPolicyURL(deleteAccessPolicyClient, deleteAccessPolicyHttpUrl, serverId, id)

	deleteAccessPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = deleteAccessPolicyClient.Request("DELETE", deleteAccessPolicyPath, &deleteAccessPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting VPN access policy")
	}

	return nil
}

func resourceAccessPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid ID format, must be <vpn_server_id>/<id>")
	}

	d.Set("vpn_server_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
