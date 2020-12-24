package huaweicloud

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"

	"github.com/huaweicloud/golangsdk/openstack/common/tags"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/lbaas_v2/listeners"
)

func ResourceListenerV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceListenerV2Create,
		Read:   resourceListenerV2Read,
		Update: resourceListenerV2Update,
		Delete: resourceListenerV2Delete,

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

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"TCP", "UDP", "HTTP", "HTTPS", "TERMINATED_HTTPS",
				}, false),
			},

			"protocol_port": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},

			"tenant_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"default_pool_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"connection_limit": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"http2_enable": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},

			"default_tls_container_ref": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"sni_container_refs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Default:  true,
				Optional: true,
			},
			"tags": tagsSchema(),
		},
	}
}

func resourceListenerV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	adminStateUp := d.Get("admin_state_up").(bool)
	http2_enable := d.Get("http2_enable").(bool)
	var sniContainerRefs []string
	if raw, ok := d.GetOk("sni_container_refs"); ok {
		for _, v := range raw.([]interface{}) {
			sniContainerRefs = append(sniContainerRefs, v.(string))
		}
	}
	createOpts := listeners.CreateOpts{
		Protocol:               listeners.Protocol(d.Get("protocol").(string)),
		ProtocolPort:           d.Get("protocol_port").(int),
		TenantID:               d.Get("tenant_id").(string),
		LoadbalancerID:         d.Get("loadbalancer_id").(string),
		Name:                   d.Get("name").(string),
		DefaultPoolID:          d.Get("default_pool_id").(string),
		Description:            d.Get("description").(string),
		DefaultTlsContainerRef: d.Get("default_tls_container_ref").(string),
		SniContainerRefs:       sniContainerRefs,
		Http2Enable:            &http2_enable,
		AdminStateUp:           &adminStateUp,
	}

	if v, ok := d.GetOk("connection_limit"); ok {
		connectionLimit := v.(int)
		createOpts.ConnLimit = &connectionLimit
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)

	// Wait for LoadBalancer to become active before continuing
	lbID := createOpts.LoadbalancerID
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Attempting to create listener")
	var listener *listeners.Listener
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		listener, err = listeners.Create(lbClient, createOpts).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error creating listener: %s", err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	d.SetId(listener.ID)

	//set tags
	tagRaw := d.Get("tags").(map[string]interface{})
	if len(tagRaw) > 0 {
		taglist := expandResourceTags(tagRaw)
		if tagErr := tags.Create(lbClient, "listeners", listener.ID, taglist).ExtractErr(); tagErr != nil {
			return fmt.Errorf("Error setting tags of elb listener %s: %s", listener.ID, tagErr)
		}
	}

	return resourceListenerV2Read(d, meta)
}

func resourceListenerV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	listener, err := listeners.Get(lbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "listener")
	}

	log.Printf("[DEBUG] Retrieved listener %s: %#v", d.Id(), listener)

	d.Set("name", listener.Name)
	d.Set("protocol", listener.Protocol)
	d.Set("tenant_id", listener.TenantID)
	d.Set("description", listener.Description)
	d.Set("protocol_port", listener.ProtocolPort)
	d.Set("admin_state_up", listener.AdminStateUp)
	d.Set("default_pool_id", listener.DefaultPoolID)
	d.Set("connection_limit", listener.ConnLimit)
	d.Set("http2_enable", listener.Http2Enable)
	d.Set("sni_container_refs", listener.SniContainerRefs)
	d.Set("default_tls_container_ref", listener.DefaultTlsContainerRef)
	d.Set("region", GetRegion(d, config))

	// fetch tags
	if resourceTags, err := tags.Get(lbClient, "listeners", d.Id()).Extract(); err == nil {
		tagmap := tagsToMap(resourceTags.Tags)
		d.Set("tags", tagmap)
	} else {
		log.Printf("[WARN] fetching tags of elb listener failed: %s", err)
	}

	return nil
}

func resourceListenerV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	if d.HasChanges("name", "description", "admin_state_up", "connection_limit",
		"default_tls_container_ref", "sni_container_refs", "http2_enable") {
		var updateOpts listeners.UpdateOpts
		if d.HasChange("name") {
			updateOpts.Name = d.Get("name").(string)
		}
		if d.HasChange("description") {
			updateOpts.Description = d.Get("description").(string)
		}
		if d.HasChange("connection_limit") {
			connLimit := d.Get("connection_limit").(int)
			updateOpts.ConnLimit = &connLimit
		}
		if d.HasChange("default_tls_container_ref") {
			updateOpts.DefaultTlsContainerRef = d.Get("default_tls_container_ref").(string)
		}
		if d.HasChange("sni_container_refs") {
			var sniContainerRefs []string
			if raw, ok := d.GetOk("sni_container_refs"); ok {
				for _, v := range raw.([]interface{}) {
					sniContainerRefs = append(sniContainerRefs, v.(string))
				}
			}
			updateOpts.SniContainerRefs = sniContainerRefs
		}
		if d.HasChange("admin_state_up") {
			asu := d.Get("admin_state_up").(bool)
			updateOpts.AdminStateUp = &asu
		}
		if d.HasChange("http2_enable") {
			http2 := d.Get("http2_enable").(bool)
			updateOpts.Http2Enable = &http2
		}

		// Wait for LoadBalancer to become active before continuing
		lbID := d.Get("loadbalancer_id").(string)
		timeout := d.Timeout(schema.TimeoutUpdate)
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Updating listener %s with options: %#v", d.Id(), updateOpts)
		//lintignore:R006
		err = resource.Retry(timeout, func() *resource.RetryError {
			_, err = listeners.Update(lbClient, d.Id(), updateOpts).Extract()
			if err != nil {
				return checkForRetryableError(err)
			}
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error updating listener %s: %s", d.Id(), err)
		}

		// Wait for LoadBalancer to become active again before continuing
		err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
		if err != nil {
			return err
		}
	}

	// update tags
	if d.HasChange("tags") {
		tagErr := UpdateResourceTags(lbClient, d, "listeners", d.Id())
		if tagErr != nil {
			return fmt.Errorf("Error updating tags of elb listener:%s, err:%s", d.Id(), tagErr)
		}
	}

	return resourceListenerV2Read(d, meta)

}

func resourceListenerV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	lbClient, err := config.ElbV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud elb client: %s", err)
	}

	// Wait for LoadBalancer to become active before continuing
	lbID := d.Get("loadbalancer_id").(string)
	timeout := d.Timeout(schema.TimeoutDelete)
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Deleting listener %s", d.Id())
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err = listeners.Delete(lbClient, d.Id()).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("Error deleting listener %s: %s", d.Id(), err)
	}

	// Wait for LoadBalancer to become active again before continuing
	err = waitForLBV2LoadBalancer(lbClient, lbID, "ACTIVE", nil, timeout)
	if err != nil {
		return err
	}

	// Wait for Listener to delete
	err = waitForLBV2Listener(lbClient, d.Id(), "DELETED", nil, timeout)
	if err != nil {
		return err
	}

	return nil
}
