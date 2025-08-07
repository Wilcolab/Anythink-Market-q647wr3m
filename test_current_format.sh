#!/bin/bash

echo "=== Testing Current API Format ==="
echo "Current format uses numeric IDs for both questions and answers:"
echo ""

echo "1. Getting available questions:"
curl -s http://localhost:8080/api/questions | jq '.[0:2]'

echo ""
echo "2. Testing quiz submission with current format:"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 1, "answerId": 1},
      {"questionId": 2, "answerId": 6},
      {"questionId": 3, "answerId": 12},
      {"questionId": 4, "answerId": 14},
      {"questionId": 5, "answerId": 18}
    ]
  }' | jq .

echo ""
echo "=== Testing Your Suggested Format ==="
echo "Your format uses string question IDs and answer text:"

echo ""
echo "3. Testing with your suggested format:"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "quizId": "123", 
    "answers": {
      "q1": "A lightweight thread managed by Go runtime",
      "q2": "type"
    }
  }' | jq .
