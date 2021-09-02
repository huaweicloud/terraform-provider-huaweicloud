package huaweicloud

import (
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/elb/v3/pools"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourcePoolV3() *schema.Resource {
	return &schema.Resource{
		Create: resourcePoolV3Create,
		Read:   resourcePoolV3Read,
		Update: resourcePoolV3Update,
		Delete: resourcePoolV3Delete,

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
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP",
				}, false),
			},

			// One of loadbalancer_id or listener_id must be provided
			"loadbalancer_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id"},
			},

			// One of loadbalancer_id or listener_id must be provided
			"listener_id": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				AtLeastOneOf: []string{"loadbalancer_id", "listener_id"},
			},

			"lb_method": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ROUND_ROBIN", "LEAST_CONNECTIONS", "SOURCE_IP", "QUIC_CID",
				}, false),
			},

			"persistence": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ForceNew: true,
							ValidateFunc: validation.StringInSlice([]string{
								"SOURCE_IP", "HTTP_COOKIE", "APP_COOKIE",
							}, false),
						},

						"cookie_name": {
							Type:     schema.TypeString,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
		},
	}
}

func resourcePoolV3Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var persistence pools.SessionPersistence
	if p, ok := d.GetOk("persistence"); ok {
		pV := (p.([]interface{}))[0].(map[string]interface{})

		persistence = pools.SessionPersistence{
			Type: pV["type"].(string),
		}

		if persistence.Type == "APP_COOKIE" {
			if pV["cookie_name"].(string) == "" {
				return fmtp.Errorf(
					"Persistence cookie_name needs to be set if using 'APP_COOKIE' persistence type")
			}
			persistence.CookieName = pV["cookie_name"].(string)
		} else {
			if pV["cookie_name"].(string) != "" {
				return fmtp.Errorf(
					"Persistence cookie_name can only be set if using 'APP_COOKIE' persistence type")
			}
		}
	}

	createOpts := pools.CreateOpts{
		Name:           d.Get("name").(string),
		Description:    d.Get("description").(string),
		Protocol:       d.Get("protocol").(string),
		LoadbalancerID: d.Get("loadbalancer_id").(string),
		ListenerID:     d.Get("listener_id").(string),
		LBMethod:       d.Get("lb_method").(string),
	}

	// Must omit if not set
	if persistence != (pools.SessionPersistence{}) {
		createOpts.Persistence = &persistence
	}

	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	pool, err := pools.Create(elbClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating pool: %s", err)
	}

	d.SetId(pool.ID)

	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForElbV3Pool(elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	return resourcePoolV3Read(d, meta)
}

func resourcePoolV3Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	pool, err := pools.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "pool")
	}

	logp.Printf("[DEBUG] Retrieved pool %s: %#v", d.Id(), pool)

	d.Set("lb_method", pool.LBMethod)
	d.Set("protocol", pool.Protocol)
	d.Set("description", pool.Description)
	d.Set("name", pool.Name)
	d.Set("region", GetRegion(d, config))

	if pool.Persistence.Type != "" {
		var persistence []map[string]interface{} = make([]map[string]interface{}, 1)
		params := make(map[string]interface{})
		params["cookie_name"] = pool.Persistence.CookieName
		params["type"] = pool.Persistence.Type
		persistence[0] = params
		if err = d.Set("persistence", persistence); err != nil {
			return fmtp.Errorf("Load balance persistence set error: %s", err)
		}
	}

	return nil
}

func resourcePoolV3Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	var updateOpts pools.UpdateOpts
	if d.HasChange("lb_method") {
		updateOpts.LBMethod = d.Get("lb_method").(string)
	}
	if d.HasChange("name") {
		updateOpts.Name = d.Get("name").(string)
	}
	if d.HasChange("description") {
		updateOpts.Description = d.Get("description").(string)
	}

	logp.Printf("[DEBUG] Updating pool %s with options: %#v", d.Id(), updateOpts)
	_, err = pools.Update(elbClient, d.Id(), updateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Unable to update pool %s: %s", d.Id(), err)
	}

	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForElbV3Pool(elbClient, d.Id(), "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	return resourcePoolV3Read(d, meta)
}

func resourcePoolV3Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	elbClient, err := config.ElbV3Client(GetRegion(d, config))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	logp.Printf("[DEBUG] Attempting to delete pool %s", d.Id())
	err = pools.Delete(elbClient, d.Id()).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Unable to delete pool %s: %s", d.Id(), err)
	}

	// Wait for Pool to delete
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForElbV3Pool(elbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}

func waitForElbV3Pool(elbClient *golangsdk.ServiceClient, id string, target string, pending []string, timeout time.Duration) error {
	logp.Printf("[DEBUG] Waiting for pool %s to become %s.", id, target)

	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Pending:    pending,
		Refresh:    resourceElbV3PoolRefreshFunc(elbClient, id),
		Timeout:    timeout,
		Delay:      5 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	_, err := stateConf.WaitForState()
	if err != nil {
		if _, ok := err.(golangsdk.ErrDefault404); ok {
			switch target {
			case "DELETED":
				return nil
			default:
				return fmtp.Errorf("Error: pool %s not found: %s", id, err)
			}
		}
		return fmtp.Errorf("Error waiting for pool %s to become %s: %s", id, target, err)
	}

	return nil
}

func resourceElbV3PoolRefreshFunc(elbClient *golangsdk.ServiceClient, poolID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		pool, err := pools.Get(elbClient, poolID).Extract()
		if err != nil {
			return nil, "", err
		}

		// The pool resource has no Status attribute, so a successful Get is the best we can do
		return pool, "ACTIVE", nil
	}
}
