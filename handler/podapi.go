package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	pod "github.com/DuanNengxin/wepass-pod/proto"
	podapi "github.com/DuanNengxin/wepass-podapi/proto"
	"strconv"
)

type PodApi struct {
	PodService pod.PodService
}

// podApi.FindPodById 通过API向外暴露为/podApi/findPodById，接收http请求
// 即：/podApi/FindPodById 请求会调用go.micro.api.podApi 服务的podApi.FindPodById 方法
func (p *PodApi) GetPodById(ctx context.Context, req *podapi.Request, resp *podapi.Response) error {
	fmt.Println("接收到 podApi.FindPodById 的请求")
	if _, ok := req.GetGet()["pod_id"]; !ok {
		resp.StatusCode = 500
		return errors.New("参数异常")
	}
	podIdString := req.GetGet()["pod_id"].Values[0]
	podId, err := strconv.ParseInt(podIdString, 10, 64)
	if err != nil {
		return err
	}
	podInfo, err := p.PodService.FindPodByID(ctx, &pod.PodID{
		Id: podId,
	})
	if err != nil {
		return err
	}
	b, _ := json.Marshal(podInfo)
	resp.StatusCode = 200
	resp.Body = string(b)
	return nil
}

func (p *PodApi) AddPod(ctx context.Context, req *podapi.Request, resp *podapi.Response) error {
	fmt.Println("接收到 podApi.AddPod 的请求")
	return nil
}

func (p *PodApi) UpdatePod(ctx context.Context, req *podapi.Request, resp *podapi.Response) error {
	fmt.Println("接收到 podApi.UpdatePod 的请求")
	return nil
}

func (p *PodApi) DeletePodById(ctx context.Context, req *podapi.Request, resp *podapi.Response) error {
	fmt.Println("接收到 podApi.DeletePodById 的请求")
	return nil
}

func (p *PodApi) Call(ctx context.Context, req *podapi.Request, resp *podapi.Response) error {
	fmt.Println("接收到 podApi.DeletePodById 的请求")
	return nil
}
