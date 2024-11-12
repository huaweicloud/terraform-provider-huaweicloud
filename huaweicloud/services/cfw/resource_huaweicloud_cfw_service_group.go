// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

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

// @API CFW POST /v1/{project_id}/service-set
// @API CFW DELETE /v1/{project_id}/service-sets/{id}
// @API CFW GET /v1/{project_id}/service-sets/{id}
// @API CFW PUT /v1/{project_id}/service-sets/{id}
func ResourceServiceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceServiceGroupCreate,
		UpdateContext: resourceServiceGroupUpdate,
		ReadContext:   resourceServiceGroupRead,
		DeleteContext: resourceServiceGroupDelete,
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
			"object_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the protected object ID.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the service group name.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the service group description.`,
			},
		},
	}
}

func resourceServiceGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createServiceGroup: Create a CFW service group.
	var (
		createServiceGroupHttpUrl = "v1/{project_id}/service-set"
		createServiceGroupProduct = "cfw"
	)
	createServiceGroupClient, err := cfg.NewServiceClient(createServiceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	createServiceGroupPath := createServiceGroupClient.Endpoint + createServiceGroupHttpUrl
	createServiceGroupPath = strings.ReplaceAll(createServiceGroupPath, "{project_id}",
		createServiceGroupClient.ProjectID)

	createServiceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createServiceGroupOpt.JSONBody = utils.RemoveNil(buildCreateServiceGroupBodyParams(d))
	createServiceGroupResp, err := createServiceGroupClient.Request("POST", createServiceGroupPath,
		&createServiceGroupOpt)
	if err != nil {
		return diag.Errorf("error creating ServiceGroup: %s", err)
	}

	createServiceGroupRespBody, err := utils.FlattenResponse(createServiceGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createServiceGroupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating ServiceGroup: ID is not found in API response")
	}
	d.SetId(id)

	return resourceServiceGroupRead(ctx, d, meta)
}

func buildCreateServiceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id":   utils.ValueIgnoreEmpty(d.Get("object_id")),
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceServiceGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getServiceGroup: Query the CFW service group detail
	var (
		getServiceGroupHttpUrl = "v1/{project_id}/service-sets/{id}"
		getServiceGroupProduct = "cfw"
	)
	getServiceGroupClient, err := cfg.NewServiceClient(getServiceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	getServiceGroupPath := getServiceGroupClient.Endpoint + getServiceGroupHttpUrl
	getServiceGroupPath = strings.ReplaceAll(getServiceGroupPath, "{project_id}",
		getServiceGroupClient.ProjectID)
	getServiceGroupPath = strings.ReplaceAll(getServiceGroupPath, "{id}", d.Id())

	getServiceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getServiceGroupResp, err := getServiceGroupClient.Request("GET", getServiceGroupPath,
		&getServiceGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving ServiceGroup",
		)
	}

	getServiceGroupRespBody, err := utils.FlattenResponse(getServiceGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("data.name", getServiceGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("data.description", getServiceGroupRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceServiceGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateServiceGroupChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateServiceGroupChanges...) {
		// updateServiceGroup: Update the configuration of CFW service group
		var (
			updateServiceGroupHttpUrl = "v1/{project_id}/service-sets/{id}"
			updateServiceGroupProduct = "cfw"
		)
		updateServiceGroupClient, err := cfg.NewServiceClient(updateServiceGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW Client: %s", err)
		}

		updateServiceGroupPath := updateServiceGroupClient.Endpoint + updateServiceGroupHttpUrl
		updateServiceGroupPath = strings.ReplaceAll(updateServiceGroupPath, "{project_id}",
			updateServiceGroupClient.ProjectID)
		updateServiceGroupPath = strings.ReplaceAll(updateServiceGroupPath, "{id}", d.Id())

		updateServiceGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateServiceGroupOpt.JSONBody = utils.RemoveNil(buildUpdateServiceGroupBodyParams(d))
		_, err = updateServiceGroupClient.Request("PUT", updateServiceGroupPath, &updateServiceGroupOpt)
		if err != nil {
			return diag.Errorf("error updating ServiceGroup: %s", err)
		}
	}
	return resourceServiceGroupRead(ctx, d, meta)
}

func buildUpdateServiceGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceServiceGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteServiceGroup: Delete an existing CFW service group
	var (
		deleteServiceGroupHttpUrl = "v1/{project_id}/service-sets/{id}"
		deleteServiceGroupProduct = "cfw"
	)
	deleteServiceGroupClient, err := cfg.NewServiceClient(deleteServiceGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	deleteServiceGroupPath := deleteServiceGroupClient.Endpoint + deleteServiceGroupHttpUrl
	deleteServiceGroupPath = strings.ReplaceAll(deleteServiceGroupPath, "{project_id}",
		deleteServiceGroupClient.ProjectID)
	deleteServiceGroupPath = strings.ReplaceAll(deleteServiceGroupPath, "{id}", d.Id())

	deleteServiceGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteServiceGroupClient.Request("DELETE", deleteServiceGroupPath, &deleteServiceGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting ServiceGroup",
		)
	}

	return nil
}
