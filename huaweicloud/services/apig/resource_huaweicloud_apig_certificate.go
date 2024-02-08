package apig

import (
	"context"
	"fmt"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/certificates"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API APIG POST /v2/{project_id}/apigw/certificates
// @API APIG DELETE /v2/{project_id}/apigw/certificates/{certificate_id}
// @API APIG GET /v2/{project_id}/apigw/certificates/{certificate_id}
// @API APIG PUT /v2/{project_id}/apigw/certificates/{certificate_id}
func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the certificate is located.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The certificate name.",
			},
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The certificate content.",
			},
			"private_key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "The private key of the certificate.",
			},
			"type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The certificate type.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The dedicated instance ID to which the certificate belongs.",
			},
			"trusted_root_ca": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: "The trusted root CA certificate.",
			},
			"effected_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The effective time of the certificate.",
			},
			"expires_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time of the certificate.",
			},
			"signature_algorithm": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "What signature algorithm the certificate uses.",
			},
			"sans": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The SAN (Subject Alternative Names) of the certificate.",
			},
		},
	}
}

func buildCertificateModifyOpts(d *schema.ResourceData) certificates.CertOpts {
	return certificates.CertOpts{
		Name:          d.Get("name").(string),
		Content:       d.Get("content").(string),
		PrivateKey:    d.Get("private_key").(string),
		Type:          d.Get("type").(string),
		InstanceId:    d.Get("instance_id").(string),
		TrustedRootCA: d.Get("trusted_root_ca").(string),
	}
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	resp, err := certificates.Create(client, buildCertificateModifyOpts(d))
	if err != nil {
		return diag.Errorf("error creating APIG SSL certificate: %s", err)
	}
	d.SetId(resp.ID)

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ApigV2Client(region)
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	certificateId := d.Id()
	resp, err := certificates.Get(client, certificateId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "APIG SSL certificate")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", resp.Name),
		d.Set("type", resp.Type),
		d.Set("instance_id", resp.InstanceId),
		d.Set("effected_at", resp.NotBefore),
		d.Set("expires_at", resp.NotAfter),
		d.Set("signature_algorithm", resp.SignatureAlgorithm),
		d.Set("sans", resp.SANs),
	)
	if mErr.ErrorOrNil() != nil {
		return diag.Errorf("error saving APIG SSL certificate (%s) fields: %s", certificateId, mErr)
	}
	return nil
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	_, err = certificates.Update(client, d.Id(), buildCertificateModifyOpts(d))
	if err != nil {
		return diag.Errorf("error updating APIG SSL certificate: %s", err)
	}

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	certificateId := d.Id()
	err = certificates.Delete(client, certificateId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting APIG SSL certificate (%s): %s",
			certificateId, err))
	}

	return nil
}
