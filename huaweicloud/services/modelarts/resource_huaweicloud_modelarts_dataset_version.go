package modelarts

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/dataset"
	"github.com/chnsz/golangsdk/openstack/modelarts/v2/version"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API ModelArts DELETE /v2/{project_id}/datasets/{datasetId}/versions/{versionId}
// @API ModelArts GET /v2/{project_id}/datasets/{datasetId}/versions/{versionId}
// @API ModelArts POST /v2/{project_id}/datasets/{datasetId}/versions
func ResourceDatasetVersion() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceDatasetVersionCreate,
		ReadContext:   ResourceDatasetVersionRead,
		DeleteContext: ResourceDatasetVersionDelete,
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
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"dataset_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"split_ratio": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Default:  "1.00",
			},

			"hard_example": {
				Type:     schema.TypeBool,
				Optional: true,
				ForceNew: true,
				Default:  false,
			},

			"version_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"verification": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"labeling_type": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"files": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"storage_path": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"is_current": {
				Type:     schema.TypeBool,
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
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func ResourceDatasetVersionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId := d.Get("dataset_id").(string)
	opts := version.CreateOpts{
		VersionName:              d.Get("name").(string),
		Description:              d.Get("description").(string),
		ClearHardProperty:        utils.Bool(!d.Get("hard_example").(bool)),
		TrainEvaluateSampleRatio: d.Get("split_ratio").(string),
	}

	rst, err := version.Create(client, datasetId, opts)
	if err != nil {
		return diag.Errorf("error creating ModelArts dataset version: %s", err)
	}

	d.SetId(fmt.Sprintf("%s/%s", datasetId, rst.VersionId))

	err = waitingforDatasetVersionCreated(ctx, client, datasetId, rst.VersionId, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return diag.FromErr(err)
	}

	return ResourceDatasetVersionRead(ctx, d, meta)
}

func ResourceDatasetVersionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId, versionId, err := ParseVersionInfoFromId(d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	detail, err := version.Get(client, datasetId, versionId)
	if err != nil {
		return common.CheckDeletedDiag(d, parseDatasetVersionErrorToError404(err),
			"error retrieving ModelArts dataset version")
	}

	mErr := multierror.Append(
		d.Set("region", region),
		d.Set("name", detail.VersionName),
		d.Set("version_id", detail.VersionId),
		d.Set("dataset_id", datasetId),
		d.Set("description", detail.Description),
		d.Set("split_ratio", detail.TrainEvaluateSampleRatio),
		d.Set("status", detail.Status),
		d.Set("verification", detail.DataValidate),
		d.Set("labeling_type", detail.LabelType),
		d.Set("files", detail.TotalSampleCount),
		d.Set("storage_path", detail.ManifestPath),
		d.Set("is_current", detail.IsCurrent),
		d.Set("created_at", utils.FormatTimeStampUTC(int64(detail.CreateTime))),
		d.Set("updated_at", utils.FormatTimeStampUTC(int64(detail.UpdateTime))),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceDatasetVersionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)
	client, err := cfg.ModelArtsV2Client(region)
	if err != nil {
		return diag.Errorf("error creating ModelArts v2 client, err=%s", err)
	}

	datasetId := d.Get("dataset_id").(string)
	versionId := d.Get("version_id").(string)
	dErr := version.Delete(client, datasetId, versionId)
	if dErr.Err != nil {
		return diag.Errorf("error deleting ModelArts dataset version, ID= %s", d.Id())
	}

	return nil
}

func waitingforDatasetVersionCreated(ctx context.Context, client *golangsdk.ServiceClient, datasetId, versionId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"0"},
		Target:  []string{"1"},
		Refresh: func() (interface{}, string, error) {
			detail, err := version.Get(client, datasetId, versionId)
			if err != nil {
				return nil, "", err
			}
			return detail, fmt.Sprint(detail.Status), nil
		},
		Timeout:      timeout,
		PollInterval: 10 * time.Second,
		Delay:        10 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for ModelArts dataset %s version (%s) to be created: %s",
			datasetId, versionId, err)
	}
	return nil
}

func parseDatasetVersionErrorToError404(respErr error) error {
	var apiError dataset.CreateResp
	if errCode, ok := respErr.(golangsdk.ErrDefault400); ok {
		pErr := json.Unmarshal(errCode.Body, &apiError)
		if pErr == nil && (apiError.ErrorCode == "ModelArts.4352" || apiError.ErrorCode == "ModelArts.4353") {
			return golangsdk.ErrDefault404(errCode)
		}
	}
	return respErr
}

func ParseVersionInfoFromId(id string) (datasetId string, versionId string, err error) {
	idParts := strings.Split(id, "/")
	if len(idParts) != 2 {
		err = fmt.Errorf("invalid format specified for dataset version. Format must be <dataset id>/<version id>")
		return
	}
	datasetId = idParts[0]
	versionId = idParts[1]
	return
}
