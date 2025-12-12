package vpcep

import (
	"context"
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

var vpcepServiceConnectionUpdateNonUpdatableParams = []string{
	"service_id", "endpoint_id",
}

// @API VPCEP PUT /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections/description
// @API VPCEP GET /v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections
func ResourceVPCEndpointServiceConnectionUpdate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCEPServiceConnectionUpdateCreateOrUpdate,
		ReadContext:   resourceVPCEPServiceConnectionUpdateRead,
		UpdateContext: resourceVPCEPServiceConnectionUpdateCreateOrUpdate,
		DeleteContext: resourceVPCEPServiceConnectionUpdateDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		CustomizeDiff: config.FlexibleForceNew(vpcepServiceConnectionUpdateNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"service_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"endpoint_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func resourceVPCEPServiceConnectionUpdateCreateOrUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	serviceId := d.Get("service_id").(string)
	endpointId := d.Get("endpoint_id").(string)

	createHttpUrl := "v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections/description"
	createPath := client.Endpoint + createHttpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpc_endpoint_service_id}", serviceId)
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		JSONBody: map[string]interface{}{
			"connections": []map[string]interface{}{
				{
					"id":          endpointId,
					"description": d.Get("description"),
				},
			},
		},
	}

	_, err = client.Request("PUT", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error updating VPCEP connection description: %s", err)
	}

	if d.IsNewResource() {
		d.SetId(serviceId + "/" + endpointId)
	}

	return resourceVPCEPServiceConnectionUpdateRead(ctx, d, meta)
}

func resourceVPCEPServiceConnectionUpdateRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.VPCEPClient(region)
	if err != nil {
		return diag.Errorf("error creating VPC endpoint client: %s", err)
	}

	parts := strings.Split(d.Id(), "/")
	if len(parts) != 2 {
		return diag.Errorf("invalid ID format, want '<service_id>/<endpoint_id>', but got '%s'", d.Id())
	}
	serviceId := parts[0]
	endpointId := parts[1]

	getHttpUrl := "v1/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/connections?limit=1000&id={endpoint_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{vpc_endpoint_service_id}", serviceId)
	getPath = strings.ReplaceAll(getPath, "{endpoint_id}", endpointId)

	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving VPC endpoint service connection")
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return diag.Errorf("error flattening VPC endpoint service connection response: %s", err)
	}

	connection := utils.PathSearch("connections|[0]", getRespBody, nil)
	if connection == nil {
		return common.CheckDeletedDiag(d, golangsdk.ErrDefault404{}, "error finding VPC endpoint service connection from API response")
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("service_id", serviceId),
		d.Set("endpoint_id", endpointId),
		d.Set("description", utils.PathSearch("description", connection, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceVPCEPServiceConnectionUpdateDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting VPC endpoint service connection update resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
