from elasticsearch import Elasticsearch
from typing import Dict, Any



class ElasticsearchService:
    def __init__(self, host: str, index_name: str):
        self.host = host
        self.index_name = index_name
        self.es = Elasticsearch(
            ["http://localhost:9200"],
            request_timeout=30,
            verify_certs=False
        )
        self.connect()
        self.ensure_index_exists()
    
    def connect(self):
        """Connect to Elasticsearch"""
        try:
              # For Elasticsearch 8.x with security disabled
            self.es = Elasticsearch(
            [self.host],
            request_timeout=30,
            verify_certs=False,    # Disable SSL verification
            ssl_show_warn=False,   # Disable SSL warnings
            basic_auth=('elastic', '') if 'elasticsearch' in self.host else None
            )
            if not self.es.ping():
                raise ConnectionError("Failed to connect to Elasticsearch")
            print(" Connected to Elasticsearch successfully")
        except Exception as e:
            print(f" Elasticsearch connection failed: {e}")
            raise
    
    def ensure_index_exists(self):
        """Create index if it doesn't exist"""
        if not self.es.indices.exists(index=self.index_name):
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
                       # "raw_text": {"type": "text", "index": False},
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
            self.es.indices.create(index=self.index_name, body=index_settings)
            print(f" Created index: {self.index_name}")
    
    def index_document(self, document: Dict[str, Any], doc_id: str = None):
        """Index a document in Elasticsearch"""
        try:
            response = self.es.index(
                index=self.index_name,
                id=doc_id,
                document=document
            )
            print(f" Document indexed with ID: {response['_id']}")
            return response
        except Exception as e:
            print(f" Failed to index document: {e}")
            raise
