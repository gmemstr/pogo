# Use latest golang image
FROM golang:latest

# Set working directory
WORKDIR /WhiteRabbit

# Add source to container so we can build
ADD . /WhiteRabbit

# 1. Install make & co.
# 2. Install project dependencies
# 3. Build binary and move to parent directory
# 4. Create podcast directory
# 5. Generate basic skeleton files
RUN apt update; apt install build-essential -y && \
	make install && \
	make linux && chmod +x whiterabbit && \
	ls -al && \
	mkdir podcasts && \
	touch feed.rss feed.json && echo '{}' >feed.json

EXPOSE 8000

CMD ./whiterabbit