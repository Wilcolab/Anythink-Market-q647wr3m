#!/bin/bash

echo "=== Testing Both API Formats ==="
echo ""

echo "1. Standard format (backward compatible):"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 1, "answerId": 1},
      {"questionId": 2, "answerId": 6}
    ]
  }' | jq '.score, .total, .passed'

echo ""
echo "2. Your suggested alternative format (complete quiz):"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "quizId": "123",
    "answers": {
      "q1": "A lightweight thread managed by Go runtime",
      "q2": "type",
      "q3": "Both 2 and 3",
      "q4": "0",
      "q5": "fmt"
    }
  }' | jq '.score, .total, .passed'

echo ""
echo "3. Alternative format with partial answers:"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "quizId": "123",
    "answers": {
      "q1": "A lightweight thread managed by Go runtime",
      "q2": "type"
    }
  }' | jq -r '.error // "No error"'

echo ""
echo "4. Alternative format with wrong answer:"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "answers": {
      "q1": "A Go function",
      "q2": "type",
      "q3": "Both 2 and 3",
      "q4": "0",
      "q5": "fmt"
    }
  }' | jq '.score, .total, .passed'
