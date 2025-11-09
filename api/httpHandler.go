package api

import (
	"net/http"

	"github.com/emaforlin/ce-document-service/core"
	"github.com/emaforlin/ce-document-service/dto"
	"github.com/gin-gonic/gin"
)

type httpHandler struct {
	documentService *core.DocumentService
}

func newHTTPHandler(service *core.DocumentService) *httpHandler {
	return &httpHandler{
		documentService: service,
	}
}

func (h *httpHandler) createDocument(c *gin.Context) {
	// This will be a repeated logic, consider moving it to a middleware
	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusForbidden, httpResponseMessage{
			Message: "document owner not provided",
		})
		return
	}

	var body dto.CreateDocumentDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: " + err.Error(),
		})
		return
	}

	body.OwnerID = ownerID
	documentID, err := h.documentService.CreateNewDocument(c.Request.Context(), body)
	if err != nil {
		// Log the error for debugging purposes
		c.Error(err)
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to create document: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"document_id": documentID,
	})
}

func (h *httpHandler) getDocuments(c *gin.Context) {
	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusForbidden, httpResponseMessage{
			Message: "document owner not provided",
		})
		return
	}

	documents, err := h.documentService.GetUserDocuments(c.Request.Context(), ownerID)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to fetch documents",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"documents": documents,
	})
}

func (h *httpHandler) getOneDocument(c *gin.Context) {
	ownerID := c.GetHeader("X-User-Id")
	if ownerID == "" {
		c.JSON(http.StatusForbidden, httpResponseMessage{
			Message: "document owner not provided",
		})
		return
	}

	documentID := c.Param("id")

	foundDoc := h.documentService.GetOneDocument(c.Request.Context(), dto.GetOneDocumentDTO{
		DocumentID: documentID,
		OwnerID:    ownerID,
	})

	if foundDoc == nil {
		c.JSON(http.StatusNotFound, httpResponseMessage{
			Message: "document not found",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"document": foundDoc,
	})
}

type httpResponseMessage struct {
	Message string `json:"message"`
}
