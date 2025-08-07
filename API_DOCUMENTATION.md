# Quiz API Documentation

## Overview
The Quiz API provides endpoints for retrieving quiz questions and submitting quiz answers with detailed scoring and feedback.

## Base URL
```
http://localhost:8080
```

## Authentication
Currently, no authentication is required. User identification is optional via `userId` parameter.

## Endpoints

### 1. Get All Questions
**GET** `/api/questions`

Returns all available quiz questions with their options in a simplified format.

**Response:**
```json
[
  {
    "questionId": "q1",
    "text": "What is a goroutine in Go?",
    "options": [
      "A lightweight thread managed by Go runtime",
      "A Go function",
      "A Go package",
      "A Go variable"
    ],
    "correctAnswer": "A lightweight thread managed by Go runtime"
  },
  {
    "questionId": "q2",
    "text": "Which keyword is used to define a new type in Go?",
    "options": [
      "struct",
      "type", 
      "var",
      "func"
    ],
    "correctAnswer": "type"
  }
]
```

### 2. Get Single Question
**GET** `/api/questions/{id}`

Returns a specific question by ID in simplified format.

**Response:**
```json
{
  "questionId": "q1",
  "text": "What is a goroutine in Go?",
  "options": [
    "A lightweight thread managed by Go runtime",
    "A Go function",
    "A Go package", 
    "A Go variable"
  ],
  "correctAnswer": "A lightweight thread managed by Go runtime"
}
```

### 3. Submit Quiz
**POST** `/api/quiz/submit`

Submit quiz answers and receive detailed scoring results. **Supports two formats** - choose the one that works best for your frontend!

#### Format 1: Standard Array Format (Original)
**Request Body:**
```json
{
  "userId": "exampleUserId",     // Optional user identifier
  "answers": [
    {
      "questionId": 1,           // Question ID (required)
      "answerId": 1              // Selected option ID (required)
    },
    {
      "questionId": 2,
      "answerId": 6
    }
  ]
}
```

#### Format 2: Alternative Object Format (Text-Based)
**Request Body:**
```json
{
  "quizId": "123",               // Optional quiz identifier
  "userId": "exampleUserId",     // Optional user identifier
  "answers": {
    "q1": "A lightweight thread managed by Go runtime",  // questionId: answer text
    "q2": "type",
    "q3": "Both 2 and 3",
    "q4": "0",
    "q5": "fmt"
  }
}
```

**Response:**
```json
{
  "score": 3,                    // Number of correct answers
  "total": 5,                    // Total number of questions
  "correctAnswers": 3,           // Same as score (for compatibility)
  "percentage": 60,              // Score percentage
  "passed": true,                // Whether user passed (>=60%)
  "sessionId": "session_1754564107",  // Unique session identifier
  "results": [                   // Detailed per-question results
    {
      "questionId": 1,
      "question": "What is a goroutine in Go?",
      "selectedId": 1,           // ID of selected option
      "correctId": 1,            // ID of correct option
      "isCorrect": true,         // Whether answer was correct
      "selectedText": "A lightweight thread managed by Go runtime",
      "correctText": "A lightweight thread managed by Go runtime"
    },
    {
      "questionId": 2,
      "question": "Which keyword is used to define a new type in Go?",
      "selectedId": 5,
      "correctId": 6,
      "isCorrect": false,
      "selectedText": "struct",
      "correctText": "type"
    }
  ]
}
```

### 4. Health Check
**GET** `/health`

Returns server health status.

**Response:**
```json
{
  "status": "ok",
  "message": "server is running"
}
```

## Validation Rules

### Quiz Submission Validation:
- ✅ All questions must be answered
- ✅ Question IDs must exist in the database
- ✅ Answer IDs must be valid options for the respective questions
- ✅ No duplicate answers for the same question
- ✅ Request body must contain valid JSON

### Passing Criteria:
- **Passing Score**: 60% or higher
- **Scoring**: 1 point per correct answer
- **Total Questions**: Currently 5 questions available

## Error Responses

### Missing Answers:
```json
{
  "error": "missing answers for questions: 2, 3, 4, 5"
}
```

### Invalid Question ID:
```json
{
  "error": "invalid question ID: 999"
}
```

### Invalid Answer ID:
```json
{
  "error": "invalid answer ID 999 for question 1"
}
```

### Method Not Allowed:
```json
{
  "error": "Method not allowed"
}
```

### Invalid JSON:
```json
{
  "error": "Invalid JSON body"
}
```

## CORS Support
- ✅ **Origins**: `*` (all origins allowed)
- ✅ **Methods**: `GET, POST, PUT, DELETE, OPTIONS`
- ✅ **Headers**: `Content-Type, Authorization`
- ✅ **Preflight**: OPTIONS requests handled

## Sample cURL Commands

### Get Questions:
```bash
curl -s "http://localhost:8080/api/questions"
```

### Submit Quiz - Standard Format (Perfect Score):
```bash
curl -X POST "http://localhost:8080/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "testUser123",
    "answers": [
      {"questionId": 1, "answerId": 1},
      {"questionId": 2, "answerId": 6},
      {"questionId": 3, "answerId": 12},
      {"questionId": 4, "answerId": 14},
      {"questionId": 5, "answerId": 18}
    ]
  }'
```

### Submit Quiz - Alternative Format (Perfect Score):
```bash
curl -X POST "http://localhost:8080/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "quizId": "123",
    "userId": "testUser456", 
    "answers": {
      "q1": "A lightweight thread managed by Go runtime",
      "q2": "type",
      "q3": "Both 2 and 3", 
      "q4": "0",
      "q5": "fmt"
    }
  }'
```

### Submit Quiz - Mixed Results (Alternative Format):
```bash
curl -X POST "http://localhost:8080/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": {
      "q1": "A Go function",     # Wrong answer
      "q2": "type",              # Correct
      "q3": "var x = 5",        # Wrong answer  
      "q4": "0",                # Correct
      "q5": "fmt"               # Correct
    }
  }'
```

## Database Integration
- **Database**: PostgreSQL
- **Connection**: Automatic with environment variables
- **Migrations**: Run automatically on startup
- **Seeding**: Initial questions populated on first run

## Environment Variables
```bash
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=postgres
DB_SSL_MODE=disable
```
