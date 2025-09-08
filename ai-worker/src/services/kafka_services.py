from confluent_kafka import Consumer, KafkaError
from typing import Callable, Dict, Any
import json

class KafkaConsumerService:
    def __init__(self, bootstrap_servers: str, group_id: str = "ai-worker-group"):
        self.conf = {
            "bootstrap.servers": bootstrap_servers,
            "group.id": group_id,
            "auto.offset.reset": "earliest"
        }
        self.consumer = Consumer(self.conf)
    
    def subscribe(self, topics: list):
        """Subscribe to Kafka topics"""
        self.consumer.subscribe(topics)
        print(f" Subscribed to topics: {topics}")
    
    def consume_messages(self, process_callback: Callable[[Dict[str, Any]], None]):
        """Consume messages and process them with callback"""
        try:
            print(" Starting Kafka consumer...")
            while True:
                msg = self.consumer.poll(1.0)
                
                if msg is None:
                    continue
                
                if msg.error():
                    if msg.error().code() == KafkaError._PARTITION_EOF:
                        continue
                    else:
                        print(f" Kafka error: {msg.error()}")
                        break
                
                try:
                    message_data = json.loads(msg.value().decode('utf-8'))
                    process_callback(message_data)
                    
                except json.JSONDecodeError as e:
                    print(f" Failed to parse JSON message: {e}")
                except Exception as e:
                    print(f" Error processing message: {e}")
                    
        except KeyboardInterrupt:
            print(" Stopping consumer...")
        finally:
            self.consumer.close()
