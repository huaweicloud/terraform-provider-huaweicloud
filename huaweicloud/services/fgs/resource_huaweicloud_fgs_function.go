package fgs

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/aliases"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/function"
	"github.com/chnsz/golangsdk/openstack/fgs/v2/versions"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

// @API FunctionGraph POST /v2/{project_id}/fgs/functions
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/config
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/versions
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/reservedinstances
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/tags/create
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/tags/delete
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/code
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/config
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/config-max-instance
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/versions
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/{function_urn}/aliases
// @API FunctionGraph POST /v2/{project_id}/fgs/functions/{function_urn}/aliases
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}/aliases/{alias_name}
// @API FunctionGraph DELETE /v2/{project_id}/fgs/functions/{function_urn}
// @API FunctionGraph PUT /v2/{project_id}/fgs/functions/{function_urn}/reservedinstances
// @API FunctionGraph GET /v2/{project_id}/fgs/functions/reservedinstanceconfigs
func ResourceFgsFunctionV2() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceFgsFunctionCreate,
		ReadContext:   resourceFgsFunctionRead,
		UpdateContext: resourceFgsFunctionUpdate,
		DeleteContext: resourceFgsFunctionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
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
				Description: utils.SchemaDesc(
					`The code type of the function.`,
					utils.SchemaDescInput{
						Required: true,
					},
				),
			},
			"handler": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: `schema: Required; The entry point of the function.`,
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
				Description:   "schema: Required",
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
				Type:     schema.TypeSet,
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
			"dns_list": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
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
			"log_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				RequiredWith: []string{
					"log_stream_id", "log_group_name", "log_stream_name"},
			},
			"log_stream_id": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
			},
			"log_group_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
			},
			"log_stream_name": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				RequiredWith: []string{"log_group_id"},
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
							Type:        schema.TypeString,
							Required:    true,
							Description: "Specifies the URL of SWR image.",
						},
						"command": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the startup commands of the SWR image.",
						},
						"args": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "Specifies the command line arguments used to start the SWR image.",
						},
						"working_dir": {
							Type:        schema.TypeString,
							Optional:    true,
							Computed:    true,
							Description: "Specifies the working directory of the SWR image.",
						},
						"user_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "schema: Internal; Specifies the user ID for running the image.",
						},
						"user_group_id": {
							Type:        schema.TypeString,
							Optional:    true,
							Description: "schema: Internal; Specifies the user group ID for running the image.",
						},
					},
				},
				ConflictsWith: []string{
					"code_type",
				},
			},
			"max_instance_num": {
				// The original type of this parameter is int, but its zero value is meaningful.
				// So, the following types of parameter passing are realized through the logic of terraform's implicit
				// conversion of int:
				//   + -1: the number of instances is unlimited.
				//   + 0: this function is disabled.
				//   + (0, +1000]: Specific value (2023.06.26).
				//   + empty: keep the default (latest updated) value.
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"versions": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:        schema.TypeString,
							Required:    true,
							Description: "The version name.",
						},
						"aliases": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:        schema.TypeString,
										Required:    true,
										Description: "The name of the version alias.",
									},
									"description": {
										Type:        schema.TypeString,
										Optional:    true,
										Description: "The description of the version alias.",
									},
								},
							},
							Description: "The aliases management for specified version.",
						},
					},
				},
				Description: "The versions management of the function.",
			},
			"tags": common.TagsSchema(),
			"reserved_instances": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"qualifier_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"qualifier_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"idle_mode": {
							Type:     schema.TypeBool,
							Optional: true,
						},
						"tactics_config": {
							Type:     schema.TypeList,
							Optional: true,
							MaxItems: 1,
							Elem:     tracticsConfigsSchema(),
						},
					},
				},
			},
			// The value in the api document is -1 to 1000, After confirmation, when the parameter set to -1 or 0,
			// the actual number of concurrent requests is 1, so the value range is set to 1 to 1000, and the document
			// will be modified later (2024.02.29).
			"concurrency_num": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"gpu_memory": {
				Type:         schema.TypeInt,
				Optional:     true,
				RequiredWith: []string{"gpu_type"},
			},
			"gpu_type": {
				Type:         schema.TypeString,
				Optional:     true,
				RequiredWith: []string{"gpu_memory"},
			},
			// Currently, the "pre_stop_timeout" and "pre_stop_timeout" are not visible on the page,
			// so they are temporarily used as internal parameters.
			"pre_stop_handler": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "schema: Internal; Specifies the pre-stop handler of a function.",
			},
			"pre_stop_timeout": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: "schema: Internal; Specifies the maximum duration that the function can be initialized.",
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

