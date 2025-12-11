package elb

import (
	"context"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var recycleBinLoadBalancerDeleteNonUpdatableParams = []string{"loadbalancer_id"}

// @API ELB DELETE /v3/{project_id}/elb/recycle-bin/loadbalancers/{loadbalancer_id}
func ResourceElbRecycleBinLoadBalancerDelete() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceElbRecycleBinLoadBalancerDeleteCreate,
		ReadContext:   resourceElbRecycleBinLoadBalancerDeleteRead,
		UpdateContext: resourceElbRecycleBinLoadBalancerDeleteUpdate,
		DeleteContext: resourceElbRecycleBinLoadBalancerDeleteDelete,

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
		},

		CustomizeDiff: config.FlexibleForceNew(recycleBinLoadBalancerDeleteNonUpdatableParams),

		Schema: map[string]*schema.Schema{
			"region": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
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

func resourceElbRecycleBinLoadBalancerDeleteCreate(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	region := cfg.GetRegion(d)

	var (
		httpUrl = "v3/{project_id}/elb/recycle-bin/loadbalancers/{loadbalancer_id}"
		product = "elb"
	)
	client, err := cfg.NewServiceClient(product, region)
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	createPath := client.Endpoint + httpUrl
	createPath = strings.ReplaceAll(createPath, "{project_id}", client.ProjectID)
	createPath = strings.ReplaceAll(createPath, "{loadbalancer_id}", d.Get("loadbalancer_id").(string))

	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
	}

	_, err = client.Request("DELETE", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating ELB recycle bin load balancer delete: %s", err)
	}

	d.SetId(d.Get("loadbalancer_id").(string))

	return nil
}

func resourceElbRecycleBinLoadBalancerDeleteRead(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceElbRecycleBinLoadBalancerDeleteUpdate(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	return nil
}

func resourceElbRecycleBinLoadBalancerDeleteDelete(_ context.Context, _ *schema.ResourceData, _ interface{}) diag.Diagnostics {
	errorMsg := "Deleting ELB recycle bin load balancer delete resource is not supported. The resource is only removed" +
		" from the state."
	return diag.Diagnostics{
		diag.Diagnostic{
			Severity: diag.Warning,
			Summary:  errorMsg,
		},
	}
}
