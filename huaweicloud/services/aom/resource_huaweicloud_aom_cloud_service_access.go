package aom

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API AOM POST /v1/{project_id}/prometheus/{prom_instance_id}/cloud-service
// @API AOM PUT /v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}
// @API AOM DELETE /v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}
// @API AOM GET /v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}
func ResourceCloudServiceAccess() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudServiceAccessCreate,
		ReadContext:   resourceCloudServiceAccessRead,
		UpdateContext: resourceCloudServiceAccessUpdate,
		DeleteContext: resourceCloudServiceAccessDelete,

		Importer: &schema.ResourceImporter{
			StateContext: resourceCloudServiceAccessImportState,
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"tag_sync": {
				Type:     schema.TypeString,
				Required: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceCloudServiceAccessCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	service := d.Get("service").(string)

	createHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{prom_instance_id}", instanceID)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
		JSONBody:         utils.RemoveNil(buildCloudServiceAccessBodyParams(cfg, client, d, "")),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM cloud service access: %s", err)
	}

	d.SetId(instanceID + "/" + service)

	return resourceCloudServiceAccessRead(ctx, d, meta)
}

func buildCloudServiceAccessBodyParams(cfg *config.Config, client *golangsdk.ServiceClient, d *schema.ResourceData,
	id string) map[string]interface{} {
	bodyParams := map[string]interface{}{
		// id is only used for update
		"id":            utils.ValueIgnoreEmpty(id),
		"provider":      d.Get("service"),
		"tag_sync":      d.Get("tag_sync"),
		"tags":          []interface{}{},
		"prometheus_id": d.Get("instance_id"),
		"ep_id":         utils.ValueIgnoreEmpty(cfg.GetEnterpriseProjectID(d)),
		"project_id":    client.ProjectID,
	}

	return bodyParams
}

func resourceCloudServiceAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	// get eps_id for import
	epsID := cfg.GetEnterpriseProjectID(d)
	if epsID == "" {
		instance, err := GetPrometheusInstanceById(client, d.Get("instance_id").(string))
		if err != nil {
			return diag.FromErr(err)
		}
		epsID = utils.PathSearch("enterprise_project_id", instance, "").(string)
		if epsID == "" {
			return diag.Errorf("error getting enterprise project ID from instance")
		}
	}

	access, err := getCloudServiceAccesss(client, d, epsID)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving cloud service access")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", utils.PathSearch("resource_id", access, nil)),
		d.Set("service", utils.PathSearch("prometheus.provider", access, nil)),
		d.Set("tag_sync", utils.PathSearch("prometheus.tag_sync", access, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("ep_id", access, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getCloudServiceAccesss(client *golangsdk.ServiceClient, d *schema.ResourceData, epsID string) (interface{}, error) {
	getHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{prom_instance_id}", d.Get("instance_id").(string))
	getPath = strings.ReplaceAll(getPath, "{provider}", d.Get("service").(string))
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type":          "application/json",
			"Enterprise-Project-Id": epsID,
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, fmt.Errorf("error retrieving cloud service access: %s", err)
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, fmt.Errorf("error flattening cloud service access: %s", err)
	}

	// API will return 200 and nil if `instance_id` and `servcie` is invalid
	if getRespBody == nil {
		return nil, golangsdk.ErrDefault404{}
	}

	return getRespBody, nil
}

func resourceCloudServiceAccessUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	access, err := getCloudServiceAccesss(client, d, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return diag.FromErr(err)
	}

	id := utils.PathSearch("id", access, "").(string)
	if id == "" {
		return diag.Errorf("unable to find ID from API response")
	}

	updateChanges := []string{
		"tag_sync",
	}

	if d.HasChanges(updateChanges...) {
		updateHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{prom_instance_id}", d.Get("instance_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{provider}", d.Get("service").(string))
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
			JSONBody:         utils.RemoveNil(buildCloudServiceAccessBodyParams(cfg, client, d, id)),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating cloud service access: %s", err)
		}
	}

	return resourceCloudServiceAccessRead(ctx, d, meta)
}

func resourceCloudServiceAccessDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	instanceID := d.Get("instance_id").(string)
	service := d.Get("service").(string)

	// precheck
	// DELETE will return 200 even deleting a non exist service
	_, err = getCloudServiceAccesss(client, d, cfg.GetEnterpriseProjectID(d))
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting cloud service access")
	}

	deleteHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{project_id}", client.ProjectID)
	deletePath = strings.ReplaceAll(deletePath, "{prom_instance_id}", instanceID)
	deletePath = strings.ReplaceAll(deletePath, "{provider}", service)
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      buildRequestMoreHeaders(cfg.GetEnterpriseProjectID(d)),
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting cloud service access")
	}

	return nil
}

func resourceCloudServiceAccessImportState(_ context.Context, d *schema.ResourceData, _ interface{}) (
	[]*schema.ResourceData, error) {
	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format of import ID, must be <instance_id>/<service>")
	}

	d.Set("instance_id", parts[0])
	d.Set("service", parts[1])

	return []*schema.ResourceData{d}, nil
}
