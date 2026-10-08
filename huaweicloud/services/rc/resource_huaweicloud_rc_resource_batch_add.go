package rc

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var resourceBatchAddNonUpdatableParams = []string{
	"group_id",
	"resource_ids",
}

// @API RC PUT /v1/resource-center/groups/{group_id}/resources
func ResourceResourceBatchAdd() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceResourceBatchAddCreate,
		ReadContext:   resourceResourceBatchAddRead,
		UpdateContext: resourceResourceBatchAddUpdate,
		DeleteContext: resourceResourceBatchAddDelete,

		CustomizeDiff: config.FlexibleForceNew(resourceBatchAddNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resources to be added are located.`,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource group.`,
			},
			"resource_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Description: `The list of resource IDs to be added.`,
				Elem:        &schema.Schema{Type: schema.TypeString},
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

func buildResourceBatchAddBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"resource_ids": utils.ExpandToStringList(d.Get("resource_ids").([]interface{})),
	}

	return bodyParams
}

func resourceResourceBatchAddCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Get("group_id").(string)
		httpUrl = "v1/resource-center/groups/{group_id}/resources"
	)

	client, err := cfg.NewServiceClient("rc", region)
	if err != nil {
		return diag.Errorf("error creating RC client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{group_id}", groupId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildResourceBatchAddBodyParams(d),
	}

	_, err = client.Request("PUT", createPath, &opt)
	if err != nil {
		return diag.Errorf("error batch adding resources to resource group: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	return nil
}

func resourceResourceBatchAddRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceBatchAddUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceResourceBatchAddDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch add resources. Deleting this resource will
not clear the corresponding request record, but will only remove the resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
