#!/bin/bash
echo "Testing /api/questions endpoint:"
curl -s http://localhost:8080/api/questions | jq '.[0]'

echo -e "\n\nTesting quiz submission:"
curl -s -X POST http://localhost:8080/api/quiz/submit \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 1, "answerId": 1},
      {"questionId": 2, "answerId": 6},
      {"questionId": 3, "answerId": 10},
      {"questionId": 4, "answerId": 14},
      {"questionId": 5, "answerId": 17}
    ]
  }'

echo ""
