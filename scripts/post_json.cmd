curl -X POST --json  "{\"id\": \"tmp0\", \"type\": \"gauge\", \"value\": 10}" localhost:8080/update/

curl -X POST --json  "{\"id\": \"tmp1\", \"type\": \"counter\", \"delta\": 20}" localhost:8080/update/

curl -X POST --json  "{\"id\": \"tmp0\", \"type\": \"gauge\"}" localhost:8080/update/

curl -X POST --json  "{\"id\": \"tmp1\", \"type\": \"counter\"}" localhost:8080/update/
PAUSE