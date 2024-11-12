---
subcategory: "Cloud Phone (CPH)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_cph_server"
description: |-
  Manages a CPH server resource within HuaweiCloud.
---

# huaweicloud_cph_server

Manages a CPH server resource within HuaweiCloud.

## Example Usage

```HCL
  variable "name" {}
  variable "server_flavor" {}
  variable "phone_flavor" {}
  variable "image_id" {}
  variable "keypair" {}
  variable "vpc_id" {}
  variable "subnet_id" {}

  resource "huaweicloud_cph_server" "test" {
    name          = var.name
    server_flavor = var.server_flavor
    phone_flavor  = var.phone_flavor
    image_id      = var.image_id
    keypair_name  = var.keypair

    vpc_id    = var.vpc_id
    subnet_id = var.subnet_id
    eip_type  = "5_bgp"

    bandwidth {
      share_type  = "0"
      charge_mode = "1"
      size        = 300
    }

    period_unit = "month"
    period      = 1
    auto_renew  = "true"

  }
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region in which to create the resource.
  If omitted, the provider-level region will be used. Changing this parameter will create a new resource.

* `name` - (Required, String) Specifies the server name.  
  The name can contain **1** to **60** characters, only English letters, Chinese characters, digits, underscore (_) and
  hyphens (-) are allowed.

* `server_flavor` - (Required, String) Specifies the CPH server flavor.

* `phone_flavor` - (Required, String) Specifies the cloud phone flavor.
  
* `image_id` - (Required, String, ForceNew) Specifies the cloud phone image ID.
  Changing this parameter will create a new resource.

* `vpc_id` - (Required, String, ForceNew) Specifies the ID of VPC which the cloud server belongs to.
  Changing this parameter will create a new resource.

* `subnet_id` - (Required, String, ForceNew) Specifies the ID of subnet which the cloud server belongs to.
  Changing this parameter will create a new resource.

* `availability_zone` - (Optional, String, ForceNew) Specifies the name of the AZ where the cloud server is located.
  Changing this parameter will create a new resource.

* `eip_id` - (Optional, String, ForceNew) Specifies the ID of an **existing** EIP assigned to the cloud server.
  This parameter and `eip_type`, `bandwidth` are alternative.
  Changing this parameter will create a new resource.

* `eip_type` - (Optional, String, ForceNew) Specifies the type of an EIP that will be automatically assigned to the
  cloud server. Changing this parameter will create a new resource. The options are as follows:
    + **5_telcom**: China Telecom.
    + **5_union**: China Unicom.
    + **5_bgp**: Dynamic BGP.
    + **5_sbgp**: Static BGP.

* `bandwidth` - (Optional, List, ForceNew) Specifies the bandwidth of an EIP that will be automatically assigned to
  the cloud server. Changing this parameter will create a new resource.
  The [bandwidth](#cphServer_BandWidth) structure is documented below.

* `period_unit` - (Required, String, ForceNew) Specifies the charging period unit.  
  The valid values are **month** and **year**.
  Changing this parameter will create a new resource.

* `period` - (Required, Int, ForceNew) Specifies the charging period.  
  If `period_unit` is set to **month**, the value ranges from **1** to **9**.
  If `period_unit` is set to **year**, the value ranges from **1** to **3**.
  Changing this parameter will create a new resource.

* `auto_renew` - (Required, String, ForceNew) Specifies whether auto renew is enabled. The valid values are **true**
  and **false**. Defaults to **false**.  
  Changing this parameter will create a new resource.

* `keypair_name` - (Optional, String, ForceNew) Specifies the key pair name, which is used for logging in to
  the cloud phone through ADB.  
  Changing this parameter will create a new resource.

* `enterprise_project_id` - (Optional, String, ForceNew) Specifies the enterprise project ID.
  Changing this parameter will create a new resource.

* `ports` - (Optional, List, ForceNew) Specifies the application port enabled by the cloud phone.
  Changing this parameter will create a new resource.
  The [ApplicationPort](#cphServer_ApplicationPort) structure is documented below.

* `tags` - (Optional, Map) Specifies the key/value pairs to associate with the CPH server.

<a name="cphServer_BandWidth"></a>
The `bandwidth` block supports:

* `share_type` - (Required, String) Specifies the bandwidth type.  
  The options are as follows:
    + **0**: Dedicated bandwidth.
    + **1**: Shared bandwidth.

* `id` - (Optional, String) Specifies the bandwidth ID.  
 You can specify an existing shared bandwidth when assigning an EIP for a shared bandwidth.
 This parameter is mandatory when you create a shared bandwidth.

* `size` - (Optional, Int) Specifies the bandwidth (Mbit/s).  
  The valid value is range from **1** to **2,000**.  
  This parameter is mandatory for a dedicated bandwidth.

* `charge_mode` - (Optional, String) Specifies which the bandwidth used by the CPH server is billed.  
 This parameter is mandatory for a dedicated bandwidth.
 The options are as follows:
   + **0**: Billed by bandwidth.
   + **1**: Billed by traffic.

<a name="cphServer_ApplicationPort"></a>
The `ApplicationPort` block supports:

* `name` - (Required, String) Specifies the application port name, which can contain a maximum of 16 bytes.  
 The key service name cannot be **adb** or **vnc**.

* `listen_port` - (Required, Int) Specifies the port number, which ranges from **10,000** to **50,000**.

* `internet_accessible` - (Required, String) Specifies whether public network access is mapped.
  The options are as follows:
    + **true**: public network access is mapped.
    + **false**: no mapping is performed.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The resource ID.

* `order_id` - The order ID.

* `addresses` - The IP addresses of the CPH server.
  The [Address](#cphServer_Address) structure is documented below.

* `security_groups` - The list of the security groups bound to the extension NIC of the CPH server.

* `status` - The CPH server status.  
  The options are as follows:
    + **0**, **1**, **3**, and **4**: Creating.
    + **2**: Abnormal.
    + **5**: Normal.
    + **8**: Frozen.
    + **10**: Stopped.
    + **11**: Being stopped.
    + **12**: Stopping failed.
    + **13**: Starting.

<a name="cphServer_Address"></a>
The `Address` block supports:

* `server_ip` - The internal IP address of the CPH server.  

* `public_ip` - The public IP address of the CPH server.  

## Import

The CPH server can be imported using the `id`, e.g.

```bash
$ terraform import huaweicloud_cph_server.test <id>
```

Note that the imported state may not be identical to your resource definition, due to some attributes missing from the
API response, security or some other reason. The missing attributes include: `image_id`, `eip_id`, `eip_type`, `auto_renew`,
`period`, and `period_unit`. It is generally recommended running `terraform plan` after importing a resource.
You can then decide if changes should be applied to the resource, or the resource definition should be updated to
align with the resource. Also you can ignore changes as below.

```hcl
resource "huaweicloud_cph_server" "test" {
    ...

    lifecycle {
      ignore_changes = [
        image_id, eip_id, eip_type, auto_renew, period, period_unit,
      ]
    }
}
