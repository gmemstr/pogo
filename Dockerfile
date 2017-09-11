# Use latest golang image
FROM golang:latest

# Set working directory
WORKDIR /POGO

# Add source to container so we can build
ADD . /POGO

# 1. Install make and dependencies
# 2. Install Go dependencies
# 3. Build named Linux binary and allow execution
# 4. List directory structure (for debugging really)\
# 5. Make empty podcast direcory
# 6. Create empty feed files
RUN apt update; apt install build-essential -y && \
	make install && \
	make linux && chmod +x whiterabbit && \
	ls -al && \
	mkdir podcasts && \
	touch feed.rss feed.json && echo '{}' >feed.json

EXPOSE 8000

CMD ./pogoapp