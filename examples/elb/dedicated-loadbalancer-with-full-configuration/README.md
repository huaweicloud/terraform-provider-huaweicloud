# Create a dedicated Elastic Load Balance service

This example provides best practice code for using Terraform to create a dedicated Elastic Load Balance (ELB) service
in HuaweiCloud, including elastic load balancer, listener, backend server group, and bind the backend server and enable
health check for each server.

## Prerequisites

* A HuaweiCloud account
* Terraform installed
* HuaweiCloud access key and secret key (AK/SK)

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

* `region_name` - The region where the resources are located
* `access_key` - The access key of the IAM user
* `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

* `vpc_name` - The name of the VPC
* `subnet_name` - The name of the subnet
* `loadbalancer_name` - The name of the loadbalancer
* `listener_name` - The name of the listener
* `security_group_name` - The name of the security group
* `instance_name` - The name of the ECS instance

#### Optional Variables

* `availability_zone` - The name of the availability zone to which the resources belong (default: "")
* `instance_flavor_id` - The flavor ID of the instance (default: "")
* `instance_flavor_performance_type` - The performance type of the instance flavor (default: "normal")
* `instance_flavor_cpu_core_count` - The CPU core count of the instance flavor (default: 2)
* `instance_flavor_memory_size` - The memory size of the instance flavor (default: 4)
* `instance_image_id` - The image ID of the instance (default: "")
* `instance_image_visibility` - The visibility of the instance image (default: "public")
* `instance_image_os` - The OS of the instance image (default: "Ubuntu")
* `vpc_cidr` - The CIDR block of the VPC (default: "172.16.0.0/16")
* `subnet_cidr` - The CIDR block of the subnet (default: "")
* `subnet_gateway_ip` - The gateway IP address of the subnet (default: "")
* `loadbalancer_cross_vpc_backend` - Whether to associate backend servers with the load balancer by using their IP
  addresses (default: false)
* `loadbalancer_description` - The description of the loadbalancer (default: null)
* `enterprise_project_id` - The enterprise project ID (default: null)
* `loadbalancer_tags` - The tags of the loadbalancer (default: {})
* `listener_protocol` - The protocol of the listener (default: "TCP")
* `listener_port` - The port of the listener (default: 8080)
* `listener_server_certificate` - The server certificate ID of the listener, required when the listener_protocol
  is HTTPS, TLS or QUIC (default: null)
* `listener_ca_certificate` - The CA certificate ID of the listener, only available when the listener_protocol
  is HTTPS (default: null)
* `listener_sni_certificates` - The SNI certificates of the listener, only available when the listener_protocol
  is HTTPS or TLS (default: [])
* `listener_sni_match_algo` - The SNI match algorithm of the listener (default: null)
* `listener_security_policy_id` - The security policy ID of the listener, only available when the listener_protocol
  is HTTPS (default: null)
* `listener_http2_enable` - Whether to enable HTTP/2, only available when the listener_protocol is HTTPS (default: null)
* `listener_port_ranges` - The port ranges of the listener (default: [])
* `listener_idle_timeout` - The idle timeout of the listener (default: 60)
* `listener_request_timeout` - The request timeout of the listener (default: null)
* `listener_response_timeout` - The response timeout of the listener (default: null)
* `listener_description` - The description of the listener (default: null)
* `listener_tags` - The tags of the listener (default: {owner = "terraform"})
* `listener_advanced_forwarding_enabled` - Whether to enable advanced forwarding (default: false)
* `pool_name` - The name of the pool (default: null)
* `pool_protocol` - The protocol of the pool (default: "TCP")
* `pool_method` - The load balancing method of the pool (default: "ROUND_ROBIN")
* `pool_any_port_enable` - Whether to enable any port for the pool (default: false)
* `pool_description` - The description of the pool (default: null)
* `pool_persistences` - The persistence configurations for the pool (default: [])
  + `type` - The persistence type
  + `cookie_name` - The cookie name for APP_COOKIE type
  + `timeout` - The timeout for the persistence
* `health_check_protocol` - The protocol of the health check (default: "TCP")
* `health_check_interval` - The interval of the health check (default: 20)
* `health_check_timeout` - The timeout of the health check (default: 15)
* `health_check_max_retries` - The maximum retries of the health check (default: 10)
* `health_check_port` - The port of the health check (default: null)
* `health_check_url_path` - The URL path of the health check (default: null)
* `health_check_status_code` - The status code of the health check (default: null)
* `health_check_http_method` - The HTTP method of the health check (default: null)
* `health_check_domain_name` - The domain name of the health check (default: null)
* `security_group_name` - The name of the security group
* `member_protocol_port` - The port of the member (default: 8080)
* `instance_name` - The name of the ECS instance
* `member_weight` - The weight of the member (default: 1)

## Usage

* Copy this example script to your `main.tf`.

* Create a `terraform.tfvars` file and fill in the required variables:

  ```hcl
  vpc_name            = "tf_test_vpc"
  subnet_name         = "tf_test_subnet"
  loadbalancer_name   = "tf_test_loadbalancer"
  listener_name       = "tf_test_listener"
  security_group_name = "tf_test_security_group"
  instance_name       = "tf_test_instance"
  ```

* Initialize Terraform:

  ```bash
  $ terraform init
  ```

* Review the Terraform plan:

  ```bash
  $ terraform plan
  ```

* Apply the configuration:

  ```bash
  $ terraform apply
  ```

* To clean up the resources:

  ```bash
  $ terraform destroy
  ```

## Protocol Configuration Options

This example supports multiple listener protocols. Here are the different protocol configurations:

### For UDP protocol listener

#### For UDP mode server group

```hcl
listener_protocol     = "UDP"
pool_protocol         = "UDP"
health_check_protocol = "UDP_CONNECT"
```

#### For QUIC mode server group

```hcl
listener_protocol = "UDP"
pool_protocol     = "QUIC"
pool_method       = "QUIC_CID"

