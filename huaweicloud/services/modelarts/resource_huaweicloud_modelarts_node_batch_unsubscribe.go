package modelarts

import (
	"context"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbc"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchUnsubscribeNonUpdatableParams = []string{
	"resource_pool_name",
	"node_ids",
}

// @API BSS POST /v2/orders/suscriptions/resources/query
// @API BSS POST /v2/orders/subscriptions/resources/unsubscribe
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodes
func ResourceV2NodeBatchUnsubscribe() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchUnsubscribeCreate,
		ReadContext:   resourceV2NodeBatchUnsubscribeRead,
		UpdateContext: resourceV2NodeBatchUnsubscribeUpdate,
		DeleteContext: resourceV2NodeBatchUnsubscribeDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchUnsubscribeNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the resource nodes are located.`,
			},
			"resource_pool_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The resource pool name to which the resource nodes belong.`,
			},
			"node_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The ID list of resource nodes to be deleted.`,
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

func waitForV2NodeBatchUnsubscribeCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string,
	unsubscribeNodeIds []interface{}, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			nodes, err := listV2ResourcePoolNodes(client, resourcePoolName)
			// 404: The resource pool does not exist.
			if _, ok := err.(golangsdk.ErrDefault404); ok {
				return nil, "COMPLETED", nil
			}

			if err != nil {
				return nil, "ERROR", err
			}

			resourceIds := utils.PathSearch(`[*].metadata.labels."os.modelarts/resource.id"`, nodes,
				make([]interface{}, 0)).([]interface{})
			if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(resourceIds),
				utils.ExpandToStringList(unsubscribeNodeIds), false, true) {
				return resourceIds, "PENDING", nil
			}
			return resourceIds, "COMPLETED", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceV2NodeBatchUnsubscribeCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		resourcePoolName = d.Get("resource_pool_name").(string)
		deleteNodeIds    = d.Get("node_ids").([]interface{})
	)

	client, err := cfg.NewServiceClient("bssv2", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating BSS client: %s", err)
	}
	// Delete
	err = cbc.UnsubscribePrePaidResources(client, deleteNodeIds)
	if err != nil {
		return diag.FromErr(err)
	}
	err = cbc.WaitForResourcesUnsubscribed(ctx, client, deleteNodeIds, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for all resources to be unsubscribed: %s ", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	client, err = cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}
	err = waitForV2NodeBatchUnsubscribeCompleted(ctx, client, resourcePoolName, deleteNodeIds,
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceV2NodeBatchUnsubscribeRead(ctx, d, meta)
}

func resourceV2NodeBatchUnsubscribeRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUnsubscribeUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchUnsubscribeDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for batch delete the ModelArts nodes. Deleting this
resource will not clear the corresponding request record, but will only remove the resource information from the tfstate
file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
