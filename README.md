# Reporting API

Run ReportingService (listen on 8090) and mock-pos (listen on 8091) under /bin directory

$ ReportingService



listen=:8090 caller=posproxy.go:30 proxy_to=POS_APIs
listen=:8090 caller=posproxy.go:48 proxy_to=[:8091]
listen=:8090 caller=main.go:35 msg=HTTP addr=:8090

$ mock-pos



callTime=2019-01-08T22:02:02-08:00 method=businesses Limit=500 Offset=0 BusinessID=businessID1 output="unsupported value type" err=null took=7.411µs



To test : open a console and run curl $ curl -XPOST -H "Authorization: Basic cG9zVXNlcjpwb3NQYXNzd29yZA==" -d'{"business_id":"businessID1", "start_date":"2018-11-12T03:04:05.123456789Z", "days_after_start_date":20,"report_type":"LCP"}' localhost:8090/reporting

$ curl -XPOST -H "Authorization: Basic cG9zVXNlcjpwb3NQYXNzd29yZA==" -d'{"business_id":"businessID1", "start_date":"2018-11-12T03:04:05.123456789Z", "days_after_start_date":20,"report_type":"LCP"}' localhost:8090/reporting



{"id":"businessID1","name":"World's Best Fried Chicken","hours":[11,12,13,14,15,16,17,18,19,20,21,22],"updated_at":"2019-01-08T21:57:09.262243-08:00","created_at":"1900-01-02T03:04:05.123456789Z"}


$ 