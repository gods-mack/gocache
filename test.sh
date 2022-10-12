

doit() {
  i="$1"
  echo $i
  curl --location --request POST 'http://localhost:4000/ram' \
--header 'Content-Type: application/json' \
--data-raw '{
    "name":"hare krinhan",
    "version":[1,3,4,5,7],
    "gend":"male"
}'
}
export -f doit
seq 12345 12550 | parallel -j100 --results ~/result/{} doit


