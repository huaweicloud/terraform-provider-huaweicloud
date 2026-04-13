package cfw

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var batchDeleteAddressGroupMemberNonUpdatableParams = []string{
	"set_id",
	"address_item_ids",
	"fw_instance_id",
	"enterprise_project_id",
}

// @API CFW DELETE /v1/{project_id}/address-items
func ResourceBatchDeleteAddressGroupMember() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceBatchDeleteAddressGroupMemberCreate,
		ReadContext:   resourceBatchDeleteAddressGroupMemberRead,
		UpdateContext: resourceBatchDeleteAddressGroupMemberUpdate,
		DeleteContext: resourceBatchDeleteAddressGroupMemberDelete,

		CustomizeDiff: config.FlexibleForceNew(batchDeleteAddressGroupMemberNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"set_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"address_item_ids": {
				Type:     schema.TypeList,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"fw_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func buildBatchDeleteAddressGroupMemberQueryParams(d *schema.ResourceData, epsId string) string {
	queryParams := ""

	if epsId != "" {
		queryParams = fmt.Sprintf("%s&enterprise_project_id=%v", queryParams, epsId)
	}
	if v, ok := d.GetOk("fw_instance_id"); ok {
		queryParams = fmt.Sprintf("%s&fw_instance_id=%v", queryParams, v)
	}
	if queryParams != "" {
		queryParams = "?" + queryParams[1:]
	}

	return queryParams
}

func buildBatchDeleteAddressGroupMemberBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"set_id":           d.Get("set_id"),
		"address_item_ids": utils.ExpandToStringList(d.Get("address_item_ids").([]interface{})),
	}

	return bodyParams
}

func resourceBatchDeleteAddressGroupMemberCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		setId   = d.Get("set_id").(string)
		epsId   = cfg.GetEnterpriseProjectID(d)
		httpUrl = "v1/{project_id}/address-items"
	)

	client, err := cfg.NewServiceClient("cfw", region)
	if err != nil {
		return diag.Errorf("error creating CFW client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildBatchDeleteAddressGroupMemberQueryParams(d, epsId)
	requestOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         buildBatchDeleteAddressGroupMemberBodyParams(d),
	}

	_, err = client.Request("DELETE", requestPath, &requestOpt)
	if err != nil {
		return diag.Errorf("error batch deleting CFW address group members: %s", err)
	}

	d.SetId(setId)

	return nil
}

func resourceBatchDeleteAddressGroupMemberRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Read()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteAddressGroupMemberUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	// No processing is performed in 'Update()' method because resource is a one-time action resource.
	return nil
}

func resourceBatchDeleteAddressGroupMemberDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is a one-time action resource used to batch delete address group members. Deleting this resource
    will not clear of corresponding request record, but will only remove resource information from
    the tf state file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