pool_persistences = [
  {
    type = "SOURCE_IP"
  }
]

health_check_protocol = "UDP_CONNECT"
```

### For TLS protocol listener

```hcl
listener_protocol           = "TLS"
listener_server_certificate = "Your_server_certificate_id"
pool_protocol               = "TLS"
health_check_protocol       = "TLS"
```

### For HTTP protocol listener

```hcl
listener_protocol     = "HTTP"
pool_protocol         = "HTTP"
health_check_protocol = "HTTP"
```

### For HTTPS protocol listener

#### One-way authentication

```hcl
listener_protocol           = "HTTPS"
listener_server_certificate = "Your_server_certificate_id"
pool_protocol               = "HTTPS"
health_check_protocol       = "HTTPS"
health_check_url_path       = "Your_health_check_url_path"
```

#### Two-way authentication

```hcl
listener_protocol           = "HTTPS"
listener_server_certificate = "Your_server_certificate_id"
listener_ca_certificate     = "Your_ca_certificate_id"
listener_sni_certificates   = "Your_sni_certificate_ids"
pool_protocol               = "HTTPS"
health_check_protocol       = "HTTPS"
health_check_url_path       = "Your_health_check_url_path"
```

### For QUIC protocol listener

```hcl
listener_protocol           = "QUIC"
pool_protocol               = "QUIC"
health_check_protocol       = "QUIC"
listener_server_certificate = "Your_server_certificate_id"
```

## Note

* Make sure to keep your credentials secure and never commit them to version control
* All resources will be created in the specified region
* Health checks ensure backend server availability
* When backend server groups enable health checks, backend server security group rules must be configured to allow ELB
  health check protocol and port
* If ELB does not enable "IP type backend", Layer 4 traffic is not restricted by security group/ACL
 rules, and ACL rules only take effect when ELB and backend servers are in the same subnet
* If health checks use UDP protocol, security group rules must also allow ICMP protocol, otherwise health checks
  cannot be performed on backend servers

## Requirements

| Name | Version |
| ---- | ---- |
| terraform | >= 1.3.0 |
| huaweicloud | >= 1.64.3 |
