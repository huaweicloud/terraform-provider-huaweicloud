// ---------------------------------------------------------------
// *** AUTO GENERATED CODE ***
// @Product CC
// ---------------------------------------------------------------

package cc

import (
	"context"
	"regexp"
	"strings"

	"github.com/chnsz/golangsdk"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/jmespath/go-jmespath"
)

func ResourceCloudConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCloudConnectionCreate,
		UpdateContext: resourceCloudConnectionUpdate,
		ReadContext:   resourceCloudConnectionRead,
		DeleteContext: resourceCloudConnectionDelete,
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
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The cloud connection name.`,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[\x{4E00}-\x{9FFC}A-Za-z-_0-9.]*$`),
						"the input is invalid"),
					validation.StringLenBetween(1, 64),
				),
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The description about the cloud connection.`,
				ValidateFunc: validation.All(
					validation.StringMatch(regexp.MustCompile(`^[^<>]+$`),
						"the input is invalid"),
					validation.StringLenBetween(0, 255),
				),
			},
			"enterprise_project_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The enterprise project id of the cloud connection.`,
			},
			"domain_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Domain ID.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the cloud connection.`,
			},
			"used_scene": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The Scenario.`,
			},
			"network_instance_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of network instances associated with the cloud connection instance.`,
			},
			"bandwidth_package_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of bandwidth packages associated with the cloud connection instance.`,
			},
			"inter_region_bandwidth_number": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: `The number of inter-domain bandwidths associated with the cloud connection instance.`,
			},
		},
	}
}

func resourceCloudConnectionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// createCloudConnection: create a Cloud Connect.
	var (
		createCloudConnectionHttpUrl = "v3/{domain_id}/ccaas/cloud-connections"
		createCloudConnectionProduct = "cc"
	)
	createCloudConnectionClient, err := config.NewServiceClient(createCloudConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CloudConnection Client: %s", err)
	}

	createCloudConnectionPath := createCloudConnectionClient.Endpoint + createCloudConnectionHttpUrl
	createCloudConnectionPath = strings.Replace(createCloudConnectionPath, "{domain_id}", config.DomainID, -1)

	createCloudConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			201,
		},
	}
	createCloudConnectionOpt.JSONBody = utils.RemoveNil(buildCreateCloudConnectionBodyParams(d, config))
	createCloudConnectionResp, err := createCloudConnectionClient.Request("POST", createCloudConnectionPath, &createCloudConnectionOpt)
	if err != nil {
		return diag.Errorf("error creating CloudConnection: %s", err)
	}

	createCloudConnectionRespBody, err := utils.FlattenResponse(createCloudConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	id, err := jmespath.Search("cloud_connection.id", createCloudConnectionRespBody)
	if err != nil {
		return diag.Errorf("error creating CloudConnection: ID is not found in API response")
	}
	d.SetId(id.(string))

	return resourceCloudConnectionRead(ctx, d, meta)
}

func buildCreateCloudConnectionBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_connection": buildCreateCloudConnectionCloudConnectionChildBody(d, config),
	}
	return bodyParams
}

func buildCreateCloudConnectionCloudConnectionChildBody(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	params := map[string]interface{}{
		"name":                  utils.ValueIngoreEmpty(d.Get("name")),
		"description":           utils.ValueIngoreEmpty(d.Get("description")),
		"enterprise_project_id": utils.ValueIngoreEmpty(common.GetEnterpriseProjectID(d, config)),
	}
	return params
}

func resourceCloudConnectionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	var mErr *multierror.Error

	// getCloudConnection: Query the Cloud Connection
	var (
		getCloudConnectionHttpUrl = "v3/{domain_id}/ccaas/cloud-connections/{id}"
		getCloudConnectionProduct = "cc"
	)
	getCloudConnectionClient, err := config.NewServiceClient(getCloudConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CloudConnection Client: %s", err)
	}

	getCloudConnectionPath := getCloudConnectionClient.Endpoint + getCloudConnectionHttpUrl
	getCloudConnectionPath = strings.Replace(getCloudConnectionPath, "{domain_id}", config.DomainID, -1)
	getCloudConnectionPath = strings.Replace(getCloudConnectionPath, "{id}", d.Id(), -1)

	getCloudConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			200,
		},
	}
	getCloudConnectionResp, err := getCloudConnectionClient.Request("GET", getCloudConnectionPath, &getCloudConnectionOpt)

	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving CloudConnection")
	}

	getCloudConnectionRespBody, err := utils.FlattenResponse(getCloudConnectionResp)
	if err != nil {
		return diag.FromErr(err)
	}

	mErr = multierror.Append(
		mErr,
		d.Set("region", region),
		d.Set("name", utils.PathSearch("cloud_connection.name", getCloudConnectionRespBody, nil)),
		d.Set("description", utils.PathSearch("cloud_connection.description", getCloudConnectionRespBody, nil)),
		d.Set("domain_id", utils.PathSearch("cloud_connection.domain_id", getCloudConnectionRespBody, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("cloud_connection.enterprise_project_id", getCloudConnectionRespBody, nil)),
		d.Set("status", utils.PathSearch("cloud_connection.status", getCloudConnectionRespBody, nil)),
		d.Set("used_scene", utils.PathSearch("cloud_connection.used_scene", getCloudConnectionRespBody, nil)),
		d.Set("network_instance_number", utils.PathSearch("cloud_connection.network_instance_number", getCloudConnectionRespBody, nil)),
		d.Set("bandwidth_package_number", utils.PathSearch("cloud_connection.bandwidth_package_number", getCloudConnectionRespBody, nil)),
		d.Set("inter_region_bandwidth_number", utils.PathSearch("cloud_connection.inter_region_bandwidth_number", getCloudConnectionRespBody, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceCloudConnectionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	updateCloudConnectionhasChanges := []string{
		"name",
		"description",
	}

	if d.HasChanges(updateCloudConnectionhasChanges...) {
		// updateCloudConnection: update the Cloud Connection
		var (
			updateCloudConnectionHttpUrl = "v3/{domain_id}/ccaas/cloud-connections/{id}"
			updateCloudConnectionProduct = "cc"
		)
		updateCloudConnectionClient, err := config.NewServiceClient(updateCloudConnectionProduct, region)
		if err != nil {
			return diag.Errorf("error creating CloudConnection Client: %s", err)
		}

		updateCloudConnectionPath := updateCloudConnectionClient.Endpoint + updateCloudConnectionHttpUrl
		updateCloudConnectionPath = strings.Replace(updateCloudConnectionPath, "{domain_id}", config.DomainID, -1)
		updateCloudConnectionPath = strings.Replace(updateCloudConnectionPath, "{id}", d.Id(), -1)

		updateCloudConnectionOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			OkCodes: []int{
				200,
			},
		}
		updateCloudConnectionOpt.JSONBody = utils.RemoveNil(buildUpdateCloudConnectionBodyParams(d, config))
		_, err = updateCloudConnectionClient.Request("PUT", updateCloudConnectionPath, &updateCloudConnectionOpt)
		if err != nil {
			return diag.Errorf("error updating CloudConnection: %s", err)
		}
	}
	return resourceCloudConnectionRead(ctx, d, meta)
}

func buildUpdateCloudConnectionBodyParams(d *schema.ResourceData, config *config.Config) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"cloud_connection": buildUpdateCloudConnectionCloudConnectionChildBody(d),
	}
	return bodyParams
}

func buildUpdateCloudConnectionCloudConnectionChildBody(d *schema.ResourceData) map[string]interface{} {
	params := map[string]interface{}{
		"name":        utils.ValueIngoreEmpty(d.Get("name")),
		"description": utils.ValueIngoreEmpty(d.Get("description")),
	}
	return params
}

func resourceCloudConnectionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)

	// deleteCloudConnection: missing operation notes
	var (
		deleteCloudConnectionHttpUrl = "v3/{domain_id}/ccaas/cloud-connections/{id}"
		deleteCloudConnectionProduct = "cc"
	)
	deleteCloudConnectionClient, err := config.NewServiceClient(deleteCloudConnectionProduct, region)
	if err != nil {
		return diag.Errorf("error creating CloudConnection Client: %s", err)
	}

	deleteCloudConnectionPath := deleteCloudConnectionClient.Endpoint + deleteCloudConnectionHttpUrl
	deleteCloudConnectionPath = strings.Replace(deleteCloudConnectionPath, "{domain_id}", config.DomainID, -1)
	deleteCloudConnectionPath = strings.Replace(deleteCloudConnectionPath, "{id}", d.Id(), -1)

	deleteCloudConnectionOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes: []int{
			204,
		},
	}
	_, err = deleteCloudConnectionClient.Request("DELETE", deleteCloudConnectionPath, &deleteCloudConnectionOpt)
	if err != nil {
		return diag.Errorf("error deleting CloudConnection: %s", err)
	}

	return nil
}
