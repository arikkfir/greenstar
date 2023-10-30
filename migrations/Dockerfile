# syntax=docker/dockerfile:1
FROM neo4j:5.13-enterprise
WORKDIR /migrations
COPY migrate.sh ./
RUN chmod +x migrate.sh
COPY *.cypher ./
ENTRYPOINT ["bash", "migrate.sh"]
