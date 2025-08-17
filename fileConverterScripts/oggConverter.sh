#!/bin/bash
# Convert all audio files in current directory to OGG
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/../assets/soundFx"
for file in *.wav *.mp3 *.m4a *.flac; do
    if [ -f "$file" ]; then
        ffmpeg -i "$file" -c:a libvorbis -q:a 5 "oggs/${file%.*}.ogg"
    fi
done