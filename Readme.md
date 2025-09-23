# WebScraper and Content Analysis  

## 📖 Project Concept  

**WebScraper and Content Analysis** is a microservice-based system designed to **scrape, analyze, and make web content searchable**.  

Users submit URLs to be scraped. The system then:  
1. **Scrapes and cleans the content** from webpages.  
2. **Analyzes the text using AI (LLM API)** to categorize articles, summarize them, and extract keywords.  
3. **Stores enriched results in Elasticsearch** for fast and powerful querying.  

This architecture is designed with **scalability**, **fault tolerance**, and **clean separation of responsibilities** in mind.  

---

## 🏗️ System Architecture  

The project is built around **event-driven microservices** connected by message queues and streaming platforms.  

1. **Orchestrator (Go)**: Accepts scraping jobs, stores metadata in PostgreSQL, splits jobs into tasks, and queues them in RabbitMQ.  
2. **Scraper Workers (Go)**: Consume tasks, fetch HTML, clean text, and publish raw content into Kafka.  
3. **AI Analysis Worker (Python)**: Consumes Kafka messages, sends text to Groq LLM API for categorization, summarization, and keyword extraction, then stores results in Elasticsearch.  
4. **Query Service (Go)**: Provides REST API endpoints to query the enriched content from Elasticsearch.  

---

## 🛠️ Technologies Used  

WebScraper and Content Analysis is built using the following technologies:  

- **Go (Golang)** → For Orchestrator, Scraper, and Query services. Provides high performance and concurrency for distributed tasks.  
- **Python** → For AI Analysis Worker. Simplifies integration with Groq LLM API.  
- **React** → Modern frontend framework for building responsive user interfaces.
- **PostgreSQL** → Stores job metadata including job ID, URL, and status.  
- **RabbitMQ** → A reliable message broker used to distribute scraping tasks to workers.  
- **Kafka** → A streaming platform used for handling large-scale raw content pipelines.  
- **Elasticsearch** → Full-text search engine storing enriched content and enabling fast queries.  
- **Groq LLM API** → Provides AI-based summarization, categorization, and keyword extraction.  
- **Docker & Docker Compose** → For containerized deployment of databases, brokers, and services.  

---

## ✨ Features  

- **Modern Web Interface**:clean, responsive React frontend for easy interaction
- **Job Orchestration**: Submit jobs with URLs, track their progress in PostgreSQL.  
- **Scalable Scraping**: Multiple workers can scrape URLs concurrently.  
- **AI Enrichment**: Articles are categorized, summarized into 3 bullet points, and keywords extracted.  
- **Full-Text Search**: Search enriched data stored in Elasticsearch via REST API.  
- **Resilient & Modular**: Each service can be scaled independently.  

---

## ⚡ Getting Started  
 


### 🔹 Setup Method : Run with Docker Compose 

This is the easiest way to get all dependencies (PostgreSQL, RabbitMQ, Kafka, Elasticsearch) and services up and running.  

1. **Clone the repository**:  
   ```bash
   git clone https://github.com/your-username/webscraperandcontentanalysis.git
   cd webscraperandcontentanalysis
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   ```
   
   Edit the `.env` file and add your configuration:
   ```env
   # Database
   DATABASE_URL=postgresql://user:password@postgres:5432/webscraper?sslmode=disable
   
   # RabbitMQ
   RABBITMQ_URL=amqp://user:password@rabbitmq:5672/
   QUEUE_NAME=scraper_jobs
   
   # Kafka
   KAFKA_BOOTSTRAP_SERVICE=kafka:9092
   KAFKA_RAW_CONTENT_TOPIC=raw_content
   
   # Groq API
   GROQ_API_KEY=your_groq_api_key_here
   
   # Elasticsearch
   ELASTICSEARCH_HOST=http://elasticsearch:9200
   ELASTICSEARCH_INDEX=web_content
   ```

3. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

---

#### 2. Set up Environment Variables
```bash
export DATABASE_URL="postgresql://user:password@localhost:5432/webscraper?sslmode=disable"
export RABBITMQ_URL="amqp://user:password@localhost:5672/"
export QUEUE_NAME="scraper_jobs"
export KAFKA_BOOTSTRAP_SERVICE="localhost:9092"
export KAFKA_RAW_CONTENT_TOPIC="raw_content"
export GROQ_API_KEY="your_groq_api_key_here"
export ELASTICSEARCH_HOST="http://localhost:9200"
export ELASTICSEARCH_INDEX="web_content"
```

#### 3. Run Backend Services in Separate Terminals

**Terminal 1 - Orchestrator**:
```bash
go run -C orchestrator cmd
```

**Terminal 2 - Scraper Worker**:
```bash
go run -C scraper-worker cmd/worker
```

**Terminal 3 - AI Analysis Worker**:
```bash
cd ai-worker
pip install -r requirements.txt
python3 main.py
```

**Terminal 4 - Query Service**:
```bash

go run -C query-service cmd/main.go
```

#### 4. Run Frontend Application

**Terminal 5 - React Frontend**:
```bash

cd client
npm install 
npm run dev
```
The frontend will be available at http://localhost:5173

## 📁 Project Structure

```
webscraperandcontentanalysis/
├── client/                 # React + Vite frontend application
│   ├── src/               # Source code
│   ├── public/            # Static assets
│   └── package.json       # Frontend dependencies
├── orchestrator/          # Go service for job management
├── scraper-worker/        # Go service for web scraping
├── ai-worker/            # Python service for AI analysis
├── query-service/        # Go service for Elasticsearch queries
├── docker-compose.yml    # Docker configuration
├── .env.example         # Environment template
└── README.md           # This file
```


---

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL connection string | - |
| `RABBITMQ_URL` | RabbitMQ connection URL | - |
| `QUEUE_NAME` | RabbitMQ queue name | `scraper_jobs` |
| `KAFKA_BOOTSTRAP_SERVICE` | Kafka broker address | - |
| `KAFKA_RAW_CONTENT_TOPIC` | Kafka topic for raw content | `raw_content` |
| `GROQ_API_KEY` | Groq API key for AI analysis | - |
| `ELASTICSEARCH_HOST` | Elasticsearch connection URL | - |
| `ELASTICSEARCH_INDEX` | Elasticsearch index name | `web_content` |

---


## 🔍 Example Usage

### Using the Web Interface:
1. Open http://localhost:5173 in your browser
2. Enter a URL in the submission form
4. Use the search functionality to find analyzed content
5. View detailed results with categories, summaries, and keywords