package cfw

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CFW POST /v1/{project_id}/address-sets/batch-delete
func ResourceBatchDeleteAddressGroups() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteAddressGroupsCreate,
		ReadContext:   resourceBatchDeleteAddressGroupsRead,
		UpdateContext: resourceBatchDeleteAddressGroupsUpdate,
		DeleteContext: resourceBatchDeleteAddressGroupsDelete,

		CustomizeDiff: config.FlexibleForceNew([]string{
			"object_id",
			"set_ids",
		}),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"object_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"set_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func buildBatchDeleteAddressGroupsBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"object_id": d.Get("object_id").(string),
		"set_ids":   utils.ExpandToStringList(d.Get("set_ids").([]interface{})),
	}

	return bodyParams
}

func resourceBatchDeleteAddressGroupsCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg      = meta.(*config.Config)
		region   = cfg.GetRegion(d)
		httpUrl  = "v1/{project_id}/address-sets/batch-delete"
		objectId = d.Get("object_id").(string)
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchDeleteAddressGroupsBodyParams(d),
	}

	resp, err := client.Request("POST", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting CFW address groups: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(objectId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("data", flattenBatchDeleteAddressGroupsDataResp(respBody)))
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenBatchDeleteAddressGroupsDataResp(respBody interface{}) []map[string]interface{} {
	responseDataResp := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
	if len(responseDataResp) == 0 {
		return nil
	}

	responseData := make([]map[string]interface{}, 0, len(responseDataResp))
	for _, v := range responseDataResp {
		item := map[string]interface{}{
			"name": utils.PathSearch("name", v, nil),
			"id":   utils.PathSearch("id", v, nil),
		}
		responseData = append(responseData, item)
	}

	return responseData
}

func resourceBatchDeleteAddressGroupsRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteAddressGroupsUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteAddressGroupsDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete address groups. Deleting this resource
    will not clear of corresponding request record, but will only remove resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
