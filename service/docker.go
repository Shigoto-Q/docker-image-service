// Package j provides ...
package service

import entity "github.com/Shigoto-Q/docker_service/entity"

type DockerService interface {
  Save(entity.DockerImage) entity.DockerImage
}


type dockerService struct {
  docker [] entity.DockerImage
}


func New()  DockerService {
  return &dockerService{}
}

func (service *dockerService) Save(docker entity.DockerImage) entity.DockerImage{
  service.docker = append(service.docker, docker)
  return docker
}
