package sms

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/chnsz/golangsdk"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/common"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/utils"
)

var sourceServerNonUpdatableParams = []string{"ip", "hostname", "os_type", "os_version", "virtualization_type",
	"linux_block_check", "firmware", "cpu_quantity", "memory", "networks.*.name", "networks.*.ip", "networks.*.ipv6",
	"networks.*.netmask", "networks.*.gateway", "networks.*.mtu", "networks.*.mac", "domain_id", "has_rsync",
	"paravirtualization", "raw_devices", "driver_files", "system_services", "account_rights", "boot_loader", "system_dir",
	"agent_version", "kernel_version", "oem_system", "start_type", "io_read_wait", "has_tc", "platform"}

// @API SMS POST /v3/sources
// @API SMS GET /v3/sources/{source_id}
// @API SMS PUT /v3/sources/{source_id}
// @API SMS DELETE /v3/sources/{source_id}
// @API SMS PUT /v3/sources/{source_id}/diskinfo
// @API SMS PUT /v3/sources/{source_id}/changestate
func ResourceSourceServer() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSourceServerCreate,
		ReadContext:   resourceSourceServerRead,
		UpdateContext: resourceSourceServerUpdate,
		DeleteContext: resourceSourceServerDelete,

		CustomizeDiff: config.FlexibleForceNew(sourceServerNonUpdatableParams),

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"hostname": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"os_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"virtualization_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"linux_block_check": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"firmware": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"cpu_quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},
			"disks": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"device_use": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"used_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"physical_volumes": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"device_use": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"file_system": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"index": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"mount_point": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"used_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"inode_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"inode_nums": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"uuid": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"size_per_cluster": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"adjust_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"need_migration": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"partition_style": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"os_disk": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"relation_name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"inode_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"adjust_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"need_migration": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"btrfs_list": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"label": {
							Type:     schema.TypeString,
							Required: true,
						},
						"uuid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"device": {
							Type:     schema.TypeString,
							Required: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"nodesize": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"sectorsize": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"data_profile": {
							Type:     schema.TypeString,
							Required: true,
						},
						"system_profile": {
							Type:     schema.TypeString,
							Required: true,
						},
						"metadata_profile": {
							Type:     schema.TypeString,
							Required: true,
						},
						"global_reserve1": {
							Type:     schema.TypeString,
							Required: true,
						},
						"g_vol_used_size": {
							Type:     schema.TypeInt,
							Required: true,
						},
						"default_subvolid": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default_subvol_name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"default_subvol_mountpath": {
							Type:     schema.TypeString,
							Required: true,
						},
						"subvolumn": {
							Type:     schema.TypeList,
							Required: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"uuid": {
										Type:     schema.TypeString,
										Required: true,
									},
									"is_snapshot": {
										Type:     schema.TypeString,
										Required: true,
									},
									"subvol_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"parent_id": {
										Type:     schema.TypeString,
										Required: true,
									},
									"subvol_name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"subvol_mount_path": {
										Type:     schema.TypeString,
										Required: true,
									},
								},
							},
						},
					},
				},
			},
			"networks": {
				Type:     schema.TypeList,
				Optional: true,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ip": {
							Type:     schema.TypeString,
							Required: true,
						},
						"netmask": {
							Type:     schema.TypeString,
							Required: true,
						},
						"gateway": {
							Type:     schema.TypeString,
							Required: true,
						},
						"mac": {
							Type:     schema.TypeString,
							Required: true,
						},
						"ipv6": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"mtu": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"domain_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"has_rsync": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"paravirtualization": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"raw_devices": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"driver_files": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"system_services": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"account_rights": {
				Type:     schema.TypeBool,
				Optional: true,
			},
			"boot_loader": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"system_dir": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"volume_groups": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"components": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"free_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"logical_volumes": {
							Type:     schema.TypeList,
							Optional: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"file_system": {
										Type:     schema.TypeString,
										Required: true,
									},
									"inode_size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"mount_point": {
										Type:     schema.TypeString,
										Required: true,
									},
									"name": {
										Type:     schema.TypeString,
										Required: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"used_size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"free_size": {
										Type:     schema.TypeInt,
										Required: true,
									},
									"block_count": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"block_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"inode_nums": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"device_use": {
										Type:     schema.TypeString,
										Optional: true,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"adjust_size": {
										Type:     schema.TypeInt,
										Optional: true,
										Computed: true,
									},
									"need_migration": {
										Type:     schema.TypeBool,
										Optional: true,
										Computed: true,
									},
								},
							},
						},
						"name": {
							Type:     schema.TypeString,
							Optional: true,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"id": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"adjust_size": {
							Type:     schema.TypeInt,
							Optional: true,
							Computed: true,
						},
						"need_migration": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
					},
				},
			},
			"agent_version": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"kernel_version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"migration_cycle": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"state": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"oem_system": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"start_type": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"io_read_wait": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"has_tc": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"platform": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"migprojectid": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"copystate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"enable_force_new": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validation.StringInSlice([]string{"true", "false"}, false),
				Description:  utils.SchemaDesc("", utils.SchemaDescInput{Internal: true}),
			},
			"enterprise_project_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"add_date": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"connected": {
				Type:     schema.TypeBool,
				Computed: true,
			},
			"init_target_server": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"disks": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"device_use": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"used_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"adjust_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"need_migration": {
										Type:     schema.TypeBool,
										Computed: true,
									},
									"physical_volumes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"device_use": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"file_system": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"index": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"mount_point": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"inode_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"used_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"uuid": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"adjust_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"need_migration": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
								},
							},
						},
						"volume_groups": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"components": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"free_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"logical_volumes": {
										Type:     schema.TypeList,
										Computed: true,
										Elem: &schema.Resource{
											Schema: map[string]*schema.Schema{
												"block_count": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"block_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"file_system": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"inode_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"inode_nums": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"device_use": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"mount_point": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"name": {
													Type:     schema.TypeString,
													Computed: true,
												},
												"size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"used_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"free_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"id": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"adjust_size": {
													Type:     schema.TypeInt,
													Computed: true,
												},
												"need_migration": {
													Type:     schema.TypeBool,
													Computed: true,
												},
											},
										},
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"id": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"adjust_size": {
										Type:     schema.TypeInt,
										Computed: true,
									},
									"need_migration": {
										Type:     schema.TypeBool,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"current_task": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"state": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"start_date": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"speed_limit": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"migrate_speed": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"start_target_server": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"vm_template_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"region_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"project_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"target_server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
						"log_collect_status": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"exist_server": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"use_public_ip": {
							Type:     schema.TypeBool,
							Computed: true,
						},
						"clone_server": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"vm_id": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"name": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"clone_error": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"clone_state": {
										Type:     schema.TypeString,
										Computed: true,
									},
									"error_msg": {
										Type:     schema.TypeString,
										Computed: true,
									},
								},
							},
						},
					},
				},
			},
			"checks": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"params": {
							Type:     schema.TypeList,
							Computed: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"result": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_or_warn": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"error_params": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"state_action_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"replicatesize": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"totalsize": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"last_visit_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"stage_action_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"adjust_disk": {
				Type:     schema.TypeBool,
				Computed: true,
			},
		},
	}
}

func resourceSourceServerCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	createHttpUrl := "v3/sources"
	createPath := client.Endpoint + createHttpUrl
	createOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(buildCreateSourceServerBodyParams(d)),
	}

	createResp, err := client.Request("POST", createPath, &createOpt)
	if err != nil {
		return diag.Errorf("error creating SMS source server: %s", err)
	}

	createRespBody, err := utils.FlattenResponse(createResp)
	if err != nil {
		return diag.Errorf("error flattening creating source server response: %s", err)
	}

	id := utils.PathSearch("id", createRespBody, "").(string)
	if id == "" {
		return diag.Errorf("error creating SMS source server: can not found source server id in return")
	}

	d.SetId(id)

	if d.Get("migprojectid").(string) != "" {
		updateHttpUrl := "v3/sources/{source_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{source_id}", d.Id())
		updateBodyParams := map[string]interface{}{
			"migprojectid": utils.ValueIgnoreEmpty(d.Get("migprojectid")),
		}
		updateOpt := golangsdk.RequestOpts{
			KeepResponseBody: true,
			MoreHeaders: map[string]string{
				"Content-Type": "application/json",
			},
			JSONBody: utils.RemoveNil(updateBodyParams),
		}

		_, err = client.Request("PUT", updatePath, &updateOpt)
		if err != nil {
			return diag.Errorf("error updating source server: %s", err)
		}
	}

	if d.Get("copystate").(string) != "" {
		err := updateSourceServerChangeState(client, d)
		if err != nil {
			return diag.Errorf("error updating source server: %s", err)
		}
	}

	return resourceSourceServerRead(ctx, d, meta)
}

