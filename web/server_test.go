package web_test

import (
	"encoding/json"
	"fmt"
	"net/http"

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
			Eventually(func(g Gomega) {
				expected, err := json.Marshal(web.OperationResult{
					Status: status,
				})
				g.Ω(err).To(Succeed())

				url := fmt.Sprintf("%s/operations/%s", gURL, operationID)
				req, err := http.NewRequestWithContext(gCtx, http.MethodGet, url, nil)
				g.Ω(err).To(Succeed())

				resp, err := http.DefaultClient.Do(req)
				g.Ω(err).To(Succeed())
				defer func() {
					g.Ω(resp.Body.Close()).To(Succeed())
				}()
				g.Ω(resp).To(HaveHTTPStatus(http.StatusOK))
				g.Ω(resp).To(HaveHTTPBody(MatchJSON(expected)))
			}, suiteConfig.Timeout, "1s").Should(Succeed())
		}
	}

	BeforeAll(func() {
		By("creating service", func() {
			url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
			req, err := http.NewRequestWithContext(gCtx, http.MethodPost, url, nil)
			Ω(err).To(Succeed())

			resp, err := http.DefaultClient.Do(req)
			Ω(err).To(Succeed())
			Ω(resp).To(HaveHTTPStatus(http.StatusOK))

			mutableResponse := new(web.MutableResponse)
			Ω(json.NewDecoder(resp.Body).Decode(mutableResponse)).To(Succeed())
			defer func() {
				Ω(resp.Body.Close()).To(Succeed())
			}()
			Ω(mutableResponse.ServiceID).To(Equal(gServiceID))
			operationID = mutableResponse.OperationID
		})

		By("waiting running state", waitStatus(web.StatusRunning))
	})

	It("runs", func() {
		url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
		req, err := http.NewRequestWithContext(gCtx, http.MethodGet, url, nil)
		Ω(err).To(Succeed())

		resp, err := http.DefaultClient.Do(req)
		Ω(err).To(Succeed())
		defer func() {
			Ω(resp.Body.Close()).To(Succeed())
		}()
		Ω(resp).To(HaveHTTPStatus(http.StatusOK))
		Ω(resp).To(HaveHTTPBody(BeEquivalentTo(web.StatusRunning)))
	})

	AfterAll(func() {
		By("deleting service", func() {
			url := fmt.Sprintf("%s/services/%s", gURL, gServiceID)
			req, err := http.NewRequestWithContext(gCtx, http.MethodDelete, url, nil)
			Ω(err).To(Succeed())

			resp, err := http.DefaultClient.Do(req)
			Ω(err).To(Succeed())
			defer func() {
				Ω(resp.Body.Close()).To(Succeed())
			}()
			Ω(resp).To(HaveHTTPStatus(http.StatusOK))
		})

		By("waiting deleted state", waitStatus(web.StatusDeleted))
	})
})
