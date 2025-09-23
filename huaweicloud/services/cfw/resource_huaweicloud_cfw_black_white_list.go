// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CFW
// ---------------------------------------------------------------

package cfw

import (
	"context"
	"fmt"
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

// @API CFW POST /v1/{project_id}/black-white-list
// @API CFW DELETE /v1/{project_id}/black-white-list/{id}
// @API CFW PUT /v1/{project_id}/black-white-list/{id}
// @API CFW GET /v1/{project_id}/black-white-lists
func ResourceBlackWhiteList() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBlackWhiteListCreate,
		UpdateContext: resourceBlackWhiteListUpdate,
		ReadContext:   resourceBlackWhiteListRead,
		DeleteContext: resourceBlackWhiteListDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourceBlackWhiteListImportState,
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
			"list_type": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the list type.`,
			},
			"direction": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the address direction.`,
			},
			"protocol": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the protocol type.`,
			},
			"port": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the destination port.`,
			},
			"address_type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: `Specifies the IP address type.`,
			},
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the address.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description.`,
			},
		},
	}
}

func resourceBlackWhiteListCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createBlackWhiteList: Create a CFW black white list.
	var (
		createBlackWhiteListHttpUrl = "v1/{project_id}/black-white-list"
		createBlackWhiteListProduct = "cfw"
	)
	createBlackWhiteListClient, err := cfg.NewServiceClient(createBlackWhiteListProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	createBlackWhiteListPath := createBlackWhiteListClient.Endpoint + createBlackWhiteListHttpUrl
	createBlackWhiteListPath = strings.ReplaceAll(createBlackWhiteListPath, "{project_id}",
		createBlackWhiteListClient.ProjectID)

	createBlackWhiteListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	createBlackWhiteListOpt.JSONBody = utils.RemoveNil(buildCreateBlackWhiteListBodyParams(d))
	createBlackWhiteListResp, err := createBlackWhiteListClient.Request("POST", createBlackWhiteListPath,
		&createBlackWhiteListOpt)
	if err != nil {
		return diag.Errorf("error creating black white list: %s", err)
	}

	createBlackWhiteListRespBody, err := utils.FlattenResponse(createBlackWhiteListResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("data.id", createBlackWhiteListRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating black white list: ID is not found in API response")
	}
	d.SetId(id)

	return resourceBlackWhiteListRead(ctx, d, meta)
}

func buildCreateBlackWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id": utils.ValueIgnoreEmpty(d.Get("object_id")),
		"list_type": utils.ValueIgnoreEmpty(d.Get("list_type")),
		// direction can be 0
		"direction": d.Get("direction"),
		"protocol":  utils.ValueIgnoreEmpty(d.Get("protocol")),
		"port":      utils.ValueIgnoreEmpty(d.Get("port")),
		// address_type can be 0
		"address_type": d.Get("address_type"),
		"address":      utils.ValueIgnoreEmpty(d.Get("address")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceBlackWhiteListRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getBlackWhiteList: Query the CFW black white list detail
	var (
		getBlackWhiteListHttpUrl = "v1/{project_id}/black-white-lists"
		getBlackWhiteListProduct = "cfw"
	)
	getBlackWhiteListClient, err := cfg.NewServiceClient(getBlackWhiteListProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	getBlackWhiteListPath := getBlackWhiteListClient.Endpoint + getBlackWhiteListHttpUrl
	getBlackWhiteListPath = strings.ReplaceAll(getBlackWhiteListPath, "{project_id}",
		getBlackWhiteListClient.ProjectID)

	getBlackWhiteListqueryParams := buildGetBlackWhiteListQueryParams(d)
	getBlackWhiteListPath += getBlackWhiteListqueryParams

	getBlackWhiteListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getBlackWhiteListResp, err := getBlackWhiteListClient.Request("GET", getBlackWhiteListPath,
		&getBlackWhiteListOpt)

	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error retrieving black white list",
		)
	}

	getBlackWhiteListRespBody, err := utils.FlattenResponse(getBlackWhiteListResp)
	if err != nil {
		return diag.FromErr(err)
	}

	lists := utils.PathSearch("data.records", getBlackWhiteListRespBody, nil)
	if lists == nil {
		return diag.Errorf("error parsing data.records from response= %#v", getBlackWhiteListRespBody)
	}

	val, ok := lists.([]interface{})
	if !ok {
		diag.Errorf("data.records is not a list, data.records= %#v", lists)
	}

	if len(val) != 1 {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error retrieving black white list")
	}

	list := val[0]

	d.SetId(utils.PathSearch("list_id", list, "").(string))
	// list_type not returned
	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("direction", utils.PathSearch("direction", list, nil)),
		d.Set("protocol", utils.PathSearch("protocol", list, nil)),
		d.Set("port", utils.PathSearch("port", list, nil)),
		d.Set("address_type", utils.PathSearch("address_type", list, nil)),
		d.Set("address", utils.PathSearch("address", list, nil)),
		d.Set("description", utils.PathSearch("description", list, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetBlackWhiteListQueryParams(d *schema.ResourceData) string {
	res := "?offset=0&limit=10"
	res = fmt.Sprintf("%s&object_id=%v", res, d.Get("object_id"))
	res = fmt.Sprintf("%s&list_type=%v", res, d.Get("list_type"))
	res = fmt.Sprintf("%s&address=%v", res, strings.Split(d.Get("address").(string), "/")[0])

	return res
}

func resourceBlackWhiteListUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateBlackWhiteListChanges := []string{
		"direction",
		"protocol",
		"port",
		"address_type",
		"address",
		"description",
	}

	if d.HasChanges(updateBlackWhiteListChanges...) {
		// updateBlackWhiteList: Update the configuration of CFW black white list
		var (
			updateBlackWhiteListHttpUrl = "v1/{project_id}/black-white-list/{id}"
			updateBlackWhiteListProduct = "cfw"
		)
		updateBlackWhiteListClient, err := cfg.NewServiceClient(updateBlackWhiteListProduct, region)
		if err != nil {
			return diag.Errorf("error creating CFW Client: %s", err)
		}

		updateBlackWhiteListPath := updateBlackWhiteListClient.Endpoint + updateBlackWhiteListHttpUrl
		updateBlackWhiteListPath = strings.ReplaceAll(updateBlackWhiteListPath, "{project_id}",
			updateBlackWhiteListClient.ProjectID)
		updateBlackWhiteListPath = strings.ReplaceAll(updateBlackWhiteListPath, "{id}", d.Id())

		updateBlackWhiteListOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateBlackWhiteListOpt.JSONBody = utils.RemoveNil(buildUpdateBlackWhiteListBodyParams(d))
		_, err = updateBlackWhiteListClient.Request("PUT", updateBlackWhiteListPath, &updateBlackWhiteListOpt)
		if err != nil {
			return diag.Errorf("error updating black white list: %s", err)
		}
	}
	return resourceBlackWhiteListRead(ctx, d, meta)
}

func buildUpdateBlackWhiteListBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"list_type": utils.ValueIgnoreEmpty(d.Get("list_type")),
		// direction can be 0
		"direction": d.Get("direction"),
		"protocol":  utils.ValueIgnoreEmpty(d.Get("protocol")),
		"port":      utils.ValueIgnoreEmpty(d.Get("port")),
		// address_type can be 0
		"address_type": d.Get("address_type"),
		"address":      utils.ValueIgnoreEmpty(d.Get("address")),
		"description":  d.Get("description"),
	}
	return bodyParams
}

func resourceBlackWhiteListDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteBlackWhiteList: Delete an existing CFW black white list
	var (
		deleteBlackWhiteListHttpUrl = "v1/{project_id}/black-white-list/{id}"
		deleteBlackWhiteListProduct = "cfw"
	)
	deleteBlackWhiteListClient, err := cfg.NewServiceClient(deleteBlackWhiteListProduct, region)
	if err != nil {
		return diag.Errorf("error creating CFW Client: %s", err)
	}

	deleteBlackWhiteListPath := deleteBlackWhiteListClient.Endpoint + deleteBlackWhiteListHttpUrl
	deleteBlackWhiteListPath = strings.ReplaceAll(deleteBlackWhiteListPath, "{project_id}",
		deleteBlackWhiteListClient.ProjectID)
	deleteBlackWhiteListPath = strings.ReplaceAll(deleteBlackWhiteListPath, "{id}", d.Id())

	deleteBlackWhiteListOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	_, err = deleteBlackWhiteListClient.Request("DELETE", deleteBlackWhiteListPath, &deleteBlackWhiteListOpt)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", "CFW.00200005"),
			"error deleting black white list",
		)
	}

	return nil
}

func resourceBlackWhiteListImportState(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	parts := strings.SplitN(d.Id(), "/", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid format specified for import id, must be <object_id>/<list_type>/<address>")
	}

	d.Set("object_id", parts[0])

	listType, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("error converting list_type: %s", err)
	}

	d.Set("list_type", listType)
	d.Set("address", parts[2])

	return []*schema.ResourceData{d}, nil
}
