package evs

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

// @API EVS GET /v3/{project_id}/recycle-bin-volumes/policy
func DataSourceRecycleBinPolicy() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRecycleBinPolicyRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"switch": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"threshold_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"keep_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceRecycleBinPolicyRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		product = "evs"
		httpUrl = "v3/{project_id}/recycle-bin-volumes/policy"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating EVS client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestOpts := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	resp, err := client.Request("GET", requestPath, &requestOpts)
	if err != nil {
		return diag.Errorf("error retrieving EVS recycle bin policy: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.FromErr(err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}

	d.SetId(generateUUID.String())

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("switch", utils.PathSearch("switch", respBody, nil)),
		d.Set("threshold_time", utils.PathSearch("threshold_time", respBody, nil)),
		d.Set("keep_time", utils.PathSearch("keep_time", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
