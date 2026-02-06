#!/bin/bash
set -e

# Update system and install Python3
export DEBIAN_FRONTEND=noninteractive
apt-get update -y
apt-get install -y python3 python3-pip

# Install pika library
pip3 install pika==1.3.2

# Create application directory
mkdir -p /opt/rabbitmq-apps
cd /opt/rabbitmq-apps

# Create producer.py
cat > producer.py <<'PRODUCER_EOF'
${producer_script}
PRODUCER_EOF

chmod +x producer.py

# Create systemd service file for producer
cat > /etc/systemd/system/rabbitmq-producer.service <<'SERVICEEOF'
[Unit]
Description=RabbitMQ Producer Service
After=network.target

[Service]
Type=simple
Environment="RABBITMQ_HOST=${rabbitmq_host}"
Environment="RABBITMQ_PORT=5672"
Environment="RABBITMQ_USER=${rabbitmq_user}"
Environment="RABBITMQ_PASS=${rabbitmq_password}"
Environment="RABBITMQ_VHOST=${rabbitmq_vhost}"
Environment="QUEUE_NAME=${queue_name}"
Environment="EXCHANGE_NAME=${exchange_name}"
Environment="EXCHANGE_TYPE=${exchange_type}"
Environment="ROUTING_KEY=${routing_key}"
Environment="MESSAGE_INTERVAL=${message_interval}"
ExecStart=/usr/bin/python3 /opt/rabbitmq-apps/producer.py
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Reload systemd and enable/start service
systemctl daemon-reload
systemctl enable rabbitmq-producer.service
systemctl start rabbitmq-producer.service

# Wait a moment and check service status
sleep 5
systemctl status rabbitmq-producer.service --no-pager || true

