#!/bin/bash

CURRENT_DIR=`pwd`
WEB_DIR="${CURRENT_DIR}/web"
LOG_FILE="${CURRENT_DIR}/gen-service.log"

# Function gets an array of services as arguments (e.g. "wire" in "service_wire.proto")
generate_files() {
    # For each sub-folder in "web"
    for d in ${WEB_DIR}/*/; do
        echo -e "# ${d}\n"
        if [ ! -f "${d}package.json" ]; then
            echo "- npm init: package.json file was not found"
            cd "${d}"
            npm init --quiet --yes >> $LOG_FILE
        fi

        if [ ! -d "${d}node_modules/" ]; then
            echo "- npm init: node_modules dir was not found"
            cd "${d}"
            npm init --quiet --yes >> $LOG_FILE
        fi

        if [ ! -d "${d}node_modules/.bin" ]; then
            echo "- npm init: .bin dir was not found"
            cd "${d}"
            npm init --quiet --yes >> $LOG_FILE
        fi

        if [ ! -f "${d}node_modules/.bin/protoc-gen-ts" ]; then
            echo "- npm install --save-dev ts-protoc-gen: protoc-gen-ts executable was not found"
            cd "${d}"
            npm install --quiet --save-dev ts-protoc-gen@0.8.0 >> $LOG_FILE
        fi

        if [ ! -d "${d}node_modules/@types/google-protobuf" ]; then
            echo "- npm install --save @types/google-protobuf: protobuf was not found"
            cd "${d}"
            npm install --quiet --save @types/google-protobuf >> $LOG_FILE
        fi

        if [ ! -d "${d}dist" ]; then
            echo "- creating ${d}dist: dist dir was not found"
            mkdir -p "${d}dist"
            echo "" > "${d}dist/.gitkeep"
        fi

        if [ ! -d "${d}src/services" ]; then
            echo "- creating ${d}src/services: src/services dir was not found"
            mkdir -p "${d}src/services"
        fi

        if [ -f "${CURRENT_DIR}/proto/models.proto" ]; then
            echo "- generating JS services for models"
            protoc -I="${CURRENT_DIR}/proto" \
                "${CURRENT_DIR}/proto/models.proto" \
                --plugin="protoc-gen-ts=${d}node_modules/.bin/protoc-gen-ts" \
                --js_out="import_style=commonjs,binary:${d}src/services" \
                --ts_out="${d}src/services"

            # proto/models.proto goes in "models" folder
            echo "# Generating GO services for models"
            protoc -I="${CURRENT_DIR}/proto" \
                "${CURRENT_DIR}/proto/models.proto" \
                --go_out=plugins=grpc,paths=source_relative:${CURRENT_DIR}/models
        fi

        # For each proto file (the argument)
        for f in $@; do
            if [ -f "${CURRENT_DIR}/proto/service_${f}.proto" ]; then
                echo "- generating JS services for ${f}"
                protoc -I="${CURRENT_DIR}/proto" \
                    "${CURRENT_DIR}/proto/service_${f}.proto" \
                    --plugin="protoc-gen-ts=${d}node_modules/.bin/protoc-gen-ts" \
                    --js_out="import_style=commonjs,binary:${d}src/services" \
                    --ts_out="service=true:${d}src/services"

                echo "# Generating GO services for ${f}"
                if [ ! -d "${CURRENT_DIR}/services/${f}" ]; then
                    mkdir -p "${CURRENT_DIR}/services/${f}"
                fi

                protoc -I="${CURRENT_DIR}/proto" \
                    "${CURRENT_DIR}/proto/service_${f}.proto" \
                    --go_out=plugins=grpc,paths=source_relative:${CURRENT_DIR}/services/${f}
            fi
        done

        echo -e ":: done\n"
    done

    echo -e "\n::Done and dusted :-)\n"
}

# -- MAIN SCRIPT STARTS HERE --

cd ${CURRENT_DIR}

dt=$(date '+%d/%m/%Y %H:%M:%S');
echo -e "# Gen-service: ${dt}\n\n" > "${LOG_FILE}"

HELP=false
ALL=false
declare -a SERVICES=() # an array

while true; do
  case $1 in
    -h | --help )
        HELP=true;
        shift;
        ;;

    -a | --all )
        ALL=true;
        shift;
        ;;

    -s | --services )
        shift;
        SERVICES="${@}"
        ;;

    -- )
        shift;
        break
        ;;

    *  )
        break
        ;;
  esac
done

# Process all service_*.proto files in proto/ folder
if [ "${ALL}" == true ]; then
    index=0
    for f in `find ${CURRENT_DIR}/proto/ -type f -iname "service_*.proto" | awk -F "service_" '{print $2}' | awk -F ".proto" '{print $1}'`; do
        SERVICES[$index]="${f}"
        index=$(($index+1))
    done
    generate_files ${SERVICES[@]}
    exit 0
fi

# Process individual service_*.proto files in proto/ folder
if [[ ${SERVICES[@]} ]]; then
    generate_files ${SERVICES[@]}
    exit 0
fi

# Display help
echo "Generation of services. Usage:"
echo ""
echo " $ ./gen-service.sh -h|--help"
echo "   Displays this help"
echo ""
echo " $ ./gen-service.sh -a|--all"
echo "   Will process all service_*.proto files in 'proto' folder"
echo ""
echo " $ ./gen-service.sh -s|--services service1 service2 ... serviceN"
echo "   Will process individual service_*.proto files in 'proto' folder (as provided in the arguments)"
echo ""
exit 0;