func tracticsConfigsSchema() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"cron_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cron": {
							Type:     schema.TypeString,
							Required: true,
						},
						"count": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"start_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"expired_time": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
			"metric_configs": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"threshold": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"min": {
							Type:     schema.TypeInt,
							Required: true,
						},
					},
				},
			},
		},
	}
}

func buildCustomImage(imageConfig []interface{}) *function.CustomImage {
	if len(imageConfig) < 1 {
		return nil
	}

	cfg := imageConfig[0].(map[string]interface{})
	return &function.CustomImage{
		Enabled:     true,
		Image:       cfg["url"].(string),
		Command:     cfg["command"].(string),
		Args:        cfg["args"].(string),
		WorkingDir:  cfg["working_dir"].(string),
		UserId:      cfg["user_id"].(string),
		UserGroupId: cfg["user_group_id"].(string),
	}
}

func buildFgsFunctionParameters(d *schema.ResourceData, cfg *config.Config) (function.CreateOpts, error) {
	// check app and package
	app, appOk := d.GetOk("app")
	pkg, pkgOk := d.GetOk("package")
	if !appOk && !pkgOk {
		return function.CreateOpts{}, fmt.Errorf("one of app or package must be configured")
	}
	packV := ""
	if appOk {
		packV = app.(string)
	} else {
		packV = pkg.(string)
	}

	// get value from agency or xrole (xrole is deplicated)
	agencyV := ""
	if v, ok := d.GetOk("agency"); ok {
		agencyV = v.(string)
	} else if v, ok := d.GetOk("xrole"); ok {
		agencyV = v.(string)
	}
	result := function.CreateOpts{
		FuncName:            d.Get("name").(string),
		Type:                d.Get("functiongraph_version").(string),
		Package:             packV,
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
		Xrole:               agencyV,
		EnterpriseProjectID: cfg.GetEnterpriseProjectID(d),
		CustomImage:         buildCustomImage(d.Get("custom_image").([]interface{})),
		GPUMemory:           d.Get("gpu_memory").(int),
		GPUType:             d.Get("gpu_type").(string),
		PreStopHandler:      d.Get("pre_stop_handler").(string),
		PreStopTimeout:      d.Get("pre_stop_timeout").(int),
	}
	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FunctionCodeOpts{
			File: utils.TryBase64EncodeString(v.(string)),
		}
		result.FuncCode = &funcCode
	}
	if v, ok := d.GetOk("log_group_id"); ok {
		logConfig := function.FuncLogConfig{
			GroupId:    v.(string),
			StreamId:   d.Get("log_stream_id").(string),
			GroupName:  d.Get("log_group_name").(string),
			StreamName: d.Get("log_stream_name").(string),
		}
		result.LogConfig = &logConfig
	}
	return result, nil
}

func resourceFgsFunctionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	fgsClient, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	createOpts, err := buildFgsFunctionParameters(d, cfg)
	if err != nil {
		return diag.FromErr(err)
	}
	log.Printf("[DEBUG] Create Options: %#v", createOpts)
	f, err := function.Create(fgsClient, createOpts).Extract()
	if err != nil {
		return diag.Errorf("error creating function: %s", err)
	}

	// The "func_urn" is the unique identifier of the function
	// in terraform, we convert to id, not using FuncUrn
	d.SetId(f.FuncUrn)
	urn := resourceFgsFunctionUrn(d.Id())
	// lintignore:R019
	if d.HasChanges("vpc_id", "func_mounts", "app_agency", "initializer_handler", "initializer_timeout", "concurrency_num") {
		err := resourceFgsFunctionMetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("depend_list") {
		err := resourceFgsFunctionCodeUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if strNum, ok := d.GetOk("max_instance_num"); ok {
		// The integer string of the maximum instance number has been already checked in the schema validation.
		maxInstanceNum, _ := strconv.Atoi(strNum.(string))
		_, err = function.UpdateMaxInstanceNumber(fgsClient, urn, maxInstanceNum)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if tagList, ok := d.GetOk("tags"); ok {
		opts := function.TagsActionOpts{
			Tags: utils.ExpandResourceTags(tagList.(map[string]interface{})),
		}
		if err := function.CreateResourceTags(fgsClient, d.Id(), opts); err != nil {
			return diag.Errorf("failed to add tags to FunctionGraph function (%s): %s", d.Id(), err)
		}
	}

	if err = createFunctionVersions(fgsClient, urn, d.Get("versions").(*schema.Set)); err != nil {
		return diag.Errorf("error creating function versions: %s", err)
	}

	if d.HasChanges("reserved_instances") {
		if err = updateReservedInstanceConfig(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFgsFunctionRead(ctx, d, meta)
}

func createFunctionVersions(client *golangsdk.ServiceClient, functionUrn string, versionSet *schema.Set) error {
	for _, v := range versionSet.List() {
		version := v.(map[string]interface{})
		versionNum := version["name"].(string) // The version name, also name as the version number.

		if versionNum != "latest" {
			createOpts := versions.CreateOpts{
				FunctionUrn: functionUrn,
				Version:     versionNum,
			}
			_, err := versions.Create(client, createOpts)
			if err != nil {
				return err
			}
		}
		// In the future, the function will support manage multiple versions, and will add the corresponding logic to
		// create versions based on the related API (Create) in this place.
		aliasCfg := version["aliases"].([]interface{})
		for _, val := range aliasCfg {
			alias := val.(map[string]interface{})
			opt := aliases.CreateOpts{
				FunctionUrn: functionUrn,
				Name:        alias["name"].(string),
				Version:     versionNum,
				Description: alias["description"].(string),
			}
			_, err := aliases.Create(client, opt)
			if err != nil {
				return err
			}
		}
	}
	return nil
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

func setFgsFunctionVpcAccess(d *schema.ResourceData, funcVpc function.FuncVpc) error {
	mErr := multierror.Append(
		d.Set("vpc_id", funcVpc.VpcId),
		d.Set("network_id", funcVpc.SubnetId),
	)
	if err := mErr.ErrorOrNil(); err != nil {
		return fmt.Errorf("error setting vault fields: %s", err)
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
			return fmt.Errorf("error setting vault fields: %s", err)
		}
	}
	return nil
}

func flattenFgsCustomImage(imageConfig function.CustomImage) []map[string]interface{} {
	if (imageConfig != function.CustomImage{}) {
		return []map[string]interface{}{
			{
				"url":           imageConfig.Image,
				"command":       imageConfig.Command,
				"args":          imageConfig.Args,
				"working_dir":   imageConfig.WorkingDir,
				"user_id":       imageConfig.UserId,
				"user_group_id": imageConfig.UserGroupId,
			},
		}
	}
	return nil
}

func queryFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]string, error) {
	queryOpts := versions.ListOpts{
		FunctionUrn: functionUrn,
	}
	versionList, err := versions.List(client, queryOpts)
	if err != nil {
		return nil, fmt.Errorf("error querying version list for the specified function URN: %s", err)
	}
	// The length of the function version list is at least 1 (when creating a function, a version named latest is
	// created by default).
	result := make([]string, len(versionList))
	for i, version := range versionList {
		result[i] = version.Version
	}
	return result, nil
}

func queryFunctionAliases(client *golangsdk.ServiceClient, functionUrn string) (map[string][]interface{}, error) {
	aliasList, err := aliases.List(client, functionUrn)
	if err != nil {
		return nil, fmt.Errorf("error querying alias list for the specified function URN: %s", err)
	}

	// Multiple version aliases may exist in the future.
	result := make(map[string][]interface{})
	for _, v := range aliasList {
		result[v.Version] = append(result[v.Version], map[string]interface{}{
			"name":        v.Name,
			"description": v.Description,
		})
	}
	return result, nil
}

func parseFunctionVersions(client *golangsdk.ServiceClient, functionUrn string) ([]map[string]interface{}, error) {
	versionList, err := queryFunctionVersions(client, functionUrn)
	if err != nil {
		return nil, err
	}
	aliasesConfig, err := queryFunctionAliases(client, functionUrn)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, 0, len(versionList))
	for _, versionNum := range versionList {
		version := map[string]interface{}{
			"name": versionNum, // The version name, also name as the version number.
		}
		if v, ok := aliasesConfig[versionNum]; ok {
			version["aliases"] = v
		} else if versionNum == "latest" {
			// If no alias is set for the default version, the corresponding structure object is not saved.
			continue
		}
		result = append(result, version)
	}

	return result, nil
}

func flattenTracticsConfigs(policyConfig function.TacticsConfigObj) []map[string]interface{} {
	if len(policyConfig.CronConfigs) == 0 && len(policyConfig.MetricConfigs) == 0 {
		return nil
	}

	cronConfigRst := make([]map[string]interface{}, len(policyConfig.CronConfigs))
	for i, v := range policyConfig.CronConfigs {
		cronConfigRst[i] = map[string]interface{}{
			"name":         v.Name,
			"cron":         v.Cron,
			"count":        v.Count,
			"start_time":   v.StartTime,
			"expired_time": v.ExpiredTime,
		}
	}

	metricConfigs := make([]map[string]interface{}, len(policyConfig.MetricConfigs))
	for i, v := range policyConfig.MetricConfigs {
		metricConfigs[i] = map[string]interface{}{
			"name":      v.Name,
			"type":      v.Type,
			"threshold": v.Threshold,
			"min":       v.Min,
		}
	}

	return []map[string]interface{}{
		{
			"cron_configs":   cronConfigRst,
			"metric_configs": metricConfigs,
		},
	}
}

func getReservedInstanceConfig(c *golangsdk.ServiceClient, d *schema.ResourceData) ([]map[string]interface{}, error) {
	opts := function.ListReservedInstanceConfigOpts{
		FunctionUrn: d.Id(),
	}
	reservedInstances, err := function.ListReservedInstanceConfigs(c, opts)
	if err != nil {
		return nil, fmt.Errorf("error getting list of the function reserved instance config: %s", err)
	}

	result := make([]map[string]interface{}, len(reservedInstances))
	for i, v := range reservedInstances {
		result[i] = map[string]interface{}{
			"count":          v.MinCount,
			"idle_mode":      v.IdleMode,
			"qualifier_name": v.QualifierName,
			"qualifier_type": v.QualifierType,
			"tactics_config": flattenTracticsConfigs(v.TacticsConfig),
		}
	}
	return result, nil
}

func getConcurrencyNum(concurrencyNum *int) int {
	return *concurrencyNum
}

func resourceFgsFunctionRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	fgsClient, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph client: %s", err)
	}

	functionUrn := resourceFgsFunctionUrn(d.Id())
	f, err := function.GetMetadata(fgsClient, functionUrn).Extract()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "FunctionGraph function")
	}

	versionConfig, err := parseFunctionVersions(fgsClient, functionUrn)
	if err != nil {
		// Not all regions support the version related API calls.
		log.Printf("[ERROR] Unable to parsing the function versions: %s", err)
	}
	log.Printf("[DEBUG] Retrieved Function %s: %+v", functionUrn, f)
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
		d.Set("urn", functionUrn),
		d.Set("app_agency", f.AppXrole),
		d.Set("depend_list", f.DependVersionList),
		d.Set("initializer_handler", f.InitializerHandler),
		d.Set("initializer_timeout", f.InitializerTimeout),
		d.Set("enterprise_project_id", f.EnterpriseProjectID),
		d.Set("functiongraph_version", f.Type),
		d.Set("custom_image", flattenFgsCustomImage(f.CustomImage)),
		d.Set("max_instance_num", strconv.Itoa(*f.StrategyConfig.Concurrency)),
		d.Set("dns_list", f.DomainNames),
		d.Set("log_group_id", f.LogGroupId),
		d.Set("log_stream_id", f.LogStreamId),
		setFgsFunctionApp(d, f.Package),
		setFgsFunctionAgency(d, f.Xrole),
		setFgsFunctionVpcAccess(d, f.FuncVpc),
		setFuncionMountConfig(d, f.MountConfig),
		d.Set("concurrency_num", getConcurrencyNum(f.StrategyConfig.ConcurrencyNum)),
		d.Set("versions", versionConfig),
		d.Set("gpu_memory", f.GPUMemory),
		d.Set("gpu_type", f.GPUType),
		d.Set("pre_stop_handler", f.PreStopHandler),
		d.Set("pre_stop_timeout", f.PreStopTimeout),
	)

	reservedInstances, err := getReservedInstanceConfig(fgsClient, d)
	if err != nil {
		return diag.Errorf("error retrieving function reserved instance: %s", err)
	}

	mErr = multierror.Append(mErr,
		d.Set("reserved_instances", reservedInstances),
	)

	if err := mErr.ErrorOrNil(); err != nil {
		return diag.Errorf("error setting function fields: %s", err)
	}

	return nil
}

