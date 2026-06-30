package das

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var emailsBatchSendNonUpdatableParams = []string{
	"task_ids",
	"email",
	"topic",
	"topic_urn",
	"obs_bucket_name",
	"service_uri",
	"access_key",
	"secret_key",
}

// @API DAS POST /v3/{project_id}/batch-inspection/batch-send-email
// @API DAS POST /v3/{project_id}/batch-inspection/check-credential
// @API DAS POST /v3/{project_id}/batch-inspection/save-credential
func ResourceEmailsBatchSend() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEmailsBatchSendCreate,
		ReadContext:   resourceEmailsBatchSendRead,
		UpdateContext: resourceEmailsBatchSendUpdate,
		DeleteContext: resourceEmailsBatchSendDelete,

		CustomizeDiff: config.FlexibleForceNew(emailsBatchSendNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the email templates are located.`,
			},

			// Required parameters.
			"task_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The list of report IDs.`,
			},

			// Optional parameters.
			"email": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The email address.`,
			},
			"topic": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The topic ID.`,
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The topic URN.`,
			},
			"obs_bucket_name": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"access_key", "secret_key"},
				Description:  `The OBS bucket name.`,
			},
			"service_uri": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The service URI of the OBS bucket.`,
			},
			"access_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The access key used to access the OBS bucket.`,
			},
			"secret_key": {
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				Description: `The secret key used to access the OBS bucket.`,
			},

			// Internal
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{Internal: true},
				),
			},
		},
	}
}

func buildEmailsBatchSendBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"task_id_list":    d.Get("task_ids"),
		"email":           utils.ValueIgnoreEmpty(d.Get("email")),
		"topic":           utils.ValueIgnoreEmpty(d.Get("topic")),
		"topic_urn":       utils.ValueIgnoreEmpty(d.Get("topic_urn")),
		"obs_bucket_name": utils.ValueIgnoreEmpty(d.Get("obs_bucket_name")),
		"service_uri":     utils.ValueIgnoreEmpty(d.Get("service_uri")),
	}
}

func checkCredential(client *golangsdk.ServiceClient, bucketName, ak, sk string) error {
	httpUrl := "v3/{project_id}/batch-inspection/check-credential"
	checkPath := client.Endpoint + httpUrl
	checkPath = strings.ReplaceAll(checkPath, "{project_id}", client.ProjectID)

	checkOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"bucket_name": bucketName,
			"ak":          ak,
			"sk":          sk,
		},
	}

	resp, err := client.Request("POST", checkPath, &checkOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	checkResult := utils.PathSearch("check_result", respBody, false).(bool)
	if !checkResult {
		// If `check_result` is false, the API still return HTTP 200
		return fmt.Errorf("%s", utils.PathSearch("error_msg", respBody, "").(string))
	}
	return nil
}

func saveCredential(client *golangsdk.ServiceClient, ak, sk string) error {
	httpUrl := "v3/{project_id}/batch-inspection/save-credential"
	savePath := client.Endpoint + httpUrl
	savePath = strings.ReplaceAll(savePath, "{project_id}", client.ProjectID)

	saveOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"ak": ak,
			"sk": sk,
		},
	}

	resp, err := client.Request("POST", savePath, &saveOpt)
	if err != nil {
		return err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return err
	}

	checkResult := utils.PathSearch("success", respBody, false).(bool)
	if !checkResult {
		// If `success` is false, the API still return HTTP 200
		return fmt.Errorf("%s", utils.PathSearch("error_msg", respBody, "").(string))
	}
	return nil
}

func resourceEmailsBatchSendCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	if v, ok := d.GetOk("obs_bucket_name"); ok {
		obsBucketName := v.(string)
		accessKey := d.Get("access_key").(string)
		secretKey := d.Get("secret_key").(string)

		if accessKey == "" || secretKey == "" {
			return diag.Errorf("the parameters 'access_key' and 'secret_key' are required when 'obs_bucket_name' is set")
		}

		if err := checkCredential(client, obsBucketName, accessKey, secretKey); err != nil {
			return diag.Errorf("error checking credential for DAS batch email: %s", err)
		}

		if err := saveCredential(client, accessKey, secretKey); err != nil {
			return diag.Errorf("error saving credential for DAS batch email: %s", err)
		}
	}

	httpUrl := "v3/{project_id}/batch-inspection/batch-send-email"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildEmailsBatchSendBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error sending DAS batch email: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	return resourceEmailsBatchSendRead(ctx, d, meta)
}

func resourceEmailsBatchSendRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEmailsBatchSendUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEmailsBatchSendDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource for batch sending emails. Deleting this resource
will not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
