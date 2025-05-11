package handler

import "github.com/gin-gonic/gin"

// Route Animal.
func RegisterAnimalRoutes(routerGroup *gin.RouterGroup, h *AnimalHandler) {
	animalRoutes := routerGroup.Group("/animal")
	{
		animalRoutes.POST("", h.CreateAnimal)       // ID dari klien
		animalRoutes.PUT("", h.UpsertAnimal)        // Upsert (ID di payload)
		animalRoutes.GET("", h.GetAllAnimals)
		animalRoutes.GET("/:id", h.GetAnimalByID)   // ID dari path
		animalRoutes.DELETE("/:id", h.DeleteAnimal) // ID dari path
	}
}