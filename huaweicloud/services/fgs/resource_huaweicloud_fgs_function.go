package fgs

import (
	"strings"
	"time"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/fmtp"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils/logp"
)

func ResourceFgsFunctionV2() *schema.Resource {
	return &schema.Resource{
		Create: resourceFgsFunctionV2Create,
		Read:   resourceFgsFunctionV2Read,
		Update: resourceFgsFunctionV2Update,
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
			"memory_size": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"runtime": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"timeout": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"code_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"handler": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"functiongraph_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice([]string{
					"v1", "v2",
				}, false), // The current default value is v1, which may be adjusted in the future.
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"package": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"app"},
				Deprecated:    "use app instead",
			},
			"app": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"package"},
			},
			"code_url": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"code_filename": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"user_data": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"encrypted_user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"xrole": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"agency"},
				Deprecated:    "use agency instead",
			},
			"agency": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"xrole"},
			},
			"app_agency": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"func_code": {
				Type:      schema.TypeString,
				Optional:  true,
				StateFunc: utils.DecodeHashAndHexEncode,
			},
			"depend_list": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"initializer_handler": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"initializer_timeout": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			"vpc_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"network_id"},
			},
			"network_id": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"vpc_id"},
			},
			"mount_user_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"mount_user_group_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"func_mounts": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"mount_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mount_resource": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mount_share_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"local_mount_path": {
							Type:     schema.TypeString,
							Required: true,
						},
						"status": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"custom_image": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"url": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
				ConflictsWith: []string{
					"code_type", "handler",
				},
			},
			"version": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"urn": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func buildCustomImage(imageConfig []interface{}) *function.CustomImage {
	if len(imageConfig) < 1 {
		return nil
	}

	config := imageConfig[0].(map[string]interface{})
	return &function.CustomImage{
		Enabled: true,
		Image:   config["url"].(string),
	}
}

func buildFgsFunctionV2Parameters(d *schema.ResourceData, config *config.Config) (function.CreateOpts, error) {
	// check app and package
	app, app_ok := d.GetOk("app")
	pkg, pkg_ok := d.GetOk("package")
	if !app_ok && !pkg_ok {
		return function.CreateOpts{}, fmtp.Errorf("One of app or package must be configured")
	}
	pack_v := ""
	if app_ok {
		pack_v = app.(string)
	} else {
		pack_v = pkg.(string)
	}

	// get value from agency or xrole (xrole is deplicated)
	agency_v := ""
	if v, ok := d.GetOk("agency"); ok {
		agency_v = v.(string)
	} else if v, ok := d.GetOk("xrole"); ok {
		agency_v = v.(string)
	}
	result := function.CreateOpts{
		FuncName:            d.Get("name").(string),
		Type:                d.Get("functiongraph_version").(string),
		Package:             pack_v,
		CodeType:            d.Get("code_type").(string),
		CodeUrl:             d.Get("code_url").(string),
		Description:         d.Get("description").(string),
		CodeFilename:        d.Get("code_filename").(string),
		Handler:             d.Get("handler").(string),
		MemorySize:          d.Get("memory_size").(int),
		Runtime:             d.Get("runtime").(string),
		Timeout:             d.Get("timeout").(int),
		UserData:            d.Get("user_data").(string),
		EncryptedUserData:   d.Get("encrypted_user_data").(string),
		Xrole:               agency_v,
		EnterpriseProjectID: config.GetEnterpriseProjectID(d),
		CustomImage:         buildCustomImage(d.Get("custom_image").([]interface{})),
	}
	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FunctionCodeOpts{
			File: utils.TryBase64EncodeToString(v.(string)),
		}
		result.FuncCode = &funcCode
	}
	return result, nil
}

func resourceFgsFunctionV2Create(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	createOpts, err := buildFgsFunctionV2Parameters(d, config)
	if err != nil {
		return err
	}
	logp.Printf("[DEBUG] Create Options: %#v", createOpts)
	f, err := function.Create(fgsClient, createOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud function: %s", err)
	}

	// The "func_urn" is the unique identifier of the function
	// in terraform, we convert to id, not using FuncUrn
	d.SetId(f.FuncUrn)
	urn := resourceFgsFunctionUrn(d.Id())
	//lintignore:R019
	if d.HasChanges("vpc_id", "func_mounts", "app_agency", "initializer_handler", "initializer_timeout") {
		err := resourceFgsFunctionV2MetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return err
		}
	}
	if d.HasChange("depend_list") {
		err := resourceFgsFunctionV2CodeUpdate(fgsClient, urn, d)
		if err != nil {
			return err
		}
	}

	return resourceFgsFunctionV2Read(d, meta)
}

func setFgsFunctionApp(d *schema.ResourceData, app string) error {
	if _, ok := d.GetOk("app"); ok {
		return d.Set("app", app)
	}
	return d.Set("package", app)
}

