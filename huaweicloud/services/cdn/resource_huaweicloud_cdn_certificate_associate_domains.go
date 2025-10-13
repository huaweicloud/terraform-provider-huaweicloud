package cdn

import (
	"context"
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

// @API CDN PUT /v1.0/cdn/domains/config-https-info
func ResourceCdnCertificateAssociateDomains() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdnCertificateAssociateDomainsCreate,
		ReadContext:   resourceCdnCertificateAssociateDomainsRead,
		UpdateContext: resourceCdnCertificateAssociateDomainsUpdate,
		DeleteContext: resourceCdnCertificateAssociateDomainsDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the CDN service is located.`,
			},

			// Required parameters.
			"domain_names": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The domain names to associate with the certificate, separated by commas.`,
			},
			"https_switch": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `The HTTPS certificate configuration switch.`,
			},

			// Optional parameters.
			"access_origin_way": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     2,
				Description: `The origin protocol configuration.`,
			},
			"force_redirect_https": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: `Whether to enable HTTPS force redirect.`,
			},
			"force_redirect_config": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"switch": {
							Type:        schema.TypeInt,
							Required:    true,
							Description: `The force redirect switch.`,
						},
						"redirect_type": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The redirect type.`,
						},
					},
				},
				Description: `The force redirect configuration.`,
			},
			"http2": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: `The HTTP/2 switch.`,
			},
			"cert_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The certificate name.`,
			},
			"certificate": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The SSL certificate content.`,
			},
			"private_key": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Sensitive:   true,
				Description: `The SSL certificate private key content.`,
			},
			"certificate_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Default:     0,
				Description: `The certificate type.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildForceRedirectConfig(forceRedirectConfigs []interface{}) map[string]interface{} {
	if len(forceRedirectConfigs) < 1 {
		return nil
	}

	forceRedirectConfig := forceRedirectConfigs[0]
	return map[string]interface{}{
		"switch":        utils.PathSearch("switch", forceRedirectConfig, nil),
		"redirect_type": utils.PathSearch("redirect_type", forceRedirectConfig, nil),
	}
}

func buildCdnCertificateAssociateDomainsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"domain_name":     d.Get("domain_names"),
		"https_switch":    d.Get("https_switch"),
		"access_origin_way": utils.ValueIgnoreEmpty(d.Get("access_origin_way")),
		"force_redirect_https": utils.ValueIgnoreEmpty(d.Get("force_redirect_https")),
		"http2":           utils.ValueIgnoreEmpty(d.Get("http2")),
		"cert_name":       utils.ValueIgnoreEmpty(d.Get("cert_name")),
		"certificate":     utils.ValueIgnoreEmpty(d.Get("certificate")),
		"private_key":     utils.ValueIgnoreEmpty(d.Get("private_key")),
		"certificate_type": utils.ValueIgnoreEmpty(d.Get("certificate_type")),
	}

	if v, ok := d.GetOk("force_redirect_config"); ok {
		bodyParams["force_redirect_config"] = buildForceRedirectConfig(v.([]interface{}))
	}

	return bodyParams
}

func resourceCdnCertificateAssociateDomainsCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cdn", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	httpUrl := "v1.0/cdn/domains/config-https-info"
	createPath := client.Endpoint + httpUrl

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"https": buildCdnCertificateAssociateDomainsBodyParams(d),
		},
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error associating CDN certificate with domains: %s", err)
	}

	// Set ID as the domain names for this action resource
	d.SetId(d.Get("domain_names").(string))

	return resourceCdnCertificateAssociateDomainsRead(ctx, d, meta)
}

func resourceCdnCertificateAssociateDomainsRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cdn", region)
	if err != nil {
		return diag.Errorf("error creating CDN client: %s", err)
	}

	// For action resources, we don't need to read the state from the server
	// as this is a one-time action. We just return the current state.
	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("domain_names", d.Id()),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCdnCertificateAssociateDomainsUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is an action resource, updates are not supported
	return diag.Errorf("updates are not supported for this action resource")
}

func resourceCdnCertificateAssociateDomainsDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	// This is an action resource, deletion is not supported
	// The action has already been performed and cannot be undone
	return nil
}
