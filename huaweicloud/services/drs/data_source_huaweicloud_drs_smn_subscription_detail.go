package drs

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

// @API DRS GET /v5/{project_id}/jobs/{job_id}/smn-subscription-detail
func DataSourceDrsSmnSubscriptionDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceDrsSmnSubscriptionDetailRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "Specifies the region in which to query the data source. If omitted, the provider-level region will be used.",
			},
			"job_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the ID of the DRS job to query SMN subscription detail.",
			},
			"topic_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The name of the SMN topic.",
			},
			"topic_urn": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique resource identifier of the SMN topic.",
			},
			"delay_time": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The delay time in seconds.",
			},
			"rpo_delay": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The RPO delay time in seconds.",
			},
			"rto_delay": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The RTO delay time in seconds.",
			},
			"is_alarm_user": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "Whether to send alarm to the subscribed user.",
			},
		},
	}
}

func dataSourceDrsSmnSubscriptionDetailRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	var (
		httpUrl = "v5/{project_id}/jobs/{job_id}/smn-subscription-detail"
		product = "drs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DRS client: %s", err)
	}

	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{job_id}", d.Get("job_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error retrieving DRS SMN subscription detail: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.FromErr(err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("topic_name", utils.PathSearch("topic_name", getRespBody, nil)),
		d.Set("topic_urn", utils.PathSearch("topic_urn", getRespBody, nil)),
		d.Set("delay_time", utils.PathSearch("delay_time", getRespBody, nil)),
		d.Set("rpo_delay", utils.PathSearch("rpo_delay", getRespBody, nil)),
		d.Set("rto_delay", utils.PathSearch("rto_delay", getRespBody, nil)),
		d.Set("is_alarm_user", utils.PathSearch("is_alarm_user", getRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
