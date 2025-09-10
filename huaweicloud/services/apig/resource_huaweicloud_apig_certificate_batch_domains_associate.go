package apig

import (
	"context"
	"fmt"
	"log"
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

var (
	certificateBatchDomainsAssociateNonUpdatableParams = []string{"instance_id", "certificate_id"}
	certificateBatchDomainsAssociateStrSliceParamKeys  = []string{"verify_disabled_domain_names", "verify_enabled_domain_names"}
)

// @API APIG POST /v2/{project_id}/apigw/certificates/{certificate_id}/domains/attach
// @API APIG POST /v2/{project_id}/apigw/certificates/{certificate_id}/domains/detach
// @API APIG GET /v2/{project_id}/apigw/certificates/{certificate_id}/attached-domains
func ResourceCertificateBatchDomainsAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateBatchDomainsAssociateCreate,
		ReadContext:   resourceCertificateBatchDomainsAssociateRead,
		UpdateContext: resourceCertificateBatchDomainsAssociateUpdate,
		DeleteContext: resourceCertificateBatchDomainsAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(certificateBatchDomainsAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceCertificateBatchDomainsAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the SSL certificate and domains are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the instance to which the certificate and domains belong.",
			},
			"certificate_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the SSL certificate.",
			},
			"verify_enabled_domain_names": {
				Type:             schema.TypeSet,
				Optional:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressStrSliceDiffs(),
				Description:      "The domain list to be enabled client certificate verification.",
			},
			"verify_disabled_domain_names": {
				Type:             schema.TypeSet,
				Optional:         true,
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
			"verify_enabled_domain_names_origin": {
				Type:             schema.TypeList,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for comparison with
the new value next time the change is made. The corresponding parameter name is 'verify_enabled_domain_names'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
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
		},
	}
}

func convertAssociatedDomains(instanceId string, domains []interface{}, enabled bool) []map[string]interface{} {
	if len(domains) == 0 {
		return make([]map[string]interface{}, 0)
	}

	result := make([]map[string]interface{}, len(domains))
	for _, domain := range domains {
		result = append(result, map[string]interface{}{
			"domain":                              domain,
			"instance_ids":                        []string{instanceId},
			"verified_client_certificate_enabled": enabled,
		})
	}
	return result
}

func buildCertificateBatchDomainsAssociateBodyParams(domains []map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"domains": domains,
	}
}

func associateDomainsToCertificate(client *golangsdk.ServiceClient, certificateId string, domains []map[string]interface{}) error {
	httpUrl := "v2/{project_id}/apigw/certificates/{certificate_id}/domains/attach"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{certificate_id}", certificateId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCertificateBatchDomainsAssociateBodyParams(domains)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error associating domains to certificate (%s): %s", certificateId, err)
	}

	return nil
}

func resourceCertificateBatchDomainsAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Lock the resource to prevent concurrent creation.
	config.MutexKV.Lock(certificateId)
	defer config.MutexKV.Unlock(certificateId)

	domains := append(
		convertAssociatedDomains(instanceId, d.Get("verify_enabled_domain_names").(*schema.Set).List(), true),
		convertAssociatedDomains(instanceId, d.Get("verify_disabled_domain_names").(*schema.Set).List(), false)...,
	)
	err = associateDomainsToCertificate(client, certificateId, domains)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%s/%s", instanceId, certificateId))

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, certificateBatchDomainsAssociateStrSliceParamKeys)
	if err != nil {
		// Don't fail the creation if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceCertificateBatchDomainsAssociateRead(ctx, d, meta)
}

func listCertificateAssociatedDomains(client *golangsdk.ServiceClient, certificateId string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/apigw/certificates/{certificate_id}/attached-domains"
		// The default limit is 20.
		limit  = 100
		offset = 0
		result = make([]interface{}, 0)
	)
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{certificate_id}", certificateId)
	listPath += fmt.Sprintf("?limit=%d", limit)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		attachedDomains := utils.PathSearch("bound_domains", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, attachedDomains...)
		if len(attachedDomains) < limit {
			break
		}

		offset += len(attachedDomains)
	}
	return result, nil
}

