#!/bin/bash

create() {

    DIRECTORY='deploy'
    find "$DIRECTORY" -type f \( -name "*.yaml" -o -name "*.yml" \) | while read -r INPUT_FILE 
    do
        TEMP_FILE=$INPUT_FILE.tmp

        cp "$INPUT_FILE" "$TEMP_FILE"

        while IFS='=' read -r name value ; do

            awk -v var_value="$value" '{gsub(/\$'"$name"'/, var_value); print}' "$TEMP_FILE" > "$TEMP_FILE.tmp"

            mv "$TEMP_FILE.tmp" "$TEMP_FILE"
        done < <(env)
    done
}

clear() {
    DIRECTORY='deploy'
    find "$DIRECTORY" -type f \( -name "*.yaml.tmp" -o -name "*.yml.tmp" \) | while read -r INPUT_FILE 
    do
        rm $INPUT_FILE
    done
}



ACTION="$1"
shift 

case "$ACTION" in
    create)
        create "$@"
        ;;
    clear)
        clear "$@"
        ;;
    *)
        echo "Acción no reconocida. Uso: $0 {create|clear}"
        exit 1
        ;;
esac