package vpcep

import (
	"context"
	"strings"

	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var endpointServiceResourceAttachNonUpdatableParams = []string{"service_id", "server_resources"}

// @API VPCEP POST /v2/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/add-server-resources
func ResourceEndpointServiceResourceAttach() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceEndpointServiceResourceAttachCreate,
		ReadContext:   resourceEndpointServiceResourceAttachRead,
		UpdateContext: resourceEndpointServiceResourceAttachUpdate,
		DeleteContext: resourceEndpointServiceResourceAttachDelete,

		CustomizeDiff: config.FlexibleForceNew(endpointServiceResourceAttachNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the VPCEP endpoint service is located.`,
			},

			// Required parameters.
			"service_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the VPCEP endpoint service.`,
			},
			"server_resources": {
				Type:        schema.TypeList,
				Required:    true,
				Elem:        endpointServiceResourceAttachServerResourceSchema(),
				Description: `The list of server resources to be added to the VPCEP endpoint service.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
		},
	}
}

func endpointServiceResourceAttachServerResourceSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"resource_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the server resource.`,
			},
			"availability_zone_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The availability zone ID of the server resource.`,
			},
		},
	}
}

func buildEndpointServiceResourceAttachServerResources(serverResources []interface{}) []map[string]interface{} {
	if len(serverResources) < 1 {
		return nil
	}

	result := make([]map[string]interface{}, 0, len(serverResources))
	for _, serverResource := range serverResources {
		result = append(result, map[string]interface{}{
			"resource_id":          utils.PathSearch("resource_id", serverResource, nil),
			"availability_zone_id": utils.PathSearch("availability_zone_id", serverResource, nil),
		})
	}

	return result
}

func buildEndpointServiceResourceAttachBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"server_resources": buildEndpointServiceResourceAttachServerResources(d.Get("server_resources").([]interface{})),
	}
}

func resourceEndpointServiceResourceAttachCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg       = meta.(*config.Config)
		region    = cfg.GetRegion(d)
		httpUrl   = "v2/{project_id}/vpc-endpoint-services/{vpc_endpoint_service_id}/add-server-resources"
		serviceId = d.Get("service_id").(string)
	)

	client, err := cfg.NewServiceClient("vpcep", region)
	if err != nil {
		return diag.Errorf("error creating VPCEP client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{vpc_endpoint_service_id}", serviceId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
		JSONBody:         buildEndpointServiceResourceAttachBodyParams(d),
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error adding server resources to endpoint service: %s", err)
	}

	randomUUID, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(randomUUID)

	return resourceEndpointServiceResourceAttachRead(ctx, d, meta)
}

func resourceEndpointServiceResourceAttachRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEndpointServiceResourceAttachUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceEndpointServiceResourceAttachDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for adding server resources to VPCEP service.
Deleting this resource will not remove the server resources from the VPCEP service, but will only remove the
resource information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
