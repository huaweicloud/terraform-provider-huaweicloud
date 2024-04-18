package elb

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/l7policies"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

// @API ELB POST /v3/{project_id}/elb/l7policies
// @API ELB GET /v3/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB PUT /v3/{project_id}/elb/l7policies/{l7policy_id}
// @API ELB DELETE /v3/{project_id}/elb/l7policies/{l7policy_id}
func ResourceL7PolicyV3() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceL7PolicyV3Create,
		ReadContext:   resourceL7PolicyV3Read,
		UpdateContext: resourceL7PolicyV3Update,
		DeleteContext: resourceL7PolicyV3Delete,
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
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"listener_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "REDIRECT_TO_POOL",
				ForceNew: true,
			},
			"priority": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"redirect_listener_id": {
				Type:     schema.TypeString,
				Optional: true,
				ExactlyOneOf: []string{"redirect_listener_id", "redirect_pool_id", "redirect_url_config",
					"fixed_response_config"},
				Computed: true,
			},
			"redirect_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"redirect_pools_extend_config": {
				Type:     schema.TypeList,
				Elem:     redirectPoolsExtendConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"redirect_url_config": {
				Type:     schema.TypeList,
				Elem:     redirectUrlConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"fixed_response_config": {
				Type:     schema.TypeList,
				Elem:     fixedResponseConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
			"provisioning_status": {
				Type:     schema.TypeString,
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
	}
}

func redirectPoolsExtendConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"rewrite_url_enabled": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"rewrite_url_config": {
				Type:     schema.TypeList,
				Elem:     rewriteUrlConfigSchema(),
				MaxItems: 1,
				Optional: true,
				Computed: true,
			},
		},
	}
	return &sc
}

func rewriteUrlConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func redirectUrlConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"protocol": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"host": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"port": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"path": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"query": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func fixedResponseConfigSchema() *schema.Resource {
	sc := schema.Resource{
		Schema: map[string]*schema.Schema{
			"status_code": {
				Type:     schema.TypeString,
				Required: true,
			},
			"content_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"message_body": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
	return &sc
}

func resourceL7PolicyV3Create(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	action := d.Get("action").(string)
	createOpts := l7policies.CreateOpts{
		Name:        d.Get("name").(string),
		Description: d.Get("description").(string),
		Action:      l7policies.Action(action),
		Priority:    int32(d.Get("priority").(int)),
		ListenerID:  d.Get("listener_id").(string),
	}
	if action == "REDIRECT_TO_POOL" {
		createOpts.RedirectPoolID = d.Get("redirect_pool_id").(string)
		createOpts.RedirectPoolsExtendConfig = buildRedirectPoolsExtendConfig(d)
	} else if action == "REDIRECT_TO_LISTENER" {
		createOpts.RedirectListenerID = d.Get("redirect_listener_id").(string)
	} else if action == "REDIRECT_TO_URL" {
		createOpts.RedirectUrlConfig = buildRedirectUrlConfig(d)
	} else {
		createOpts.FixedResponseConfig = buildFixedResponseConfig(d)
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	l7Policy, err := l7policies.Create(elbClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating L7 Policy: %s", err)
	}

	timeout := d.Timeout(schema.TimeoutCreate)
	// Wait for L7 Policy to become active before continuing
	err = waitForElbV3Policy(ctx, elbClient, l7Policy.ID, "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(l7Policy.ID)

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func buildRedirectPoolsExtendConfig(d *schema.ResourceData) *l7policies.RedirectPoolsExtendConfig {
	var redirectPoolsExtendConfig *l7policies.RedirectPoolsExtendConfig
	redirectPoolsExtendConfigRaw := d.Get("redirect_pools_extend_config").([]interface{})
	log.Printf("[DEBUG] redirectPoolsExtendConfigRaw: %+v", redirectPoolsExtendConfigRaw)
	if len(redirectPoolsExtendConfigRaw) == 1 {
		if v, ok := redirectPoolsExtendConfigRaw[0].(map[string]interface{}); ok {
			redirectPoolsExtendConfig = &l7policies.RedirectPoolsExtendConfig{
				RewriteUrlEnable: v["rewrite_url_enabled"].(bool),
				RewriteUrlConfig: buildRewriteUrlConfig(v["rewrite_url_config"]),
			}
		}
	}
	log.Printf("[DEBUG] redirectPoolsExtendConfig: %+v", redirectPoolsExtendConfig)
	return redirectPoolsExtendConfig
}

func buildRewriteUrlConfig(data interface{}) *l7policies.RewriteUrlConfig {
	var rewriteUrlConfig *l7policies.RewriteUrlConfig
	rewriteUrlConfigRaw := data.([]interface{})
	log.Printf("[DEBUG] rewriteUrlConfigRaw: %+v", rewriteUrlConfigRaw)
	if len(rewriteUrlConfigRaw) == 1 {
		if v, ok := rewriteUrlConfigRaw[0].(map[string]interface{}); ok {
			rewriteUrlConfig = &l7policies.RewriteUrlConfig{
				Host:  v["host"].(string),
				Path:  v["path"].(string),
				Query: v["query"].(string),
			}
		}
	}
	log.Printf("[DEBUG] rewriteUrlConfig: %+v", rewriteUrlConfig)
	return rewriteUrlConfig
}

func buildRedirectUrlConfig(d *schema.ResourceData) *l7policies.RedirectUrlConfig {
	var redirectUrlConfig *l7policies.RedirectUrlConfig
	redirectUrlConfigRaw := d.Get("redirect_url_config").([]interface{})
	log.Printf("[DEBUG] redirectUrlConfigRaw: %+v", redirectUrlConfigRaw)
	if len(redirectUrlConfigRaw) == 1 {
		if v, ok := redirectUrlConfigRaw[0].(map[string]interface{}); ok {
			redirectUrlConfig = &l7policies.RedirectUrlConfig{
				Protocol:   v["protocol"].(string),
				Host:       v["host"].(string),
				Port:       v["port"].(string),
				Path:       v["path"].(string),
				Query:      v["query"].(string),
				StatusCode: v["status_code"].(string),
			}
		}
	}
	log.Printf("[DEBUG] redirectUrlConfig: %+v", redirectUrlConfig)
	return redirectUrlConfig
}

func buildFixedResponseConfig(d *schema.ResourceData) *l7policies.FixedResponseConfig {
	var fixedResponseConfig *l7policies.FixedResponseConfig
	fixedResponseConfigRaw := d.Get("fixed_response_config").([]interface{})
	log.Printf("[DEBUG] fixedResponseConfigRaw: %+v", fixedResponseConfigRaw)
	if len(fixedResponseConfigRaw) == 1 {
		if v, ok := fixedResponseConfigRaw[0].(map[string]interface{}); ok {
			fixedResponseConfig = &l7policies.FixedResponseConfig{
				StatusCode:  v["status_code"].(string),
				ContentType: v["content_type"].(string),
				MessageBody: v["message_body"].(string),
			}
		}
	}
	log.Printf("[DEBUG] fixedResponseConfig: %+v", fixedResponseConfig)
	return fixedResponseConfig
}

func resourceL7PolicyV3Read(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	l7Policy, err := l7policies.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "L7 Policy")
	}

	log.Printf("[DEBUG] Retrieved L7 Policy %s: %#v", d.Id(), l7Policy)

	mErr := multierror.Append(nil,
		d.Set("description", l7Policy.Description),
		d.Set("name", l7Policy.Name),
		d.Set("action", l7Policy.Action),
		d.Set("priority", l7Policy.Priority),
		d.Set("listener_id", l7Policy.ListenerID),
		d.Set("redirect_pool_id", l7Policy.RedirectPoolID),
		d.Set("redirect_listener_id", l7Policy.RedirectListenerID),
		d.Set("redirect_pools_extend_config", flattenRedirectPoolsExtendConfig(l7Policy)),
		d.Set("redirect_url_config", flattenRedirectUrlConfig(l7Policy)),
		d.Set("fixed_response_config", flattenFixedResponseConfig(l7Policy)),
		d.Set("region", cfg.GetRegion(d)),
		d.Set("created_at", l7Policy.CreatedAt),
		d.Set("updated_at", l7Policy.UpdatedAt),
		d.Set("provisioning_status", l7Policy.ProvisioningStatus),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting Dedicated ELB l7policy fields: %s", err)
	}

	return nil
}

func flattenRedirectPoolsExtendConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectPoolsExtendConfig []map[string]interface{}
	if l7policy.RedirectPoolsExtendConfig != nil {
		redirectPoolsExtendConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["rewrite_url_enabled"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlEnable
		params["rewrite_url_config"] = flattenRewriteUrlConfig(l7policy)
		redirectPoolsExtendConfig[0] = params
	}
	return redirectPoolsExtendConfig
}

func flattenRewriteUrlConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var rewriteUrlConfig []map[string]interface{}
	if l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig != nil {
		rewriteUrlConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["host"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Host
		params["path"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Path
		params["query"] = l7policy.RedirectPoolsExtendConfig.RewriteUrlConfig.Query
		rewriteUrlConfig[0] = params
	}
	return rewriteUrlConfig
}

func flattenRedirectUrlConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var redirectUrlConfig []map[string]interface{}
	if l7policy.RedirectUrlConfig != nil {
		redirectUrlConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["protocol"] = l7policy.RedirectUrlConfig.Protocol
		params["host"] = l7policy.RedirectUrlConfig.Host
		params["port"] = l7policy.RedirectUrlConfig.Port
		params["path"] = l7policy.RedirectUrlConfig.Path
		params["query"] = l7policy.RedirectUrlConfig.Query
		params["status_code"] = l7policy.RedirectUrlConfig.StatusCode
		redirectUrlConfig[0] = params
	}
	return redirectUrlConfig
}

func flattenFixedResponseConfig(l7policy *l7policies.L7Policy) []map[string]interface{} {
	var fixedResponseConfig []map[string]interface{}
	if l7policy.FixedResponseConfig != nil {
		fixedResponseConfig = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["status_code"] = l7policy.FixedResponseConfig.StatusCode
		params["content_type"] = l7policy.FixedResponseConfig.ContentType
		params["message_body"] = l7policy.FixedResponseConfig.MessageBody
		fixedResponseConfig[0] = params
	}
	return fixedResponseConfig
}

func resourceL7PolicyV3Update(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	var updateOpts l7policies.UpdateOpts

	if d.HasChange("name") {
		name := d.Get("name").(string)
		updateOpts.Name = &name
	}
	if d.HasChange("priority") {
		priority := d.Get("priority").(int)
		updateOpts.Priority = int32(priority)
	}
	if d.HasChange("description") {
		description := d.Get("description").(string)
		updateOpts.Description = &description
	}
	if d.HasChange("redirect_pool_id") {
		redirectPoolID := d.Get("redirect_pool_id").(string)
		updateOpts.RedirectPoolID = &redirectPoolID
	}
	if d.HasChange("redirect_pools_extend_config") {
		updateOpts.RedirectPoolsExtendConfig = buildRedirectPoolsExtendConfig(d)
	}
	if d.HasChange("redirect_listener_id") {
		redirectListenerID := d.Get("redirect_listener_id").(string)
		updateOpts.RedirectListenerID = &redirectListenerID
	}
	if d.HasChange("redirect_url_config") {
		updateOpts.RedirectUrlConfig = buildRedirectUrlConfig(d)
	}
	if d.HasChange("fixed_response_config") {
		updateOpts.FixedResponseConfig = buildFixedResponseConfig(d)
	}

	log.Printf("[DEBUG] Updating L7 Policy %s with options: %#v", d.Id(), updateOpts)
	_, err = l7policies.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return diag.Errorf("unable to update L7 Policy %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Policy(ctx, elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return resourceL7PolicyV3Read(ctx, d, meta)
}

func resourceL7PolicyV3Delete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	elbClient, err := cfg.ElbV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating ELB client: %s", err)
	}

	log.Printf("[DEBUG] Attempting to delete L7 Policy %s", d.Id())
	err = l7policies.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting L7 Policy")
	}

	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Policy(ctx, elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func waitForElbV3Policy(ctx context.Context, elbClient *golangsdk.ServiceClient,
	id string, target string, pending []string, timeout time.Duration) error {
	log.Printf("[DEBUG] Waiting for policy %s to become %s", id, target)

	stateConf := &resource.StateChangeConf{
		Target:       []string{target},
		Pending:      pending,
		Refresh:      resourceElbV3PolicyRefreshFunc(elbClient, id),
		Timeout:      timeout,
		Delay:        5 * time.Second,
		PollInterval: 3 * time.Second,
	}

	_, err := stateConf.WaitForStateContext(ctx)
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmt.Errorf("error: policy %s not found: %s", id, err)
			}
		}
		return fmt.Errorf("error waiting for policy %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3PolicyRefreshFunc(elbClient *golangsdk.ServiceClient, id string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		policy, err := l7policies.Get(elbClient, id).Extract()
		if err != nil {
			return nil, "", err
		}

		return policy, policy.ProvisioningStatus, nil
	}
}
