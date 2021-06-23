package controller

import (
  "os"
  "io"
  "bufio"
  "encoding/json"
  "fmt"
  "errors"
  "log"
  "time"
  "encoding/base64"
  "strings"
  "context"
  "github.com/go-git/go-git/v5"
  "github.com/mitchellh/go-homedir"
  "github.com/docker/docker/api/types"
  "github.com/docker/docker/client"
  "github.com/docker/docker/pkg/archive"
  "github.com/joho/godotenv"
  "github.com/gin-gonic/gin"
  "github.com/Shigoto-Q/docker_service/service"
  "github.com/Shigoto-Q/docker_service/entity"

)
var dockerRegistryUserID = "shigoto"

type DockerController interface {
  Save(ctx *gin.Context, cli *client.Client) entity.DockerImage
}


type ErrorLine struct {
    Error       string      `json:"error"`
    ErrorDetail ErrorDetail `json:"errorDetail"`
}

type ErrorDetail struct {
    Message string `json:"message"`
}


func cloneRepo(url string, name string) error {
    _, err := git.PlainClone(name, false, &git.CloneOptions{
    URL: url,
    Progress: os.Stdout,
  })
  if err != nil {
    return err
  }
  return nil
}

func dprint(rd io.Reader) error {
    var lastLine string

    scanner := bufio.NewScanner(rd)
    for scanner.Scan() {
        lastLine = scanner.Text()
        fmt.Println(scanner.Text())
    }

    errLine := &ErrorLine{}
    json.Unmarshal([]byte(lastLine), errLine)
    if errLine.Error != "" {
        return errors.New(errLine.Error)
    }

    if err := scanner.Err(); err != nil {
        return err
  }

  return nil
}

func imagePush(cli *client.Client, imageName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
    defer cancel()
    err := godotenv.Load(".env")
    if err != nil {
      fmt.Println(err)
    }
    user := os.Getenv("DOCKERHUBUSER")
    token := os.Getenv("DOCKERHUBTOKEN")

    var authConfig = types.AuthConfig{
      Username: user,
      Password: token, 
      ServerAddress: "https://index.docker.io/v1/",
    }
    authConfigBytes, _ := json.Marshal(authConfig)
    authConfigEncoded := base64.URLEncoding.EncodeToString(authConfigBytes)
    tag := dockerRegistryUserID + imageName
    opts := types.ImagePushOptions{RegistryAuth: authConfigEncoded}
    rd, err := cli.ImagePush(ctx, tag, opts)
    if err != nil {
      return err
    }
  
    defer rd.Close()
    err = dprint(rd)
    if err != nil {
      return err
    }
    return nil
}

func GetContext(filePath string) io.Reader {
    filePath, err := homedir.Expand(filePath)
    if err != nil {
      log.Println(err)
    }
    ctx, _ := archive.TarWithOptions(filePath, &archive.TarOptions{})
    return ctx
}



func imageBuild(cli *client.Client, filePath string, imageName string) error {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*120)
    defer cancel()

    opts := types.ImageBuildOptions{
        Dockerfile: "Dockerfile",
        Tags:       []string{dockerRegistryUserID + imageName},
        Remove:     true,
    }
    res, err := cli.ImageBuild(ctx, GetContext(filePath), opts)
    if err != nil {
        return err
    }
    defer res.Body.Close()
    err = dprint(res.Body)
    if err != nil {
        return err
    }
  return nil
}
type controller struct {
  service service.DockerService
}

func New(service service.DockerService) DockerController {
  return &controller{
    service: service,
  }
}
func (c *controller) Save(ctx *gin.Context, cli *client.Client) entity.DockerImage {
  var docker entity.DockerImage

  ctx.BindJSON(&docker)
  c.service.Save(docker)

  name := docker.FullName
  url := docker.RepoUrl
  err := cloneRepo(url, "/tmp/" + name)
  if err != nil {
      log.Fatal(err)
  }
  imageNameS := strings.Split(name, "/")
  imageName :=  "/" + imageNameS[len(imageNameS)-1]
  err = imageBuild(cli, "/tmp/" + name, imageName) 
  if err != nil {
    fmt.Println(err.Error())
  }
  
  err = imagePush(cli, imageName)
  if err != nil {
    fmt.Println(err.Error())
  }
  docker.ImageName = dockerRegistryUserID + imageName
  err = os.RemoveAll("/tmp/" + name)
  if err != nil {
      fmt.Println(err.Error())
  }

  return docker

}