func resourceSourceServerRead(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	sourceServer, err := GetSourceServer(client, d.Id())
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error retrieving source server")
	}

	mErr := multierror.Append(nil,
		d.Set("ip", utils.PathSearch("ip", sourceServer, nil)),
		d.Set("name", utils.PathSearch("name", sourceServer, nil)),
		d.Set("hostname", utils.PathSearch("hostname", sourceServer, nil)),
		d.Set("enterprise_project_id", utils.PathSearch("enterprise_project_id", sourceServer, nil)),
		d.Set("add_date", utils.PathSearch("add_date", sourceServer, nil)),
		d.Set("os_type", utils.PathSearch("os_type", sourceServer, nil)),
		d.Set("os_version", utils.PathSearch("os_version", sourceServer, nil)),
		d.Set("oem_system", utils.PathSearch("oem_system", sourceServer, nil)),
		d.Set("state", utils.PathSearch("state", sourceServer, nil)),
		d.Set("connected", utils.PathSearch("connected", sourceServer, nil)),
		d.Set("firmware", utils.PathSearch("firmware", sourceServer, nil)),
		d.Set("init_target_server", flattenSmsSourceServerInitTargetServer(
			utils.PathSearch("init_target_server", sourceServer, nil))),
		d.Set("cpu_quantity", utils.PathSearch("cpu_quantity", sourceServer, nil)),
		d.Set("memory", utils.PathSearch("memory", sourceServer, nil)),
		d.Set("current_task", flattenSmsSourceServerCurrentTask(
			utils.PathSearch("current_task", sourceServer, nil))),
		d.Set("disks", flattenSmsSourceServerDisks(
			utils.PathSearch("disks", sourceServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("volume_groups", flattenSmsSourceServerVolumeGroups(
			utils.PathSearch("volume_groups", sourceServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("btrfs_list", flattenSmsSourceServerBtrfsList(
			utils.PathSearch("btrfs_list", sourceServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("networks", flattenSmsSourceServerNetworks(
			utils.PathSearch("networks", sourceServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("checks", flattenSmsSourceServerChecks(
			utils.PathSearch("checks", sourceServer, make([]interface{}, 0)).([]interface{}))),
		d.Set("migration_cycle", utils.PathSearch("migration_cycle", sourceServer, nil)),
		d.Set("state_action_time", utils.PathSearch("state_action_time", sourceServer, nil)),
		d.Set("replicatesize", utils.PathSearch("replicatesize", sourceServer, nil)),
		d.Set("totalsize", utils.PathSearch("totalsize", sourceServer, nil)),
		d.Set("last_visit_time", utils.PathSearch("last_visit_time", sourceServer, nil)),
		d.Set("stage_action_time", utils.PathSearch("stage_action_time", sourceServer, nil)),
		d.Set("agent_version", utils.PathSearch("agent_version", sourceServer, nil)),
		d.Set("has_tc", utils.PathSearch("has_tc", sourceServer, nil)),
		d.Set("adjust_disk", utils.PathSearch("adjust_disk", sourceServer, nil)),
	)

	return diag.FromErr(mErr.ErrorOrNil())
}

func flattenSmsSourceServerInitTargetServer(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"disks": flattenSmsSourceServerInitTargetServerDisks(
				utils.PathSearch("disks", param, make([]interface{}, 0)).([]interface{})),
			"volume_groups": flattenSmsSourceServerInitTargetServerVolumeGroups(
				utils.PathSearch("volume_groups", param, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenSmsSourceServerInitTargetServerDisks(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"device_use":     utils.PathSearch("device_use", params, nil),
			"used_size":      utils.PathSearch("used_size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
			"physical_volumes": flattenSmsSourceServerInitTargetServerDisksPhysicalVolumes(
				utils.PathSearch("physical_volumes", params, make([]interface{}, 0)).([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerInitTargetServerDisksPhysicalVolumes(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"device_use":     utils.PathSearch("device_use", params, nil),
			"file_system":    utils.PathSearch("file_system", params, nil),
			"index":          utils.PathSearch("index", params, nil),
			"mount_point":    utils.PathSearch("mount_point", params, nil),
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"inode_size":     utils.PathSearch("inode_size", params, nil),
			"used_size":      utils.PathSearch("used_size", params, nil),
			"uuid":           utils.PathSearch("uuid", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerInitTargetServerVolumeGroups(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"components": utils.PathSearch("components", params, nil),
			"free_size":  utils.PathSearch("free_size", params, nil),
			"logical_volumes": flattenSmsSourceServerInitTargetServerVolumeGroupsLogicalVolumes(
				utils.PathSearch("logical_volumes", params, make([]interface{}, 0)).([]interface{})),
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerInitTargetServerVolumeGroupsLogicalVolumes(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"block_count":    utils.PathSearch("block_count", params, nil),
			"block_size":     utils.PathSearch("block_size", params, nil),
			"file_system":    utils.PathSearch("file_system", params, nil),
			"inode_size":     utils.PathSearch("inode_size", params, nil),
			"inode_nums":     utils.PathSearch("inode_nums", params, nil),
			"device_use":     utils.PathSearch("device_use", params, nil),
			"mount_point":    utils.PathSearch("mount_point", params, nil),
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"used_size":      utils.PathSearch("used_size", params, nil),
			"free_size":      utils.PathSearch("free_size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerCurrentTask(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"id":                  utils.PathSearch("id", param, nil),
			"name":                utils.PathSearch("name", param, nil),
			"type":                utils.PathSearch("type", param, nil),
			"state":               utils.PathSearch("state", param, nil),
			"start_date":          utils.PathSearch("start_date", param, nil),
			"speed_limit":         utils.PathSearch("speed_limit", param, nil),
			"migrate_speed":       utils.PathSearch("migrate_speed", param, nil),
			"start_target_server": utils.PathSearch("start_target_server", param, nil),
			"vm_template_id":      utils.PathSearch("vm_template_id", param, nil),
			"region_id":           utils.PathSearch("region_id", param, nil),
			"project_name":        utils.PathSearch("project_name", param, nil),
			"project_id":          utils.PathSearch("project_id", param, nil),
			"target_server": flattenSmsSourceServerCurrentTaskTargetServer(
				utils.PathSearch("target_server", param, nil)),
			"log_collect_status": utils.PathSearch("log_collect_status", param, nil),
			"exist_server":       utils.PathSearch("exist_server", param, nil),
			"use_public_ip":      utils.PathSearch("use_public_ip", param, nil),
			"clone_server": flattenSmsSourceServerCurrentTaskCloneServer(
				utils.PathSearch("clone_server", param, make([]interface{}, 0)).([]interface{})),
		},
	}

	return rst
}

func flattenSmsSourceServerCurrentTaskTargetServer(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"vm_id": utils.PathSearch("vm_id", param, nil),
			"name":  utils.PathSearch("name", param, nil),
		},
	}

	return rst
}

func flattenSmsSourceServerCurrentTaskCloneServer(param interface{}) interface{} {
	if param == nil {
		return nil
	}
	rst := []map[string]interface{}{
		{
			"vm_id":       utils.PathSearch("vm_id", param, nil),
			"name":        utils.PathSearch("name", param, nil),
			"clone_error": utils.PathSearch("clone_error", param, nil),
			"clone_state": utils.PathSearch("clone_state", param, nil),
			"error_msg":   utils.PathSearch("error_msg", param, nil),
		},
	}

	return rst
}

func flattenSmsSourceServerDisks(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":            utils.PathSearch("name", params, nil),
			"partition_style": utils.PathSearch("partition_style", params, nil),
			"device_use":      utils.PathSearch("device_use", params, nil),
			"size":            utils.PathSearch("size", params, nil),
			"used_size":       utils.PathSearch("used_size", params, nil),
			"physical_volumes": flattenSmsSourceServerDisksPhysicalVolumes(
				utils.PathSearch("physical_volumes", params, make([]interface{}, 0)).([]interface{})),
			"os_disk":        utils.PathSearch("os_disk", params, nil),
			"relation_name":  utils.PathSearch("relation_name", params, nil),
			"inode_size":     utils.PathSearch("inode_size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerDisksPhysicalVolumes(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"device_use":       utils.PathSearch("device_use", params, nil),
			"file_system":      utils.PathSearch("file_system", params, nil),
			"index":            utils.PathSearch("index", params, nil),
			"mount_point":      utils.PathSearch("mount_point", params, nil),
			"name":             utils.PathSearch("name", params, nil),
			"size":             utils.PathSearch("size", params, nil),
			"used_size":        utils.PathSearch("used_size", params, nil),
			"inode_size":       utils.PathSearch("inode_size", params, nil),
			"inode_nums":       utils.PathSearch("inode_nums", params, nil),
			"uuid":             utils.PathSearch("uuid", params, nil),
			"size_per_cluster": utils.PathSearch("size_per_cluster", params, nil),
			"id":               utils.PathSearch("id", params, nil),
			"adjust_size":      utils.PathSearch("adjust_size", params, nil),
			"need_migration":   utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerVolumeGroups(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"components": utils.PathSearch("components", params, nil),
			"free_size":  utils.PathSearch("free_size", params, nil),
			"logical_volumes": flattenSmsSourceServerVolumeGroupsLogicalVolumes(
				utils.PathSearch("logical_volumes", params, make([]interface{}, 0)).([]interface{})),
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerVolumeGroupsLogicalVolumes(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"block_count":    utils.PathSearch("block_count", params, nil),
			"block_size":     utils.PathSearch("block_size", params, nil),
			"file_system":    utils.PathSearch("file_system", params, nil),
			"inode_size":     utils.PathSearch("inode_size", params, nil),
			"inode_nums":     utils.PathSearch("inode_nums", params, nil),
			"device_use":     utils.PathSearch("device_use", params, nil),
			"mount_point":    utils.PathSearch("mount_point", params, nil),
			"name":           utils.PathSearch("name", params, nil),
			"size":           utils.PathSearch("size", params, nil),
			"used_size":      utils.PathSearch("used_size", params, nil),
			"free_size":      utils.PathSearch("free_size", params, nil),
			"id":             utils.PathSearch("id", params, nil),
			"adjust_size":    utils.PathSearch("adjust_size", params, nil),
			"need_migration": utils.PathSearch("need_migration", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerBtrfsList(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":                     utils.PathSearch("name", params, nil),
			"label":                    utils.PathSearch("label", params, nil),
			"uuid":                     utils.PathSearch("uuid", params, nil),
			"device":                   utils.PathSearch("device", params, nil),
			"size":                     utils.PathSearch("size", params, nil),
			"nodesize":                 utils.PathSearch("nodesize", params, nil),
			"sectorsize":               utils.PathSearch("sectorsize", params, nil),
			"data_profile":             utils.PathSearch("data_profile", params, nil),
			"system_profile":           utils.PathSearch("system_profile", params, nil),
			"metadata_profile":         utils.PathSearch("metadata_profile", params, nil),
			"global_reserve1":          utils.PathSearch("global_reserve1", params, nil),
			"g_vol_used_size":          utils.PathSearch("g_vol_used_size", params, nil),
			"default_subvolid":         utils.PathSearch("default_subvolid", params, nil),
			"default_subvol_name":      utils.PathSearch("default_subvol_name", params, nil),
			"default_subvol_mountpath": utils.PathSearch("default_subvol_mountpath", params, nil),
			"subvolumn": flattenSmsSourceServerBtrfsListSubvolumn(
				utils.PathSearch("subvolumn", params, make([]interface{}, 0)).([]interface{})),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerBtrfsListSubvolumn(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"uuid":              utils.PathSearch("uuid", params, nil),
			"is_snapshot":       utils.PathSearch("is_snapshot", params, nil),
			"subvol_id":         utils.PathSearch("subvol_id", params, nil),
			"parent_id":         utils.PathSearch("parent_id", params, nil),
			"subvol_name":       utils.PathSearch("subvol_name", params, nil),
			"subvol_mount_path": utils.PathSearch("subvol_mount_path", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerNetworks(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"name":    utils.PathSearch("name", params, nil),
			"ip":      utils.PathSearch("ip", params, nil),
			"ipv6":    utils.ValueIgnoreEmpty(utils.PathSearch("ipv6", params, nil)),
			"netmask": utils.PathSearch("netmask", params, nil),
			"gateway": utils.PathSearch("gateway", params, nil),
			"mtu":     utils.ValueIgnoreEmpty(utils.PathSearch("mtu", params, nil)),
			"mac":     utils.PathSearch("mac", params, nil),
			"id":      utils.ValueIgnoreEmpty(utils.PathSearch("id", params, nil)),
		}
		rst = append(rst, m)
	}

	return rst
}

func flattenSmsSourceServerChecks(paramsList []interface{}) []interface{} {
	if len(paramsList) == 0 {
		return nil
	}
	rst := make([]interface{}, 0, len(paramsList))
	for _, params := range paramsList {
		m := map[string]interface{}{
			"id":            utils.PathSearch("id", params, nil),
			"params":        utils.PathSearch("params", params, nil),
			"name":          utils.PathSearch("name", params, nil),
			"result":        utils.PathSearch("result", params, nil),
			"error_code":    utils.PathSearch("error_code", params, nil),
			"error_or_warn": utils.PathSearch("error_or_warn", params, nil),
			"error_params":  utils.PathSearch("error_params", params, nil),
		}
		rst = append(rst, m)
	}

	return rst
}

func GetSourceServer(client *golangsdk.ServiceClient, sourceServerId string) (interface{}, error) {
	getHttpUrl := "v3/sources/{source_id}"
	getPath := client.Endpoint + getHttpUrl
	getPath = strings.ReplaceAll(getPath, "{source_id}", sourceServerId)
	getOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	getResp, err := client.Request("GET", getPath, &getOpt)
	if err != nil {
		return nil, err
	}
	getRespBody, err := utils.FlattenResponse(getResp)
	if err != nil {
		return nil, err
	}

	return getRespBody, nil
}

func resourceSourceServerUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	changeList := []string{
		"name", "migprojectid", "disks", "volume_groups",
	}
	if d.HasChanges(changeList...) {
		updateHttpUrl := "v3/sources/{source_id}"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{source_id}", d.Id())
		updateBodyParam := utils.RemoveNil(buildUpdateSourceServerBodyParams(d))
		if !reflect.DeepEqual(updateBodyParam, map[string]interface{}{}) {
			updateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				JSONBody: updateBodyParam,
			}

			_, err = client.Request("PUT", updatePath, &updateOpt)
			if err != nil {
				return diag.Errorf("error updating source server: %s", err)
			}
		}
	}

	changeList = []string{
		"disks", "volume_groups", "btrfs_list",
	}
	if d.HasChanges(changeList...) {
		updateHttpUrl := "v3/sources/{source_id}/diskinfo"
		updatePath := client.Endpoint + updateHttpUrl
		updatePath = strings.ReplaceAll(updatePath, "{source_id}", d.Id())
		updateBodyParams, isChange := buildUpdateSourceServerDiskInfoBodyParams(d)
		if isChange {
			updateOpt := golangsdk.RequestOpts{
				KeepResponseBody: true,
				MoreHeaders: map[string]string{
					"Content-Type": "application/json",
				},
				JSONBody: utils.RemoveNil(updateBodyParams),
			}

			_, err = client.Request("PUT", updatePath, &updateOpt)
			if err != nil {
				return diag.Errorf("error updating source server: %s", err)
			}
		}
	}

	changeList = []string{
		"copystate", "migration_cycle",
	}
	if d.HasChanges(changeList...) {
		err := updateSourceServerChangeState(client, d)
		if err != nil {
			return diag.Errorf("error updating source server: %s", err)
		}
	}

	return resourceSourceServerRead(ctx, d, meta)
}

func buildCreateSourceServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	bodyParams := map[string]interface{}{
		"ip":                  utils.ValueIgnoreEmpty(d.Get("ip")),
		"name":                utils.ValueIgnoreEmpty(d.Get("name")),
		"hostname":            utils.ValueIgnoreEmpty(d.Get("hostname")),
		"os_type":             utils.ValueIgnoreEmpty(d.Get("os_type")),
		"os_version":          utils.ValueIgnoreEmpty(d.Get("os_version")),
		"virtualization_type": utils.ValueIgnoreEmpty(d.Get("virtualization_type")),
		"linux_block_check":   utils.ValueIgnoreEmpty(d.Get("linux_block_check")),
		"firmware":            utils.ValueIgnoreEmpty(d.Get("firmware")),
		"cpu_quantity":        utils.ValueIgnoreEmpty(d.Get("cpu_quantity")),
		"memory":              utils.ValueIgnoreEmpty(d.Get("memory")),
		"disks":               buildSourceServerDiskParamsBody(d.Get("disks")),
		"btrfs_list":          buildSourceServerBtrfsFileSystemParamsBody(d.Get("btrfs_list")),
		"networks":            buildSourceServerNetWorkParamsBody(d.Get("networks")),
		"domain_id":           utils.ValueIgnoreEmpty(d.Get("domain_id")),
		"has_rsync":           utils.ValueIgnoreEmpty(d.Get("has_rsync")),
		"paravirtualization":  utils.ValueIgnoreEmpty(d.Get("paravirtualization")),
		"raw_devices":         utils.ValueIgnoreEmpty(d.Get("raw_devices")),
		"driver_files":        utils.ValueIgnoreEmpty(d.Get("driver_files")),
		"system_services":     utils.ValueIgnoreEmpty(d.Get("system_services")),
		"account_rights":      utils.ValueIgnoreEmpty(d.Get("account_rights")),
		"boot_loader":         utils.ValueIgnoreEmpty(d.Get("boot_loader")),
		"system_dir":          utils.ValueIgnoreEmpty(d.Get("system_dir")),
		"volume_groups":       buildSourceServerVolumeGroupParamsBody(d.Get("volume_groups")),
		"agent_version":       utils.ValueIgnoreEmpty(d.Get("agent_version")),
		"kernel_version":      utils.ValueIgnoreEmpty(d.Get("kernel_version")),
		"migration_cycle":     utils.ValueIgnoreEmpty(d.Get("migration_cycle")),
		"state":               utils.ValueIgnoreEmpty(d.Get("state")),
		"oem_system":          utils.ValueIgnoreEmpty(d.Get("oem_system")),
		"start_type":          utils.ValueIgnoreEmpty(d.Get("start_type")),
		"io_read_wait":        utils.ValueIgnoreEmpty(d.Get("io_read_wait")),
		"has_tc":              utils.ValueIgnoreEmpty(d.Get("has_tc")),
		"platform":            utils.ValueIgnoreEmpty(d.Get("platform")),
	}

	return bodyParams
}

func buildSourceServerDiskParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":             raw["name"],
				"device_use":       raw["device_use"],
				"size":             raw["size"],
				"used_size":        raw["used_size"],
				"physical_volumes": buildSourceServerDiskPhysicalVolumeParamsBody(raw["physical_volumes"]),
				"partition_style":  utils.ValueIgnoreEmpty(raw["partition_style"]),
				"os_disk":          utils.ValueIgnoreEmpty(raw["os_disk"]),
				"relation_name":    utils.ValueIgnoreEmpty(raw["relation_name"]),
				"inode_size":       utils.ValueIgnoreEmpty(raw["inode_size"]),
			}
		}
		return params
	}

	return nil
}

func buildSourceServerDiskPhysicalVolumeParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"device_use":       utils.ValueIgnoreEmpty(raw["device_use"]),
				"file_system":      utils.ValueIgnoreEmpty(raw["file_system"]),
				"index":            utils.ValueIgnoreEmpty(raw["index"]),
				"mount_point":      utils.ValueIgnoreEmpty(raw["mount_point"]),
				"name":             utils.ValueIgnoreEmpty(raw["name"]),
				"size":             utils.ValueIgnoreEmpty(raw["size"]),
				"used_size":        utils.ValueIgnoreEmpty(raw["used_size"]),
				"inode_size":       utils.ValueIgnoreEmpty(raw["inode_size"]),
				"inode_nums":       utils.ValueIgnoreEmpty(raw["inode_nums"]),
				"uuid":             utils.ValueIgnoreEmpty(raw["uuid"]),
				"size_per_cluster": utils.ValueIgnoreEmpty(raw["size_per_cluster"]),
			}
		}
		return params
	}

	return nil
}

func buildSourceServerBtrfsFileSystemParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":                     raw["name"],
				"label":                    raw["label"],
				"uuid":                     raw["uuid"],
				"device":                   raw["device"],
				"size":                     raw["size"],
				"nodesize":                 raw["nodesize"],
				"sectorsize":               raw["sectorsize"],
				"data_profile":             raw["data_profile"],
				"system_profile":           raw["system_profile"],
				"metadata_profile":         raw["metadata_profile"],
				"global_reserve1":          raw["global_reserve1"],
				"g_vol_used_size":          raw["g_vol_used_size"],
				"default_subvolid":         raw["default_subvolid"],
				"default_subvol_name":      raw["default_subvol_name"],
				"default_subvol_mountpath": raw["default_subvol_mountpath"],
				"subvolumn":                buildSourceServerBtrfsFileSystemSubvolumnParamsBody(raw["subvolumn"]),
			}
		}
		return params
	}

	return nil
}

func buildSourceServerBtrfsFileSystemSubvolumnParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"uuid":              raw["uuid"],
				"is_snapshot":       raw["is_snapshot"],
				"subvol_id":         raw["subvol_id"],
				"parent_id":         raw["parent_id"],
				"subvol_name":       raw["subvol_name"],
				"subvol_mount_path": raw["subvol_mount_path"],
			}
		}
		return params
	}

	return nil
}

func buildSourceServerNetWorkParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"name":    raw["name"],
				"ip":      raw["ip"],
				"netmask": raw["netmask"],
				"gateway": raw["gateway"],
				"mac":     raw["mac"],
				"ipv6":    utils.ValueIgnoreEmpty(raw["ipv6"]),
				"mtu":     utils.ValueIgnoreEmpty(raw["mtu"]),
				"id":      utils.ValueIgnoreEmpty(raw["id"]),
			}
		}
		return params
	}

	return nil
}

func buildSourceServerVolumeGroupParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"components":      utils.ValueIgnoreEmpty(raw["components"]),
				"free_size":       utils.ValueIgnoreEmpty(raw["free_size"]),
				"logical_volumes": buildSourceServerVolumeGroupLogicalVolumeParamsBody(raw["logical_volumes"]),
				"name":            utils.ValueIgnoreEmpty(raw["name"]),
				"size":            utils.ValueIgnoreEmpty(raw["size"]),
			}
		}
		return params
	}

	return nil
}

