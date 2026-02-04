package ccm

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var mapParamKeys = []string{
	"tags",
}

// @API CCM POST /v3/scm/certificates/buy
// @API CCM GET /v3/scm/certificates/{certificate_id}
// @API CCM DELETE /v3/scm/certificates/{cert_id}/unsubscribe
// @API CCM POST /v3/scm/{resource_id}/tags/action
// @API CCM GET /v3/scm/{resource_id}/tags
// @API EPS POST /v1.0/enterprise-projects/{enterprise_project_id}/resources-migrate
func ResourceCCMCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCCMCertificateCreate,
		UpdateContext: resourceCCMCertificateUpdate,
		ReadContext:   resourceCCMCertificateRead,
		DeleteContext: resourceCCMCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cert_brand": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the certificate authority.`,
			},
			"cert_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the certificate type.`,
			},
			"domain_type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the type of domain name.`,
			},
			"effective_time": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the validity period (year).`,
			},
			"domain_numbers": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the quantity of domain name.`,
			},
			"primary_domain_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the type of primary domain name in multiple domains.`,
			},
			"single_domain_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of additional single domain names.`,
			},
			"wildcard_domain_number": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the number of additional wildcard domain names.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the enterprise project ID.`,
			},
			"tags": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressMapDiffs(),
			},
			// Internal attributes.
			"tags_origin": {
				Type:             schema.TypeMap,
				Optional:         true,
				Computed:         true,
				Elem:             &schema.Schema{Type: schema.TypeString},
				DiffSuppressFunc: utils.SuppressDiffAll,
				Description: utils.SchemaDesc(
					`The script configuration value of this change is also the original value used for
					comparison with the new value next time the change is made. The corresponding parameter name is
					'tags'.`,
					utils.SchemaDescInput{
						Internal: true,
					},
				),
			},
			"validity_period": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The validity period (month).`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate status.`,
			},
			"order_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The order ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate name.`,
			},
			"push_support": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Whether the certificate supports push.`,
			},
			"revoke_reason": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The reason for certificate revocation.`,
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The signature algorithm.`,
			},
			"issue_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate issuance time.`,
			},
			"not_before": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate validity time.`,
			},
			"not_after": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The certificate expiration time.`,
			},
			"validation_method": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The authentication method of domain name.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name bound to the certificate.`,
			},
			"sans": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The information of additional domain name for the bound certificate.`,
			},
			"fingerprint": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SHA-1 fingerprint of the certificate.`,
			},
			"authentification": {
				Type:        schema.TypeList,
				Elem:        authentificationSchema(),
				Computed:    true,
				Description: `The ownership certification information of domain name.`,
			},
		},
	}
}

func authentificationSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"record_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The name of the domain name check value.`,
			},
			"record_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the domain name check value.`,
			},
			"record_value": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name check value.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain name corresponding to the check value.`,
			},
		},
	}
	return &sc
}

// The certificate status can only be used normally after it changes to `PAID`.
func waitingForCCMCertificatePaid(ctx context.Context, client *golangsdk.ServiceClient, certID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getRespBody, err := ReadCCMCertificate(client, certID)
			if err != nil {
				return nil, "ERROR", err
			}

			status := utils.PathSearch("status", getRespBody, "").(string)
			if status == "" {
				return nil, "ERROR", fmt.Errorf("status is not found in API response")
			}

			if status == "PAID" {
				return "success", "COMPLETED", nil
			}

			return getRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceCCMCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/scm/certificates/buy"
		product = "scm"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateCCMCertificateBodyParams(d, cfg)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CCM certificate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certId := utils.PathSearch("cert|[0].cert_id", createRespBody, "").(string)
	if certId == "" {
		return diag.Errorf("unable to find the CCM certificate ID from the API response")
	}
	d.SetId(certId)

	if err := waitingForCCMCertificatePaid(ctx, client, d.Id(), d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for CCM certificate (%s) creation to PAID: %s", d.Id(), err)
	}

	tagsRaw := d.Get("tags").(map[string]interface{})
	if err := createOrUpdateCCMCertificateTags(client, d, "create", tagsRaw); err != nil {
		return diag.Errorf("error creating CCM certificate tags in create operation: %s", err)
	}

	// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
	// corresponding '_origin' attributes for subsequent determination and construction of the request body during
	// next updates.
	// And whether corresponding parameters are changed, the origin values must be refreshed.
	err = utils.RefreshObjectParamOriginValues(d, mapParamKeys)
	if err != nil {
		return diag.Errorf("unable to refresh the origin values: %s", err)
	}

	return resourceCCMCertificateRead(ctx, d, meta)
}

// createOrUpdateCCMCertificateTags The valid value of action is "create" or "delete".
func createOrUpdateCCMCertificateTags(client *golangsdk.ServiceClient, d *schema.ResourceData, action string,
	tagRaw map[string]interface{}) error {
	if len(tagRaw) == 0 {
		return nil
	}

	requestPath := client.Endpoint + "v3/scm/{resource_id}/tags/action"
	requestPath = strings.ReplaceAll(requestPath, "{resource_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"action": action,
			"tags":   utils.ExpandResourceTags(tagRaw),
		},
	}

	_, err := client.Request("POST", requestPath, &requestOpt)
	return err
}

func buildCreateCCMCertificateBodyParams(d *schema.ResourceData, cfg *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cert_brand":               d.Get("cert_brand"),
		"cert_type":                d.Get("cert_type"),
		"domain_type":              d.Get("domain_type"),
		"effective_time":           d.Get("effective_time"),
		"domain_numbers":           d.Get("domain_numbers"),
		"primary_domain_type":      utils.ValueIgnoreEmpty(d.Get("primary_domain_type")),
		"single_domain_number":     utils.ValueIgnoreEmpty(d.Get("single_domain_number")),
		"wildcard_domain_number":   utils.ValueIgnoreEmpty(d.Get("wildcard_domain_number")),
		"enterprise_project_id":    utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"order_number":             1,
		"agree_privacy_protection": true,
		"is_auto_pay":              true,
	}
	return bodyParams
}

func flattenPrimaryDomainTypeAttribute(respBody interface{}) string {
	multiDomainType := utils.PathSearch("multi_domain_type", respBody, "").(string)
	switch multiDomainType {
	case "primary_single":
		return "SINGLE_DOMAIN"
	case "primary_wildcard":
		return "WILDCARD_DOMAIN"
	}

	return ""
}

func flattenAuthentification(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("authentification", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"record_name":  utils.PathSearch("record_name", v, nil),
			"record_type":  utils.PathSearch("record_type", v, nil),
			"record_value": utils.PathSearch("record_value", v, nil),
			"domain":       utils.PathSearch("domain", v, nil),
		})
	}
	return rst
}

// ReadCCMCertificate Test cases use this method, so the first letter is capitalized
func ReadCCMCertificate(client *golangsdk.ServiceClient, certID string) (interface{}, error) {
	getPath := client.Endpoint + "v3/scm/certificates/{certificate_id}"
	getPath = strings.ReplaceAll(getPath, "{certificate_id}", certID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(getResp)
}

func setCCMCertificateTagsToState(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v3/scm/{resource_id}/tags"
	requestPath = strings.ReplaceAll(requestPath, "{resource_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		log.Printf("[WARN] Error fetching tags of %s: %s", d.Id(), err)
		return nil
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		log.Printf("[WARN] Error flattening response tags of %s: %s", d.Id(), err)
		return nil
	}

	tagsList := utils.PathSearch("tags", respBody, make([]interface{}, 0)).([]interface{})
	if len(tagsList) == 0 {
		return nil
	}

	return d.Set("tags", utils.FlattenTagsToMap(tagsList))
}

func resourceCCMCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		mErr    *multierror.Error
		product = "scm"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	getRespBody, err := ReadCCMCertificate(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CCM certificate")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("cert_brand", utils.PathSearch("brand", getRespBody, nil)),
		d.Set("cert_type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("domain_type", utils.PathSearch("domain_type", getRespBody, nil)),
		d.Set("validity_period", utils.PathSearch("validity_period", getRespBody, nil)),
		d.Set("domain_numbers", utils.PathSearch("domain_count", getRespBody, nil)),
		d.Set("primary_domain_type", flattenPrimaryDomainTypeAttribute(getRespBody)),
		d.Set("wildcard_domain_number", utils.PathSearch("wildcard_count", getRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("order_id", utils.PathSearch("order_id", getRespBody, nil)),
		d.Set("name", utils.PathSearch("name", getRespBody, nil)),
		d.Set("push_support", utils.PathSearch("push_support", getRespBody, nil)),
		d.Set("revoke_reason", utils.PathSearch("revoke_reason", getRespBody, nil)),
		d.Set("signature_algorithm", utils.PathSearch("signature_algorithm", getRespBody, nil)),
		d.Set("issue_time", utils.PathSearch("issue_time", getRespBody, nil)),
		d.Set("not_before", utils.PathSearch("not_before", getRespBody, nil)),
		d.Set("not_after", utils.PathSearch("not_after", getRespBody, nil)),
		d.Set("validation_method", utils.PathSearch("validation_method", getRespBody, nil)),
		d.Set("domain", utils.PathSearch("domain", getRespBody, nil)),
		d.Set("sans", utils.PathSearch("sans", getRespBody, nil)),
		d.Set("fingerprint", utils.PathSearch("fingerprint", getRespBody, nil)),
		d.Set("authentification", flattenAuthentification(getRespBody)),
		setCCMCertificateTagsToState(client, d),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCCMCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	if d.HasChange("enterprise_project_id") {
		migrateOpts := config.MigrateResourceOpts{
			ResourceId:   d.Id(),
			ResourceType: "scm",
			RegionId:     region,
			ProjectId:    cfg.GetProjectID(region),
		}
		if err := cfg.MigrateEnterpriseProject(ctx, d, migrateOpts); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		client, err := cfg.NewServiceClient("scm", region)
		if err != nil {
			return diag.Errorf("error creating CCM client: %s", err)
		}

		oRaw, nRaw := d.GetChange("tags")
		// remove old tags
		if err := createOrUpdateCCMCertificateTags(client, d, "delete", oRaw.(map[string]interface{})); err != nil {
			return diag.Errorf("error deleting CCM certificate tags in update operation: %s", err)
		}

		// set new tags
		if err := createOrUpdateCCMCertificateTags(client, d, "create", nRaw.(map[string]interface{})); err != nil {
			return diag.Errorf("error creating CCM certificate tags in update operation: %s", err)
		}

		// If the request is successful, obtain the values of all JSON|object parameters first and save them to the
		// corresponding '_origin' attributes for subsequent determination and construction of the request body during
		// next updates.
		// And whether corresponding parameters are changed, the origin values must be refreshed.
		err = utils.RefreshObjectParamOriginValues(d, mapParamKeys)
		if err != nil {
			return diag.Errorf("unable to refresh the origin values: %s", err)
		}
	}
	return resourceCCMCertificateRead(ctx, d, meta)
}

func waitingForCCMCertificateDeleted(ctx context.Context, client *golangsdk.ServiceClient, certID string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			getRespBody, err := ReadCCMCertificate(client, certID)
			if err != nil {
				var errDefault404 golangsdk.ErrDefault404
				if errors.As(err, &errDefault404) {
					return "success", "COMPLETED", nil
				}
				return nil, "ERROR", err
			}

			return getRespBody, "PENDING", nil
		},
		Timeout:      timeout,
		Delay:        20 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

// buildDeleteRetryFunc When the time interval between creation and deletion is too short, deletion will fail,
// so deletion needs to be retried.
// An example of the response body for deletion retry is as follows: {"error_code": "SCM.0016","error_msg": "订单异常"}
func buildDeleteRetryFunc(client *golangsdk.ServiceClient, deletePath string, deleteOpt *golangsdk.RequestOpts) common.RetryFunc {
	retryFunc := func() (interface{}, bool, error) {
		deleteResp, err := client.Request("DELETE", deletePath, deleteOpt)
		if err != nil {
			var errCode golangsdk.ErrDefault500
			if errors.As(err, &errCode) {
				var apiError interface{}
				if jsonErr := json.Unmarshal(errCode.Body, &apiError); jsonErr != nil {
					return nil, false, err
				}

				errorCode := utils.PathSearch("error_code", apiError, "").(string)
				if errorCode == "SCM.0016" {
					return nil, true, err
				}
			}
			return nil, false, err
		}

		deleteRespBody, err := utils.FlattenResponse(deleteResp)
		return deleteRespBody, false, err
	}
	return retryFunc
}

func resourceCCMCertificateDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/scm/certificates/{cert_id}/unsubscribe"
		product = "scm"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{cert_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	deleteRespBody, err := common.RetryContextWithWaitForState(&common.RetryContextWithWaitForStateParam{
		Ctx:          ctx,
		RetryFunc:    buildDeleteRetryFunc(client, deletePath, &deleteOpt),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		DelayTimeout: 10 * time.Second,
		PollInterval: 10 * time.Second,
	})

	if err != nil {
		return diag.Errorf("error deleting CCM certificate: %s", err)
	}

	unsubscribeResult := utils.PathSearch("unsubscribe_results", deleteRespBody, "").(string)
	if unsubscribeResult != "SUCCESS" {
		return diag.Errorf("error deleting CCM certificate: Unsubscribe result is not SUCCESS in deletion API response")
	}

	if err := waitingForCCMCertificateDeleted(ctx, client, d.Id(), d.Timeout(schema.TimeoutDelete)); err != nil {
		return diag.Errorf("error waiting for CCM certificate (%s) deleted: %s", d.Id(), err)
	}

	return nil
}
