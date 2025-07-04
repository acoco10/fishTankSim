#!/usr/bin/env bash

# Directory containing your spritesheets
INPUT_DIR="./assets/images/fishSpriteSheets/"

# Directory to store the generated normal maps


# Create the output directory if it doesn't exist
# Path to your Laigter binary
# If laigter is in your PATH, you can just do: LAIGTER=laigter
LAIGTER="laigter"


# Iterate over all PNG files in the input directory
for f in "$INPUT_DIR"/*.png; do
    # Get the base filename without extension
    filename=$(basename -- "$f")
    name="${filename%.*}"

    # Skip files ending with _n.png or _normal.png
    if [[ "$filename" == *_n.png || "$filename" == *_normal.png ]]; then
        echo "Skipping $filename (already a normal map)"
        continue
    fi

    echo "Processing $filename..."

    # Run Laigter to generate the normal map
    "$LAIGTER" --no-gui -d "$f" -n
    generated_normal="assets/images/fishSpriteSheets/assets/${filename}_n.png"
    if [[ -f "$generated_normal" ]]; then
        # Move to output dir
        mv "$generated_normal" "./my-sprites/normals/${filename}_n.png"
        echo "✅ Moved normal map for $name"
        # Optionally delete the generated assets directory
        rm -rf "assets/images/fishSpriteSheets/assets"
    else
        echo "⚠️ Normal map not found for $name"
    fi
done





echo "Done"