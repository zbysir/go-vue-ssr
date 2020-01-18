const Koa = require('koa');
const bodyParser = require('koa-bodyparser')
const {Parser} = require("acorn")

const app = new Koa();
app.use(bodyParser())

app.use(async ctx => {
  let code = ctx.request.body.code

  if (!code) {
    code = ctx.request.query.code
  }
  if (!code) {
    ctx.response.body = "null";
    return
  }

  try {
    let ast = Parser.parse(code)
    ctx.response.body = ast
    ctx.response.status = 200
  } catch (e) {
    ctx.response.body = {code: 400, err: e}
    ctx.response.status = 400
  }
});

app.listen(3000);