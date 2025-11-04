package dew

import (
	"context"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW POST /v1.0/{project_id}/kms/list-retirable-grants
func DataSourceRetirableGrantsGrants() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRetirableGrantsGrantsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"sequence": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"grants": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"key_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grant_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grantee_principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"grantee_principal_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"operations": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"issuing_principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"retiring_principal": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"creation_date": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceRetirableGrantsGrantsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg               = meta.(*config.Config)
		region            = cfg.GetRegion(d)
		listGrantsHttpUrl = "v1.0/{project_id}/kms/list-retirable-grants"
		allGrants         = make([]interface{}, 0)
		requestBody       = utils.RemoveNil(buildListRetirableGrantsGrantsBody(d))
		nextMarker        string
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	listGrantsPath := client.Endpoint + listGrantsHttpUrl
	listGrantsPath = strings.ReplaceAll(listGrantsPath, "{project_id}", client.ProjectID)
	listGrantOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody:         requestBody,
	}

	for {
		listGrantResp, err := client.Request("POST", listGrantsPath, &listGrantOpt)
		if err != nil {
			return diag.Errorf("error retrieving retirable grants: %s", err)
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

		requestBody["marker"] = nextMarker
	}

	datasourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(datasourceId)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("grants", flattenRetirableGrantsGrants(allGrants)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildListRetirableGrantsGrantsBody(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"sequence": utils.ValueIgnoreEmpty(d.Get("sequence")),
		"limit":    100,
	}

	return bodyParams
}

func flattenRetirableGrantsGrants(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	rst := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		rst = append(rst, map[string]interface{}{
			"key_id":                 utils.PathSearch("key_id", v, nil),
			"grant_id":               utils.PathSearch("grant_id", v, nil),
			"name":                   utils.PathSearch("name", v, nil),
			"grantee_principal":      utils.PathSearch("grantee_principal", v, nil),
			"grantee_principal_type": utils.PathSearch("grantee_principal_type", v, nil),
			"operations":             utils.PathSearch("operations", v, nil),
			"issuing_principal":      utils.PathSearch("issuing_principal", v, nil),
			"retiring_principal":     utils.PathSearch("retiring_principal", v, nil),
			"creation_date":          utils.PathSearch("creation_date", v, ""),
		})
	}
	return rst
}
