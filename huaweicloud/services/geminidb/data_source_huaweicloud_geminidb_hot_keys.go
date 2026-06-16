package geminidb

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

// @API GeminiDB GET /v3/{project_id}/instances/{instance_id}/hot-keys
func DataSourceHotKeys() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceHotKeysRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"keys": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"command": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"qps": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"db_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func listHotKeys(client *golangsdk.ServiceClient, d *schema.ResourceData) ([]interface{}, error) {
	var (
		httpUrl    = "v3/{project_id}/instances/{instance_id}/hot-keys?limit={limit}"
		instanceId = d.Get("instance_id").(string)
		limit      = 50
		offset     = 0
		result     = make([]interface{}, 0)
	)

	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", instanceId)
	listPath = strings.ReplaceAll(listPath, "{limit}", strconv.Itoa(limit))
	listOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", listPath, offset)
		listResp, err := client.Request("GET", currentPath, &listOpt)
		if err != nil {
			return nil, err
		}

		resp, err := utils.FlattenResponse(listResp)
		if err != nil {
			return nil, err
		}

		respBody := utils.PathSearch("rules", resp, make([]interface{}, 0)).([]interface{})
		result = append(result, respBody...)
		if len(respBody) < limit {
			break
		}

		offset += len(respBody)
	}

	return result, nil
}

func dataSourceHotKeysRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)

	client, err := cfg.NewServiceClient("geminidb", region)
	if err != nil {
		return diag.Errorf("error creating GeminiDB client: %s", err)
	}

	resp, err := listHotKeys(client, d)
	if err != nil {
		return diag.Errorf("error retrieving GeminiDB Redis instance hot keys: %s", err)
	}

	randomUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(randomUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("keys", flattenHotKeys(resp)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenHotKeys(resp []interface{}) []interface{} {
	if len(resp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(resp))
	for _, v := range resp {
		result = append(result, map[string]interface{}{
			"name":    utils.PathSearch("name", v, nil),
			"type":    utils.PathSearch("type", v, nil),
			"command": utils.PathSearch("command", v, nil),
			"qps":     utils.PathSearch("qps", v, nil),
			"db_id":   utils.PathSearch("db_id", v, nil),
		})
	}

	return result
}
