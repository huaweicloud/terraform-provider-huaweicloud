package apig

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API APIG PUT /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl
// @API APIG GET /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl
// @API APIG DELETE /v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl
func ResourceApplicationAcl() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceApplicationAclCreate,
		ReadContext:   resourceApplicationAclRead,
		UpdateContext: resourceApplicationAclUpdate,
		DeleteContext: resourceApplicationAclDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceApplicationAclImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "The region where the application and ACL rules are located.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the dedicated instance to which the application belongs.",
			},
			"application_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the application to which the ACL rules belong.",
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ACL type.",
			},
			"values": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "The ACL values.",
			},
		},
	}
}

func buildCreateACLBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"app_acl_type":   d.Get("type"),
		"app_acl_values": d.Get("values").([]interface{}),
	}
}

func settingACLsToApplication(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl"
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
	)

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateACLBodyParams(d)),
	}

	_, err := client.Request("PUT", createPath, &opt)
	if err != nil {
		return fmt.Errorf("error setting the ACL rules to the application (%s): %s", appId, err)
	}
	return nil
}

func resourceApplicationAclCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		appId  = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = settingACLsToApplication(client, d)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(appId)

	return resourceApplicationAclRead(ctx, d, meta)
}

func resourceApplicationAclRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl"
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)
	getPath = strings.ReplaceAll(getPath, "{app_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving the ACL rules from the application (%s)", appId))
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if respBody == nil || (reflect.TypeOf(respBody).Kind() == reflect.Map && len(respBody.(map[string]interface{})) == 0) {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{},
			fmt.Sprintf("no ACL rule found in the application (%s)", appId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("type", utils.PathSearch("app_acl_type", respBody, nil)),
		d.Set("values", utils.PathSearch("app_acl_values", respBody, nil)),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving the fields of the application ACL rules: %s", err)
	}
	return nil
}

func resourceApplicationAclUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	err = settingACLsToApplication(client, d)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceApplicationAclRead(ctx, d, meta)
}

func resourceApplicationAclDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		httpUrl    = "v2/{project_id}/apigw/instances/{instance_id}/apps/{app_id}/app-acl"
		instanceId = d.Get("instance_id").(string)
		appId      = d.Get("application_id").(string)
	)

	client, err := cfg.NewServiceClient("apig", region)
	if err != nil {
		return diag.Errorf("error creating APIG client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)
	createPath = strings.ReplaceAll(createPath, "{app_id}", appId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &opt)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting the ACL rules from the application (%s)", appId))
	}
	return nil
}

func resourceApplicationAclImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importedId := d.Id()
	parts := strings.Split(importedId, "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<instance_id>/<id>', but got '%s'", importedId)
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("instance_id", parts[0]),
		d.Set("application_id", parts[1]),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return []*schema.ResourceData{d},
			fmt.Errorf("error saving the fields of the application ACL rules during import: %s", err)
	}
	return []*schema.ResourceData{d}, nil
}
