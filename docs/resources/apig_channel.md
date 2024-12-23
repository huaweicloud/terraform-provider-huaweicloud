---
subcategory: "API Gateway (Dedicated APIG)"
layout: "huaweicloud"
page_title: "HuaweiCloud: huaweicloud_apig_channel"
description: ""
---

# huaweicloud_apig_channel

Manages a channel resource within HuaweiCloud.

-> After creating a channel of type server, you can configure it for an API of an HTTP/HTTPS backend service.

## Example Usage

### Create a channel of type server and use the default group to manage servers

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "backend_servers" {
  type = list(object({
    group_name = string
    id         = string
    weight     = number
  }))
}

resource "huaweicloud_apig_channel" "test" {
  instance_id = var.instance_id
  name        = var.channel_name
  port        = 8080

  dynamic "member" {
    for_each = var.backend_servers

    content {
      id     = member.value["id"]
      weight = member.value["weight"]
    }
  }
}
```

### Create a channel of type server and use the custom group to manage servers

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "backend_server_groups" {
  type = list(object({
    name        = string
    description = string
    weight      = number
  }))
}
variable "backend_servers" {
  type = list(object({
    group_name = string
    id         = string
    weight     = number
  }))
}

resource "huaweicloud_apig_channel" "test" {
  instance_id = var.instance_id
  name        = var.channel_name
  port        = 8080

  # The length of group list cannot be 0 if you want to use dynamic syntax
  dynamic "member_group" {
    for_each = var.backend_server_groups

    content {
      name        = member.value["name"]
      description = member.value["description"]
      weight      = member.value["weight"]
    }
  }

  dynamic "member" {
    for_each = var.backend_servers

    content {
      group_name = member.value["group_name"]
      id         = member.value["id"]
      weight     = member.value["weight"]
    }
  }
}
```

### Create a channel of type microservice

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "cluster_id" {}
variable "stateless_workload_name" {}
variable "member_groups_config" {
  type = list(object({
    name                 = string
    weight               = number
    microservice_port    = number
    microservice_labels  = map(string)
  }))
}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = var.instance_id
  name             = var.channel_name
  port             = 80
  balance_strategy = 1
  member_type      = "ip"
  type             = "microservice"

  dynamic "member_group" {
    for_each = var.member_groups_config

    content {
      name                 = member_group.value["name"]
      weight               = member_group.value["weight"]
      microservice_port    = member_group.value["microservice_port"]
      microservice_labels  = member_group.value["microservice_labels"]
    }
  }

  health_check {
    protocol           = "TCP"
    threshold_normal   = 2
    threshold_abnormal = 2
    interval           = 5
    timeout            = 2
    port               = 65530
    path               = "/"
    method             = "GET"
    http_codes         = "200,201,208-209"
    enable_client_ssl  = false
    status             = 1
  }

  microservice {
    cce_config {
      cluster_id    = var.cluster_id
      namespace     = "default"
      workload_type = "deployment"
      label_key     = "app"
      label_value   = var.stateless_workload_name
    }
  }
}
```

### Create a channel of type reference

```hcl
variable "instance_id" {}
variable "channel_name" {}
variable "member_group_name" {}
variable "reference_channel_id" {}

