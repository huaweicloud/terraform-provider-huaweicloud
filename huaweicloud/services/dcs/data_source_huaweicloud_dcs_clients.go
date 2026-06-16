package dcs

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk/pagination"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API DCS GET /v2/{project_id}/instances/{instance_id}/clients
func DataSourceDcsClients() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDcsClientsRead,

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
			"node_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"addr": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sort": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"order": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"clients": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     dcsClientSchema(),
			},
		},
	}
}

func dcsClientSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"addr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fd": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cmd": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"age": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"idle": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"db": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"flags": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sub": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"psub": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"multi": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"qbuf": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"qbuf_free": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"obl": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"oll": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"omem": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"events": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"peer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceDcsClientsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("dcs", region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	httpUrl := "v2/{project_id}/instances/{instance_id}/clients"
	listPath := client.Endpoint + httpUrl
	listPath = strings.ReplaceAll(listPath, "{project_id}", client.ProjectID)
	listPath = strings.ReplaceAll(listPath, "{instance_id}", d.Get("instance_id").(string))
	listPath += buildGetClientsQueryParams(d)

	listResp, err := pagination.ListAllItems(
		client,
		"offset",
		listPath,
		&pagination.QueryOpts{MarkerField: ""})
	if err != nil {
		return diag.Errorf("error retrieving DCS clients: %s", err)
	}

	listRespJson, err := json.Marshal(listResp)
	if err != nil {
		return diag.FromErr(err)
	}

	var listRespBody interface{}
	err = json.Unmarshal(listRespJson, &listRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(id.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("clients", flattenGetClientsBody(listRespBody)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetClientsQueryParams(d *schema.ResourceData) string {
	res := fmt.Sprintf("?node_id=%v", d.Get("node_id"))

	if v, ok := d.GetOk("addr"); ok {
		res = fmt.Sprintf("%s&addr=%v", res, v)
	}
	if v, ok := d.GetOk("sort"); ok {
		res = fmt.Sprintf("%s&sort=%v", res, v)
	}
	if v, ok := d.GetOk("order"); ok {
		res = fmt.Sprintf("%s&order=%v", res, v)
	}

	return res
}

func flattenGetClientsBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("clients", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))

	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"id":        utils.PathSearch("id", v, nil),
			"addr":      utils.PathSearch("addr", v, nil),
			"fd":        utils.PathSearch("fd", v, nil),
			"name":      utils.PathSearch("name", v, nil),
			"cmd":       utils.PathSearch("cmd", v, nil),
			"age":       utils.PathSearch("age", v, nil),
			"idle":      utils.PathSearch("idle", v, nil),
			"db":        utils.PathSearch("db", v, nil),
			"flags":     utils.PathSearch("flags", v, nil),
			"sub":       utils.PathSearch("sub", v, nil),
			"psub":      utils.PathSearch("psub", v, nil),
			"multi":     utils.PathSearch("multi", v, nil),
			"qbuf":      utils.PathSearch("qbuf", v, nil),
			"qbuf_free": utils.PathSearch("qbuf_free", v, nil),
			"obl":       utils.PathSearch("obl", v, nil),
			"oll":       utils.PathSearch("oll", v, nil),
			"omem":      utils.PathSearch("omem", v, nil),
			"events":    utils.PathSearch("events", v, nil),
			"network":   utils.PathSearch("network", v, nil),
			"peer":      utils.PathSearch("peer", v, nil),
			"user":      utils.PathSearch("user", v, nil),
		})
	}
	return res
}
