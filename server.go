package main

import (
	"fmt"
	"log"

	"github.com/Shigoto-Q/docker_service/controller"
	"github.com/Shigoto-Q/docker_service/service"
	"github.com/docker/docker/client"
	"github.com/gin-gonic/gin"
)

var (
	dockerService        service.DockerService       = service.New()
	dockerController     controller.DockerController = controller.New(dockerService)
	dockerRegistryUserID                             = ""
)

func main() {
	log.SetFlags(log.Lshortfile)
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server := gin.Default()
	server.POST("/docker", func(ctx *gin.Context) {
		ctx.JSON(200, dockerController.Save(ctx, cli))
	})
	server.Run(":5050")
}
