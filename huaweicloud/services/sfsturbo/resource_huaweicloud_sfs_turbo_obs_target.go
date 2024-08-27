package sfsturbo

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SFSTurbo POST /v1/{project_id}/sfs-turbo/shares/{share_id}/targets
// @API SFSTurbo GET /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}
// @API SFSTurbo DELETE /v1/{project_id}/sfs-turbo/shares/{share_id}/targets/{target_id}
func ResourceOBSTarget() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOBSTargetCreate,
		ReadContext:   resourceOBSTargetRead,
		DeleteContext: resourceOBSTargetDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceOBSTargetImportState,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"share_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"file_system_path": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"obs": {
				Type:     schema.TypeList,
				Required: true,
				ForceNew: true,
				MaxItems: 1,
				Elem:     obsSchema(),
			},
			"delete_data_in_file_system": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func obsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"bucket": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
	return &sc
}

func buildCreateOBSTargetBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"file_system_path": d.Get("file_system_path"),
		"obs":              buildOBSBody(d.Get("obs.0").(map[string]interface{})),
	}
	return bodyParams
}

func buildOBSBody(obsData map[string]interface{}) map[string]interface{} {
	return map[string]interface{}{
		"bucket":   obsData["bucket"],
		"endpoint": obsData["endpoint"],
	}
}

func resourceOBSTargetCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 client: %s", err)
	}

	createObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets"
	createObsTargetPath := client.ResourceBaseURL() + createObsTargetHttpUrl
	createObsTargetPath = strings.ReplaceAll(createObsTargetPath, "{share_id}", d.Get("share_id").(string))

	createObsTargetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	createObsTargetOpt.JSONBody = utils.RemoveNil(buildCreateOBSTargetBodyParams(d))
	createObsTargetResp, err := client.Request("POST", createObsTargetPath, &createObsTargetOpt)
	if err != nil {
		return diag.Errorf("error creating OBS target to the SFS Turbo: %s", err)
	}

	createObsTargetRespBody, err := utils.FlattenResponse(createObsTargetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("target_id", createObsTargetRespBody)
	if err != nil {
		return diag.Errorf("error creating OBS target to the SFS Turbo: ID is not found in API response")
	}

	d.SetId(id.(string))

	err = obsTargetWaitingForStateCompleted(ctx, d, meta, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.Errorf("error waiting for the creation of OBS target (%s) to complete: %s", d.Id(), err)
	}

	return resourceOBSTargetRead(ctx, d, meta)
}

func getOBSTargetInfo(d *schema.ResourceData, meta interface{}) (*http.Response, error) {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return nil, fmt.Errorf("error creating SFS v1 client: %s", err)
	}

	getObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets/{target_id}"
	getObsTargetPath := client.ResourceBaseURL() + getObsTargetHttpUrl
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{share_id}", d.Get("share_id").(string))
	getObsTargetPath = strings.ReplaceAll(getObsTargetPath, "{target_id}", d.Id())
	getObsTargetOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	resp, err := client.Request("GET", getObsTargetPath, &getObsTargetOpts)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func resourceOBSTargetRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	getObsTargetResp, err := getOBSTargetInfo(d, meta)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "SFS Turbo OBS target")
	}

	getObsTargetRespBody, err := utils.FlattenResponse(getObsTargetResp)
	if err != nil {
		return diag.FromErr(err)
	}

	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("file_system_path", utils.PathSearch("file_system_path", getObsTargetRespBody, nil)),
		d.Set("obs", flattenGetOBSDataResponseBody(utils.PathSearch("obs", getObsTargetRespBody, nil))),
		d.Set("status", utils.PathSearch("lifecycle", getObsTargetRespBody, nil)),
		d.Set("created_at", utils.PathSearch("creation_time", getObsTargetRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetOBSDataResponseBody(resp interface{}) []map[string]interface{} {
	if resp == nil {
		return nil
	}

	return []map[string]interface{}{
		{
			"bucket":   utils.PathSearch("bucket", resp, nil),
			"endpoint": utils.PathSearch("endpoint", resp, nil),
		},
	}
}

func resourceOBSTargetDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.SfsV1Client(region)
	if err != nil {
		return diag.Errorf("error creating SFS v1 Client: %s", err)
	}

	deleteObsTargetHttpUrl := "sfs-turbo/shares/{share_id}/targets/{target_id}"
	deleteObsTargetPath := client.ResourceBaseURL() + deleteObsTargetHttpUrl
	deleteObsTargetPath = strings.ReplaceAll(deleteObsTargetPath, "{share_id}", d.Get("share_id").(string))
	deleteObsTargetPath = strings.ReplaceAll(deleteObsTargetPath, "{target_id}", d.Id())

	if v, ok := d.GetOk("delete_data_in_file_system"); ok {
		deleteObsTargetPath += fmt.Sprintf("?delete_data_in_file_system=%v", v)
	}

	deleteObsTargetOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	_, err = client.Request("DELETE", deleteObsTargetPath, &deleteObsTargetOpt)
	if err != nil {
		return diag.Errorf("error deleting OBS target from SFS Turbo: %s", err)
	}

	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"DELETED"},
		Refresh:      obsTargetStatusRefreshFunc(d, meta, true),
		Timeout:      d.Timeout(schema.TimeoutDelete),
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err = stateConf.WaitForStateContext(ctx)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func obsTargetWaitingForStateCompleted(ctx context.Context, d *schema.ResourceData, meta interface{}, t time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      obsTargetStatusRefreshFunc(d, meta, false),
		Timeout:      t,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func obsTargetStatusRefreshFunc(d *schema.ResourceData, meta interface{}, isDelete bool) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		resp, err := getOBSTargetInfo(d, meta)
		if err != nil {
			if _, ok := err.(golangsdk.ErrDefault404); ok && isDelete {
				return "Resource Not Found", "DELETED", nil
			}

			return nil, "ERROR", err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, "ERROR", err
		}

		state, err := jmespath.Search("lifecycle", respBody)
		if err != nil {
			return nil, "ERROR", fmt.Errorf("error parsing %s from response body", state)
		}

		statusRaw := fmt.Sprintf("%v", state)

		if utils.StrSliceContains([]string{"MISCONFIGURED", "FAILED"}, statusRaw) {
			return respBody, "ERROR", fmt.Errorf("unexpected status: '%s'", statusRaw)
		}

		if utils.StrSliceContains([]string{"AVAILABLE"}, statusRaw) {
			return respBody, "COMPLETED", nil
		}

		return respBody, "PENDING", nil
	}
}

func resourceOBSTargetImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData,
	error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format for import ID, want '<share_id>/<id>', but got '%s'", d.Id())
	}

	d.SetId(parts[1])
	return []*schema.ResourceData{d}, d.Set("share_id", parts[0])
}