func GetCertificateAssociatedDomainsByDomain(client *golangsdk.ServiceClient, certificateId string, originEnabledDomainNames,
	originDisabledDomainNames []interface{}) (verifyEnabledDomains, verifyDisabledDomains []interface{}, err error) {
	attachedDomains, err := listCertificateAssociatedDomains(client, certificateId)
	if err != nil {
		return
	}

	if len(attachedDomains) < 1 {
		err = golangsdk.ErrDefault404{}
		return
	}

	// If origin exists, the console only deletes the domains created by the script, not the manually created domains.
	verifyEnabledDomains, verifyDisabledDomains = flattenCertificateAssociatedDomains(attachedDomains)
	if len(originEnabledDomainNames) > 0 || len(originDisabledDomainNames) > 0 {
		if len(utils.FildSliceIntersection(verifyEnabledDomains, originEnabledDomainNames)) == 0 &&
			len(utils.FildSliceIntersection(verifyDisabledDomains, originDisabledDomainNames)) == 0 {
			err = golangsdk.ErrDefault404{}
			return
		}
	}

	return
}

func resourceCertificateBatchDomainsAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	enabledDomains, disabledDomains, err := GetCertificateAssociatedDomainsByDomain(
		client,
		certificateId,
		d.Get("verify_enabled_domain_names_origin").([]interface{}),
		d.Get("verify_disabled_domain_names_origin").([]interface{}),
	)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving the associated domains with the SSL certificate (%s)", certificateId))
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("verify_enabled_domain_names", enabledDomains),
		d.Set("verify_disabled_domain_names", disabledDomains),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

// If the instance has the custom inbound port feature, the same domain name will be bound to one certificate at the same time,
// so must be deduplicated.
// The behavior of `verified_client_certificate_enabled` is consistent with the corresponding domain.
func uniqueCertificateAssociatedDomains(domains []interface{}, elem interface{}) []interface{} {
	if !utils.SliceContains(domains, elem) {
		return append(domains, elem)
	}
	return domains
}

func flattenCertificateAssociatedDomains(associatedDomains []interface{}) (verifyEnabledDomains, verifyDisabledDomains []interface{}) {
	for _, attachedDomain := range associatedDomains {
		enabled := utils.PathSearch("verified_client_certificate_enabled", attachedDomain, false).(bool)
		domainName := utils.PathSearch("url_domain", attachedDomain, "").(string)
		if domainName == "" {
			continue
		}

		if enabled {
			verifyEnabledDomains = uniqueCertificateAssociatedDomains(verifyEnabledDomains, domainName)
		} else {
			verifyDisabledDomains = uniqueCertificateAssociatedDomains(verifyDisabledDomains, domainName)
		}
	}

	return verifyEnabledDomains, verifyDisabledDomains
}

// Before disassociating domains, get the all associated domains under the certificate.
// Prevent errors when unbinding an unbound domain name when operating multiple resources simultaneously.
func findBoundDomainsUnderCertificate(client *golangsdk.ServiceClient, certificateId string,
	domains []map[string]interface{}) ([]interface{}, error) {
	attachedDomains, err := listCertificateAssociatedDomains(client, certificateId)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			log.Printf("[DEBUG] All domains have been disassociated from the certificate (%s)", certificateId)
			return nil, nil
		}
		return nil, fmt.Errorf("error querying bound domains for certificate (%s)", certificateId)
	}

	deleteDomainNames := utils.FildSliceIntersection(
		utils.PathSearch("[].domain", domains, make([]interface{}, 0)).([]interface{}),
		utils.PathSearch("[].url_domain", attachedDomains, make([]interface{}, 0)).([]interface{}),
	)
	return deleteDomainNames, nil
}

