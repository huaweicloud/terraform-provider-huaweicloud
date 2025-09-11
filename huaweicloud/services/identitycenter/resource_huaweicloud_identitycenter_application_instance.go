package identitycenter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
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

var identityCenterApplicationInstanceNonUpdateParams = []string{"instance_id", "name", "template_id"}

// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances
// @API IdentityCenter GET /v1/instances/{instance_id}/application-instances/{application_instance_id}
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/display-data
// @API IdentityCenter POST /v1/instances/{instance_id}/application-instances/{application_instance_id}/metadata
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/response-configuration
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/response-schema-configuration
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/service-provider-configuration
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/status
// @API IdentityCenter PUT /v1/instances/{instance_id}/application-instances/{application_instance_id}/security-configuration
// @API IdentityCenter DELETE /v1/instances/{instance_id}/application-instances/{application_instance_id}
func ResourceIdentityCenterApplicationInstance() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceIdentityCenterApplicationInstanceCreate,
		UpdateContext: resourceIdentityCenterApplicationInstanceUpdate,
		ReadContext:   resourceIdentityCenterApplicationInstanceRead,
		DeleteContext: resourceIdentityCenterApplicationInstanceDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceIdentityCenterApplicationInstanceImportState,
		},

		CustomizeDiff: config.FlexibleForceNew(identityCenterApplicationInstanceNonUpdateParams),

		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the IAM Identity Center instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the IAM Identity Center application instance.`,
			},
			"template_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the IAM Identity Center application template.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"description": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"display_name"},
				Description:  `Specifies the description of the IAM Identity Center application instance.`,
			},
			"display_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"description"},
				Description:  `Specifies the display name of the IAM Identity Center application instance.`,
			},
			"metadata": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the string of the application instance metadata.`,
			},
			"status": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the status of the application instance.`,
			},
			"response_config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: configDiffSuppressFunc,
				Description:      `Specifies the response configuration of the IAM Identity Center application instance.`,
			},
			"response_schema_config": {
				Type:             schema.TypeString,
				Optional:         true,
				Computed:         true,
				ValidateFunc:     validation.StringIsJSON,
				DiffSuppressFunc: configDiffSuppressFunc,
				Description:      `Specifies the response schema configuration of the IAM Identity Center application instance.`,
			},
			"security_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the security configuration of the IAM Identity Center application instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"ttl": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the ttl of the IAM Identity Center application instance certificate.`,
						},
					},
				},
			},
			"service_provider_config": {
				Type:        schema.TypeList,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the service provider configuration of the IAM Identity Center application instance.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"audience": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the audience of the IAM Identity Center application instance.`,
						},
						"require_request_signature": {
							Type:        schema.TypeBool,
							Optional:    true,
							Computed:    true,
							Description: `Whether the IAM Identity Center application instance require request signature.`,
						},
						"consumers": {
							Type:     schema.TypeList,
							Optional: true,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"location": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the location url of the IAM Identity Center application instance.`,
									},
									"binding": {
										Type:        schema.TypeString,
										Optional:    true,
										Computed:    true,
										Description: `Specifies the binding method of the IAM Identity Center application instance.`,
									},
									"default_value": {
										Type:        schema.TypeBool,
										Optional:    true,
										Computed:    true,
										Description: `Whether the IAM Identity Center application instance consumer is default.`,
									},
								},
							},
						},
						"start_url": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: `Specifies the start url of the IAM Identity Center application instance.`,
						},
					},
				},
			},
			"identity_provider_config": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"issuer_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the issuer url of the IAM Identity Center.`,
						},
						"metadata_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the metadata url of the IAM Identity Center.`,
						},
						"remote_login_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the remote login url of the IAM Identity Center.`,
						},
						"remote_logout_url": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the remote logout url of the IAM Identity Center.`,
						},
					},
				},
			},
			"active_certificate": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"algorithm": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the algorithm of the IAM Identity Center application instance certificate.`,
						},
						"certificate": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the certificate of the IAM Identity Center application instance.`,
						},
						"certificate_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the ID of the IAM Identity Center application instance certificate.`,
						},
						"expiry_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the expiry date of the IAM Identity Center application instance certificate.`,
						},
						"status": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the status of the IAM Identity Center application instance certificate.`,
						},
						"key_size": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the key size of the IAM Identity Center application instance certificate.`,
						},
						"issue_date": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `Specifies the issue date of the IAM Identity Center application instance certificate.`,
						},
					},
				},
			},
			"visible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the IAM Identity Center application instance is visible.`,
			},
			"client_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the ID of the IAM Identity Center application instance client.`,
			},
			"end_user_visible": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the IAM Identity Center application instance is visible for end user.`,
			},
			"managed_account": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Specifies the ID of the IAM Identity Center application instance account.`,
			},
		},
	}
}

func resourceIdentityCenterApplicationInstanceCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	instanceId := d.Get("instance_id").(string)

	var (
		createApplicationInstanceHttpUrl  = "v1/instances/{instance_id}/application-instances"
		updateDisplayDataHttpUrl          = "v1/instances/{instance_id}/application-instances/{application_instance_id}/display-data"
		importMetadataHttpUrl             = "v1/instances/{instance_id}/application-instances/{application_instance_id}/metadata"
		updateResponseConfigHttpUrl       = "v1/instances/{instance_id}/application-instances/{application_instance_id}/response-configuration"
		updateResponseSchemaConfigHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}" +
			"/response-schema-configuration"
		updateServiceProviderConfigHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}" +
			"/service-provider-configuration"
		updateStatusHttpUrl              = "v1/instances/{instance_id}/application-instances/{application_instance_id}/status"
		updateSecurityConfigHttpUrl      = "v1/instances/{instance_id}/application-instances/{application_instance_id}/security-configuration"
		createApplicationInstanceProduct = "identitycenter"
	)

	client, err := cfg.NewServiceClient(createApplicationInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	createApplicationInstancePath := client.Endpoint + createApplicationInstanceHttpUrl
	createApplicationInstancePath = strings.ReplaceAll(createApplicationInstancePath, "{instance_id}", instanceId)

	createApplicationInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateApplicationInstanceBodyParams(d)),
	}

	createApplicationInstanceResp, err := client.Request("POST",
		createApplicationInstancePath, &createApplicationInstanceOpt)
	if err != nil {
		return diag.Errorf("error creating Identity Center application instance: %s", err)
	}

	createApplicationInstanceRespBody, err := utils.FlattenResponse(createApplicationInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	applicationInstanceId := utils.PathSearch("application_instance.application_instance_id", createApplicationInstanceRespBody, "").(string)
	if applicationInstanceId == "" {
		return diag.Errorf("unable to find the Identity Center application instance ID from the API response")
	}
	d.SetId(applicationInstanceId)

	if _, ok := d.GetOk("display_name"); ok {
		if _, ok := d.GetOk("description"); ok {
			updateDisplayDataPath := client.Endpoint + updateDisplayDataHttpUrl
			updateDisplayDataPath = strings.ReplaceAll(updateDisplayDataPath, "{instance_id}", d.Get("instance_id").(string))
			updateDisplayDataPath = strings.ReplaceAll(updateDisplayDataPath, "{application_instance_id}", d.Id())

			updateDisplayDataOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				JSONBody:         utils.RemoveNil(buildUpdateDisplayDataBodyParams(d)),
			}

			_, err = client.Request("PUT", updateDisplayDataPath, &updateDisplayDataOpt)
			if err != nil {
				return diag.Errorf("error updating IdentityCenter application instance display data: %s", err)
			}
		}
	}

	if _, ok := d.GetOk("metadata"); ok {
		importMetadataPath := client.Endpoint + importMetadataHttpUrl
		importMetadataPath = strings.ReplaceAll(importMetadataPath, "{instance_id}", d.Get("instance_id").(string))
		importMetadataPath = strings.ReplaceAll(importMetadataPath, "{application_instance_id}", d.Id())

		importMetadataOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildImportMetadataBodyParams(d)),
		}

		_, err = client.Request("POST", importMetadataPath, &importMetadataOpt)
		if err != nil {
			return diag.Errorf("error importing IdentityCenter application instance metadata: %s", err)
		}
	}

	if _, ok := d.GetOk("status"); ok {
		updateStatusPath := client.Endpoint + updateStatusHttpUrl
		updateStatusPath = strings.ReplaceAll(updateStatusPath, "{instance_id}", d.Get("instance_id").(string))
		updateStatusPath = strings.ReplaceAll(updateStatusPath, "{application_instance_id}", d.Id())

		updateStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateStatusBodyParams(d)),
		}

		_, err = client.Request("PUT", updateStatusPath, &updateStatusOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance status: %s", err)
		}
	}

	if _, ok := d.GetOk("security_config"); ok {
		updateSecurityConfigPath := client.Endpoint + updateSecurityConfigHttpUrl
		updateSecurityConfigPath = strings.ReplaceAll(updateSecurityConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateSecurityConfigPath = strings.ReplaceAll(updateSecurityConfigPath, "{application_instance_id}", d.Id())

		updateSecurityConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateSecurityConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateSecurityConfigPath, &updateSecurityConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance security configuration: %s", err)
		}
	}

	if _, ok := d.GetOk("response_config"); ok {
		updateResponseConfigPath := client.Endpoint + updateResponseConfigHttpUrl
		updateResponseConfigPath = strings.ReplaceAll(updateResponseConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateResponseConfigPath = strings.ReplaceAll(updateResponseConfigPath, "{application_instance_id}", d.Id())

		updateResponseConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateResponseConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateResponseConfigPath, &updateResponseConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance response configuration: %s", err)
		}
	}

	if _, ok := d.GetOk("response_schema_config"); ok {
		updateResponseSchemaConfigPath := client.Endpoint + updateResponseSchemaConfigHttpUrl
		updateResponseSchemaConfigPath = strings.ReplaceAll(updateResponseSchemaConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateResponseSchemaConfigPath = strings.ReplaceAll(updateResponseSchemaConfigPath, "{application_instance_id}", d.Id())

		updateResponseSchemaConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateResponseSchemaConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateResponseSchemaConfigPath, &updateResponseSchemaConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance response schema configuration: %s", err)
		}
	}

	if _, ok := d.GetOk("service_provider_config"); ok {
		updateServiceProviderConfigPath := client.Endpoint + updateServiceProviderConfigHttpUrl
		updateServiceProviderConfigPath = strings.ReplaceAll(updateServiceProviderConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateServiceProviderConfigPath = strings.ReplaceAll(updateServiceProviderConfigPath, "{application_instance_id}", d.Id())

		updateServiceProviderConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateServiceProviderConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateServiceProviderConfigPath, &updateServiceProviderConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance service provider configuration: %s", err)
		}
	}

	return resourceIdentityCenterApplicationInstanceRead(ctx, d, meta)
}

func buildCreateApplicationInstanceBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"template_id": utils.ValueIgnoreEmpty(d.Get("template_id")),
	}
	return bodyParams
}

func resourceIdentityCenterApplicationInstanceRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		getApplicationInstanceHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}"
		getApplicationInstanceProduct = "identitycenter"
	)
	client, err := cfg.NewServiceClient(getApplicationInstanceProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	getApplicationInstancePath := client.Endpoint + getApplicationInstanceHttpUrl
	getApplicationInstancePath = strings.ReplaceAll(getApplicationInstancePath, "{instance_id}", d.Get("instance_id").(string))
	getApplicationInstancePath = strings.ReplaceAll(getApplicationInstancePath, "{application_instance_id}", d.Id())

	getApplicationInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getApplicationInstanceResp, err := client.Request("GET",
		getApplicationInstancePath, &getApplicationInstanceOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Identity Center application instance.")
	}

	getApplicationInstanceRespBody, err := utils.FlattenResponse(getApplicationInstanceResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("application_instance.name", getApplicationInstanceRespBody, nil)),
		d.Set("display_name", utils.PathSearch("application_instance.display.display_name", getApplicationInstanceRespBody, nil)),
		d.Set("description", utils.PathSearch("application_instance.display.description", getApplicationInstanceRespBody, nil)),
		d.Set("active_certificate",
			flattenActiveCertificate(utils.PathSearch("application_instance.active_certificate", getApplicationInstanceRespBody, nil))),
		d.Set("identity_provider_config",
			flattenIdentityProviderConfig(utils.PathSearch("application_instance.identity_provider_config", getApplicationInstanceRespBody, nil))),
		d.Set("visible", utils.PathSearch("application_instance.visible", getApplicationInstanceRespBody, nil)),
		d.Set("status", utils.PathSearch("application_instance.status", getApplicationInstanceRespBody, nil)),
		d.Set("client_id", utils.PathSearch("application_instance.client_id", getApplicationInstanceRespBody, nil)),
		d.Set("end_user_visible", utils.PathSearch("application_instance.end_user_visible", getApplicationInstanceRespBody, nil)),
		d.Set("managed_account", utils.PathSearch("application_instance.managed_account", getApplicationInstanceRespBody, nil)),
		d.Set("security_config",
			flattenSecurityConfig(utils.PathSearch("application_instance.security_config", getApplicationInstanceRespBody, nil))),
		d.Set("service_provider_config",
			flattenServiceProviderConfig(utils.PathSearch("application_instance.service_provider_config", getApplicationInstanceRespBody, nil))),
		d.Set("response_config",
			marshalJsonFormatParams("response config",
				utils.PathSearch("application_instance.response_config", getApplicationInstanceRespBody, nil))),
		d.Set("response_schema_config",
			marshalJsonFormatParams("response schema config",
				utils.PathSearch("application_instance.response_schema_config", getApplicationInstanceRespBody, nil))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenActiveCertificate(data interface{}) []map[string]interface{} {
	if data == nil || len(data.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"algorithm":      utils.PathSearch("algorithm", data, nil),
			"certificate":    utils.PathSearch("certificate", data, nil),
			"certificate_id": utils.PathSearch("certificate_id", data, nil),
			"status":         utils.PathSearch("status", data, nil),
			"key_size":       utils.PathSearch("key_size", data, nil),
			"expiry_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("expiry_date", data, float64(0)).(float64))/1000, false),
			"issue_date": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("issue_date", data, float64(0)).(float64))/1000, false),
		},
	}
}