func buildSourceServerVolumeGroupLogicalVolumeParamsBody(rawParams interface{}) []map[string]interface{} {
	if rawArray, ok := rawParams.([]interface{}); ok {
		if len(rawArray) == 0 {
			return nil
		}

		params := make([]map[string]interface{}, len(rawArray))
		for i, v := range rawArray {
			raw := v.(map[string]interface{})
			params[i] = map[string]interface{}{
				"file_system": raw["file_system"],
				"inode_size":  raw["inode_size"],
				"mount_point": raw["mount_point"],
				"name":        raw["name"],
				"size":        raw["size"],
				"used_size":   raw["used_size"],
				"free_size":   raw["free_size"],
				"block_count": utils.ValueIgnoreEmpty(raw["block_count"]),
				"block_size":  utils.ValueIgnoreEmpty(raw["block_size"]),
				"inode_nums":  utils.ValueIgnoreEmpty(raw["inode_nums"]),
				"device_use":  utils.ValueIgnoreEmpty(raw["device_use"]),
			}
		}
		return params
	}

	return nil
}

func buildUpdateSourceServerBodyParams(d *schema.ResourceData) map[string]interface{} {
	disksParams := buildUpdateSourceServerDisksBodyParams(d)
	volumeGroupsParams := buildUpdateSourceServerVolumeGroupsBodyParams(d)

	name := ""
	if d.HasChange("name") {
		name = d.Get("name").(string)
	}

	migprojectid := ""
	if d.HasChange("migprojectid") {
		migprojectid = d.Get("migprojectid").(string)
	}

	bodyParams := map[string]interface{}{
		"name":          utils.ValueIgnoreEmpty(name),
		"migprojectid":  utils.ValueIgnoreEmpty(migprojectid),
		"disks":         disksParams,
		"volume_groups": volumeGroupsParams,
	}

	return bodyParams
}

