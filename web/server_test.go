package web_test

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/waffleboot/ginkgo-learn/web"
)

var _ = Describe("Service", Serial, Ordered, func() {
	var operationID uuid.UUID

	waitStatus := func(status web.Status) func() {
		return func() {
			suiteConfig, _ := GinkgoConfiguration()
			Eventually(func(g Gomega) web.Status {
				url := fmt.Sprintf("%s/operations/%s", gURL, operationID)
				req, err := http.NewRequestWithContext(gCtx, http.MethodGet, url, nil)
				g.Expect(err).To(Succeed())

				resp, err := http.DefaultClient.Do(req)
				g.Expect(err).To(Succeed())
				g.Expect(resp.StatusCode).To(Equal(http.StatusOK))

				operationResult := new(web.OperationResult)
				g.Expect(json.NewDecoder(resp.Body).Decode(operationResult)).To(Succeed())
				g.Expect(resp.Body.Close()).To(Succeed())

				return operationResult.Status
			}).WithTimeout(suiteConfig.Timeout).WithPolling(1 * time.Second).Should(Equal(status))
		}
	}

	BeforeAll(func() {
		By("creating service", func() {
			body := strings.NewReader("{}")

			url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
			req, err := http.NewRequestWithContext(gCtx, http.MethodPost, url, body)
			Expect(err).To(Succeed())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))

			mutableResponse := new(web.MutableResponse)

			Expect(json.NewDecoder(resp.Body).Decode(mutableResponse)).To(Succeed())
			Expect(mutableResponse.ServiceID).To(Equal(gServiceID))

			Expect(resp.Body.Close()).To(Succeed())

			operationID = mutableResponse.OperationID
		})

		By("waiting running state", waitStatus(web.StatusRunning))
	})

	It("runs", func() {
		url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
		req, err := http.NewRequestWithContext(gCtx, http.MethodGet, url, nil)
		Expect(err).To(Succeed())

		resp, err := http.DefaultClient.Do(req)
		Expect(err).To(Succeed())
		Expect(resp.StatusCode).To(Equal(http.StatusOK))

		data, err := io.ReadAll(resp.Body)
		Expect(err).To(Succeed())

		status := web.Status(string(data))
		Expect(status).To(Equal(web.StatusRunning))

		Expect(resp.Body.Close()).To(Succeed())
	})

	AfterAll(func() {
		By("deleting service", func() {
			url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
			req, err := http.NewRequestWithContext(gCtx, http.MethodDelete, url, nil)
			Expect(err).To(Succeed())

			resp, err := http.DefaultClient.Do(req)
			Expect(err).To(Succeed())
			Expect(resp.StatusCode).To(Equal(http.StatusOK))
		})

		By("waiting deleted state", waitStatus(web.StatusDeleted))
	})
})