func updateFunctionTags(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		oRaw, nRaw  = d.GetChange("tags")
		oMap        = oRaw.(map[string]interface{})
		nMap        = nRaw.(map[string]interface{})
		functionUrn = d.Id()
	)

	if len(oMap) > 0 {
		opts := function.TagsActionOpts{
			Tags: utils.ExpandResourceTags(oMap),
		}
		if err := function.DeleteResourceTags(client, functionUrn, opts); err != nil {
			return fmt.Errorf("failed to delete tags from FunctionGraph function (%s): %s", functionUrn, err)
		}
	}

	if len(nMap) > 0 {
		opts := function.TagsActionOpts{
			Tags: utils.ExpandResourceTags(nMap),
		}
		if err := function.CreateResourceTags(client, functionUrn, opts); err != nil {
			return fmt.Errorf("failed to add tags to FunctionGraph function (%s): %s", functionUrn, err)
		}
	}
	return nil
}

func deleteFunctionVersions(client *golangsdk.ServiceClient, functionUrn string, versionSet *schema.Set) error {
	// In the future, the function will support manage multiple versions.
	for _, v := range versionSet.List() {
		version := v.(map[string]interface{})
		versionNum := version["name"].(string) // The version name, also name as the version number.
		if versionNum != "latest" {
			// Deletes a function version, also deleting all aliases beneath it.
			err := function.Delete(client, fmt.Sprintf("%s:%s", functionUrn, versionNum)).ExtractErr()
			if err != nil {
				return err
			}
			continue
		}
		// Since the latest version cannot be deleted, only the version alias under it can be deleted.
		aliasCfg := version["aliases"].([]interface{})
		for _, val := range aliasCfg {
			alias := val.(map[string]interface{})
			err := aliases.Delete(client, functionUrn, alias["name"].(string))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func updateFunctionVersions(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	var (
		functionUrn = resourceFgsFunctionUrn(d.Id())

		oldSet, newSet = d.GetChange("versions")
		decrease       = oldSet.(*schema.Set).Difference(newSet.(*schema.Set))
		increase       = newSet.(*schema.Set).Difference(oldSet.(*schema.Set))
	)

	err := deleteFunctionVersions(client, functionUrn, decrease)
	if err != nil {
		return fmt.Errorf("error deleting function versions: %s", err)
	}

	err = createFunctionVersions(client, functionUrn, increase)
	if err != nil {
		return fmt.Errorf("error creating function versions: %s", err)
	}

	return nil
}

func buildCronConfigs(cronConfigs []interface{}) []function.CronConfigObj {
	if len(cronConfigs) < 1 {
		return nil
	}

	result := make([]function.CronConfigObj, len(cronConfigs))
	for i, v := range cronConfigs {
		cronConfig := v.(map[string]interface{})
		result[i] = function.CronConfigObj{
			Name:        cronConfig["name"].(string),
			Cron:        cronConfig["cron"].(string),
			Count:       cronConfig["count"].(int),
			StartTime:   cronConfig["start_time"].(int),
			ExpiredTime: cronConfig["expired_time"].(int),
		}
	}
	return result
}

func buildMetricConfigs(metricConfigs []interface{}) []function.MetricConfigObj {
	if len(metricConfigs) < 1 {
		return nil
	}
	result := make([]function.MetricConfigObj, len(metricConfigs))
	for i, v := range metricConfigs {
		metricConfig := v.(map[string]interface{})
		result[i] = function.MetricConfigObj{
			Name:      metricConfig["name"].(string),
			Type:      metricConfig["type"].(string),
			Threshold: metricConfig["threshold"].(int),
			Min:       metricConfig["min"].(int),
		}
	}
	return result
}

func buildTracticsConfigs(tacticsConfigs []interface{}) *function.TacticsConfigObj {
	if len(tacticsConfigs) < 1 {
		return nil
	}

	tacticsConfig := tacticsConfigs[0].(map[string]interface{})
	result := function.TacticsConfigObj{
		CronConfigs:   buildCronConfigs(tacticsConfig["cron_configs"].([]interface{})),
		MetricConfigs: buildMetricConfigs(tacticsConfig["metric_configs"].([]interface{})),
	}
	return &result
}

func getVersionUrn(client *golangsdk.ServiceClient, functionUrn string, qualifierName string) (string, error) {
	queryOpts := versions.ListOpts{
		FunctionUrn: functionUrn,
	}
	versionList, err := versions.List(client, queryOpts)
	if err != nil {
		return "", fmt.Errorf("error querying version list for the specified function URN: %s", err)
	}

	for _, val := range versionList {
		if val.Version == qualifierName {
			return val.FuncUrn, nil
		}
	}

	return "", nil
}

func getReservedInstanceUrn(client *golangsdk.ServiceClient, functionUrn string, policy map[string]interface{}) (string, error) {
	qualifierName := policy["qualifier_name"].(string)
	if policy["qualifier_type"].(string) == "version" {
		urn, err := getVersionUrn(client, functionUrn, qualifierName)
		if err != nil {
			return "", err
		}
		return urn, nil
	}

	aliasList, err := aliases.List(client, functionUrn)
	if err != nil {
		return "", fmt.Errorf("error querying alias list for the specified function URN: %s", err)
	}
	for _, val := range aliasList {
		if val.Name == qualifierName {
			return val.AliasUrn, nil
		}
	}

	return "", nil
}

func removeReservedInstances(client *golangsdk.ServiceClient, functionUrn string, policies []interface{}) error {
	for _, v := range policies {
		policy := v.(map[string]interface{})
		urn, err := getReservedInstanceUrn(client, functionUrn, policy)
		if err != nil {
			return err
		}
		// Deleting the alias will also delete the corresponding reserved instance.
		if urn == "" {
			return nil
		}
		opts := function.UpdateReservedInstanceObj{
			FunctionUrn: urn,
			Count:       utils.Int(0),
			IdleMode:    utils.Bool(false),
		}
		_, err = function.UpdateReservedInstanceConfig(client, opts)
		if err != nil {
			return fmt.Errorf("error removing function reversed instance: %s", err)
		}
	}

	return nil
}

func addReservedInstances(client *golangsdk.ServiceClient, functionUrn string, addPolicies []interface{}) error {
	for _, v := range addPolicies {
		addPolicy := v.(map[string]interface{})
		urn, err := getReservedInstanceUrn(client, functionUrn, addPolicy)

		if err != nil {
			return err
		}

		opts := function.UpdateReservedInstanceObj{
			FunctionUrn:   urn,
			Count:         utils.Int(addPolicy["count"].(int)),
			IdleMode:      utils.Bool(addPolicy["idle_mode"].(bool)),
			TacticsConfig: buildTracticsConfigs(addPolicy["tactics_config"].([]interface{})),
		}
		_, err = function.UpdateReservedInstanceConfig(client, opts)
		if err != nil {
			return fmt.Errorf("error updating function reversed instance: %s", err)
		}
	}

	return nil
}

func updateReservedInstanceConfig(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	oldRaw, newRaw := d.GetChange("reserved_instances")
	addRaw := newRaw.(*schema.Set).Difference(oldRaw.(*schema.Set))
	removeRaw := oldRaw.(*schema.Set).Difference(newRaw.(*schema.Set))
	functionUrn := resourceFgsFunctionUrn(d.Id())
	if removeRaw.Len() > 0 {
		if err := removeReservedInstances(client, functionUrn, removeRaw.List()); err != nil {
			return err
		}
	}

	if addRaw.Len() > 0 {
		if err := addReservedInstances(client, functionUrn, addRaw.List()); err != nil {
			return err
		}
	}

	return nil
}

func resourceFgsFunctionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	fgsClient, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	// lintignore:R019
	if d.HasChanges("code_type", "code_url", "code_filename", "depend_list", "func_code") {
		err := resourceFgsFunctionCodeUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	// lintignore:R019
	if d.HasChanges("app", "handler", "memory_size", "timeout", "encrypted_user_data",
		"user_data", "agency", "app_agency", "description", "initializer_handler", "initializer_timeout",
		"vpc_id", "network_id", "dns_list", "mount_user_id", "mount_user_group_id", "func_mounts", "custom_image",
		"log_group_id", "log_stream_id", "log_group_name", "log_stream_name", "concurrency_num", "gpu_memory", "gpu_type") {
		err := resourceFgsFunctionMetadataUpdate(fgsClient, urn, d)
		if err != nil {
			return diag.FromErr(err)
		}
	}
	if d.HasChange("max_instance_num") {
		// The integer string of the maximum instance number has been already checked in the schema validation.
		maxInstanceNum, _ := strconv.Atoi(d.Get("max_instance_num").(string))
		_, err = function.UpdateMaxInstanceNumber(fgsClient, urn, maxInstanceNum)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("tags") {
		if err = updateFunctionTags(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChange("versions") {
		if err = updateFunctionVersions(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	if d.HasChanges("reserved_instances") {
		if err = updateReservedInstanceConfig(fgsClient, d); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceFgsFunctionRead(ctx, d, meta)
}

func resourceFgsFunctionDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	fgsClient, err := cfg.FgsV2Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating FunctionGraph v2 client: %s", err)
	}

	urn := resourceFgsFunctionUrn(d.Id())

	err = function.Delete(fgsClient, urn).ExtractErr()
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting function")
	}
	return nil
}

func resourceFgsFunctionMetadataUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	// check app and package
	app, appOk := d.GetOk("app")
	pkg, pkgOk := d.GetOk("package")
	if !appOk && !pkgOk {
		return fmt.Errorf("one of app or package must be configured")
	}
	packV := ""
	if appOk {
		packV = app.(string)
	} else {
		packV = pkg.(string)
	}

	// get value from agency or xrole
	agencyV := ""
	if v, ok := d.GetOk("agency"); ok {
		agencyV = v.(string)
	} else if v, ok := d.GetOk("xrole"); ok {
		agencyV = v.(string)
	}

	updateMetadateOpts := function.UpdateMetadataOpts{
		Handler:            d.Get("handler").(string),
		MemorySize:         d.Get("memory_size").(int),
		Timeout:            d.Get("timeout").(int),
		Runtime:            d.Get("runtime").(string),
		Package:            packV,
		Description:        d.Get("description").(string),
		UserData:           d.Get("user_data").(string),
		EncryptedUserData:  d.Get("encrypted_user_data").(string),
		Xrole:              agencyV,
		AppXrole:           d.Get("app_agency").(string),
		InitializerHandler: d.Get("initializer_handler").(string),
		InitializerTimeout: d.Get("initializer_timeout").(int),
		CustomImage:        buildCustomImage(d.Get("custom_image").([]interface{})),
		DomainNames:        d.Get("dns_list").(string),
		GPUMemory:          d.Get("gpu_memory").(int),
		GPUType:            d.Get("gpu_type").(string),
		PreStopHandler:     d.Get("pre_stop_handler").(string),
		PreStopTimeout:     d.Get("pre_stop_timeout").(int),
	}

	if _, ok := d.GetOk("vpc_id"); ok {
		updateMetadateOpts.FuncVpc = resourceFgsFunctionFuncVpc(d)
	}

	if _, ok := d.GetOk("func_mounts"); ok {
		updateMetadateOpts.MountConfig = resourceFgsFunctionMountConfig(d)
	}

	// check name here as it will only save to sate if specified before
	if v, ok := d.GetOk("log_group_name"); ok {
		logConfig := function.FuncLogConfig{
			GroupId:    d.Get("log_group_id").(string),
			StreamId:   d.Get("log_stream_id").(string),
			GroupName:  v.(string),
			StreamName: d.Get("log_stream_name").(string),
		}
		updateMetadateOpts.LogConfig = &logConfig
	}

	if v, ok := d.GetOk("concurrency_num"); ok {
		strategyConfig := function.StrategyConfig{
			ConcurrencyNum: utils.Int(v.(int)),
		}
		updateMetadateOpts.StrategyConfig = &strategyConfig
	}

	log.Printf("[DEBUG] Metaddata Update Options: %#v", updateMetadateOpts)
	_, err := function.UpdateMetadata(fgsClient, urn, updateMetadateOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating metadata of function: %s", err)
	}

	return nil
}

func resourceFgsFunctionCodeUpdate(fgsClient *golangsdk.ServiceClient, urn string, d *schema.ResourceData) error {
	updateCodeOpts := function.UpdateCodeOpts{
		CodeType:     d.Get("code_type").(string),
		CodeUrl:      d.Get("code_url").(string),
		CodeFileName: d.Get("code_filename").(string),
	}

	if v, ok := d.GetOk("depend_list"); ok {
		dependListRaw := v.(*schema.Set)
		dependList := make([]string, 0, dependListRaw.Len())
		for _, depend := range dependListRaw.List() {
			dependList = append(dependList, depend.(string))
		}
		updateCodeOpts.DependList = dependList
	}

	if v, ok := d.GetOk("func_code"); ok {
		funcCode := function.FunctionCodeOpts{
			File: utils.TryBase64EncodeString(v.(string)),
		}
		updateCodeOpts.FuncCode = funcCode
	}

	log.Printf("[DEBUG] Code Update Options: %#v", updateCodeOpts)
	_, err := function.UpdateCode(fgsClient, urn, updateCodeOpts).Extract()
	if err != nil {
		return fmt.Errorf("error updating code of function: %s", err)
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

/*
 * Parse urn according from fun_urn.
 * If the separator is not ":" then return to the original value.
 */
func resourceFgsFunctionUrn(urn string) string {
	// urn = urn:fss:ru-moscow-1:0910fc31530026f82fd0c018a303517e:function:default:func_2:latest
	index := strings.LastIndex(urn, ":")
	if index != -1 {
		urn = urn[0:index]
	}
	return urn
}
