from typing import Dict, Any
from services.groq_services import GroqService
from services.elasticSearch_services import ElasticsearchService
from models.schemas import ScrapedData, ElasticsearchDocument

class ContentProcessor:
    def __init__(self, groq_service: GroqService, es_service: ElasticsearchService):
        self.groq_service = groq_service
        self.es_service = es_service
    
    def process_message(self, message_data: Dict[str, Any]):
        """Process a single Kafka message"""
        try:
            # Validate and parse incoming data
            scraped_data = ScrapedData(**message_data)
            
            if scraped_data.error:
                print(f" Scraper error for TaskID {scraped_data.task_id}: {scraped_data.error}")
                return
            
            print(f" Processing TaskID {scraped_data.task_id} from {scraped_data.url}")
            
            # Analyze content with AI
            analysis_result = self.groq_service.analyze_content(scraped_data.raw_text)
            print(f" Analysis complete: {analysis_result.content_type} - {analysis_result.domain_category}")
            
            # Prepare document for Elasticsearch
            es_document = ElasticsearchDocument.create_from_analysis(scraped_data, analysis_result)
            doc_id = f"{scraped_data.job_id}_{scraped_data.task_id}"
            
            # Store in Elasticsearch
            self.es_service.index_document(es_document.dict(), doc_id)
            print(f" Stored in Elasticsearch with ID: {doc_id}")
            
        except Exception as e:
            print(f" Error processing message: {e}")
