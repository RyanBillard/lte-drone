
function delay {
  sleep 1
  go run relay/testproducer.go -a ":8080" & nc localhost 5050
}

go run relay/tcprelay.go & delay
