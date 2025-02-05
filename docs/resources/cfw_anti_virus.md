---
subcategory: "Cloud Firewall (CFW)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cfw_anti_virus"
description: |-
  Manages a CFW anti virus resource within HuaweiCloud.
---

# huaweicloud_cfw_anti_virus

Manages a CFW anti virus resource within HuaweiCloud.

## Example Usage

```hcl
variable "object_id" {}

resource "huaweicloud_cfw_anti_virus" "test" {
  object_id = var.object_id

  scan_protocol_configs {
    protocol_type =  3
    action        =  1
  }

  scan_protocol_configs {
    protocol_type =  2
    action        =  1
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) The region in which to create the resource.
  If omitted, the provider-level region will be used.
  Changing this creates a new resource.

* `object_id` - (Required, String, NonUpdatable) Specifies the protected object ID.

* `scan_protocol_configs` - (Required, List) Specifies the scan protocol configurations.
  The [scan_protocol_configs](#ScanProtocolConfigs) structure is documented below.

<a name="ScanProtocolConfigs"></a>
The `scan_protocol_configs` block supports:

* `action` - (Required, Int) The anti virus action. The valid value can be **0** (observe) or **1** (block).

* `protocol_type` - (Required, Int) The protocol type.
  The valid values are as follows:
  + **0**: HTTP;
  + **1**: SMTP;
  + **2**: POP3;
  + **3**: IMAP4;
  + **4**: FTP;
  + **5**: SMB;
  + **6**: Malicious Access Control;

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID. The value is the same as the `object_id`.

## Import

The anti virus can be imported using `id`, e.g.

```bash
$ terraform import huaweicloud_cfw_anti_virus.test <id>
```
