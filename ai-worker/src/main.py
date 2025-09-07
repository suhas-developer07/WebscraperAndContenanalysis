
#TODO: Consume a message stream comming from Kafka from this (raw_content) topic.
#for each message call openai api to process the raw content in the message stream.
#Trigger a prompt and get the response from the openAI 
#Take it down the response and store it in the  ElasticSearch.

#_______Must be finish within Monday._________

import os
from dotenv import load_dotenv
from pathlib import Path
import json
from confluent_kafka import Consumer, KafkaError
from groq import Groq
from elasticsearch import ElasticSearch,helpers
from Groq import analyze_with_groq

BASE_DIR = Path(__file__).resolve().parent.parent
dotenv_path = BASE_DIR.parent / ".env"

load_dotenv(dotenv_path)

kafka_bootstrap_service = os.getenv("KAFKA_BOOTSTRAP_SERVICE")
kafka_topic = os.getenv("KAFKA_RAW_CONTENT_TOPIC")
GROQ_API_KEY = os.getenv("GROQ_API_KEY")
elasticSearch_host = os.getenv("ELASTICSEARCH_HOST","localhost")


#kafka initialization
consumer_conf = {
    "bootstrap.servers":kafka_bootstrap_service,
    "group.id":"ai-worker-group",
    "auto.offset.reset":"earliest"
}

consumer = Consumer(consumer_conf)
consumer.subscribe([kafka_topic])
print(f" Listening for messages on topic: {kafka_topic}")

#Groq client 
client = Groq(api_key=GROQ_API_KEY)
print("Connected to Groq api")

es = ElasticSearch(
    [f"http://{elasticSearch_host}:9200"],
    request_timeout=30
) 

try:
    while True:
        msg = consumer.poll(1.0) 
        if msg is None:
            continue
        if msg.error():
            if msg.error().code() == KafkaError._PARTITION_EOF:
                continue
            else:
                print(f" Kafka error: {msg.error()}")
                break

        raw_value = msg.value().decode("utf-8")
        try:
            data = json.loads(raw_value)
           
            job_id = data.get("job_id")
            task_id = data.get("task_id")
            url = data.get("url")
            raw_text = data.get("raw_text")
            error = data.get("error")

            print(f" JobID: {job_id}, TaskID: {task_id}, URL: {url}")

            if error and raw_text == "":
                print(f"Error:{error}")
            else:
                analysis_result = analyze_with_groq(raw_text,client)
                print(analysis_result)
        
        except json.JSONDecodeError as e:
            print("f Failed to parse JSON:{raw_value},error:{e}")
            
except KeyboardInterrupt:
    print(" Stopping consumer...")

finally:
    consumer.close()
