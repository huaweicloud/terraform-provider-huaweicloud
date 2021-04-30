package huaweicloud

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/huaweicloud/golangsdk/openstack/fgs/v2/function"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func resourceFgsFunctionV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceFgsFunctionV2Create,
		Read:   resourceFgsFunctionV2Read,
		Delete: resourceFgsFunctionV2Delete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
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
				Required: true,
				ForceNew: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"package": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"app"},
				Deprecated:    "use app instead",
			},
			"app": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"package"},
			},
			"code_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"code_url": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"code_filename": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"handler": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"memory_size": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
				ForceNew: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"xrole": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"agency"},
				Deprecated:    "use agency instead",
			},
			"agency": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"xrole"},
			},
			"func_code": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceFgsFunctionV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	// check app and package
	app, app_ok := d.GetOk("app")
	pak, pak_ok := d.GetOk("package")
	if !app_ok && !pak_ok {
		return fmt.Errorf("One of app or package must be configured")
	}
	pack_v := ""
	if app_ok {
		pack_v = app.(string)
	} else {
		pack_v = pak.(string)
	}

	// get value from agency or xrole
	agency_v := ""
	if v, ok := d.GetOk("agency"); ok {
		agency_v = v.(string)
	} else if v, ok := d.GetOk("xrole"); ok {
		agency_v = v.(string)
	}

	func_code := function.FunctionCodeOpts{
		File: d.Get("func_code").(string),
	}

	createOpts := function.CreateOpts{
		FuncName:     d.Get("name").(string),
		Package:      pack_v,
		CodeType:     d.Get("code_type").(string),
		CodeUrl:      d.Get("code_url").(string),
		Description:  d.Get("description").(string),
		CodeFilename: d.Get("code_filename").(string),
		Handler:      d.Get("handler").(string),
		MemorySize:   d.Get("memory_size").(int),
		Runtime:      d.Get("runtime").(string),
		Timeout:      d.Get("timeout").(int),
		UserData:     d.Get("user_data").(string),
		Xrole:        agency_v,
		FuncCode:     func_code,
	}

	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	f, err := function.Create(fgsClient, createOpts).Extract()
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud function: %s", err)
	}

	d.SetId(f.FuncUrn)

	return resourceFgsFunctionV2Read(d, meta)
}

func resourceFgsFunctionV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	f, err := function.GetMetadata(fgsClient, d.Id()).Extract()
	if err != nil {
		return CheckDeleted(d, err, "function")
	}

	log.Printf("[DEBUG] Retrieved Function %s: %+v", d.Id(), f)

	d.Set("name", f.FuncName)
	d.Set("code_type", f.CodeType)
	d.Set("code_url", f.CodeUrl)
	d.Set("description", f.Description)
	d.Set("code_filename", f.CodeFileName)
	d.Set("handler", f.Handler)
	d.Set("memory_size", f.MemorySize)
	d.Set("runtime", f.Runtime)
	d.Set("timeout", f.Timeout)
	d.Set("user_data", f.UserData)

	if _, ok := d.GetOk("app"); ok {
		d.Set("app", f.Package)
	} else {
		d.Set("package", f.Package)
	}

	if _, ok := d.GetOk("agency"); ok {
		d.Set("agency", f.Xrole)
	} else {
		d.Set("xrole", f.Xrole)
	}

	return nil
}

func resourceFgsFunctionV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(GetRegion(d, config))
	if err != nil {
		return fmt.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	urn := d.Id()
	if strings.HasSuffix(urn, ":latest") {
		urn = urn[0 : len(urn)-7]
	}

	err = function.Delete(fgsClient, urn).ExtractErr()
	if err != nil {
		return fmt.Errorf("Error deleting HuaweiCloud function: %s", err)
	}
	d.SetId("")
	return nil
}
