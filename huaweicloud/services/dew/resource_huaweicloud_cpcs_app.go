package dew

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

// @API DEW POST /v1/{project_id}/dew/cpcs/apps
// @API DEW GET /v1/{project_id}/dew/cpcs/apps
// @API DEW DELETE /v1/{project_id}/dew/cpcs/apps/{app_id}
func ResourceCpcsApp() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCpcsAppCreate,
		ReadContext:   resourceCpcsAppRead,
		UpdateContext: resourceCpcsAppUpdate,
		DeleteContext: resourceCpcsAppDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceCpcsAppImportState,
		},

		CustomizeDiff: config.FlexibleForceNew([]string{
			"app_name",
			"vpc_id",
			"vpc_name",
			"subnet_id",
			"subnet_name",
			"description",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the application name.`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the VPC to which the application belongs.`,
			},
			"vpc_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the VPC to which the application belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the subnet to which the application belongs.`,
			},
			"subnet_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the name of the subnet to which the application belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the application description.`,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"create_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The creation time of the application, UNIX timestamp in milliseconds.`,
			},
		},
	}
}

func buildCreateAppBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"app_name":    d.Get("app_name"),
		"vpc_id":      d.Get("vpc_id"),
		"vpc_name":    d.Get("vpc_name"),
		"subnet_id":   d.Get("subnet_id"),
		"subnet_name": d.Get("subnet_name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceCpcsAppCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW Client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildCreateAppBodyParams(d)),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CPCS application: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	appId := utils.PathSearch("app_id", respBody, "").(string)
	if appId == "" {
		return diag.Errorf("unable to find the CPCS application ID from the API response")
	}
	d.SetId(appId)

	return resourceCpcsAppRead(ctx, d, meta)
}

// The value of `app_name` can be used to accurately find the target value.
func QueryCpcsAppByAppName(client *golangsdk.ServiceClient, appName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/dew/cpcs/apps"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?app_name=%s", appName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return nil, err
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return nil, err
	}

	appDetail := utils.PathSearch("result|[0]", respBody, nil)
	if appDetail == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return appDetail, nil
}

func resourceCpcsAppRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW Client: %s", err)
	}

	appDetail, err := QueryCpcsAppByAppName(client, d.Get("app_name").(string))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CPCS application")
	}

	// Make sure that the ID value can be written back normally when importing.
	d.SetId(utils.PathSearch("app_id", appDetail, "").(string))
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("app_name", utils.PathSearch("app_name", appDetail, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", appDetail, nil)),
		d.Set("vpc_name", utils.PathSearch("vpc_name", appDetail, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", appDetail, nil)),
		d.Set("subnet_name", utils.PathSearch("subnet_name", appDetail, nil)),
		d.Set("description", utils.PathSearch("description", appDetail, nil)),
		d.Set("create_time", utils.PathSearch("create_time", appDetail, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCpcsAppUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceCpcsAppDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/dew/cpcs/apps/{app_id}"
		product = "kms"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DEW Client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{app_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error deleting CPCS application: %s", err)
	}

	return nil
}

func resourceCpcsAppImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	return []*schema.ResourceData{d}, d.Set("app_name", d.Id())
}
