// 中间件函数定义
package middle

import "github.com/kataras/iris/v12"

/*
	如果是NGINX转发，Nginx 从1.7.5开始，add_header支持“always”参数 ，这样后端返回4xx或5xx状态代码，则允许CORS工作。参考资料：https://gist.github.com/Stanback/7145487
	nginx配置401、404、500状态码跨域问题
	在做项目权限的时候，后端那边做接口的权限，做完后出现的问题是，接口正常跑返回200的时候没有问题，但是如果接口报错的时候，并没有正确的提示错误信息，而是提示了跨域，最后发现是运维nginx配置的问题。
	参考 http://nginx.org/en/docs/http/ngx_http_headers_module.html
	如果响应代码等于 200,201,204,206,301,302,303,304,307,308，时nginx会则将指定字段添加到响应标头，
	而 401、404、500，都没有添加header，如果指定了always参数，则无论响应代码如何，都将添加头字段。
	由于 Access-Control-Allow-Origin * 选项未生效,导致返回的response header里面没有允许跨域请求的选项。
	只要在Access-Control-Allow-Origin *后加上always，即解决
	Access-Control-Allow-Origin *;  后面添加always，变为：Access-Control-Allow-Origin * always; 如：
		if ($request_method = 'OPTIONS') {
		add_header Access-Control-Allow-Origin '$http_origin' always;
		add_header Access-Control-Allow-Credentials 'true' always;
		add_header Access-Control-Allow-Methods 'GET,POST,PUT,DELETE,OPTIONS' always;
		add_header Access-Control-Allow-Headers 'Authorization,DNT,User-Agent,Keep-Alive,Content-Type,accept,origin,X-Requested-With' always;
		add_header Access-Control-Max-Age 3600;
		add_header Content-Length 0;
		return 200;
		}
*/

// 于处理 CORS 的中间件
func MiddlewareCORS(ctx iris.Context) {
	// iris服务端自身实现跨域，可以通过Iris cors 中间件实现，也可以自己代码实现。
	// 当然也可以通过nginx作为反向代理服务器解决跨域问题
	ctx.Header("Access-Control-Allow-Origin", "*")         // 给所有请求添加头信息Access-Control-Allow-Origin * 表明它允许所有域发起跨域请求
	ctx.Header("Access-Control-Allow-Credentials", "true") // 允许携带 用户认证凭据（也就是允许客户端发送的请求携带Cookie）

	// 如果是OPTIONS方法的预请求  ctx.Request().Method == "OPTIONS"
	if ctx.Method() == iris.MethodOptions {
		// 允许哪些方法的跨域请求
		ctx.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")

		// 允许跨域请求包含content-type头
		ctx.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Origin, Content-Type, Accept, Authorization")

		// 用于CORS预检请求的缓存时间 - 不需要再发送预检验请求 （单位为秒） 不指定时即使用默认值，Chrome默认5秒。
		// 常用浏览器有不同的最大值限制，Firefox上限是24小时 （即86400秒），Chrom是10分钟（即600秒）。
		ctx.Header("Access-Control-Max-Age", "86400")

		// CORS caching for CDNs - To cache CORS responses in CDNs and other proxies between the browser and your API server, add:
		// Cache-Control 用于广泛的一般上下文,以指定资源被视为新鲜的最长时间.
		ctx.Header("Cache-Control", "public, max-age=86400")

		// Always set Vary headers.
		// Vary头应该在设置304 Not Modified完全一样会被设定在相当的反应200 OK响应。
		ctx.Header("Vary", "Origin, Access-Control-Request-Method, Access-Control-Request-Headers")

		// 成功响应并添加对应的头信息，状态码响应204或者200值 ctx.StatusCode(204)
		ctx.StatusCode(iris.StatusNoContent)

		// 停止继续执行项目代码
		return
	}

	ctx.Next()
}
