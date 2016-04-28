package util

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

//Nginx log format
//'[$time_iso8601] - $app_name - $remote_addr - $remote_user - $status - "$request" - $bytes_sent - "$http_referer" - "$http_user_agent" - "$server_name" - $upstream_addr - $http_host - $upstream_response_time - $request_time';
var _ = Describe("Utils", func() {
	Context("log message", func() {
		Context("with valid json", func() {
			var (
				validLogMessage  = `{"log":"[2016-05-06T20:04:23+00:00] - deis/deis-monitor-grafana - 10.240.0.15 - - - 200 - \"GET /foo HTTP/1.1\" - 31947 - \"http://grafana.mydomain.com/dashboard/db/deis-router\" - \"Mozilla/5.0 (Macintosh;\" - \"~^grafana\\x5C.(?<domain>.+)$\" - 10.167.252.252:80 - grafana.domain.com - 0.064 - 0.064\n","stream":"stderr","docker":{"container_id":"b495e59972c6d6030bc234a41951fc564dd33bead9b92915896e06abf8a8c156"},"kubernetes":{"namespace_name":"deis","pod_id":"825fbad6-13c5-11e6-9d94-42010a800021","pod_name":"deis-router-9a41n","container_name":"deis-router","labels":{"app":"deis-router"},"host":"some-valid-host"}}`
				messageJSON, err = ParseMessage(validLogMessage)
				parsedMessage, _ = ParseNginxLog(messageJSON["log"].(string))
			)

			It("should correctly parse into a map", func() {
				Expect(messageJSON).To(HaveLen(4))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be from deis-router", func() {
				Expect(FromContainer(messageJSON, "deis-router"), true)
			})

			It("should not be from deis-logger", func() {
				Expect(FromContainer(messageJSON, "deis-logger"), true)
			})

			It("should be able to return host", func() {
				Expect(GetHost(messageJSON)).To(Equal("some-valid-host"))
			})

			It("should parse the nginx log line", func() {
				Expect(parsedMessage).To(HaveLen(6))
			})

			It("should correctly parse app name", func() {
				Expect(parsedMessage["app"]).To(Equal("deis/deis-monitor-grafana"))
			})
		})
	})

	Context("is a plain string", func() {
		It("should return an error", func() {
			_, err := ParseMessage("invalid log message")
			Expect(err).To(HaveOccurred())
		})
	})

	Context("with invalid log structure", func() {
		It("ParseNginxLog should return an error", func() {
			invalidLogMessage := `{"log":"2016/05/04 18:37:13 INFO: Router configuration has changed in k8s.\n","stream":"stderr","docker":{"container_id":"436bec6d64b57c072b074e8704a2551a1197f5a3f438e45639a33aac3f7a615d"},"kubernetes":{"namespace_name":"deis","pod_id":"ac2ab3b8-1216-11e6-bb85-42010a800021","pod_name":"deis-router-dqgcj","container_name":"deis-router","labels":{"app":"deis-router"},"host":"gke-jchauncey-default-pool-7ae1c279-c4zl"}}`
			messageJSON, _ := ParseMessage(invalidLogMessage)
			_, err := ParseNginxLog(messageJSON["log"].(string))
			Expect(err).To(HaveOccurred())
		})
	})

	Context("malformed json", func() {
		It("should return an error", func() {
			_, err := ParseMessage(`{"log":invalid json}`)
			Expect(err).To(HaveOccurred())
		})
	})
})