func buildUpdateSourceServerVolumeGroupsBodyParams(d *schema.ResourceData) []map[string]interface{} {
	volumeGroupsAllParams := make([]map[string]interface{}, 0)
	isChange := false
	if d.HasChange("volume_groups") {
		for i := range d.Get("volume_groups").([]interface{}) {
			idIndex := fmt.Sprintf("volume_groups.%d.id", i)
			needMigrationIndex := fmt.Sprintf("volume_groups.%d.need_migration", i)
			adjustSizeIndex := fmt.Sprintf("volume_groups.%d.adjust_size", i)
			logicalVolumesIndex := fmt.Sprintf("volume_groups.%d.logical_volumes", i)

			if d.HasChanges(needMigrationIndex, adjustSizeIndex) {
				isChange = true
			}

			volumeGroupsLogicalVolumesAllParams := make([]map[string]interface{}, 0)
			paramsList := d.Get(logicalVolumesIndex).([]interface{})
			for j, rawParams := range paramsList {
				_, ok := rawParams.(map[string]interface{})
				if ok {
					idLogicalVolumesIndex := fmt.Sprintf("volume_groups.%d.logical_volumes.%d.id", i, j)
					needMigrationLogicalVolumesIndex := fmt.Sprintf("volume_groups.%d.logical_volumes.%d.need_migration", i, j)
					adjustSizeLogicalVolumesIndex := fmt.Sprintf("volume_groups.%d.logical_volumes.%d.adjust_size", i, j)
					m := utils.RemoveNil(map[string]interface{}{
						"id":             d.Get(idLogicalVolumesIndex),
						"need_migration": d.Get(needMigrationLogicalVolumesIndex),
						"adjust_size":    d.Get(adjustSizeLogicalVolumesIndex),
					})
					if !reflect.DeepEqual(m, map[string]interface{}{}) {
						if d.HasChanges(needMigrationLogicalVolumesIndex, adjustSizeLogicalVolumesIndex) {
							isChange = true
						}
						volumeGroupsLogicalVolumesAllParams = append(volumeGroupsLogicalVolumesAllParams, m)
					}
				}
			}

			volumeGroupsAllParams = append(volumeGroupsAllParams, utils.RemoveNil(map[string]interface{}{
				"id":              d.Get(idIndex),
				"adjust_size":     d.Get(adjustSizeIndex),
				"need_migration":  d.Get(needMigrationIndex),
				"logical_volumes": volumeGroupsLogicalVolumesAllParams,
			}))
		}
	}

	if isChange {
		return volumeGroupsAllParams
	}
	return make([]map[string]interface{}, 0)
}

