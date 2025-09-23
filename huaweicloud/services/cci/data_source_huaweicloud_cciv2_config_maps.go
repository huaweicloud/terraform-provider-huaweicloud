package cci

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

// @API CCI GET /apis/cci/v2/namespaces/{namespace}/configmaps
func DataSourceV2ConfigMaps() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceV2ConfigMapsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"config_maps": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"namespace": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"annotations": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"labels": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"creation_timestamp": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"resource_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"uid": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"binary_data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"data": {
							Type:     schema.TypeMap,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"immutable": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func dataSourceV2ConfigMapsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("cci", region)
	if err != nil {
		return diag.Errorf("error creating CCI client: %s", err)
	}
	listConfigMapsHttpUrl := "apis/cci/v2/namespaces/{namespace}/configmaps"
	listConfigMapsPath := client.Endpoint + listConfigMapsHttpUrl
	listConfigMapsPath = strings.ReplaceAll(listConfigMapsPath, "{namespace}", d.Get("namespace").(string))
	listConfigMapsOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	listConfigMapsResp, err := client.Request("GET", listConfigMapsPath, &listConfigMapsOpt)
	if err != nil {
		return diag.Errorf("error getting CCI configMaps list: %s", err)
	}

	listConfigMapsRespBody, err := utils.FlattenResponse(listConfigMapsResp)
	if err != nil {
		return diag.Errorf("error retrieving CCI configMaps: %s", err)
	}

	uuid, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(uuid)

	configMaps := utils.PathSearch("items", listConfigMapsRespBody, make([]interface{}, 0)).([]interface{})
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("config_maps", flattenConfigMaps(configMaps)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenConfigMaps(configMaps []interface{}) []interface{} {
	if len(configMaps) == 0 {
		return nil
	}

	rst := make([]interface{}, len(configMaps))
	for i, v := range configMaps {
		rst[i] = map[string]interface{}{
			"name":               utils.PathSearch("metadata.name", v, nil),
			"namespace":          utils.PathSearch("metadata.namespace", v, nil),
			"annotations":        utils.PathSearch("metadata.annotations", v, nil),
			"labels":             utils.PathSearch("metadata.labels", v, nil),
			"creation_timestamp": utils.PathSearch("metadata.creationTimestamp", v, nil),
			"resource_version":   utils.PathSearch("metadata.resourceVersion", v, nil),
			"uid":                utils.PathSearch("metadata.uid", v, nil),
			"binary_data":        utils.PathSearch("binaryData", v, nil),
			"data":               utils.PathSearch("data", v, nil),
			"immutable":          utils.PathSearch("immutable", v, nil),
		}
	}
	return rst
}
