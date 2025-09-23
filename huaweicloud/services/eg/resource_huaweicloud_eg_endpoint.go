// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product EG
// ---------------------------------------------------------------

package eg

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API EG POST /v1/{project_id}/endpoints
// @API EG GET /v1/{project_id}/endpoints
// @API EG PUT /v1/{project_id}/endpoints/{endpoint_id}
// @API EG DELETE /v1/{project_id}/endpoints/{endpoint_id}
func ResourceEndpoint() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointCreate,
		UpdateContext: resourceEndpointUpdate,
		ReadContext:   resourceEndpointRead,
		DeleteContext: resourceEndpointDelete,
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
				ForceNew: true,
				Description: `Specifies the name of the endpoint.
`,
			},
			"vpc_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the VPC to which the endpoint belongs.`,
			},
			"subnet_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `Specifies the ID of the subnet to which the endpoint belongs.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the description of the endpoint.`,
			},
			"domain": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The domain of the endpoint.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the endpoint.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the endpoint.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The last update time of the endpoint.`,
			},
		},
	}
}

func resourceEndpointCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// createEndpoint: Create an EG Endpoint.
	var (
		createEndpointHttpUrl = "v1/{project_id}/endpoints"
		createEndpointProduct = "eg"
	)
	createEndpointClient, err := cfg.NewServiceClient(createEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	createEndpointPath := createEndpointClient.Endpoint + createEndpointHttpUrl
	createEndpointPath = strings.ReplaceAll(createEndpointPath, "{project_id}", createEndpointClient.ProjectID)

	createEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	createEndpointOpt.JSONBody = utils.RemoveNil(buildCreateEndpointBodyParams(d))
	createEndpointResp, err := createEndpointClient.Request("POST", createEndpointPath, &createEndpointOpt)
	if err != nil {
		return diag.Errorf("error creating Endpoint: %s", err)
	}

	createEndpointRespBody, err := utils.FlattenResponse(createEndpointResp)
	if err != nil {
		return diag.FromErr(err)
	}

	endpointId := utils.PathSearch("id", createEndpointRespBody, "").(string)
	if endpointId == "" {
		return diag.Errorf("unable to find the EG endpoint ID from the API response")
	}
	d.SetId(endpointId)

	return resourceEndpointRead(ctx, d, meta)
}

func buildCreateEndpointBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"name":        d.Get("name"),
		"vpc_id":      d.Get("vpc_id"),
		"subnet_id":   d.Get("subnet_id"),
		"description": utils.ValueIgnoreEmpty(d.Get("description")),
	}
	return bodyParams
}

func resourceEndpointRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getEndpoint: Query the EG Endpoint detail
	var (
		getEndpointHttpUrl = "v1/{project_id}/endpoints"
		getEndpointProduct = "eg"
	)
	getEndpointClient, err := cfg.NewServiceClient(getEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	getEndpointPath := getEndpointClient.Endpoint + getEndpointHttpUrl
	getEndpointPath = strings.ReplaceAll(getEndpointPath, "{project_id}", getEndpointClient.ProjectID)

	getEndpointqueryParams := BuildGetEndpointQueryParams(d.Get("name").(string))
	getEndpointPath += getEndpointqueryParams

	getEndpointResp, err := pagination.ListAllItems(
		getEndpointClient,
		"offset",
		getEndpointPath,
		&pagination.QueryOpts{MarkerField: ""})

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving Endpoint")
	}

	getEndpointRespJson, err := json.Marshal(getEndpointResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var getEndpointRespBody interface{}
	err = json.Unmarshal(getEndpointRespJson, &getEndpointRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	jsonPath := fmt.Sprintf("items[?id =='%s']|[0]", d.Id())
	getEndpointRespBody = utils.PathSearch(jsonPath, getEndpointRespBody, nil)
	if getEndpointRespBody == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "no data found")
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("name", getEndpointRespBody, nil)),
		d.Set("vpc_id", utils.PathSearch("vpc_id", getEndpointRespBody, nil)),
		d.Set("subnet_id", utils.PathSearch("subnet_id", getEndpointRespBody, nil)),
		d.Set("description", utils.PathSearch("description", getEndpointRespBody, nil)),
		d.Set("domain", utils.PathSearch("domain", getEndpointRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getEndpointRespBody, nil)),
		d.Set("created_at", utils.PathSearch("created_time", getEndpointRespBody, nil)),
		d.Set("updated_at", utils.PathSearch("updated_time", getEndpointRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func BuildGetEndpointQueryParams(name string) string {
	res := "?limit=100"

	if name != "" {
		res += fmt.Sprintf("&name=%s", name)
	}

	return res
}

func resourceEndpointUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	updateEndpointChanges := []string{
		"description",
	}

	if d.HasChanges(updateEndpointChanges...) {
		// updateEndpoint: Update the EG Endpoint.
		var (
			updateEndpointHttpUrl = "v1/{project_id}/endpoints/{id}"
			updateEndpointProduct = "eg"
		)
		updateEndpointClient, err := cfg.NewServiceClient(updateEndpointProduct, region)
		if err != nil {
			return diag.Errorf("error creating EG client: %s", err)
		}

		updateEndpointPath := updateEndpointClient.Endpoint + updateEndpointHttpUrl
		updateEndpointPath = strings.ReplaceAll(updateEndpointPath, "{project_id}", updateEndpointClient.ProjectID)
		updateEndpointPath = strings.ReplaceAll(updateEndpointPath, "{id}", d.Id())

		updateEndpointOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
			MoreHeaders: map[string]string{"Content-Type": "application/json"},
		}

		updateEndpointOpt.JSONBody = utils.RemoveNil(buildUpdateEndpointBodyParams(d))
		_, err = updateEndpointClient.Request("PUT", updateEndpointPath, &updateEndpointOpt)
		if err != nil {
			return diag.Errorf("error updating Endpoint: %s", err)
		}
	}
	return resourceEndpointRead(ctx, d, meta)
}

func buildUpdateEndpointBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"description": d.Get("description"),
	}
	return bodyParams
}

func resourceEndpointDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	// deleteEndpoint: Delete an existing EG Endpoint
	var (
		deleteEndpointHttpUrl = "v1/{project_id}/endpoints/{id}"
		deleteEndpointProduct = "eg"
	)
	deleteEndpointClient, err := cfg.NewServiceClient(deleteEndpointProduct, region)
	if err != nil {
		return diag.Errorf("error creating EG client: %s", err)
	}

	deleteEndpointPath := deleteEndpointClient.Endpoint + deleteEndpointHttpUrl
	deleteEndpointPath = strings.ReplaceAll(deleteEndpointPath, "{project_id}", deleteEndpointClient.ProjectID)
	deleteEndpointPath = strings.ReplaceAll(deleteEndpointPath, "{id}", d.Id())

	deleteEndpointOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
		MoreHeaders: map[string]string{"Content-Type": "application/json"},
	}

	_, err = deleteEndpointClient.Request("DELETE", deleteEndpointPath, &deleteEndpointOpt)
	if err != nil {
		return diag.Errorf("error deleting Endpoint: %s", err)
	}

	return nil
}
