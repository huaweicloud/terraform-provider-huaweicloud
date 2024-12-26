package live

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

// @API LIVE GET /v1/{project_id}/domain/delay
// @API LIVE PUT /v1/{project_id}/domain/delay
func ResourceStreamDelay() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceStreamDelayCreate,
		ReadContext:   resourceStreamDelayRead,
		UpdateContext: resourceStreamDelayUpdate,
		DeleteContext: resourceStreamDelayDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceStreamDelayImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				ForceNew:    true,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to create the resource. If omitted, the provider-level region will be used.`,
			},
			"domain_name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the streaming domain name.`,
			},
			"delay": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the delay time, in ms.`,
			},
			"app_name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Computed:    true,
				Description: `Specifies the application name.`,
			},
		},
	}
}

func buildStreamDelayBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"play_domain": d.Get("domain_name"),
		"delay":       d.Get("delay"),
		"app":         utils.ValueIgnoreEmpty(d.Get("app_name")),
	}
}

func updateStreamDelay(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	requestPath := client.Endpoint + "v1/{project_id}/domain/delay"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildStreamDelayBodyParams(d)),
	}

	_, err := client.Request("PUT", requestPath, &requestOpt)
	return err
}

func resourceStreamDelayCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateStreamDelay(client, d); err != nil {
		return diag.Errorf("error configuring Live stream delay time in creation operation: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)

	return resourceStreamDelayRead(ctx, d, meta)
}

func ReadStreamDelay(client *golangsdk.ServiceClient, domainName string) (interface{}, error) {
	requestPath := client.Endpoint + "v1/{project_id}/domain/delay"
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += fmt.Sprintf("?play_domain=%s", domainName)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		// When the `domain_name` does not exist, calling the query API will return a `400` status code.
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", domainNameNotExistsCode)
	}

	return utils.FlattenResponse(resp)
}

func resourceStreamDelayRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		domainName = d.Get("domain_name").(string)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	respBody, err := ReadStreamDelay(client, domainName)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Live stream delay time")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("domain_name", utils.PathSearch("play_domain", respBody, nil)),
		d.Set("app_name", utils.PathSearch("delay_config|[0].app", respBody, nil)),
		d.Set("delay", utils.PathSearch("delay_config|[0].delay", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceStreamDelayUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("live", region)
	if err != nil {
		return diag.Errorf("error creating Live client: %s", err)
	}

	if err := updateStreamDelay(client, d); err != nil {
		return diag.Errorf("error configuring Live stream delay time in update operation: %s", err)
	}

	return resourceStreamDelayRead(ctx, d, meta)
}

// This resource always has value and cannot be deleted.
func resourceStreamDelayDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceStreamDelayImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	importedId := d.Id()
	if importedId == "" {
		return nil, fmt.Errorf("invalid format specified for import ID, `domain_name` is empty")
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return nil, fmt.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(resourceId)
	return []*schema.ResourceData{d}, d.Set("domain_name", importedId)
}
