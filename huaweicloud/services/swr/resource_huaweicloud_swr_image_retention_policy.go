// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product SWR
// ---------------------------------------------------------------

package swr

import (
	"context"
	"errors"
	"fmt"
	"log"
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

// @API SWR POST /v2/manage/namespaces/{namespace}/repos/{repository}/retentions
// @API SWR DELETE /v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}
// @API SWR GET /v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}
// @API SWR PATCH /v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}
func ResourceSwrImageRetentionPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSwrImageRetentionPolicyCreate,
		UpdateContext: resourceSwrImageRetentionPolicyUpdate,
		ReadContext:   resourceSwrImageRetentionPolicyRead,
		DeleteContext: resourceSwrImageRetentionPolicyDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceSwrImageRetentionPolicyImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"organization": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the organization.`,
			},
			"repository": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the name of the repository.`,
			},
			"type": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the retention policy type.`,
			},
			"number": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the number of retention.`,
			},
			"tag_selectors": {
				Type:        schema.TypeList,
				Elem:        ImageRetentionPolicyTagSelectorSchema(),
				Optional:    true,
				Computed:    true,
				Description: `Specifies the image tags that are not counted in the retention policy`,
			},
			"retention_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the retention ID.`,
			},
		},
	}
}

func ImageRetentionPolicyTagSelectorSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"kind": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the Matching rule.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the Matching pattern.`,
			},
		},
	}
	return &sc
}

func resourceSwrImageRetentionPolicyCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createSwrImageRetentionPolicy: create SWR image retention policy
	var (
		createSwrImageRetentionPolicyHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/retentions"
		createSwrImageRetentionPolicyProduct = "swr"
	)
	createSwrImageRetentionPolicyClient, err := cfg.NewServiceClient(createSwrImageRetentionPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)

	createSwrImageRetentionPolicyPath := createSwrImageRetentionPolicyClient.Endpoint + createSwrImageRetentionPolicyHttpUrl
	createSwrImageRetentionPolicyPath = strings.ReplaceAll(createSwrImageRetentionPolicyPath, "{namespace}",
		organization)
	createSwrImageRetentionPolicyPath = strings.ReplaceAll(createSwrImageRetentionPolicyPath, "{repository}",
		strings.ReplaceAll(repository, "/", "$"))

	createSwrImageRetentionPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createSwrImageRetentionPolicyOpt.JSONBody = buildSwrImageRetentionPolicyBodyParams(d)
	createSwrImageRetentionPolicyResp, err := createSwrImageRetentionPolicyClient.Request("POST",
		createSwrImageRetentionPolicyPath, &createSwrImageRetentionPolicyOpt)
	if err != nil {
		return diag.Errorf("error creating SWR image retention policy: %s", err)
	}

	createSwrImageRetentionPolicyRespBody, err := utils.FlattenResponse(createSwrImageRetentionPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", createSwrImageRetentionPolicyRespBody, float64(-1)).(float64)
	if id == -1 {
		return diag.Errorf("error creating SWR image retention policy: ID is not found in API response")
	}

	d.SetId(organization + "/" + repository + "/" + strconv.Itoa(int(id)))
	d.Set("retention_id", strconv.Itoa(int(id)))

	return resourceSwrImageRetentionPolicyRead(ctx, d, meta)
}

func buildSwrImageRetentionPolicyBodyParams(d *schema.ResourceData) map[string]interface{} {
	template := utils.ValueIgnoreEmpty(d.Get("type"))
	number := utils.ValueIgnoreEmpty(d.Get("number")).(int)
	param := map[string]interface{}{
		"template":      utils.ValueIgnoreEmpty(d.Get("type")),
		"tag_selectors": buildSwrImageRetentionPolicyTagSelectorsChildBody(d),
	}
	if template == "date_rule" {
		param["params"] = map[string]interface{}{
			"days": strconv.Itoa(number),
		}
	} else {
		param["params"] = map[string]interface{}{
			"num": strconv.Itoa(number),
		}
	}
	bodyParams := map[string]interface{}{
		"algorithm": "or",
		"rules":     []map[string]interface{}{param},
	}
	return bodyParams
}

func buildSwrImageRetentionPolicyTagSelectorsChildBody(d *schema.ResourceData) []map[string]interface{} {
	params := make([]map[string]interface{}, 0)
	rawParams := d.Get("tag_selectors").([]interface{})
	if len(rawParams) == 0 {
		return params
	}

	for _, rawParam := range rawParams {
		raw := rawParam.(map[string]interface{})
		param := map[string]interface{}{
			"kind":    utils.ValueIgnoreEmpty(raw["kind"]),
			"pattern": utils.ValueIgnoreEmpty(raw["pattern"]),
		}
		params = append(params, param)
	}

	return params
}

func resourceSwrImageRetentionPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getSwrImageRetentionPolicy: Query SWR image retention policy
	var (
		getSwrImageRetentionPolicyHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}"
		getSwrImageRetentionPolicyProduct = "swr"
	)
	getSwrImageRetentionPolicyClient, err := cfg.NewServiceClient(getSwrImageRetentionPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	retentionId := d.Get("retention_id").(string)

	getSwrImageRetentionPolicyPath := getSwrImageRetentionPolicyClient.Endpoint + getSwrImageRetentionPolicyHttpUrl
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{namespace}", organization)
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{repository}", strings.ReplaceAll(repository, "/", "$"))
	getSwrImageRetentionPolicyPath = strings.ReplaceAll(getSwrImageRetentionPolicyPath, "{retention_id}", retentionId)

	getSwrImageRetentionPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getSwrImageRetentionPolicyResp, err := getSwrImageRetentionPolicyClient.Request("GET",
		getSwrImageRetentionPolicyPath, &getSwrImageRetentionPolicyOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving SWR image retention policy")
	}

	getSwrImageRetentionPolicyRespBody, err := utils.FlattenResponse(getSwrImageRetentionPolicyResp)
	if err != nil {
		return diag.FromErr(err)
	}

	policies := utils.PathSearch("rules", getSwrImageRetentionPolicyRespBody,
		make([]interface{}, 0)).([]interface{})
	if len(policies) == 0 {
		log.Printf("[WARN] failed to get SWR image retention policy by organization(%s),"+
			"repository(%s) and retention_id(%s)", organization, repository, retentionId)
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "")
	}

	template := utils.PathSearch("template", policies[0], nil)
	var number int
	if template.(string) == "date_rule" {
		number, _ = strconv.Atoi(utils.PathSearch("params.days", policies[0], "0").(string))
	} else {
		number, _ = strconv.Atoi(utils.PathSearch("params.num", policies[0], "0").(string))
	}
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("organization", organization),
		d.Set("repository", repository),
		d.Set("type", template),
		d.Set("number", number),
		d.Set("tag_selectors", flattenGetImageRetentionPolicyResponseBodyTagSelector(policies[0])),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetImageRetentionPolicyResponseBodyTagSelector(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("tag_selectors", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		rst = append(rst, map[string]interface{}{
			"kind":    utils.PathSearch("kind", v, nil),
			"pattern": utils.PathSearch("pattern", v, nil),
		})
	}
	return rst
}

func resourceSwrImageRetentionPolicyUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateSwrImageRetentionPolicyHasChanges := []string{
		"number",
		"tag_selectors",
	}

	if d.HasChanges(updateSwrImageRetentionPolicyHasChanges...) {
		// updateSwrImageRetentionPolicy: update SWR image retention policy
		var (
			updateSwrImageRetentionPolicyHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}"
			updateSwrImageRetentionPolicyProduct = "swr"
		)
		updateSwrImageRetentionPolicyClient, err := cfg.NewServiceClient(updateSwrImageRetentionPolicyProduct, region)
		if err != nil {
			return diag.Errorf("error creating SWR client: %s", err)
		}

		organization := d.Get("organization").(string)
		repository := d.Get("repository").(string)
		retentionId := d.Get("retention_id").(string)

		updateSwrImageRetentionPolicyPath := updateSwrImageRetentionPolicyClient.Endpoint + updateSwrImageRetentionPolicyHttpUrl
		updateSwrImageRetentionPolicyPath = strings.ReplaceAll(updateSwrImageRetentionPolicyPath, "{namespace}", organization)
		updateSwrImageRetentionPolicyPath = strings.ReplaceAll(updateSwrImageRetentionPolicyPath, "{repository}",
			strings.ReplaceAll(repository, "/", "$"))
		updateSwrImageRetentionPolicyPath = strings.ReplaceAll(updateSwrImageRetentionPolicyPath, "{retention_id}", retentionId)

		updateSwrImageRetentionPolicyOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				201,
			},
		}
		updateSwrImageRetentionPolicyOpt.JSONBody = buildSwrImageRetentionPolicyBodyParams(d)
		_, err = updateSwrImageRetentionPolicyClient.Request("PATCH", updateSwrImageRetentionPolicyPath, &updateSwrImageRetentionPolicyOpt)
		if err != nil {
			return diag.Errorf("error updating SWR image retention policy: %s", err)
		}
	}
	return resourceSwrImageRetentionPolicyRead(ctx, d, meta)
}

func resourceSwrImageRetentionPolicyDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteSwrImageRetentionPolicy: Delete SWR image retention policy
	var (
		deleteSwrImageRetentionPolicyHttpUrl = "v2/manage/namespaces/{namespace}/repos/{repository}/retentions/{retention_id}"
		deleteSwrImageRetentionPolicyProduct = "swr"
	)
	deleteSwrImageRetentionPolicyClient, err := cfg.NewServiceClient(deleteSwrImageRetentionPolicyProduct, region)
	if err != nil {
		return diag.Errorf("error creating SWR client: %s", err)
	}

	organization := d.Get("organization").(string)
	repository := d.Get("repository").(string)
	retentionId := d.Get("retention_id").(string)

	deleteSwrImageRetentionPolicyPath := deleteSwrImageRetentionPolicyClient.Endpoint + deleteSwrImageRetentionPolicyHttpUrl
	deleteSwrImageRetentionPolicyPath = strings.ReplaceAll(deleteSwrImageRetentionPolicyPath, "{namespace}", organization)
	deleteSwrImageRetentionPolicyPath = strings.ReplaceAll(deleteSwrImageRetentionPolicyPath, "{repository}",
		strings.ReplaceAll(repository, "/", "$"))
	deleteSwrImageRetentionPolicyPath = strings.ReplaceAll(deleteSwrImageRetentionPolicyPath, "{retention_id}", retentionId)

	deleteSwrImageRetentionPolicyOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteSwrImageRetentionPolicyClient.Request("DELETE",
		deleteSwrImageRetentionPolicyPath, &deleteSwrImageRetentionPolicyOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "errors|[0].errorCode", "SVCSTG.SWR.4000306"),
			"error deleting SWR image retention policy")
	}

	return nil
}

func resourceSwrImageRetentionPolicyImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), ",")
	if len(parts) != 3 {
		parts = strings.Split(d.Id(), "/")
		if len(parts) != 3 {
			return nil, errors.New("invalid id format, must be <organization_name>/<repository_name>/<retention_id> or " +
				"<organization_name>,<repository_name>,<retention_id>")
		}
	} else {
		// reform ID to be separated by slashes
		id := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
		d.SetId(id)
	}

	d.Set("organization", parts[0])
	d.Set("repository", parts[1])
	d.Set("retention_id", parts[2])

	return []*schema.ResourceData{d}, nil
}
