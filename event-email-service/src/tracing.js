const traceHeaders = [
  'Traceparent',
  'Tracestate',
  'Baggage',
  'X-Request-Id',
  'X-B3-Traceid',
  'X-B3-Spanid',
  'X-B3-Parentspanid',
  'X-B3-Sampled',
  'X-B3-Flags',
  'X-Ot-Span-Context'
];

module.exports.propagateTracingHeaders = function(headers, downstreamReq) {
  for (var h in traceHeaders) {
    let headerVal = headers[h];
    if (headerVal !== undefined) {
      downstreamReq.headers[h] = headerVal;
    }
  }
};
