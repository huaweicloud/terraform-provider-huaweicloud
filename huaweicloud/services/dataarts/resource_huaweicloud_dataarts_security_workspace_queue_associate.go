package dataarts

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

var securityWorkspaceQueueAssociateNonUpdatableParams = []string{
	"workspace_id",
	"source_type",
	"queue_name",
	"connection_id",
	"cluster_id",
	"description",
}

// @API DataArtsStudio POST /v1/{project_id}/security/permission/queue/assigned-source
// @API DataArtsStudio GET /v1/{project_id}/security/permission/queue/assigned-source
// @API DataArtsStudio DELETE /v1/{project_id}/security/permission/queue/assigned-source/{id}
func ResourceSecurityWorkspaceQueueAssociate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSecurityWorkspaceQueueAssociateCreate,
		ReadContext:   resourceSecurityWorkspaceQueueAssociateRead,
		UpdateContext: resourceSecurityWorkspaceQueueAssociateUpdate,
		DeleteContext: resourceSecurityWorkspaceQueueAssociateDelete,

		CustomizeDiff: config.FlexibleForceNew(securityWorkspaceQueueAssociateNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceSecurityWorkspaceQueueAssociateImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the workspace queue associate is located.`,
			},

			// Required parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the workspace to which the queue resource is assigned.`,
			},
			"source_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The queue resource service type. Currently, only MRS and DLI are supported.`,
			},
			"queue_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The queue name.`,
			},
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The data connection ID.`,
			},

			// Optional parameters.
			"cluster_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The ID of MRS cluster.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description of the queue resource assigned to the workspace.`,
			},

			// Attributes
			"queue_attribute": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The attribute of the queue.`,
			},
			"queue_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the queue.`,
			},
			"connection_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The data connection name.`,
			},
			"cluster_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The cluster name.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the queue was added to the workspace.`,
			},
			"create_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who added the queue to the workspace.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The time when the queue resource under the workspace was updated.`,
			},
			"update_user": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user who updated the queue resource under the workspace.`,
			},

			// Internal parameters.
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

func buildSecurityWorkspaceQueueAssociateCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"source_type": d.Get("source_type"),
		"queue_name":  []string{d.Get("queue_name").(string)},
		"conn_id":     d.Get("connection_id"),
		"cluster_id":  utils.ValueIgnoreEmpty(d.Get("cluster_id")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceSecurityWorkspaceQueueAssociateCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg         = meta.(*config.Config)
		region      = cfg.GetRegion(d)
		httpUrl     = "v1/{project_id}/security/permission/queue/assigned-source"
		product     = "dataarts"
		workspaceId = d.Get("workspace_id").(string)
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
		JSONBody:         utils.RemoveNil(buildSecurityWorkspaceQueueAssociateCreateBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Security workspace queue associate: %s", err)
	}

	queueName := d.Get("queue_name").(string)
	d.SetId(fmt.Sprintf("%s/%s", workspaceId, queueName))

	return resourceSecurityWorkspaceQueueAssociateRead(ctx, d, meta)
}

func GetSecurityWorkspaceAssociatedQueueByName(client *golangsdk.ServiceClient, workspaceId, queueName string) (interface{}, error) {
	result, err := listSecurityWorkspaceAssociatedQueues(client, workspaceId)
	if err != nil {
		return nil, err
	}

	queue := utils.PathSearch(fmt.Sprintf("[?queue_name=='%s']|[0]", queueName), result, nil)
	if queue == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/security/permission/queue/assigned-source",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the workspace queue associate with queue name '%s' has been removed", queueName)),
			},
		}
	}

	return queue, nil
}

func resourceSecurityWorkspaceQueueAssociateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dataarts"

		workspaceId = d.Get("workspace_id").(string)
		queueName   = d.Get("queue_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetSecurityWorkspaceAssociatedQueueByName(client, workspaceId, queueName)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving DataArts Security workspace associated queue (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("workspace_id", workspaceId),
		d.Set("source_type", utils.PathSearch("source_type", respBody, nil)),
		d.Set("queue_name", utils.PathSearch("queue_name", respBody, nil)),
		d.Set("connection_id", utils.PathSearch("conn_id", respBody, nil)),
		d.Set("cluster_id", utils.PathSearch("cluster_id", respBody, nil)),
		d.Set("queue_attribute", utils.PathSearch("queue_attr", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("queue_type", utils.PathSearch("queue_type", respBody, nil)),
		d.Set("connection_name", utils.PathSearch("conn_name", respBody, nil)),
		d.Set("cluster_name", utils.PathSearch("cluster_name", respBody, nil)),
		d.Set("create_user", utils.PathSearch("create_user", respBody, nil)),
		d.Set("update_user", utils.PathSearch("update_user", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("create_time", respBody, float64(0)).(float64))/1000, false)),
		d.Set("updated_at", utils.FormatTimeStampRFC3339(
			int64(utils.PathSearch("update_time", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceSecurityWorkspaceQueueAssociateUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceSecurityWorkspaceQueueAssociateDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "dataarts"

		workspaceId = d.Get("workspace_id").(string)
		queueName   = d.Get("queue_name").(string)
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	respBody, err := GetSecurityWorkspaceAssociatedQueueByName(client, workspaceId, queueName)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error retrieving DataArts Security workspace associated queue (%s)", d.Id()))
	}

	queueId := utils.PathSearch("id", respBody, "").(string)
	if queueId == "" {
		return diag.Errorf("unable to find the queue resource ID for delete")
	}

	httpUrl := "v1/{project_id}/security/permission/queue/assigned-source/{id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{id}", queueId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildSecurityMoreHeaders(workspaceId),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err,
			fmt.Sprintf("error deleting DataArts Security workspace queue associate (%s)", d.Id()))
	}

	return nil
}

func resourceSecurityWorkspaceQueueAssociateImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import ID, want '<workspace_id>/<queue_name>', got '%s'", d.Id())
	}

	return []*schema.ResourceData{d}, multierror.Append(nil,
		d.Set("workspace_id", parts[0]),
		d.Set("queue_name", parts[1]),
	).ErrorOrNil()
}
