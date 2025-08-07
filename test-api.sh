#!/bin/bash

# Quiz API Test Script
echo "🧪 Testing Go Quiz API with PostgreSQL"
echo "======================================"

BASE_URL="http://localhost:8080"

# Test 1: Health check
echo ""
echo "1️⃣  Testing health endpoint..."
curl -s "$BASE_URL/health" | jq .

# Test 2: Get all questions
echo ""
echo "2️⃣  Testing questions endpoint..."
curl -s "$BASE_URL/api/questions" | jq '[.[] | {id: .id, question: .question, correctAnswer: .answer}]'

# Test 3: Get single question
echo ""
echo "3️⃣  Testing single question endpoint..."
curl -s "$BASE_URL/api/questions/1" | jq '{id: .id, question: .question, options: [.options[] | .text]}'

# Test 4: Perfect score submission
echo ""
echo "4️⃣  Testing quiz submission (perfect score)..."
curl -X POST "$BASE_URL/api/quiz/submit" \
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
  }' | python3 -m json.tool

# Test 5: Failing score submission
echo ""
echo "5️⃣  Testing quiz submission (failing score)..."
curl -X POST "$BASE_URL/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 1, "answerId": 2},
      {"questionId": 2, "answerId": 5},
      {"questionId": 3, "answerId": 9},
      {"questionId": 4, "answerId": 13},
      {"questionId": 5, "answerId": 17}
    ]
  }' | jq .

# Test 6: Validation error - missing answers
echo ""
echo "6️⃣  Testing validation (missing answers)..."
curl -X POST "$BASE_URL/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 1, "answerId": 1},
      {"questionId": 2, "answerId": 6}
    ]
  }' | jq .

# Test 7: Validation error - invalid question
echo ""
echo "7️⃣  Testing validation (invalid question)..."
curl -X POST "$BASE_URL/api/quiz/submit" \
  -H "Content-Type: application/json" \
  -d '{
    "answers": [
      {"questionId": 999, "answerId": 1}
    ]
  }' | jq .

echo ""
echo "✅ All tests completed!"
echo ""
echo "🌐 Quiz submission endpoint URL:"
echo "$BASE_URL/api/quiz/submit"
