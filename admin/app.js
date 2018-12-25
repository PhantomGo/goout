const app = require('tcb-admin-node');
app.init({
    secretId: 'AKIDJg1fXCAOqBvaHu66hMBJqLVmzMarrcEu',
    secretKey: 'nnFFE6sYDlqba6VVAyytVcHulS2BQfaN',
    env: 'wx-f1b839'
});
const db = app.database();

const express = require('express')
const server = express()

server.get('/levels', async(req, res) => {
    const item = await db.collection('userlevels').field({
        level: true,
        exp: true
      }).get()
    res.json(item.data)
})

server.listen(3000, () => console.log('Example app listening on port 3000!'))


