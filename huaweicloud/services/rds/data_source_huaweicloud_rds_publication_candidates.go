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

// @API RDS POST /v3/{project_id}/instances/{instance_id}/replication/publication-candidates
func DataSourceRdsPublicationCandidates() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceRdsPublicationCandidatesRead,
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
			"publication_instance_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publication_instance_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"instance_publications": {
				Type:     schema.TypeList,
				Elem:     instancePublicationsSchema(),
				Computed: true,
			},
		},
	}
}

func instancePublicationsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"instance_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"publication_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRdsPublicationCandidatesRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var (
		httpUrl = "v3/{project_id}/instances/{instance_id}/replication/publication-candidates"
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

	var publicationCandidates []interface{}
	offset := 0
	for {
		getOpt.JSONBody = buildPublicationCandidatesQueryBody(d, offset)
		getResp, err := client.Request("POST", getPath, &getOpt)
		if err != nil {
			return diag.Errorf("error retrieving RDS publication candidates: %s", err)
		}

		getRespBody, err := utils.FlattenResponse(getResp)
		if err != nil {
			return diag.FromErr(err)
		}
		res := flattenGetPublicationCandidatesResponseBody(getRespBody)

		if len(res) == 0 {
			break
		}
		publicationCandidates = append(publicationCandidates, res...)
		offset += len(publicationCandidates)
	}

	dataSourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(dataSourceId)

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("instance_publications", publicationCandidates),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func buildPublicationCandidatesQueryBody(d *schema.ResourceData, offset int) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"limit":                     1000,
		"offset":                    offset,
		"publication_instance_id":   utils.ValueIgnoreEmpty(d.Get("publication_instance_id").(string)),
		"publication_instance_name": utils.ValueIgnoreEmpty(d.Get("publication_instance_name").(string)),
		"publication_name":          utils.ValueIgnoreEmpty(d.Get("publication_name").(string)),
	}

	return bodyParams
}

func flattenGetPublicationCandidatesResponseBody(resp interface{}) []interface{} {
	curJson := utils.PathSearch("instance_publications", resp, make([]interface{}, 0))
	curArray := curJson.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"instance_id":      utils.PathSearch("instance_id", v, nil),
			"instance_name":    utils.PathSearch("instance_name", v, nil),
			"publication_id":   utils.PathSearch("publication_id", v, nil),
			"publication_name": utils.PathSearch("publication_name", v, nil),
		})
	}
	return res
}
