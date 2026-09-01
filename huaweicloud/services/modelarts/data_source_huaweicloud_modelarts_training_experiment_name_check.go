package modelarts

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts GET /v2/{project_id}/training-experiment-names
func DataSourceTrainingExperimentNameCheck() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceTrainingExperimentNameCheckRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The region where the ModelArts service is located.`,
			},

			// Required parameters.
			"experiment_name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the training experiment to be checked.`,
			},

			// Optional parameters.
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The ID of the workspace to which the training experiment belongs.`,
			},

			// Attributes.
			"is_duplicate": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: `Whether the training experiment name is duplicate.`,
			},
		},
	}
}

func buildTrainingExperimentNameCheckQueryParams(d *schema.ResourceData) string {
	queryParams := fmt.Sprintf("?experiment_name=%s", d.Get("experiment_name").(string))

	if v, ok := d.GetOk("workspace_id"); ok {
		queryParams = fmt.Sprintf("%s&workspace_id=%s", queryParams, v.(string))
	}

	return queryParams
}

func dataSourceTrainingExperimentNameCheckRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		httpUrl = "v2/{project_id}/training-experiment-names"
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	requestPath := client.Endpoint + httpUrl
	requestPath = strings.ReplaceAll(requestPath, "{project_id}", client.ProjectID)
	requestPath += buildTrainingExperimentNameCheckQueryParams(d)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", requestPath, &opt)
	if err != nil {
		return diag.Errorf("error checking training experiment name: %s", err)
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error flattening response: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("is_duplicate", utils.PathSearch("is_duplicate", respBody, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}
