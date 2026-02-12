package internal

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	documentService *DocumentService
}

func NewHTTPHandler(service *DocumentService) *HTTPHandler {
	return &HTTPHandler{
		documentService: service,
	}
}

func (h *HTTPHandler) deleteDocument(c *gin.Context) {
	resCode := http.StatusOK
	documentID := c.GetString("documentID")

	if err := h.documentService.DeleteDocument(c.Request.Context(), documentID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			resCode = http.StatusNotFound
		} else {
			resCode = http.StatusBadRequest
		}
		c.JSON(resCode, httpResponseMessage{
			Message: err.Error(),
		})
		return
	}
	c.JSON(resCode, gin.H{
		"message": "document deleted",
	})
}

func (h *HTTPHandler) updateDocument(c *gin.Context) {
	documentID := c.GetString("documentID")

	var body UpdateDocumentDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to marshal body",
		})
		return
	}
	body.DocumentID = documentID

	if err := h.documentService.UpdateDocumentMetadata(c.Request.Context(), body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to update document",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "document metadata updated",
	})
}

func (h *HTTPHandler) getDocumentCollaborators(c *gin.Context) {
	documentID := c.GetString("documentID")

	collaborators, err := h.documentService.getDocumentCollaborators(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, httpResponseMessage{
			Message: err.Error(),
		})
		return
	}

	// Convert to DTO response
	response := ToCollaboratorResponseList(collaborators)
	c.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) removeDocumentCollaborator(c *gin.Context) {
	var body RemoveCollaboratorDTO

	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: " + err.Error(),
		})
		return
	}

	body.DocumentID = c.GetString("documentID")

	if err := h.documentService.RemoveDocumentCollaborator(c.Request.Context(), body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: couldn't remove collaborator",
		})
		return
	}
	c.JSON(http.StatusOK, httpResponseMessage{
		Message: "collaborator removed",
	})
}

func (h *HTTPHandler) addDocumentCollaborator(c *gin.Context) {
	// Get values from middleware context
	ownerID := c.GetString("userID")
	documentID := c.GetString("documentID")

	var body AddCollaboratorDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: " + err.Error(),
		})
		return
	}

	body.OwnerID = ownerID
	body.DocumentID = documentID

	if err := h.documentService.AddCollaboratorToDocument(c.Request.Context(), body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, httpResponseMessage{
		Message: "document permission created",
	})
}

func (h *HTTPHandler) createDocument(c *gin.Context) {
	// Get userID from middleware context
	ownerID := c.GetString("userID")

	var body CreateDocumentDTO
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "bad request: " + err.Error(),
		})
		return
	}

	body.OwnerID = ownerID
	document, err := h.documentService.CreateNewDocument(c.Request.Context(), body)
	if err != nil {
		// Log the error for debugging purposes
		c.Error(err)
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to create document: " + err.Error(),
		})
		return
	}

	response := ToDocumentResponse(document)

	c.JSON(http.StatusCreated, response)
}

func (h *HTTPHandler) getDocuments(c *gin.Context) {
	// Get userID from middleware context
	userID := c.GetString("userID")

	documents, err := h.documentService.GetUserDocuments(c.Request.Context(), userID, false)
	if err != nil {
		c.JSON(http.StatusBadRequest, httpResponseMessage{
			Message: "failed to fetch documents",
		})
		return
	}

	// Convert to DTO response
	response := ToDocumentResponseList(documents)
	c.JSON(http.StatusOK, response)
}

func (h *HTTPHandler) getOneDocument(c *gin.Context) {
	// Get document from middleware context (already validated)
	document, exists := c.Get("document")
	if !exists {
		c.JSON(http.StatusInternalServerError, httpResponseMessage{
			Message: "document not found in context",
		})
		return
	}

	// Type assert and convert to DTO response
	doc, ok := document.(*Document)
	if !ok {
		c.JSON(http.StatusInternalServerError, httpResponseMessage{
			Message: "invalid document type in context",
		})
		return
	}

	response := ToDocumentDetailResponse(doc)
	c.JSON(http.StatusOK, response)
}

type httpResponseMessage struct {
	Message string `json:"message"`
}
