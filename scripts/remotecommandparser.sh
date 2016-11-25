 go run remotecommandreceiver.go | while IFS= read -r line
  do
    eval "$line"
  done
