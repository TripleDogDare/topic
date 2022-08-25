package topic

import (
	"context"
	"strings"
	"testing"
)

const TestData = `2019-05-28T18:26:11-05:00 hello
2019-05-28T18:26:13-05:00 hello
2019-05-28T18:35:20-05:00 hello, world!
2019-05-29T10:18:09-05:00 time tracker
2019-05-29T10:31:33-05:00 time tracker
2019-06-03T15:32:25-04:00 client mtg
2019-06-04T12:51:37-04:00 client work
2019-06-06T14:20:54-05:00 client node issues
2019-06-10T13:04:21-05:00 other client
2019-06-10T13:36:16-05:00 manager 1x1
2019-06-10T13:36:21-05:00 timesheets
2019-06-10T15:17:10-05:00 client meeting
2019-06-10T15:17:19-05:00 client status report
2019-06-10T15:17:25-05:00 pick up car
2019-07-10T13:32:04-05:00 timesheets
2019-07-15T13:07:01-05:00 kubernetes hpa testing
2019-08-14T14:27:12-05:00 client tracing
2019-08-14T14:35:30-05:00 
2019-08-14T14:35:56-05:00 new pipeline
2019-12-02T15:47:33-06:00 install tools
2019-12-02T16:20:20-06:00 coffee
2019-12-04T12:24:38-06:00 project status
2019-12-04T14:27:07-06:00 break
2019-12-09T05:10:34-06:00 status updates
2019-12-09T05:27:08-06:00 merge pipeline pr
2019-12-09T05:27:15-06:00 test new build status
2019-12-09T06:14:39-06:00 gitlab demo
2019-12-09T13:48:02-06:00 gitlab
2019-12-09T14:30:44-06:00 opportunity meeting
2019-12-09T15:00:44-06:00 opportunity demo
2019-12-09T19:12:44-06:00 break
2019-12-10T03:12:00-06:00 other client
`

func TestReport(t *testing.T) {
	ctx := context.Background()
	w := new(strings.Builder)
	r := strings.NewReader(TestData)
	if err := generateReport(ctx, r, w); err != nil {
		t.Error(err)
	}
	t.Log(w.String())
}

func BenchmarkReport(t *testing.B) {
	ctx := context.Background()
	w := new(strings.Builder)
	r := strings.NewReader(TestData)
	t.ResetTimer()
	if err := generateReport(ctx, r, w); err != nil {
		t.Error(err)
	}
}
