package main

import (
  "os"
  "fmt"
  "bytes"
  "archive/tar"
  "io"
  "io/ioutil"
  "log"
  "context"
  "github.com/gin-gonic/gin"
  "github.com/Shigoto-Q/docker_service/controller"
  "github.com/Shigoto-Q/docker_service/service"
  "github.com/docker/docker/client"
  "github.com/docker/docker/api/types"
  "github.com/go-git/go-git/v5"
)

var (
  dockerService service.DockerService = service.New()
  dockerController controller.DockerController = controller.New(dockerService)
  dockerRegistryUserID = "registry"
)


func imageBuild(cli *client.Client) {
    ctx := context.Background()
    buf := new(bytes.Buffer)
    tw := tar.NewWriter(buf)
    defer tw.Close()

    dockerFile := "myDockerfile"
    dockerFileReader, err := os.Open("/home/shins/shigoto/docker_service/docker_deploy_test/Dockerfile")
    if err != nil {
        log.Fatal(err, " :unable to open Dockerfile")
    }
    readDockerFile, err := ioutil.ReadAll(dockerFileReader)
    if err != nil {
        log.Fatal(err, " :unable to read dockerfile")
    }

    tarHeader := &tar.Header{
        Name: dockerFile,
        Size: int64(len(readDockerFile)),
    }
    err = tw.WriteHeader(tarHeader)
    if err != nil {
        log.Fatal(err, " :unable to write tar header")
    }
    _, err = tw.Write(readDockerFile)
    if err != nil {
        log.Fatal(err, " :unable to write tar body")
    }
    dockerFileTarReader := bytes.NewReader(buf.Bytes())

    imageBuildResponse, err := cli.ImageBuild(
        ctx,
        dockerFileTarReader,
        types.ImageBuildOptions{
            Context:    dockerFileTarReader,
            Dockerfile: dockerFile,
            Remove:     true})
    if err != nil {
        log.Fatal(err, " :unable to build docker image")
    }
    defer imageBuildResponse.Body.Close()
    _, err = io.Copy(os.Stdout, imageBuildResponse.Body)
    if err != nil {
        log.Fatal(err, " :unable to read image build response")
    }
}


func main() {
  cli, err := client.NewClientWithOpts(client.FromEnv)
  if err != nil {
    panic(err)
  }

  if err != nil {
    fmt.Println(err.Error())
  }
  repo, err := git.PlainClone("/home/shins/shigoto/docker_service/docker_deploy_test/", false, &git.CloneOptions{
    URL: "https://github.com/Shigoto-Q/docker_deploy_test",
    Progress: os.Stdout,
  })
  fmt.Println(&repo)
  imageBuild(cli)
  if err != nil {
    panic(err)
  }

  server := gin.Default()

  server.POST("/docker", func(ctx *gin.Context) {
    ctx.JSON(200, dockerController.Save(ctx))
  })

  server.Run(":5050")
}
