package css

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/thesaurus"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API CSS DELETE /v1.0/{project_id}/clusters/{clusterId}/thesaurus
// @API CSS GET /v1.0/{project_id}/clusters/{clusterId}/thesaurus
// @API CSS POST /v1.0/{project_id}/clusters/{clusterId}/thesaurus
func ResourceCssthesaurus() *schema.Resource {
	return &schema.Resource{
		CreateContext: ResourceCssthesaurusCreate,
		ReadContext:   ResourceCssthesaurusRead,
		UpdateContext: ResourceCssthesaurusUpdate,
		DeleteContext: ResourceCssthesaurusDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

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
				ForceNew: true,
			},
			"bucket_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"main_object": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				AtLeastOneOf: []string{"main_object", "stop_object", "synonym_object"},
			},
			"stop_object": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"synonym_object": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func ResourceCssthesaurusCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	clusterID := d.Get("cluster_id").(string)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	addCustomThesaurusHttpUrl := "v1.0/{project_id}/clusters/{cluster_id}/thesaurus"
	addCustomThesaurusPath := cssV1Client.Endpoint + addCustomThesaurusHttpUrl
	addCustomThesaurusPath = strings.ReplaceAll(addCustomThesaurusPath, "{project_id}", cssV1Client.ProjectID)
	addCustomThesaurusPath = strings.ReplaceAll(addCustomThesaurusPath, "{cluster_id}", clusterID)

	addCustomThesaurusOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}
	addCustomThesaurusOpt.JSONBody = buildThesaurusCreateParameters(d)
	_, err = cssV1Client.Request("POST", addCustomThesaurusPath, &addCustomThesaurusOpt)
	if err != nil {
		return diag.Errorf("error adding CSS cluster custom thesaurus: %s", err)
	}

	d.SetId(clusterID)

	createResultErr := checkThesaurusLoadResult(ctx, cssV1Client, clusterID, d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	return ResourceCssthesaurusRead(ctx, d, meta)
}

func ResourceCssthesaurusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceCssthesaurusCreate(ctx, d, meta)
}

func buildThesaurusCreateParameters(d *schema.ResourceData) map[string]interface{} {
	opts := map[string]interface{}{
		"bucketName": d.Get("bucket_name").(string),
	}

	// nil or "Default" indicates no change, "" or "Unused" indicates that this value is cleared.
	if obj, ok := d.GetOk("main_object"); ok && obj.(string) != "Default" {
		opts["mainObject"] = convertUnusedToNullString(obj.(string))
	}
	if obj, ok := d.GetOk("stop_object"); ok && obj.(string) != "Default" {
		opts["stopObject"] = convertUnusedToNullString(obj.(string))
	}
	if obj, ok := d.GetOk("synonym_object"); ok && obj.(string) != "Default" {
		opts["synonymObject"] = convertUnusedToNullString(obj.(string))
	}

	return opts
}

func convertUnusedToNullString(str string) string {
	if str == "Unused" {
		return ""
	}
	return str
}

func ResourceCssthesaurusRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	detail, err := thesaurus.Get(cssV1Client, d.Id())
	if err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error querying CSS cluster thesaurus")
	}

	mErr := multierror.Append(
		d.Set("cluster_id", detail.ClusterId),
		d.Set("bucket_name", detail.Bucket),
		d.Set("main_object", detail.MainObj),
		d.Set("stop_object", detail.StopObj),
		d.Set("synonym_object", detail.SynonymObj),
		d.Set("status", detail.Status),
		d.Set("update_time", time.Unix(int64(detail.UpdateTime/1000), 0).UTC().Format(time.RFC3339)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func ResourceCssthesaurusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	conf := meta.(*config.Config)
	region := conf.GetRegion(d)
	cssV1Client, err := conf.CssV1Client(region)
	if err != nil {
		return diag.Errorf("error creating CSS V1 client: %s", err)
	}

	clusterId := d.Id()

	errResult := thesaurus.Delete(cssV1Client, clusterId)
	if errResult.Err != nil {
		// "CSS.0015": The cluster does not exist. Status code is 403.
		err = common.ConvertExpected403ErrInto404Err(errResult.Err, "errCode", "CSS.0015")
		return common.CheckDeletedDiag(d, err, "error deleting CSS cluster thesaurus")
	}

	errCheckRt := checkThesaurusDeleteResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return diag.Errorf("failed to check the result of deletion %s", errCheckRt)
	}
	return nil
}

func checkThesaurusLoadResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Loading"},
		Target:  []string{"Loaded"},
		Refresh: func() (interface{}, string, error) {
			resp, err := thesaurus.Get(client, clusterId)
			if err != nil {
				return nil, "failed", err
			}
			if resp.Status == "Failed" {
				return nil, "failed", fmt.Errorf("load thesaurus failed in cluster_id: %s", clusterId)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		return fmt.Errorf("error waiting for CSS (%s) to load thesaurus: %s", clusterId, err)
	}
	return nil
}

func checkThesaurusDeleteResult(ctx context.Context, client *golangsdk.ServiceClient, clusterId string,
	timeout time.Duration) error {
	stateConf := &resource.StateChangeConf{
		Pending: []string{"Pending"},
		Target:  []string{"Done"},
		Refresh: func() (interface{}, string, error) {
			resp, err := thesaurus.Get(client, clusterId)
			if err != nil {
				if _, ok := err.(golangsdk.ErrDefault404); ok {
					return nil, "Done", nil
				}
				return nil, "failed", err
			}
			if resp != nil && resp.MainObj == "Unused" && resp.StopObj == "Unused" && resp.SynonymObj == "Unused" {
				return resp, "Done", nil
			}
			return resp, "Pending", nil
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		return fmt.Errorf("error waiting for CSS thesaurus (%s) to be delete: %s", clusterId, err)
	}
	return nil
}
