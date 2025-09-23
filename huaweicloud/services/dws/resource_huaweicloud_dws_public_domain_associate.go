package dws

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS POST /v1.0/{project_id}/clusters/{cluster_id}/dns
// @API DWS GET /v1.0/{project_id}/clusters/{cluster_id}/dns
// @API DWS DELETE /v1.0/{project_id}/clusters/{cluster_id}/dns
func ResourcePublicDomainAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePublicDomainAssociateCreate,
		ReadContext:   resourcePublicDomainAssociateRead,
		UpdateContext: resourcePublicDomainAssociateUpdate,
		DeleteContext: resourcePublicDomainAssociateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourcePublicDomainAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the DWS cluster ID to which the public domain name to be associated belongs.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the public domain name.`,
			},
			"ttl": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `Specifies cache period of the SOA record set, in seconds.`,
			},
		},
	}
}

func resourcePublicDomainAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/clusters/{cluster_id}/dns"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", d.Get("cluster_id").(string))
	createOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildPublicDomainAssociateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating public domain name for DWS cluster: %s", err)
	}

	d.SetId(d.Get("domain_name").(string))

	return resourcePublicDomainAssociateRead(ctx, d, meta)
}

func buildPublicDomainAssociateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name": d.Get("domain_name"),
		"type": "public",
		"ttl":  utils.ValueIgnoreEmpty(d.Get("ttl")),
	}
}

func GetDomainNameInfoByClusterId(client *golangsdk.ServiceClient, clusterId string) (interface{}, error) {
	resp, err := GetClusterInfoByClusterId(client, clusterId)
	if err != nil {
		return nil, err
	}

	domainNameInfo := utils.PathSearch("cluster.public_endpoints[*].public_connect_info|[0]", resp, nil)
	// The format after binding a domain name is `${domian_name}.dws.huaweicloud.com.`, eg. `terraform.dws.huaweicloud.com.`
	// The format after unbinding a domain name is `${eip_address}:${dws_cluster_port}`, eg. `120.46.40.224:8000`
	// The net.ParseIP method is used to determine whether the IP is legal.
	// If it is legal, it returns the IP, otherwise it returns nil.
	if domainNameInfo == nil || net.ParseIP(strings.Split(domainNameInfo.(string), ":")[0]) != nil {
		return nil, golangsdk.ErrDefault404{}
	}
	return domainNameInfo, nil
}

func resourcePublicDomainAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	domainNameInfo, err := GetDomainNameInfoByClusterId(client, d.Get("cluster_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving public domain name")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", parsePublicDomainName(domainNameInfo)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parsePublicDomainName(domainNameExpression interface{}) interface{} {
	if domainNameExpression == nil {
		return nil
	}
	// The format of public connection information is "${domain_name}.dws.huaweicloud.com."
	return strings.Split(domainNameExpression.(string), ".")[0]
}

func resourcePublicDomainAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dws", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	if d.HasChanges("domain_name", "ttl") {
		opt := buildPublicDomainAssociateBodyParams(d)
		if _, err = UpdateDomainName(client, d.Get("cluster_id").(string), opt); err != nil {
			return diag.Errorf("error updating public domain name for the DWS cluster: %s", err)
		}
	}

	return resourcePublicDomainAssociateRead(ctx, d, meta)
}

// UpdateDomainName is a method used to update the public and private domain name of the DWS cluster.
func UpdateDomainName(client *golangsdk.ServiceClient, clusterId string, opt map[string]interface{}) (*http.Response, error) {
	httpUrl := "v1.0/{project_id}/clusters/{cluster_id}/dns"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterId)
	updateOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(opt),
	}
	return client.Request("PUT", updatePath, &updateOpt)
}

func resourcePublicDomainAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/clusters/{cluster_id}/dns?type=public"
	)

	client, err := cfg.NewServiceClient("dws", region)
	if err != nil {
		return diag.Errorf("error creating DWS client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", d.Get("cluster_id").(string))
	deleteOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	// "DWS.0047": Cluster ID does not exist, the status code is 404.
	// "DWS.5212": Resource does not exist, the status code is 400.
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", "DWS.5212"),
			fmt.Sprintf("error unbinding public domain name (%s) for the DWS cluster", d.Id()))
	}

	return nil
}

func resourcePublicDomainAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<cluster_id>/<domain_name>', but got '%s'",
			importId)
	}

	d.Set("cluster_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
