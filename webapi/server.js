// require our dependencies
var express        = require('express');
var app            = express();
var port           = process.env.PORT || 8080;
var bodyParser = require('body-parser');
var router = require('./app/router');



// set static files (css and images, etc) location
app.use(express.static(__dirname + '/public'));
router(app);

// parse requests of content-type - application/json
app.use(bodyParser.json())
// start the server

// app.get('/', function(req, res){
//   res.json({"message": "Welcome to CandyChain"});
// });


app.listen(port, function() {
  console.log('app started');
});