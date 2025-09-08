
#TODO: Consume a message stream comming from Kafka from this (raw_content) topic.
#for each message call openai api to process the raw content in the message stream.
#Trigger a prompt and get the response from the openAI 
#Take it down the response and store it in the  ElasticSearch.

#_______Must be finish within Monday._________

# import os
# from dotenv import load_dotenv
# from pathlib import Path
# import json
# from confluent_kafka import Consumer, KafkaError
# from groq import Groq
# from elasticsearch import ElasticSearch,helpers
# from Groq import analyze_with_groq

# BASE_DIR = Path(__file__).resolve().parent.parent
# dotenv_path = BASE_DIR.parent / ".env"

# load_dotenv(dotenv_path)

# kafka_bootstrap_service = os.getenv("KAFKA_BOOTSTRAP_SERVICE")
# kafka_topic = os.getenv("KAFKA_RAW_CONTENT_TOPIC")
# GROQ_API_KEY = os.getenv("GROQ_API_KEY")
# elasticSearch_host = os.getenv("ELASTICSEARCH_HOST","localhost")


# #kafka initialization
# consumer_conf = {
#     "bootstrap.servers":kafka_bootstrap_service,
#     "group.id":"ai-worker-group",
#     "auto.offset.reset":"earliest"
# }

# consumer = Consumer(consumer_conf)
# consumer.subscribe([kafka_topic])
# print(f" Listening for messages on topic: {kafka_topic}")

# #Groq client 
# client = Groq(api_key=GROQ_API_KEY)
# print("Connected to Groq api")

# es = ElasticSearch(
#     [f"http://{elasticSearch_host}:9200"],
#     request_timeout=30
# ) 

# try:
#     while True:
#         msg = consumer.poll(1.0) 
#         if msg is None:
#             continue
#         if msg.error():
#             if msg.error().code() == KafkaError._PARTITION_EOF:
#                 continue
#             else:
#                 print(f" Kafka error: {msg.error()}")
#                 break

#         raw_value = msg.value().decode("utf-8")
#         try:
#             data = json.loads(raw_value)
           
#             job_id = data.get("job_id")
#             task_id = data.get("task_id")
#             url = data.get("url")
#             raw_text = data.get("raw_text")
#             error = data.get("error")

#             print(f" JobID: {job_id}, TaskID: {task_id}, URL: {url}")

#             if error and raw_text == "":
#                 print(f"Error:{error}")
#             else:
#                 analysis_result = analyze_with_groq(raw_text,client)
#                 print(analysis_result)
        
#         except json.JSONDecodeError as e:
#             print("f Failed to parse JSON:{raw_value},error:{e}")
            
# except KeyboardInterrupt:
#     print(" Stopping consumer...")

# finally:
#     consumer.close()
# Remove all src. prefixes from imports
# src/main.py - USE ABSOLUTE IMPORTS WITHOUT DOTS
from config.settings import load_environment, get_settings
from services.kafka_services import KafkaConsumerService
from services.groq_services import GroqService
from services.elasticSearch_services import ElasticsearchService
from processors.content_processor import ContentProcessor

def main():
    # Load configuration
    load_environment()
    settings = get_settings()
    
    print("🚀 Starting AI Worker Service...")
    print(f"📊 Settings: {settings['elasticsearch_host']}, {settings['kafka_topic']}")
    
    try:
        # Initialize services
        groq_service = GroqService(settings['groq_api_key'])
        es_service = ElasticsearchService(
            f"http://{settings['elasticsearch_host']}:9200",
            settings['elasticsearch_index']
        )
        processor = ContentProcessor(groq_service, es_service)
        
        # Initialize and start Kafka consumer
        kafka_consumer = KafkaConsumerService(settings['kafka_bootstrap_servers'])
        kafka_consumer.subscribe([settings['kafka_topic']])
        
        # Start consuming messages
        kafka_consumer.consume_messages(processor.process_message)
        
    except Exception as e:
        print(f"❌ Failed to start service: {e}")
        raise

if __name__ == "__main__":
    main()