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
			StateContext: schema.ImportStatePassthroughContext,
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
			"tags": {
				Type:     schema.TypeSet,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"sync": {
							Type:     schema.TypeBool,
							Required: true,
						},
						"key": {
							Type:     schema.TypeString,
							Required: true,
						},
						"values": {
							Type:     schema.TypeSet,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCloudServiceAccessBodyParams(d)),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating AOM cloud service access: %s", err)
	}

	d.SetId(instanceID + "/" + service)

	return resourceCloudServiceAccessRead(ctx, d, meta)
}

func buildCloudServiceAccessBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"provider": d.Get("service"),
		"tag_sync": d.Get("tag_sync"),
		"tags":     buildCloudServiceAccessTags(d),
	}

	return bodyParams
}

func buildCloudServiceAccessTags(d *schema.ResourceData) []map[string]interface{} {
	rawParams := d.Get("tags").(*schema.Set).List()
	rst := make([]map[string]interface{}, 0, len(rawParams))
	for _, val := range rawParams {
		raw := val.(map[string]interface{})
		params := map[string]interface{}{
			"sync":   raw["sync"],
			"key":    raw["key"],
			"values": raw["values"].(*schema.Set).List(),
		}
		rst = append(rst, params)
	}

	return rst
}

func resourceCloudServiceAccessRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.NewServiceClient("aom", region)
	if err != nil {
		return diag.Errorf("error creating AOM client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, must be <instance_id>/<service>")
	}
	instanceID := parts[0]
	service := parts[1]

	access, err := getCloudServiceAccesss(client, instanceID, service)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving cloud service access")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("instance_id", instanceID),
		d.Set("service", service),
		d.Set("tag_sync", utils.PathSearch("prometheus.tag_sync", access, nil)),
		d.Set("tags", utils.PathSearch("tags", access, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func getCloudServiceAccesss(client *golangsdk.ServiceClient, instanceID, service string) (interface{}, error) {
	getHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{prom_instance_id}", instanceID)
	getPath = strings.ReplaceAll(getPath, "{provider}", service)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
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

	updateChanges := []string{
		"tag_sync",
		"tags",
	}

	if d.HasChanges(updateChanges...) {
		updateHttpUrl := "v1/{project_id}/prometheus/{prom_instance_id}/cloud-service/{provider}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{project_id}", client.ProjectID)
		updatePath = strings.ReplaceAll(updatePath, "{prom_instance_id}", d.Get("instance_id").(string))
		updatePath = strings.ReplaceAll(updatePath, "{provider}", d.Get("service").(string))
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(buildCloudServiceAccessBodyParams(d)),
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
	_, err = getCloudServiceAccesss(client, instanceID, service)
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
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting cloud service access")
	}

	return nil
}
