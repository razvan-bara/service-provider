function sendReq(){
    curl --request POST \
    --url http://localhost:8080/file \
    --header 'Content-Type: multipart/form-data' \
    --header 'User-Agent: insomnia/9.1.0' \
    --form file=@/home/razvan/Desktop/facultate/master/a1/sem2/da/service-provider/generated/input_100k.csv 
}

export -f sendReq
function sendParallelReq(){
    seq 3 | parallel -j 3 --joblog generated/my.log sendReq
    sleep 1s
    seq 5 | parallel -j 5 --joblog generated/my.log sendReq
    sleep 1s
    seq 4 | parallel -j 4 --joblog generated/my.log sendReq
}