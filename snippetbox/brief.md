1. Đây là project dùng tech: golang, mysql, pat, middleware, testing, logger, wrapper errors.
2. Mọi thứ điều được khai báo ở main.go bằng application structure, models cũng phải khai báo ở main.
3. Có render template
4. Dùng session để lưu authentication, và check info authen mỗi khi user login vào web.
5. Flow. main -> router -> middleware -> handler -> validation input ->(models) send data to function get data from database or from api 3rd -> return data into handler -> return response api or render data to template.