resource "huaweicloud_apig_channel" "test" {
  instance_id      = var.instance_id
  name             = var.channel_name
  port             = 82
  balance_strategy = 2
  member_type      = "ecs"
  type             = "reference"

  member_group {
    name                     = var.member_group_name
    description              = "Created by terraform script"
    weight                   = 2
    reference_vpc_channel_id = var.reference_channel_id
  }

  health_check {
    protocol           = "HTTPS"
    threshold_normal   = 2
    threshold_abnormal = 5
    interval           = 10
    timeout            = 5
    path               = "/terraform/"
    method             = "GET"
    port               = "50"
    http_codes         = "500"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional, String, ForceNew) Specifies the region where the channel is located.  
  If omitted, the provider-level region will be used. Changing this will create a new resource.

* `instance_id` - (Required, String, ForceNew) Specifies the ID of the dedicated instance to which the channel
  belongs.  
  Changing this will create a new resource.

* `name` - (Required, String) Specifies the channel name.  
  The valid length is limited from `3` to `64`, only chinese characters, english letters, digits, hyphens (-),
  underscores (_) and dots (.) are allowed.  
  The name must start with a Chinese character or English letter.

* `port` - (Required, Int) Specifies the default port for health check in channel.  
  The valid value ranges from `1` to `65,535`.

* `balance_strategy` - (Required, Int) Specifies the distribution algorithm.  
  The valid values are as follows:
  + **1**: Weighted round robin (WRR).
  + **2**: Weighted least connections (WLC).
  + **3**: Source hashing.
  + **4**: URI hashing.

* `member_type` - (Optional, String) Specifies the member type of the channel.  
  The valid values are as follows:
  + **ip**.
  + **ecs**.

* `type` - (Optional, String) Specifies the type of the channel.  
  The valid values are as follows:
  + **builtin**: Server type.
  + **microservice**: Microservice type.
  + **reference**: Reference load balance channel type.

  Defaults to `builtin` (server type).

* `member_group` - (Optional, List) Specifies the backend (server) groups of the channel.  
  The [object](#channel_member_group) structure is documented below.

* `member` - (Optional, List) Specifies the backend servers of the channel.  
  This parameter is required and only available if the `type` is `builtin`.  
  The [object](#channel_members) structure is documented below.

* `health_check` - (Optional, List) Specifies the health configuration of cloud servers associated with the load balance
  channel for APIG regularly check.  
  The [object](#channel_health_check) structure is documented below.

* `microservice` - (Optional, List) Specifies the configuration of the microservice.  
  The [object](#channel_microservice) structure is documented below.

<a name="channel_member_group"></a>
The `member_group` block supports:

* `name` - (Required, String) Specifies the name of the member group.
  The valid length is limited from `3` to `64`, only chinese and english letters, digits, hyphens (-), underscores (_)
  and dots (.) are allowed.  
  The name must start with a Chinese or English letter.

* `description` - (Optional, String) Specifies the description of the member group.

* `weight` - (Optional, String) Specifies the weight of the current member group.

* `microservice_version` - (Optional, String) Specifies the microservice version of the backend server group.

* `microservice_port` - (Optional, Int) Specifies the microservice port of the backend server group.  
  The valid value ranges from `0` to `65,535`.

* `microservice_labels` - (Optional, Map) Specifies the microservice tags of the backend server group.

* `reference_vpc_channel_id` - (Optional, String) Specifies the ID of the reference load balance channel.
  This parameter is only available if the `type` is **reference**.

<a name="channel_members"></a>
The `member` block supports:

* `host` - (Optional, String) Specifies the IP address each backend servers.

* `id` - (Optional, String) Specifies the ECS ID for each backend servers.

  -> One of the parameter `member.host` and `member.id` must be set if `member_type` is **ecs**.
     The parameter `member.host` and `member.id` are alternative.

* `name` - (Optional, String) Specifies the name of the backend server.  
  Required if the parameter `member.id` is set.

* `weight` - (Optional, Int) Specifies the weight of current backend server.  
  The valid value ranges from `0` to `10,000`, defaults to `0`.

* `is_backup` - (Optional, Bool) Specifies whether this member is the backup member.  
  Defaults to **false**.

* `group_name` - (Optional, String) Specifies the IP address each backend servers.
  If omitted, means that all backend servers are both in one group.

* `status` - (Optional, Int) Specifies the status of the backend server.  
  The valid values are as follows:
  + **1**: Normal.
  + **2**: Abnormal.

  Defaults to **1** (normal).

* `port` - (Optional, Int) Specifies the port of the backend server.  
  The valid value ranges from `0` to `65,535`.
  If omitted, the default port of channel will be used.

<a name="channel_health_check"></a>
The `health_check` block supports:

* `protocol` - (Required, String) Specifies the microservice for performing health check on backend servers.  
  The valid values are **TCP**, **HTTP** and **HTTPS**, defaults to **TCP**.

* `threshold_normal` - (Required, Int) Specifies the the healthy threshold, which refers to the number of consecutive
  successful checks required for a backend server to be considered healthy.  
  The valid value ranges from `1` to `10`.

* `threshold_abnormal` - (Required, Int) Specifies the unhealthy threshold, which refers to the number of consecutive
  failed checks required for a backend server to be considered unhealthy.  
  The valid value ranges from `1` to `10`.

* `interval` - (Required, Int) Specifies the interval between consecutive checks, in second.  
  The valid value ranges from `1` to `300`.

* `timeout` - (Required, Int) Specifies the timeout for determining whether a health check fails, in second.  
  The value must be less than the value of the time `interval`.
  The valid value ranges from `1` to `30`.

* `path` - (Optional, String) Specifies the destination path for health checks.  
  Required if the `protocol` is **HTTP** or **HTTPS**.

* `method` - (Optional, String) Specifies the request method for health check.  
  The valid values are **GET** and **HEAD**.

* `port` - (Optional, Int) Specifies the destination host port for health check.  
  The valid value ranges from `0` to `65,535`.

* `http_codes` - (Optional, String) Specifies the response codes for determining a successful HTTP response.  
  The valid value ranges from `100` to `599` and the valid formats are as follows:
  + The multiple values, for example, **200,201,202**.
  + The range, for example, **200-299**.
  + Both multiple values and ranges, for example, **201,202,210-299**.

* `enable_client_ssl` - (Optional, Bool) Specifies whether to enable two-way authentication.  
  Defaults to **false**.

* `status` - (Optional, Int) Specifies the status of health check.  
  The valid values are as follows:
  + **1**: Normal.
  + **2**: Abnormal.

  Defaults to `1` (normal).

<a name="channel_microservice"></a>
The `microservice` block supports:

* `cce_config` - (Optional, List) Specifies the CCE microservice details.  
  The [object](#microservice_cce_config) structure is documented below.

<a name="microservice_cce_config"></a>
The `cce_config` block supports:

* `cluster_id` - (Required, String) Specifies the CCE cluster ID.

* `namespace` - (Required, String) Specifies the namespace, such as the default namespace for CCE cluster: **default**.

* `workload_type` - (Required, String) Specifies the workload type.
  + **deployment**: Stateless load.
  + **statefulset**: Stateful load.
  + **daemonset**: Daemons set.

* `label_key` - (Required, String) Specifies the service label key.

* `label_value` - (Required, String) Specifies the service label value.

## Attribute Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The ID of the channel.

* `created_at` - The time when the channel was created.

* `status` - The current status of the channel.
  + **1**: Normal.
  + **2**: Abnormal.

## Import

Channels can be imported using their `id` and the ID of the related dedicated instance, separated by a slash, e.g.

```bash
$ terraform import huaweicloud_apig_channel.test <instance_id>/<id>
```
