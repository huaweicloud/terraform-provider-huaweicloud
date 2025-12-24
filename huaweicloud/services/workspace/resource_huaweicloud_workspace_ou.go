package workspace

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

var (
	ouNonUpdatableParams = []string{"domain"}
	// WKS.00010099: The Workspace service is disabled.
	workspaceServiceDisabledErrCode = "WKS.00010099"
)

// @API Workspace POST /v2/{project_id}/ous
// @API Workspace GET /v2/{project_id}/ous/{ou_id}
// @API Workspace PUT /v2/{project_id}/ous/{ou_id}
// @API Workspace DELETE /v2/{project_id}/ous/{ou_id}
func ResourceOu() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOuCreate,
		ReadContext:   resourceOuRead,
		UpdateContext: resourceOuUpdate,
		DeleteContext: resourceOuDelete,

		CustomizeDiff: config.FlexibleForceNew(ouNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: resourceOuImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the OU is located.`,
			},
			"ou_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the OU.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The AD domain name to which the OU belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the OU.`,
			},
			"ou_dn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The distinguished name (DN) of the OU.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the AD domain.`,
			},
			// Internal parameter(s).
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildCreateOuBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ou_name":     d.Get("ou_name"),
		"domain":      d.Get("domain"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceOuCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/ous"
	)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildCreateOuBodyParams(d)),
	}

	resp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating OU: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	ouId := utils.PathSearch("id", respBody, "").(string)
	if ouId == "" {
		return diag.Errorf("unable to find OU ID from API response")
	}
	d.SetId(ouId)

	return resourceOuRead(ctx, d, meta)
}

func listOus(client *golangsdk.ServiceClient, queryParams ...string) ([]interface{}, error) {
	var (
		httpUrl = "v2/{project_id}/ous"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
		listOpt = golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json;charset=UTF-8",
			},
		}
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = fmt.Sprintf("%s?limit=%d", listPath, limit)

	if len(queryParams) > 0 && queryParams[0] != "" {
		listPath += queryParams[0]
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		ous := utils.PathSearch("ous", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, ous...)
		if len(ous) < limit {
			break
		}

		offset += len(ous)
	}

	return result, nil
}

func GetOuByName(client *golangsdk.ServiceClient, ouId string) (interface{}, error) {
	ous, err := listOus(client)
	if err != nil {
		return nil, err
	}

	ouInfo := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", ouId), ous, nil)
	if ouInfo == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v2/{project_id}/ous",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the OU '%s' not found", ouId)),
			},
		}
	}

	return ouInfo, nil
}

func resourceOuRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
		ouId   = d.Id()
	)
	client, err := cfg.NewServiceClient("workspace", region)
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	ou, err := GetOuByName(client, ouId)
	if err != nil {
		// WKS.0202: The resource not found.
		// WKS.00010099: The Workspace service is disabled.
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected500ErrInto404Err(
				common.ConvertExpected403ErrInto404Err(err, "error_code", workspaceServiceDisabledErrCode),
				"error_code",
				"WKS.0202",
			),
			fmt.Sprintf("error retrieving OU (%s)", ouId),
		)
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("ou_name", utils.PathSearch("ou_name", ou, nil)),
		d.Set("domain", utils.PathSearch("domain", ou, nil)),
		d.Set("description", utils.PathSearch("description", ou, nil)),
		d.Set("ou_dn", utils.PathSearch("ou_dn", ou, nil)),
		d.Set("domain_id", utils.PathSearch("domain_id", ou, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateOuBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"ou_name":     d.Get("ou_name"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
}

func resourceOuUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/ous/{ou_id}"
	)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{ou_id}", d.Id())

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
		JSONBody: utils.RemoveNil(buildUpdateOuBodyParams(d)),
	}

	_, err = client.Request("PUT", updatePath, &updateOpt)
	if err != nil {
		return diag.Errorf("error updating Workspace OU (%s): %s", d.Id(), err)
	}

	return resourceOuRead(ctx, d, meta)
}

func resourceOuDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		httpUrl = "v2/{project_id}/ous/{ou_id}"
		ouId    = d.Id()
	)
	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating Workspace client: %s", err)
	}

	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{ou_id}", ouId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=UTF-8",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		// WKS.0915: The resource not found.
		// WKS.00010099: The Workspace service is disabled.
		return common.CheckDeletedDiag(
			d,
			common.ConvertExpected500ErrInto404Err(
				common.ConvertExpected403ErrInto404Err(err, "error_code", workspaceServiceDisabledErrCode),
				"error_code",
				"WKS.0915",
			),
			fmt.Sprintf("error deleting OU (%s)", ouId),
		)
	}

	return nil
}

func resourceOuImportState(_ context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	var (
		cfg      = meta.(*config.Config)
		importId = d.Id()
	)

	client, err := cfg.NewServiceClient("workspace", cfg.GetRegion(d))
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace client: %s", err)
	}

	ous, err := listOus(client)
	if err != nil {
		return nil, fmt.Errorf("error retrieving OUs: %s", err)
	}

	// Both the OU name and its ID can be 32-bit UUIDs.
	// The name may be the same as the ID, so import by ID is preferred. If the OU is not found, try importing by name.
	ouInfo := utils.PathSearch(fmt.Sprintf("[?id=='%s']|[0]", importId), ous, nil)
	if ouInfo != nil {
		// Import by OU ID.
		return []*schema.ResourceData{d}, nil
	}

	// Import by OU name.
	ouId := utils.PathSearch(fmt.Sprintf("[?ou_name=='%s']|[0].id", importId), ous, "").(string)
	if ouId == "" {
		return nil, fmt.Errorf("unable to find the OU ID using its name (%s): %s", importId, err)
	}

	d.SetId(ouId)
	return []*schema.ResourceData{d}, d.Set("ou_name", importId)
}
