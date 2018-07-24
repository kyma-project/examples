const express = require('express');
const bodyParser = require('body-parser');
//const tracing = require('./tracing.js');
const app = express();
const port = 3000;

app.use(bodyParser.urlencoded({ extended: false }));
app.use(bodyParser.json());

// Listen to events and send an email
app.post('/v1/events/register', (req, res) => {
  let uid = getCustomerID(req.body);
  let email = getCustomerEmail(req.body);

  if (uid === undefined || email === undefined) {
    console.log('No customer ID or Email received!');
    res.sendStatus(400);
  } else {
    console.log('Customer created with ID: ' + uid + ' and Email: ' + email);
    res.sendStatus(200);
  }

  /*
        TODO:
         - propagate tracing headers to email service: tracing.propagateTracingHeaders(req.headers, req)
         - send email to event.customer.uid
    */
});

var server = app.listen(port, () =>
  console.log('Example app listening on port ' + port + '!')
);

app.stop = function() {
  server.close();
};

module.exports = app;

function getCustomerID(body) {
  if (body.event === undefined || body.event.customer === undefined) {
    return undefined;
  }
  return body.event.customer.customerID;
}

function getCustomerEmail(body) {
  if (body.event === undefined || body.event.customer === undefined) {
    return undefined;
  }
  return body.event.customer.uid;
}
