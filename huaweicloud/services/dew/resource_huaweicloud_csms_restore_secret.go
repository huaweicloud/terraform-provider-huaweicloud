package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1/{project_id}/secrets/restore
func ResourceRestoreSecret() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceRestoreSecretCreate,
		ReadContext:   resourceRestoreSecretRead,
		UpdateContext: resourceRestoreSecretUpdate,
		DeleteContext: resourceRestoreSecretDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"secret_blob",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secret_blob": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"scheduled_delete_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_rotation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rotation_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"next_rotation_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"event_subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_func_urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceRestoreSecretCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v1/{project_id}/secrets/restore"
		product    = "kms"
		secretBlob = d.Get("secret_blob").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestBody := map[string]interface{}{
		"secret_blob": secretBlob,
	}

	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error restoring DEW CSMS secret: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	secretId := utils.PathSearch("secret.id", respBody, "").(string)
	if secretId == "" {
		return diag.Errorf("error restoring DEW CSMS secret: secret ID is not found in API response")
	}

	d.SetId(secretId)

	mErr := multierror.Append(
		d.Set("name", utils.PathSearch("secret.name", respBody, nil)),
		d.Set("state", utils.PathSearch("secret.state", respBody, nil)),
		d.Set("kms_key_id", utils.PathSearch("secret.kms_key_id", respBody, nil)),
		d.Set("description", utils.PathSearch("secret.description", respBody, nil)),
		d.Set("create_time", utils.PathSearch("secret.create_time", respBody, nil)),
		d.Set("update_time", utils.PathSearch("secret.update_time", respBody, nil)),
		d.Set("scheduled_delete_time", utils.PathSearch("secret.scheduled_delete_time", respBody, nil)),
		d.Set("secret_type", utils.PathSearch("secret.secret_type", respBody, nil)),
		d.Set("auto_rotation", utils.PathSearch("secret.auto_rotation", respBody, nil)),
		d.Set("rotation_period", utils.PathSearch("secret.rotation_period", respBody, nil)),
		d.Set("rotation_config", utils.PathSearch("secret.rotation_config", respBody, nil)),
		d.Set("rotation_time", utils.PathSearch("secret.rotation_time", respBody, nil)),
		d.Set("next_rotation_time", utils.PathSearch("secret.next_rotation_time", respBody, nil)),
		d.Set("event_subscriptions", utils.PathSearch("secret.event_subscriptions", respBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("secret.enterprise_project_id", respBody, nil)),
		d.Set("rotation_func_urn", utils.PathSearch("secret.rotation_func_urn", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceRestoreSecretRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Read()' method because the resource is a one-time action resource.
	return nil
}

func resourceRestoreSecretUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in the 'Update()' method because the resource is a one-time action resource.
	return nil
}

func resourceRestoreSecretDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to restore a secret from a backup blob.
Deleting this resource will not delete the restored secret, but will only remove the resource information from
the tfstate file. The restored secret will continue to exist in the DEW service.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
