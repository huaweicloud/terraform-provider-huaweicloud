package cts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CTS PUT /v3/{project_id}/configuration
// @API CTS GET /v3/{project_id}/configuration
func ResourceConfiguration() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConfigurationCreate,
		ReadContext:   resourceConfigurationRead,
		UpdateContext: resourceConfigurationUpdate,
		DeleteContext: resourceConfigurationDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region in which to manage the CTS configuration.`,
			},

			// Required parameters.
			"is_sync_global_trace": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to synchronize global service logs from the central region.`,
			},
			"is_support_read_only": {
				Type:        schema.TypeBool,
				Required:    true,
				Description: `Whether to enable the reporting of read-only audit logs for all cloud services.`,
			},

			// Optional parameters.
			"support_read_only_services": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `Cloud services that enable read-only audit logs.`,
			},
		},
	}
}

func buildConfigurationCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"is_sync_global_trace":       d.Get("is_sync_global_trace"),
		"is_support_read_only":       d.Get("is_support_read_only"),
		"support_read_only_services": utils.ValueIgnoreEmpty(utils.ExpandToStringList(d.Get("support_read_only_services").([]interface{}))),
	}
}

func updateConfiguration(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v3/{project_id}/configuration"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildConfigurationCreateBodyParams(d)),
		OkCodes:          []int{204},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return fmt.Errorf("error updating CTS configuration: %s", err)
	}
	return nil
}

func resourceConfigurationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	err = updateConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error creating CTS configuration: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceConfigurationRead(ctx, d, meta)
}

func GetConfiguration(client *golangsdk.ServiceClient) (interface{}, error) {
	httpUrl := "v3/{project_id}/configuration"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceConfigurationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	configuration, err := GetConfiguration(client)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CTS configuration")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("is_sync_global_trace", utils.PathSearch("is_sync_global_trace", configuration, nil)),
		d.Set("is_support_read_only", utils.PathSearch("is_support_read_only", configuration, nil)),
		d.Set("support_read_only_services", utils.PathSearch("support_read_only_services", configuration, make([]interface{}, 0))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceConfigurationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cts", region)
	if err != nil {
		return diag.Errorf("error creating CTS client: %s", err)
	}

	err = updateConfiguration(client, d)
	if err != nil {
		return diag.Errorf("error updating CTS configuration: %s", err)
	}

	return resourceConfigurationRead(ctx, d, meta)
}

func resourceConfigurationDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a configuration resource for managing CTS configuration. Deleting this resource
will not restore the CTS configuration to the default value, but will only remove the resource information from the
tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
