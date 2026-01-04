package dataarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/go-uuid"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var exclusiveApiNotPublishErrors = []string{
	"DLM.4001", // Workspace is not exist.
	"DLM.4018", // Api is not exist.
	"DLM.4215", // Instance is not exist.
}

// @API DataArtsStudio POST /v1/{project_id}/service/apis/{api_id}/instances/{instance_id}/publish
// @API DataArtsStudio GET /v1/{project_id}/service/apis/{api_id}/publish-info
func ResourceDataServiceApiPublishment() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDataServiceApiPublishmentCreate,
		ReadContext:   resourceDataServiceApiPublishmentRead,
		DeleteContext: resourceDataServiceApiPublishmentDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the API to be published is located.`,
			},

			// Parameter in request header
			"workspace_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The workspace ID to which the published API belongs.`,
			},

			// Arguments
			"api_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The ID of the API to be published.`,
			},
			"instance_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: `The exclusive cluster ID to which the published API belongs in Data Service side.`,
			},
			"apig_type": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The type of the APIG object.`,
			},
			"apig_instance_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The  APIG instance ID to which the API is published simultaneously in APIG service.`,
			},
			"apig_group_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The APIG group ID to which the published API belongs.`,
			},
			"roma_app_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: `The application ID for ROMA APIC.`,
			},

			// Attributes
			"publish_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The publish ID.`,
			},
		},
	}
}

func resourceDataServiceApiPublishmentCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = publishApi(client, d)
	if err != nil {
		return diag.Errorf("error publishing API: %s", err)
	}

	resourceId, err := uuid.GenerateUUID()
	if err != nil {
		return diag.Errorf("unable to generate ID: %s", err)
	}
	d.SetId(resourceId)

	return resourceDataServiceApiPublishmentRead(ctx, d, meta)
}

func queryApiPublishInfo(client *golangsdk.ServiceClient, workspaceId, apiId string) ([]interface{}, error) {
	var (
		httpUrl = "v1/{project_id}/service/apis/{api_id}/publish-info"
	)

	queryPath := client.Endpoint + httpUrl
	queryPath = strings.ReplaceAll(queryPath, "{project_id}", client.ProjectID)
	queryPath = strings.ReplaceAll(queryPath, "{api_id}", apiId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json;charset=utf8",
			"workspace":    workspaceId,
			"Dlm-Type":     "EXCLUSIVE",
		},
	}

	requestResp, err := client.Request("GET", queryPath, &opt)
	if err != nil {
		return nil, common.ConvertExpected400ErrInto404Err(err, "error_code", exclusiveApiNotPublishErrors...)
	}
	respBody, err := utils.FlattenResponse(requestResp)
	if err != nil {
		return nil, err
	}
	return utils.PathSearch("publish_messages", respBody, make([]interface{}, 0)).([]interface{}), nil
}

// QueryApiPublishInfoByInstanceId is a method that used to query the publish record in the specified cluster ID.
// Returns a 404 error if no publish record found.
func QueryApiPublishInfoByInstanceId(client *golangsdk.ServiceClient, workspaceId, apiId,
	instanceId string) (interface{}, error) {
	publishRecords, err := queryApiPublishInfo(client, workspaceId, apiId)
	if err != nil {
		return nil, err
	}
	publishRecord := utils.PathSearch(fmt.Sprintf("[?instance_id=='%s']|[0]", instanceId), publishRecords, nil)
	// Since this resource does not provide the import function, the release operation must have been performed before entering the query stage.
	// So, the 'API_STATUS_CREATED' status does not exist.
	// If the status response is empty (means not found) or value 'API_STATUS_OFFLINE', that means the API has been unpublished.
	apiStatus := utils.PathSearch("api_status", publishRecord, "").(string)
	if apiStatus == "" || apiStatus == "API_STATUS_OFFLINE" {
		return nil, golangsdk.ErrDefault404{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/service/apis/{api_id}/publish-info",
				RequestId: "NONE",
				Body:      []byte(fmt.Sprintf("the API has been unpublished (%s)", apiId)),
			},
		}
	}
	if apiStatus != "API_STATUS_PUBLISHED" {
		return nil, golangsdk.ErrDefault500{
			ErrUnexpectedResponseCode: golangsdk.ErrUnexpectedResponseCode{
				Method:    "GET",
				URL:       "/v1/{project_id}/service/apis/{api_id}/publish-info",
				RequestId: "NONE",
				Body: []byte(fmt.Sprintf("the API status is not in expect, want 'API_STATUS_PUBLISHED', "+
					"but got '%s'", apiStatus)),
			},
		}
	}
	return publishRecord, nil
}

func resourceDataServiceApiPublishmentRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Get("api_id").(string)
		instanceId  = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	publishRecord, err := QueryApiPublishInfoByInstanceId(client, workspaceId, apiId, instanceId)
	if err != nil {
		return common.CheckDeletedDiag(d, err, fmt.Sprintf("error getting publish information for API (%s)", d.Id()))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		d.Set("publish_id", utils.PathSearch("id", publishRecord, nil)),
	)
	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceDataServiceApiPublishmentDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg    = meta.(*config.Config)
		region = cfg.GetRegion(d)

		workspaceId = d.Get("workspace_id").(string)
		apiId       = d.Get("api_id").(string)
		instanceId  = d.Get("instance_id").(string)
	)
	client, err := cfg.NewServiceClient("dataarts", region)
	if err != nil {
		return diag.Errorf("error creating DataArts Studio client: %s", err)
	}

	err = doApiAction(client, workspaceId, apiId, instanceId, "UNPUBLISH")
	if err != nil {
		// If the workspace or API has been removed before unpublish request sent, remove the tfstate record for this resource.
		parsedErr := common.ConvertExpected400ErrInto404Err(err, "error_code", exclusiveApiNotPublishErrors...)
		if _, ok := parsedErr.(golangsdk.ErrDefault404); ok {
			return common.CheckDeletedDiag(d, parsedErr, "")
		}
		// The DLM.4182 error code (this action is illegal) means that the API may have been offline before.
		parsedErr = common.ConvertExpected400ErrInto404Err(err, "error_code", "DLM.4182")
		_, ok := parsedErr.(golangsdk.ErrDefault404)
		if !ok {
			// If not, means the unpublish request triggers an error.
			return diag.Errorf(`error unpublishing API (%s): %s`, apiId, err)
		}
		// If the API has been unpublished before unpublish request sent, remove the tfstate record for this resource.
		_, queryErr := QueryApiPublishInfoByInstanceId(client, workspaceId, apiId, instanceId)
		if err != nil {
			if _, ok := queryErr.(golangsdk.ErrDefault404); ok {
				return common.CheckDeletedDiag(d, queryErr, "")
			}
		}
	}
	return nil
}
