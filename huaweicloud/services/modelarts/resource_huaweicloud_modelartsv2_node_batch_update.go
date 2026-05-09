package modelarts

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchUpdateNonUpdatableParams = []string{
	"pool_id",
	"node_names",
	"action",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-update
func ResourceV2NodeBatchUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchUpdateCreate,
		ReadContext:   resourceV2NodeBatchUpdateRead,
		UpdateContext: resourceV2NodeBatchUpdateUpdate,
		DeleteContext: resourceV2NodeBatchUpdateDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(30 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource nodes are located.`,
			},

			// Required parameters.
			"pool_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the resource pool to which the resource nodes belong.`,
			},
			"node_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of resource nodes to be updated.`,
			},
			"action": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The type of the node update action.`,
			},

			// Optional parameters.
			"ha_redundant_effect": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The effect of the high availability redundancy.`,
			},
			"driver": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The driver version and status information of the node.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"version": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The version of the driver on the node.`,
						},
						"update_strategy": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: `The driver upgrade strategy of the node.`,
						},
					},
				},
			},
			"tags": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: `The list of resource tags to be operated in batch.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The key of the tag.`,
						},
						"value": {
							Type:        schema.TypeString,
							Required:    true,
							Description: `The value of the tag.`,
						},
					},
				},
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

func buildNodeBatchUpdateDriverParams(drivers []interface{}) map[string]interface{} {
	if len(drivers) < 1 {
		return nil
	}

	driver := drivers[0]
	return map[string]interface{}{
		"version":        utils.ValueIgnoreEmpty(utils.PathSearch("version", driver, nil)),
		"updateStrategy": utils.ValueIgnoreEmpty(utils.PathSearch("update_strategy", driver, nil)),
	}
}

func buildNodeBatchUpdateTagsParams(tags []interface{}) []map[string]interface{} {
	if len(tags) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(tags))
	for _, tag := range tags {
		result = append(result, map[string]interface{}{
			"key":   utils.PathSearch("key", tag, nil),
			"value": utils.PathSearch("value", tag, nil),
		})
	}
	return result
}

func buildNodeBatchUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameters.
		"nodeNames": d.Get("node_names"),
		"action":    d.Get("action"),

		// Optional parameters.
		"haRedundantEffect": utils.ValueIgnoreEmpty(d.Get("ha_redundant_effect")),
		"driver":            buildNodeBatchUpdateDriverParams(d.Get("driver").([]interface{})),
		"tags":              buildNodeBatchUpdateTagsParams(d.Get("tags").([]interface{})),
	}
}

func resourceV2NodeBatchUpdateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/pools/{pool_name}/nodes/batch-update"
		poolId  = d.Get("pool_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{pool_name}", poolId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: buildNodeBatchUpdateBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &opt)
	if err != nil {
		return diag.Errorf("error executing batch update operation: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceV2NodeBatchUpdateRead(ctx, d, meta)
}

func resourceV2NodeBatchUpdateRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUpdateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch update the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
