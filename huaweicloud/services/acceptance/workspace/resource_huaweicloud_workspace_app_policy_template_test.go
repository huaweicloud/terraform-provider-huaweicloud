package workspace

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"

	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/acceptance"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/workspace"
)

func getAppPolicyTemplateFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	client, err := conf.NewServiceClient("appstream", acceptance.HW_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating Workspace APP client: %s", err)
	}
	return workspace.GetAppPolicyTemplateById(client, state.Primary.ID)
}

func TestAccAppPolicyTemplate_basic(t *testing.T) {
	var (
		policyTemplate interface{}
		resourceName   = "huaweicloud_workspace_app_policy_template.test"
		name           = acceptance.RandomAccResourceName()
		updateName     = acceptance.RandomAccResourceName()
	)

	rc := acceptance.InitResourceCheck(
		resourceName,
		&policyTemplate,
		getAppPolicyTemplateFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccAppPolicyTemplate_basic(name, "Created by terraform script"),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", name),
					resource.TestCheckResourceAttr(resourceName, "description", "Created by terraform script"),
					resource.TestCheckResourceAttrSet(resourceName, "policies"),
				),
			},
			{
				Config: testAccAppPolicyTemplate_basic(updateName, ""),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", updateName),
					resource.TestCheckResourceAttr(resourceName, "description", ""),
					resource.TestCheckResourceAttrSet(resourceName, "policies"),
				),
			},
		},
	})
}

