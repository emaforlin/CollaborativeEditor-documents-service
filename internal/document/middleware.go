package internal

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// UserHeaderMiddleware validates the X-User-Id header and stores it in the context
func UserHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.GetHeader("X-User-Id")
		if userID == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "missing or invalid user ID header",
			})
			c.Abort()
			return
		}

		// Store userID in context for use in handlers
		c.Set("userID", userID)
		c.Next()
	}
}

// DocumentOwnershipMiddleware validates that the user owns the document
// This middleware depends on UserHeaderMiddleware being called first
func DocumentOwnershipMiddleware(service *DocumentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "user ID not found in context",
			})
			c.Abort()
			return
		}

		documentID := c.Param("id")
		if documentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "document ID is required",
			})
			c.Abort()
			return
		}

		// Check if document exists and user is the owner
		foundDoc := service.GetOneDocument(c.Request.Context(), GetOneDocumentDTO{
			DocumentID: documentID,
			OwnerID:    userID.(string),
		})

		if foundDoc == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "document not found or access denied",
			})
			c.Abort()
			return
		}

		// Store document and documentID in context for handlers
		c.Set("document", foundDoc)
		c.Set("documentID", documentID)
		c.Next()
	}
}

// CollaboratorAccessMiddleware validates that the user has access to the document (owner or collaborator)
// This is more permissive than DocumentOwnershipMiddleware
func CollaboratorAccessMiddleware(service *DocumentService) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "user ID not found in context",
			})
			c.Abort()
			return
		}

		documentID := c.Param("id")
		if documentID == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "document ID is required",
			})
			c.Abort()
			return
		}

		// TODO: Implement logic to check if user has access as owner or collaborator
		// For now, this uses the same logic as DocumentOwnershipMiddleware
		// You'll need to extend this when you implement the collaborator checking logic
		foundDoc := service.GetOneDocument(c.Request.Context(), GetOneDocumentDTO{
			DocumentID: documentID,
			OwnerID:    userID.(string),
		})

		if foundDoc == nil {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "document not found or access denied",
			})
			c.Abort()
			return
		}

		c.Set("document", foundDoc)
		c.Set("documentID", documentID)
		c.Next()
	}
}
