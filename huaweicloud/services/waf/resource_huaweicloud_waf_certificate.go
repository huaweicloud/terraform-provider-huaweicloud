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
func ResourceWafCertificateV1() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceWafCertificateV1Create,
		ReadContext:   resourceWafCertificateV1Read,
		UpdateContext: resourceWafCertificateV1Update,
		DeleteContext: resourceWafCertificateV1Delete,
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
			"expiration": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceWafCertificateV1Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	return resourceWafCertificateV1Read(ctx, d, meta)
}

func resourceWafCertificateV1Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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

	expires := time.Unix(int64(n.ExpireTime/1000), 0).UTC().Format("2006-01-02 15:04:05 MST")
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", n.Name),
		d.Set("expiration", expires),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

// Field `name` is required for updating operation.
// Field `certificate` and `private_key` must be specified together.
// Ignore the updating operation when field `certificate` does not change.
func resourceWafCertificateV1Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
	return resourceWafCertificateV1Read(ctx, d, meta)
}

func resourceWafCertificateV1Delete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
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
