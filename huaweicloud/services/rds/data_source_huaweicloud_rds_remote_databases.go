package rds

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

// @API RDS POST /v3/{project_id}/instances/{instance_id}/replication/remote-databases
func DataSourceRemoteDatabases() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRemoteDatabasesRead,

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
			"server_ip": {
				Type:     schema.TypeString,
				Required: true,
			},
			"server_port": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_user_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"login_user_password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"databases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     remoteDatabasesDatabasesSchema(),
			},
		},
	}
}

func remoteDatabasesDatabasesSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"character_set": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRemoteDatabasesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/remote-databases"
		product = "rds"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating RDS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	var databases []interface{}
	offset := 0
	for {
		getOpt.JSONBody = buildRemoteDatabasesQueryBody(d, offset)
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS remote databases: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res := flattenGetRemoteDatabasesResponseBody(getRespBody)
		if len(res) == 0 {
			break
		}
		databases = append(databases, res...)
		offset += 100
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("databases", databases),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildRemoteDatabasesQueryBody(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":               100,
		"offset":              offset,
		"server_ip":           d.Get("server_ip").(string),
		"server_port":         d.Get("server_port").(string),
		"login_user_name":     d.Get("login_user_name").(string),
		"login_user_password": d.Get("login_user_password").(string),
	}

	return bodyParams
}

func flattenGetRemoteDatabasesResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("databases", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"name":          utils.PathSearch("name", v, nil),
			"character_set": utils.PathSearch("character_set", v, nil),
			"state":         utils.PathSearch("state", v, nil),
		})
	}
	return res
}
