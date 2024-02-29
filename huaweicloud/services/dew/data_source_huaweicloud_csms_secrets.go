package dew

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW GET /v1/{project_id}/secrets
func DataSourceDewCsmsSecrets() *schema.Resource {
	return &schema.Resource{
		ReadContext: resourceDewCsmsSecretsRead,
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"event_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secrets": {
				Type:     schema.TypeList,
				Elem:     secretsSchema(),
				Computed: true,
			},
		},
	}
}

func secretsSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"kms_key_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scheduled_deleted_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"secret_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"auto_rotation": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"rotation_period": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_config": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rotation_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"next_rotation_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"event_subscriptions": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
	return &sc
}

func resourceDewCsmsSecretsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		listSecretsHttpUrl = "v1/{project_id}/secrets"
		listSecretsProduct = "kms"
	)
	listSecretsClient, err := cfg.NewServiceClient(listSecretsProduct, region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	listSecretsPath := listSecretsClient.Endpoint + listSecretsHttpUrl
	listSecretsPath = strings.ReplaceAll(listSecretsPath, "{project_id}", listSecretsClient.ProjectID)

	listSecretsQueryParams := buildListSecretsQueryParams(d)
	listSecretsPath += listSecretsQueryParams

	listSecretsResp, err := pagination.ListAllItems(
		listSecretsClient,
		"marker",
		listSecretsPath,
		&pagination.QueryOpts{MarkerField: "name"})

	if err != nil {
		return diag.Errorf("error retrieving CSMS secrets: %s", err)
	}

	listSecretsRespJson, err := json.Marshal(listSecretsResp)
	if err != nil {
		return diag.FromErr(err)
	}
	var listSecretsRespBody interface{}
	err = json.Unmarshal(listSecretsRespJson, &listSecretsRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("secrets", filterListSecretsBody(
			flattenListSecretsBody(listSecretsRespBody), d)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListSecretsBody(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}
	curJson := utils.PathSearch("secrets", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		createAt := utils.PathSearch("create_time", v, float64(0))
		updateAt := utils.PathSearch("update_time", v, float64(0))
		scheduledDeletedAt := utils.PathSearch("scheduled_delete_time", v, float64(0))
		rotationAt := utils.PathSearch("rotation_time", v, float64(0))
		nextRotationAt := utils.PathSearch("next_rotation_time", v, float64(0))
		rst = append(rst, map[string]interface{}{
			"id":                    utils.PathSearch("id", v, nil),
			"name":                  utils.PathSearch("name", v, nil),
			"status":                utils.PathSearch("state", v, nil),
			"kms_key_id":            utils.PathSearch("kms_key_id", v, nil),
			"description":           utils.PathSearch("description", v, nil),
			"created_at":            utils.FormatTimeStampRFC3339(int64(createAt.(float64))/1000, false),
			"updated_at":            utils.FormatTimeStampRFC3339(int64(updateAt.(float64))/1000, false),
			"scheduled_deleted_at":  utils.FormatTimeStampRFC3339(int64(scheduledDeletedAt.(float64))/1000, false),
			"secret_type":           utils.PathSearch("secret_type", v, nil),
			"auto_rotation":         utils.PathSearch("auto_rotation", v, nil),
			"rotation_period":       utils.PathSearch("rotation_period", v, nil),
			"rotation_config":       utils.PathSearch("rotation_config", v, nil),
			"rotation_at":           utils.FormatTimeStampRFC3339(int64(rotationAt.(float64))/1000, false),
			"next_rotation_at":      utils.FormatTimeStampRFC3339(int64(nextRotationAt.(float64))/1000, false),
			"event_subscriptions":   utils.PathSearch("event_subscriptions", v, nil),
			"enterprise_project_id": utils.PathSearch("enterprise_project_id", v, nil),
		})
	}
	return rst
}

func filterListSecretsBody(all []interface{}, d *schema.ResourceData) []interface{} {
	rst := make([]interface{}, 0, len(all))

	for _, v := range all {
		if param, ok := d.GetOk("name"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("name", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("secret_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("id", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("status"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("status", v, nil)) {
			continue
		}
		if param, ok := d.GetOk("enterprise_project_id"); ok &&
			fmt.Sprint(param) != fmt.Sprint(utils.PathSearch("enterprise_project_id", v, nil)) {
			continue
		}

		rst = append(rst, v)
	}

	return rst
}

func buildListSecretsQueryParams(d *schema.ResourceData) string {
	res := ""
	if v, ok := d.GetOk("event_name"); ok {
		res = fmt.Sprintf("%s&event_name=%v", res, v)
	}

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
