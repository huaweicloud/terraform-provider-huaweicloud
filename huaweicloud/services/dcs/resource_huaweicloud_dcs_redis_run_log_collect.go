package dcs

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

var redisRunLogCollectNonUpdatableParams = []string{
	"instance_id",
	"log_type",
	"query_time",
	"replication_id",
}

// @API DCS POST /v2/{project_id}/instances/{instance_id}/redislog
func ResourceDcsRedisRunLogCollect() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDcsRedisRunLogCollectCreate,
		ReadContext:   resourceDcsRedisRunLogCollectRead,
		UpdateContext: resourceDcsRedisRunLogCollectUpdate,
		DeleteContext: resourceDcsRedisRunLogCollectDelete,

		CustomizeDiff: config.FlexibleForceNew(redisRunLogCollectNonUpdatableParams),

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
			},
			"log_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_time": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"replication_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceDcsRedisRunLogCollectCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v2/{project_id}/instances/{instance_id}/redislog"
		product = "dcs"
	)

	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating DCS client: %s", err)
	}

	instanceId := d.Get("instance_id").(string)
	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{instance_id}", instanceId)

	createPath = buildCreateRedisRunLogCollectQueryParams(createPath, d)

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		OkCodes:          []int{204},
		MoreHeaders:      map[string]string{"Content-Type": "application/json"},
	}

	_, err = client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating DCS redis run log collect: %s", err)
	}

	d.SetId(instanceId)

	return nil
}

func buildCreateRedisRunLogCollectQueryParams(url string, d *schema.ResourceData) string {
	queryParams := make([]string, 0)

	if v, ok := d.GetOk("query_time"); ok {
		queryParams = append(queryParams, fmt.Sprintf("query_time=%d", v.(int)))
	}

	if v, ok := d.GetOk("log_type"); ok {
		queryParams = append(queryParams, fmt.Sprintf("log_type=%s", v.(string)))
	}

	if v, ok := d.GetOk("replication_id"); ok {
		queryParams = append(queryParams, fmt.Sprintf("replication_id=%s", v.(string)))
	}

	if len(queryParams) > 0 {
		url = url + "?" + strings.Join(queryParams, "&")
	}

	return url
}

func resourceDcsRedisRunLogCollectRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsRedisRunLogCollectUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceDcsRedisRunLogCollectDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting DCS redis run log collect resource is not supported. The resource is only removed from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
