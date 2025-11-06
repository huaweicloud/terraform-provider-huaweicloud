package dew

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DEW GET /v1.0/{project_id}/keystores
func DataSourceDedicatedKeystores() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDedicatedKeystoresRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"keystores": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"keystore_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"domain_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keystore_alias": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"keystore_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hsm_cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cluster_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"create_time": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceDedicatedKeystoresRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v1.0/{project_id}/keystores?limit=10"
		offset  = 0
		result  = make([]interface{}, 0)
	)

	client, err := cfg.NewServiceClient("kms", region)
	if err != nil {
		return diag.Errorf("error creating KMS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	for {
		currentPath := fmt.Sprintf("%s&offset=%v", getPath, offset)
		getResp, err := client.Request("GET", currentPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving dedicated keystores: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}

		keystores := utils.PathSearch("keystores", getRespBody, make([]interface{}, 0)).([]interface{})
		if len(keystores) == 0 {
			break
		}

		result = append(result, keystores...)
		offset += len(keystores)
	}

	generateUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID)

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("keystores", flattenDedicatedKeystores(result)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenDedicatedKeystores(dataResp []interface{}) []interface{} {
	if len(dataResp) == 0 {
		return nil
	}

	result := make([]interface{}, 0, len(dataResp))
	for _, v := range dataResp {
		result = append(result, map[string]interface{}{
			"keystore_id":    utils.PathSearch("keystore_id", v, nil),
			"domain_id":      utils.PathSearch("domain_id", v, nil),
			"keystore_alias": utils.PathSearch("keystore_alias", v, nil),
			"keystore_type":  utils.PathSearch("keystore_type", v, nil),
			"hsm_cluster_id": utils.PathSearch("hsm_cluster_id", v, nil),
			"cluster_id":     utils.PathSearch("cluster_id", v, nil),
			"create_time":    utils.PathSearch("create_time", v, nil),
		})
	}

	return result
}
