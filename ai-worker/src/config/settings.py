import os
from pathlib import Path
from dotenv import load_dotenv

def load_environment():
    """Load environment variables from .env file"""
    BASE_DIR = Path(__file__).resolve().parent.parent.parent
    dotenv_path = BASE_DIR / ".env"
    load_dotenv(dotenv_path)

def get_settings():
    """Get all application settings"""
    return {
        "kafka_bootstrap_servers": os.getenv("KAFKA_BOOTSTRAP_SERVERS"),
        "kafka_topic": os.getenv("KAFKA_RAW_CONTENT_TOPIC"),
        "groq_api_key": os.getenv("GROQ_API_KEY"),
        "elasticsearch_host": os.getenv("ELASTICSEARCH_HOST", "localhost"),
        "elasticsearch_index": os.getenv("ELASTICSEARCH_INDEX", "analyzed_content")
    }
