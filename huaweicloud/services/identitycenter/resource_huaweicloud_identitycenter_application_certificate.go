package identitycenter

import (
	"context"
	"errors"
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

var identityCenterApplicationCertificateNonUpdateParams = []string{"instance_id", "application_instance_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates
// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}
// @API IdentityCenter GET /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates
// @API IdentityCenter DELETE /v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}
func ResourceIdentityCenterApplicationCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationCertificateCreate,
		UpdateContext: resourceIdentityCenterApplicationCertificateUpdate,
		ReadContext:   resourceIdentityCenterApplicationCertificateRead,
		DeleteContext: resourceIdentityCenterApplicationCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterApplicationCertificateImport,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationCertificateNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"application_instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"algorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"key_size": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"issue_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceIdentityCenterApplicationCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		createHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		updateHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		product       = "identitycenter"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	createPath = strings.ReplaceAll(createPath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center application certificate: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("application_instance_certificate.certificate_id", createRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the Identity Center application certificate ID from the API response")
	}
	d.SetId(certificateId)

	if status, ok := d.GetOk("status"); ok && status == "ACTIVE" {
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
		updatePath = strings.ReplaceAll(updatePath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))
		updatePath = strings.ReplaceAll(updatePath, "{certificate_id}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating Identity Center application certificate status: %s", err)
		}
	}

	return resourceIdentityCenterApplicationCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterApplicationCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates"
		product     = "identitycenter"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	listPath := client.Endpoint + listHttpUrl
	listPath = strings.ReplaceAll(listPath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	listPath = strings.ReplaceAll(listPath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	listResp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center service provider certificate.")
	}

	listRespBody, err := utils.FlattenResponse(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificate := utils.PathSearch(fmt.Sprintf("application_instance_certificates[?certificate_id =='%s']|[0]", d.Id()), listRespBody, nil)

	if certificate == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no application certificate found.")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("algorithm", utils.PathSearch("algorithm", certificate, nil)),
		d.Set("certificate", utils.PathSearch("certificate", certificate, nil)),
		d.Set("status", utils.PathSearch("status", certificate, nil)),
		d.Set("key_size", utils.PathSearch("key_size", certificate, nil)),
		d.Set("expiry_date", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("expiry_date", certificate, float64(0)).(float64))/1000, false)),
		d.Set("issue_date", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("issue_date", certificate, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceIdentityCenterApplicationCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		updateProduct = "identitycenter"
	)

	if d.HasChange("status") {
		if d.Get("status").(string) != "ACTIVE" {
			return diag.Errorf("status must be set to 'ACTIVE'")
		}

		client, err := cfg.NewServiceClient(updateProduct, region)
		if err != nil {
			return diag.Errorf("error creating Identity Center Client: %s", err)
		}

		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
		updatePath = strings.ReplaceAll(updatePath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))
		updatePath = strings.ReplaceAll(updatePath, "{certificate_id}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating Identity Center application certificate status: %s", err)
		}
	}

	return resourceIdentityCenterApplicationCertificateRead(ctx, d, meta)
}

func resourceIdentityCenterApplicationCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/certificates/{certificate_id}"
		deleteProduct = "identitycenter"
	)

	client, err := cfg.NewServiceClient(deleteProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center Client: %s", err)
	}

	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{instance_id}", fmt.Sprintf("%v", d.Get("instance_id")))
	deletePath = strings.ReplaceAll(deletePath, "{application_instance_id}", fmt.Sprintf("%v", d.Get("application_instance_id")))
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting Identity Center application certificate: %s", err)
	}

	return nil
}

func resourceIdentityCenterApplicationCertificateImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 3 {
		return nil, errors.New("invalid id format, must be <instance_id>/<application_instance_id>/<certificate_id>")
	}
	d.SetId(parts[2])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_instance_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