func unbindDomainsFromCertificate(client *golangsdk.ServiceClient, certificateId string, domains []map[string]interface{}) error {
	httpUrl := "v2/{project_id}/apigw/certificates/{certificate_id}/domains/detach"
	detachPath := client.Endpoint + httpUrl
	detachPath = strings.ReplaceAll(detachPath, "{project_id}", client.ProjectID)
	detachPath = strings.ReplaceAll(detachPath, "{certificate_id}", certificateId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCertificateBatchDomainsAssociateBodyParams(domains)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("POST", detachPath, &opt)
	return err
}

func getUpdateCertificateBatchDomainsAssociateDomainNames(d *schema.ResourceData, key string) (newDomainNames, rmDomainNames []interface{}) {
	var (
		consoleDomainNames, scriptDomainNames = d.GetChange(key)
		consoleDomainNamesList                = consoleDomainNames.(*schema.Set).List()
		scriptDomainNamesList                 = scriptDomainNames.(*schema.Set).List()
		originDomainNamesList                 = d.Get(fmt.Sprintf("%s_origin", key)).([]interface{})
	)

	newDomainNames = utils.FindSliceElementsNotInAnother(scriptDomainNamesList, consoleDomainNamesList)
	rmDomainNames = utils.FindSliceElementsNotInAnother(originDomainNamesList, scriptDomainNamesList)
	return newDomainNames, rmDomainNames
}

func resourceCertificateBatchDomainsAssociateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		certificateId = d.Get("certificate_id").(string)
		addRaws       = make([]map[string]interface{}, 0)
		removeRaws    = make([]map[string]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Lock the resource to prevent concurrent updates.
	config.MutexKV.Lock(certificateId)
	defer config.MutexKV.Unlock(certificateId)

	if d.HasChange("verify_enabled_domain_names") {
		newDomainNames, rmDomainNames := getUpdateCertificateBatchDomainsAssociateDomainNames(d, "verify_enabled_domain_names")
		addRaws = convertAssociatedDomains(instanceId, newDomainNames, true)
		removeRaws = convertAssociatedDomains(instanceId, rmDomainNames, true)
	}

	if d.HasChange("verify_disabled_domain_names") {
		newDomainNames, rmDomainNames := getUpdateCertificateBatchDomainsAssociateDomainNames(d, "verify_disabled_domain_names")
		addRaws = append(addRaws, convertAssociatedDomains(instanceId, newDomainNames, false)...)
		removeRaws = append(removeRaws, convertAssociatedDomains(instanceId, rmDomainNames, false)...)
	}

	// Must disassociate domains first, prevent the bound data from containing domains to be unbound.
	if len(removeRaws) > 0 {
		if err = disassociateDomainsFromCertificate(client, certificateId, removeRaws); err != nil {
			return diag.FromErr(err)
		}
	}

	if len(addRaws) > 0 {
		if err = associateDomainsToCertificate(client, certificateId, addRaws); err != nil {
			return diag.FromErr(err)
		}
	}

	// If the request is successful, obtain the values of all slice parameters first and save them to the corresponding
	// '_origin' attributes for subsequent determination and construction of the request body during next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshSliceParamOriginValues(d, certificateBatchDomainsAssociateStrSliceParamKeys)
	if err != nil {
		// Don't fail the update if origin refresh fails
		log.Printf("[WARN] Unable to refresh the origin values: %s", err)
	}

	return resourceCertificateBatchDomainsAssociateRead(ctx, d, meta)
}

func disassociateDomainsFromCertificate(client *golangsdk.ServiceClient, certificateId string, rmDomains []map[string]interface{}) error {
	deleteDomainNames, err := findBoundDomainsUnderCertificate(client, certificateId, rmDomains)
	if err != nil {
		return err
	}

	if len(deleteDomainNames) == 0 {
		log.Printf("[DEBUG] unable to find any domains to bind under certificate (%s)", certificateId)
		return nil
	}

	return unbindDomainsFromCertificate(client, certificateId, rmDomains)
}

// getCertificateBoundDomains retrieves certificate associated domains from configuration or origin.
func getCertificateBoundDomains(d *schema.ResourceData, key string) []interface{} {
	// Fallback to origin (last known configuration).
	if origin, ok := d.Get(fmt.Sprintf("%s_origin", key)).([]interface{}); ok && len(origin) > 0 {
		log.Printf("[DEBUG] Found %d domain(s) from the origin attribute: %v", len(origin), origin)
		return origin
	}

	log.Printf("[DEBUG] Unable to find the domain(s) from the origin attribute, so try to get from current state")
	// After resource imported, the origin attribute is not set, so try to get from current state.
	current := d.Get(key).(*schema.Set).List()
	log.Printf("[DEBUG] Found %d domain(s) from the current state: %v", len(current), current)

	return current
}

func resourceCertificateBatchDomainsAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		instanceId    = d.Get("instance_id").(string)
		certificateId = d.Get("certificate_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	// Lock the resource to prevent concurrent deletion.
	config.MutexKV.Lock(certificateId)
	defer config.MutexKV.Unlock(certificateId)

	deleteDomains := append(
		convertAssociatedDomains(instanceId, getCertificateBoundDomains(d, "verify_enabled_domain_names"), true),
		convertAssociatedDomains(instanceId, getCertificateBoundDomains(d, "verify_disabled_domain_names"), false)...,
	)
	// { "error_code": "APIG.2000", "error_msg": "no valid domain name"}: All domains have been unbound from the certificate.
	if err = unbindDomainsFromCertificate(client, certificateId, deleteDomains); err != nil {
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "APIG.2000"),
			fmt.Sprintf("error disassociating domains from certificate (%s)", certificateId),
		)
	}

	return nil
}

func resourceCertificateBatchDomainsAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<certificate_id>', "+
			"but got '%s'", importedId)
	}

	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("certificate_id", parts[1]),
	)

	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
