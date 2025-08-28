package coc

import (
	"context"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var groupResourceRelationNonUpdatableParams = []string{"group_id", "cmdb_resource_id_list"}

// @API COC POST /v1/group-resource-relations
// @API COC DELETE /v1/group-resource-relations
func ResourceGroupResourceRelation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceGroupResourceRelationCreate,
		ReadContext:   resourceGroupResourceRelationRead,
		UpdateContext: resourceGroupResourceRelationUpdate,
		DeleteContext: resourceGroupResourceRelationDelete,

		CustomizeDiff: config.FlexibleForceNew(groupResourceRelationNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"group_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"cmdb_resource_id_list": {
				Type:     schema.TypeSet,
				Required: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"relation_id_list": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceGroupResourceRelationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	groupID := d.Get("group_id").(string)

	createHttpUrl := "v1/group-resource-relations"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: utils.RemoveNil(
			buildCreateGroupResourceRelationBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating COC group resource relation: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.FromErr(err)
	}

	relationIdList := utils.PathSearch("data", createRespBody, make([]interface{}, 0)).([]interface{})
	if len(relationIdList) < 1 {
		return diag.Errorf("unable to creating group resource relation ID from the API response")
	}

	d.SetId(groupID)

	mErr := multierror.Append(nil,
		d.Set("relation_id_list", relationIdList),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildCreateGroupResourceRelationBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"group_id":              d.Get("group_id"),
		"cmdb_resource_id_list": d.Get("cmdb_resource_id_list").(*schema.Set).List(),
	}

	return bodyParams
}

func resourceGroupResourceRelationRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGroupResourceRelationUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceGroupResourceRelationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

	if _, ok := d.GetOk("relation_id_list"); !ok {
		return diag.Errorf("error deleting COC group resource relation: The relationship ID list does not exist.")
	}

	deleteHttpUrl := "v1/group-resource-relations?" +
		buildQueryStringParams("id_list", d.Get("relation_id_list").([]interface{}))[1:]
	deletePath := client.Endpoint + deleteHttpUrl
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, common.ConvertExpected400ErrInto404Err(err, "error_code",
			"COC.00101031"), "error deleting COC group resource relation")
	}

	return nil
}
