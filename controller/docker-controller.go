package controller

import (
  "github.com/gin-gonic/gin"
  service "github.com/Shigoto-Q/docker_service/service"
  entity "github.com/Shigoto-Q/docker_service/entity"

)
type DockerController interface {
  Save(ctx *gin.Context) entity.DockerImage
}


type controller struct {
  service service.DockerService
}

func New(service service.DockerService) {
}
func (c *controller) Save(ctx *gin.Context) entity.DockerImage {
  var docker entity.DockerImage
  return docker

}
