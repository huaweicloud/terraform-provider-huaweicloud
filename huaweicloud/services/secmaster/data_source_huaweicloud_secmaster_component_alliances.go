package secmaster

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API SecMaster GET /v1/{project_id}/workspaces/{workspace_id}/soc/components/alliance
func DataSourceComponentAlliances() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceComponentAlliancesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"is_built_in": {
				Type:     schema.TypeBool,
				Required: true,
			},
			"data": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     componentAllianceSchema(),
			},
		},
	}
}

func componentAllianceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"alliance_code": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alliance_description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alliance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"alliance_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"logo": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceComponentAlliancesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		product   = "secmaster"
		httpUrl   = "v1/{project_id}/workspaces/{workspace_id}/soc/components/alliance"
		result    = make([]interface{}, 0)
		limit     = 1000
		offset    = 0
		isBuiltIn = d.Get("is_built_in").(bool)
		mErr      *multierror.Error
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating SecMaster client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{workspace_id}", d.Get("workspace_id").(string))

	reqOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"content-type": "application/json;charset=UTF-8",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s?offset=%d&limit=%d&is_built_in=%t", requestPath, offset, limit, isBuiltIn)

		resp, err := client.Request("GET", currentPath, &reqOpt)
		if err != nil {
			return diag.Errorf("error retrieving SecMaster component alliances: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		alliancesRaw := utils.PathSearch("data", respBody, make([]interface{}, 0)).([]interface{})
		if len(alliancesRaw) == 0 {
			break
		}

		result = append(result, alliancesRaw...)

		offset += len(alliancesRaw)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("data", flattenComponentAlliances(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenComponentAlliances(result []interface{}) []interface{} {
	if len(result) == 0 {
		return nil
	}

	alliances := make([]interface{}, 0, len(result))
	for _, item := range result {
		alliances = append(alliances, map[string]interface{}{
			"alliance_code":        utils.PathSearch("alliance_code", item, nil),
			"alliance_description": utils.PathSearch("alliance_description", item, nil),
			"alliance_name":        utils.PathSearch("alliance_name", item, nil),
			"alliance_type":        utils.PathSearch("alliance_type", item, nil),
			"logo":                 utils.PathSearch("logo", item, nil),
			"id":                   utils.PathSearch("id", item, nil),
			"create_time":          utils.PathSearch("create_time", item, nil),
			"creator_name":         utils.PathSearch("creator_name", item, nil),
			"update_time":          utils.PathSearch("update_time", item, nil),
		})
	}

	return alliances
}
