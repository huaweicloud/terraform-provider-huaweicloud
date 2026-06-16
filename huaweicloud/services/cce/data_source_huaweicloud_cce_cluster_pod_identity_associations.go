package cce

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations
func DataSourceCCEClusterPodIdentityAssociations() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceCCEClusterPodIdentityAssociationsRead,

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cluster_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Specifies the cluster ID.",
			},
			"associations": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     podIdentityAssociationSchema(),
			},
		},
	}
}

func podIdentityAssociationSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"uid": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"service_account": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": common.TagsComputedSchema(),
		},
	}
}

func dataSourceCCEClusterPodIdentityAssociationsRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var mErr *multierror.Error

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	httpUrl := "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", d.Get("cluster_id").(string))

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return diag.Errorf("error getting CCE cluster pod identity associations: %s", err)
	}

	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error retrieving CCE cluster pod identity associations: %s", err)
	}

	generateUUID, err := uuid.NewRandom()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(generateUUID.String())

	mErr = multierror.Append(
		d.Set("region", region),
		d.Set("associations", flattenPodIdentityAssociationsBody(getRespBody)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenPodIdentityAssociationsBody(resp interface{}) []interface{} {
	curArray := resp.([]interface{})
	res := make([]interface{}, 0, len(curArray))
	for _, v := range curArray {
		res = append(res, map[string]interface{}{
			"uid":             utils.PathSearch("uid", v, nil),
			"cluster_id":      utils.PathSearch("clusterId", v, nil),
			"namespace":       utils.PathSearch("namespace", v, nil),
			"service_account": utils.PathSearch("serviceAccount", v, nil),
			"agency_name":     utils.PathSearch("agencyName", v, nil),
			"created_at":      utils.PathSearch("createdAt", v, nil),
			"updated_at":      utils.PathSearch("updatedAt", v, nil),
			"tags":            utils.FlattenTagsToMap(utils.PathSearch("tags", v, nil)),
		})
	}
	return res
}
