package cce

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/tidwall/gjson"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/httphelper"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/helper/schemas"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/clusters/configuration/detail
func DataSourceClusterConfigurationDetails() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterConfigurationDetailsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"network_mode": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"configurations": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The configurations of the cce cluster.",
			},
		},
	}
}

type ClusterConfigurationDetailsDSWrapper struct {
	*schemas.ResourceDataWrapper
	Config *config.Config
}

func newClusterConfigurationDetailsDSWrapper(d *schema.ResourceData, meta interface{}) *ClusterConfigurationDetailsDSWrapper {
	return &ClusterConfigurationDetailsDSWrapper{
		ResourceDataWrapper: schemas.NewSchemaWrapper(d),
		Config:              meta.(*config.Config),
	}
}

func dataSourceClusterConfigurationDetailsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	wrapper := newClusterConfigurationDetailsDSWrapper(d, meta)
	clusterConfigurationDetailsRst, err := wrapper.ShowClusterConfigurationDetails()
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.GenerateUUID()
	if err != nil {
		return diag.FromErr(err)
	}
	d.SetId(id)

	err = wrapper.showClusterConfigurationDetailsToSchema(clusterConfigurationDetailsRst)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

// @API CCE GET /api/v3/clusters/configuration/detail
func (w *ClusterConfigurationDetailsDSWrapper) ShowClusterConfigurationDetails() (*gjson.Result, error) {
	client, err := w.NewClient(w.Config, "cce")
	if err != nil {
		return nil, err
	}

	uri := "/api/v3/clusters/configuration/detail"
	params := map[string]any{
		"clusterType":    w.Get("cluster_type"),
		"clusterVersion": w.Get("cluster_version"),
		"clusterID":      w.Get("cluster_id"),
		"networkMode":    w.Get("network_mode"),
	}
	params = utils.RemoveNil(params)
	return httphelper.New(client).
		Method("GET").
		URI(uri).
		Query(params).
		Request().
		Result()
}

func (w *ClusterConfigurationDetailsDSWrapper) showClusterConfigurationDetailsToSchema(body *gjson.Result) error {
	d := w.ResourceData
	mErr := multierror.Append(nil,
		d.Set("region", w.Config.GetRegion(w.ResourceData)),
		d.Set("configurations", w.setConfigurationsValues(body)),
	)
	return mErr.ErrorOrNil()
}

func (*ClusterConfigurationDetailsDSWrapper) setConfigurationsValues(data *gjson.Result) string {
	configurationsValues := data.Get("configurations")
	jsonBytes, _ := json.Marshal(configurationsValues)

	return string(jsonBytes)
}
