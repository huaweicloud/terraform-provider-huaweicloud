package servicestage

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

// @API ServiceStage POST /v3/{project_id}/cas/runtimestacks/action
// @API ServiceStage GET /v3/{project_id}/cas/runtimestacks
func ResourceV3RuntimeStackBatchRelease() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceV3RuntimeStackBatchReleaseCreate,
		ReadContext:   resourceV3RuntimeStackBatchReleaseRead,
		UpdateContext: resourceV3RuntimeStackBatchReleaseUpdate,
		DeleteContext: resourceV3RuntimeStackBatchReleaseDelete,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the runtime stacks are located.`,
			},

			// Required parameter(s).
			"runtime_stack_ids": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The runtime stack IDs to be released.`,
			},
		},
	}
}

// The behavior of parameter 'runtime_stack_ids' is 'Required', so the slice length is never less than 1.
func buildV3RuntimeStackReleaseActionParams(runtimeStackIds []interface{}) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(runtimeStackIds))

	for _, runtimeStackId := range runtimeStackIds {
		result = append(result, map[string]interface{}{
			"id": runtimeStackId,
		})
	}

	return result
}

func doV3RuntimeStackBatchReleaseAction(client *golangsdk.ServiceClient, action string, runtimeStackIds []interface{}) error {
	httpUrl := "v3/{project_id}/cas/runtimestacks/action"
	actionPath := client.Endpoint + httpUrl
	actionPath = strings.ReplaceAll(actionPath, "{project_id}", client.ProjectID)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: map[string]interface{}{
			"action": action,
			"parameters": map[string]interface{}{
				"runtimestacks": buildV3RuntimeStackReleaseActionParams(runtimeStackIds),
			},
		},
		OkCodes: []int{
			204,
		},
	}

	_, err := client.Request("POST", actionPath, &opt)
	return err
}

func resourceV3RuntimeStackBatchReleaseCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	err = doV3RuntimeStackBatchReleaseAction(client, "Supported", d.Get("runtime_stack_ids").([]interface{}))
	if err != nil {
		return diag.Errorf("error releasing runtime stacks: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	return resourceV3RuntimeStackBatchReleaseRead(ctx, d, meta)
}

func FilterV3ReleasedRuntimeStacks(client *golangsdk.ServiceClient, runtimeStackIds []interface{}) ([]interface{}, error) {
	runtimeStacks, err := listV3RuntimeStacks(client)
	if err != nil {
		return nil, err
	}

	result := make([]interface{}, 0, len(runtimeStackIds))
	for _, runtimeStackId := range runtimeStackIds {
		result = append(result, utils.PathSearch(fmt.Sprintf("[?id=='%s'&&status=='Supported']", runtimeStackId),
			runtimeStacks, make([]interface{}, 0)).([]interface{})...)
	}

	if len(result) < 1 {
		return result, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v3/{project_id}/cas/runtimestacks",
				RequestId: "NONE",
				Body:      []byte("all runtime stack releases are canceled or deleted"),
			},
		}
	}
	return result, nil
}

func resourceV3RuntimeStackBatchReleaseRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg             = meta.(*config.Config)
		region          = cfg.GetRegion(d)
		runtimeStackIds = d.Get("runtime_stack_ids").([]interface{})
	)

	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	releasedRuntimeStacks, err := FilterV3ReleasedRuntimeStacks(client, runtimeStackIds)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error querying runtime stack releases")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("runtime_stack_ids", utils.PathSearch("[*].id", releasedRuntimeStacks, make([]interface{}, 0))),
	)
	if err = mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error saving resource fields of runtime stack batch release: %s", err)
	}
	return nil
}

func resourceV3RuntimeStackBatchReleaseUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	oldVal, newVal := d.GetChange("runtime_stack_ids")
	if canceled := utils.FindSliceExtraElems(oldVal.([]interface{}), newVal.([]interface{})); len(canceled) > 0 {
		err = doV3RuntimeStackBatchReleaseAction(client, "Disable", canceled)
		if err != nil {
			return diag.Errorf("error canceling the release of runtime stacks: %s", err)
		}
	}

	if released := utils.FindSliceExtraElems(newVal.([]interface{}), oldVal.([]interface{})); len(released) > 0 {
		err = doV3RuntimeStackBatchReleaseAction(client, "Supported", released)
		if err != nil {
			return diag.Errorf("error releasing runtime stacks: %s", err)
		}
	}

	return resourceV3RuntimeStackBatchReleaseRead(ctx, d, meta)
}

func resourceV3RuntimeStackBatchReleaseDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("servicestage", region)
	if err != nil {
		return diag.Errorf("error creating ServiceStage client: %s", err)
	}

	releasedRuntimeStackIds, _ := d.GetChange("runtime_stack_ids")
	err = doV3RuntimeStackBatchReleaseAction(client, "Disable", releasedRuntimeStackIds.([]interface{}))
	if err != nil {
		return diag.Errorf("error canceling the release of runtime stacks: %s", err)
	}
	return nil
}
