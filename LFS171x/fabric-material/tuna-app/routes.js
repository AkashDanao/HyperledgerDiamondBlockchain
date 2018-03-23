//SPDX-License-Identifier: Apache-2.0

var tuna = require('./controller.js');
var reqBodyParser    = require('body-parser');

module.exports = function(app){

  app.use(reqBodyParser.json({limit:'50mb'}));
  app.use(reqBodyParser.urlencoded({extended:true, limit:'50mb'})); 

  app.get('/get_tuna/:id', function(req, res){
    tuna.get_tuna(req, res);
  });
  app.post('/add_tuna', function(req, res){
    tuna.add_tuna(req, res);
  });
  app.get('/get_all_tuna', function(req, res){
    tuna.get_all_tuna(req, res);
  });
  app.get('/change_holder/:holder', function(req, res){
    tuna.change_holder(req, res);
  });
  app.get('/update_location/:location',function(req,res){
    tuna.update_location(req,res);
  });
}
