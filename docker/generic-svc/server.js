const express = require('express')
const fs = require("fs")
const app = express()
const port = process.env.PORT || 4001
// respond with "hello world" when a GET request is made to the homepage
const db_file = process.env.DATA_FILE
if(!fs.existsSync(db_file)){
  console.error("DB File is not exist")
  process.exit(1)
}
const db = require(db_file)
var logger = (req, res, next) => {
  next()
  console.log([req.method, req.url, res.statusCode].join(" - "))
}

app.use(logger)
app.use(express.json());


function queryMatcher(query){
  return (data) => {
    for(const key in query) {
      const qval = query[key]
      const dval = data[key]

      if(typeof qval == 'string' && qval != dval)
        return false;
      
      if(Array.isArray(qval) && qval.indexOf(dval) < 0 )
        return false;
    }
    return true
  }
}

app.get('/:id', (req, res) => {
  req.setTimeout(200)
  const {id} = req.params
  const data = db.data.find( (e) => e.id == id)
  if(data == undefined){
    res.status(404)
    return res.send({data: null})
  }
  res.send({data})
})

app.put('/:id', (req, res) => {
  res.send({data: req.body})
})

app.get('/', (req, res) => {
  res.send({ data: db.data.filter( queryMatcher(req.query) )})
})

app.post('/', (req, res) => {
  res.send({data: req.body})
})

app.listen(port, () => {
  console.log(`Example app listening on port ${port}`)
})

