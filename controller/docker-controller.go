package controller

import (
  "github.com/gin-gonic/gin"
  "github.com/Shigoto-Q/docker_service/service"
  "github.com/Shigoto-Q/docker_service/entity"

)
type DockerController interface {
  Save(ctx *gin.Context) entity.DockerImage
}


type controller struct {
  service service.DockerService
}

func New(service service.DockerService) DockerController {
  return &controller{
    service: service,
  }
}
func (c *controller) Save(ctx *gin.Context) entity.DockerImage {
  var docker entity.DockerImage
  ctx.BindJSON(&docker)
  c.service.Save(docker)
  return docker

}
