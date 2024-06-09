echo -en "\nDOCUMENT_LIST=" >> web/.env.production
for file in ./docs/*; do
    if [ -f "$file" ]; then
        filename=$(basename -- "${file}");
        title=$(grep "title:" "$file" | sed "s|^title: \(.*\)$|\1|" -)
        echo -n "${title}:docs/${filename%.*}.html " >> web/.env.production
    fi
done