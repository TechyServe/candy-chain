var cnd = require('../candies/candy.js');

exports.queryallcandies = function(req,res)
{
    cnd.queryallcandies().then((err,data) => {
        if(err)
        res.send(err);
        else
        res.send(data);
    });
   
   
};

exports.findCandy = function(req,res)
{
    cnd.findcandy(req.params.candy).then((err,data) => {
        if(err)
        res.send(err);
        else
        res.json(data);
    });
   
}