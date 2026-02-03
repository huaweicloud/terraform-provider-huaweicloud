# Create a Message Producer and Consumer Best Practice for RabbitMQ

This example provides best practice code for using Terraform to create a complete message queue application scenario
using RabbitMQ and ECS instances in Huaweicloud.

## Prerequisites

- A HuaweiCloud account
- Terraform installed (>= 0.14.0)
- HuaweiCloud access key and secret key (AK/SK)
- (Optional) SSH key pair for accessing ECS instances

## Architecture

```
┌─────────────────┐
│  Producer ECS   │
│  (producer.py)  │
└────────┬────────┘
         │
         │ AMQP (5672)
         │
         ▼
┌─────────────────┐
│   RabbitMQ      │
│   Instance      │
│  ┌───────────┐  │
│  │ Exchange  │  │
│  └─────┬─────┘  │
│        │        │
│  ┌─────▼─────┐  │
│  │  Queue    │  │
│  └─────┬─────┘  │
└────────┼────────┘
         │
         │ AMQP (5672)
         │
         ▼
┌─────────────────┐
│  Consumer ECS   │
│  (consumer.py)  │
└─────────────────┘
```

### Common Commands

```bash
# Check service status
systemctl status rabbitmq-producer
systemctl status rabbitmq-consumer

# View logs
journalctl -u rabbitmq-producer -f
journalctl -u rabbitmq-consumer -f

# Restart service
systemctl restart rabbitmq-producer
systemctl restart rabbitmq-consumer

# Stop service
systemctl stop rabbitmq-producer
systemctl stop rabbitmq-consumer

# Start service
systemctl start rabbitmq-producer
systemctl start rabbitmq-consumer
```

## Variable Introduction

The following variables need to be configured:

### Authentication Variables

- `region_name` - The region where the RAM service is located
- `access_key` - The access key of the IAM user
- `secret_key` - The secret key of the IAM user

### Resource Variables

#### Required Variables

- `vpc_name` - The VPC name
- `subnet_name` - The subnet name
- `security_group_name` - The security group name
- `instance_name` - The RabbitMQ instance name
- `instance_access_user_name` - The access user of the RabbitMQ instance
- `instance_password` - The access password of the RabbitMQ instance
- `producer_instance_name` - The name of the producer ECS instance
- `consumer_instance_name` - The name of the consumer ECS instance

#### Optional Variables

- `instance_flavor_type` - The flavor type of the RabbitMQ instance (default: "cluster")
- `instance_storage_spec_code` - The storage specification code of the RabbitMQ instance
  (default: "dms.physical.storage.ultra.v2")
- `ecs_image_name` - The image name of the ECS instances (default: "Ubuntu 20.04 server 64bit")
- `vpc_cidr` - The CIDR block of the VPC (default: "192.168.0.0/16")
- `subnet_cidr` - The CIDR block of the subnet (default: "192.168.0.0/24")
- `subnet_gateway_ip` - The gateway IP of the subnet (default: "192.168.0.1")
- `security_group_name` - The name of the security group
- `security_group_rule_configurations` - The list of security group rule configurations
  (default: two ingress rules for RabbitMQ and SSH from "192.168.0.0/16")
- `instance_engine_version` - The engine version of the RabbitMQ instance (default: "3.8.35")
- `instance_broker_num` - The number of brokers of the RabbitMQ instance (default: 3)
- `instance_storage_space` - The storage space of the RabbitMQ instance in GB (default: 600)
- `instance_ssl_enable` - Whether to enable SSL for the RabbitMQ instance (default: false)
- `instance_description` - The description of the RabbitMQ instance (default: "")
- `enterprise_project_id` - The ID of the enterprise project to which the RabbitMQ instance belongs (default: null)
- `instance_tags` - The key/value pairs to associate with the RabbitMQ instance (default: {})
- `charging_mode` - The charging mode of the RabbitMQ instance (default: "postPaid")
- `period_unit` - The period unit of the RabbitMQ instance (default: null)
- `period` - The period of the RabbitMQ instance (default: null)
- `auto_renew` - The auto renew of the RabbitMQ instance (default: "false")
- `vhost_name` - The name of the RabbitMQ virtual host (default: "app_vhost")
- `exchange_name` - The name of the RabbitMQ exchange (default: "app_exchange")
- `exchange_type` - The type of the RabbitMQ exchange (default: "direct")
- `queue_name` - The name of the RabbitMQ queue (default: "app_queue")
- `eip_type` - The type of the ECS EIP (default: "5_bgp")
- `eip_share_type` - The share type of the ECS EIP (default: "PER")
- `eip_size` - The size of the ECS EIP (default: 5)
- `eip_charge_mode` - The charge mode of the ECS EIP (default: "traffic")
- `message_interval` - The interval in seconds between messages sent by the producer (default: 5)

## Usage

- Copy this example script to your working directory.

- Create a `terraform.tfvars` file and fill in the required variables:

   ```hcl
   # Network variables
   vpc_name                  = "tf_test_vpc_rabbitmq"
   subnet_name               = "tf_test_subnet_rabbitmq"
   security_group_name       = "tf_test_sg_rabbitmq"
   # RabbitMQ instance variables
   instance_name             = "tf_test_rabbitmq_instance"
   # Replace your account info
   instance_access_user_name = "admin"
   instance_password         = "YourPassword@123"
   # ECS instance variables
   producer_instance_name    = "tf_test_producer"
   consumer_instance_name    = "tf_test_consumer"
   ```

- Initialize Terraform:

   ```bash
   terraform init
   ```

- Review the Terraform plan:

   ```bash
   terraform plan
   ```

- Apply the configuration:

   ```bash
   terraform apply
   ```

- Verify the deployment:

   After the deployment completes, you can verify that the services are running:

   + **SSH to the producer ECS instance** (if keypair is configured):

     ```bash
     ssh -i ~/.ssh/your_keypair ubuntu@<producer_ip>
     ```

   + **Check producer service status:**

     ```bash
     systemctl status rabbitmq-producer
     ```

   + **View producer logs:**

     ```bash
     journalctl -u rabbitmq-producer -f
     ```

   + **SSH to the consumer ECS instance:**

     ```bash
     ssh -i ~/.ssh/your_keypair ubuntu@<consumer_ip>
     ```

   + **Check consumer service status:**

     ```bash
     systemctl status rabbitmq-consumer
     ```

   + **View consumer logs:**

     ```bash
     journalctl -u rabbitmq-consumer -f
     ```

- To clean up the resources (note: this will not undo accepted/rejected invitations):

  ```bash
  $ terraform destroy
  ```

## Notes

- Never commit `terraform.tfvars` with real credentials to version control
- Configure keypair for secure SSH access to ECS instances
- Security group rules restrict access to RabbitMQ port (5672) only from the VPC subnet
- Application uses a dedicated RabbitMQ user with limited permissions (only on the specified vhost)
- RabbitMQ instance creation takes about 20-50 minutes depending on the flavor and broker number
- ECS instances will automatically deploy and start the producer/consumer applications via user_data scripts
- Applications are configured to automatically restart on failure
- All logs are available via systemd journal (`journalctl`)

## Requirements

| Name | Version |
| ---- | ------- |
| terraform | >= 1.1.0 |
| huaweicloud | >= 1.69.0 |
