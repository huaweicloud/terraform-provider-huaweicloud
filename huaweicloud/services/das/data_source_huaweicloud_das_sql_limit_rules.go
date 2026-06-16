package das

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/instances/{instance_id}/sql-limit/rules
func DataSourceSqlLimitRules() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSqlLimitRulesRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the SQL limit rules are located.`,
			},

			// Required parameters.
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the instance to which the SQL limit rules belong.`,
			},

			// Attributes.
			"rules": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        sqlLimitRulesElem(),
				Description: `The list of SQL limit rules that matched filter parameters.`,
			},
		},
	}
}

func sqlLimitRulesElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The ID of the SQL limit rule.`,
			},
			"sql_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the SQL.`,
			},
			"pattern": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The pattern of the SQL limit rule.`,
			},
			"max_concurrency": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum concurrency of the SQL limit rule.`,
			},
			"max_waiting": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The maximum waiting time of the SQL limit rule.`,
			},
		},
	}
	return &sc
}

func listSqlLimitRules(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/sql-limit/rules?datastore_type={datastore_type}&limit={limit}"
		limit   = 100
		offset  = 0
		result  = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath = strings.ReplaceAll(listPath, "{datastore_type}", "MySQL")
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		listPathWithOffset := fmt.Sprintf("%s&offset=%d", listPath, offset)
		requestResp, err := client.Request("GET", listPathWithOffset, &opt)
		if err != nil {
			return nil, err
		}

		respBody, err := utils.FlattenResponse(requestResp)
		if err != nil {
			return nil, err
		}

		rules := utils.PathSearch("sql_limit_rules", respBody, make([]interface{}, 0)).([]interface{})
		result = append(result, rules...)
		if len(rules) < limit {
			break
		}

		offset += len(rules)
	}

	return result, nil
}

func flattenSqlLimitRules(rules []interface{}) []map[string]interface{} {
	if len(rules) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(rules))
	for _, rule := range rules {
		result = append(result, map[string]interface{}{
			"id":              utils.PathSearch("id", rule, nil),
			"sql_type":        utils.PathSearch("sql_type", rule, nil),
			"pattern":         utils.PathSearch("pattern", rule, nil),
			"max_concurrency": utils.PathSearch("max_concurrency", rule, nil),
			"max_waiting":     utils.PathSearch("max_waiting", rule, nil),
		})
	}
	return result
}

func dataSourceSqlLimitRulesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	rules, err := listSqlLimitRules(client, d)
	if err != nil {
		return diag.Errorf("error querying DAS SQL limit rules: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("rules", flattenSqlLimitRules(rules)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
