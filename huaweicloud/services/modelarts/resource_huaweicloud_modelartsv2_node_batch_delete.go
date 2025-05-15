package modelarts

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var v2NodeBatchDeleteNonUpdatableParams = []string{
	"resource_pool_name",
	"node_names",
}

// @API ModelArts POST /v2/{project_id}/pools/{pool_name}/nodes/batch-delete
// @API ModelArts GET /v2/{project_id}/pools/{pool_name}/nodes
func ResourceV2NodeBatchDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV2NodeBatchDeleteCreate,
		ReadContext:   resourceV2NodeBatchDeleteRead,
		UpdateContext: resourceV2NodeBatchDeleteUpdate,
		DeleteContext: resourceV2NodeBatchDeleteDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(45 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(v2NodeBatchDeleteNonUpdatableParams),

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
			"node_names": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The name list of resource nodes to be deleted.`,
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

func waitForV2NodeBatchDeleteCompleted(ctx context.Context, client *golangsdk.ServiceClient, resourcePoolName string,
	deleteNodeNames []interface{}, timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"PENDING"},
		Target:  []string{"COMPLETED"},
		Refresh: func() (interface{}, string, error) {
			nodes, err := listV2ResourcePoolNodes(client, resourcePoolName)
			if err != nil {
				return nil, "ERROR", err
			}

			names := utils.PathSearch(`[*].metadata.name`, nodes, make([]interface{}, 0)).([]interface{})
			if utils.IsSliceContainsAnyAnotherSliceElement(utils.ExpandToStringList(names),
				utils.ExpandToStringList(deleteNodeNames), false, true) {
				return names, "PENDING", nil
			}
			return names, "COMPLETED", nil
		},
		Timeout:      timeout,
		Delay:        30 * time.Second,
		PollInterval: 30 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceV2NodeBatchDeleteCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg              = meta.(*config.Config)
		region           = cfg.GetRegion(d)
		httpUrl          = "v2/{project_id}/pools/{pool_name}/nodes/batch-delete"
		resourcePoolName = d.Get("resource_pool_name").(string)
		deleteNodeNames  = d.Get("node_names").([]interface{})
	)
	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	batchDeletePath := client.Endpoint + httpUrl
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{project_id}", client.ProjectID)
	batchDeletePath = strings.ReplaceAll(batchDeletePath, "{pool_name}", resourcePoolName)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"deleteNodeNames": deleteNodeNames,
		},
	}

	_, err = client.Request("POST", batchDeletePath, &opt)
	if err != nil {
		return diag.Errorf("error executing batch delete operation: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	err = waitForV2NodeBatchDeleteCompleted(ctx, client, resourcePoolName, deleteNodeNames, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceV2NodeBatchDeleteRead(ctx, d, meta)
}

func resourceV2NodeBatchDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceV2NodeBatchDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
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
