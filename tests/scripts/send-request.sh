#!/bin/bash

# Comprobar argumentos
if [ "$#" -ne 3 ]; then
    echo "Uso: $0 <ffmpegParams> <[files-download]> <[files-upload]>"
    exit 1
fi

FFMPEG_PARAMS="$1"
DOWNLOAD_FILES="$2"
UPLOAD_FILES="$3"

for i in $(seq 1 $N) 
do
    TIMESTAMP=$(date +%s%3N)  
    
    curl -s -X POST $K8S_URL \
        -H "Ce-createdtime: $TIMESTAMP" \
        -H 'Content-Type: application/json' \
        -H 'Ce-Type: encoder' \
        -H 'Ce-Specversion: 1.0' \
        -H 'Ce-Source: /HttpEventSource' \
        -H "X-set-response-delay-ms: $DELAY" \
        -H "Ce-Id: $i" \
        -d "{
    \"ffmpegParams\": \"$FFMPEG_PARAMS\",
    \"datamesh\": {
        \"downloadFiles\": [$DOWNLOAD_FILES],
        \"uploadFiles\": [$UPLOAD_FILES],
        \"uploadUrl\": \"$UPLOAD_SERVER_URL\",
        \"timesFile\": \"$TIMES_FILE\"
    }
    }"

done

#echo "Todas las peticiones han sido enviadas."