func buildUpdateSourceServerDisksBodyParams(d *schema.ResourceData) []map[string]interface{} {
	disksAllParams := make([]map[string]interface{}, 0)
	isChange := false
	if d.HasChange("disks") {
		for i := range d.Get("disks").([]interface{}) {
			needMigrationIndex := fmt.Sprintf("disks.%d.need_migration", i)
			idIndex := fmt.Sprintf("disks.%d.id", i)
			adjustSizeIndex := fmt.Sprintf("disks.%d.adjust_size", i)
			physicalVolumesIndex := fmt.Sprintf("disks.%d.physical_volumes", i)

			if d.HasChanges(needMigrationIndex, adjustSizeIndex) {
				isChange = true
			}

			disksPhysicalVolumesAllParams := make([]map[string]interface{}, 0)
			paramsList := d.Get(physicalVolumesIndex).([]interface{})
			for j, rawParams := range paramsList {
				_, ok := rawParams.(map[string]interface{})
				if ok {
					idPhysicalVolumesIndex := fmt.Sprintf("disks.%d.physical_volumes.%d.id", i, j)
					needMigrationPhysicalVolumesIndex := fmt.Sprintf("disks.%d.physical_volumes.%d.need_migration", i, j)
					adjustSizePhysicalVolumesIndex := fmt.Sprintf("disks.%d.physical_volumes.%d.adjust_size", i, j)
					m := utils.RemoveNil(map[string]interface{}{
						"id":             d.Get(idPhysicalVolumesIndex),
						"need_migration": d.Get(needMigrationPhysicalVolumesIndex),
						"adjust_size":    d.Get(adjustSizePhysicalVolumesIndex),
					})
					if !reflect.DeepEqual(m, map[string]interface{}{}) {
						if d.HasChanges(needMigrationPhysicalVolumesIndex, adjustSizePhysicalVolumesIndex) {
							isChange = true
						}
						disksPhysicalVolumesAllParams = append(disksPhysicalVolumesAllParams, m)
					}
				}
			}

			disksAllParams = append(disksAllParams, utils.RemoveNil(map[string]interface{}{
				"id":               d.Get(idIndex),
				"adjust_size":      d.Get(adjustSizeIndex),
				"need_migration":   d.Get(needMigrationIndex),
				"physical_volumes": disksPhysicalVolumesAllParams,
			}))
		}
	}

	if isChange {
		return disksAllParams
	}
	return make([]map[string]interface{}, 0)
}

