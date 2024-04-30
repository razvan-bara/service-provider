function sendReq(){

    file=$1

    curl --request POST \
    --url http://localhost:8080/file \
    --header 'Content-Type: multipart/form-data' \
    --header 'User-Agent: insomnia/9.1.0' \
    --form file=@/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/generated/input_$file.csv
}