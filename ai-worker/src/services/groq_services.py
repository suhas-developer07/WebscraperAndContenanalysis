from groq import Groq
import json
from typing import Dict, Any
from models.schemas import AnalysisResult  # Remove any src. prefix

class GroqService:
    def __init__(self, api_key: str):
        self.client = Groq(api_key=api_key)
    
    def truncate_text(self, text: str, max_tokens: int = 3000) -> str:
        """Truncate text to fit within token limits"""
        max_chars = max_tokens * 4
        truncated = text[:max_chars]
        if len(text) > max_chars:
            truncated += "... [text truncated]"
        return truncated
    
    def create_prompt(self, text: str) -> str:
        """Create analysis prompt"""
        return f"""
### ROLE ###
You are a skilled content analyst. Your task is to analyze provided text from any webpage and extract structured information.

### PRIMARY TASK ###
First, determine if the text contains a substantive article, post, or meaningful content.

### DECISION TREE ###

IF the text is:
- A full article or blog post
- A meaningful news story, tutorial, or opinion piece
- A product page with descriptions
- Any other substantive content

THEN perform a DETAILED ANALYSIS with the following structure:

### DETAILED ANALYSIS INSTRUCTIONS ###
1.  **Content Type:** Categorize into: ["News", "Opinion", "Tutorial/Guide", "Product/Service", "Entertainment", "Academic", "Technical", "Sports", "Other"]
2.  **Domain Category:** Identify the primary topic: ["Technology", "Politics", "Business", "Science", "Health", "Sports", "Entertainment", "Lifestyle", "Education", "Other"]
3.  **Detailed Summary:** Provide a comprehensive 4-5 sentence summary that includes:
    - Main topic and thesis of the content
    - Key points, arguments, or findings discussed
    - Important names, entities, products, or events mentioned
    - Conclusion or significance of the content
4.  **Key Entities:** Extract specific names, companies, products, or important terms mentioned.
5.  **Sentiment Tone:** Describe the overall tone: ["Informative", "Positive", "Negative", "Neutral", "Critical", "Promotional", "Opinionated"]

ELSE IF the text is:
- Mostly navigation menus, ads, or boilerplate
- An error page (404, login required, etc.)
- A list of links without substantial content
- Too short or fragmented to analyze
- Clearly not meaningful content

THEN use the FALLBACK RESPONSE:
- Set "content_type" to "Non-Article"
- Set "domain_category" to "Unknown" 
- Set "summary" to "Content not analyzable - appears to be navigation, error page, or non-substantive content"
- Set "key_entities" to an empty array
- Set "sentiment_tone" to "Neutral"

### TEXT TO ANALYZE ###
{text}

### CRITICAL RESPONSE RULES ###
- Always respond with valid JSON using the exact structure below
- Be specific and detailed in summaries for analyzable content
- For lists, identify the main topic of the list items
- For product pages, focus on the product description and features
- For error pages, be honest about the lack of analyzable content

### RESPONSE FORMAT ###
{{
    "content_type": "News | Opinion | Tutorial/Guide | Product/Service | Entertainment | Academic | Technical | Sports | Other | Non-Article",
    "domain_category": "Technology | Politics | Business | Science | Health | Sports | Entertainment | Lifestyle | Education | Other | Unknown",
    "summary": "A comprehensive 4-5 sentence summary for substantive content, or the fallback text for non-articles.",
    "key_entities": ["Entity 1", "Entity 2", "Entity 3", "Entity 4", "Entity 5"],
    "sentiment_tone": "Informative | Positive | Negative | Neutral | Critical | Promotional | Opinionated"
}}
"""
    
    def analyze_content(self, text: str) -> AnalysisResult:
        """Analyze content using Groq API"""
        try:
            processed_text = self.truncate_text(text)
            prompt = self.create_prompt(processed_text)
            
            response = self.client.chat.completions.create(
                messages=[{"role": "user", "content": prompt}],
                model="llama-3.1-8b-instant",
                temperature=0.1,
                max_tokens=800,
                response_format={"type": "json_object"}
            )
            
            result_data = json.loads(response.choices[0].message.content)
            return AnalysisResult(**result_data)
            
        except Exception as e:
            print(f" Groq analysis failed: {e}")
            # Return fallback result
            return AnalysisResult(
                content_type="Non-Article",
                domain_category="Unknown",
                summary=f"Analysis failed: {str(e)}",
                key_entities=[],
                sentiment_tone="Neutral"
            )
