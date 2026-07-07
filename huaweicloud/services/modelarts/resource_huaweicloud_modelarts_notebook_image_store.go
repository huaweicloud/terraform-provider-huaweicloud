package modelarts

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var (
	notebookImageStoreNonUpdatableParams = []string{
		"notebook_id",
		"name",
		"namespace",
		"tag",
		"description",
		"workspace_id",
		"swr_instance_id",
		"swr_instance_domain",
	}
	storedNotebookImageNotFoundErrCodes = []string{
		"ModelArts.6404",
	}
)

// @API ModelArts POST /v1/{project_id}/notebooks/{id}/create-image
// @API ModelArts GET /v1/{project_id}/images/{id}
func ResourceNotebookImageStore() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceNotebookImageStoreCreate,
		ReadContext:   resourceNotebookImageStoreRead,
		UpdateContext: resourceNotebookImageStoreUpdate,
		DeleteContext: resourceNotebookImageStoreDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(notebookImageStoreNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: `The region where the notebook is located.`,
			},

			// Required parameters.
			"notebook_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The ID of the notebook instance.`,
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The name of the image.`,
			},
			"namespace": {
				Type:        schema.TypeString,
				Required:    true,
				Description: `The namespace of the image.`,
			},

			// Optional parameters.
			"tag": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The tag of the image.`,
			},
			"description": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: `The description of the image.`,
			},
			"workspace_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `The workspace ID to which the image belongs.`,
			},

			// Attributes.
			"swr_path": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The SWR path of the image.`,
			},
			"arch": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The processor architecture type supported by the image.`,
			},
			"origin": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The origin of the image.`,
			},
			"resource_categories": {
				Type:        schema.TypeList,
				Computed:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: `The resource categories supported by the image.`,
			},
			"service_type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The service type supported by the image.`,
			},
			"visibility": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The visibility of the image.`,
			},
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The status of the image.`,
			},
			"type": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The type of the image.`,
			},
			"created_at": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `The creation time of the image, in RFC3339 format.`,
			},

			// Internal parameters.
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description: utils.SchemaDesc(
					`Whether to allow parameters that do not support changes to have their change-triggered behavior set to 'ForceNew'.`,
					utils.SchemaDescInput{
						Internal: true,
						Required: true,
					}),
			},
		},
	}
}

func buildNotebookImageStoreBodyParams(d *schema.ResourceData) map[string]interface{} {
	return map[string]interface{}{
		"name":         d.Get("name"),
		"namespace":    d.Get("namespace"),
		"tag":          utils.ValueIgnoreEmpty(d.Get("tag")),
		"description":  utils.ValueIgnoreEmpty(d.Get("description")),
		"workspace_id": utils.ValueIgnoreEmpty(d.Get("workspace_id")),
	}
}

func createNotebookImageStore(client *golangsdk.ServiceClient, notebookId string, d *schema.ResourceData) (interface{}, error) {
	httpURL := "v1/{project_id}/notebooks/{id}/create-image"

	createPath := client.Endpoint + httpURL
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{id}", notebookId)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildNotebookImageStoreBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(createResp)
}

func refreshNotebookImageStoreStatusFunc(client *golangsdk.ServiceClient, imageId string) retry.StateRefreshFunc {
	return func() (interface{}, string, error) {
		respBody, err := getNotebookImageStoreById(client, imageId)
		if err != nil {
			return nil, "ERROR", err
		}

		status := utils.PathSearch("status", respBody, "").(string)
		unexpectedStatus := []string{
			"CREATE_FAILED",
			"ERROR",
		}
		if utils.StrSliceContains(unexpectedStatus, status) {
			return respBody, "ERROR", fmt.Errorf("unexpected status (%s)", status)
		}

		if status == "ACTIVE" {
			return respBody, "COMPLETED", nil
		}
		return "continue", "PENDING", nil
	}
}

func waitForNotebookImageStoreStatusCompleted(ctx context.Context, client *golangsdk.ServiceClient, imageId string, timeout time.Duration) error {
	stateConf := &retry.StateChangeConf{
		Pending:      []string{"PENDING"},
		Target:       []string{"COMPLETED"},
		Refresh:      refreshNotebookImageStoreStatusFunc(client, imageId),
		Timeout:      timeout,
		Delay:        10 * time.Second,
		PollInterval: 10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	return err
}

func resourceNotebookImageStoreCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg        = meta.(*config.Config)
		region     = cfg.GetRegion(d)
		notebookId = d.Get("notebook_id").(string)
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := createNotebookImageStore(client, notebookId, d)
	if err != nil {
		return diag.Errorf("error saving image from notebook (%s): %s", notebookId, err)
	}

	imageId := utils.PathSearch("id", respBody, "").(string)
	if imageId == "" {
		return diag.Errorf("unable to find the image ID from the API response")
	}
	d.SetId(imageId)

	if err = waitForNotebookImageStoreStatusCompleted(ctx, client, imageId, d.Timeout(schema.TimeoutCreate)); err != nil {
		return diag.Errorf("error waiting for the image (%s) saved from notebook (%s) to complete: %s", imageId, notebookId, err)
	}

	return resourceNotebookImageStoreRead(ctx, d, meta)
}

func getNotebookImageStoreById(client *golangsdk.ServiceClient, imageId string) (interface{}, error) {
	httpUrl := "v1/{project_id}/images/{id}"
	getPath := client.Endpoint + httpUrl
	getPath = strings.ReplaceAll(getPath, "{project_id}", client.ProjectID)
	getPath = strings.ReplaceAll(getPath, "{id}", imageId)

	opt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	requestResp, err := client.Request("GET", getPath, &opt)
	if err != nil {
		return nil, err
	}

	return utils.FlattenResponse(requestResp)
}

func resourceNotebookImageStoreRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var (
		cfg     = meta.(*config.Config)
		region  = cfg.GetRegion(d)
		imageId = d.Id()
	)

	client, err := cfg.NewServiceClient("modelarts", region)
	if err != nil {
		return diag.Errorf("error creating ModelArts client: %s", err)
	}

	respBody, err := getNotebookImageStoreById(client, imageId)
	if err != nil {
		return common.CheckDeletedDiag(d,
			common.ConvertExpected400ErrInto404Err(err, "error_code", storedNotebookImageNotFoundErrCodes...),
			fmt.Sprintf("error retrieving ModelArts notebook image storage (%s)", imageId))
	}

	mErr := multierror.Append(nil,
		d.Set("region", region),
		// Required parameters.
		d.Set("name", utils.PathSearch("name", respBody, nil)),
		d.Set("namespace", utils.PathSearch("namespace", respBody, nil)),
		// Optional parameters.
		d.Set("workspace_id", utils.PathSearch("workspace_id", respBody, nil)),
		d.Set("description", utils.PathSearch("description", respBody, nil)),
		d.Set("tag", utils.PathSearch("tag", respBody, nil)),
		// Attributes.
		d.Set("swr_path", utils.PathSearch("swr_path", respBody, nil)),
		d.Set("arch", utils.PathSearch("arch", respBody, nil)),
		d.Set("origin", utils.PathSearch("origin", respBody, nil)),
		d.Set("resource_categories", utils.PathSearch("resource_categories", respBody, nil)),
		d.Set("service_type", utils.PathSearch("service_type", respBody, nil)),
		d.Set("visibility", utils.PathSearch("visibility", respBody, nil)),
		d.Set("status", utils.PathSearch("status", respBody, nil)),
		d.Set("type", utils.PathSearch("type", respBody, nil)),
		d.Set("created_at", utils.FormatTimeStampRFC3339(int64(utils.PathSearch("create_at", respBody, float64(0)).(float64))/1000, false)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func resourceNotebookImageStoreUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceNotebookImageStoreDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := `This resource is only a one-time action resource for storing the notebook image in ModelArts side.
Deleting this resource will not clear the corresponding image (version) record, but will only remove the resource
information from the tfstate file.`
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