func flattenIdentityProviderConfig(data interface{}) []map[string]interface{} {
	if data == nil || len(data.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"issuer_url":        utils.PathSearch("issuer_url", data, nil),
			"metadata_url":      utils.PathSearch("metadata_url", data, nil),
			"remote_login_url":  utils.PathSearch("remote_login_url", data, nil),
			"remote_logout_url": utils.PathSearch("remote_logout_url", data, nil),
		},
	}
}

func flattenSecurityConfig(data interface{}) []map[string]interface{} {
	if data == nil || len(data.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"ttl": utils.PathSearch("ttl", data, nil),
		},
	}
}

func flattenServiceProviderConfig(data interface{}) []map[string]interface{} {
	if data == nil || len(data.(map[string]interface{})) == 0 {
		return nil
	}

	return []map[string]interface{}{
		{
			"audience":                  utils.PathSearch("audience", data, nil),
			"require_request_signature": utils.PathSearch("require_request_signature", data, nil),
			"consumers":                 flattenConsumers(utils.PathSearch("consumers", data, nil)),
			"start_url":                 utils.PathSearch("start_url", data, nil),
		},
	}
}

func flattenConsumers(data interface{}) []map[string]interface{} {
	if data == nil || len(data.([]interface{})) == 0 {
		return nil
	}
	consumer := utils.PathSearch("[0]", data, nil)
	return []map[string]interface{}{
		{
			"location":      utils.PathSearch("location", consumer, nil),
			"binding":       utils.PathSearch("binding", consumer, nil),
			"default_value": utils.PathSearch("default_value", consumer, nil),
		},
	}
}

func resourceIdentityCenterApplicationInstanceUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		updateDisplayDataHttpUrl          = "v1/instances/{instance_id}/application-instances/{application_instance_id}/display-data"
		importMetadataHttpUrl             = "v1/instances/{instance_id}/application-instances/{application_instance_id}/metadata"
		updateResponseConfigHttpUrl       = "v1/instances/{instance_id}/application-instances/{application_instance_id}/response-configuration"
		updateResponseSchemaConfigHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}" +
			"/response-schema-configuration"
		updateServiceProviderConfigHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}" +
			"/service-provider-configuration"
		updateStatusHttpUrl         = "v1/instances/{instance_id}/application-instances/{application_instance_id}/status"
		updateSecurityConfigHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}/security-configuration"
		updateProduct               = "identitycenter"
	)
	client, err := cfg.NewServiceClient(updateProduct, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	if d.HasChanges("display_name", "description") {
		updateDisplayDataPath := client.Endpoint + updateDisplayDataHttpUrl
		updateDisplayDataPath = strings.ReplaceAll(updateDisplayDataPath, "{instance_id}", d.Get("instance_id").(string))
		updateDisplayDataPath = strings.ReplaceAll(updateDisplayDataPath, "{application_instance_id}", d.Id())

		updateDisplayDataOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateDisplayDataBodyParams(d)),
		}

		_, err = client.Request("PUT", updateDisplayDataPath, &updateDisplayDataOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance display data: %s", err)
		}
	}

	if d.HasChange("metadata") {
		importMetadataPath := client.Endpoint + importMetadataHttpUrl
		importMetadataPath = strings.ReplaceAll(importMetadataPath, "{instance_id}", d.Get("instance_id").(string))
		importMetadataPath = strings.ReplaceAll(importMetadataPath, "{application_instance_id}", d.Id())

		importMetadataOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildImportMetadataBodyParams(d)),
		}

		_, err = client.Request("POST", importMetadataPath, &importMetadataOpt)
		if err != nil {
			return diag.Errorf("error importing IdentityCenter application instance metadata: %s", err)
		}
	}

	if d.HasChange("status") {
		updateStatusPath := client.Endpoint + updateStatusHttpUrl
		updateStatusPath = strings.ReplaceAll(updateStatusPath, "{instance_id}", d.Get("instance_id").(string))
		updateStatusPath = strings.ReplaceAll(updateStatusPath, "{application_instance_id}", d.Id())

		updateStatusOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateStatusBodyParams(d)),
		}

		_, err = client.Request("PUT", updateStatusPath, &updateStatusOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance status: %s", err)
		}
	}

	if d.HasChange("security_config") {
		updateSecurityConfigPath := client.Endpoint + updateSecurityConfigHttpUrl
		updateSecurityConfigPath = strings.ReplaceAll(updateSecurityConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateSecurityConfigPath = strings.ReplaceAll(updateSecurityConfigPath, "{application_instance_id}", d.Id())

		updateSecurityConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateSecurityConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateSecurityConfigPath, &updateSecurityConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance security configuration: %s", err)
		}
	}

	if d.HasChange("response_config") {
		updateResponseConfigPath := client.Endpoint + updateResponseConfigHttpUrl
		updateResponseConfigPath = strings.ReplaceAll(updateResponseConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateResponseConfigPath = strings.ReplaceAll(updateResponseConfigPath, "{application_instance_id}", d.Id())

		updateResponseConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateResponseConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateResponseConfigPath, &updateResponseConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance response configuration: %s", err)
		}
	}

	if d.HasChange("response_schema_config") {
		updateResponseSchemaConfigPath := client.Endpoint + updateResponseSchemaConfigHttpUrl
		updateResponseSchemaConfigPath = strings.ReplaceAll(updateResponseSchemaConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateResponseSchemaConfigPath = strings.ReplaceAll(updateResponseSchemaConfigPath, "{application_instance_id}", d.Id())

		updateResponseSchemaConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateResponseSchemaConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateResponseSchemaConfigPath, &updateResponseSchemaConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance response schema configuration: %s", err)
		}
	}

	if d.HasChange("service_provider_config") {
		updateServiceProviderConfigPath := client.Endpoint + updateServiceProviderConfigHttpUrl
		updateServiceProviderConfigPath = strings.ReplaceAll(updateServiceProviderConfigPath, "{instance_id}", d.Get("instance_id").(string))
		updateServiceProviderConfigPath = strings.ReplaceAll(updateServiceProviderConfigPath, "{application_instance_id}", d.Id())

		updateServiceProviderConfigOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			JSONBody:         utils.RemoveNil(buildUpdateServiceProviderConfigBodyParams(d)),
		}

		_, err = client.Request("PUT", updateServiceProviderConfigPath, &updateServiceProviderConfigOpt)
		if err != nil {
			return diag.Errorf("error updating IdentityCenter application instance service provider configuration: %s", err)
		}
	}

	return resourceIdentityCenterApplicationInstanceRead(ctx, d, meta)
}

func buildUpdateDisplayDataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"display_name": utils.ValueIgnoreEmpty(d.Get("display_name")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func buildImportMetadataBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"metadata": utils.ValueIgnoreEmpty(d.Get("metadata")),
	}
	return bodyParams
}

func buildUpdateStatusBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"status": utils.ValueIgnoreEmpty(d.Get("status")),
	}
	return bodyParams
}

func buildUpdateSecurityConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	securityConfigs := d.Get("security_config").([]interface{})
	if len(securityConfigs) == 0 {
		return nil
	}
	bodyParams := map[string]interface{}{
		"security_config": map[string]interface{}{
			"ttl": utils.ValueIgnoreEmpty(buildUpdateSecurityConfigTTLBodyParams(securityConfigs)),
		},
	}
	return bodyParams
}

func buildUpdateSecurityConfigTTLBodyParams(data []interface{}) interface{} {
	securityConfig := data[0].(map[string]interface{})
	if _, ok := securityConfig["ttl"]; ok {
		return securityConfig["ttl"].(string)
	}
	return nil
}

func buildUpdateResponseConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	updateValue := d.Get("response_config").(string)
	bodyParams := map[string]interface{}{
		"response_config": utils.ValueIgnoreEmpty(unmarshalJsonFormatParams("response config", updateValue)),
	}
	return bodyParams
}

func buildUpdateResponseSchemaConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	updateValue := d.Get("response_schema_config").(string)
	bodyParams := map[string]interface{}{
		"response_schema_config": utils.ValueIgnoreEmpty(unmarshalJsonFormatParams("response schema config", updateValue)),
	}
	return bodyParams
}

func buildUpdateServiceProviderConfigBodyParams(d *schema.ResourceData) map[string]interface{} {
	configs := d.Get("service_provider_config").([]interface{})
	if len(configs) == 0 {
		return nil
	}
	bodyParams := map[string]interface{}{
		"service_provider_config": map[string]interface{}{
			"audience":                  utils.ValueIgnoreEmpty(buildServiceProviderConfigAudienceBodyParams(configs)),
			"require_request_signature": utils.ValueIgnoreEmpty(buildServiceProviderConfigRequireRequestSignatureBodyParams(configs)),
			"consumers":                 utils.ValueIgnoreEmpty(buildServiceProviderConfigConsumersBodyParams(configs)),
			"start_url":                 utils.ValueIgnoreEmpty(buildServiceProviderConfigStartUrlBodyParams(configs)),
		},
	}
	return bodyParams
}

