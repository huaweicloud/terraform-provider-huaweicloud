---
layout: "huaweicloud"
page_title: "OpenStack: huaweicloud_elb_healthcheck"
sidebar_current: "docs-huaweicloud-resource-elb-healthcheck"
description: |-
  Manages an elastic loadbalancer healthcheck resource within huawei cloud.
---

# huaweicloud\_elb\_healthcheck

Manages an elastic loadbalancer healthcheck resource within huawei cloud.

## Example Usage

```hcl
resource "huaweicloud_elb_loadbalancer" "elb" {
  name = "elb"
  type = "External"
  description = "test elb"
  vpc_id = "e346dc4a-d9a6-46f4-90df-10153626076e"
  admin_state_up = 1
  bandwidth = 5
}

resource "huaweicloud_elb_listener" "listener" {
  name = "test-elb-listener"
  description = "great listener"
  protocol = "TCP"
  backend_protocol = "TCP"
  port = 12345
  backend_port = 8080
  lb_algorithm = "roundrobin"
  loadbalancer_id = "${huaweicloud_elb_loadbalancer.elb.id}"
  timeouts {
	create = "5m"
	update = "5m"
	delete = "5m"
  }
}

resource "huaweicloud_elb_healthcheck" "healthcheck" {
  listener_id = "${huaweicloud_elb_listener.listener.id}"
  healthcheck_protocol = "TCP"
  healthcheck_connect_porta = 22
  healthy_threshold = 5
  healthcheck_timeout = 25
  healthcheck_interval = 3
  timeouts {
    create = "5m"
    update = "5m"
    delete = "5m"
  }
}
```

## Argument Reference

The following arguments are supported:

* `listener_id` - (Required) Specifies the ID of the listener to which the health
    check task belongs.

* `healthcheck_protocol` - (Optional) Specifies the protocol used for the health
    check. The value can be HTTP or TCP (case-insensitive).

* `healthcheck_uri` - (Optional) Specifies the URI for health check. This parameter
    is valid when healthcheck_ protocol is HTTP. The value is a string of 1 to 80
    characters that must start with a slash (/) and can only contain letters, digits,
    and special characters, such as -/.%?#&.

* `healthcheck_connect_port` - (Optional) Specifies the port used for the health
    check. The value ranges from 1 to 65535.

* `healthy_threshold` - (Optional) Specifies the threshold at which the health
    check result is success, that is, the number of consecutive successful health
    checks when the health check result of the backend server changes from fail
    to success. The value ranges from 1 to 10.

* `unhealthy_threshold` - (Optional) Specifies the threshold at which the health
    check result is fail, that is, the number of consecutive failed health checks
    when the health check result of the backend server changes from success to fail.
    The value ranges from 1 to 10.

* `healthcheck_timeout` - (Optional) Specifies the maximum timeout duration
    (s) for the health check. The value ranges from 1 to 50.

* `healthcheck_interval` - (Optional) Specifies the maximum interval (s) for
    health check. The value ranges from 1 to 5.

## Attributes Reference

The following attributes are exported:

* `listener_id` - See Argument Reference above.
* `healthcheck_protocol` - See Argument Reference above.
* `healthcheck_uri` - See Argument Reference above.
* `healthcheck_connect_port` - See Argument Reference above.
* `healthy_threshold` - See Argument Reference above.
* `unhealthy_threshold` - See Argument Reference above.
* `healthcheck_timeout` - See Argument Reference above.
* `healthcheck_interval` - See Argument Reference above.
* `id` - Specifies the health check task ID.
* `update_time` - Specifies the time when information about the health check
    task was updated.
* `create_time` - Specifies the time when the health check task was created.