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
5. **React Frontend (Vite)**: Modern web interface for submitting URLs and searching analyzed content.  

---

## 🛠️ Technologies Used  

WebScraper and Content Analysis is built using the following technologies:  

- **Go (Golang)** → For Orchestrator, Scraper, and Query services. Provides high performance and concurrency for distributed tasks.  
- **Python** → For AI Analysis Worker. Simplifies integration with Groq LLM API.  
- **React + Vite** → Modern frontend framework for building responsive user interfaces.  
- **PostgreSQL** → Stores job metadata including job ID, URL, and status.  
- **RabbitMQ** → A reliable message broker used to distribute scraping tasks to workers.  
- **Kafka** → A streaming platform used for handling large-scale raw content pipelines.  
- **Elasticsearch** → Full-text search engine storing enriched content and enabling fast queries.  
- **Groq LLM API** → Provides AI-based summarization, categorization, and keyword extraction.  
- **Docker & Docker Compose** → For containerized deployment of databases, brokers, and services.  

---

## ✨ Features  

- **Modern Web Interface**: Clean, responsive React frontend for easy interaction
- **Job Orchestration**: Submit jobs with URLs, track their progress in PostgreSQL.  
- **Scalable Scraping**: Multiple workers can scrape URLs concurrently.  
- **AI Enrichment**: Articles are categorized, summarized into 3 bullet points, and keywords extracted.  
- **Full-Text Search**: Search enriched data stored in Elasticsearch via REST API.  
- **Real-time Updates**: Frontend displays job progress and results in real-time
- **Resilient & Modular**: Each service can be scaled independently.  

---

## ⚡ Getting Started  

There are two ways you can run **WebScraper and Content Analysis** on your machine. Choose the one which suits you the most.  

---

### 🔹 Setup Method 1: Run with Docker Compose (Recommended)

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
   
   # Frontend
   VITE_API_BASE_URL=http://localhost:8080
   ```

3. **Run with Docker Compose**:
   ```bash
   docker-compose up -d
   ```

4. **Access the application**:
   - Frontend: http://localhost:5173
   - API Services: http://localhost:8080 (Orchestrator), http://localhost:8081 (Query Service)

5. **Verify services are running**:
   ```bash
   docker-compose ps
   ```

---

### 🔹 Setup Method 2: Manual Setup

If you prefer to run services manually or for development:

#### Prerequisites
- Go 1.19+
- Python 3.8+
- Node.js 16+
- Docker and Docker Compose

#### 1. Start Infrastructure Services
```bash
# Start PostgreSQL, RabbitMQ, Kafka, and Elasticsearch
docker-compose up -d postgres rabbitmq kafka elasticsearch
```

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
cd orchestrator
go run cmd/main.go
```

**Terminal 2 - Scraper Worker**:
```bash
cd scraper-worker
go run cmd/worker/main.go
```

**Terminal 3 - AI Analysis Worker**:
```bash
cd ai-worker
pip install -r requirements.txt
python main.py
```

**Terminal 4 - Query Service**:
```bash
cd query-service
go run cmd/main.go
```

#### 4. Run Frontend Application

**Terminal 5 - React Frontend**:
```bash
cd client
npm install
npm run dev
```

The frontend will be available at http://localhost:5173

---

## 🖥️ Frontend Features

The React frontend built with Vite provides:

- **📱 Responsive Design**: Works seamlessly on desktop and mobile devices
- **🔗 URL Submission Form**: Easy interface to submit new scraping jobs
- **📊 Job Status Tracking**: Real-time updates on scraping and analysis progress
- **🔍 Advanced Search**: Powerful search interface with filters and sorting
- **📈 Results Visualization**: Clean display of categorized content with summaries
- **⚡ Fast Development**: Hot reload for efficient development with Vite

### Frontend Commands:
```bash
cd client
npm install          # Install dependencies
npm run dev         # Start development server
npm run build       # Build for production
npm run preview     # Preview production build
```

---

## 🔍 Example Usage

### Using the Web Interface:
1. Open http://localhost:5173 in your browser
2. Enter a URL in the submission form
3. Monitor the job status in real-time
4. Use the search functionality to find analyzed content
5. View detailed results with categories, summaries, and keywords

### API Endpoints (for direct access):

**Submit a scraping job**:
```bash
curl -X POST http://localhost:8080/api/jobs \
  -H "Content-Type: application/json" \
  -d '{"url": "https://example-blog.com"}'
```

**Search analyzed content**:
```bash
curl "http://localhost:8081/api/search?query=technology"
```

---

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
| `VITE_API_BASE_URL` | Frontend API base URL | `http://localhost:8080` |

---

## 🚀 Deployment

### Production Deployment
For production deployment, consider:

1. **Use managed services** for databases and message brokers
2. **Configure proper monitoring** and logging
3. **Set up reverse proxy** (nginx/traefik) for frontend and APIs
4. **Implement SSL/TLS** encryption
5. **Configure proper resource limits** in Docker Compose
6. **Build optimized frontend**: `cd client && npm run build`

### Scaling Services
```bash
# Scale scraper workers
docker-compose up -d --scale scraper-worker=3

# Scale AI workers
docker-compose up -d --scale ai-worker=2
```

---

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Open a Pull Request

### Development Workflow:
```bash
# Backend development
cd orchestrator && go run cmd/main.go

# Frontend development
cd client && npm run dev

# Run tests
cd client && npm test
```

---

## 📄 License

This project is licensed under the MIT License.

---

## 🆘 Support

If you encounter any issues:

1. **Check service logs**: `docker-compose logs [service-name]`
2. **Verify environment variables** are set correctly
3. **Ensure all dependencies** are running
4. **Frontend issues**: Check browser console for errors
5. **Backend issues**: Verify API endpoints are accessible

For additional help, open an issue on GitHub.

---

## 🌟 Quick Start for Developers

```bash
# Clone and start everything
git clone <repository-url>
cd webscraperandcontentanalysis
cp .env.example .env
# Edit .env with your API keys
docker-compose up -d

# Or for development with hot reload:
cd client && npm run dev        # Frontend (port 5173)
cd orchestrator && go run cmd/main.go    # Backend API (port 8080)
```

The application will be ready at **http://localhost:5173** 🚀