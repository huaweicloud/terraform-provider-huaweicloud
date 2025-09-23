package live

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API LIVE PUT /v1/{project_id}/guard/ip
// @API LIVE GET /v1/{project_id}/guard/ip
func ResourceIpAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIpAclCreate,
		ReadContext:   resourceIpAclRead,
		UpdateContext: resourceIpAclUpdate,
		DeleteContext: resourceIpAclDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceIpAclImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ingest or streaming domain name.`,
			},
			"auth_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the authentication mode.`,
			},
			"ip_auth_list": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the blacklist or whitelist IP addresses.`,
			},
		},
	}
}

func buildUpdateIPAddressAclBodyParams(d *schema.ResourceData, authType string) map[string]interface{} {
	return map[string]interface{}{
		"domain":       d.Get("domain_name"),
		"auth_type":    authType,
		"ip_auth_list": d.Get("ip_auth_list"),
	}
}

func updateIPAddressAcl(client *golangsdk.ServiceClient, d *schema.ResourceData, authType string) error {
	requestPath := client.Endpoint + "v1/{project_id}/guard/ip"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildUpdateIPAddressAclBodyParams(d, authType),
	}
	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceIpAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "live"
		authType = d.Get("auth_type").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	resourceID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("error generating Live IP address acl resource ID: %s", err)
	}

	if err := updateIPAddressAcl(client, d, authType); err != nil {
		return diag.Errorf("error creating Live IP address acl: %s", err)
	}

	d.SetId(resourceID)

	return resourceIpAclRead(ctx, d, meta)
}

func ReadIPAddressAcl(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/guard/ip"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?domain=%s", domainName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceIpAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	respBody, err := ReadIPAddressAcl(client, d.Get("domain_name").(string))
	if err != nil {
		// When the domain does not exist, it will respond with a `400` status code.
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "LIVE.100011001"),
			"error retrieving Live IP address acl")
	}

	authType := utils.PathSearch("auth_type", respBody, "").(string)
	if authType == "NONE" {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("domain", respBody, nil)),
		d.Set("auth_type", authType),
		d.Set("ip_auth_list", utils.PathSearch("ip_auth_list", respBody, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIpAclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		product  = "live"
		authType = d.Get("auth_type").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateIPAddressAcl(client, d, authType); err != nil {
		return diag.Errorf("error updating Live IP address acl: %s", err)
	}

	return resourceIpAclRead(ctx, d, meta)
}

func resourceIpAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "live"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateIPAddressAcl(client, d, "NONE"); err != nil {
		return diag.Errorf("error deleting Live IP address acl: %s", err)
	}
	return nil
}

func resourceIpAclImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)
	return []*schema.ResourceData{d}, d.Set("domain_name", importedId)
}
