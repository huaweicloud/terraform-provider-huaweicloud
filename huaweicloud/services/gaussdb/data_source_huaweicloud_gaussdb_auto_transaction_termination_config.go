package gaussdb

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API GaussDB GET /v3/{project_id}/instances/{instance_id}/transactions/get-auto-kill-config
func DataSourceAutoTransactionTerminationConfig() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAutoTransactionTerminationConfigRead,

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
			"type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"usernames": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"threshold": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"auto_stop": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"database_names": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"database_names_select": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceAutoTransactionTerminationConfigRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("opengauss", region)
	if err != nil {
		return diag.Errorf("error creating GaussDB client: %s", err)
	}

	httpUrl := "v3/{project_id}/instances/{instance_id}/transactions/get-auto-kill-config"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{instance_id}", d.Get("instance_id").(string))
	getPath += buildGetAutoTransactionTerminationConfigQueryParams(d)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving GaussDB auto transaction termination config: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	dataSourceId, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("type", utils.PathSearch("type", getRespBody, nil)),
		d.Set("usernames", utils.PathSearch("usernames", getRespBody, nil)),
		d.Set("threshold", utils.PathSearch("threshold", getRespBody, nil)),
		d.Set("auto_stop", utils.PathSearch("auto_stop", getRespBody, nil)),
		d.Set("database_names", utils.PathSearch("database_names", getRespBody, nil)),
		d.Set("database_names_select", utils.PathSearch("database_names_select", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildGetAutoTransactionTerminationConfigQueryParams(d *schema.ResourceData) string {
	res := ""
	res = fmt.Sprintf("%s&type=%v", res, d.Get("type"))

	if res != "" {
		res = "?" + res[1:]
	}
	return res
}
