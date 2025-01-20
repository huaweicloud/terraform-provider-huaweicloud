package cpts

import (
	"context"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CPTS PUT /v1/{project_id}/test-suites/{test_suite_id}
// @API CPTS DELETE /v1/{project_id}/test-suites/{test_suite_id}
// @API CPTS GET /v1/{project_id}/test-suites/{test_suite_id}
// @API CPTS POST /v1/{project_id}/test-suites
func ResourceProject() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceProjectCreate,
		UpdateContext: resourceProjectUpdate,
		DeleteContext: resourceProjectDelete,
		ReadContext:   resourceProjectRead,
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

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCreateProjectBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
}

func resourceProjectCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/test-suites"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildCreateProjectBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error creating CPTS project: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	projectID := utils.PathSearch("project_id", respBody, nil)
	if projectID == nil {
		return diag.Errorf("error creating CPTS project: ID is not found in API response")
	}

	// The `project_id` field is a numeric type.
	d.SetId(strconv.Itoa(int(projectID.(float64))))
	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/test-suites/{test_suite_id}"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{test_suite_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", "SVCSTG.CPTS.4032002"),
			"error retrieving CPTS project")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	createTime := utils.PathSearch("project.create_time", respBody, "").(string)
	updateTime := utils.PathSearch("project.update_time", respBody, "").(string)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("project.name", respBody, nil)),
		d.Set("description", utils.PathSearch("project.description", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampUTC(utils.ConvertTimeStrToNanoTimestamp(createTime)/1000)),
		d.Set("updated_at", utils.FormatTimeStampUTC(utils.ConvertTimeStrToNanoTimestamp(updateTime)/1000)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildUpdateProjectBodyParams(d *schema.ResourceData, idInt64 int64) map[string]interface{} {
	return map[string]interface{}{
		"id":          idInt64,
		"name":        d.Get("name"),
		"description": d.Get("description"),
	}
}

func resourceProjectUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/test-suites/{test_suite_id}"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	idInt64, err := strconv.ParseInt(d.Id(), 10, 32)
	if err != nil {
		return diag.Errorf("the project ID must be integer: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{test_suite_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{200, 201, 204},
		JSONBody:         buildUpdateProjectBodyParams(d, idInt64),
	}

	_, err = client.Request("PUT", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error updating CPTS project: %s", err)
	}

	return resourceProjectRead(ctx, d, meta)
}

func resourceProjectDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1/{project_id}/test-suites/{test_suite_id}"
		product = "cpts"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating CPTS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{test_suite_id}", d.Id())
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected403ErrInto404Err(err, "code", "SVCSTG.CPTS.4032002"),
			"error deleting CPTS project")
	}

	return nil
}
