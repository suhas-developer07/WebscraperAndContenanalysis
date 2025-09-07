
def InitElasticSearch(es):
    try:
        if not es.ping():
            raise ValueError("Failed to Connect ElasticSearch! Check if it's running")
        print("Successfully connected to ElasticSearch")

        index_name = "analyzed_content"
        if not es.indices.exists(index=index_name):
            index_settings = {
            "settings": {
                "number_of_shards": 1,
                "number_of_replicas": 0
            },
            "mappings": {
                "properties": {
                    "job_id": {"type": "integer"},
                    "task_id": {"type": "integer"},
                    "url": {"type": "keyword"},
                    "raw_text": {"type": "text", "index": False},  
                    "content_type": {"type": "keyword"},
                    "domain_category": {"type": "keyword"},
                    "summary": {"type": "text"},
                    "key_entities": {"type": "keyword"},
                    "sentiment_tone": {"type": "keyword"},
                    "processed_at": {"type": "date"},
                    "ai_model": {"type": "keyword"}
                }
            }
        }
            es.indices.create(index=index_name, body=index_settings)
            print(f" Created Elasticsearch index: {index_name}")
        else:
            print(f" Elasticsearch index {index_name} already exists")
        
    except Exception as e:
        print(f" Elasticsearch initialization error: {e}")
        exit(1)


