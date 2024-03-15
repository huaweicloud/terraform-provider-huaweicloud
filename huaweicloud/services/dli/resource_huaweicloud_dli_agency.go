// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product DLI
// ---------------------------------------------------------------

package dli

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

const DliAgencyID = "dli_admin_agency"

// @API DLI POST /v2/{project_id}/agency
// @API DLI GET /v2/{project_id}/agency
func ResourceDliAgency() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDliAgencyCreateOrUpdate,
		UpdateContext: resourceDliAgencyCreateOrUpdate,
		ReadContext:   resourceDliAgencyRead,
		DeleteContext: resourceDliAgencyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"roles": {
				Type:        schema.TypeSet,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Required:    true,
				Description: `The list of roles.`,
			},
			"version": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Agency version information.`,
			},
		},
	}
}

func resourceDliAgencyCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAgency: create a Agency.
	var (
		createAgencyHttpUrl = "v2/{project_id}/agency"
		createAgencyProduct = "dli"
	)
	createAgencyClient, err := cfg.NewServiceClient(createAgencyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	createAgencyPath := createAgencyClient.Endpoint + createAgencyHttpUrl
	createAgencyPath = strings.ReplaceAll(createAgencyPath, "{project_id}", createAgencyClient.ProjectID)

	createAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAgencyOpt.JSONBody = utils.RemoveNil(buildDliAgencyBodyParams(d))
	requestResp, err := createAgencyClient.Request("POST", createAgencyPath, &createAgencyOpt)
	if err != nil {
		return diag.Errorf("error creating DLI Agency: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to modify the agency: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}

	d.SetId(DliAgencyID)

	return resourceDliAgencyRead(ctx, d, meta)
}

func buildDliAgencyBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"roles": utils.ExpandToStringListBySet(d.Get("roles").(*schema.Set)),
	}
	return bodyParams
}

func resourceDliAgencyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAgency: Query the Agency.
	var (
		getAgencyHttpUrl = "v2/{project_id}/agency"
		getAgencyProduct = "dli"
	)
	getAgencyClient, err := cfg.NewServiceClient(getAgencyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	getAgencyPath := getAgencyClient.Endpoint + getAgencyHttpUrl
	getAgencyPath = strings.ReplaceAll(getAgencyPath, "{project_id}", getAgencyClient.ProjectID)

	getAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAgencyResp, err := getAgencyClient.Request("GET", getAgencyPath, &getAgencyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DLI Agency")
	}

	getAgencyRespBody, err := utils.FlattenResponse(getAgencyResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", getAgencyRespBody, true).(bool) {
		return diag.Errorf("unable to query the agency: %s",
			utils.PathSearch("message", getAgencyRespBody, "Message Not Found"))
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("roles", utils.PathSearch("current_roles", getAgencyRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getAgencyRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDliAgencyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAgency: delete Agency
	var (
		deleteAgencyHttpUrl = "v2/{project_id}/agency"
		deleteAgencyProduct = "dli"
	)
	deleteAgencyClient, err := cfg.NewServiceClient(deleteAgencyProduct, region)
	if err != nil {
		return diag.Errorf("error creating DLI Client: %s", err)
	}

	deleteAgencyPath := deleteAgencyClient.Endpoint + deleteAgencyHttpUrl
	deleteAgencyPath = strings.ReplaceAll(deleteAgencyPath, "{project_id}", deleteAgencyClient.ProjectID)

	deleteAgencyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	deleteAgencyOpt.JSONBody = map[string]interface{}{
		"roles": make([]string, 0),
	}

	requestResp, err := deleteAgencyClient.Request("POST", deleteAgencyPath, &deleteAgencyOpt)
	if err != nil {
		return diag.Errorf("error deleting DLI Agency: %s", err)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return diag.FromErr(err)
	}
	if !utils.PathSearch("is_success", respBody, true).(bool) {
		return diag.Errorf("unable to delete the agency: %s",
			utils.PathSearch("message", respBody, "Message Not Found"))
	}
	return nil
}
