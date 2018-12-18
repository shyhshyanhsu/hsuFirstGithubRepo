# Reporting API

Run ReportingService (listen on 8090) and mock-pos (listen on 8091) under /bin directory

$ ReportingService



listen=:8090 caller=posproxy.go:31 proxy_to=POS_APIs
listen=:8090 caller=posproxy.go:49 proxy_to=[:8091]
listen=:8090 caller=main.go:35 msg=HTTP addr=:8090

$ mock-pos



callTime=2018-12-18T01:31:19-06:00 method=businesses input="hello, world" output="/businesses : HELLO, WORLD" err=null took=1.178712ms

To test : open a console and run curl $ curl -XPOST -d'{"s":"hello, world"}' localhost:8090/reporting

$ curl -XPOST -d'{"s":"hello, world"}' localhost:8090/reporting



{"v":"/businesses : HELLO, WORLD"}
$ 