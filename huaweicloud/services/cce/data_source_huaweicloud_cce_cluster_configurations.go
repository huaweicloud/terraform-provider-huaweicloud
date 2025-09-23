package cce

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

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/configuration/detail
func DataSourceClusterConfigurations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterConfigurationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"configurations": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceClusterConfigurationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		getConfigurationsHttpUrl = "api/v3/projects/{project_id}/clusters/{cluster_id}/configuration/detail"
		getConfigurationsProduct = "cce"
	)
	getConfigurationsClient, err := cfg.NewServiceClient(getConfigurationsProduct, region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getConfigurationsHttpPath := getConfigurationsClient.Endpoint + getConfigurationsHttpUrl
	getConfigurationsHttpPath = strings.ReplaceAll(getConfigurationsHttpPath, "{project_id}", getConfigurationsClient.ProjectID)
	getConfigurationsHttpPath = strings.ReplaceAll(getConfigurationsHttpPath, "{cluster_id}", d.Get("cluster_id").(string))

	getConfigurationsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	getConfigurationsResp, err := getConfigurationsClient.Request("GET", getConfigurationsHttpPath, &getConfigurationsOpt)
	if err != nil {
		return diag.Errorf("error retrieving CCE configurations: %s", err)
	}

	getConfigurationsRespBody, err := utils.FlattenResponse(getConfigurationsResp)
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
		d.Set("configurations", flattenConfigurations(getConfigurationsRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigurations(resp interface{}) map[string]interface{} {
	configurations := resp.(map[string]interface{})
	if len(configurations) == 0 {
		return nil
	}

	res := make(map[string]interface{}, len(configurations))
	for k, v := range configurations {
		res[k] = utils.JsonToString(v)
	}

	return res
}
