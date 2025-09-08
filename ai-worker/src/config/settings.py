import os
from pathlib import Path
from dotenv import load_dotenv

def load_environment():
    """Load environment variables from .env file"""
    # Go up 2 levels from this file to reach WebscraperAndContentanalysis directory
    BASE_DIR = Path(__file__).resolve().parent.parent.parent.parent
    dotenv_path = BASE_DIR / ".env"
    print(f"Looking for .env at: {dotenv_path}")  # Debug line
    if dotenv_path.exists():
        print(".env file found!")  # Debug line
        load_dotenv(dotenv_path)
    else:
        print(".env file NOT found!")  # Debug line

def get_settings():
    """Get all application settings"""
    # Debug: print all environment variables
    print("Environment variables loaded:")
    for key in ["KAFKA_BOOTSTRAP_SERVICE", "KAFKA_RAW_CONTENT_TOPIC", "GROQ_API_KEY"]:
        value = os.getenv(key)
        print(f"{key}: {'SET' if value else 'NOT SET'}")
    
    return {
        "kafka_bootstrap_servers": os.getenv("KAFKA_BOOTSTRAP_SERVICE"),
        "kafka_topic": os.getenv("KAFKA_RAW_CONTENT_TOPIC"),
        "groq_api_key": os.getenv("GROQ_API_KEY"),
        "elasticsearch_host": os.getenv("ELASTICSEARCH_HOST", "localhost"),
        "elasticsearch_index": os.getenv("ELASTICSEARCH_INDEX", "analyzed_content")
    }