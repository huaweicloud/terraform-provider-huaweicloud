package huaweicloud

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"

	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/elb/listeners"
)

const nameELBListener = "ELB-Listener"

func resourceELBListener() *schema.Resource {
	return &schema.Resource{
		Create:             resourceELBListenerCreate,
		Read:               resourceELBListenerRead,
		Update:             resourceELBListenerUpdate,
		Delete:             resourceELBListenerDelete,
		DeprecationMessage: "use ELB(Enhanced) resource instead",

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[a-zA-Z0-9-_]{1,64}$"),
					"Input is a string of 1 to 64 characters that consist of letters, digits, underscores (_), and hyphens (-)"),
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: validation.StringMatch(regexp.MustCompile("^[^<>]{0,128}$"),
					"Input is a string of 0 to 128 characters and cannot contain angle brackets (<>)"),
			},

			"loadbalancer_id": {
				Type:     schema.TypeString,
				Required: true,
			},

			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "HTTP":
					case "TCP":
					case "HTTPS":
					case "SSL":
					case "UDP":
					default:
						errors = append(errors, fmt.Errorf("The valid value of %s is: HTTP, TCP, HTTPS, SSL, UDP", k))
					}
					return
				},
			},

			"port": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 65535 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [1, 65535]", k))
					}
					return
				},
			},

			"backend_protocol": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "HTTP":
					case "TCP":
					case "UDP":
					default:
						errors = append(errors, fmt.Errorf("The valid value of %s is: HTTP, TCP, UDP", k))
					}
					return
				},
			},

			"backend_port": {
				Type:     schema.TypeInt,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 65535 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [1, 65535]", k))
					}
					return
				},
			},

			"lb_algorithm": {
				Type:     schema.TypeString,
				Required: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "roundrobin":
					case "leastconn":
					case "source":
					default:
						errors = append(errors, fmt.Errorf("The valid value of %s is: roundrobin, leastconn, source", k))
					}
					return
				},
			},

			"session_sticky": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"sticky_session_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					if value != "insert" {
						errors = append(errors, fmt.Errorf("The valid value of %s is: insert", k))
					}
					return
				},
			},

			"cookie_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 1440 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [1, 1440]", k))
					}
					return
				},
			},

			"tcp_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 5 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [1, 5]", k))
					}
					return
				},
			},

			"tcp_draining": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},

			"tcp_draining_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 0 || value > 60 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [0, 60]", k))
					}
					return
				},
			},

			"certificate_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"certificates": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},

			"udp_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(int)
					if value < 1 || value > 1440 {
						errors = append(errors, fmt.Errorf("The value of %s must be in [1, 1440]", k))
					}
					return
				},
			},

			"ssl_protocols": {
				Type:     schema.TypeString,
				Optional: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "TLSv1.2":
					case "TLSv1.1":
					case "TLSv1":
					default:
						errors = append(errors, fmt.Errorf("The valid value of %s is: TLSv1.2 TLSv1.1 TLSv1", k))
					}
					return
				},
				Default: "TLSv1.2",
			},

			"ssl_ciphers": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					switch value {
					case "Default":
					case "Extended":
					case "Strict":
					default:
						errors = append(errors, fmt.Errorf("The valid value of %s is: Default, Extended, Strict", k))
					}
					return
				},
			},

			"update_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"create_time": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"admin_state_up": {
				Type:     schema.TypeBool,
				Computed: true,
			},

			"member_number": {
				Type:     schema.TypeInt,
				Computed: true,
			},

			"healthcheck_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceELBListenerCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.elasticLBClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	var opts listeners.CreateOpts
	not_pass_params, err := buildCreateParam(&opts, d, nil)
	if err != nil {
		return fmt.Errorf("Error creating %s: building parameter failed:%s", nameELBListener, err)
	}
	log.Printf("[DEBUG] Create %s Options: %#v", nameELBListener, opts)

	switch {
	case (opts.Protocol == "HTTPS" || opts.Protocol == "SSL") && !hasFilledOpt(d, "certificate_id"):
		return fmt.Errorf("certificate_id is mandatory when protocol is set to HTTPS or SSL")
	}
	l, err := listeners.Create(elbClient, opts, not_pass_params).Extract()
	if err != nil {
		return fmt.Errorf("Error creating %s: %s", nameELBListener, err)
	}
	log.Printf("[DEBUG] Create %s: %#v", nameELBListener, *l)

	// Wait for Listener to become active before continuing
	timeout := d.Timeout(schema.TimeoutCreate)
	err = waitForELBListenerActive(elbClient, l.ID, timeout)
	if err != nil {
		return err
	}

	// If all has been successful, set the ID on the resource
	d.SetId(l.ID)

	return resourceELBListenerRead(d, meta)
}

func resourceELBListenerRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.elasticLBClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	l, err := listeners.Get(elbClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "listener")
	}
	log.Printf("[DEBUG] Retrieved %s %s: %#v", nameELBListener, d.Id(), l)

	sp := d.Get("ssl_protocols")
	if l.SslProtocols == "" && sp != nil {
		l.SslProtocols = sp.(string)
	}
	return refreshResourceData(l, d, nil)
}

func resourceELBListenerUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.elasticLBClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lId := d.Id()

	var opts listeners.UpdateOpts
	not_pass_params, err := buildUpdateParam(&opts, d, nil)
	if err != nil {
		return fmt.Errorf("Error updating %s %s: building parameter failed:%s", nameELBListener, lId, err)
	}

	protocol := d.Get("protocol").(string)
	switch {
	case (protocol == "HTTPS" || protocol == "SSL") && !hasFilledOpt(d, "certificate_id"):
		return fmt.Errorf("certificate_id is mandatory when protocol is set to HTTPS or SSL")
	}
	// Wait for Listener to become active before continuing
	timeout := d.Timeout(schema.TimeoutUpdate)
	err = waitForELBListenerActive(elbClient, lId, timeout)
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating %s %s with options: %#v", nameELBListener, lId, opts)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		_, err := listeners.Update(elbClient, lId, opts, not_pass_params).Extract()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("Error updating %s %s: %s", nameELBListener, lId, err)
	}

	return resourceELBListenerRead(d, meta)
}

func resourceELBListenerDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)
	elbClient, err := config.elasticLBClient(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud networking client: %s", err)
	}

	lId := d.Id()
	log.Printf("[DEBUG] Deleting %s %s", nameELBListener, lId)

	timeout := d.Timeout(schema.TimeoutDelete)
	//lintignore:R006
	err = resource.Retry(timeout, func() *resource.RetryError {
		err := listeners.Delete(elbClient, lId).ExtractErr()
		if err != nil {
			return checkForRetryableError(err)
		}
		return nil
	})
	if err != nil {
		if isResourceNotFound(err) {
			log.Printf("[INFO] deleting an unavailable %s: %s", nameELBListener, lId)
			return nil
		}
		return fmt.Errorf("Error deleting %s %s: %s", nameELBListener, lId, err)
	}

	return nil
}
