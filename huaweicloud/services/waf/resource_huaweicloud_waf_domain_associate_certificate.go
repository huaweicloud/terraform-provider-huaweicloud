package waf

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API WAF POST /v1/{project_id}/waf/certificate/{certificate_id}/apply-to-hosts

var associateCertificateNonUpdatableParams = []string{"certificate_id", "cloud_host_ids", "preminum_host_ids", "enterprise_project_id"}

func ResourceDomainAssociateCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDomainAssociateCertificateCreate,
		ReadContext:   resourceDomainAssociateCertificateRead,
		UpdateContext: resourceDomainAssociateCertificateUpdate,
		DeleteContext: resourceDomainAssociateCertificateDelete,

		CustomizeDiff: config.FlexibleForceNew(associateCertificateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"certificate_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cloud_host_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"premium_host_ids": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildDomainAssociateCertificateBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_host_ids":   utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("cloud_host_ids").([]interface{}))),
		"premium_host_ids": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("premium_host_ids").([]interface{}))),
	}

	return bodyParams
}

func resourceDomainAssociateCertificateCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg           = meta.(*config.Config)
		region        = cfg.GetRegion(d)
		certificateId = d.Get("certificate_id").(string)
		epsId         = cfg.GetEnterpriseProjectID(d)
		httpUrl       = "v1/{project_id}/waf/certificate/{certificate_id}/apply-to-hosts"
	)

	client, err := cfg.NewServiceClient("waf", region)
	if err != nil {
		return diag.Errorf("error creating WAF client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{certificate_id}", certificateId)
	if epsId != "" {
		createPath += fmt.Sprintf("?enterprise_project_id=%s", epsId)
	}

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
		},
		JSONBody: utils.RemoveNil(buildDomainAssociateCertificateBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error associating the certificate to domain: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	certId := utils.PathSearch("id", respBody, "").(string)
	if certId == "" {
		return diag.Errorf("unable to find the certificate ID from the API response")
	}

	d.SetId(certId)

	return nil
}

func resourceDomainAssociateCertificateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceDomainAssociateCertificateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a action resource.
	return nil
}

func resourceDomainAssociateCertificateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Delete()' method because the resource is a action resource.
	return nil
}