func buildUpdateSourceServerDiskInfoBodyParams(d *schema.ResourceData) (map[string]interface{}, bool) {
	disksParams, isChangeDisks := buildUpdateSourceServerDiskInfoDisksBodyParams(d)
	volumeGroupsParams, isChangeVolumeGroups := buildUpdateSourceServerDiskInfoVolumeGroupsBodyParams(d)
	btrfsListParams, isChangeBtrfsList := buildUpdateSourceServerDiskInfoBtrfsListBodyParams(d)

	bodyParams := map[string]interface{}{
		"disks":        disksParams,
		"volumegroups": volumeGroupsParams,
		"btrfs_list":   btrfsListParams,
	}

	return bodyParams, isChangeDisks || isChangeVolumeGroups || isChangeBtrfsList
}

func buildUpdateSourceServerDiskInfoDisksBodyParams(d *schema.ResourceData) ([]map[string]interface{}, bool) {
	disksAllParams := make([]map[string]interface{}, 0)
	isChange := false
	if d.HasChange("disks") {
		for i := range d.Get("disks").([]interface{}) {
			nameIndex := fmt.Sprintf("disks.%d.name", i)
			partitionStyleIndex := fmt.Sprintf("disks.%d.partition_style", i)
			deviceUseIndex := fmt.Sprintf("disks.%d.device_use", i)
			sizeIndex := fmt.Sprintf("disks.%d.size", i)
			usedSizeIndex := fmt.Sprintf("disks.%d.used_size", i)
			osDiskIndex := fmt.Sprintf("disks.%d.os_disk", i)
			relationNameIndex := fmt.Sprintf("disks.%d.relation_name", i)
			inodeSizeIndex := fmt.Sprintf("disks.%d.inode_size", i)
			physicalVolumesIndex := fmt.Sprintf("disks.%d.physical_volumes", i)

			if d.HasChanges(nameIndex, partitionStyleIndex, deviceUseIndex, sizeIndex, usedSizeIndex, osDiskIndex,
				relationNameIndex, inodeSizeIndex, physicalVolumesIndex) {
				isChange = true
			}

			disksPhysicalVolumesAllParams := make([]map[string]interface{}, 0)
			paramsList := d.Get(physicalVolumesIndex).([]interface{})
			for _, rawParams := range paramsList {
				params, ok := rawParams.(map[string]interface{})
				if ok {
					m := utils.RemoveNil(map[string]interface{}{
						"device_use":       utils.ValueIgnoreEmpty(params["device_use"]),
						"file_system":      utils.ValueIgnoreEmpty(params["file_system"]),
						"index":            utils.ValueIgnoreEmpty(params["index"]),
						"mount_point":      utils.ValueIgnoreEmpty(params["mount_point"]),
						"name":             utils.ValueIgnoreEmpty(params["name"]),
						"size":             utils.ValueIgnoreEmpty(params["size"]),
						"used_size":        utils.ValueIgnoreEmpty(params["used_size"]),
						"inode_size":       utils.ValueIgnoreEmpty(params["inode_size"]),
						"inode_nums":       utils.ValueIgnoreEmpty(params["inode_nums"]),
						"uuid":             utils.ValueIgnoreEmpty(params["uuid"]),
						"size_per_cluster": utils.ValueIgnoreEmpty(params["size_per_cluster"]),
					})
					if !reflect.DeepEqual(m, map[string]interface{}{}) {
						disksPhysicalVolumesAllParams = append(disksPhysicalVolumesAllParams, m)
					}
				}
			}

			disksAllParams = append(disksAllParams, utils.RemoveNil(map[string]interface{}{
				"name":             d.Get(nameIndex),
				"partition_style":  utils.ValueIgnoreEmpty(d.Get(partitionStyleIndex)),
				"device_use":       d.Get(deviceUseIndex),
				"size":             d.Get(sizeIndex),
				"used_size":        d.Get(usedSizeIndex),
				"os_disk":          utils.ValueIgnoreEmpty(d.Get(osDiskIndex)),
				"relation_name":    utils.ValueIgnoreEmpty(d.Get(relationNameIndex)),
				"inode_size":       utils.ValueIgnoreEmpty(d.Get(inodeSizeIndex)),
				"physical_volumes": disksPhysicalVolumesAllParams,
			}))
		}
	}

	return disksAllParams, isChange
}

func buildUpdateSourceServerDiskInfoVolumeGroupsBodyParams(d *schema.ResourceData) ([]map[string]interface{}, bool) {
	volumeGroupsAllParams := make([]map[string]interface{}, 0)
	isChange := false
	if d.HasChange("volumegroups") {
		for i := range d.Get("volumegroups").([]interface{}) {
			componentsIndex := fmt.Sprintf("volumegroups.%d.components", i)
			freeSizeIndex := fmt.Sprintf("volumegroups.%d.free_size", i)
			nameIndex := fmt.Sprintf("volumegroups.%d.name", i)
			sizeIndex := fmt.Sprintf("volumegroups.%d.size", i)
			logicalVolumesIndex := fmt.Sprintf("volumegroups.%d.logical_volumes", i)

			if d.HasChanges(componentsIndex, freeSizeIndex, nameIndex, sizeIndex, logicalVolumesIndex) {
				isChange = true
			}

			volumeGroupsLogicalVolumesAllParams := make([]map[string]interface{}, 0)
			paramsList := d.Get(logicalVolumesIndex).([]interface{})
			for _, rawParams := range paramsList {
				params, ok := rawParams.(map[string]interface{})
				if ok {
					m := utils.RemoveNil(map[string]interface{}{
						"block_count": utils.ValueIgnoreEmpty(params["block_count"]),
						"block_size":  utils.ValueIgnoreEmpty(params["block_size"]),
						"file_system": params["file_system"],
						"inode_size":  params["inode_size"],
						"inode_nums":  utils.ValueIgnoreEmpty(params["inode_nums"]),
						"device_use":  utils.ValueIgnoreEmpty(params["device_use"]),
						"mount_point": params["mount_point"],
						"name":        params["name"],
						"size":        params["size"],
						"used_size":   params["used_size"],
						"free_size":   params["free_size"],
					})
					if !reflect.DeepEqual(m, map[string]interface{}{}) {
						volumeGroupsLogicalVolumesAllParams = append(volumeGroupsLogicalVolumesAllParams, m)
					}
				}
			}

			volumeGroupsAllParams = append(volumeGroupsAllParams, utils.RemoveNil(map[string]interface{}{
				"components":      utils.ValueIgnoreEmpty(d.Get(componentsIndex)),
				"free_size":       utils.ValueIgnoreEmpty(d.Get(freeSizeIndex)),
				"name":            utils.ValueIgnoreEmpty(d.Get(nameIndex)),
				"size":            utils.ValueIgnoreEmpty(d.Get(sizeIndex)),
				"logical_volumes": volumeGroupsLogicalVolumesAllParams,
			}))
		}
	}

	return volumeGroupsAllParams, isChange
}