func setFgsFunctionAgency(d *schema.ResourceData, agency string) error {
	if _, ok := d.GetOk("agency"); ok {
		return d.Set("agency", agency)
	}
	return d.Set("xrole", agency)
}

func setFgsFunctionVpcAccess(d *schema.ResourceData, FuncVpc function.FuncVpc) error {
	mErr := multierror.Append(
		d.Set("vpc_id", FuncVpc.VpcId),
		d.Set("network_id", FuncVpc.SubnetId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}
	return nil
}

func setFuncionMountConfig(d *schema.ResourceData, mountConfig function.MountConfig) error {
	// set mount_config
	if mountConfig.MountUser != (function.MountUser{}) {
		funcMounts := make([]map[string]string, 0, len(mountConfig.FuncMounts))
		for _, v := range mountConfig.FuncMounts {
			funcMount := map[string]string{
				"mount_type":       v.MountType,
				"mount_resource":   v.MountResource,
				"mount_share_path": v.MountSharePath,
				"local_mount_path": v.LocalMountPath,
				"status":           v.Status,
			}
			funcMounts = append(funcMounts, funcMount)
		}
		mErr := multierror.Append(
			d.Set("func_mounts", funcMounts),
			d.Set("mount_user_id", mountConfig.MountUser.UserId),
			d.Set("mount_user_group_id", mountConfig.MountUser.UserGroupId),
		)
		if err := mErr.ErrorOrNil(); err != nil {
			return fmtp.Errorf("Error setting vault fields: %s", err)
		}
	}
	return nil
}

func flattenFgsCustomImage(imageConfig function.CustomImage) []map[string]interface{} {
	if (imageConfig != function.CustomImage{}) {
		return []map[string]interface{}{
			{
				"url": imageConfig.Image,
			},
		}
	}
	return nil
}

func resourceFgsFunctionV2Read(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	f, err := function.GetMetadata(fgsClient, d.Id()).Extract()
	if err != nil {
		return common.CheckDeleted(d, err, "function")
	}

	logp.Printf("[DEBUG] Retrieved Function %s: %+v", d.Id(), f)
	mErr := multierror.Append(
		d.Set("name", f.FuncName),
		d.Set("code_type", f.CodeType),
		d.Set("code_url", f.CodeUrl),
		d.Set("description", f.Description),
		d.Set("code_filename", f.CodeFileName),
		d.Set("handler", f.Handler),
		d.Set("memory_size", f.MemorySize),
		d.Set("runtime", f.Runtime),
		d.Set("timeout", f.Timeout),
		d.Set("user_data", f.UserData),
		d.Set("encrypted_user_data", f.EncryptedUserData),
		d.Set("version", f.Version),
		d.Set("urn", resourceFgsFunctionUrn(d.Id())),
		d.Set("app_agency", f.AppXrole),
		d.Set("depend_list", f.DependList),
		d.Set("initializer_handler", f.InitializerHandler),
		d.Set("initializer_timeout", f.InitializerTimeout),
		d.Set("enterprise_project_id", f.EnterpriseProjectID),
		d.Set("functiongraph_version", f.Type),
		d.Set("custom_image", flattenFgsCustomImage(f.CustomImage)),
		setFgsFunctionApp(d, f.Package),
		setFgsFunctionAgency(d, f.Xrole),
		setFgsFunctionVpcAccess(d, f.FuncVpc),
		setFuncionMountConfig(d, f.MountConfig),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmtp.Errorf("Error setting vault fields: %s", err)
	}

	return nil
}

func resourceFgsFunctionV2Update(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	//lintignore:R019
	if d.HasChanges("code_type", "code_url", "code_filename", "depend_list", "func_code") {
		err := resourceFgsFunctionV2CodeUpdate(fgsClient, urn, d)
		if err != nil {
			return err
		}
	}
	//lintignore:R019
	if d.HasChanges("app", "handler", "depend_list", "memory_size", "timeout", "encrypted_user_data",
		"user_data", "agency", "app_agency", "description", "initializer_handler", "initializer_timeout",
		"vpc_id", "network_id", "mount_user_id", "mount_user_group_id", "func_mounts") {
		err := resourceFgsFunctionV2MetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return err
		}
	}

	return resourceFgsFunctionV2Read(d, meta)
}

