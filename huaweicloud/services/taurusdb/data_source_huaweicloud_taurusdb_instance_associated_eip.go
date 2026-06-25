package taurusdb

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API TaurusDB GET /v3/{project_id}/instances/{instance_id}/eip
func DataSourceTaurusDBInstanceAssociatedEip() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTaurusDBInstanceAssociatedEipRead,

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
			"can_enable_public_access": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"eip_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"port_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip_address": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"bandwidth_size": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"bandwidth_share_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"profile": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     instanceEipProfileSchema(),
			},
		},
	}
}

func instanceEipProfileSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"order_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"product_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTaurusDBInstanceAssociatedEipRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		instanceId = d.Get("instance_id").(string)
		httpUrl    = "v3/{project_id}/instances/{instance_id}/eip"
	)
	client, err := cfg.NewServiceClient("gaussdb", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", instanceId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error querying TaurusDB instance (%s) EIP: %s", instanceId, err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("can_enable_public_access", utils.PathSearch("can_enable_public_access", getRespBody, nil)),
		d.Set("eip_id", utils.PathSearch("id", getRespBody, nil)),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("port_id", utils.PathSearch("port_id", getRespBody, nil)),
		d.Set("public_ip_address", utils.PathSearch("public_ip_address", getRespBody, nil)),
		d.Set("private_ip_address", utils.PathSearch("private_ip_address", getRespBody, nil)),
		d.Set("status", utils.PathSearch("status", getRespBody, nil)),
		d.Set("create_time", utils.PathSearch("create_time", getRespBody, nil)),
		d.Set("bandwidth_id", utils.PathSearch("bandwidth_id", getRespBody, nil)),
		d.Set("bandwidth_name", utils.PathSearch("bandwidth_name", getRespBody, nil)),
		d.Set("bandwidth_size", utils.PathSearch("bandwidth_size", getRespBody, float64(0)).(float64)),
		d.Set("bandwidth_share_type", utils.PathSearch("bandwidth_share_type", getRespBody, nil)),
		d.Set("profile", flattenInstanceEipProfile(getRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenInstanceEipProfile(resp interface{}) []interface{} {
	curJson := utils.PathSearch("profile", resp, nil)
	if curJson == nil {
		return nil
	}

	orderId := utils.PathSearch("order_id", curJson, nil)
	productId := utils.PathSearch("product_id", curJson, nil)

	// If both fields are nil, return nil to avoid empty object
	if orderId == nil && productId == nil {
		return nil
	}

	return []interface{}{
		map[string]interface{}{
			"order_id":   orderId,
			"product_id": productId,
		},
	}
}
