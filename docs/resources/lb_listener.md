---
subcategory: "Elastic Load Balance (ELB)"
---

# huaweicloud\_lb\_listener

Manages an ELB listener resource within HuaweiCloud.
This is an alternative to `huaweicloud_lb_listener_v2`

## Example Usage

```hcl
resource "huaweicloud_lb_listener" "listener_1" {
  protocol        = "HTTP"
  protocol_port   = 8080
  loadbalancer_id = "d9415786-5f1a-428b-b35f-2f1523e146d2"

  tags = {
    key = "value"
  }
}
```

## Argument Reference

The following arguments are supported:

* `region` - (Optional) The region in which to create the listener resource.
    If omitted, the provider-level region will be used.
    Changing this creates a new listener.

* `protocol` - (Required) The protocol can either be TCP, HTTP, HTTPS or TERMINATED_HTTPS.
    Changing this creates a new listener.

* `protocol_port` - (Required) The port on which to listen for client traffic.
    Changing this creates a new listener.

* `tenant_id` - (Optional) Required for admins. The UUID of the tenant who owns
    the listener.  Only administrative users can specify a tenant UUID
    other than their own. Changing this creates a new listener.

* `loadbalancer_id` - (Required) The load balancer on which to provision this
    listener. Changing this creates a new listener.

* `name` - (Optional) Human-readable name for the listener. Does not have
    to be unique.

* `default_pool_id` - (Optional) The ID of the default pool with which the
    listener is associated. Changing this creates a new listener.

* `description` - (Optional) Human-readable description for the listener.

* `connection_limit` - (Optional) The maximum number of connections allowed
    for the listener. The value ranges from -1 to 2,147,483,647.
    This parameter is reserved and has been not used.
    Only the administrator can specify the maximum number of connections.

* `http2_enable` - (Optional) Specifies whether to use HTTP/2. The default value is false.
    This parameter is valid only when the protocol is set to *TERMINATED_HTTPS*.

* `default_tls_container_ref` - (Optional) Specifies the ID of the server certificate
    used by the listener. This parameter is mandatory when protocol is set to *TERMINATED_HTTPS*.

* `sni_container_refs` - (Optional) Lists the IDs of SNI certificates (server certificates
    with a domain name) used by the listener. This parameter is valid when protocol is set to *TERMINATED_HTTPS*.

* `admin_state_up` - (Optional) The administrative state of the listener.
    A valid value is true (UP) or false (DOWN).

* `tags` - (Optional) The key/value pairs to associate with the listener.

## Attributes Reference

The following attributes are exported:

* `id` - The unique ID for the listener.
* `protocol` - See Argument Reference above.
* `protocol_port` - See Argument Reference above.
* `tenant_id` - See Argument Reference above.
* `name` - See Argument Reference above.
* `default_port_id` - See Argument Reference above.
* `description` - See Argument Reference above.
* `connection_limit` - See Argument Reference above.
* `http2_enable` - See Argument Reference above.
* `default_tls_container_ref` - See Argument Reference above.
* `sni_container_refs` - See Argument Reference above.
* `admin_state_up` - See Argument Reference above.
* `tags` - See Argument Reference above.

## Timeouts
This resource provides the following timeouts configuration options:
- `create` - Default is 10 minute.
- `update` - Default is 10 minute.
- `delete` - Default is 10 minute.
