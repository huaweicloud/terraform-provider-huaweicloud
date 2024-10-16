package dataarts

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

var workspaceIdNotFound = "DLG.0818"

// @API DataArtsStudio POST /v2/{project_id}/design/approvals/users
// @API DataArtsStudio GET /v2/{project_id}/design/approvals/users
// @API DataArtsStudio DELETE /v2/{project_id}/design/approvals/users
func ResourceDataArtsArchitectureReviewer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceArchitectureReviewerCreate,
		ReadContext:   resourceArchitectureReviewerRead,
		DeleteContext: resourceArchitectureReviewerDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceArchitectureReviewerImport,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"user_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"email": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"phone_number": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"reviewer_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceArchitectureReviewerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	createReviewerHttpUrl := "v2/{project_id}/design/approvals/users"
	createReviewerProduct := "dataarts"

	reviewerClient, err := cfg.NewServiceClient(createReviewerProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	createReviewerPath := reviewerClient.Endpoint + createReviewerHttpUrl
	createReviewerPath = strings.ReplaceAll(createReviewerPath, "{project_id}", reviewerClient.ProjectID)
	createReviewerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": d.Get("workspace_id").(string)},
	}
	createReviewerOpt.JSONBody = utils.RemoveNil(buildCreateReviewerParams(d))
	createReviewerResp, err := reviewerClient.Request("POST", createReviewerPath, &createReviewerOpt)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio architecture reviewer: %s", err)
	}

	createReviewerBody, err := utils.FlattenResponse(createReviewerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	userName := utils.PathSearch("data.value.user_name", createReviewerBody, "").(string)
	if userName == "" {
		return diag.Errorf("unable to find the user name of the DataArts Studio architecture reviewer from the API response")
	}
	d.SetId(userName)

	return resourceArchitectureReviewerRead(ctx, d, meta)
}

func buildCreateReviewerParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"approver_name": d.Get("user_name"),
		"user_id":       d.Get("user_id"),
		"email":         utils.ValueIgnoreEmpty(d.Get("email")),
		"phone_number":  utils.ValueIgnoreEmpty(d.Get("phone_number")),
		"email_notify":  isFieldExist(d, "email"),
		"sms_notify":    isFieldExist(d, "phone_number"),
	}
	return bodyParams
}

func isFieldExist(d *schema.ResourceData, key string) interface{} {
	_, ok := d.GetOk(key)
	return ok
}

func resourceArchitectureReviewerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	workspaceID := d.Get("workspace_id").(string)

	getArchitectureReviewerHttpUrl := "v2/{project_id}/design/approvals/users?approver_name={user_name}"
	getArchitectureReviewerProduct := "dataarts"

	getArchitectureReviewerClient, err := cfg.NewServiceClient(getArchitectureReviewerProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	getArchitectureReviewerPath := getArchitectureReviewerClient.Endpoint + getArchitectureReviewerHttpUrl
	getArchitectureReviewerPath = strings.ReplaceAll(getArchitectureReviewerPath, "{project_id}", getArchitectureReviewerClient.ProjectID)
	getArchitectureReviewerPath = strings.ReplaceAll(getArchitectureReviewerPath, "{user_name}", d.Id())

	getArchitectureReviewerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	getArchitectureReviewerResp, err := getArchitectureReviewerClient.Request("GET", getArchitectureReviewerPath, &getArchitectureReviewerOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "errors|[0].error_code", workspaceIdNotFound),
			"DataArts Studio architecture reviewer")
	}

	getArchitectureReviewerRespBody, err := utils.FlattenResponse(getArchitectureReviewerResp)
	if err != nil {
		return diag.FromErr(err)
	}

	reviewer := utils.PathSearch("data.value.records|[0]", getArchitectureReviewerRespBody, nil)
	if reviewer == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "DataArts Studio architecture reviewer")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("user_name", utils.PathSearch("approver_name", reviewer, nil)),
		d.Set("user_id", utils.PathSearch("user_id", reviewer, nil)),
		d.Set("reviewer_id", utils.PathSearch("id", reviewer, nil)),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting DataArts Studio architecture reviewer fields: %s", err)
	}

	return nil
}

func resourceArchitectureReviewerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	reviewerProduct := "dataarts"
	reviewerClient, err := cfg.NewServiceClient(reviewerProduct, region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio Client: %s", err)
	}
	delReviewerHttpUrl := "v2/{project_id}/design/approvals/users?approver_ids={reviewer_id}"
	delReviewerPath := reviewerClient.Endpoint + delReviewerHttpUrl
	delReviewerPath = strings.ReplaceAll(delReviewerPath, "{project_id}", reviewerClient.ProjectID)
	delReviewerPath = strings.ReplaceAll(delReviewerPath, "{reviewer_id}", d.Get("reviewer_id").(string))

	workspaceID := d.Get("workspace_id").(string)
	delReviewerOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"workspace": workspaceID},
	}
	_, err = reviewerClient.Request("DELETE", delReviewerPath, &delReviewerOpt)

	if err != nil {
		return diag.Errorf("error deleting DataArts Studio architecture reviewer: %s", err)
	}
	return nil
}

func resourceArchitectureReviewerImport(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <workspace_id>/<user_name>")
	}
	err := d.Set("workspace_id", parts[0])
	if err != nil {
		return nil, fmt.Errorf("error setting DataArts Studio architecture reviewer fields: %s", err)
	}
	d.SetId(parts[1])
	return []*schema.ResourceData{d}, nil
}
