// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product IdentityCenter
// ---------------------------------------------------------------

package identitycenter

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API IdentityCenter GET /v1/identity-stores/{identity_store_id}/groups
func DataSourceIdentityCenterGroups() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceIdentityCenterGroupsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"identity_store_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the ID of the identity store that associated with IAM Identity Center.`,
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `Specifies the name of the group.`,
			},
			"groups": {
				Type:        schema.TypeList,
				Elem:        identityCenterGroupsGroupSchema(),
				Computed:    true,
				Description: `Indicates the list of IdentityCenter group.`,
			},
		},
	}
}

func identityCenterGroupsGroupSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the ID of the group.`,
			},
			"name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the name of the group.`,
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the description of the group.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the creation time.`,
			},
			"updated_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Indicates the last update time.`,
			},
		},
	}
	return &sc
}

func resourceIdentityCenterGroupsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	// getIdentityCenterGroups: Query Identity Center groups
	var (
		getIdentityCenterGroupsHttpUrl = "v1/identity-stores/{identity_store_id}/groups"
		getIdentityCenterGroupsProduct = "identitystore"
	)
	getIdentityCenterGroupsClient, err := cfg.NewServiceClient(getIdentityCenterGroupsProduct, region)
	if err != nil {
		return diag.Errorf("error creating Identity Center client: %s", err)
	}

	getIdentityCenterGroupsBasePath := getIdentityCenterGroupsClient.Endpoint + getIdentityCenterGroupsHttpUrl
	getIdentityCenterGroupsBasePath = strings.ReplaceAll(getIdentityCenterGroupsBasePath, "{identity_store_id}",
		d.Get("identity_store_id").(string))

	getIdentityCenterGroupsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	var groups []interface{}
	var marker string
	var getIdentityCenterGroupsPath string
	for {
		getIdentityCenterGroupsPath = getIdentityCenterGroupsBasePath + buildGetIdentityCenterGroupsQueryParams(d, marker)
		getIdentityCenterGroupsResp, err := getIdentityCenterGroupsClient.Request("GET",
			getIdentityCenterGroupsPath, &getIdentityCenterGroupsOpt)

		if err != nil {
			return common.CheckDeletedDiag(d, err, "error retrieving Identity Center groups")
		}

		getIdentityCenterGroupsRespBody, err := utils.FlattenResponse(getIdentityCenterGroupsResp)
		if err != nil {
			return diag.FromErr(err)
		}
		groups = append(groups, flattenGetIdentityCenterGroupsResponseBodyGroup(getIdentityCenterGroupsRespBody)...)
		marker = utils.PathSearch("page_info.next_marker", getIdentityCenterGroupsRespBody, "").(string)
		if marker == "" {
			break
		}
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("groups", groups),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenGetIdentityCenterGroupsResponseBodyGroup(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("groups", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		createAt := utils.PathSearch("created_at", v, float64(0)).(float64)
		updateAt := utils.PathSearch("updated_at", v, float64(0)).(float64)
		rst = append(rst, map[string]interface{}{
			"id":          utils.PathSearch("group_id", v, nil),
			"name":        utils.PathSearch("display_name", v, nil),
			"description": utils.PathSearch("description", v, nil),
			"created_at":  utils.FormatTimeStampRFC3339(int64(createAt)/1000, false),
			"updated_at":  utils.FormatTimeStampRFC3339(int64(updateAt)/1000, false),
		})
	}
	return rst
}

func buildGetIdentityCenterGroupsQueryParams(d *schema.ResourceData, marker string) string {
	res := "?limit=100"
	if v, ok := d.GetOk("name"); ok {
		res = fmt.Sprintf("%s&display_name=%v", res, v)
	}

	if marker != "" {
		res = fmt.Sprintf("%s&marker=%v", res, marker)
	}
	return res
}
