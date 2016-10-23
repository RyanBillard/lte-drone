
function delay {
  sleep(1000)
  go run relay/testproducer.go -a ":8080" & go run relay/tcpconsumer.go -a ":5050"
}

go run relay/tcprelay.go & delay
