
// require express
var express = require('express');
var path    = require('path');

// create our router object
var router = express.Router();


// export our router
module.exports = function(app)
{
  var fca = require('./controller/fabric-ca-controller.js');
  var candy = require('./controller/candies-controller.js');
  app.get('/',function(req,res){
    res.json({"message": "Welcome to CandyChain"});
  });
  app.get('/Admin/',fca.enrollAdmin);
  app.get('/registeruser/:id',fca.enrollandregisteruser);
  app.get('/candies',candy.queryallcandies);
  app.get('/candies/:candy',candy.findCandy)


};
