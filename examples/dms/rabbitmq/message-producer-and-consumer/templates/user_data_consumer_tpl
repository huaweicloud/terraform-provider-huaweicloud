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

# Create consumer.py
cat > consumer.py <<'CONSUMER_EOF'
${consumer_script}
CONSUMER_EOF

chmod +x consumer.py

# Create systemd service file for consumer
cat > /etc/systemd/system/rabbitmq-consumer.service <<'SERVICEEOF'
[Unit]
Description=RabbitMQ Consumer Service
After=network.target

[Service]
Type=simple
Environment="RABBITMQ_HOST=${rabbitmq_host}"
Environment="RABBITMQ_PORT=5672"
Environment="RABBITMQ_USER=${rabbitmq_user}"
Environment="RABBITMQ_PASS=${rabbitmq_password}"
Environment="RABBITMQ_VHOST=${rabbitmq_vhost}"
Environment="QUEUE_NAME=${queue_name}"
ExecStart=/usr/bin/python3 /opt/rabbitmq-apps/consumer.py
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Reload systemd and enable/start service
systemctl daemon-reload
systemctl enable rabbitmq-consumer.service
systemctl start rabbitmq-consumer.service

# Wait a moment and check service status
sleep 5
systemctl status rabbitmq-consumer.service --no-pager || true