func resourceFgsFunctionV2Delete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*config.Config)
	fgsClient, err := config.FgsV2Client(config.GetRegion(d))
	if err != nil {
		return fmtp.Errorf("Error creating HuaweiCloud FGS V2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	err = function.Delete(fgsClient, urn).ExtractErr()
	if err != nil {
		return fmtp.Errorf("Error deleting HuaweiCloud function: %s", err)
	}
	d.SetId("")
	return nil
}

func resourceFgsFunctionV2MetadataUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	// check app and package
	app, app_ok := d.GetOk("app")
	pkg, pkg_ok := d.GetOk("package")
	if !app_ok && !pkg_ok {
		return fmtp.Errorf("One of app or package must be configured")
	}
	pack_v := ""
	if app_ok {
		pack_v = app.(string)
	} else {
		pack_v = pkg.(string)
	}

	// get value from agency or xrole
	agency_v := ""
	if v, ok := d.GetOk("agency"); ok {
		agency_v = v.(string)
	} else if v, ok := d.GetOk("xrole"); ok {
		agency_v = v.(string)
	}

	updateMetadateOpts := function.UpdateMetadataOpts{
		Handler:            d.Get("handler").(string),
		MemorySize:         d.Get("memory_size").(int),
		Timeout:            d.Get("timeout").(int),
		Runtime:            d.Get("runtime").(string),
		Package:            pack_v,
		Description:        d.Get("description").(string),
		UserData:           d.Get("user_data").(string),
		EncryptedUserData:  d.Get("encrypted_user_data").(string),
		Xrole:              agency_v,
		AppXrole:           d.Get("app_agency").(string),
		InitializerHandler: d.Get("initializer_handler").(string),
		InitializerTimeout: d.Get("initializer_timeout").(int),
	}

	if _, ok := d.GetOk("vpc_id"); ok {
		updateMetadateOpts.FuncVpc = resourceFgsFunctionFuncVpc(d)
	}

	if _, ok := d.GetOk("func_mounts"); ok {
		updateMetadateOpts.MountConfig = resourceFgsFunctionMountConfig(d)
	}

	logp.Printf("[DEBUG] Metaddata Update Options: %#v", updateMetadateOpts)
	_, err := function.UpdateMetadata(fgsClient, urn, updateMetadateOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating metadata of HuaweiCloud function: %s", err)
	}

	return nil
}

func resourceFgsFunctionV2CodeUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	updateCodeOpts := function.UpdateCodeOpts{
		CodeType:     d.Get("code_type").(string),
		CodeUrl:      d.Get("code_url").(string),
		CodeFileName: d.Get("code_filename").(string),
	}

	if v, ok := d.GetOk("depend_list"); ok {
		dependListRaw := v.([]interface{})
		dependList := make([]string, 0, len(dependListRaw))
		for _, depend := range dependListRaw {
			dependList = append(dependList, depend.(string))
		}
		updateCodeOpts.DependList = dependList
	}

	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FunctionCodeOpts{
			File: utils.TryBase64EncodeToString(v.(string)),
		}
		updateCodeOpts.FuncCode = funcCode
	}

	logp.Printf("[DEBUG] Code Update Options: %#v", updateCodeOpts)
	_, err := function.UpdateCode(fgsClient, urn, updateCodeOpts).Extract()
	if err != nil {
		return fmtp.Errorf("Error updating code of HuaweiCloud function: %s", err)
	}

	return nil
}

func resourceFgsFunctionFuncVpc(d *schema.ResourceData) *function.FuncVpc {
	var funcVpc function.FuncVpc
	funcVpc.VpcId = d.Get("vpc_id").(string)
	funcVpc.SubnetId = d.Get("network_id").(string)
	return &funcVpc
}

func resourceFgsFunctionMountConfig(d *schema.ResourceData) *function.MountConfig {
	var mountConfig function.MountConfig
	funcMountsRaw := d.Get("func_mounts").([]interface{})
	if len(funcMountsRaw) >= 1 {
		funcMounts := make([]function.FuncMount, 0, len(funcMountsRaw))
		for _, funcMountRaw := range funcMountsRaw {
			var funcMount function.FuncMount
			funcMountMap := funcMountRaw.(map[string]interface{})
			funcMount.MountType = funcMountMap["mount_type"].(string)
			funcMount.MountResource = funcMountMap["mount_resource"].(string)
			funcMount.MountSharePath = funcMountMap["mount_share_path"].(string)
			funcMount.LocalMountPath = funcMountMap["local_mount_path"].(string)

			funcMounts = append(funcMounts, funcMount)
		}

		mountConfig.FuncMounts = funcMounts

		mountUser := function.MountUser{
			UserId:      -1,
			UserGroupId: -1,
		}

		if v, ok := d.GetOk("mount_user_id"); ok {
			mountUser.UserId = v.(int)
		}

		if v, ok := d.GetOk("mount_user_group_id"); ok {
			mountUser.UserGroupId = v.(int)
		}

		mountConfig.MountUser = mountUser
	}
	return &mountConfig
}

/**
 * Parse urn according from fun_urn.
 * If the separator is not ":" then return to the original value.
 */
func resourceFgsFunctionUrn(urn string) string {
	//urn = urn:fss:ru-moscow-1:0910fc31530026f82fd0c018a303517e:function:default:func_2:latest
	index := strings.LastIndex(urn, ":")
	if index != -1 {
		urn = urn[0:index]
	}
	return urn
}
