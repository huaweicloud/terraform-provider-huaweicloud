package ccm

import (
	"archive/zip"
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCM GET /v1/private-certificates/{certificate_id}/export
func DataSourcePrivateCertificateExport() *schema.Resource {
	return &schema.Resource{
		ReadContext: privateCertificateExport,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
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
			"keystore_pass": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_pfx": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"server_jks": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

// whetherCompressed Return true to obtain the compressed file of the certificate, and return false to obtain the
// certificate content in json format.
// Obtain the compressed file of the certificate because the certificate structure of `TOMCAT` and `IIS` types is not defined.
func whetherCompressed(serverType string) bool {
	if serverType == "IIS" || serverType == "TOMCAT" {
		return true
	}
	return false
}

func buildExpCertificateBodyParams(d *schema.ResourceData, serverType string, isCompressed bool) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"is_compressed":  isCompressed,
		"type":           serverType,
		"is_sm_standard": d.Get("sm_standard"),
		"password":       d.Get("password"),
	}
	return bodyParams
}

func saveCertificateFromZipFile(f *zip.File, d *schema.ResourceData) error {
	fileName := f.Name
	inFile, err := f.Open()
	if err != nil {
		return fmt.Errorf("error openning certificate file: %s", err)
	}
	defer inFile.Close()
	content, err := io.ReadAll(inFile)
	if err != nil {
		return fmt.Errorf("error reading unzip certificate file: %s", err)
	}
	var mErr *multierror.Error
	if fileName == "keystorePass.txt" {
		mErr = multierror.Append(mErr, d.Set("keystore_pass", string(content)))
	}
	if fileName == "server.pfx" {
		mErr = multierror.Append(mErr, d.Set("server_pfx", utils.Base64EncodeString(string(content))))
	}
	if fileName == "server.jks" {
		mErr = multierror.Append(mErr, d.Set("server_jks", utils.Base64EncodeString(string(content))))
	}
	return mErr.ErrorOrNil()
}

func handleCompressedCertificateResponse(d *schema.ResourceData, expCertificateResp *http.Response) diag.Diagnostics {
	data, err := io.ReadAll(expCertificateResp.Body)
	if err != nil {
		return diag.Errorf("error reading response certificate: %s", err)
	}
	zipReader, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return diag.Errorf("error reading zip certificate: %s", err)
	}
	for _, f := range zipReader.File {
		err := saveCertificateFromZipFile(f, d)
		if err != nil {
			return diag.Errorf("error saving unzip certificate: %s", err)
		}
	}
	return nil
}

func handleJSONCertificateResponse(d *schema.ResourceData, expCertificateResp *http.Response) diag.Diagnostics {
	expCertificateRespBody, err := utils.FlattenResponse(expCertificateResp)
	if err != nil {
		return diag.FromErr(err)
	}
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

func privateCertificateExport(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg                   = meta.(*config.Config)
		product               = "ccm"
		expCertificateHttpUrl = "v1/private-certificates/{certificate_id}/export"
		certId                = d.Get("certificate_id").(string)
		serverType            = d.Get("type").(string)
		isCompressed          = whetherCompressed(serverType)
	)

	client, err := cfg.NewServiceClient(product, cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CCM client: %s", err)
	}

	expCertificatePath := client.Endpoint + expCertificateHttpUrl
	expCertificatePath = strings.ReplaceAll(expCertificatePath, "{certificate_id}", certId)
	expCertificateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	expCertificateOpt.JSONBody = utils.RemoveNil(buildExpCertificateBodyParams(d, serverType, isCompressed))
	expCertificateResp, err := client.Request("POST", expCertificatePath, &expCertificateOpt)
	if err != nil {
		return diag.Errorf("error exporting CCM private certificate: %s", err)
	}
	d.SetId(certId)

	if isCompressed {
		return handleCompressedCertificateResponse(d, expCertificateResp)
	}
	return handleJSONCertificateResponse(d, expCertificateResp)
}
