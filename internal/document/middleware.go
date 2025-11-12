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

// DocumentAccessMiddleware - unified middleware to validate access on-demand
// requiredPermission can be: "owner", "editor", "viewer"
func DocumentAccessMiddleware(service *DocumentService, requiredPermission string) gin.HandlerFunc {
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

		// Verify cached permissions
		if permissions, exists := c.Get("userPermissions"); exists {
			permMap := permissions.(map[string]string)
			if permission, hasAccess := permMap[documentID]; hasAccess {
				if validatePermission(permission, requiredPermission) {
					c.Set("documentID", documentID)
					c.Set("userPermission", permission)
					c.Next()
					return
				}
			}
		}

		// Specific query for this document
		document, permission := service.GetDocumentWithPermission(c.Request.Context(), userID.(string), documentID)

		if document == nil || !validatePermission(permission, requiredPermission) {
			c.JSON(http.StatusNotFound, gin.H{
				"message": "document not found or access denied",
			})
			c.Abort()
			return
		}

		c.Set("document", document)
		c.Set("documentID", documentID)
		c.Set("userPermission", permission)
		c.Next()
	}
}

// validatePermission verify if the user has the required permission level
func validatePermission(userPermission, required string) bool {
	// owner > editor > viewer
	permissionLevels := map[string]int{
		"owner":  3,
		"editor": 2,
		"viewer": 1,
	}

	userLevel, userExists := permissionLevels[userPermission]
	requiredLevel, requiredExists := permissionLevels[required]

	if !userExists || !requiredExists {
		return false
	}

	return userLevel >= requiredLevel
}

// Convenience functions for common permission checks
func RequireOwnerAccess(service *DocumentService) gin.HandlerFunc {
	return DocumentAccessMiddleware(service, "owner")
}

func RequireEditorAccess(service *DocumentService) gin.HandlerFunc {
	return DocumentAccessMiddleware(service, "editor")
}

func RequireViewerAccess(service *DocumentService) gin.HandlerFunc {
	return DocumentAccessMiddleware(service, "viewer")
}
