// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/jmespath/go-jmespath"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW POST /v1/{project_id}/address-items
// @API CFW GET /v1/{project_id}/address-items
// @API CFW DELETE /v1/{project_id}/address-items/{id}
func ResourceAddressGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAddressGroupMemberCreate,
		ReadContext:   resourceAddressGroupMemberRead,
		DeleteContext: resourceAddressGroupMemberDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceAddressGroupMemberImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"group_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the IP address group.`,
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the IP address.`,
			},
			// address_type is an attribute in the doc
			// because it only works in cn-south-4
			"address_type": {
				Type:        schema.TypeInt,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies the address type.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `Specifies address description.`,
			},
			// Deprecated
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `schema: Deprecated; Specifies the address name.`,
			},
		},
	}
}

func resourceAddressGroupMemberCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createAddressGroupMember: Create a CFW IP address group member.
	var (
		createAddressGroupMemberHttpUrl = "v1/{project_id}/address-items"
		createAddressGroupMemberProduct = "cfw"
	)
	createAddressGroupMemberClient, err := cfg.NewServiceClient(createAddressGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	createAddressGroupMemberPath := createAddressGroupMemberClient.Endpoint + createAddressGroupMemberHttpUrl
	createAddressGroupMemberPath = strings.ReplaceAll(createAddressGroupMemberPath, "{project_id}",
		createAddressGroupMemberClient.ProjectID)

	createAddressGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createAddressGroupMemberOpt.JSONBody = utils.RemoveNil(buildCreateAddressGroupMemberBodyParams(d))
	createAddressGroupMemberResp, err := createAddressGroupMemberClient.Request("POST",
		createAddressGroupMemberPath, &createAddressGroupMemberOpt)
	if err != nil {
		return diag.Errorf("error creating AddressGroupMember: %s", err)
	}

	createAddressGroupMemberRespBody, err := utils.FlattenResponse(createAddressGroupMemberResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("data.items[0].id", createAddressGroupMemberRespBody)
	if err != nil {
		return diag.Errorf("error creating AddressGroupMember: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceAddressGroupMemberRead(ctx, d, meta)
}

func buildCreateAddressGroupMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	item := map[string]interface{}{
		"name":         utils.ValueIngoreEmpty(d.Get("name")),
		"address":      utils.ValueIngoreEmpty(d.Get("address")),
		"address_type": utils.ValueIngoreEmpty(d.Get("address_type")),
		"description":  utils.ValueIngoreEmpty(d.Get("description")),
	}

	bodyParams := map[string]interface{}{
		"set_id":        utils.ValueIngoreEmpty(d.Get("group_id")),
		"address_items": []map[string]interface{}{item},
	}
	return bodyParams
}

func resourceAddressGroupMemberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAddressGroupMember: Query the CFW IP address group member detail
	var (
		getAddressGroupMemberHttpUrl = "v1/{project_id}/address-items"
		getAddressGroupMemberProduct = "cfw"
	)
	getAddressGroupMemberClient, err := cfg.NewServiceClient(getAddressGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	getAddressGroupMemberPath := getAddressGroupMemberClient.Endpoint + getAddressGroupMemberHttpUrl
	getAddressGroupMemberPath = strings.ReplaceAll(getAddressGroupMemberPath, "{project_id}",
		getAddressGroupMemberClient.ProjectID)

	getAddressGroupMemberqueryParams := buildGetAddressGroupMemberQueryParams(d)
	getAddressGroupMemberPath += getAddressGroupMemberqueryParams

	getAddressGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getAddressGroupMemberResp, err := getAddressGroupMemberClient.Request("GET", getAddressGroupMemberPath,
		&getAddressGroupMemberOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving AddressGroupMember")
	}

	getAddressGroupMemberRespBody, err := utils.FlattenResponse(getAddressGroupMemberResp)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(utils.PathSearch("data.records[0].item_id", getAddressGroupMemberRespBody, "").(string))

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("data.set_id", getAddressGroupMemberRespBody, nil)),
		d.Set("name", utils.PathSearch("data.records[0].name", getAddressGroupMemberRespBody, nil)),
		d.Set("address", utils.PathSearch("data.records[0].address", getAddressGroupMemberRespBody, nil)),
		d.Set("address_type", utils.PathSearch("data.records[0].address_type", getAddressGroupMemberRespBody, nil)),
		d.Set("description", utils.PathSearch("data.records[0].description", getAddressGroupMemberRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAddressGroupMemberQueryParams(d *schema.ResourceData) string {
	res := "?offset=0&limit=10"

	res = fmt.Sprintf("%s&set_id=%v", res, d.Get("group_id"))

	res = fmt.Sprintf("%s&address=%v", res, d.Get("address"))

	return res
}

func resourceAddressGroupMemberDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteAddressGroupMember: Delete an existing CFW IP address group
	var (
		deleteAddressGroupMemberHttpUrl = "v1/{project_id}/address-items/{id}"
		deleteAddressGroupMemberProduct = "cfw"
	)
	deleteAddressGroupMemberClient, err := cfg.NewServiceClient(deleteAddressGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	deleteAddressGroupMemberPath := deleteAddressGroupMemberClient.Endpoint + deleteAddressGroupMemberHttpUrl
	deleteAddressGroupMemberPath = strings.ReplaceAll(deleteAddressGroupMemberPath, "{project_id}",
		deleteAddressGroupMemberClient.ProjectID)
	deleteAddressGroupMemberPath = strings.ReplaceAll(deleteAddressGroupMemberPath, "{id}", d.Id())

	deleteAddressGroupMemberOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteAddressGroupMemberClient.Request("DELETE", deleteAddressGroupMemberPath, &deleteAddressGroupMemberOpt)
	if err != nil {
		return diag.Errorf("error deleting AddressGroupMember: %s", err)
	}

	return nil
}

func resourceAddressGroupMemberImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <group_id>/<address>")
	}

	d.Set("group_id", parts[0])
	d.Set("address", parts[1])

	return []*schema.ResourceData{d}, nil
}
