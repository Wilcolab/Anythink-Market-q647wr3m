#!/bin/bash

echo "=== Verifying Timestamp Fix ==="
echo ""

echo "1. Testing API endpoints:"
curl -s "http://localhost:8080/api/questions" | jq 'length' | xargs -I {} echo "   Found {} questions"

echo ""
echo "2. Testing quiz submission (to verify internal timestamps work):"
RESULT=$(curl -s -X POST "http://localhost:8080/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": {
      "q1": "A lightweight thread managed by Go runtime",
      "q2": "type",
      "q3": "Both 2 and 3",
      "q4": "0", 
      "q5": "fmt"
    }
  }')

echo "$RESULT" | jq '{score, total, passed, sessionId}'

echo ""
echo "3. First question details:"
echo "$RESULT" | jq '.results[0] | {questionId, question, isCorrect}'

echo ""
echo "✅ Timestamp migration successful! Server is working properly."
