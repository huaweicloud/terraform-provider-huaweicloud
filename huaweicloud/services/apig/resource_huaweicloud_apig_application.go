package apig

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/apigw/dedicated/v2/applications"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

type SecretAction string

const (
	SecretActionReset SecretAction = "RESET"
)

// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes/{app_code_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-codes
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}
// @API APIG POST /v2/{project_id}/apigw/instances/{instance_id}/apps
// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/secret/{app_id}
func ResourceApigApplicationV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationCreate,
		ReadContext:   resourceApplicationRead,
		UpdateContext: resourceApplicationUpdate,
		DeleteContext: resourceApplicationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the application is located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the application belongs.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The application name.",
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The application description.",
			},
			"app_codes": {
				Type:        schema.TypeSet,
				Optional:    true,
				Computed:    true,
				MaxItems:    5,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The array of one or more application codes that the application has.",
			},
			"secret_action": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(SecretActionReset),
				}, false),
				Description: "The secret action to be done for the application.",
			},
			"registration_time": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The registration time.",
			},
			"app_key": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The APP key.",
			},
			"app_secret": {
				Type:        schema.TypeString,
				Computed:    true,
				Sensitive:   true,
				Description: "The APP secret.",
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The latest update time of the application.",
			},
		},
	}
}

func createApplicationCodes(client *golangsdk.ServiceClient, instanceId, appId string, codes []interface{}) error {
	for _, code := range codes {
		opt := applications.AppCodeOpts{
			AppCode: code.(string),
		}
		_, err := applications.CreateAppCode(client, instanceId, appId, opt).Extract()
		if err != nil {
			return fmt.Errorf("error creating application code: %s", err)
		}
	}
	return nil
}

func resourceApplicationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)

		opts = applications.AppOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
	)
	resp, err := applications.Create(client, instanceId, opts).Extract()
	if err != nil {
		return diag.Errorf("error creating dedicated application: %s", err)
	}
	d.SetId(resp.Id)

	if v, ok := d.GetOk("app_codes"); ok {
		if err := createApplicationCodes(client, instanceId, d.Id(), v.(*schema.Set).List()); err != nil {
			return diag.FromErr(err)
		}
	}
	return resourceApplicationRead(ctx, d, meta)
}

func queryApplicationCodes(client *golangsdk.ServiceClient, instanceId, appId string) ([]applications.AppCode, error) {
	allPages, err := applications.ListAppCode(client, instanceId, appId, applications.ListCodeOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	return applications.ExtractAppCodes(allPages)
}

func flattenApplicationCodes(codes []applications.AppCode) []interface{} {
	if len(codes) < 1 {
		return nil
	}
	result := make([]interface{}, len(codes))
	for i, v := range codes {
		result[i] = v.Code
	}
	return result
}

func resourceApplicationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)
	resp, err := applications.Get(client, instanceId, appId).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "dedicated application")
	}

	mErr := multierror.Append(nil,
		d.Set("region", cfg.GetRegion(d)),
		d.Set("name", resp.Name),
		d.Set("description", resp.Description),
		d.Set("registration_time", resp.RegistrationTime),
		d.Set("updated_at", resp.UpdateTime),
		d.Set("app_key", resp.AppKey),
		d.Set("app_secret", resp.AppSecret),
	)
	if codes, err := queryApplicationCodes(client, instanceId, appId); err != nil {
		mErr = multierror.Append(mErr, err)
	} else {
		// The application code is sort by create time on server, not code.
		mErr = multierror.Append(d.Set("app_codes",
			schema.NewSet(schema.HashString, flattenApplicationCodes(codes))))
	}

	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving dedicated application fields: %s", err)
	}
	return nil
}

func isCodeInApplication(codes []applications.AppCode, code string) (string, bool) {
	for _, s := range codes {
		if s.Code == code {
			return s.Id, true
		}
	}
	return "", false
}

func removeApplicationCodes(client *golangsdk.ServiceClient, instanceId, appId string, codes []interface{}) error {
	results, err := queryApplicationCodes(client, instanceId, appId)
	if err != nil {
		return fmt.Errorf("error retrieving application codes: %s", err)
	}
	for _, code := range codes {
		codeId, ok := isCodeInApplication(results, code.(string))
		if !ok {
			continue
		}
		if err := applications.RemoveAppCode(client, instanceId, appId, codeId).ExtractErr(); err != nil {
			return fmt.Errorf("error removing code (%v) form the application (%s) : %s", code, appId, err)
		}
	}
	return nil
}

func updateApplicationCodes(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		instanceId       = d.Get("instance_id").(string)
		appId            = d.Id()
		oldRaws, newRaws = d.GetChange("app_codes")

		addRaws    = newRaws.(*schema.Set).Difference(oldRaws.(*schema.Set))
		removeRaws = oldRaws.(*schema.Set).Difference(newRaws.(*schema.Set))
	)
	if removeRaws.Len() > 0 {
		if err := removeApplicationCodes(client, instanceId, appId, removeRaws.List()); err != nil {
			return err
		}
	}

	if addRaws.Len() > 0 {
		if err := createApplicationCodes(client, instanceId, appId, addRaws.List()); err != nil {
			return err
		}
	}
	return nil
}

func resourceApplicationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)

	if d.HasChanges("name", "description") {
		opt := applications.AppOpts{
			Name:        d.Get("name").(string),
			Description: d.Get("description").(string),
		}
		_, err = applications.Update(client, instanceId, appId, opt).Extract()
		if err != nil {
			return diag.Errorf("error updating dedicated application (%s): %s", appId, err)
		}
	}
	if d.HasChange("app_codes") {
		err = updateApplicationCodes(client, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("secret_action") {
		if v, ok := d.GetOk("secret_action"); ok && v.(string) == string(SecretActionReset) {
			if _, err := applications.ResetAppSecret(client, instanceId, appId,
				applications.SecretResetOpts{}).Extract(); err != nil {
				return diag.Errorf("error reseting application secret: %s", err)
			}
		}
	}
	return resourceApplicationRead(ctx, d, meta)
}

func resourceApplicationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.ApigV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating APIG v2 client: %s", err)
	}

	var (
		instanceId = d.Get("instance_id").(string)
		appId      = d.Id()
	)
	err = applications.Delete(client, instanceId, appId).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting application (%s) from the instance (%s)", appId, instanceId))
	}

	return nil
}

func resourceApplicationImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, must be <instance_id>/<id>")
	}
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, d.Set("instance_id", parts[0])
}
