package dli

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

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DLI GET /v3/{project_id}/datasource/auth-infos
func DataSourceAuths() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAuthsRead,

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
			"auths": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     authsSchema(),
			},
		},
	}
}

func authsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"username": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"certificate_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"truststore_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keystore_location": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"keytab": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"krb5_conf": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"owner": {
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
		},
	}
}

func buildListAuthsQueryParams(d *schema.ResourceData) string {
	queryParam := ""
	if v, ok := d.GetOk("name"); ok {
		queryParam += fmt.Sprintf("&auth_info_name=%v", v)
	}

	if queryParam != "" {
		queryParam = "?" + queryParam
	}

	return queryParam
}

func dataSourceAuthsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v3/{project_id}/datasource/auth-infos"
	)
	client, err := cfg.NewServiceClient("dli", region)
	if err != nil {
		return diag.Errorf("error creating DLI v3 client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)

	listAuthsParams := buildListAuthsQueryParams(d)
	getPath += listAuthsParams

	resp, err := pagination.ListAllItems(
		client,
		"offset",
		getPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving DLI datasource authentications")
	}

	listRespJson, err := json.Marshal(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID)

	mErr := multierror.Append(
		nil,
		d.Set("region", region),
		d.Set("auths", flattenListAuths(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenListAuths(resp interface{}) []interface{} {
	if resp == nil {
		return nil
	}

	curJson := utils.PathSearch("auth_infos", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	rst := make([]interface{}, len(curArray))
	for i, v := range curArray {
		rst[i] = map[string]interface{}{
			"name":                 utils.PathSearch("auth_info_name", v, nil),
			"type":                 utils.PathSearch("datasource_type", v, nil),
			"username":             utils.PathSearch("user_name", v, nil),
			"certificate_location": utils.PathSearch("certificate_location", v, nil),
			"truststore_location":  utils.PathSearch("truststore_location", v, nil),
			"keystore_location":    utils.PathSearch("keystore_location", v, nil),
			"keytab":               utils.PathSearch("keytab", v, nil),
			"krb5_conf":            utils.PathSearch("krb5_conf", v, nil),
			"owner":                utils.PathSearch("owner", v, nil),
			"created_at":           utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_time", v, float64(0)).(float64))/1000, false),
			"updated_at":           utils.FormatTimeStampRFC3339(int64(utils.PathSearch("update_time", v, float64(0)).(float64))/1000, false),
		}
	}
	return rst
}
