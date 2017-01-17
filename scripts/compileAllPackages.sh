set -e
go build -i github.com/RyanBillard/lte-drone/shared
go build -i github.com/RyanBillard/lte-drone/groundmav
go build -i github.com/RyanBillard/lte-drone/groundrtp
go build -i github.com/RyanBillard/lte-drone/relaymav
go build -i github.com/RyanBillard/lte-drone/relayrtp
