package web_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/waffleboot/ginkgo-learn/web"

	"github.com/gin-gonic/gin"
)

var (
	gStatus    web.Status
	gServiceID uuid.UUID
	gCtx       context.Context
	gURL       string
)

var _ = BeforeSuite(func() {
	gServiceID = uuid.New()

	GinkgoWriter.Printf("ServiceID=%s", gServiceID)

	suiteConfig, _ := GinkgoConfiguration()

	ctx, cancel := context.WithTimeout(context.Background(), suiteConfig.Timeout)
	DeferCleanup(cancel)

	router := gin.Default()

	router.POST("/services/:serviceID", func(c *gin.Context) {
		serviceID, err := uuid.Parse(c.Param("serviceID"))
		if err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		operationID, err := createService()
		if err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.JSON(http.StatusOK, web.MutableResponse{
			ServiceID:   serviceID,
			OperationID: operationID,
		})
	})

	router.GET("/operations/:operationID", func(c *gin.Context) {
		if _, err := uuid.Parse(c.Param("operationID")); err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		c.JSON(http.StatusOK, web.OperationResult{
			Status: gStatus,
		})
	})

	router.DELETE("/services/:serviceID", func(c *gin.Context) {
		if _, err := uuid.Parse(c.Param("serviceID")); err != nil {
			c.Error(err)
			c.String(http.StatusBadRequest, err.Error())
			return
		}

		if err := deleteService(); err != nil {
			c.Error(err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		c.String(http.StatusOK, "deleting")
	})

	srv := httptest.NewServer(router)
	DeferCleanup(srv.Close)

	gURL = srv.URL

	gCtx = ctx
})

func TestWeb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Web Suite")
}

func createService() (uuid.UUID, error) {
	gStatus = web.StatusCreating

	go func() {
		select {
		case <-time.After(3 * time.Second):
			gStatus = web.StatusRunning
			return
		case <-gCtx.Done():
			return
		}
	}()

	for {
		select {
		case <-time.After(1 * time.Second):
			return uuid.New(), nil
		case <-gCtx.Done():
			return uuid.Nil, gCtx.Err()
		}
	}
}

func deleteService() error {
	gStatus = web.StatusDeleting

	go func() {
		select {
		case <-time.After(3 * time.Second):
			gStatus = web.StatusDeleted
			return
		case <-gCtx.Done():
			return
		}
	}()

	for {
		select {
		case <-time.After(1 * time.Second):
			return nil
		case <-gCtx.Done():
			return gCtx.Err()
		}
	}
}
