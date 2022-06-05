package web_test

import (
	"bytes"
	"context"
	"encoding/binary"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
	"github.com/waffleboot/ginkgo-learn/web"

	"github.com/gin-gonic/gin"
)

var (
	gHttpServer *httptest.Server
	gStatus     web.Status
	gServiceID  uuid.UUID
	gCtx        context.Context
	gURL        string
)

var _ = SynchronizedBeforeSuite(func() []byte {
	router := gin.New()

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

	gHttpServer = httptest.NewServer(router)

	return []byte(gHttpServer.URL)
}, func(b []byte) {
	gURL = string(b)

	gServiceID = uuid.New()

	suiteConfig, _ := GinkgoConfiguration()

	ctx, cancel := context.WithTimeout(context.Background(), suiteConfig.Timeout)
	DeferCleanup(cancel)

	gCtx = ctx
})

var _ = SynchronizedAfterSuite(func() {}, func() {
	gHttpServer.Close()
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
			buf := new(bytes.Buffer)
			seed := GinkgoRandomSeed()
			for i := 0; i < 2; i++ {
				if err := binary.Write(buf, binary.LittleEndian, seed); err != nil {
					return uuid.Nil, errors.WithMessagef(err, "write ginkgo seed to buffer: seed=%d", seed)
				}
			}
			operationID, err := uuid.NewRandomFromReader(buf)
			if err != nil {
				return uuid.Nil, errors.WithMessage(err, "generate random uuid from buffer")
			}
			return operationID, nil
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
