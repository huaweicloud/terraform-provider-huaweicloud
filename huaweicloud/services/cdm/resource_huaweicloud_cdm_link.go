package cdm

import (
	"context"
	"fmt"
	"strings"

	"github.com/chnsz/golangsdk/openstack/cdm/v1/link"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

const (
	configPref = "linkConfig."
)

func ResourceCdmLink() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCdmLinkCreate,
		ReadContext:   resourceCdmLinkRead,
		UpdateContext: resourceCdmLinkUpdate,
		DeleteContext: resourceCdmLinkDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"connector": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"config": {
				Type:     schema.TypeMap,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},

			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},

			"access_key": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"secret_key"},
			},

			"secret_key": {
				Type:         schema.TypeString,
				Optional:     true,
				Sensitive:    true,
				RequiredWith: []string{"access_key"},
			},

			"enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  true,
			},
		},
	}
}

func resourceCdmLinkCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	linkConfigValues, err := buildLinkConfigParamter(d)
	if err != nil {
		return diag.FromErr(err)
	}

	opts := link.LinkCreateOpts{
		Links: []link.Link{
			{
				Name:             d.Get("name").(string),
				ConnectorName:    d.Get("connector").(string),
				Enabled:          utils.Bool(d.Get("enabled").(bool)),
				LinkConfigValues: *linkConfigValues,
			},
		},
	}

	clusterId := d.Get("cluster_id").(string)

	rst, createErr := link.Create(client, clusterId, opts)
	if createErr != nil {
		return fmtp.DiagErrorf("Error creating CDM link: %s", createErr)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, rst.Name))

	return resourceCdmLinkRead(ctx, d, meta)
}

func buildLinkConfigParamter(d *schema.ResourceData) (*link.LinkConfigs, error) {
	var configs []link.Input
	configRaw := d.Get("config").(map[string]interface{})

	if len(configRaw) < 1 {
		return nil, fmtp.Errorf("The config is Required.")
	}

	for k, v := range configRaw {
		config := link.Input{
			Name:  fmt.Sprintf("%s%s", configPref, k),
			Value: v.(string),
		}
		configs = append(configs, config)
	}

	connector := d.Get("connector").(string)

	if v, ok := d.GetOk("password"); ok {
		if connector == link.GenericJdbcConnector || connector == link.HdfsConnector ||
			connector == link.HbaseConnector || connector == link.SftpConnector ||
			connector == link.MongodbConnector || connector == link.ElasticsearchConnector {
			input := link.Input{
				Name:  fmt.Sprintf("%s%s", configPref, "password"),
				Value: v.(string),
			}
			configs = append(configs, input)
		}
	}

	if v, ok := d.GetOk("secret_key"); ok {
		if connector == link.ObsConnector || connector == link.ThirdpartyObsConnector ||
			connector == link.HbaseConnector {
			ak := link.Input{
				Name:  fmt.Sprintf("%s%s", configPref, "accessKey"),
				Value: d.Get("access_key").(string),
			}
			sk := link.Input{
				Name:  fmt.Sprintf("%s%s", configPref, "securityKey"),
				Value: v.(string),
			}
			configs = append(configs, ak, sk)
		} else if connector == link.DisConnector || connector == link.DliConnector ||
			connector == link.OpentsdbConnector || connector == link.DmsKafkaConnector {
			ak := link.Input{
				Name:  fmt.Sprintf("%s%s", configPref, "ak"),
				Value: d.Get("ak").(string),
			}
			sk := link.Input{
				Name:  fmt.Sprintf("%s%s", configPref, "sk"),
				Value: d.Get("sk").(string),
			}
			configs = append(configs, ak, sk)
		}
	}

	LinkConfigValues := link.LinkConfigs{
		Configs: []link.Configs{
			{
				Name:   "linkConfig",
				Inputs: configs,
			},
		},
	}

	return &LinkConfigValues, nil

}

func resourceCdmLinkRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	clusterId, linkName, err := ParseLinkInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	resp, gErr := link.Get(client, clusterId, linkName)
	if gErr != nil {
		return common.CheckDeletedDiag(d, err, "Error retrieving CDM link")
	}

	detail := resp.Links[0]
	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.Name),
		d.Set("cluster_id", clusterId),
		d.Set("connector", detail.ConnectorName),
		d.Set("enabled", detail.Enabled),
		setConfigToState(d, detail.LinkConfigValues.Configs),
	)

	if mErr.ErrorOrNil() != nil {
		return fmtp.DiagErrorf("Error setting CDM link fields: %s", mErr)
	}

	return nil
}

func setConfigToState(d *schema.ResourceData, configs []link.Configs) error {
	if len(configs) == 0 {
		return nil
	}

	result := make(map[string]string)
	for _, item := range configs {
		if item.Name == "linkConfig" {
			for _, v := range item.Inputs {
				if v.Value != "" {
					key := strings.Replace(v.Name, configPref, "", 1)

					if key == "password" {
						d.Set("password", v.Value)
					} else if key == "securityKey" || key == "sk" {
						d.Set("secret_key", v.Value)
					} else if key == "accessKey" || key == "ak" {
						d.Set("access_key", v.Value)
					} else {
						result[key] = v.Value
					}
				}
			}
		}
	}
	return d.Set("config", result)
}

func resourceCdmLinkUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	linkConfigValues, err := buildLinkConfigParamter(d)
	if err != nil {
		return diag.FromErr(err)
	}

	clusterId, linkName, err := ParseLinkInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	newName := d.Get("name").(string)

	opts := link.LinkCreateOpts{
		Links: []link.Link{
			{
				Name:             newName,
				ConnectorName:    d.Get("connector").(string),
				Enabled:          utils.Bool(d.Get("enabled").(bool)),
				LinkConfigValues: *linkConfigValues,
			},
		},
	}

	_, updateErr := link.Update(client, clusterId, linkName, opts)
	if updateErr != nil {
		return fmtp.DiagErrorf("Error update CDM link: %s", updateErr)
	}

	d.SetId(fmt.Sprintf("%s/%s", clusterId, newName))
	return resourceCdmLinkRead(ctx, d, meta)
}

func resourceCdmLinkDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	client, err := config.CdmV11Client(region)
	if err != nil {
		return diag.Errorf("error creating CDM v1 client, err=%s", err)
	}

	clusterId, linkName, err := ParseLinkInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	_, dErr := link.Delete(client, clusterId, linkName)
	if dErr != nil {
		return fmtp.DiagErrorf("delete CDM link failed. %q:%s", d.Id(), dErr)
	}

	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func ParseLinkInfoFromId(id string) (clusterId, linkName string, err error) {
	idArrays := strings.SplitN(id, "/", 2)
	if len(idArrays) != 2 {
		err = fmtp.Errorf("Invalid format specified for ID. Format must be <cluster_id>/<link_name>")
		return
	}

	clusterId = idArrays[0]
	linkName = idArrays[1]
	return
}
