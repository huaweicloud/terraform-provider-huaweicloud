package das

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DAS GET /v3/{project_id}/connections/{connection_id}/get-shared-list
func DataSourceSharedConnections() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSharedConnectionsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the shared connections are located.`,
			},

			// Required parameters.
			"connection_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the connection to which the shared connection belongs.`,
			},

			// Optional parameters.
			"keywords": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Keywords to search for shared connections.",
			},

			// Attributes.
			"shared_connections": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        sharedConnectionsElem(),
				Description: `The list of shared connections that matched filter parameters.`,
			},
		},
	}
}

func sharedConnectionsElem() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"user_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The user ID of the shared connection.`,
			},
			"user_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The user name of the shared connection.",
			},
			"expired_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The expiration time of the shared connection, in RFC3339 format.",
			},
			"shared_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The creation time of the shared connection, in RFC3339 format.",
			},
		},
	}
	return &sc
}

func buildSharedConnectionsQueryParams(d *schema.ResourceData) string {
	res := ""

	if v, ok := d.GetOk("keywords"); ok {
		res = fmt.Sprintf("%s&keywords=%v", res, v)
	}

	return res
}

func flattenSharedConnections(sharedConnections []interface{}) []map[string]interface{} {
	if len(sharedConnections) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(sharedConnections))
	for _, sharedConnection := range sharedConnections {
		result = append(result, map[string]interface{}{
			"user_id":   utils.PathSearch("user_id", sharedConnection, nil),
			"user_name": utils.PathSearch("user_name", sharedConnection, nil),
			"shared_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("shared_time", sharedConnection, float64(0)).(float64))/1000, false),
			"expired_at": utils.FormatTimeStampRFC3339(
				int64(utils.PathSearch("expired_time", sharedConnection, float64(0)).(float64))/1000, false),
		})
	}
	return result
}

func dataSourceSharedConnectionsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("das", region)
	if err != nil {
		return diag.Errorf("error creating DAS client: %s", err)
	}

	sharedConnections, err := listSharedConnections(client, d.Get("connection_id").(string), buildSharedConnectionsQueryParams(d))
	if err != nil {
		return diag.Errorf("error querying DAS shared connections: %s", err)
	}

	randUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("shared_connections", flattenSharedConnections(sharedConnections)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}