func testAccAppPolicyTemplate_basic(name, description string) string {
	return fmt.Sprintf(`
resource "huaweicloud_workspace_app_policy_template" "test" {
  name        = "%[1]s"
  description = "%[2]s"
  policies    = jsonencode({
    "peripherals": {
      "usb_port_redirection": {
        "usb_enable": true,
        "options": {
          "usb_image_enable": true,
          "usb_video_enable": true,
          "usb_printer_enable": true,
          "usb_storage_enable": false,
          "wireless_devices_enable": false,
          "network_devices_enable": false,
          "usb_smart_card_enable": true,
          "other_usb_devices_enable": false,
          "usb_redirection_customization_policy": "",
          "usb_redirection_mode": "Classical mode"
        }
      },
      "device_redirection": {
        "printer_redirection": {
          "printer_enable": true,
          "options": {
            "sync_client_default_printer_enable": true,
            "universal_printer_driver": "Default"
          }
        },
        "session_printer": {
          "session_printer_enable": false,
          "options": {
            "session_printer_customization_policy": ""
          }
        },
        "camera_redirection": {
          "video_compress_enable": true,
          "options": {
            "camera_frame_rate": 15,
            "camera_max_width": 3000,
            "camera_max_heigth": 3000,
            "camera_compression_method": "H.264"
          }
        },
        "twain_redirection_enable": true,
        "image_compression_level": "medium"
      },
      "usb_device_common": {
        "pcsc_smart_card_enable": "Disable",
        "common_options": {
          "remove_smart_card_disconnect_enable": false
        }
      },
      "serial_port_redirection": {
        "serial_port_enable": false,
        "options": {
          "auto_connect_enable": false
        }
      }
    },
    "audio": {
      "audio_redirection_enable": true,
      "play_redirection_enable": true,
      "play_classification": "Music Play",
      "record_redirection_enable": true,
      "record_classification": "Speech Call"
    },
    "client": {
      "automatic_reconnection_interval": 5,
      "session_persistence_time": 180,
      "forbid_screen_capture": true
    },
    "display": {
      "display_level": "LEVEL4",
      "options": {
        "display_bandwidth": 20000,
        "frame_rate": 60,
        "video_frame_rate": 60,
        "min_image_cache": 200,
        "smoothing_factor": 60,
        "lossless_compression_mode": "Basic Compression",
        "deep_compression_options": {
          "deep_compression_level": "Compression grade 0"
        },
        "lossy_compression_quality": 85,
        "color_enhancement_enable": false,
        "quality_bandwidth_first": "Bandwidth First",
        "video_bit_rate_options": {
          "average_video_bit_rate": 18000
        },
        "peak_video_bit_rate": 18000,
        "video_quality_options": {
          "average_video_quality": 15,
          "lowest_video_quality": 25,
          "highest_video_quality": 7
        },
        "gop_size": 100,
        "encoding_preset": "Preset 1"
      },
      "rendering_acceleration_enable": true,
      "rendering_acceleration_options": {
        "video_acceleration_enhancement_enable": true,
        "video_optimization_enable": true,
        "gpu_color_optimization_enable": true
      },
      "video_card_memory_size": 64,
      "driver_delegation_mode_enable": true,
      "driver_delegation_latency": 80,
      "video_latency": 80,
      "change_resolution_vm": true
    },
    "file_and_clipboard": {
      "bypass_in_remote_app_enable": true,
      "file_redirection": {
        "redirection_mode": "READ_AND_WRITE",
        "options": {
          "fluid_control_switch_enable": true,
          "fluid_control_options": {
            "good_network_latency": 30,
            "normal_network_latency": 70,
            "poor_network_latency": 100,
            "reducing_step": 20,
            "slow_increasing_step": 10,
            "quick_increasing_step": 20,
            "start_speed": 1024,
            "test_block_size": 64,
            "test_time_gap": 10000
          },
          "compression_switch_enable": true,
          "compression_switch_options": {
            "compression_threshold": 512,
            "minimum_compression_rate": 900
          },
          "linux_file_size_supported_enable": true,
          "linux_file_size_supported_options": {
            "linux_file_size_supported_threshold": 100
          },
          "linux_root_mount_switch_enable": true,
          "linux_root_dir_list": "\\var\\log",
          "linux_file_mount_path": "\\media|\\Volumes|\\swdb\\mnt|\\home|\\storage|\\tmp|\\run\\media",
          "linux_fixed_drive_file_system_format": "",
          "linux_removable_drive_file_system_format": "vfat|ntfs|msdos|fuseblk|sdcardfs|exfat|fuse.fdredir",
          "linux_cdrom_drive_file_system_format": "cd9660|iso9660|udf",
          "linux_network_drive_file_system_format": "smbfs|afpfs|cifs",
          "path_separator": "|",
          "fixed_drive_enable": true,
          "removable_drive_enable": true,
          "cd_rom_drive_enable": true,
          "network_drive_enable": true
        },
        "vm_send_file_client": true,
        "redirection_send_file_options": {
          "read_write_speed": 0
        }
      },
      "fd_mobile_client_redir_enable": true,
      "clipboard_redirection": "TWO_WAY_ENABLED",
      "clipboard_redirection_options": {
        "rich_text_redirection_enable": true,
        "rich_text_clipboard_redirection": "TWO_WAY_ENABLED",
        "clipboard_file_redirection_enable": true,
        "file_clipboard_redirection": "TWO_WAY_ENABLED",
        "clipboard_length_limit_cts_enable": true,
        "clipboard_length_limit_cts": 1,
        "clipboard_length_limit_stc_enable": true,
        "clipboard_length_limit_stc": 1
      }
    },
    "session": {
      "sbc": {
        "sbc_automatic_disconnection": "AUTO_DISCONNECT",
        "sbc_automatic_disconnection_options": {
          "disconnection_waiting_time": 15,
          "sbc_auto_logout": true,
          "auto_logout_options": {
            "sbc_logout_waiting_time": 60
          }
        }
      }
    },
    "virtual_channel": {
      "virtual_channel_control_enable": false,
      "options": {
        "custom_virtual_channel_name": "",
        "virtual_channel_plugin_details": "",
        "third_party_plugin_name": ""
      }
    },
    "keyboard_mouse": {
      "mouse_feedback": "SELFADAPTION",
      "mouse_simulation_mode": "ABSOLUTE_POSITION",
      "external_cursor_feedback": false
    },
    "bandwidth": {
      "intelligent_data_transport_flag": "DISABLE",
      "total_bandwidth_control_enable": false,
      "options": {
        "total_bandwidth_control_value": 30000,
        "display_bandwidth_percentage_enable": false,
        "display_bandwidth_percentage_options": {
          "display_bandwidth_percentage_value": 65
        },
        "multimedia_bandwidth_percentage_enable": false,
        "multimedia_bandwidth_percentage_options": {
          "multimedia_bandwidth_percentage_value": 50
        },
        "usb_bandwidth_percentage_enable": false,
        "usb_bandwidth_percentage_options": {
          "usb_bandwidth_percentage_value": 100
        },
        "pcsc_bandwidth_percentage_enable": false,
        "pcsc_bandwidth_percentage_options": {
          "pcsc_bandwidth_percentage_value": 5
        },
        "twain_bandwidth_percentage_enable": false,
        "twain_bandwidth_percentage_options": {
          "twain_bandwidth_percentage_value": 15
        },
        "printer_bandwidth_percentage_enable": false,
        "printer_bandwidth_percentage_options": {
          "printer_bandwidth_percentage_value": 5
        },
        "com_bandwidth_percentage_enable": false,
        "com_bandwidth_percentage_options": {
          "com_bandwidth_percentage_value": 3
        },
        "file_redirection_bandwidth_percentage_enable": false,
        "file_redirection_bandwidth_percentage_options": {
          "file_redirection_bandwidth_percentage_value": 30
        },
        "clipboard_bandwidth_percentage_enable": false,
        "clipboard_bandwidth_percentage_options": {
          "clipboard_bandwidth_percentage_value": 3
        },
        "secure_channel_bandwidth_percentage_enable": false,
        "secure_channel_bandwidth_percentage_options": {
          "secure_channel_bandwidth_percentage_value": 30
        },
        "camera_bandwidth_percentage_enable": false,
        "camera_bandwidth_percentage_options": {
          "camera_bandwidth_percentage_value": 30
        },
        "virtual_channel_bandwidth_percentage_enable": false,
        "virtual_channel_bandwidth_percentage_options": {
          "virtual_channel_bandwidth_percentage_value": 65
        }
      },
      "display_bandwidth_control_enable": false,
      "display_bandwidth_control_options": {
        "display_bandwidth_control_value": 20000
      },
      "multimedia_bandwidth_control_enable": false,
      "multimedia_bandwidth_control_options": {
        "multimedia_bandwidth_control_value": 15000
      },
      "usb_bandwidth_control_enable": false,
      "usb_bandwidth_control_options": {
        "usb_bandwidth_control_value": 30000
      },
      "pcsc_bandwidth_control_enable": false,
      "pcsc_bandwidth_control_options": {
        "pcsc_bandwidth_control_value": 2000
      },
      "twain_bandwidth_control_enable": false,
      "twain_bandwidth_control_options": {
        "twain_bandwidth_control_value": 5000
      },
      "printer_bandwidth_control_enable": false,
      "printer_bandwidth_control_options": {
        "printer_bandwidth_control_value": 2000
      },
      "com_bandwidth_control_enable": false,
      "com_bandwidth_control_options": {
        "com_bandwidth_control_value": 1000
      },
      "file_redirection_bandwidth_control_enable": false,
      "file_redirection_bandwidth_control_options": {
        "file_redirection_bandwidth_control_value": 10000
      },
      "clipboard_bandwidth_control_enable": false,
      "clipboard_bandwidth_control_options": {
        "clipboard_bandwidth_control_value": 1000
      },
      "secure_channel_bandwidth_control_enable": false,
      "secure_channel_bandwidth_control_options": {
        "secure_channel_bandwidth_control_value": 10000
      },
      "camera_bandwidth_control_enable": false,
      "camera_bandwidth_control_options": {
        "camera_bandwidth_control_value": 10000
      },
      "virtual_channel_bandwidth_control_enable": false,
      "virtual_channel_bandwidth_control_options": {
        "virtual_channel_bandwidth_control_value": 20000
      }
    },
    "custom": {
      "custom_configuration1_enable": true,
      "options": {
        "custom_configuration1_rule": "",
        "rail_transparent_config": {
          "select_policy": 0,
          "transparent_list_custom": ""
        }
      }
    },
    "cloud_storage": {
      "cloud_storage_enable": true,
      "options": {
        "cloud_storage_rule": jsonencode({
          "personal_folder": {
            "storage_paths": [
              {
                "storage_type": "HUAWEI_CLOUD_SFS",
                "project_config_id": ""
              }
            ]
          },
          "shared_folder": {}
        })
      }
    },
    "user_profile": {
      "user_profile_enable": true,
      "options": {
        "user_profile_rule": jsonencode({
          "ProfileEnable": {
            "status": 1,
            "value": 0
          },
          "ProfileLocations": {
            "status": 1,
            "type": "SMBLocation",
            "user_profile_locations": [{
              "id": 0,
              "container_type": "HUAWEI_CLOUD_SFS",
              "path": ""
            }]
          },
          "ProfileVhdType": {
            "status": 1,
            "value": "vhd"
          },
          "ProfileIsDynamic": {
            "status": 1,
            "value": 1
          },
          "ProfileVhdFileSize": {
            "status": 1,
            "value": 30000
          },
          "ProfileDiffDiskParentFolderPath": {
            "status": 2,
            "value": "%%TEMP%%" # Please watching this format handling about %% characters
          },
          "ProfileSidDirNameRule": {
            "status": 2,
            "value": "%%sid%%_%%username%%" # Please watching this format handling about %% characters
          },
          "ProfileSidDirMatchRule": {
            "status": 2,
            "value": "%%sid%%_%%username%%" # Please watching this format handling about %% characters
          },
          "ProfileVhdFileNameRule": {
            "status": 2,
            "value": "Profile_%%username%%" # Please watching this format handling about %% characters
          },
          "ProfileVhdFileMatchRule": {
            "status": 2,
            "value": "Profile*"
          },
          "ProfileNoProfileContainingFolder": {
            "status": 2,
            "value": 0
          },
          "ProfileFlipFlopProfileDirectoryName": {
            "status": 2,
            "value": 0
          },
          "ProfileRemoveOrphanedOSTFilesOnLogoff": {
            "status": 2,
            "value": 0
          },
          "ProfileVolumeWaitTimeMS": {
            "status": 2,
            "value": 20000
          },
          "ProfileVHDXSectorSize": {
            "status": 2,
            "value": 0
          },
          "ProfileRebootOnUserLogoff": {
            "status": 2,
            "value": 0
          },
          "ProfileRedirectType": {
            "status": 2,
            "value": 2
          },
          "ProfileSetTempToLocalPath": {
            "status": 2,
            "value": 3
          },
          "ProfileRoamSearch": {
            "status": 2,
            "value": 0
          },
          "ProfileLockedRetryInterval": {
            "status": 2,
            "value": 5
          },
          "ProfileVhdAttachedRetryTimes": {
            "status": 2,
            "value": 60
          },
          "ProfileShutdownOnUserLogoff": {
            "status": 2,
            "value": 0
          },
          "ProfileOutlookCachedMode": {
            "status": 2,
            "value": 1
          },
          "ProfileKeepLocalDir": {
            "status": 2,
            "value": 0
          },
          "ProfilePreventLoginWithFailure": {
            "status": 2,
            "value": 0
          },
          "ProfileDirSDDL": {
            "status": 2,
            "value": "N/A"
          },
          "ProfileRedirXmlLocation": {
            "status": 2,
            "value": "N/A"
          },
          "ProfileAttachVHDSDDL": {
            "status": 2,
            "value": "N/A"
          },
          "ProfileDeleteLocalProfile": {
            "status": 2,
            "value": 0
          },
          "ProfileInstallAppxPackages": {
            "status": 2,
            "value": 1
          },
          "ProfilePreventLoginWithTempProfile": {
            "status": 2,
            "value": 0
          },
          "ProfileVhdReattachInterval": {
            "status": 2,
            "value": 10
          },
          "SIDDirSDDL": {
            "status": 2,
            "value": "N/A"
          },
          "ProfileCleanOutNotifications": {
            "status": 2,
            "value": 1
          },
          "ProfileRoamIdentity": {
            "status": 2,
            "value": 0
          },
          "ProfileAccessNetworkAsComputerObject": {
            "status": 2,
            "value": 0
          },
          "ProfileType": {
            "status": 2,
            "value": 0
          },
          "ProfileLockedRetryCount": {
            "status": 2,
            "value": 12
          },
          "ODFCEnable": {
            "status": 1,
            "value": 0
          },
          "ODFCLocations": {
            "status": 1,
            "type": "SMBLocation",
            "user_profile_locations": [{
              "id": 0,
              "container_type": "HUAWEI_CLOUD_SFS",
              "path": ""
            }]
          },
          "ODFCVHDType": {
            "status": 1,
            "value": "vhd"
          },
          "ODFCIsDynamic": {
            "status": 1,
            "value": 1
          },
          "ODFCVHDFileSize": {
            "status": 1,
            "value": 30000
          },
          "ODFCVHDNamePattern": {
            "status": 2,
            "value": "ODFC_%%username%%" # Please watching this format handling about %% characters
          },
          "ODFCSIDDirNamePattern": {
            "status": 2,
            "value": "%%sid%%_%%username%%" # Please watching this format handling about %% characters
          },
          "ODFCSIDDirNameMatch": {
            "status": 2,
            "value": "%%sid%%_%%username%%" # Please watching this format handling about %% characters
          },
          "ODFCDiffDiskParentFolderPath": {
            "status": 2,
            "value": "%%TEMP%%" # Please watching this format handling about %% characters
          },
          "ODFCFlipFlopProfileDirectoryName": {
            "status": 2,
            "value": 0
          },
          "ODFCVHDNameMatch": {
            "status": 2,
            "value": "ODFC*"
          },
          "ODFCNoProfileContainingFolder": {
            "status": 2,
            "value": 0
          },
          "IncludeOfficeActivation": {
            "status": 2,
            "value": 1
          },
          "IncludeOneNote": {
            "status": 2,
            "value": 0
          },
          "ODFCVHDXSectorSize": {
            "status": 2,
            "value": 0
          },
          "OutlookFolderPath": {
            "status": 2,
            "value": "%%userprofile%%\\AppData\\Local\\Microsoft\\Outlook" # Please watching this format handling about %% characters
          },
          "ODFCVolumeWaitTimeMS": {
            "status": 2,
            "value": 20000
          },
          "IncludeOutlookPersonalization": {
            "status": 2,
            "value": 1
          },
          "ODFCVHDReAttachIntervalSeconds": {
            "status": 2,
            "value": 10
          },
          "ODFCMirrorLocalOSTToVHD": {
            "status": 2,
            "value": 0
          },
          "ODFCIncludeTeams": {
            "status": 2,
            "value": 0
          },
          "ODFCRemoveOrphanedOSTFilesOnLogoff": {
            "status": 2,
            "value": 0
          },
          "ODFCVHDReAttachRetryCount": {
            "status": 2,
            "value": 60
          },
          "ODFCLockedRetryInterval": {
            "status": 2,
            "value": 5
          },
          "ODFCIncludeOneNote_UWP": {
            "status": 2,
            "value": 0
          },
          "ODFCRedirectType": {
            "status": 2,
            "value": 2
          },
          "ODFCLockedRetryCount": {
            "status": 2,
            "value": 12
          },
          "IncludeOneDrive": {
            "status": 2,
            "value": 1
          },
          "ODFCRefreshUserPolicy": {
            "status": 2,
            "value": 0
          },
          "IncludeSkype": {
            "status": 2,
            "value": 1
          },
          "ODFCRoamSearch": {
            "status": 2,
            "value": 0
          },
          "OutlookCachedMode": {
            "status": 2,
            "value": 1
          },
          "ODFCIncludeSharepoint": {
            "status": 2,
            "value": 1
          },
          "IncludeOutlook": {
            "status": 2,
            "value": 1
          },
          "ODFCAttachVHDSDDL": {
            "status": 2,
            "value": "N/A"
          },
          "ODFCAccessNetworkAsComputerObject": {
            "status": 2,
            "value": 0
          },
          "ODFCPreventLoginWithTempProfile": {
            "status": 2,
            "value": 0
          },
          "ODFCPreventLoginWithFailure": {
            "status": 2,
            "value": 0
          },
          "ODFCNumSessionVHDsToKeep": {
            "status": 2,
            "value": 2
          },
          "ODFCVHDAccessMode": {
            "status": 2,
            "value": 0
          },
          "CacheEnable": {
            "status": 1,
            "value": 0
          },
          "CacheDirectory": {
            "status": 2,
            "value": "C:\\ProgramData\\FSLogix\\Cache"
          },
          "WriteCacheDirectory": {
            "status": 2,
            "value": "C:\\ProgramData\\FSLogix\\Cache"
          },
          "ProxyDirectory": {
            "status": 2,
            "value": "C:\\ProgramData\\FSLogix\\Proxy"
          },
          "ProfileCcdUnregisterTimeout": {
            "status": 2,
            "value": 0
          },
          "ProfileHealthyProvidersRequiredForRegister": {
            "status": 2,
            "value": 0
          },
          "ProfileCcdMaxCacheSizeInMBs": {
            "status": 2,
            "value": 0
          },
          "ProfileHealthyProvidersRequiredForUnregister": {
            "status": 2,
            "value": 1
          },
          "ProfileClearCacheOnLogoff": {
            "status": 2,
            "value": 0
          },
          "ProfileClearCacheOnForcedUnregister": {
            "status": 2,
            "value": 0
          },
          "ODFCClearCacheOnForcedUnregister": {
            "status": 2,
            "value": 0
          },
          "ODFCClearCacheOnLogoff": {
            "status": 2,
            "value": 0
          },
          "ODFCHealthyProvidersRequiredForRegister": {
            "status": 2,
            "value": 0
          },
          "ODFCHealthyProvidersRequiredForUnregister": {
            "status": 2,
            "value": 1
          },
          "ODFCCcdMaxCacheSizeInMBs": {
            "status": 2,
            "value": 0
          },
          "ODFCCcdUnregisterTimeout": {
            "status": 2,
            "value": 0
          },
          "LoggingEnabled": {
            "status": 1,
            "value": 2
          },
          "LogLevel": {
            "status": 1,
            "value": 1
          },
          "LogFileKeepingPeriod": {
            "status": 1,
            "value": 2
          },
          "LogLocation": {
            "status": 1,
            "value": "%%ProgramData%%FSLogixLogs" # Please watching this format handling about %% characters
          },
          "RobocopyLogPath": {
            "status": 1,
            "value": ""
          },
          "DriverInterface": {
            "status": 1,
            "value": 0
          },
          "RuleCompilation": {
            "status": 1,
            "value": 0
          },
          "ProcessStart": {
            "status": 1,
            "value": 0
          },
          "LoggingProfile": {
            "status": 1,
            "value": 1
          },
          "LoggingService": {
            "status": 1,
            "value": 0
          },
          "LoggingPrinter": {
            "status": 1,
            "value": 0
          },
          "LoggingFont": {
            "status": 1,
            "value": 0
          },
          "ConfigTool": {
            "status": 1,
            "value": 1
          },
          "LoggingODFC": {
            "status": 1,
            "value": 0
          },
          "AdsComputerGroup": {
            "status": 1,
            "value": 0
          },
          "LoggingNetwork": {
            "status": 1,
            "value": 0
          },
          "LoggingSearch": {
            "status": 1,
            "value": 1
          },
          "FrxLauncher": {
            "status": 1,
            "value": 0
          },
          "SearchPlugin": {
            "status": 1,
            "value": 0
          }
        }),
        "redir_exclude_common_folders": jsonencode({
          "status": 1,
          "value": 0
        }),
        "redir_exclude_copy0s": "[]",
        "redir_exclude_copy1s": "[]",
        "redir_exclude_copy2s": "[]",
        "redir_exclude_copy3s": "[]",
        "redir_exclude_includes": "[]"
      }
    }
  })
}
`, name, description)
}
