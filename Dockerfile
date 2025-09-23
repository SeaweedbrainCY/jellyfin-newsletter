FROM python:3.12-slim

# Set environment variables
ENV LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8 \
    PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1

WORKDIR /app

# Install only runtime dependencies and locales in one layer, clean up immediately
RUN apt-get update && apt-get install -y --no-install-recommends \
    locales \
    gosu \
    ca-certificates \
    && echo "fr_FR.UTF-8 UTF-8" >> /etc/locale.gen \
    && echo "en_US.UTF-8 UTF-8" >> /etc/locale.gen \
    && locale-gen \
    && rm -rf /var/lib/apt/lists/* \
    && apt-get clean


# Install Python dependencies
COPY requirements.txt /app/requirements.txt
RUN pip install --no-cache-dir --upgrade pip \
 && pip install --no-cache-dir --no-compile -r /app/requirements.txt

# Create config directory
RUN mkdir -p /app/config

# Copy application code
COPY source /app/source
COPY main.py /app
COPY template /app/template
COPY assets /app/assets
COPY entrypoint.sh /app/entrypoint.sh
COPY VERSION /app
COPY config/config-example.yml /app/default/config-example.yml

RUN chmod +x /app/entrypoint.sh

ENTRYPOINT ["/app/entrypoint.sh"]