
# huaweicloud\_mls\_instance

mls instance

## Example Usage

### Mls Instance Example

```hcl
resource "huaweicloud_mrs_cluster" "cluster1" {
  cluster_name          = "mrs-cluster-acc"
  region                = "en-OS_REGION_NAME"
  billing_type          = 12
  master_node_num       = 2
  core_node_num         = 3
  master_node_size      = "s1.4xlarge.linux.mrs"
  core_node_size        = "s1.xlarge.linux.mrs"
  available_zone_id     = "{{ availability_zone }}"
  vpc_id                = "{{ vpc_id }}"
  subnet_id             = "{{ network_id }}"
  cluster_version       = "MRS 1.3.0"
  volume_type           = "SATA"
  volume_size           = 100
  safe_mode             = 0
  cluster_type          = 0
  node_public_cert_name = "KeyPair-ci"
  cluster_admin_secret  = ""
  component_list {
    component_name = "Hadoop"
  }
  component_list {
    component_name = "Spark"
  }
  component_list {
    component_name = "Hive"
  }
  timeouts {
    create = "60m"
  }
}

resource "huaweicloud_mls_instance" "instance" {
  name    = "terraform-mls-instancei"
  version = "1.5.0"
  flavor  = "mls.c2.2xlarge.common"
  network {
    vpc_id         = "{{ vpc_id }}"
    network_id     = "{{ network_id }}"
    available_zone = "{{ availability_zone }}"
    public_ip {
      bind_type = "not_use"
    }
  }
  mrs_cluster {
    id = huaweicloud_mrs_cluster.cluster1.id
  }

  timeouts {
    create = "60m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `flavor` -
  (Required)
  Instance flavor

* `mrs_cluster` -
  (Required)
  A nested object resource Structure is documented below.

* `name` -
  (Required)
  Instance name. A tenant has a unique name of the instance of one
  type.  Value range:  An instance name must contain 4 to 64 characters
  and must start with a letter. The name is case insensitive and
  contains only letters, digits, and hyphens (-) or underscores (_),
  excluding other special characters.

* `network` -
  (Required)
  A nested object resource Structure is documented below.

* `version` -
  (Required)
  Instance version

The `mrs_cluster` block supports:

* `id` -
  (Required)
  MRS cluster ID

* `user_name` -
  (Optional)
  MRS cluster username. This parameter is mandatory only when the
  MRS cluster is in the security mode

* `user_password` -
  (Optional)
  Password of the MRS cluster user

The `network` block supports:

* `available_zone` -
  (Required)
  az

* `network_id` -
  (Required)
  ID of the subnet where the instance resides

* `public_ip` -
  (Required)
  A nested object resource Structure is documented below.

* `security_group_id` -
  (Optional)
  ID of the security group of the instance

* `vpc_id` -
  (Required)
  ID of the virtual private cloud (VPC) where the instance resides

The `public_ip` block supports:

* `bind_type` -
  (Required)
  Bind type. Possible values: auto_assign, not_use


- - -

* `agency` -
  (Optional)
  Agency name. This parameter is mandatory only when you bind an
  instance to an elastic IP address (EIP). An instance must be bound to
  an EIP to grant MLS rights to obtain a tenant's token.

## Attributes Reference

In addition to the arguments listed above, the following computed attributes are exported:

* `created` -
  Time when the instance is created. The parameter format is yyyy-mm-dd
  Thh:mm:ssZ. In the format, T indicates a time start point and Z
  specifies a UTC offset, for example, the Beijing time offset is +0800

* `current_task` -
  Instance task status. Possible values: UNFREEZING FREEZING RESTORING
  SNAPSHOTTING GROWING REBOOTING REBOOT_FAILURE RESIZE_FAILURE

* `id` -
  instance id

* `inner_endpoint` -
  URL for accessing the instance. Only machines in the same VPC and
  subnet as the instance can access the URL

* `public_endpoint` -
  URL for accessing the instance. The URL can be accessed from the
  Internet. The URL is created only after the instance is bound to an
  EIP.

* `status` -
  Instance status. Possible values: CREATING AVAILABLE FAILED CREATION
  FAILED

* `updated` -
  Time when the instance is updated. The parameter format is the same
  as the format of the created parameter

* `eip_id` -
  EIP ID. This parameter value is returned only when bindType
  is set to auto_assign

## Timeouts

This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `delete` - Default is 10 minute.

## Import

Instance can be imported using the following format:

```
$ terraform import huaweicloud_mls_instance.default {{ resource id}}
```
