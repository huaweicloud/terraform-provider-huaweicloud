package ccm

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

func DataSourceCcmPrivateCertificateExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: ccmPrivateCertificateExport,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"sm_standard": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_chain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_certificate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_private_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"enc_sm2_enveloped_key": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signed_and_enveloped_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func ccmPrivateCertificateExport(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	product := "ccm"

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	expCertificateHttpUrl := "v1/private-certificates/{certificate_id}/export"
	expCertificatePath := client.Endpoint + expCertificateHttpUrl
	expCertificatePath = strings.ReplaceAll(expCertificatePath, "{certificate_id}", d.Get("certificate_id").(string))

	expCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	expCertificateOpt.JSONBody = utils.RemoveNil(buildExpCertificateBodyParams(d))
	expCertificateResp, err := client.Request("POST", expCertificatePath, &expCertificateOpt)
	if err != nil {
		return diag.Errorf("error export CCM private certificate: %s", err)
	}
	expCertificateRespBody, err := utils.FlattenResponse(expCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(d.Get("certificate_id").(string))
	mErr := multierror.Append(nil,
		d.Set("private_key", utils.PathSearch("private_key", expCertificateRespBody, nil)),
		d.Set("certificate", utils.PathSearch("certificate", expCertificateRespBody, nil)),
		d.Set("certificate_chain", utils.PathSearch("certificate_chain", expCertificateRespBody, nil)),
		d.Set("enc_private_key", utils.PathSearch("enc_private_key", expCertificateRespBody, nil)),
		d.Set("enc_certificate", utils.PathSearch("enc_certificate", expCertificateRespBody, nil)),
		d.Set("enc_sm2_enveloped_key", utils.PathSearch("enc_sm2_enveloped_key", expCertificateRespBody, nil)),
		d.Set("signed_and_enveloped_data", utils.PathSearch("signed_and_enveloped_data", expCertificateRespBody, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting CCM private certificate export fields: %s", err)
	}
	return nil
}

func buildExpCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"is_compressed":  "false",
		"type":           d.Get("type"),
		"is_sm_standard": d.Get("sm_standard"),
		"password":       d.Get("password"),
	}
	return bodyParams
}
