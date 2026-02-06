#!/usr/bin/env python3
"""
RabbitMQ Consumer Application

This script connects to RabbitMQ and continuously consumes messages from a queue.
It reads connection information from environment variables and supports graceful shutdown.
"""

import os
import sys
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
            channel.stop_consuming()
        except Exception as e:
            logger.error(f"Error stopping consumption: {e}")
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


def process_message(ch, method, properties, body):
    """Process a received message."""
    try:
        # Parse message body
        message = json.loads(body.decode('utf-8'))
        logger.info(f"Received message: {json.dumps(message, indent=2)}")

        # Process the message (here you can add your business logic)
        # For example: save to database, call API, etc.
        logger.info(f"Processing message ID: {message.get('message_id', 'unknown')}")

        # Acknowledge the message
        ch.basic_ack(delivery_tag=method.delivery_tag)
        logger.info("Message acknowledged successfully")

    except json.JSONDecodeError as e:
        logger.error(f"Failed to parse message as JSON: {e}")
        logger.error(f"Message body: {body}")
        # Reject the message and don't requeue it
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)
    except Exception as e:
        logger.error(f"Error processing message: {e}")
        # Reject the message and requeue it for retry
        ch.basic_nack(delivery_tag=method.delivery_tag, requeue=True)


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

        # Declare queue (idempotent operation)
        channel.queue_declare(queue=queue_name, durable=True)
        logger.info(f"Queue '{queue_name}' declared")

        # Set QoS to process one message at a time
        channel.basic_qos(prefetch_count=1)

        logger.info("Successfully connected to RabbitMQ")
        return True

    except Exception as e:
        logger.error(f"Failed to connect to RabbitMQ: {e}")
        return False


def main():
    """Main function."""
    global shutdown_flag, connection, channel

    # Register signal handlers
    signal.signal(signal.SIGTERM, signal_handler)
    signal.signal(signal.SIGINT, signal_handler)

    logger.info("Starting RabbitMQ Consumer...")

    # Connect to RabbitMQ
    if not connect_rabbitmq():
        logger.error("Failed to establish connection to RabbitMQ")
        sys.exit(1)

    queue_name = get_env_var("QUEUE_NAME")

    try:
        logger.info(f"Waiting for messages in queue '{queue_name}'. To exit press CTRL+C")

        # Start consuming messages
        channel.basic_consume(
            queue=queue_name,
            on_message_callback=process_message,
            auto_ack=False  # Manual acknowledgment
        )

        channel.start_consuming()

    except KeyboardInterrupt:
        logger.info("Received keyboard interrupt")
        shutdown_flag = True
        if channel and channel.is_open:
            channel.stop_consuming()
    except Exception as e:
        logger.error(f"Unexpected error: {e}")
        shutdown_flag = True
    finally:
        logger.info("Consumer shutting down...")
        if channel and channel.is_open:
            try:
                channel.stop_consuming()
                channel.close()
            except Exception:
                pass
        if connection and connection.is_open:
            try:
                connection.close()
            except Exception:
                pass
        logger.info("Consumer stopped")


if __name__ == "__main__":
    main()

