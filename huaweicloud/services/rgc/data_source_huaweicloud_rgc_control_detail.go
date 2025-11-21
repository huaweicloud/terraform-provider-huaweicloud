package rgc

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API RGC GET /v1/governance/controls/{control_id}
func DataSourceControlDetail() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceControlRead,
		Schema: map[string]*schema.Schema{
			"control_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"identifier": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"implementation": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"guidance": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"resource": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"service": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"behavior": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"control_objective": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"framework": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"artifacts": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"content": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"en": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"ch": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"aliases": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"owner": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"severity": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"release_date": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceControlRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	var getControlProduct = "rgc"
	getControlClient, err := cfg.NewServiceClient(getControlProduct, region)
	if err != nil {
		return diag.Errorf("error creating RGC client: %s", err)
	}

	getControlRespBody, err := getControl(getControlClient, d)

	if err != nil {
		return diag.Errorf("error retrieving RGC control detail: %s", err)
	}

	artifacts, err := parseArtifacts(getControlRespBody)
	if err != nil {
		return diag.FromErr(err)
	}

	randUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randUUID)

	mErr = multierror.Append(mErr,
		d.Set("identifier", utils.PathSearch("identifier", getControlRespBody, nil)),
		d.Set("implementation", utils.PathSearch("implementation", getControlRespBody, nil)),
		d.Set("guidance", utils.PathSearch("guidance", getControlRespBody, nil)),
		d.Set("resource", utils.PathSearch("resource", getControlRespBody, nil)),
		d.Set("service", utils.PathSearch("service", getControlRespBody, nil)),
		d.Set("behavior", utils.PathSearch("behavior", getControlRespBody, nil)),
		d.Set("control_objective", utils.PathSearch("control_objective", getControlRespBody, nil)),
		d.Set("framework", utils.PathSearch("framework", getControlRespBody, nil)),
		d.Set("artifacts", artifacts),
		d.Set("aliases", utils.PathSearch("aliases", getControlRespBody, nil)),
		d.Set("owner", utils.PathSearch("owner", getControlRespBody, nil)),
		d.Set("severity", utils.PathSearch("severity", getControlRespBody, nil)),
		d.Set("version", utils.PathSearch("version", getControlRespBody, nil)),
		d.Set("release_date", utils.PathSearch("release_date", getControlRespBody, nil)),
		d.Set("region", region),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func parseArtifacts(body interface{}) ([]interface{}, error) {
	artifacts := make([]interface{}, 0)
	artifactSlice := utils.PathSearch("artifacts", body, make([]interface{}, 0)).([]interface{})
	if len(artifactSlice) == 0 {
		return nil, errors.New("can not get artifacts from RGC")
	}

	for _, artifact := range artifactSlice {
		artifactMap := artifact.(map[string]interface{})
		if v, ok := artifactMap["content"]; ok {
			artifactMap["content"] = []interface{}{v}
		}

		if v, ok := artifactMap["type"]; ok {
			artifactMap["type"] = v.(string)
		}

		artifacts = append(artifacts, artifactMap)
	}

	return artifacts, nil
}

func getControl(client *golangsdk.ServiceClient, d *schema.ResourceData) (interface{}, error) {
	controlId := d.Get("control_id").(string)
	var (
		getControlHttpUrl = "v1/governance/controls/{control_id}"
	)
	getControlHttpPath := client.Endpoint + getControlHttpUrl
	getControlHttpPath = strings.ReplaceAll(getControlHttpPath, "{control_id}", controlId)

	getControlHttpOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}
	getControlHttpResp, err := client.Request("GET", getControlHttpPath, &getControlHttpOpt)
	if err != nil {
		return nil, err
	}
	getControlRespBody, err := utils.FlattenResponse(getControlHttpResp)
	if err != nil {
		return nil, err
	}
	return getControlRespBody, nil
}
