#!/usr/bin/env python3
"""
RabbitMQ Producer Application

This script connects to RabbitMQ and periodically sends messages to a queue.
It reads connection information from environment variables and supports graceful shutdown.
"""

import os
import sys
import time
import signal
import json
import logging
from datetime import datetime

try:
    import pika
except ImportError:
    print("ERROR: pika library is not installed. Please run: pip3 install pika==1.3.2")
    sys.exit(1)

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
    datefmt='%Y-%m-%d %H:%M:%S'
)
logger = logging.getLogger(__name__)

# Global variables for graceful shutdown
shutdown_flag = False
connection = None
channel = None


def signal_handler(signum, frame):
    """Handle shutdown signals gracefully."""
    global shutdown_flag, connection, channel
    logger.info(f"Received signal {signum}, initiating graceful shutdown...")
    shutdown_flag = True
    if channel and channel.is_open:
        try:
            channel.close()
        except Exception as e:
            logger.error(f"Error closing channel: {e}")
    if connection and connection.is_open:
        try:
            connection.close()
        except Exception as e:
            logger.error(f"Error closing connection: {e}")
    sys.exit(0)


def get_env_var(name, default=None, required=True):
    """Get environment variable with optional default value."""
    value = os.environ.get(name, default)
    if required and value is None:
        logger.error(f"Required environment variable {name} is not set")
        sys.exit(1)
    return value


def connect_rabbitmq():
    """Establish connection to RabbitMQ."""
    global connection, channel

    # Get connection parameters from environment variables
    host = get_env_var("RABBITMQ_HOST")
    port = int(get_env_var("RABBITMQ_PORT", "5672"))
    user = get_env_var("RABBITMQ_USER")
    password = get_env_var("RABBITMQ_PASS")
    vhost = get_env_var("RABBITMQ_VHOST", "/")
    queue_name = get_env_var("QUEUE_NAME")
    exchange_name = get_env_var("EXCHANGE_NAME", "")
    routing_key = get_env_var("ROUTING_KEY", queue_name)

    # Create connection parameters
    credentials = pika.PlainCredentials(user, password)
    parameters = pika.ConnectionParameters(
        host=host,
        port=port,
        virtual_host=vhost,
        credentials=credentials,
        heartbeat=600,
        blocked_connection_timeout=300
    )

    try:
        logger.info(f"Connecting to RabbitMQ at {host}:{port} (vhost: {vhost})...")
        connection = pika.BlockingConnection(parameters)
        channel = connection.channel()

        # Declare exchange if specified
        if exchange_name:
            channel.exchange_declare(
                exchange=exchange_name,
                exchange_type=get_env_var("EXCHANGE_TYPE", "direct"),
                durable=True
            )
            logger.info(f"Exchange '{exchange_name}' declared")

        # Declare queue
        channel.queue_declare(queue=queue_name, durable=True)
        logger.info(f"Queue '{queue_name}' declared")

        # Bind queue to exchange if exchange is specified
        if exchange_name:
            channel.queue_bind(
                exchange=exchange_name,
                queue=queue_name,
                routing_key=routing_key
            )
            logger.info(f"Queue '{queue_name}' bound to exchange '{exchange_name}' with routing key '{routing_key}'")

        logger.info("Successfully connected to RabbitMQ")
        return True

    except Exception as e:
        logger.error(f"Failed to connect to RabbitMQ: {e}")
        return False


def send_message(message_count):
    """Send a message to the queue."""
    global channel, shutdown_flag

    if not channel or not channel.is_open:
        logger.error("Channel is not open")
        return False

    queue_name = get_env_var("QUEUE_NAME")
    exchange_name = get_env_var("EXCHANGE_NAME", "")
    routing_key = get_env_var("ROUTING_KEY", queue_name)

    # Create message payload
    message_body = {
        "message_id": message_count,
        "timestamp": datetime.now().isoformat(),
        "data": f"Message #{message_count} from producer",
        "source": "rabbitmq-producer"
    }

    try:
        # Publish message
        if exchange_name:
            channel.basic_publish(
                exchange=exchange_name,
                routing_key=routing_key,
                body=json.dumps(message_body),
                properties=pika.BasicProperties(
                    delivery_mode=2,  # Make message persistent
                    content_type='application/json'
                )
            )
        else:
            channel.basic_publish(
                exchange='',
                routing_key=queue_name,
                body=json.dumps(message_body),
                properties=pika.BasicProperties(
                    delivery_mode=2,  # Make message persistent
                    content_type='application/json'
                )
            )

        logger.info(f"Sent message #{message_count}: {json.dumps(message_body)}")
        return True

    except Exception as e:
        logger.error(f"Failed to send message: {e}")
        return False


def main():
    """Main function."""
    global shutdown_flag, connection, channel

    # Register signal handlers
    signal.signal(signal.SIGTERM, signal_handler)
    signal.signal(signal.SIGINT, signal_handler)

    logger.info("Starting RabbitMQ Producer...")

    # Connect to RabbitMQ
    if not connect_rabbitmq():
        logger.error("Failed to establish connection to RabbitMQ")
        sys.exit(1)

    # Get message interval from environment variable
    interval = int(get_env_var("MESSAGE_INTERVAL", "5"))
    message_count = 0

    logger.info(f"Producer started. Sending messages every {interval} seconds...")

    try:
        while not shutdown_flag:
            message_count += 1
            send_message(message_count)
            time.sleep(interval)

    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt")
        shutdown_flag = True
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
    finally:
        logger.info("Producer shutting down...")
        if channel and channel.is_open:
            try:
                channel.close()
            except Exception:
                pass
        if connection and connection.is_open:
            try:
                connection.close()
            except Exception:
                pass
        logger.info("Producer stopped")


if __name__ == "__main__":
    main()

