---
subcategory: "Server Migration Service (SMS)"
---

# huaweicloud\_sms\_template

Use this data source to get information about templates in the SMS.

## Example Usage

 ```hcl
resource "huaweicloud_sms_template" "test" {
  name                = "template_test_old"
  is_template         = "true"
  region              = "ap-southeast-1"
  projectid           = "0e1dae200a00f3772f62c00030321606"
  availability_zone   = "ap-southeast-1a"
  target_server_name  = "serverName"
  flavor              = "flavor"
  volumetype          = "SAS"
  data_volume_type    = "SAS"
  target_password     = "123456"

  vpc {
    id = "ee3ef9f0-0c49-46d6-8da3-5ba364adc791"
    name = "vpc123"
    cidr = "192.168.0.0/16"
  }

  nics {
    id = "autoCreate"
    name = "nics1"
    cidr = "192.168.0.0/16"
  }

  security_groups {
    id = "autoCreate"
    name = "sg1"
  }

  publicip {
    type = "5_bgp"
    bandwidth_size = 1
  }

  disk {
    index = 1
    name = "disk1"
    disktype = "SAS"
    size = 40
  }
}

 ```

## Argument Reference


The following arguments are supported:

* `name` - (Required, String) template name.

* `is_template` - (Required, boolean) Indicates whether a template is generic. If a template is associated with a task, it is not generic.

* `region` - (Required, String) Obtaining Region Information.

* `projectid` - (Required, String) The project ID.

* `target_server_name` - (Optional, String) Name of the target serve.

* `id` - (Required, String) This is the ID of the template whose information needs to be modified.

* `availability_zone` - (Required, String) Available area, if not written will cause resources can not be displayed

* `volumetype` - (Optional, String) Disk type.

* `flavor` - (Optional, String) Disk type.

* `vpc` - (Optional, VpcObject object) VPC objectn.

* `nics` - (Optional, Array of Nics objects) Nic information. Multiple nics are supported. If the nic is created automatically, enter only one NIC and use "autoCreate" as the ID.

* `security_groups` - (Optional, Array of SgObject objects) Security group: Multiple security groups are supported. If the security group is created automatically, enter only one security group and use the "autoCreate" id.

* `publicip` - (Optional, PublicIp object) public wireless mobile IP network.

* `disk` - (Optional, Array of TemplateDisk objects) Disk information.

* `data_volume_type` - (Optional, String) Data disk Disk type.

* `target_password` - (Optional, String) Destination password.

## Attributes template.vpc
* `id` - (Required, String) Virtual private cloud ID. If it is created automatically, enter "autoCreate".

* `name` - (Required, String) Name of the virtual private cloud.

* `cidr` - (Optional, String) VPC network segment: 192.168.0.0/16 by default.

## Attributes template.nics
* `id` - (Required, String) Subnet ID, if it's created automatically, use "autoCreate".

* `name` - (Required, String) The name of the subnet.

* `cidr` - (Required, String) Subnet Gateway/mask.

* `ip` - (Optional, String) Vm IP address. If this field is not displayed, an IP address is automatically assigned.

## Attributes security_groups
* `id` - (Required, String) Security group ID.

* `name` - (Required, String) Security group name.

## Attributes template.publicip
* `type` - (Required, String) Elastic public IP address type. The default value is 5_BGP.

* `bandwidth_size` - (Required, Integer) Bandwidth size, in Mbit/s
The minimum bandwidth adjustment unit varies according to the bandwidth range.
The value must be smaller than or equal to 300Mbit/s. The default minimum unit is 1Mbit/s. The default minimum unit is 50Mbit/s. The value ranges from 300Mbit/s to 1000Mbit/s. Greater than 1000Mbit/s: the default minimum unit is 500Mbit/s.

## Attributes template.disk
* `index` - (Required, Integer) Disk number, starting from 0.

* `name` - (Required, String) The name of the disk.

* `disktype` - (Required, String) Disk type, which is the same as volumeType.

* `size` - (Optional, String) The unit of disk size is GB。

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - This is the id of the template that is returned when the template is successfully created, and is required for delete and query operations.
