package coc

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var groupResourceRelationNonUpdatableParams = []string{"group_id", "cmdb_resource_id"}

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
			"cmdb_resource_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceGroupResourceRelationCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)

	client, err := cfg.NewServiceClient("coc", cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating COC client: %s", err)
	}

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

	id := utils.PathSearch("data[0]", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("unable to find the COC group resource relation ID from the API response")
	}

	d.SetId(id)

	return nil
}

func buildCreateGroupResourceRelationBodyParams(d *schema.ResourceData) map[string]interface{} {
	cmdbResourceIdList := make([]interface{}, 0)
	cmdbResourceIdList = append(cmdbResourceIdList, d.Get("cmdb_resource_id"))

	bodyParams := map[string]interface{}{
		"group_id":              d.Get("group_id"),
		"cmdb_resource_id_list": cmdbResourceIdList,
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

	deleteHttpUrl := "v1/group-resource-relations?id_list={id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{id}", d.Id())
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
