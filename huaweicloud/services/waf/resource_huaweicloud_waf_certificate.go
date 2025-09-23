package waf

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/waf/v1/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF GET /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF PUT /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF DELETE /v1/{project_id}/waf/certificate/{certificate_id}
// @API WAF POST /v1/{project_id}/waf/certificate
func ResourceWafCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafCertificateCreate,
		ReadContext:   resourceWafCertificateRead,
		UpdateContext: resourceWafCertificateUpdate,
		DeleteContext: resourceWafCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceWAFImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: utils.SuppressTrimSpace,
				Sensitive:        true,
			},
			"private_key": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expired_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// Deprecated; Reasons for abandonment are as follows:
			// `expiration`: Uniformly use dates in RFC3339 format.
			"expiration": {
				Type:        schema.TypeString,
				Computed:    true,
				Deprecated:  "Use 'expired_at' instead. ",
				Description: `schema: Deprecated; The certificate expiration time.`,
			},
		},
	}
}

func resourceWafCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createOpts := certificates.CreateOpts{
		Name:                d.Get("name").(string),
		Content:             strings.TrimSpace(d.Get("certificate").(string)),
		Key:                 strings.TrimSpace(d.Get("private_key").(string)),
		EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
	}

	certificate, err := certificates.Create(client, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating WAF certificate: %s", err)
	}
	d.SetId(certificate.Id)

	return resourceWafCertificateRead(ctx, d, meta)
}

func resourceWafCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.WafV1Client(region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	n, err := certificates.GetWithEpsID(client, d.Id(), epsID).Extract()
	if err != nil {
		// If the certificate does not exist, the response HTTP status code of the details API is 404.
		return common.CheckDeletedDiag(d, err, "error retrieving WAF certificate")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(n.TimeStamp/1000), true)),
		d.Set("expired_at", utils.FormatTimeStampRFC3339(int64(n.ExpireTime/1000), true)),
		// Keep historical code logic
		d.Set("expiration", utils.FormatTimeStampRFC3339(int64(n.ExpireTime/1000), true, "2006-01-02 15:04:05 MST")),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// Field `name` is required for updating operation.
// Field `certificate` and `private_key` must be specified together.
// Ignore the updating operation when field `certificate` does not change.
func resourceWafCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	if d.HasChanges("name", "certificate") {
		updateOpts := certificates.UpdateOpts{
			Name:                d.Get("name").(string),
			EnterpriseProjectId: cfg.GetEnterpriseProjectID(d),
		}
		if d.HasChange("certificate") {
			updateOpts.Content = strings.TrimSpace(d.Get("certificate").(string))
			updateOpts.Key = strings.TrimSpace(d.Get("private_key").(string))
		}
		_, err = certificates.Update(client, d.Id(), updateOpts).Extract()
		if err != nil {
			return diag.Errorf("error updating WAF certificate: %s", err)
		}
	}
	return resourceWafCertificateRead(ctx, d, meta)
}

func resourceWafCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.WafV1Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	epsID := cfg.GetEnterpriseProjectID(d)
	err = certificates.DeleteWithEpsID(client, d.Id(), epsID).ExtractErr()
	if err != nil {
		// If the certificate does not exist, the response HTTP status code of the deletion API is 404.
		return common.CheckDeletedDiag(d, err, "error deleting WAF certificate")
	}
	return nil
}