func buildServiceProviderConfigAudienceBodyParams(data []interface{}) interface{} {
	spConfig := data[0].(map[string]interface{})
	if _, ok := spConfig["audience"]; ok {
		return spConfig["audience"].(string)
	}
	return nil
}

func buildServiceProviderConfigRequireRequestSignatureBodyParams(data []interface{}) interface{} {
	spConfig := data[0].(map[string]interface{})
	if _, ok := spConfig["require_request_signature"]; ok {
		return spConfig["require_request_signature"].(bool)
	}
	return nil
}

func buildServiceProviderConfigStartUrlBodyParams(data []interface{}) interface{} {
	spConfig := data[0].(map[string]interface{})
	if _, ok := spConfig["start_url"]; ok {
		return spConfig["start_url"].(string)
	}
	return nil
}

func buildServiceProviderConfigConsumersBodyParams(data []interface{}) []map[string]interface{} {
	spConfig := data[0].(map[string]interface{})
	if len(spConfig["consumers"].([]interface{})) == 0 {
		return nil
	}
	consumer := spConfig["consumers"].([]interface{})[0].(map[string]interface{})
	var location interface{}
	var binding interface{}
	var defaultValue interface{}
	if _, ok := consumer["location"]; ok {
		location = consumer["location"].(string)
	}
	if _, ok := consumer["binding"]; ok {
		binding = consumer["binding"].(string)
	}
	if _, ok := consumer["default_value"]; ok {
		defaultValue = consumer["default_value"].(bool)
	}
	return []map[string]interface{}{
		{
			"location":      location,
			"binding":       binding,
			"default_value": defaultValue,
		},
	}
}

func unmarshalJsonFormatParams(paramName, paramVal string) map[string]interface{} {
	parseResult := make(map[string]interface{})
	err := json.Unmarshal([]byte(paramVal), &parseResult)
	if err != nil {
		log.Printf("[ERROR] Invalid type of the %s, not json format.", paramName)
	}
	return parseResult
}

func marshalJsonFormatParams(paramName string, paramVal interface{}) interface{} {
	jsonFilter, err := json.Marshal(paramVal)
	if err != nil {
		log.Printf("[ERROR] unable to convert the %s, not json format.", paramName)
	}
	return string(jsonFilter)
}

func resourceIdentityCenterApplicationInstanceDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		deleteApplicationInstanceHttpUrl = "v1/instances/{instance_id}/application-instances/{application_instance_id}"
		product                          = "identitycenter"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating IdentityCenter client: %s", err)
	}

	deleteApplicationInstancePath := client.Endpoint + deleteApplicationInstanceHttpUrl
	deleteApplicationInstancePath = strings.ReplaceAll(deleteApplicationInstancePath, "{instance_id}", d.Get("instance_id").(string))
	deleteApplicationInstancePath = strings.ReplaceAll(deleteApplicationInstancePath, "{application_instance_id}", d.Id())

	deleteApplicationInstanceOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", deleteApplicationInstancePath, &deleteApplicationInstanceOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting IdentityCenter application instance.")
	}

	return nil
}

func resourceIdentityCenterApplicationInstanceImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid id format, must be <instance_id>/<application_instance_id>")
	}
	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return nil, fmt.Errorf("failed to set value to state when import, %s", err)
	}
	return []*schema.ResourceData{d}, nil
}

func configDiffSuppressFunc(_, o, n string, _ *schema.ResourceData) bool {
	equal, _ := utils.CompareJsonTemplateAreEquivalent(o, n)
	return equal
}
