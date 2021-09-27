package css

import (
	"context"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/css/v1/thesaurus"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
)

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
				AtLeastOneOf: []string{"main_object", "stop_object", "synonym_object"},
			},
			"stop_object": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"synonym_object": {
				Type:     schema.TypeString,
				Optional: true,
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
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating CSS V1 client: %s", err)
	}
	opts := buildThesaurusCreateParameters(d)
	clusterId := d.Get("cluster_id").(string)

	loadErr := thesaurus.Load(cssV1Client, clusterId, *opts)
	if loadErr.Err != nil {
		return fmtp.DiagErrorf("load thesaurus to css cluster failed. cluster_id=%s,error=%s", clusterId, loadErr.Err)
	}

	d.SetId(clusterId)

	createResultErr := checkThesaurusLoadResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutCreate))
	if createResultErr != nil {
		return diag.FromErr(createResultErr)
	}

	return ResourceCssthesaurusRead(ctx, d, meta)
}

func ResourceCssthesaurusUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return ResourceCssthesaurusCreate(ctx, d, meta)
}

func buildThesaurusCreateParameters(d *schema.ResourceData) *thesaurus.LoadThesaurusReq {
	opts := thesaurus.LoadThesaurusReq{
		BucketName: d.Get("bucket_name").(string),
	}

	if obj, ok := d.GetOk("main_object"); ok {
		opts.MainObject = obj.(string)
	}
	if obj, ok := d.GetOk("stop_object"); ok {
		opts.StopObject = obj.(string)
	}
	if obj, ok := d.GetOk("synonym_object"); ok {
		opts.SynonymObject = obj.(string)
	}

	return &opts
}

func ResourceCssthesaurusRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating CSS V1 client: %s", err)
	}

	detail, err := thesaurus.Get(cssV1Client, d.Id())
	if err != nil {
		return fmtp.DiagErrorf("Query cluster thesaurus failed,cluster_id=%s,err=%s", d.Id(), err)
	}

	mErr := multierror.Append(
		d.Set("cluster_id", detail.ClusterId),
		d.Set("bucket_name", detail.Bucket),
		d.Set("main_object", detail.MainObj),
		d.Set("stop_object", detail.StopObj),
		d.Set("stop_object", detail.StopObj),
		d.Set("synonym_object", detail.SynonymObj),
		d.Set("status", detail.Status),
		d.Set("update_time", time.Unix(int64(detail.UpdateTime/1000), 0).UTC().Format(time.RFC3339)),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.DiagErrorf("Error setting vault fields: %s", err)
	}

	return nil
}

func ResourceCssthesaurusDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	config := meta.(*config.Config)
	region := config.GetRegion(d)
	cssV1Client, err := config.CssV1Client(region)
	if err != nil {
		return fmtp.DiagErrorf("Error creating CSS V1 client: %s", err)
	}

	clusterId := d.Id()

	errResult := thesaurus.Delete(cssV1Client, clusterId)
	if errResult.Err != nil {
		return fmtp.DiagErrorf("Delete CSS Cluster thesaurus failed. %s", errResult.Err)
	}

	errCheckRt := checkThesaurusDeleteResult(ctx, cssV1Client, clusterId, d.Timeout(schema.TimeoutDelete))
	if errCheckRt != nil {
		return fmtp.DiagErrorf("Failed to check the result of deletion %s", errCheckRt)
	}
	d.SetId("")
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
				return nil, "failed", fmtp.Errorf("load thesaurus failed in cluster_id=%s", clusterId)
			}
			return resp, resp.Status, err
		},
		Timeout:      timeout,
		PollInterval: 10 * timeout,
		Delay:        10 * time.Second,
	}
	_, err := stateConf.WaitForStateContext(ctx)

	if err != nil {
		return fmtp.Errorf("error waiting for CSS (%s) to load thesaurus: %s", clusterId, err)
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
			if resp != nil && resp.MainObj == "" && resp.StopObj == "" && resp.SynonymObj == "" {
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
		return fmtp.Errorf("error waiting for CSS thesaurus (%s) to be delete: %s", clusterId, err)
	}
	return nil
}
