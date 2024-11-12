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

// @API CFW POST /v1/{project_id}/address-set
// @API CFW DELETE /v1/{project_id}/address-sets/{id}
// @API CFW GET /v1/{project_id}/address-sets/{id}
// @API CFW PUT /v1/{project_id}/address-sets/{id}
func ResourceAddressGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddressGroupCreate,
		UpdateContext: resourceAddressGroupUpdate,
		ReadContext:   resourceAddressGroupRead,
		DeleteContext: resourceAddressGroupDelete,
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
				Description: `Specifies the IP address group name.`,
			},
			// address_type is an attribute in the doc
			// because it only works in cn-south-4
			"address_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `schema: Computed; Specifies the Address type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the Address group description.`,
			},
		},
	}
}

func resourceAddressGroupCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAddressGroup: Create a CFW IP address group.
	var (
		createAddressGroupHttpUrl = "v1/{project_id}/address-set"
		createAddressGroupProduct = "cfw"
	)
	createAddressGroupClient, err := cfg.NewServiceClient(createAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	createAddressGroupPath := createAddressGroupClient.Endpoint + createAddressGroupHttpUrl
	createAddressGroupPath = strings.ReplaceAll(createAddressGroupPath, "{project_id}",
		createAddressGroupClient.ProjectID)

	createAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAddressGroupOpt.JSONBody = utils.RemoveNil(buildCreateAddressGroupBodyParams(d))
	createAddressGroupResp, err := createAddressGroupClient.Request("POST", createAddressGroupPath,
		&createAddressGroupOpt)
	if err != nil {
		return diag.Errorf("error creating AddressGroup: %s", err)
	}

	createAddressGroupRespBody, err := utils.FlattenResponse(createAddressGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createAddressGroupRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating AddressGroup: ID is not found in API response")
	}
	d.SetId(id)

	return resourceAddressGroupRead(ctx, d, meta)
}

func buildCreateAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id":    utils.ValueIgnoreEmpty(d.Get("object_id")),
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"address_type": utils.ValueIgnoreEmpty(d.Get("address_type")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceAddressGroupRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAddressGroup: Query the CFW IP address group detail
	var (
		getAddressGroupHttpUrl = "v1/{project_id}/address-sets/{id}"
		getAddressGroupProduct = "cfw"
	)
	getAddressGroupClient, err := cfg.NewServiceClient(getAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	getAddressGroupPath := getAddressGroupClient.Endpoint + getAddressGroupHttpUrl
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{project_id}",
		getAddressGroupClient.ProjectID)
	getAddressGroupPath = strings.ReplaceAll(getAddressGroupPath, "{id}", d.Id())

	getAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAddressGroupResp, err := getAddressGroupClient.Request("GET", getAddressGroupPath,
		&getAddressGroupOpt)

	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving AddressGroup",
		)
	}

	getAddressGroupRespBody, err := utils.FlattenResponse(getAddressGroupResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("data.name", getAddressGroupRespBody, nil)),
		d.Set("address_type", utils.PathSearch("data.address_type", getAddressGroupRespBody, nil)),
		d.Set("description", utils.PathSearch("data.description", getAddressGroupRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceAddressGroupUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateAddressGroupChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateAddressGroupChanges...) {
		// updateAddressGroup: Update the configuration of CFW IP address group
		var (
			updateAddressGroupHttpUrl = "v1/{project_id}/address-sets/{id}"
			updateAddressGroupProduct = "cfw"
		)
		updateAddressGroupClient, err := cfg.NewServiceClient(updateAddressGroupProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW Client: %s", err)
		}

		updateAddressGroupPath := updateAddressGroupClient.Endpoint + updateAddressGroupHttpUrl
		updateAddressGroupPath = strings.ReplaceAll(updateAddressGroupPath, "{project_id}",
			updateAddressGroupClient.ProjectID)
		updateAddressGroupPath = strings.ReplaceAll(updateAddressGroupPath, "{id}", d.Id())

		updateAddressGroupOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateAddressGroupOpt.JSONBody = utils.RemoveNil(buildUpdateAddressGroupBodyParams(d))
		_, err = updateAddressGroupClient.Request("PUT", updateAddressGroupPath, &updateAddressGroupOpt)
		if err != nil {
			return diag.Errorf("error updating AddressGroup: %s", err)
		}
	}
	return resourceAddressGroupRead(ctx, d, meta)
}

func buildUpdateAddressGroupBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        utils.ValueIgnoreEmpty(d.Get("name")),
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceAddressGroupDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAddressGroup: Delete an existing CFW IP address group
	var (
		deleteAddressGroupHttpUrl = "v1/{project_id}/address-sets/{id}"
		deleteAddressGroupProduct = "cfw"
	)
	deleteAddressGroupClient, err := cfg.NewServiceClient(deleteAddressGroupProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	deleteAddressGroupPath := deleteAddressGroupClient.Endpoint + deleteAddressGroupHttpUrl
	deleteAddressGroupPath = strings.ReplaceAll(deleteAddressGroupPath, "{project_id}",
		deleteAddressGroupClient.ProjectID)
	deleteAddressGroupPath = strings.ReplaceAll(deleteAddressGroupPath, "{id}", d.Id())

	deleteAddressGroupOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteAddressGroupClient.Request("DELETE", deleteAddressGroupPath, &deleteAddressGroupOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting AddressGroup",
		)
	}

	return nil
}
