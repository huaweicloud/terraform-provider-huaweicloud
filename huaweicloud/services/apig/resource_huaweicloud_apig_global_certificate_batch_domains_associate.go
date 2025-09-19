package apig

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	globalCertificateBatchDomainsAssociateNonUpdatableParams = []string{"certificate_id"}
	globalCertificateBatchDomainsAssociateStrSliceParamKeys  = []string{"verify_disabled_domain_names"}
)

// @API APIG POST /v2/{project_id}/apigw/certificates/{certificate_id}/domains/attach
// @API APIG POST /v2/{project_id}/apigw/certificates/{certificate_id}/domains/detach
// @API APIG GET /v2/{project_id}/apigw/certificates/{certificate_id}/attached-domains
func ResourceGlobalCertificateBatchDomainsAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGlobalCertificateBatchDomainsAssociateCreate,
		ReadContext:   resourceGlobalCertificateBatchDomainsAssociateRead,
		UpdateContext: resourceGlobalCertificateBatchDomainsAssociateUpdate,
		DeleteContext: resourceGlobalCertificateBatchDomainsAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(globalCertificateBatchDomainsAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the global SSL certificate and domains are located.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the global SSL certificate.",
			},
			"verify_disabled_domain_names": {
				Type:             schema.TypeSet,
				Required:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
				Description:      "The domain list to be disabled client certificate verification.",
			},
			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"verify_disabled_domain_names_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'verify_disabled_domain_names'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			// Attributes.
			"domains": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the associated domain.`,
						},
						"url_domain": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The associated domain name.`,
						},
						"instance_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the dedicated instance to which the domain belongs.`,
						},
						"status": {
							Type:        schema.TypeInt,
							Computed:    true,
							Description: `The CNAME resolution status of the domain name.`,
						},
						"min_ssl_version": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "The minimum SSL protocol version of the domain.",
						},
						"api_group_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the API group to which the domain belongs.`,
						},
						"api_group_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the API group to which the domain belongs.`,
						},
					},
				},
				Description: "The domain list associated with the global SSL certificate.",
			},
		},
	}
}

func buildGlobalCertificateAssociatedDomains(domains []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, len(domains))
	for _, domain := range domains {
		result = append(result, map[string]interface{}{
			"domain":                              domain,
			"verified_client_certificate_enabled": false,
		})
	}
	return result
}

func resourceGlobalCertificateBatchDomainsAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = associateDomainsToCertificate(
		client,
		certificateId,
		buildGlobalCertificateAssociatedDomains(d.Get("verify_disabled_domain_names").(*schema.Set).List()),
	)
	if err != nil {
		return diag.Errorf("error associating domains to global SSL certificate (%s): %s", certificateId, err)
	}

	d.SetId(certificateId)

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, globalCertificateBatchDomainsAssociateStrSliceParamKeys)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceGlobalCertificateBatchDomainsAssociateRead(ctx, d, meta)
}

func GetGlobalCertificateAssociatedDomains(client *golangsdk.ServiceClient, cerId string, originDomains []interface{}) (
	verifyDisabledDomains, verifyDisabledDomainNames []interface{}, err error) {
	attachedDomains, err := listCertificateAssociatedDomains(client, cerId)
	if err != nil {
		return nil, nil, err
	}

	if len(attachedDomains) < 1 {
		return nil, nil, golangsdk.ErrDefault404{}
	}

	verifyDisabledDomains = utils.PathSearch("[?verified_client_certificate_enabled == `false`]",
		attachedDomains, make([]interface{}, 0)).([]interface{})
	verifyDisabledDomainNames = utils.PathSearch("[*].url_domain", verifyDisabledDomains, make([]interface{}, 0)).([]interface{})
	// If origin exists, the console only deletes the domains created by the script, not the manually created domains.
	if len(originDomains) > 0 && len(utils.FildSliceIntersection(verifyDisabledDomainNames, originDomains)) == 0 {
		return nil, nil, golangsdk.ErrDefault404{}
	}

	return verifyDisabledDomains, verifyDisabledDomainNames, nil
}

func resourceGlobalCertificateBatchDomainsAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Id()
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	domains, verifyDisabledDomains, err := GetGlobalCertificateAssociatedDomains(
		client,
		certificateId,
		d.Get("verify_disabled_domain_names_origin").([]interface{}),
	)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving the associated domains with the global SSL certificate (%s)", certificateId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("certificate_id", certificateId),
		d.Set("verify_disabled_domain_names", verifyDisabledDomains),
		// Attributes.
		d.Set("domains", flattenGlobalCertificateAssociatedDomains(domains)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGlobalCertificateAssociatedDomains(domains []interface{}) []interface{} {
	if len(domains) < 1 {
		return nil
	}

	result := make([]interface{}, 0, len(domains))
	for _, domain := range domains {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", domain, nil),
			"url_domain":      utils.PathSearch("url_domain", domain, nil),
			"instance_id":     utils.PathSearch("instance_id", domain, nil),
			"status":          utils.PathSearch("status", domain, nil),
			"min_ssl_version": utils.PathSearch("min_ssl_version", domain, nil),
			"api_group_id":    utils.PathSearch("api_group_id", domain, nil),
			"api_group_name":  utils.PathSearch("api_group_name", domain, nil),
		})
	}
	return result
}

func resourceGlobalCertificateBatchDomainsAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                 = meta.(*config.Config)
		region              = cfg.GetRegion(d)
		certificateId       = d.Get("certificate_id").(string)
		addRaws, removeRaws = getUpdateCertificateBatchDomainsAssociateDomainNames(d, "verify_disabled_domain_names")
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Must disassociate domains first, prevent the bound data from containing domains to be unbound.
	if len(removeRaws) > 0 {
		if err = disassociateDomainsFromCertificate(client, certificateId, buildGlobalCertificateAssociatedDomains(removeRaws)); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addRaws) > 0 {
		if err = associateDomainsToCertificate(client, certificateId, buildGlobalCertificateAssociatedDomains(addRaws)); err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, globalCertificateBatchDomainsAssociateStrSliceParamKeys)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceGlobalCertificateBatchDomainsAssociateRead(ctx, d, meta)
}

func resourceGlobalCertificateBatchDomainsAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	deleteDomains := buildGlobalCertificateAssociatedDomains(getCertificateBoundDomains(d, "verify_disabled_domain_names"))
	// { "error_code": "APIG.2000", "error_msg": "no valid domain name"}: All domains have been unbound from the certificate.
	if err = unbindDomainsFromCertificate(client, certificateId, deleteDomains); err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "APIG.2000"),
			fmt.Sprintf("error disassociating domains from global certificate (%s)", certificateId),
		)
	}

	return nil
}
