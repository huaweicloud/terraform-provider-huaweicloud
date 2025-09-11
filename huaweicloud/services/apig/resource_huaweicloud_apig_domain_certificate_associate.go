package apig

import (
	"context"
	"fmt"
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

var domainCertificateAssociateNonUpdatableParams = []string{
	"instance_id",
	"group_id",
	"domain_id",
	"certificate_ids",
}

// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}/certificates/attach
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}/certificates/detach
func ResourceDomainCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainCertificateCreate,
		ReadContext:   resourceDomainCertificateRead,
		UpdateContext: resourceDomainCertificateUpdate,
		DeleteContext: resourceDomainCertificateDelete,

		CustomizeDiff: config.FlexibleForceNew(domainCertificateAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceDomainCertificateAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the domain and certificates are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the dedicated instance to which the domain belongs.",
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the API group to which the domain belongs.",
			},
			"domain_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the domain.",
			},
			"certificate_ids": {
				Type:        schema.TypeSet,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The list of the certificate IDs to associate with the domain.",
			},
			"verified_client_certificate_enabled": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Whether to enable verified client certificate.",
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildAssociateCertificatesToDomainBodyParams(certificateIds []interface{}, verifiedClientCertificateEnabled bool) map[string]interface{} {
	return map[string]interface{}{
		"certificate_ids":                     certificateIds,
		"verified_client_certificate_enabled": verifiedClientCertificateEnabled,
	}
}

func associateCertificatesToDomain(client *golangsdk.ServiceClient, instanceId, groupId, domainId string, certificateIds []interface{},
	verifiedClientCertificateEnabled bool) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}/certificates/attach"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)
	createPath = strings.ReplaceAll(createPath, "{domain_id}", domainId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes:  []int{204},
		JSONBody: buildAssociateCertificatesToDomainBodyParams(certificateIds, verifiedClientCertificateEnabled),
	}

	_, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error associating certificates to domain (%s): %s", domainId, err)
	}

	return nil
}

func resourceDomainCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = associateCertificatesToDomain(
		client,
		instanceId,
		groupId,
		d.Get("domain_id").(string),
		d.Get("certificate_ids").(*schema.Set).List(),
		d.Get("verified_client_certificate_enabled").(bool),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s/%s", instanceId, groupId, d.Get("domain_id").(string)))

	return resourceDomainCertificateRead(ctx, d, meta)
}

func getDomainAssociatedCertificates(client *golangsdk.ServiceClient, instanceId, groupId, domainId string) (interface{}, error) {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{group_id}", groupId)

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
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}

	return utils.PathSearch(fmt.Sprintf("url_domains[?id=='%s']|[0]", domainId), respBody, nil), nil
}

func GetDomainAssociatedCertificatesByDomainId(client *golangsdk.ServiceClient, instanceId, groupId, domainId string) (interface{}, error) {
	domainInfos, err := getDomainAssociatedCertificates(client, instanceId, groupId, domainId)
	if err != nil {
		return nil, err
	}

	if len(utils.PathSearch("ssl_infos", domainInfos, make([]interface{}, 0)).([]interface{})) == 0 {
		return nil, golangsdk.ErrDefault404{}
	}

	return domainInfos, nil
}

func resourceDomainCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		groupId    = d.Get("group_id").(string)
		domainId   = d.Get("domain_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	certificates, err := GetDomainAssociatedCertificatesByDomainId(client, instanceId, groupId, domainId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("unable to query the certificate associated with the domain (%s)", domainId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("certificate_ids", utils.PathSearch("ssl_infos[*].ssl_id", certificates, nil)),
		d.Set("verified_client_certificate_enabled", utils.PathSearch("verified_client_certificate_enabled", certificates, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// Currently, a certificate can only be associated with one domain name. After multiple certificates are supported in the future,
// the logic for updating certificate IDs will be supplemented.
func resourceDomainCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		groupId  = d.Get("group_id").(string)
		domainId = d.Get("domain_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	if d.HasChanges("verified_client_certificate_enabled") {
		err = associateCertificatesToDomain(
			client,
			d.Get("instance_id").(string),
			groupId,
			domainId,
			d.Get("certificate_ids").(*schema.Set).List(),
			d.Get("verified_client_certificate_enabled").(bool),
		)
		if err != nil {
			return diag.FromErr(fmt.Errorf("error updating client verified function domain (%s) under group (%s): %s", domainId, groupId, err))
		}
	}

	return resourceDomainCertificateRead(ctx, d, meta)
}

func disassociateCertificatesFromDomain(client *golangsdk.ServiceClient, instanceId, groupId, domainId string, certificateIds []interface{},
	verifiedClientCertificateEnabled bool) error {
	httpUrl := "v2/{project_id}/apigw/instances/{instance_id}/api-groups/{group_id}/domains/{domain_id}/certificates/detach"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)
	createPath = strings.ReplaceAll(createPath, "{domain_id}", domainId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes:  []int{204},
		JSONBody: buildAssociateCertificatesToDomainBodyParams(certificateIds, verifiedClientCertificateEnabled),
	}

	_, err := client.Request("POST", createPath, &opt)
	return err
}

func resourceDomainCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		domainId = d.Get("domain_id").(string)
	)
	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = disassociateCertificatesFromDomain(
		client,
		d.Get("instance_id").(string),
		d.Get("group_id").(string),
		domainId,
		d.Get("certificate_ids").(*schema.Set).List(),
		d.Get("verified_client_certificate_enabled").(bool),
	)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error disassociating certificates from domain (%s)", domainId))
	}
	return nil
}

func resourceDomainCertificateAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<group_id>/<domain_id>', but got '%s'",
			importedId)
	}

	mErr := multierror.Append(
		d.Set("instance_id", parts[0]),
		d.Set("group_id", parts[1]),
		d.Set("domain_id", parts[2]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
