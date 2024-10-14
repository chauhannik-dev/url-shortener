# url-shortener
curl http://localhost:8080/jlrtPU -i
HTTP/1.1 302 Found
content-length: 0
location: https://www.example.com

curl -L http://localhost:8080/1NIm5dudJgU

curl -i http://localhost:8080/3KX4g2v040e

curl -X POST http://localhost:8080/ \
-H "Content-Type: application/json" \
-d '{"url": "https://www.youtube.com"}'

curl -X DELETE http://localhost:8080/3KX4g2v040e -i