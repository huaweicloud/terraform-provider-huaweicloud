package cae

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var certificateResourceNotFoundCodes = []string{
	"CAE.01500404", // The environment not found.
	"CAE.01500005", // The resource not found.
}

// @API CAE POST /v1/{project_id}/cae/certificates
// @API CAE GET /v1/{project_id}/cae/certificates
// @API CAE PUT /v1/{project_id}/cae/certificates/{certificate_id}
// @API CAE DELETE /v1/{project_id}/cae/certificates/{certificate_id}
func ResourceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCertificateCreate,
		ReadContext:   resourceCertificateRead,
		UpdateContext: resourceCertificateUpdate,
		DeleteContext: resourceCertificateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCertificateImportState,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"environment_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the CAE environment.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the certificate.`,
			},
			"crt": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The content of the certificate.`,
			},
			"key": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: `The private key of the certificate.`,
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The ID of the enterprise project to which the certificate belongs.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the certificate, in RFC3339 format.`,
			},
		},
	}
}

func buildCertificateBodyParams(d *schema.ResourceData, isCreated bool) map[string]interface{} {
	params := map[string]interface{}{
		"api_version": "v1",
		"kind":        "Certificate",
		"spec": map[string]interface{}{
			"crt": d.Get("crt"),
			"key": d.Get("key"),
		},
	}

	if isCreated {
		params["metadata"] = map[string]interface{}{
			"name": d.Get("name"),
		}
	}
	return params
}

func resourceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		httpUrl       = "v1/{project_id}/cae/certificates"
		environmentId = d.Get("environment_id").(string)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, cfg.GetEnterpriseProjectID(d)),
		JSONBody:         buildCertificateBodyParams(d, true),
	}
	createResp, err := client.Request("POST", createPath, &createOpts)
	if err != nil {
		return diag.Errorf("error creating certificate for environment (%s): %s", environmentId, err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("items|[0].metadata.id", createRespBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("unable to find the certificate ID from the API response")
	}

	d.SetId(certificateId)
	return resourceCertificateRead(ctx, d, meta)
}

func GetCertificateById(client *golangsdk.ServiceClient, environmentId, certificateId, epsId string) (interface{}, error) {
	certificateInfos, err := getCertificates(client, environmentId, epsId)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", envResourceNotFoundCodes...)
	}

	certificateInfo := utils.PathSearch(fmt.Sprintf("items[?metadata.id=='%s']|[0]", certificateId), certificateInfos, nil)
	if certificateInfo == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return certificateInfo, nil
}

func getCertificates(client *golangsdk.ServiceClient, environmentId, epsId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/cae/certificates"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(environmentId, epsId),
	}
	resp, err := client.Request("GET", listPath, &listOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(resp)
}

func resourceCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("cae", region)
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	certificateInfo, err := GetCertificateById(client, d.Get("environment_id").(string), d.Id(), cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving certificate")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", utils.PathSearch("metadata.name", certificateInfo, nil)),
		d.Set("crt", utils.PathSearch("spec.crt", certificateInfo, nil)),
		d.Set("key", utils.PathSearch("spec.key", certificateInfo, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(utils.ConvertTimeStrToNanoTimestamp(utils.PathSearch("spec.created_at",
			certificateInfo, "").(string))/1000, false)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	if d.HasChanges("crt", "key") {
		httpUrl := "v1/{project_id}/cae/certificates/{certificate_id}"
		updatePath := client.Endpoint + httpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{certificate_id}", d.Id())
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
			JSONBody:         buildCertificateBodyParams(d, false),
			OkCodes:          []int{204},
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating certificate (%s): %s", d.Get("name").(string), err)
		}
	}

	return resourceCertificateRead(ctx, d, meta)
}

func resourceCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v1/{project_id}/cae/certificates/{certificate_id}"
	)
	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating CAE client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(d.Get("environment_id").(string), cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", certificateResourceNotFoundCodes...),
			fmt.Sprintf("error deleting certificate (%s)", d.Get("name").(string)))
	}

	return nil
}

func resourceCertificateImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg        = meta.(*config.Config)
		importedId = d.Id()
		parts      = strings.Split(importedId, "/")
		mErr       *multierror.Error
	)

	switch len(parts) {
	case 2:
		mErr = multierror.Append(d.Set("environment_id", parts[0]))
	case 3:
		mErr = multierror.Append(
			d.Set("environment_id", parts[0]),
			d.Set("enterprise_project_id", parts[2]),
		)
	default:
		return nil, fmt.Errorf("invalid format specified for import ID, want '<environment_id>/<name>' or "+
			"'<environment_id>/<name>/<enterprise_project_id>', but got '%s'",
			importedId)
	}

	if mErr.ErrorOrNil() != nil {
		return nil, mErr
	}

	client, err := cfg.NewServiceClient("cae", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating CAE client: %s", err)
	}

	certificates, err := getCertificates(client, parts[0], cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return nil, fmt.Errorf("error retrieving the list of the certificates: %s", err)
	}

	certificateId := utils.PathSearch(fmt.Sprintf("items[?metadata.name=='%s']|[0].metadata.id", parts[1]), certificates, "").(string)
	if certificateId == "" {
		return nil, fmt.Errorf("unable to find the ID of the certificate (%s) from API response : %s", parts[1], err)
	}

	d.SetId(certificateId)
	return []*schema.ResourceData{d}, nil
}
