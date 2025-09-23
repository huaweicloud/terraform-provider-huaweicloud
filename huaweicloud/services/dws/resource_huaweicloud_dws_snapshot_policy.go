// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DWS
// ---------------------------------------------------------------

package dws

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DWS PUT /v2/{project_id}/clusters/{cluster_id}/snapshot-policies
// @API DWS GET /v2/{project_id}/clusters/{cluster_id}/snapshot-policies
// @API DWS DELETE /v1.0/{project_id}/clusters/{cluster_id}/snapshot-policies/{id}
func ResourceDwsSnapshotPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDwsSnapshotPolicyCreate,
		ReadContext:   resourceDwsSnapshotPolicyRead,
		DeleteContext: resourceDwsSnapshotPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceDwsSnapshotPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The cluster ID of which the automated snapshot policy belongs to.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The name of the automated snapshot policy.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The type of the automated snapshot policy.`,
			},
			"strategy": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The strategy of the automated snapshot policy.`,
			},
		},
	}
}

func resourceDwsSnapshotPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createDwsSnapshotPolicy: create a DWS snapshot policy.
	var (
		createDwsSnapshotPolicyHttpUrl = "v2/{project_id}/clusters/{cluster_id}/snapshot-policies"
		createDwsSnapshotPolicyProduct = "dws"
	)
	createDwsSnapshotPolicyClient, err := cfg.NewServiceClient(createDwsSnapshotPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	createDwsSnapshotPolicyPath := createDwsSnapshotPolicyClient.Endpoint + createDwsSnapshotPolicyHttpUrl
	createDwsSnapshotPolicyPath = strings.ReplaceAll(createDwsSnapshotPolicyPath, "{project_id}", createDwsSnapshotPolicyClient.ProjectID)
	createDwsSnapshotPolicyPath = strings.ReplaceAll(createDwsSnapshotPolicyPath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	createDwsSnapshotPolicyOpt := golangsdk.RequestOpts{
		MoreHeaders:      requestOpts.MoreHeaders,
		KeepResponseBody: true,
	}
	createDwsSnapshotPolicyOpt.JSONBody = utils.RemoveNil(buildCreateDwsSnapshotPolicyBodyParams(d))
	createDwsSnapshotPolicyResp, err := createDwsSnapshotPolicyClient.Request("PUT", createDwsSnapshotPolicyPath, &createDwsSnapshotPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating DWS snapshot policy: %s", err)
	}

	_, err = utils.FlattenResponse(createDwsSnapshotPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	// getDwsSnapshotPolicy: read ID after creation.
	var (
		getDwsSnapshotPolicyHttpUrl = "v2/{project_id}/clusters/{cluster_id}/snapshot-policies"
		getDwsSnapshotPolicyProduct = "dws"
	)
	getDwsSnapshotPolicyClient, err := cfg.NewServiceClient(getDwsSnapshotPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsSnapshotPolicyPath := getDwsSnapshotPolicyClient.Endpoint + getDwsSnapshotPolicyHttpUrl
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{project_id}", getDwsSnapshotPolicyClient.ProjectID)
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	getDwsSnapshotPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	getDwsSnapshotPolicyResp, err := getDwsSnapshotPolicyClient.Request("GET", getDwsSnapshotPolicyPath, &getDwsSnapshotPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating DWS snapshot policy: %s", err)
	}

	getDwsSnapshotPolicyRespBody, err := utils.FlattenResponse(getDwsSnapshotPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("backup_strategies[?policy_name=='%s']|[0].policy_id", d.Get("name"))
	id := utils.PathSearch(jsonPath, getDwsSnapshotPolicyRespBody, nil)
	if id == nil {
		return diag.Errorf("error creating DWS snapshot policy: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceDwsSnapshotPolicyRead(ctx, d, meta)
}

func buildCreateDwsSnapshotPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"backup_strategies": []map[string]interface{}{
			{
				"policy_name":     utils.ValueIgnoreEmpty(d.Get("name")),
				"cluster_id":      utils.ValueIgnoreEmpty(d.Get("cluster_id")),
				"backup_type":     utils.ValueIgnoreEmpty(d.Get("type")),
				"backup_strategy": utils.ValueIgnoreEmpty(d.Get("strategy")),
				"backup_level":    "cluster",
			},
		},
	}
	return bodyParams
}

func resourceDwsSnapshotPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getDwsSnapshotPolicy: Query the DWS snapshot policy.
	var (
		getDwsSnapshotPolicyHttpUrl = "v2/{project_id}/clusters/{cluster_id}/snapshot-policies"
		getDwsSnapshotPolicyProduct = "dws"
	)
	getDwsSnapshotPolicyClient, err := cfg.NewServiceClient(getDwsSnapshotPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	getDwsSnapshotPolicyPath := getDwsSnapshotPolicyClient.Endpoint + getDwsSnapshotPolicyHttpUrl
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{project_id}", getDwsSnapshotPolicyClient.ProjectID)
	getDwsSnapshotPolicyPath = strings.ReplaceAll(getDwsSnapshotPolicyPath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))

	getDwsSnapshotPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	getDwsSnapshotPolicyResp, err := getDwsSnapshotPolicyClient.Request("GET", getDwsSnapshotPolicyPath, &getDwsSnapshotPolicyOpt)

	if err != nil {
		// The cluster ID does not exist.
		// "DWS.0001": The cluster ID is a non-standard UUID, the status code is 400.
		// "DWS.0047": The cluster ID is a standard UUID, the status code is 404.
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code", ClusterIdIllegalErrCode),
			"error retrieving DWS snapshot policy")
	}

	getDwsSnapshotPolicyRespBody, err := utils.FlattenResponse(getDwsSnapshotPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("backup_strategies[?policy_id=='%s']|[0]", d.Id())
	rawData := utils.PathSearch(jsonPath, getDwsSnapshotPolicyRespBody, nil)
	if rawData == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("policy_name", rawData, nil)),
		d.Set("type", utils.PathSearch("backup_type", rawData, nil)),
		d.Set("strategy", utils.PathSearch("backup_strategy", rawData, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDwsSnapshotPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteDwsSnapshotPolicy: delete DWS snapshot policy
	var (
		deleteDwsSnapshotPolicyHttpUrl = "v1.0/{project_id}/clusters/{cluster_id}/snapshot-policies/{id}"
		deleteDwsSnapshotPolicyProduct = "dws"
	)
	deleteDwsSnapshotPolicyClient, err := cfg.NewServiceClient(deleteDwsSnapshotPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DWS Client: %s", err)
	}

	deleteDwsSnapshotPolicyPath := deleteDwsSnapshotPolicyClient.Endpoint + deleteDwsSnapshotPolicyHttpUrl
	deleteDwsSnapshotPolicyPath = strings.ReplaceAll(deleteDwsSnapshotPolicyPath, "{project_id}", deleteDwsSnapshotPolicyClient.ProjectID)
	deleteDwsSnapshotPolicyPath = strings.ReplaceAll(deleteDwsSnapshotPolicyPath, "{cluster_id}", fmt.Sprintf("%v", d.Get("cluster_id")))
	deleteDwsSnapshotPolicyPath = strings.ReplaceAll(deleteDwsSnapshotPolicyPath, "{id}", d.Id())

	deleteDwsSnapshotPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      requestOpts.MoreHeaders,
	}
	_, err = deleteDwsSnapshotPolicyClient.Request("DELETE", deleteDwsSnapshotPolicyPath, &deleteDwsSnapshotPolicyOpt)
	if err != nil {
		return diag.Errorf("error deleting DWS SnapshotPolicy: %s", err)
	}

	return nil
}

func resourceDwsSnapshotPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <cluster_id>/<id>")
	}

	d.Set("cluster_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
