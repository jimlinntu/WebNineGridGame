你好！
很冒昧用这样的方式来和你沟通，如有打扰请忽略我的提交哈。我是光年实验室（gnlab.com）的HR，在招Golang开发工程师，我们是一个技术型团队，技术氛围非常好。全职和兼职都可以，不过最好是全职，工作地点杭州。
我们公司是做流量增长的，Golang负责开发SAAS平台的应用，我们做的很多应用是全新的，工作非常有挑战也很有意思，是国内很多大厂的顾问。
如果有兴趣的话加我微信：13515810775  ，也可以访问 https://gnlab.com/，联系客服转发给HR。
# WebNineGridGame

## Setup
* Compiled Vue (I build it by `@vue/cli 4.2.3`) by `npm run build` and move the compiled folder into `./dist`.
* Run your mongo db server at `mongodb://172.17.0.1:17990` (ex. `docker run --rm -p 172.17.0.1:17990:27017 mongo`)
* Setup the Golang docker image by: `docker build -t jimlin7777/webninegrid .`
* Run the Golang docker image by: `docker run -it --rm -p 17989:80 -v $(pwd):/root/WebNineGridGame jimlin7777/webninegrid bash`

## System Logical View
```
Vue (./game -> ./dist) <-> Golang (main.go) <-> MongoDB
```

## Demo

* About
![about](./demo/about.png)

* Login
![login](./demo/login.png)

* Login Success
![login\_success](./demo/login_success.png)

* User Drag-and-drop Interative Interface
![dnd](./demo/drag-n-drop.png)

* User Question Selection
![user_selection](./demo/question-selection.png)

* User Question Answering
![user_answering](./demo/question-answering.png)

* Admin Management Page
![admin](./demo/admin_page.png)

* Admin Approval or Rejection Page
![admin_approval](./demo/admin_page_2.png)

* Admin Teams' Status Page
![admin_team_status](./demo/admin_page_3.png)

* User Question Passed
![user_passed](./demo/user-answer-passed.png)
