from datetime import datetime
from typing import List, Optional
from pydantic import BaseModel

class ScrapedData(BaseModel):
    """Schema for data received from Kafka"""
    job_id: int
    task_id: int
    url: str
    raw_text: str
    error: Optional[str] = None

class AnalysisResult(BaseModel):
    """Schema for AI analysis result"""
    content_type: str
    domain_category: str
    summary: str
    key_entities: List[str]
    sentiment_tone: str

class ElasticsearchDocument(BaseModel):
    """Schema for document to be stored in Elasticsearch"""
    job_id: int
    task_id: int
    url: str
    raw_text: str
    content_type: str
    domain_category: str
    summary: str
    key_entities: List[str]
    sentiment_tone: str
    processed_at: datetime
    ai_model: str = "llama-3.1-8b-instant"

    @classmethod
    def create_from_analysis(cls, scraped_data: ScrapedData, analysis_result: AnalysisResult):
        """Create Elasticsearch document from scraped data and analysis result"""
        return cls(
            job_id=scraped_data.job_id,
            task_id=scraped_data.task_id,
            url=scraped_data.url,
            raw_text=scraped_data.raw_text,
            content_type=analysis_result.content_type,
            domain_category=analysis_result.domain_category,
            summary=analysis_result.summary,
            key_entities=analysis_result.key_entities,
            sentiment_tone=analysis_result.sentiment_tone,
            processed_at=datetime.utcnow()
        )
