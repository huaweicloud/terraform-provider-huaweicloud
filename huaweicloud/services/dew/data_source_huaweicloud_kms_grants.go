package dew

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/list-grants
func DataSourceKmsGrants() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDewKmsGrantsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `Specifies the region in which to query the resource.`,
			},
			"key_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the key ID to which the grants belong.`,
			},
			"grants": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: `The list of the grants.`,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the grant.`,
						},
						"key_id": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The key ID to which the grant belongs.`,
						},
						"name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The name of the grant.`,
						},
						"type": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The authorization type.`,
						},
						"creator": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the user who created the grant.`,
						},
						"grantee_principal": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The ID of the authorized user or account.`,
						},
						"operations": {
							Type:        schema.TypeList,
							Computed:    true,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Description: `List of granted operations.`,
						},
						"created_at": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: `The creation time of the grant.`,
						},
					},
				},
			},
		},
	}
}

func dataSourceDewKmsGrantsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		listGrantsHttpUrl = "v1.0/{project_id}/kms/list-grants"
		listGrantsProduct = "kms"
	)
	listGrantsClient, err := cfg.NewServiceClient(listGrantsProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	listGrantsPath := listGrantsClient.Endpoint + listGrantsHttpUrl
	listGrantsPath = strings.ReplaceAll(listGrantsPath, "{project_id}", listGrantsClient.ProjectID)

	listGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	allGrants := make([]interface{}, 0)
	var nextMarker string

	listGrantOpt.JSONBody = utils.RemoveNil(buildListGrantsBody(d, cfg))
	listGrantJSONBody := listGrantOpt.JSONBody.(map[string]interface{})

	for {
		listGrantResp, err := listGrantsClient.Request("POST", listGrantsPath, &listGrantOpt)
		if err != nil {
			return diag.Errorf("error retrieving grants: %s", err)
		}
		listGrantRespBody, err := utils.FlattenResponse(listGrantResp)
		if err != nil {
			return diag.FromErr(err)
		}

		grants := utils.PathSearch("grants", listGrantRespBody, make([]interface{}, 0)).([]interface{})
		if len(grants) == 0 {
			return nil
		}
		allGrants = append(allGrants, grants...)

		nextMarker = utils.PathSearch("next_marker", listGrantRespBody, "").(string)
		if nextMarker == "" {
			break
		}
		listGrantJSONBody["marker"] = nextMarker
	}

	uuId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("grants", flattenListSnatRuleResponseBody(allGrants)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListGrantsBody(d *schema.ResourceData, _ *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"key_id": d.Get("key_id"),
		"limit":  1000,
	}
	return bodyParams
}

func flattenListSnatRuleResponseBody(resp []interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"id":                utils.PathSearch("grant_id", v, nil),
			"key_id":            utils.PathSearch("key_id", v, nil),
			"name":              utils.PathSearch("name", v, nil),
			"type":              utils.PathSearch("grantee_principal_type", v, nil),
			"creator":           utils.PathSearch("issuing_principal", v, nil),
			"grantee_principal": utils.PathSearch("grantee_principal", v, nil),
			"operations":        utils.PathSearch("operations", v, nil),
			"created_at":        utils.FormatTimeStampRFC3339(convertStrToInt(utils.PathSearch("creation_date", v, "").(string))/1000, false),
		})
	}
	return rst
}

func convertStrToInt(str string) int64 {
	resp, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("error convert the string (%s) to int", str)
		return 0
	}

	return int64(resp)
}
