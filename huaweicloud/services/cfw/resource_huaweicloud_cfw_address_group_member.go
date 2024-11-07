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

	id := utils.PathSearch("data.items[0].id", createAddressGroupMemberRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating AddressGroupMember: ID is not found in API response")
	}
	d.SetId(id)

	return resourceAddressGroupMemberRead(ctx, d, meta)
}

func buildCreateAddressGroupMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	item := map[string]interface{}{
		"name":         utils.ValueIgnoreEmpty(d.Get("name")),
		"address":      utils.ValueIgnoreEmpty(d.Get("address")),
		"address_type": utils.ValueIgnoreEmpty(d.Get("address_type")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}

	bodyParams := map[string]interface{}{
		"set_id":        utils.ValueIgnoreEmpty(d.Get("group_id")),
		"address_items": []map[string]interface{}{item},
	}
	return bodyParams
}

func resourceAddressGroupMemberRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getAddressGroupMember: Query the CFW IP address group member detail
	getAddressGroupMemberProduct := "cfw"
	getAddressGroupMemberClient, err := cfg.NewServiceClient(getAddressGroupMemberProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	addressGroupMembers, respBody, err := ReadAddressGroupMembers(d.Get("group_id").(string), getAddressGroupMemberClient)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving address group member",
		)
	}

	findAddressGroupMemberExpr := fmt.Sprintf("[?item_id == '%s']|[0]", d.Id())
	addressGroupMember := utils.PathSearch(findAddressGroupMemberExpr, addressGroupMembers, nil)
	if addressGroupMember == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving address group member")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("group_id", utils.PathSearch("data.set_id", respBody, nil)),
		d.Set("name", utils.PathSearch("name", addressGroupMember, nil)),
		d.Set("address", utils.PathSearch("address", addressGroupMember, nil)),
		d.Set("address_type", utils.PathSearch("address_type", addressGroupMember, nil)),
		d.Set("description", utils.PathSearch("description", addressGroupMember, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ReadAddressGroupMembers(setID string, client *golangsdk.ServiceClient) ([]interface{}, interface{}, error) {
	httpUrl := "v1/{project_id}/address-items"
	basePath := client.Endpoint + httpUrl
	basePath = strings.ReplaceAll(basePath, "{project_id}", client.ProjectID)
	var result []interface{}
	var respBody interface{}

	offset := 0
	for {
		path := fmt.Sprintf("%s?limit=10&offset=%d&set_id=%s", basePath, offset, setID)
		opt := golangsdk.RequestOpts{
			KeepResponseBody: true,
		}
		resp, err := client.Request("GET", path, &opt)
		if err != nil {
			return nil, nil, err
		}

		respBody, err = utils.FlattenResponse(resp)
		if err != nil {
			return nil, nil, err
		}

		curJson := utils.PathSearch("data.records[*]", respBody, make([]interface{}, 0))
		curArray := curJson.([]interface{})

		if len(curArray) == 0 {
			break
		}

		result = append(result, curArray...)

		offset += 10
	}
	return result, respBody, nil
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
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting address group member",
		)
	}

	return nil
}

func resourceAddressGroupMemberImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <group_id>/<id>")
	}

	d.Set("group_id", parts[0])
	d.SetId(parts[1])

	return []*schema.ResourceData{d}, nil
}
