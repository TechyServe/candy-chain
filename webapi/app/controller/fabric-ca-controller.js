var ca = require('../fabric-ca/enroll.js');

exports.enrollAdmin = function(req,res)
{
   res.send(ca.enrollAdmin());
};

exports.enrollandregisteruser = function(req,res)
{
    res.send(ca.enrollandregisteruser(req.params.id))
};