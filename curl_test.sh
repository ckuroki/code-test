curl -H "Content-Type: application/json" \
  -X POST \
  --data '{"url":"http://www.google.com","code":"google"}' \
  http://localhost:8080/urls
curl -H "Content-Type: application/json" \
  -X POST \
  --data '{"code":"example"}' \
  http://localhost:8080/urls
curl -v http://localhost:8080/google
