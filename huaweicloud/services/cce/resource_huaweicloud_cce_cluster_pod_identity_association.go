package cce

import (
	"context"
	"errors"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var podIdentityAssociationNoneUpdatableParams = []string{
	"cluster_id", "namespace", "service_account", "tags",
}

// @API CCE POST /api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations
// @API CCE GET /api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}
// @API CCE PUT /api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}
// @API CCE DELETE /api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}
func ResourceClusterPodIdentityAssociation() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterPodIdentityAssociationCreate,
		ReadContext:   resourceClusterPodIdentityAssociationRead,
		UpdateContext: resourceClusterPodIdentityAssociationUpdate,
		DeleteContext: resourceClusterPodIdentityAssociationDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceClusterPodIdentityAssociationImport,
		},

		CustomizeDiff: config.FlexibleForceNew(podIdentityAssociationNoneUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"cluster_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"namespace": {
				Type:     schema.TypeString,
				Required: true,
			},
			"service_account": {
				Type:     schema.TypeString,
				Required: true,
			},
			"agency_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"tags": common.TagsSchema(),
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"created_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"updated_at": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceClusterPodIdentityAssociationCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterID = d.Get("cluster_id").(string)
		createURL = "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations"
	)

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	createPath := client.Endpoint + createURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{cluster_id}", clusterID)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildPodIdentityAssociationCreateParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating CCE cluster pod identity association: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error retrieving CCE cluster pod identity association creation response: %s", err)
	}

	uid := utils.PathSearch("uid", createRespBody, "").(string)
	if uid == "" {
		return diag.Errorf("error getting uid from CCE cluster pod identity association creation response")
	}

	d.SetId(uid)

	return resourceClusterPodIdentityAssociationRead(ctx, d, meta)
}

func buildPodIdentityAssociationCreateParams(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"namespace":      d.Get("namespace").(string),
		"serviceAccount": d.Get("service_account").(string),
		"agencyName":     d.Get("agency_name").(string),
	}
	if tags := d.Get("tags").(map[string]interface{}); len(tags) > 0 {
		params["tags"] = utils.ExpandResourceTags(tags)
	}
	return params
}

func resourceClusterPodIdentityAssociationRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterID = d.Get("cluster_id").(string)
		getURL    = "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}"
	)

	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	getPath := client.Endpoint + getURL
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{cluster_id}", clusterID)
	getPath = strings.ReplaceAll(getPath, "{association_id}", d.Id())

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	resp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error getting CCE cluster pod identity association")
	}

	respBody, err := utils.FlattenResponse(resp)
	if err != nil {
		return diag.Errorf("error retrieving CCE cluster pod identity association: %s", err)
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("cluster_id", clusterID),
		d.Set("namespace", utils.PathSearch("namespace", respBody, nil)),
		d.Set("service_account", utils.PathSearch("serviceAccount", respBody, nil)),
		d.Set("agency_name", utils.PathSearch("agencyName", respBody, nil)),
		d.Set("tags", utils.FlattenTagsToMap(utils.PathSearch("tags", respBody, nil))),
		d.Set("created_at", utils.PathSearch("createdAt", respBody, nil)),
		d.Set("updated_at", utils.PathSearch("updatedAt", respBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceClusterPodIdentityAssociationUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChanges("agency_name") {
		cfg := meta.(*config.Config)
		region := cfg.GetRegion(d)
		clusterID := d.Get("cluster_id").(string)
		updateURL := "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}"

		client, err := cfg.NewServiceClient("cce", region)
		if err != nil {
			return diag.Errorf("error creating CCE client: %s", err)
		}
		updatePath := client.Endpoint + updateURL
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{cluster_id}", clusterID)
		updatePath = strings.ReplaceAll(updatePath, "{association_id}", d.Id())

		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: map[string]interface{}{
				"agencyName": d.Get("agency_name").(string),
			},
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating CCE cluster pod identity association: %s", err)
		}
	}

	return resourceClusterPodIdentityAssociationRead(ctx, d, meta)
}

func resourceClusterPodIdentityAssociationDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		clusterID = d.Get("cluster_id").(string)
		deleteURL = "api/v3/projects/{project_id}/clusters/{cluster_id}/pod-identity-associations/{association_id}"
	)
	client, err := cfg.NewServiceClient("cce", region)
	if err != nil {
		return diag.Errorf("error creating CCE client: %s", err)
	}

	deletePath := client.Endpoint + deleteURL
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{cluster_id}", clusterID)
	deletePath = strings.ReplaceAll(deletePath, "{association_id}", d.Id())

	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		OkCodes: []int{204},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return diag.Errorf("error deleting CCE cluster pod identity association: %s", err)
	}

	return nil
}

func resourceClusterPodIdentityAssociationImport(_ context.Context, d *schema.ResourceData, _ interface{}) ([]*schema.ResourceData, error) {
	importId := d.Id()
	parts := strings.Split(importId, "/")
	if len(parts) != 2 {
		return nil, errors.New("invalid format specified for import ID, must be <cluster_id>/<id>")
	}

	d.SetId(parts[1])
	mErr := multierror.Append(nil,
		d.Set("cluster_id", parts[0]),
	)
	return []*schema.ResourceData{d}, mErr.ErrorOrNil()
}
