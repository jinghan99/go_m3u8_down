import service from 'axios' // 引入axios

const api = service.create({
  baseURL: "http://127.0.0.1:8080",
  timeout: 2000                   //请求超时设置，单位ms
})
// http request 拦截器
service.interceptors.request.use(

)

// http response 拦截器
service.interceptors.response.use(
  response => {
    if (response.data.code === 0 || response.headers.success === 'true') {
      if (response.headers.msg) {
        response.data.msg = decodeURI(response.headers.msg)
      }
      return response.data
    } else {
      return response.data.msg ? response.data : response
    }
  },
  error => {
    if (!error.response) {
      return
    }

    switch (error.response.status) {
      case 500:
        break
      case 404:
        break
    }

    return error
  }
)
export default api
