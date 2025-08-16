package http

import (
	"fmt"
	"strconv"
	"strings"
)

type Method string

const (
	MethodGet     Method = "GET"
	MethodHead    Method = "HEAD"
	MethodPost    Method = "POST"
	MethodPut     Method = "PUT"
	MethodPatch   Method = "PATCH"
	MethodDelete  Method = "DELETE"
	MethodConnect Method = "CONNECT"
	MethodOptions Method = "OPTIONS"
	MethodTrace   Method = "TRACE"
)

func IsMethodValid(m string) (Method, error) {
	switch m {
	case "GET":
		return MethodGet, nil
	case "HEAD":
		return MethodHead, nil
	case "POST":
		return MethodPost, nil
	case "PUT":
		return MethodPut, nil
	case "PATCH":
		return MethodPatch, nil
	case "DELETE":
		return MethodDelete, nil
	case "CONNECT":
		return MethodConnect, nil
	case "OPTIONS":
		return MethodOptions, nil
	case "TRACE":
		return MethodTrace, nil
	default:
		return "", fmt.Errorf("unsupported method %q", m)
	}
}

func (m Method) String() string {
	return string(m)
}

func handleHeadMethod(req Request, res Response, resource map[Method]Handler) {
	handler, handlerOk := resource[MethodGet]
	if !handlerOk {
		HTTPError(res, StatusNotFound)
		return
	}
	handler.ServeHTTP(req, res)
}

func handlerOptionsMethod(req Request, res Response, router *DefaultRouter, resource map[Method]Handler) {
	if req.Path() == "*" {
		handleGeneralOptionsMethod(res, router)
		return
	}

	handleSpecificOptionsMethod(res, router, resource)
}

func handleGeneralOptionsMethod(res Response, router *DefaultRouter) {
	methods := router.getAllAvailableMethods()
	mtdStrArr := make([]string, len(methods))
	for i, m := range methods {
		mtdStrArr[i] = m.String()
	}
	corsMtdStrArr := make([]string, len(router.cors.AllowedMethods))
	for i, m := range router.cors.AllowedMethods {
		corsMtdStrArr[i] = m.String()
	}

	res.WriteHeader(StatusOk)
	res.SetHeader("Allow", strings.Join(mtdStrArr, ", "))
	res.SetHeader("Access-Control-Allow-Origin", strings.Join(router.cors.AllowedOrigins, ", "))
	res.SetHeader("Access-Control-Allow-Methods", strings.Join(corsMtdStrArr, ", "))
	res.SetHeader("Access-Control-Allow-Headers", strings.Join(router.cors.AllowedHeaders, ", "))
	res.SetHeader("Access-Control-Allow-Credentials", strconv.FormatBool(router.cors.AllowedCredentials))
	res.SetHeader("Access-Control-Max-Age", strconv.Itoa(router.cors.AllowedMaxAge))
}

func handleSpecificOptionsMethod(res Response, router *DefaultRouter, resource map[Method]Handler) {
	methods := make([]Method, 0)
	for k := range resource {
		methods = append(methods, k)
	}
	mtdStrArr := make([]string, len(methods))
	for i, m := range methods {
		mtdStrArr[i] = m.String()
	}
	corsMtdStrArr := make([]string, len(router.cors.AllowedMethods))
	for i, m := range router.cors.AllowedMethods {
		corsMtdStrArr[i] = m.String()
	}

	res.WriteHeader(StatusOk)
	res.SetHeader("Allow", strings.Join(mtdStrArr, ", "))
	res.SetHeader("Access-Control-Allow-Origin", strings.Join(router.cors.AllowedOrigins, ", "))
	res.SetHeader("Access-Control-Allow-Methods", strings.Join(corsMtdStrArr, ", "))
	res.SetHeader("Access-Control-Allow-Headers", strings.Join(router.cors.AllowedHeaders, ", "))
	res.SetHeader("Access-Control-Allow-Credentials", strconv.FormatBool(router.cors.AllowedCredentials))
	res.SetHeader("Access-Control-Max-Age", strconv.Itoa(router.cors.AllowedMaxAge))
}
