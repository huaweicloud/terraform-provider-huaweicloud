package taurusdb

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

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/htap/process
func DataSourceTaurusDBHtapSessions() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBHtapSessionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the resource.",
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `Specifies the HTAP instance ID.`,
			},
			"process_list": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The list of sessions of the HTAP instance.",
				Elem:        htapSessionsSchema(),
			},
		},
	}
}

func htapSessionsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"database": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sql_statement": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"duration": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"command": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBHtapSessionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/instances/{instance_id}/htap/process"
		product = "gaussdb"
		result  = make([]interface{}, 0)
		limit   = 100
		offset  = 0
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating TaurusDB Client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProviderClient.ProjectID)
	requestPath = strings.ReplaceAll(requestPath, "{instance_id}", d.Get("instance_id").(string))

	listOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	for {
		currentPath := fmt.Sprintf("%s?limit=%d&offset=%v", requestPath, limit, offset)
		resp, err := client.Request("GET", currentPath, &listOpts)
		if err != nil {
			return diag.Errorf("error retrieving TaurusDB HTAP instance sessions: %s", err)
		}

		respBody, err := utils.FlattenResponse(resp)
		if err != nil {
			return diag.FromErr(err)
		}

		processList := utils.PathSearch("process_list", respBody, make([]interface{}, 0)).([]interface{})
		if len(processList) == 0 {
			break
		}
		result = append(result, processList...)

		if len(processList) < limit {
			break
		}

		offset += len(processList)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("process_list", flattenHtapSessionsResponseBody(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHtapSessionsResponseBody(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	res := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		res = append(res, map[string]interface{}{
			"id":            utils.PathSearch("id", v, nil),
			"user":          utils.PathSearch("user", v, nil),
			"host":          utils.PathSearch("host", v, nil),
			"state":         utils.PathSearch("state", v, nil),
			"database":      utils.PathSearch("database", v, nil),
			"sql_statement": utils.PathSearch("sql_statement", v, nil),
			"duration":      utils.PathSearch("duration", v, nil),
			"command":       utils.PathSearch("command", v, nil),
		})
	}
	return res
}
