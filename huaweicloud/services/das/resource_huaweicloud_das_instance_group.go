package das

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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

var instanceGroupNonUpdatableParams = []string{
	"datastore_type",
}

// @API DAS POST /v3/{project_id}/batch-inspection/instance-group
// @API DAS GET /v3/{project_id}/batch-inspection/instance-group
// @API DAS PUT /v3/{project_id}/batch-inspection/instance-group
// @API DAS DELETE /v3/{project_id}/batch-inspection/instance-group
func ResourceInstanceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceInstanceGroupCreate,
		ReadContext:   resourceInstanceGroupRead,
		UpdateContext: resourceInstanceGroupUpdate,
		DeleteContext: resourceInstanceGroupDelete,

		CustomizeDiff: config.FlexibleForceNew(instanceGroupNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceInstanceGroupImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the instance group is located.`,
			},

			// Required parameters.
			"datastore_type": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The database type.`,
			},
			"group_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The instance group name.`,
			},

			// Required parameters.
			"description": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The description of the instance group.`,
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

func buildInstanceGroupCreateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"datastore_type": d.Get("datastore_type"),
		"group_name":     d.Get("group_name"),
		"description":    d.Get("description"),
	}
}

func resourceInstanceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/batch-inspection/instance-group"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildInstanceGroupCreateBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DAS instance group: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error parsing response: %s", err)
	}

	groupId := utils.PathSearch("group_id", respBody, float64(0)).(float64)
	if groupId == 0 {
		return diag.Errorf("error creating DAS instance group: %s", err)
	}
	d.SetId(strconv.Itoa(int(groupId)))

	return resourceInstanceGroupRead(ctx, d, meta)
}

// GetInstanceGroupById queries the instance group detail by datastore type and group ID.
func GetInstanceGroupById(client *golangsdk.ServiceClient, datastoreType, groupId string) (interface{}, error) {
	queryParams := fmt.Sprintf("&datastore_type=%s", datastoreType)
	groups, err := ListInstanceGroups(client, queryParams)
	if err != nil {
		return nil, err
	}

	for _, group := range groups {
		currentGroupId := fmt.Sprintf("%v", utils.PathSearch("group_id", group, ""))
		if currentGroupId == groupId {
			return group, nil
		}
	}

	return nil, golangsdk.ErrDefault404{
		ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
			Method:    "GET",
			URL:       "/v3/{project_id}/batch-inspection/instance-group",
			RequestId: "NONE",
			Body:      []byte(fmt.Sprintf("the instance group (%s) has been removed", groupId)),
		},
	}
}

func resourceInstanceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		datastoreType = d.Get("datastore_type").(string)
		groupId       = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	respBody, err := GetInstanceGroupById(client, datastoreType, groupId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error retrieving DAS instance group (%s)", groupId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("datastore_type", datastoreType),
		d.Set("group_name", utils.PathSearch("group_name", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildInstanceGroupUpdateBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		// Required parameter
		"group_id": d.Id(),

		// Optional parameters
		"group_name":  utils.ValueIgnoreEmpty(d.Get("group_name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceInstanceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/batch-inspection/instance-group"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildInstanceGroupUpdateBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating DAS instance group: %s", err)
	}

	return resourceInstanceGroupRead(ctx, d, meta)
}

func resourceInstanceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		groupId = d.Id()
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	httpUrl := "v3/{project_id}/batch-inspection/instance-group"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"group_id": groupId,
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error deleting DAS instance group (%s)", groupId))
	}

	return nil
}

func resourceInstanceGroupImportState(_ context.Context, d *schema.ResourceData,
	_ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <datastore_type>/<group_id>")
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("datastore_type", parts[0])
}
