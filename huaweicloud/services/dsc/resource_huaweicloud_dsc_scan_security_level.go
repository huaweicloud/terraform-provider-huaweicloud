package dsc

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

var scanSecurityLevelNotFoundErrCodes = []string{
	"dsc.10000009", // The DSC instance does not exist.
}

// @API DSC POST /v1/{project_id}/scan-security-levels
// @API DSC GET /v1/{project_id}/scan-security-levels
// @API DSC PUT /v1/{project_id}/scan-security-levels/{level_id}
// @API DSC DELETE /v1/{project_id}/scan-security-levels/{level_id}
func ResourceScanSecurityLevel() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceScanSecurityLevelCreate,
		ReadContext:   resourceScanSecurityLevelRead,
		UpdateContext: resourceScanSecurityLevelUpdate,
		DeleteContext: resourceScanSecurityLevelDelete,

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
			"security_level_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"color_number": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"security_level_desc": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"is_deleted": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"category": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"used_count": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"sort_weight": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func buildScanSecurityLevelBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"security_level_name": d.Get("security_level_name"),
		"color_number":        utils.ValueIgnoreEmpty(d.Get("color_number")),
		"security_level_desc": utils.ValueIgnoreEmpty(d.Get("security_level_desc")),
	}
}

func createScanSecurityLevel(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/scan-security-levels"
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildScanSecurityLevelBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("POST", createPath, &createOpt)
	return err
}

func resourceScanSecurityLevelCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg  = meta.(*config.Config)
		name = d.Get("security_level_name").(string)
	)

	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = createScanSecurityLevel(client, d)
	if err != nil {
		return diag.Errorf("error creating security level (%s): %s", name, err)
	}

	// The create API does not return the security level ID, query the list to get it.
	securityLevels, err := listScanSecurityLevels(client)
	if err != nil {
		return diag.Errorf("error getting security level (%s): %s", name, err)
	}

	levelId := utils.PathSearch(fmt.Sprintf("[?security_level_name=='%s']|[0].level_id", name), securityLevels, "").(string)
	if levelId == "" {
		return diag.Errorf("unable to find the security level ID from API response")
	}

	d.SetId(levelId)

	return resourceScanSecurityLevelRead(ctx, d, meta)
}

func listScanSecurityLevels(client *golangsdk.ServiceClient) ([]interface{}, error) {
	var (
		httpUrl           = "v1/{project_id}/scan-security-levels"
		limit             = 1000
		offset            = 0
		allSecurityLevels = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)

	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		// The 'is_deleted' parameter is set to true to query all security levels, including the disabled ones.
		listPathWithOffset := fmt.Sprintf("%s?limit=%d&offset=%d&is_deleted=true", listPath, limit, offset)
		resp, err := client.Request("GET", listPathWithOffset, &listOpt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return nil, err
		}

		securityLevels := utils.PathSearch("security_levels_list", respBody, make([]interface{}, 0)).([]interface{})
		allSecurityLevels = append(allSecurityLevels, securityLevels...)
		if len(securityLevels) < limit {
			break
		}

		offset += len(securityLevels)
	}

	return allSecurityLevels, nil
}

func GetScanSecurityLevelById(client *golangsdk.ServiceClient, levelId string) (interface{}, error) {
	securityLevels, err := listScanSecurityLevels(client)
	if err != nil {
		return nil, err
	}

	securityLevel := utils.PathSearch(fmt.Sprintf("[?level_id=='%s']|[0]", levelId), securityLevels, nil)
	if securityLevel == nil {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/scan-security-levels",
				RequestId: "NONE",
				Body:      fmt.Appendf(nil, "the security level (%s) does not exist", levelId),
			},
		}
	}

	return securityLevel, nil
}

func resourceScanSecurityLevelRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		levelId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", region)
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	securityLevel, err := GetScanSecurityLevelById(client, levelId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", scanSecurityLevelNotFoundErrCodes...),
			fmt.Sprintf("error retrieving security level (%s)", levelId),
		)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("security_level_name", utils.PathSearch("security_level_name", securityLevel, nil)),
		d.Set("color_number", utils.PathSearch("color_number", securityLevel, nil)),
		d.Set("security_level_desc", utils.PathSearch("security_level_desc", securityLevel, nil)),
		d.Set("is_deleted", utils.PathSearch("is_deleted", securityLevel, nil)),
		d.Set("category", utils.PathSearch("category", securityLevel, nil)),
		d.Set("used_count", utils.PathSearch("used_count", securityLevel, nil)),
		d.Set("sort_weight", utils.PathSearch("sort_weight", securityLevel, nil)),
		d.Set("create_time", utils.PathSearch("create_time", securityLevel, nil)),
		d.Set("update_time", utils.PathSearch("update_time", securityLevel, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func updateScanSecurityLevel(client *golangsdk.ServiceClient, levelId string, d *schema.ResourceData) error {
	httpUrl := "v1/{project_id}/scan-security-levels/{level_id}"
	updatePath := client.Endpoint + httpUrl
	updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
	updatePath = strings.ReplaceAll(updatePath, "{level_id}", levelId)

	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         utils.RemoveNil(buildScanSecurityLevelBodyParams(d)),
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceScanSecurityLevelUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		levelId = d.Id()
	)

	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	if d.HasChanges("security_level_name", "color_number", "security_level_desc") {
		err = updateScanSecurityLevel(client, levelId, d)
		if err != nil {
			return diag.Errorf("error updating security level (%s): %s", levelId, err)
		}
	}

	return resourceScanSecurityLevelRead(ctx, d, meta)
}

func deleteScanSecurityLevel(client *golangsdk.ServiceClient, levelId string) error {
	httpUrl := "v1/{project_id}/scan-security-levels/{level_id}"
	deletePath := client.Endpoint + httpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{level_id}", levelId)

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err := client.Request("DELETE", deletePath, &deleteOpt)
	return err
}

func resourceScanSecurityLevelDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.NewServiceClient("dsc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating DSC client: %s", err)
	}

	err = deleteScanSecurityLevel(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected401ErrInto404Err(err, "error_code", scanSecurityLevelNotFoundErrCodes...),
			fmt.Sprintf("error deleting security level (%v)", d.Get("security_level_name")),
		)
	}

	return nil
}