func buildUpdateSourceServerDiskInfoBtrfsListBodyParams(d *schema.ResourceData) ([]map[string]interface{}, bool) {
	btrfsListAllParams := make([]map[string]interface{}, 0)
	isChange := false
	if d.HasChange("btrfs_list") {
		for i := range d.Get("btrfs_list").([]interface{}) {
			nameIndex := fmt.Sprintf("btrfs_list.%d.name", i)
			labelIndex := fmt.Sprintf("btrfs_list.%d.label", i)
			uuidIndex := fmt.Sprintf("btrfs_list.%d.uuid", i)
			deviceIndex := fmt.Sprintf("btrfs_list.%d.device", i)
			sizeIndex := fmt.Sprintf("btrfs_list.%d.size", i)
			nodeSizeIndex := fmt.Sprintf("btrfs_list.%d.nodesize", i)
			sectorsizeIndex := fmt.Sprintf("btrfs_list.%d.sectorsize", i)
			dataProfileIndex := fmt.Sprintf("btrfs_list.%d.data_profile", i)
			systemProfileIndex := fmt.Sprintf("btrfs_list.%d.system_profile", i)
			metadataProfileIndex := fmt.Sprintf("btrfs_list.%d.metadata_profile", i)
			globalReserve1Index := fmt.Sprintf("btrfs_list.%d.global_reserve1", i)
			gvolUsedSizeIndex := fmt.Sprintf("btrfs_list.%d.g_vol_used_size", i)
			defaultSubvolidIndex := fmt.Sprintf("btrfs_list.%d.default_subvolid", i)
			defaultSubvolNameIndex := fmt.Sprintf("btrfs_list.%d.default_subvol_name", i)
			defaultSubvolMountpathIndex := fmt.Sprintf("btrfs_list.%d.default_subvol_mountpath", i)
			subvolumnIndex := fmt.Sprintf("btrfs_list.%d.subvolumn", i)

			if d.HasChanges(nameIndex, labelIndex, uuidIndex, deviceIndex, sizeIndex, nodeSizeIndex, sectorsizeIndex,
				dataProfileIndex, systemProfileIndex, metadataProfileIndex, globalReserve1Index, gvolUsedSizeIndex,
				defaultSubvolidIndex, defaultSubvolNameIndex, defaultSubvolMountpathIndex, subvolumnIndex) {
				isChange = true
			}

			btrfsListSubvolumnAllParams := make([]map[string]interface{}, 0)
			paramsList := d.Get(subvolumnIndex).([]interface{})
			for _, rawParams := range paramsList {
				params, ok := rawParams.(map[string]interface{})
				if ok {
					m := utils.RemoveNil(map[string]interface{}{
						"uuid":              params["uuid"],
						"is_snapshot":       params["is_snapshot"],
						"subvol_id":         params["subvol_id"],
						"parent_id":         params["parent_id"],
						"subvol_name":       params["subvol_name"],
						"subvol_mount_path": params["subvol_mount_path"],
					})
					if !reflect.DeepEqual(m, map[string]interface{}{}) {
						btrfsListSubvolumnAllParams = append(btrfsListSubvolumnAllParams, m)
					}
				}
			}

			btrfsListAllParams = append(btrfsListAllParams, utils.RemoveNil(map[string]interface{}{
				"name":                     d.Get(nameIndex),
				"label":                    d.Get(labelIndex),
				"uuid":                     d.Get(uuidIndex),
				"device":                   d.Get(deviceIndex),
				"size":                     d.Get(sizeIndex),
				"nodesize":                 d.Get(nodeSizeIndex),
				"sectorsize":               d.Get(sectorsizeIndex),
				"data_profile":             d.Get(dataProfileIndex),
				"system_profile":           d.Get(systemProfileIndex),
				"metadata_profile":         d.Get(metadataProfileIndex),
				"global_reserve1":          d.Get(globalReserve1Index),
				"g_vol_used_size":          d.Get(gvolUsedSizeIndex),
				"default_subvolid":         d.Get(defaultSubvolidIndex),
				"default_subvol_name":      d.Get(defaultSubvolNameIndex),
				"default_subvol_mountpath": d.Get(defaultSubvolMountpathIndex),
				"subvolumn":                btrfsListSubvolumnAllParams,
			}))
		}
	}

	return btrfsListAllParams, isChange
}

func updateSourceServerChangeState(client *golangsdk.ServiceClient, d *schema.ResourceData) error {
	updateHttpUrl := "v3/sources/{source_id}/changestate"
	updatePath := client.Endpoint + updateHttpUrl
	updatePath = strings.ReplaceAll(updatePath, "{source_id}", d.Id())
	updateBodyParams := map[string]interface{}{
		"copystate":      utils.ValueIgnoreEmpty(d.Get("copystate")),
		"migrationcycle": utils.ValueIgnoreEmpty(d.Get("migration_cycle")),
	}
	updateOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
		JSONBody: utils.RemoveNil(updateBodyParams),
	}

	_, err := client.Request("PUT", updatePath, &updateOpt)
	return err
}

func resourceSourceServerDelete(_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	cfg := meta.(*config.Config)
	client, err := cfg.SmsV3Client(cfg.GetRegion(d))
	if err != nil {
		return diag.Errorf("error creating SMS client: %s", err)
	}

	deleteHttpUrl := "v3/sources/{source_id}"
	deletePath := client.Endpoint + deleteHttpUrl
	deletePath = strings.ReplaceAll(deletePath, "{source_id}", d.Id())
	deleteOpt := golangsdk.RequestOpts{
		KeepResponseBody: true,
		MoreHeaders: map[string]string{
			"Content-Type": "application/json",
		},
	}

	_, err = client.Request("DELETE", deletePath, &deleteOpt)
	if err != nil {
		return common.CheckDeletedDiag(d, err, "error deleting SMS source server")
	}

	return nil
}
