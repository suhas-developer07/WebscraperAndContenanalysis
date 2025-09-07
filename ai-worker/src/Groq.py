import json

def create_analysis_prompt(raw_text):
    """
    Creates a universal, robust prompt for AI analysis of any web content.
    Handles articles, error pages, lists, and non-text content intelligently.
    """
    prompt = f"""
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
{raw_text}

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
    return prompt

def truncate_text_for_ai(raw_text, max_tokens=4000):
    """
    Pre-processes the raw text from the scraper to make it suitable for AI analysis.
    1. Extracts only the visible article text (ignoring scripts, navbars, etc.)
    2. Truncates it to a token limit safe for the AI API.
    """
    # A simple but effective way to get cleaner text: take the first N characters.
    # This assumes the most relevant content is at the beginning of the scraped text.
    # A more advanced method would use a library like `readability` or `bs4` to extract article body.
    
    # Convert max_tokens to a rough character count (approx 4 chars per token)
    max_characters = max_tokens * 4
    
    # Truncate the text to the calculated character limit
    truncated_text = raw_text[:max_characters]
    
    # Optional: Add an ellipsis to indicate the text was truncated
    if len(raw_text) > max_characters:
        truncated_text += "... [text truncated due to length]"
    
    return truncated_text

def analyze_with_groq(raw_text,client):
    """
    Analyzes raw text from any URL. Handles errors and junk data gracefully.
    """

    processed_text = truncate_text_for_ai(raw_text,max_tokens=3000)
    try:
       
        prompt = create_analysis_prompt(processed_text)
        
        # 2. Call the Groq API
        chat_completion = client.chat.completions.create(
            messages=[{"role": "user", "content": prompt}],
            model="llama-3.1-8b-instant",
            temperature=0.1,
            max_tokens=800,
            response_format={"type": "json_object"}   
        )
        
       
        response_json = json.loads(chat_completion.choices[0].message.content)
        return response_json

    except json.JSONDecodeError:
        # This catches if the AI doesn't return valid JSON
        print("ERROR: AI response was not valid JSON.")
        return {"category": "Unknown", "summary": "Analysis error: invalid JSON", "keywords": ""}
    except Exception as e:
        # This catches network errors, API errors, etc.
        print(f"ERROR: Groq API call failed: {e}")
        return {"category": "Unknown", "summary": f"Analysis error: {str(e)}", "keywords": ""}# Example usage:

