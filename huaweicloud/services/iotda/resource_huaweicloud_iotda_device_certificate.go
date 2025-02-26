package iotda

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

// @API IoTDA POST /v5/iot/{project_id}/certificates
// @API IoTDA GET /v5/iot/{project_id}/certificates
// @API IoTDA POST /v5/iot/{project_id}/certificates/{certificate_id}/action
// @API IoTDA DELETE /v5/iot/{project_id}/certificates/{certificate_id}
func ResourceDeviceCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceCertificateCreate,
		ReadContext:   resourceDeviceCertificateRead,
		UpdateContext: resourceDeviceCertificateUpdate,
		DeleteContext: resourceDeviceCertificateDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"space_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"verify_content": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cn": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"verify_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"effective_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"expiry_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildDeviceCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	certificateParams := map[string]interface{}{
		"content": d.Get("content"),
		"app_id":  utils.ValueIgnoreEmpty(d.Get("space_id")),
	}

	return certificateParams
}

func resourceDeviceCertificateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/certificates"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildDeviceCertificateBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating IoTDA device CA certificate: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	certificateId := utils.PathSearch("certificate_id", respBody, "").(string)
	if certificateId == "" {
		return diag.Errorf("error creating IoTDA device CA certificate: ID is not found in API response")
	}

	d.SetId(certificateId)

	return resourceDeviceCertificateRead(ctx, d, meta)
}

func resourceDeviceCertificateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	respBody, err := QueryDeviceCertificate(client, d.Id(), d.Get("space_id").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving IoTDA device CA certificate")
	}

	status := utils.PathSearch("status", respBody, false).(bool)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("cn", utils.PathSearch("cn_name", respBody, nil)),
		d.Set("owner", utils.PathSearch("owner", respBody, nil)),
		d.Set("status", flattenStatus(status)),
		d.Set("verify_code", utils.PathSearch("verify_code", respBody, nil)),
		d.Set("effective_date", utils.PathSearch("effective_date", respBody, nil)),
		d.Set("expiry_date", utils.PathSearch("expiry_date", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenStatus(status bool) string {
	if status {
		return "Verified"
	}

	return "Unverified"
}

func resourceDeviceCertificateUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/certificates/{certificate_id}/action?action_id=verify"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	verifyPath := client.Endpoint + httpUrl
	verifyPath = strings.ReplaceAll(verifyPath, "{project_id}", client.ProjectID)
	verifyPath = strings.ReplaceAll(verifyPath, "{certificate_id}", d.Id())
	verifyOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"verify_content": d.Get("verify_content"),
		},
	}

	_, err = client.Request("POST", verifyPath, &verifyOpts)
	if err != nil {
		return diag.Errorf("error verifing IoTDA device CA certificate: %s", err)
	}

	return resourceDeviceCertificateRead(ctx, d, meta)
}

func resourceDeviceCertificateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		isDerived = WithDerivedAuth(cfg, region)
		httpUrl   = "v5/iot/{project_id}/certificates/{certificate_id}"
	)

	client, err := cfg.NewServiceClientWithDerivedAuth("iotda", region, isDerived)
	if err != nil {
		return diag.Errorf("error creating IoTDA client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{certificate_id}", d.Id())
	deleteOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpts)
	if err != nil {
		// When the resource does not exist, delete API will return `404`.
		return common.CheckDeletedDiag(d, err, "error deleting IoTDA device CA certificate")
	}

	return nil
}

func QueryDeviceCertificate(client *golangsdk.ServiceClient, certificateId, spaceId string) (interface{}, error) {
	var (
		httpUrl    = "v5/iot/{project_id}/certificates?limit=50"
		nextMarker string
	)

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	if spaceId != "" {
		getPath = fmt.Sprintf("%s&app_id=%s", getPath, spaceId)
	}
	getOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		getPathWithMarker := getPath
		if nextMarker != "" {
			getPathWithMarker = fmt.Sprintf("%s&marker=%s", getPathWithMarker, nextMarker)
		}

		resp, err := client.Request("GET", getPathWithMarker, &getOpts)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		certificates := utils.PathSearch("certificates", respBody, make([]interface{}, 0)).([]interface{})
		if len(certificates) == 0 {
			break
		}

		certificate := utils.PathSearch(fmt.Sprintf("[?certificate_id=='%s']|[0]", certificateId), certificates, nil)
		if certificate != nil {
			return certificate, nil
		}

		nextMarker = utils.PathSearch("page.marker", respBody, "").(string)
		if nextMarker == "" {
			break
		}
	}

	return nil, golangsdk.ErrDefault404{}
}
